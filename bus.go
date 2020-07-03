package goeventbus

import (
	"errors"
	"sync"
)

// Bus type contains all topics and subscribers
type Bus struct {
	subNode map[string]*node
	rw      sync.RWMutex
}

func NewBus() Bus {
	return Bus{
		subNode: make(map[string]*node),
		rw:      sync.RWMutex{},
	}
}

// Subscribe means a subscriber follw a topic
func (b *Bus) Subscribe(topic string, sub Sub) {
	b.rw.Lock()
	if n, ok := b.subNode[topic]; ok {
		// found the node
		b.rw.Unlock()
		n.rw.Lock()
		defer n.rw.Unlock()
		n.subs = append(n.subs, sub)
	} else {
		defer b.rw.Unlock()
		n := NewNode()
		b.subNode[topic] = &n
		n.subs = append(n.subs, sub)
	}
}

func (b *Bus) UnSubscribe(topic string, sub Sub) {
	b.rw.Lock()
	if n, ok := b.subNode[topic]; ok && (n.SubsLen() > 0) {
		b.rw.Unlock()
		b.subNode[topic].removeSub(sub)
	} else {
		b.rw.Unlock()
		return
	}
}

// Publish msg to all subscriber who have subscribed to the topic
func (b *Bus) Publish(topic string, msg interface{}) error {
	b.rw.Lock()
	if n, ok := b.subNode[topic]; ok {
		// found the node
		b.rw.Unlock()
		n.rw.RLock()
		defer n.rw.RUnlock()
		// got the subs list and publish msg
		go func(subs []Sub, msg interface{}) {
			for _, sub := range subs {
				sub.receive(msg)
			}
		}(n.subs, msg)
		// successed return null
		return nil
	} else {
		// topic not exist
		defer b.rw.Unlock()
		return errors.New("topic not exist")
	}

}

// PubFunc return a function that publish msg to one topic
func (b *Bus) PubFunc(topic string) func(msg interface{}) {
	return func(msg interface{}) {
		b.Publish(topic, msg)
	}
}

// SubsLen return the length of subs contained topic node
func (b *Bus) SubsLen(topic string) (int, error) {
	b.rw.Lock()
	if n, ok := b.subNode[topic]; ok {
		// found the node
		b.rw.Unlock()
		n.rw.RLock()
		defer n.rw.RUnlock()
		return n.SubsLen(), nil
	} else {
		// topic not exist
		defer b.rw.Unlock()
		return 0, errors.New("topic not exist")
	}
}
