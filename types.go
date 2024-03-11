package throttledqueue

import (
	"sync"
	"time"
)

type FuctionWithParams struct {
	Function func(...interface{})
	Params   []any
}

type ThrottledQueue struct {
	queue        []FuctionWithParams
	interval     time.Duration
	maxRequests  int
	numRequests  int
	lastExecuted time.Time
	timeout      *time.Timer
	wg           sync.WaitGroup
}
