package goeventbus

import "sync"

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
