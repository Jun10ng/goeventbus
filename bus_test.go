package goeventbus

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
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

func TestFindSubIdx(t *testing.T) {
	sub1 := NewSub()
	sub2 := NewSub()
	sub3 := NewSub()
	sub4 := NewSub()

	bus := NewBus()

	bus.Subscribe("topic1", sub1)
	bus.Subscribe("topic1", sub2)
	bus.Subscribe("topic1", sub3)

	idx1 := bus.subNode["topic1"].findSubIdx(sub1)
	idx2 := bus.subNode["topic1"].findSubIdx(sub2)
	idx3 := bus.subNode["topic1"].findSubIdx(sub3)
	idx4 := bus.subNode["topic1"].findSubIdx(sub4)

	assert.Equal(t, 0, idx1, fmt.Sprintf("sub1 index should be 0, not %v", idx1))
	assert.Equal(t, 1, idx2, fmt.Sprintf("sub2 index should be 1, not %v", idx2))
	assert.Equal(t, 2, idx3, fmt.Sprintf("sub3 index should be 2, not %v", idx3))

	assert.Equal(t, -1, idx4, fmt.Sprintf("sub4 not in subs, index should be -1, not %v", idx4))
}

func TestRemoveSub(t *testing.T) {
	sub1 := NewSub()
	sub2 := NewSub()
	sub3 := NewSub()

	bus := NewBus()

	bus.Subscribe("topic1", sub1)
	bus.Subscribe("topic1", sub2)
	bus.Subscribe("topic1", sub3)

	// PubTopic1 := bus.PubFunc("topic1")

	bus.subNode["topic1"].removeSub(sub3)

	idx := bus.subNode["topic1"].findSubIdx(sub3)

	assert.Equal(t, -1, idx, "sub3 is deleted, should be -1")
}

func TestUnScribe(t *testing.T) {
	sub := NewSub()

	bus := NewBus()

	bus.Subscribe("topic1", sub)
	PubTopic1 := bus.PubFunc("topic1")

	bus.Subscribe("topic2", sub)
	PubTopic2 := bus.PubFunc("topic2")

	bus.UnSubscribe("topic1", sub)

	go func() {
		msg := sub.Out().(string)
		// fmt.Println(msg)
		assert.Equal(t, "form topic2", msg, "sub not unsubsrible topic1 success")
	}()

	PubTopic1("form topic1")
	PubTopic2("form topic2")

	time.Sleep(time.Duration(1) * time.Second)
}
