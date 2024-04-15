// Package pubsub is a simple publish-subscribe implementation using generics.
package pubsub

import "context"

type Server[T any] struct {
	publishCh     chan<- T
	subscribeCh   chan<- chan<- T
	unsubscribeCh chan<- chan<- T
}

func NewServer[T any](ctx context.Context) *Server[T] {
	publishCh := make(chan T)
	subscribeCh := make(chan chan<- T)
	unsubscribeCh := make(chan chan<- T)
	s := &Server[T]{
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
			}
		}
	}()
	return s
}

func (s *Server[T]) Close() {
	close(s.publishCh)
}

func (s *Server[T]) Publish(value T) {
	s.publishCh <- value
}

func (s *Server[T]) Subscribe(subscriber chan<- T) {
	s.subscribeCh <- subscriber
}

func (s *Server[T]) Unsubscribe(subscriber chan<- T) {
	s.unsubscribeCh <- subscriber
}
