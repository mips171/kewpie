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
BenchmarkDequeue1-10                 100               110.1 ns/op             0 B/op          0 allocs/op
BenchmarkDequeue10-10                100               940.0 ns/op           128 B/op          4 allocs/op
BenchmarkDequeue100-10               100              2998 ns/op            1024 B/op          7 allocs/op
BenchmarkDequeue1000-10              100             11769 ns/op            8192 B/op         10 allocs/op
BenchmarkDequeue10000-10             100             82939 ns/op          131074 B/op         14 allocs/op
BenchmarkDequeue100000-10            100            582988 ns/op         1048579 B/op         17 allocs/op
BenchmarkDequeue1000000-10           100           5343473 ns/op         8388608 B/op         20 allocs/op
BenchmarkEnqueue1-10                 100               172.9 ns/op            20 B/op          0 allocs/op
BenchmarkEnqueue10-10                100               208.4 ns/op           163 B/op          0 allocs/op
BenchmarkEnqueue100-10               100               974.2 ns/op          2621 B/op          0 allocs/op
BenchmarkEnqueue1000-10              100              7812 ns/op           20971 B/op          0 allocs/op
BenchmarkEnqueue10000-10             100             67582 ns/op          167775 B/op          0 allocs/op
BenchmarkEnqueue100000-10            100            779382 ns/op         2684354 B/op          0 allocs/op
BenchmarkEnqueue1000000-10           100           6926040 ns/op        21474838 B/op          0 allocs/op
BenchmarkEnqueueMessages/1-10        100               864.2 ns/op           258 B/op          2 allocs/op
BenchmarkEnqueueMessages/10-10       100              6680 ns/op            2310 B/op         20 allocs/op
BenchmarkEnqueueMessages/100-10                      100             64748 ns/op           28576 B/op        200 allocs/op
BenchmarkEnqueueMessages/1000-10                     100            673420 ns/op          253876 B/op       2000 allocs/op
BenchmarkEnqueueMessages/10000-10                    100           6548411 ns/op         2286715 B/op      20000 allocs/op
BenchmarkEnqueueMessages/100000-10                   100          66954072 ns/op        28906175 B/op     200000 allocs/op
BenchmarkEnqueueMessages/1000000-10                  100         821247701 ns/op        256849066 B/op   2000000 allocs/op
BenchmarkDequeueMessages/1-10                        100               682.1 ns/op           126 B/op          1 allocs/op
BenchmarkDequeueMessages/10-10                       100             15450 ns/op            3539 B/op         27 allocs/op
BenchmarkDequeueMessages/100-10                      100            148808 ns/op           32237 B/op        211 allocs/op
BenchmarkDequeueMessages/1000-10                     100            807242 ns/op          282730 B/op       1999 allocs/op
BenchmarkDequeueMessages/10000-10                    100           6874234 ns/op         3620356 B/op      19828 allocs/op
BenchmarkDequeueMessages/100000-10                   100          67830701 ns/op        31430091 B/op     198034 allocs/op
BenchmarkDequeueMessages/1000000-10                  100         671253912 ns/op        276717879 B/op   1980040 allocs/op
BenchmarkEnqueueBatch/BatchSize10-10                 100               107.1 ns/op          1024 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize100-10                100              1022 ns/op           15710 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize1000-10               100             30048 ns/op          125337 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize10000-10              100            298482 ns/op          998769 B/op          0 allocs/op
BenchmarkDequeueBatch/BatchSize10-10                 100           6044221 ns/op        22456293 B/op      50017 allocs/op
BenchmarkDequeueBatch/BatchSize100-10                100           4779223 ns/op        21992290 B/op       8017 allocs/op
BenchmarkDequeueBatch/BatchSize1000-10               100           3067474 ns/op        18515481 B/op       1117 allocs/op
BenchmarkDequeueBatch/BatchSize10000-10              100           5025982 ns/op        29554694 B/op        207 allocs/op
PASS
ok      github.com/mips171/kewpie       175.586s
```

# Contributing
Contributions are welcome! Please submit issues and pull requests via GitHub, and ensure your code follows the Go coding standards.

# License
Kewpie is released under the MIT License. See the LICENSE file for more details.
