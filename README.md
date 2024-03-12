# go-throttled-queue

This is a golang library for creating throttled queues. It allows you to control the rate at which items are processed. Unlike the `throttle`, each item is placed in a queue and emptied at the specified rate limit. May be useful for throttling API requests.  

Inspired by a throttled-queue library for JavaScript by Shaun Persad: https://github.com/shaunpersad/throttled-queue

## Installation
```shell
go get github.com/uselesss/go-throttled-queue
```

## Usage
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

## Contributing
If you'd like to contribute to this library, please fork the repository and create a pull request. We welcome any contributions, including bug fixes, new features, and documentation improvements.

## License
This library is licensed under the MIT License. See the LICENSE file for details.
