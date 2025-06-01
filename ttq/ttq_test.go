package ttq

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestQueueCreation(t *testing.T) {
	/**
	* Case: Queue creation
	*
	* Tests if queue object is initialized with correct params
	 */

	var emptyTimeout *time.Timer
	var tests = []struct {
		name string
		want interface{}
		get  func(*ThrottledQueue) interface{}
	}{
		{
			name: "interval",
			want: time.Second * 2,
			get:  func(q *ThrottledQueue) interface{} { return q.interval },
		},
		{
			name: "maxRequests",
			want: 1,
			get:  func(q *ThrottledQueue) interface{} { return q.maxRequests },
		},
		{
			name: "numRequests",
			want: 0,
			get:  func(q *ThrottledQueue) interface{} { return q.numRequests },
		},
		{
			name: "lastExecuted",
			want: time.Now().Unix(),
			get:  func(q *ThrottledQueue) interface{} { return q.lastExecuted.Unix() },
		},
		{
			name: "timeout",
			want: emptyTimeout,
			get:  func(q *ThrottledQueue) interface{} { return q.timeout },
		},
		{
			name: "wg",
			want: &sync.WaitGroup{},
			get:  func(q *ThrottledQueue) interface{} { return &q.wg },
		},
	}

	queue := New(time.Second*2, 1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := tt.get(queue)
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v for %s", ans, tt.want, tt.name)
			}
		})
	}

}

func TestThrottledQueueBasic(t *testing.T) {
	/**
	* Case: Basic test
	* Queue config: 10 requests / 100 ms, 100 total requests
	*
	* Expected outcome: 10 batches of 10 requests, with 100 ms intervals
	 */

	var results sync.Map
	queue := New(time.Millisecond*100, 10)
	tasksCount := 100

	callback := func(params ...interface{}) {
		callbackId, ok := params[0].(int)
		if !ok {
			panic("Could not convert param to int")
		}
		results.Store(callbackId, time.Now().UnixMilli())
	}

	for i := 0; i < tasksCount; i++ {
		queue.Enqueue(callback, i)
	}

	startTimeMS := time.Now().UnixMilli()

	queue.Wait()

	// Check results
	for i := 0; i < tasksCount; i++ {
		name := "task #" + strconv.Itoa(i)
		t.Run(name, func(t *testing.T) {
			result, _ := results.Load(i)
			res, _ := result.(int64)

			batch := int64(i / 10)
			expectedTime := startTimeMS + batch*100

			// allowed 10ms error (10% of the interval)
			if res-expectedTime > 10 {
				t.Errorf("got %v, want %v for %s", result, expectedTime, name)
			}
		})
	}

}

func TestComputeHeavyCallback(t *testing.T) {
	/**
	* Case: compute-heavy callback test
	* Queue config: 10 requests / 1000 ms, 100 total requests
	*
	* Expected outcome: 10 batches of 10 requests, with 1000 ms intervals
	 */

	var results sync.Map
	queue := New(time.Millisecond*1000, 10)
	tasksCount := 100

	callback := func(params ...interface{}) {
		callbackId, ok := params[0].(int)
		if !ok {
			panic("Could not convert param to int")
		}

		// write callback start time

		results.Store(callbackId, time.Now().UnixMilli())

		// simulate heavy task

		ch := make(chan bool)

		go func() {
			time.Sleep(5 * time.Second)
			ch <- true
		}()

		<-ch
		close(ch)
	}

	for i := 0; i < tasksCount; i++ {
		queue.Enqueue(callback, i)
	}

	startTimeMS := time.Now().UnixMilli()

	queue.Wait()

	// Check results
	for i := 0; i < tasksCount; i++ {
		name := "task #" + strconv.Itoa(i)
		t.Run(name, func(t *testing.T) {
			result, _ := results.Load(i)
			res, _ := result.(int64)

			batch := int64(i / 10)
			expectedTime := startTimeMS + batch*1000

			// allowed 10ms error (1% of the interval)
			if res-expectedTime > 10 {
				t.Errorf("got %v, want %v for %s", result, expectedTime, name)
			}
		})
	}

}
