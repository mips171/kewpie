package kewpie

import (
	"errors"
	"log"
)

// Queue represents a generic FIFO (first in, first out) queue using a ring buffer.
type Queue[T any] struct {
	data             []T
	head, tail, size int
}

// NewQueue creates a new instance of a Queue for elements of type T with an initial capacity.
func NewQueue[T any](sizes ...int) *Queue[T] {
	var size int
	if len(sizes) > 0 {
		size = sizes[0]
	} else {
		size = 1
	}

	if size <= 0 {
		size = 1
	}
	return &Queue[T]{data: make([]T, size), head: 0, tail: 0, size: 0}
}

// Enqueue adds an element of type T to the end of the queue.
// TODO add soft limit (percentage) before resize is triggreed.
// TODO If resize fails after soft limit, then go into degraded perf mode and warn.
func (queue *Queue[T]) Enqueue(data T) {
	if queue.size == len(queue.data) {
		queue.resize(len(queue.data) * 2) // Double the size when full like a normal Go slice or map
	}
	queue.data[queue.tail] = data
	queue.tail = (queue.tail + 1) % len(queue.data)
	queue.size++
}

// Dequeue removes and returns the element at the front of the queue.
// It returns an error if the queue is empty.
func (queue *Queue[T]) Dequeue() (T, error) {
	if queue.size == 0 {
		var zero T
		return zero, errors.New("kewpie: queue is empty")
	}

	element := queue.data[queue.head]
	var zero T
	queue.data[queue.head] = zero // Clearing the reference to avoid memory leak from stale struct
	queue.head = (queue.head + 1) % len(queue.data)
	queue.size--

	// shrink queue size if too large for current needs
	if len(queue.data) > 1 && queue.size <= len(queue.data)/4 {
		queue.resize(len(queue.data) / 2)
	}

	return element, nil
}

// Peek returns the element at the front of the queue without removing it.
// It returns an error if the queue is empty.
func (queue *Queue[T]) Peek() (T, error) {
	if queue.size == 0 {
		var zero T
		return zero, errors.New("kewpie: queue is empty")
	}
	return queue.data[queue.head], nil
}

// Returns the queue's size
// Mostly just for the stress test
func (queue *Queue[T]) Size() int {
	return queue.size
}

// Resize changes the size of the queue's data slice prioritising data integrity.
func (queue *Queue[T]) resize(newCapacity int) {
	// Attempt to allocate a new slice with the new capacity.
	// Use a defer-recover mechanism to catch any panic (e.g., out of memory).
	defer func() {
		if err := recover(); err != nil {
			// If we're here, allocation failed. Don't proceed with resizing.
			// The operation is aborted but the existing data in the queue remains intact.
			log.Printf("kewpie: failed to resize the queue: %v", err)
			log.Printf("kewpie: the resize operation was aborted but the existing data in the queue remains intact")
			return
		}
	}()

	// Safety check on the capacity before copying the data
	if newCapacity <= queue.size {
		newCapacity = max(queue.size, 1)
	}

	newData := make([]T, newCapacity)
	for i := 0; i < queue.size; i++ {
		newData[i] = queue.data[(queue.head+i)%len(queue.data)]
	}

	queue.data = newData
	queue.head = 0
	queue.tail = queue.size
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
