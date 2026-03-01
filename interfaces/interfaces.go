package interfaces

import "fmt"


type Sender interface{
	Send(message string) error
}

type EmailSender struct{
	Email string
}

type TelegramSender struct{	
	Alias string
}

type SMSSender struct{
	Phone int
}

func (es EmailSender) Send(message string) error {
	fmt.Printf("Sending email: %s, via email from %s", message, es.Email)
	return nil
}

func (ts TelegramSender) Send(message string) error {
	fmt.Printf("Sending email : %s, via telegram from %s", message, ts.Alias)
	return nil
}

func (ss SMSSender) Send(message string) error {
	fmt.Printf("Sending email : %s, via phone from %d", message, ss.Phone)
	return nil
}

func NotifyUser(s Sender, msg string) {
	s.Send(msg)
}

func main() {
	senders := []Sender{
		&EmailSender{
			"sobaka@mail.ru",
		},
		&TelegramSender{
			"sobaka",
		},
		&SMSSender{
			89526953134,
		},
	}
	for _, sender := range senders{
		NotifyUser(sender, "hello")
	}
}
