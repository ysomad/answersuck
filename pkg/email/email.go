package email

type Sender interface {
	Send(l Letter) error
}
