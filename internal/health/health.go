// Package health manages the health state of the server
package health

import "sync/atomic"

type Checker struct {
	ready atomic.Bool
}

func New() *Checker {
	return &Checker{}
}

func (chk *Checker) SetReady() {
	chk.ready.Store(true)
}

func (chk *Checker) SetNotReady() {
	chk.ready.Store(false)
}

func (chk *Checker) Ready() bool {
	return chk.ready.Load()
}
