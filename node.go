package goeventbus

import (
	"sync"
)

// node contains a slice of subs that subscribe same topic
type node struct {
	subs []Sub
	// Note node's rw will not be used when bus's rw is helded.
	rw sync.RWMutex
}

// NewNode return new node
func NewNode() node {
	return node{
		subs: []Sub{},
		rw:   sync.RWMutex{},
	}
}

func (n *node) SubsLen() int {
	return len(n.subs)
}

// if sub not exisit, no change
func (n *node) removeSub(s Sub) {
	// lenOfSubs should placed before Lock Statement, otherwise it will cause a deadlock.
	lenOfSubs := len(n.subs)
	n.rw.Lock()
	defer n.rw.Unlock()
	idx := n.findSubIdx(s)
	if idx < 0 {
		return
	}
	copy(n.subs[idx:], n.subs[idx+1:])
	n.subs[lenOfSubs-1] = Sub{}
	n.subs = n.subs[:lenOfSubs-1]
}

//findSubIdx return index of sub, if sub not exisit return -1.
func (n *node) findSubIdx(s Sub) int {
	for idx, sub := range n.subs {
		if sub == s {
			return idx
		}
	}
	return -1
}
