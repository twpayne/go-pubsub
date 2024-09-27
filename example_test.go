package pubsub_test

import (
	"context"
	"fmt"

	"github.com/twpayne/go-pubsub"
)

func ExampleTopic() {
	ctx := context.Background()

	topic := pubsub.NewTopic[int](ctx)

	ch := make(chan int)
	topic.Subscribe(ch)
	go func() {
		defer topic.Close()
		for i := range 3 {
			topic.Publish(i)
		}
	}()

	for i := range ch {
		fmt.Println(i)
	}

	// Output:
	// 0
	// 1
	// 2
}
