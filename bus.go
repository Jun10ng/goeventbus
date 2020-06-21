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
		return errors.New("topic not exist")
	}

}
