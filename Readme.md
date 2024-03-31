# Kewpie Queue Library
Kewpie is a Go package implementing a generic, dynamic, and efficient First-In-First-Out (FIFO) queue data structure using a ring buffer. This implementation allows for enqueueing and dequeueing with amortised constant time complexity, making it suitable for various applications including breadth-first search algorithms, event processing systems, and more.

## Features
* Generic Implementation: Works with any data type.
* Dynamic Resizing: Automatically resises based on the queue's current size, ensuring efficient use of memory.
* Ring Buffer: Minimises the overhead of enqueueing and dequeueing operations.
* Error Handling: Provides clear error messages for operations that cannot be completed.


## Installation
To use the Kewpie queue in your Go project, install it by running:

```shell
go get -u github.com/mips171/kewpie
```

## Usage

Here's a simple example demonstrating how to use the Kewpie queue:

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
    bfs(root)
}
```

# Contributing
Contributions are welcome! Please submit issues and pull requests via GitHub, and ensure your code follows the Go coding standards.

# License
Kewpie is released under the MIT License. See the LICENSE file for more details.
