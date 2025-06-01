# go-throttled-queue

[![coverage](https://camo.githubusercontent.com/b167c2e10fd87d7d2d8afc04fd0bad9c8dadce0986e43365d200863d193d94aa/68747470733a2f2f696d672e736869656c64732e696f2f636f6465636f762f632f6769746875622f6477796c2f686170692d617574682d6a7774322e7376673f6d61784167653d32353932303030)](./cover.html)

This is a golang library for creating throttled queues. It allows you to control the rate at which items are processed. Unlike the `throttle`, each item is placed in a queue and emptied at the specified rate limit. May be useful for throttling API requests.  

Inspired by a throttled-queue library for JavaScript by Shaun Persad: https://github.com/shaunpersad/throttled-queue

## Installation
```shell
go get github.com/uselesss/go-throttled-queue
```

## Usage

#### Simple queue
```go
package main

import (
	"fmt"
	"time"

	"github.com/uselesss/go-throttled-queue/ttq"
)

func main() {
	// Define a callback function
	throttledCallback := func(params ...interface{}) {
		fmt.Println("Executed callback with id:", params[0])
	}

	// Create a new throttled queue with a maximum of 5 calls/second
	throttledQueue := ttq.New(time.Second, 5)

	// Enqueue items
	for i := 0; i < 20; i++ {
		throttledQueue.Enqueue(throttledCallback, i)
	}

	// Wait until items are processed
	throttledQueue.Wait()
	fmt.Println("Done")
}
```
  
#### Concurrent access
```go
package main

import (
	"fmt"
	"time"

	"github.com/uselesss/go-throttled-queue/ttq"
)

func main() {
	var results sync.Map

	// Define a callback function
	throttledCallback := func(params ...interface{}) {
		fmt.Println("Started callback with id:", params[0])

		// Simulate heavy tasks
		ch := make(chan bool)

		go func() {
			time.Sleep(1 * time.Second)
			ch <- true
		}()

		<-ch
		close(ch)

		results.Store(params[0], "computed value")

		fmt.Println("Finished callback with id:", params[0])
	}

	// Create a new throttled queue with a maximum of 5 calls/second
	throttledQueue := ttq.New(time.Second, 5)

	// Enqueue items
	for i := 0; i < 20; i++ {
		throttledQueue.Enqueue(throttledCallback, i)
	}

	// Wait until items are processed
	throttledQueue.Wait()
	fmt.Println("Done")
}
```


## Contributing
If you'd like to contribute to this library, please fork the repository and create a pull request. We welcome any contributions, including bug fixes, new features, and documentation improvements.

## License
This library is licensed under the MIT License. See the LICENSE file for details.
