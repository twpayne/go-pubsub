# go-pubsub

[![PkgGoDev](https://pkg.go.dev/badge/github.com/twpayne/go-pubsub)](https://pkg.go.dev/github.com/twpayne/go-pubsub)

Package pubsub is a simple publish-subscribe implementation using generics.

## Example

```go
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
```

## License

MIT
