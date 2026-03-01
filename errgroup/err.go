package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func fetchPriceFromAPI(ctx context.Context, ID int) (float64, error) {
	if ID == 3 {
		return 0, fmt.Errorf("ID is 3")
	}
	select {
	case <-time.After(time.Second * 2):
		return rand.Float64() * float64(ID), nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}

}

// FetchPrices - ЭТУ ФУНКЦИЮ НУЖНО РЕАЛИЗОВАТЬ
func FetchPrices(ctx context.Context, itemIDs []int) (map[int]float64, error) {
	g, ctx := errgroup.WithContext(ctx)
	mu := sync.Mutex{}

	prices := make(map[int]float64)
	g.SetLimit(5)
	for _, ID := range itemIDs {
		g.Go(func() error {
			price, err := fetchPriceFromAPI(ctx, ID)
			if err == nil {
				mu.Lock()
				prices[ID] = price
				mu.Unlock()
				return nil
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return prices, nil
}

func main() {
	fmt.Println("--- Тест 1: Успешный запрос ---")
	prices1, err := FetchPrices(context.Background(), []int{1, 2, 4, 5, 6, 7, 8, 9, 19})
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Успех! Цены: %v\n", prices1)
	}

	fmt.Println("\n--- Тест 2: Запрос с ошибкой ---")
	prices2, err := FetchPrices(context.Background(), []int{1, 2, 3, 4, 5})
	if err != nil {
		fmt.Printf("Ожидаемая ошибка: %v\n", err)
	} else {
		fmt.Printf("Успех! Цены: %v\n", prices2)
	}
	fmt.Println("\n--- Тест 3: Лимит горутин (DDoS тест) ---")
	// Проверяем, что при лимите 2 и 5 товарах мы не уложимся в 2 секунды (если они идут пачками)
	// Или просто проверяем работоспособность на большом списке
	start := time.Now()
	bigList := []int{10, 11, 12, 13, 14, 15}
	prices3, err := FetchPrices(context.Background(), bigList)
	duration := time.Since(start)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Успех! Обработано %d товаров за %v\n", len(prices3), duration)
		if duration > time.Second*4 {
			fmt.Println("Бонус-трек: Похоже, лимит горутин работает (время выполнения увеличилось)")
		}
	}

	fmt.Println("\n--- Тест 4: Внешняя отмена (Timeout) ---")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	_, err = FetchPrices(ctx, []int{1, 2, 4})
	if err == context.DeadlineExceeded {
		fmt.Println("Успех! Функция корректно завершилась по таймауту")
	} else {
		fmt.Printf("Что-то пошло не так: ожидался таймаут, получили: %v\n", err)
	}
}
