package goeventbus

import "testing"

func TestSub(t *testing.T) {
	sub := NewSub()

	go func() {
		msg := sub.Out().(int)
		if msg != 7 {
			t.Fatalf("got wrong number:%d", msg)
		}
	}()

	sub.receive(7)
}

func TestBus(t *testing.T) {
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
}

func TestManySubscribe(t *testing.T) {
	sub1 := NewSub()
	sub2 := NewSub()
	sub3 := NewSub()

	bus := NewBus()

	bus.Subscribe("topic1", sub1)
	bus.Subscribe("topic1", sub2)
	bus.Subscribe("topic1", sub3)

	bus.Subscribe("topic2", sub2)
	bus.Subscribe("topic2", sub3)

	bus.Subscribe("topic3", sub3)

	go func() {
		// topic1
		sum := 0
		sum = sub1.Out().(int) + sub2.Out().(int) + sub3.Out().(int)
		if sum != 21 {
			t.Fatalf("got wrong %d", sum)
		}

		// topic2
		sum = 0
		sum = sub2.Out().(int) + sub3.Out().(int)
		if sum != 14 {
			t.Fatalf("got wrong %d", sum)
		}

		// topic3
		sum = 0
		sum = sub3.Out().(int)
		if sum != 7 {
			t.Fatalf("got wrong %d", sum)
		}
	}()

	bus.Publish("topic1", 7)
	bus.Publish("topic2", 7)
	bus.Publish("topic3", 7)
}
