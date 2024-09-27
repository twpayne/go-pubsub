// Package pubsub is a simple publish-subscribe implementation using generics.
package pubsub

import "context"

// A Topic is a pub-sub server that handles messages of type T.
type Topic[T any] struct {
	publishCh     chan<- T
	subscribeCh   chan<- chan<- T
	unsubscribeCh chan<- chan<- T
}

// NewTopic returns a new Topic. It will terminate when ctx is done or when
// Close is called.
func NewTopic[T any](ctx context.Context) *Topic[T] {
	publishCh := make(chan T)
	subscribeCh := make(chan chan<- T)
	unsubscribeCh := make(chan chan<- T)
	t := &Topic[T]{
		publishCh:     publishCh,
		subscribeCh:   subscribeCh,
		unsubscribeCh: unsubscribeCh,
	}
	go func() {
		subscribers := make(map[chan<- T]struct{})
		defer func() {
			for subscriber := range subscribers {
				close(subscriber)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-publishCh:
				if !ok {
					return
				}
				for subscriber := range subscribers {
					subscriber <- value
				}
			case subscriber := <-subscribeCh:
				subscribers[subscriber] = struct{}{}
			case subscriber := <-unsubscribeCh:
				delete(subscribers, subscriber)
				close(subscriber)
			}
		}
	}()
	return t
}

// Close closes t.
func (t *Topic[T]) Close() {
	close(t.publishCh)
}

// PublishContext publishes value to all subscribers.
func (t *Topic[T]) Publish(value T) {
	t.publishCh <- value
}

// PublishContext publishes value to all subscribers. If ctx is done then it
// returns without publishing value.
func (t *Topic[T]) PublishContext(ctx context.Context, value T) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case t.publishCh <- value:
		return nil
	}
}

// Subscribe adds ch as a subscriber. t takes ownership of ch and will close it
// when t terminates.
func (t *Topic[T]) Subscribe(ch chan<- T) {
	t.subscribeCh <- ch
}

// SubscribeContext adds ch as a subscriber. t takes ownership of ch and will
// close it when t terminates. If ctx is done then it returns without
// subscribing ch and does not take ownership of ch.
func (t *Topic[T]) SubscribeContext(ctx context.Context, ch chan<- T) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case t.subscribeCh <- ch:
		return nil
	}
}

// Unsubscribe removes ch as a subscriber. t will close ch when the
// unsubscription is complete.
func (t *Topic[T]) Unsubscribe(ch chan<- T) {
	t.unsubscribeCh <- ch
}

// UnsubscribeContext removes ch as a subscriber. t will close ch when the
// unsubscription is complete. If ctx is done then it returns without
// unsubscribing ch.
func (t *Topic[T]) UnsubscribeContext(ctx context.Context, ch chan<- T) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case t.unsubscribeCh <- ch:
		return nil
	}
}
