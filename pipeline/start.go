package pipeline


import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	in := make(chan interface{})
	firstCh := in // Сохраняем начало трубы

	wg := sync.WaitGroup{}
	for _, c := range cmds {
		out := make(chan interface{})
		wg.Add(1)

		// Передаем каналы и команду как аргументы, чтобы избежать гонки в RunPipeline
		go func(entryCmd cmd, currIn, currOut chan interface{}) {
		defer wg.Done()
			// ВАЖНО: RunPipeline закрывает канал, когда воркер закончил работу
			defer close(currOut)
			entryCmd(currIn, currOut)
		}(c, in, out)

		in = out
	}

	close(firstCh) // Закрываем вход, чтобы запустить первую функцию
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	// Можно использовать map[uint64]bool, это быстрее и надежнее
	lay := make(map[uint64]bool)

	for user := range in {
		wg.Add(1)
		go func(u interface{}) {
			defer wg.Done()
			us := GetUser(u.(string))

			mu.Lock()
			// Используем 'Check and Set' паттерн
			if _, exists := lay[us.ID]; !exists {
				lay[us.ID] = true
				mu.Unlock()
				out <- us
			} else {
				mu.Unlock()
			}
		}(user)
	}

	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	buff := []User{}

	getmessagefunc := func(usersBatch []User) {
		defer wg.Done()
		msgs, err := GetMessages(usersBatch...)
		if err == nil {
			for _, msg := range msgs {
				out <- msg
			}
		}
	}

	for user := range in {
		buff = append(buff, user.(User))
		if len(buff) == GetMessagesMaxUsersBatch {
			wg.Add(1)
			// Важно: создаем новый слайс, чтобы отвязаться от старого массива
			batchCopy := make([]User, len(buff))
			copy(batchCopy, buff)

			go getmessagefunc(batchCopy)
			buff = []User{} // Сбрасываем буфер
		}
	}

	if len(buff) != 0 {
		wg.Add(1)
		go getmessagefunc(buff)
	}

	wg.Wait()
	// УБРАЛИ close(out)
}

func CheckSpam(in, out chan interface{}) {
	sem := make(chan struct{}, HasSpamMaxAsyncRequests)
	wg := sync.WaitGroup{}

	for id := range in {
		sem <- struct{}{}
		wg.Add(1)
		// ИСПРАВЛЕНИЕ RACE: передаем id как аргумент (val)
		go func(val interface{}) {
			defer func() { <-sem }()
			defer wg.Done()

			msgID := val.(MsgID)
			isSpam, err := HasSpam(msgID)
			if err != nil {
				return
			}
			out <- MsgData{
				ID:      msgID,
				HasSpam: isSpam,
			}
		}(id)
	}

	wg.Wait()
	// УБРАЛИ close(out)
}

func CombineResults(in, out chan interface{}) {
	sl := []MsgData{}

	for msg := range in {
		if data, ok := msg.(MsgData); ok {
			sl = append(sl, data)
		}
	}

	sort.Slice(sl, func(i, j int) bool {
		// Сначала true (спам), потом false
		if sl[i].HasSpam != sl[j].HasSpam {
			return sl[i].HasSpam && !sl[j].HasSpam
		}
		// Внутри групп по возрастанию ID
		return sl[i].ID < sl[j].ID
	})

	for _, data := range sl {
		out <- fmt.Sprintf("%t %d", data.HasSpam, data.ID)
	}
	// УБРАЛИ close(out), хотя здесь это не критично, но для порядка лучше убрать
}
