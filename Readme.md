# Kewpie Queue Library

[Chinese/中文](Readme_zh.md)

Kewpie is a Go package implementing a generic, dynamic, and efficient First-In-First-Out (FIFO) queue data structure using a ring buffer. This implementation allows for enqueueing and dequeueing with amortised constant time complexity, making it suitable for various applications including breadth-first search algorithms, event processing systems, and more.

## Features
* Generic Implementation: Works with any data type.
* Dynamic Resizing: Like Go maps or slices, queue data automatically resizes based on the queue's current size, ensuring efficient use of memory.
* Ring Buffer: Minimises the overhead of enqueueing and dequeueing operations because it just increments a pointer.
* Error Handling: Provides clear error messages for operations that cannot be completed. There are only a couple things that can go wrong anyway, and they mainly should only happen if your system is out of memory.

![kewpie](https://github.com/mips171/kewpie/assets/18670565/c48e43a5-927a-4dea-82d8-85589989ff37)

*Sushi queueing up*

## Installation
To use the Kewpie queue in your Go project, install it by running:

```shell
go get -u github.com/mips171/kewpie
```

## Usage

Here's a simple example demonstrating how to use the Kewpie queue to queeue some integers.

```go
package main

import (
    "fmt"
    "github.com/mips171/kewpie"
)

func main() {
    // Create a new queue of integers with an initial capacity of 10.
    // Note: the 10 is optional. If unspecified kewpie will default to a size of 1.
    queue := kewpie.NewQueue[int](10)

    // Enqueue some elements.
    queue.Enqueue(1)
    queue.Enqueue(2)
    queue.Enqueue(3)

    // Peek at the front element.
    front, err := queue.Peek()
    if err != nil {
        fmt.Println("Error peeking at the queue:", err)
        return
    }
    fmt.Println("The front element is:", front)

    // Dequeue and print elements.
    for {
        element, err := queue.Dequeue()
        if err != nil {
            // If we encounter an error, it means the queue is empty.
            fmt.Println("Queue is empty:", err)
            break
        }
        fmt.Println("Dequeued:", element)
    }
}
```
Output:
```sh
The front element is: 1
Dequeued: 1
Dequeued: 2
Dequeued: 3
Queue is empty: kewpie: queue is empty
```

# Kewpie in BFS

Now let's use Kewpie to queue a custom TreeNode type for a Breadth-First Search:

```go
package main

import (
    "fmt"
    "github.com/mips171/kewpie"
)

type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}


func bfs(root *TreeNode) {
    if root == nil {
        return
    }
    
    queue := kewpie.NewQueue[*TreeNode]() // Initialize a new queue for TreeNode pointers
    queue.Enqueue(root) // Start BFS with the root node
    
    for queue.Size() > 0 {
        node, err := queue.Dequeue() // Remove the next node from the queue
        if err != nil {
            fmt.Println("Error dequeuing:", err)
            return
        }
        
        fmt.Println(node.Val) // Process the current node
        
        // Add child nodes to the queue for processing
        if node.Left != nil {
            queue.Enqueue(node.Left)
        }
        if node.Right != nil {
            queue.Enqueue(node.Right)
        }
    }
}

func main() {
    // Example binary tree:
    //     1
    //    / \
    //   2   3
    //  / \   \
    // 4   5   6
    root := &TreeNode{1,
        &TreeNode{2,
            &TreeNode{4, nil, nil},
            &TreeNode{5, nil, nil},
        },
        &TreeNode{3,
            nil,
            &TreeNode{6, nil, nil},
        },
    }
    
    fmt.Println("BFS traversal of the tree:")
    bfs(root) // [1 2 3 4 5 6]
}
```

## Batch Operations
This function is ideal for larger queues, somewhere beyond 10K and 100K items at a time, increasing efficiency by minimising the number of resize operations needed. It checks if the queue's current capacity can accommodate the additional items and doubles the capacity until it is sufficient. This approach significantly reduces the overhead associated with adding multiple elements to the queue one by one, especially when dealing with large numbers of elements.

```go
package main

import (
    "fmt"
    "github.com/mips171/kewpie"
)

type Message struct {
    ID   string
    Text string
}

func main() {
    // Initialize a queue for Messages.
    queue := kewpie.NewQueue[Message]()

    // Define a batch of messages to enqueue.
    messagesToEnqueue := []Message{
        {ID: "1", Text: "Hello"},
        {ID: "2", Text: "World"},
        {ID: "3", Text: "How are you?"},
    }

    // Enqueue the batch of messages.
    queue.EnqueueBatch(messagesToEnqueue)
    
    fmt.Println("Enqueued messages")

    // Specify the batch size for dequeueing.
    batchSize := 2

    // Dequeue a batch of messages.
    messagesDequeued, err := queue.DequeueBatch(batchSize)
    if err != nil {
        fmt.Println("Error dequeuing messages:", err)
        return
    }

    for _, msg := range messagesDequeued {
        fmt.Printf("Dequeued Message ID: %s, Text: %s\n", msg.ID, msg.Text)
    }
}
```

## Benchmarks

```sh
go test -bench=. github.com/mips171/kewpie  -benchmem -benchtime=100x
```

This output shows that when using large batch sizes (somewhere between 10,000 and 100,000), that it's more efficient to use the Batch version of the function (EnqueueBatch/DequeueBatch) with an appropriate batch size to reduce the number of allocations required by Go.
Output:

```sh
goos: darwin
goarch: arm64
pkg: github.com/mips171/kewpie
BenchmarkDequeue1-10                 100               103.0 ns/op             0 B/op          0 allocs/op
BenchmarkDequeue10-10                100               983.3 ns/op           128 B/op          4 allocs/op
BenchmarkDequeue100-10               100              1959 ns/op            1024 B/op          7 allocs/op
BenchmarkDequeue1000-10              100             11058 ns/op            8192 B/op         10 allocs/op
BenchmarkDequeue10000-10             100             88383 ns/op          131074 B/op         14 allocs/op
BenchmarkEnqueue1-10                 100               139.6 ns/op            20 B/op          0 allocs/op
BenchmarkEnqueue10-10                100               177.1 ns/op           163 B/op          0 allocs/op
BenchmarkEnqueue100-10               100               915.9 ns/op          2621 B/op          0 allocs/op
BenchmarkEnqueue1000-10              100              8073 ns/op           20971 B/op          0 allocs/op
BenchmarkEnqueue10000-10             100             75846 ns/op          167773 B/op          0 allocs/op
BenchmarkEnqueueMessages/1-10        100               912.9 ns/op           258 B/op          2 allocs/op
BenchmarkEnqueueMessages/10-10       100              7733 ns/op            2310 B/op         20 allocs/op
BenchmarkEnqueueMessages/100-10                      100             67778 ns/op           28576 B/op        200 allocs/op
BenchmarkEnqueueMessages/1000-10                     100            662272 ns/op          253876 B/op       2000 allocs/op
BenchmarkEnqueueMessages/10000-10                    100           6698946 ns/op         2286749 B/op      20000 allocs/op
BenchmarkDequeueMessages/1-10                        100               806.2 ns/op           126 B/op          1 allocs/op
BenchmarkDequeueMessages/10-10                       100              7211 ns/op            3539 B/op         27 allocs/op
BenchmarkDequeueMessages/100-10                      100             66341 ns/op           32237 B/op        211 allocs/op
BenchmarkDequeueMessages/1000-10                     100            651548 ns/op          282732 B/op       1999 allocs/op
BenchmarkDequeueMessages/10000-10                    100           6671030 ns/op         3620320 B/op      19828 allocs/op
BenchmarkEnqueueBatch/BatchSize10-10                 100               115.8 ns/op          1024 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize100-10                100               897.1 ns/op         15710 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize1000-10               100             31023 ns/op          125337 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize10000-10              100            292337 ns/op          998803 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize100000-10                     100           3284543 ns/op        16043213 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize1000000-10                    100          26368697 ns/op        128345702 B/op         0 allocs/op
BenchmarkDequeueBatch/BatchSize10-10                         100          93716526 ns/op        1285310464 B/op  1000020 allocs/op
BenchmarkDequeueBatch/BatchSize100-10                        100          70284867 ns/op        1291704576 B/op   100017 allocs/op
BenchmarkDequeueBatch/BatchSize1000-10                       100          66622307 ns/op        1296777216 B/op    10014 allocs/op
BenchmarkDequeueBatch/BatchSize10000-10                      100          61337613 ns/op        1287847936 B/op     1010 allocs/op
BenchmarkDequeueBatch/BatchSize100000-10                     100          60451120 ns/op        1279066112 B/op      107 allocs/op
BenchmarkDequeueBatch/BatchSize1000000-10                    100          57964251 ns/op        1235025920 B/op       14 allocs/op
PASS
ok      github.com/mips171/kewpie       114.039s
```

# Contributing
Contributions are welcome! Please submit issues and pull requests via GitHub, and ensure your code follows the Go coding standards.

# License
Kewpie is released under the MIT License. See the LICENSE file for more details.
