package event

import "context"

type Bus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(subscriber Subscriber) error
}

type bus struct {
	subscribers map[Topic][]Subscriber
}

func (bus *bus) Publish(ctx context.Context, event Event) error {
	subscribers := bus.subscribers[event.Topic]
	for _, sub := range subscribers {
		if err := sub.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (bus *bus) Subscribe(subscriber Subscriber) error {
	topic := subscriber.Topic()
	if bus.subscribers[topic] == nil {
		bus.subscribers[topic] = make([]Subscriber, 0)
	}
	bus.subscribers[topic] = append(bus.subscribers[topic], subscriber)
	return nil
}

func NewBus() Bus {
	return &bus{
		subscribers: map[Topic][]Subscriber{},
	}
}
