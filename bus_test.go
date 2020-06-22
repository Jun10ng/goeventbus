package goeventbus

import (
	"sync"
	"testing"
	"time"
)

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

func TestSubscribe(t *testing.T) {
	bus := NewBus()
	sub := NewSub()
	bus.Subscribe("topic1", sub)
	bus.Subscribe("topic2", sub)

	subslen, err := bus.SubsLen("topic1")
	if err != nil && subslen != 1 {
		t.Errorf("topic1 subslen wrong")
	}
	t.Logf("topic1 len  = %v", subslen)

	subslen, err = bus.SubsLen("topic2")
	if err != nil && subslen != 1 {
		t.Errorf("topic2 subslen wrong")
	}
	t.Logf("topic2 len  = %v", subslen)

	sub2 := NewSub()
	bus.Subscribe("topic2", sub2)
	subslen, err = bus.SubsLen("topic2")
	if err != nil && subslen != 1 {
		t.Errorf("topic2 subslen wrong")
	}
	t.Logf("topic2 len  = %v", subslen)
}

func TestBus(t *testing.T) {
	bus := NewBus()
	sub := NewSub()
	bus.Subscribe("topic1", sub)

	var msg interface{}
	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		defer wait.Done()
		msg = sub.Out().(int)
	}()

	bus.Publish("topic1", 7)

	wait.Wait()
	t.Logf("msg is %v", msg)
	if msg != 7 {
		t.Fatalf("got wrong number:%d", msg)
	}
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
		t.Log(sum)
		if sum != 21 {
			t.Fatalf("got wrong %d", sum)
		}

		// topic2
		sum = 0
		sum = sub2.Out().(int) + sub3.Out().(int)
		t.Log(sum)
		if sum != 14 {
			t.Fatalf("got wrong %d", sum)
		}

		// topic3
		sum = 0
		sum = sub3.Out().(int)
		t.Log(sum)
		if sum != 7 {
			t.Fatalf("got wrong %d", sum)
		}
	}()

	bus.Publish("topic1", 7)
	bus.Publish("topic2", 7)
	bus.Publish("topic3", 7)

	time.Sleep(time.Duration(2) * time.Second)
}

func TestPubFunc(t *testing.T) {
	bus := NewBus()
	sub := NewSub()

	bus.Subscribe("topic1", sub)

	go func() {
		msg := sub.Out().(int)
		t.Logf("msg is %v", msg)
		if msg != 7 {
			t.Errorf("wrong num %v", msg)
		}
	}()
	PubTopic1 := bus.PubFunc("topic1")
	PubTopic1(7)

	time.Sleep(time.Duration(1) * time.Second)
}
