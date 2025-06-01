package ttq

import (
	"time"
)

// Create a new throttled queue
//
// Param: `interval` throttle time interval
//
// Param: `maxRequests` maximum amount of requests per interval
func New(interval time.Duration, maxRequests int) *ThrottledQueue {
	return &ThrottledQueue{
		interval:     interval,
		maxRequests:  maxRequests,
		numRequests:  0,
		lastExecuted: time.Now(),
	}
}

// Enqueue new throttled callback function
//
// Param: `callback` function that will be called when its turn has come
//
// Param: `args` arguments passed to the callback function
func (q *ThrottledQueue) Enqueue(callback func(...interface{}), args ...any) {

	if q.numRequests < q.maxRequests {
		q.numRequests++
		go callback(args...)
	} else {
		q.queue = append(q.queue, FuctionWithParams{callback, args})

		if q.timeout == nil {
			q.wg.Add(1)
			q.timeout = time.AfterFunc(time.Until(q.lastExecuted.Add(q.interval)), q.dequeue)
		}
	}
}

// Wait blocks until queue is empty
func (q *ThrottledQueue) Wait() {
	q.wg.Wait()
}

// Dequeue all available requests and queue up requests that didn't qualify
func (q *ThrottledQueue) dequeue() {
	defer q.wg.Done()
	q.lastExecuted = time.Now()
	q.numRequests = 0

	for _, function := range q.queue[0:q.maxRequests] {
		q.numRequests++
		go function.Function(function.Params...)
	}
	q.queue = q.queue[q.maxRequests:]

	if len(q.queue) > 0 {
		q.wg.Add(1)
		q.timeout = time.AfterFunc(q.interval, q.dequeue)
	} else {
		q.timeout = nil
	}
}
