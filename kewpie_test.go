package kewpie_test

import (
	"testing"

	"github.com/mips171/kewpie"

	"github.com/stretchr/testify/assert"
)

func TestEnqueueDequeue(t *testing.T) {
	q := kewpie.NewQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)

	val, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 2, val)

	_, err = q.Dequeue()
	assert.Error(t, err, "queue is empty")
}

func TestPeek(t *testing.T) {
	q := kewpie.NewQueue[string]()
	q.Enqueue("first")
	q.Enqueue("second")

	val, err := q.Peek()
	assert.NoError(t, err)
	assert.Equal(t, "first", val)

	_, err = q.Dequeue() // Remove "first"
	assert.NoError(t, err)

	val, err = q.Peek() // "second" is now at the front
	assert.NoError(t, err)
	assert.Equal(t, "second", val)

	_, err = q.Dequeue() // Remove "second"
	assert.NoError(t, err)

	_, err = q.Peek() // Queue is empty now
	assert.Error(t, err, "queue is empty")
}

func TestEmptyQueue(t *testing.T) {
	q := kewpie.NewQueue[float64]()

	_, err := q.Dequeue()
	assert.Error(t, err, "queue is empty")

	_, err = q.Peek()
	assert.Error(t, err, "queue is empty")
}

func TestQueueStress(t *testing.T) {
	const count = 1000000000    // Number of elements to test with
	q := kewpie.NewQueue[int]() // Initialize a new queue

	// Enqueue a large number of elements
	for i := 0; i < count; i++ {
		q.Enqueue(i)
	}

	// Check the size of the queue to ensure all elements were added
	if q.Size() != count {
		t.Errorf("expected queue size %d, got %d", count, q.Size())
	}

	// Dequeue elements and check that they are in the correct order
	for i := 0; i < count; i++ {
		val, err := q.Dequeue()
		if err != nil {
			t.Fatal("dequeue failed:", err)
		}
		if val != i {
			t.Errorf("expected %d, got %d", i, val)
			break
		}
	}

	// Ensure the queue is empty after all operations
	if q.Size() != 0 {
		t.Errorf("expected queue to be empty, size is %d", q.Size())
	}
}

// TreeNode represents a node in a binary tree.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// TestBFS uses the kewpie Queue to perform a Breadth-First Search (BFS) on a binary tree.
func TestBFS(t *testing.T) {
	// Create a simple binary tree for testing.
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

	expectedOrder := []int{1, 2, 3, 4, 5, 6} // The expected BFS order of node values
	var resultOrder []int                    // To store the order of node values obtained via BFS

	queue := kewpie.NewQueue[*TreeNode]()
	queue.Enqueue(root)

	for queue.Size() > 0 {
		node, err := queue.Dequeue()
		if err != nil {
			t.Fatal("Error dequeuing:", err)
		}

		resultOrder = append(resultOrder, node.Val) // Add the current node's value to the result slice

		// Enqueue child nodes
		if node.Left != nil {
			queue.Enqueue(node.Left)
		}
		if node.Right != nil {
			queue.Enqueue(node.Right)
		}
	}
	// [1 2 3 4 5 6]
	assert.Equal(t, expectedOrder, resultOrder, "The BFS traversal order did not match the expected order.")
}
