package event

type Topic string

type Event struct {
	Topic Topic
	Data  interface{}
}
