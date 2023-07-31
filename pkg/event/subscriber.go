package event

import "context"

type Subscriber interface {
	Topic() Topic
	Handle(ctx context.Context, event Event) error
}
