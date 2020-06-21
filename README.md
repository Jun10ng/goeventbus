## What

a simple event bus in golang

## Use

```golang
	bus := NewBus()
	sub := NewSub()

	bus.Subscribe("topic1", sub)

	go func() {
		msg := sub.Out().(int)
		if msg != 7 {
			t.Fatalf("got wrong number:%d", msg)
		}
	}()
	bus.Publish("topic1", 7)
```
For more details,see [testfile](https://github.com/Jun10ng/goeventbus/blob/master/bus_test.go)
