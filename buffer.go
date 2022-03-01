package giotgo

import "sync"

type buffer struct {
	sync.Mutex
	bytes []byte
}

func newBuffer(size uint32) *buffer {
	b := &buffer{}
	b.bytes = make([]byte, size)
	return b
}
