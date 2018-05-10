package messaging

type Publisher interface {
	Publish(message interface{})
}
