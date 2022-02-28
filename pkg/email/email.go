package email

type Sender interface {
	Send(to, subject, body string) error
}
