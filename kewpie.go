package kewpie

import (
	"errors"
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

// Enqueue adds an element of type T to the end of the q.
// TODO add soft limit (percentage) before resize is triggreed.
// TODO If resize fails after soft limit, then go into degraded perf mode and warn.
func (q *Queue[T]) Enqueue(data T) {
	if q.size == len(q.data) {
		q.resize(len(q.data) * 2) // Double the size when full like a normal Go slice or map
	}
	q.data[q.tail] = data
	q.tail = (q.tail + 1) % len(q.data)
	q.size++
}

// EnqueueBatch adds multiple elements of type T to the end of the queue, minimising number of resize operations.
func (q *Queue[T]) EnqueueBatch(items []T) {
	batchSize := len(items)
	if batchSize == 0 {
		return
	}

	requiredCapacity := q.size + batchSize
	currentCapacity := len(q.data)

	// Check if resizing is necessary
	if requiredCapacity > currentCapacity {
		newCapacity := currentCapacity
		// double the capacity until it can fit the new items
		for newCapacity < requiredCapacity {
			newCapacity *= 2
		}

		q.resize(newCapacity)
	}

	for _, item := range items {
		q.data[q.tail] = item
		q.tail = (q.tail + 1) % len(q.data) // Ensure the tail wraps around correctly
		q.size++
	}
}

// Dequeue removes and returns the element at the front of the q.
// It returns an error if the queue is empty.
func (q *Queue[T]) Dequeue() (T, error) {
	if q.size == 0 {
		var zero T
		return zero, errors.New("kewpie: queue is empty")
	}

	element := q.data[q.head]
	var zero T
	q.data[q.head] = zero // Clearing the reference to avoid memory leak from stale struct
	q.head = (q.head + 1) % len(q.data)
	q.size--

	// shrink queue size if too large for current needs
	if len(q.data) > 1 && q.size <= len(q.data)/4 {
		q.resize(len(q.data) / 2)
	}

	return element, nil
}

// DequeueBatch dequeues messages up to the specified batchSize.
func (q *Queue[Message]) DequeueBatch(batchSize int) ([]Message, error) {
	var batch []Message
	for i := 0; i < batchSize; i++ {
		if q.size == 0 {
			break
		}
		msg, err := q.Dequeue()
		if err != nil {
			return nil, err
		}
		batch = append(batch, msg)
	}
	return batch, nil
}

// Peek returns the element at the front of the queue without removing it.
// It returns an error if the queue is empty.
func (q *Queue[T]) Peek() (T, error) {
	if q.size == 0 {
		var zero T
		return zero, errors.New("kewpie: queue is empty")
	}
	return q.data[q.head], nil
}

// Returns the queue's size
// Mostly just for the stress test
func (q *Queue[T]) Size() int {
	return q.size
}

// Resize changes the size of the queue's data slice prioritising data integrity.
func (q *Queue[T]) resize(newCapacity int) {
	// Attempt to allocate a new slice with the new capacity.
	// Use a defer-recover mechanism to catch any panic (e.g., out of memory).
	defer func() {
		if err := recover(); err != nil {
			// If we're here, allocation failed. Don't proceed with resizing.
			return
		}
	}()

	// Safety check on the capacity before copying the data
	if newCapacity <= q.size {
		newCapacity = max(q.size, 1)
	}

	newData := make([]T, newCapacity)
	for i := 0; i < q.size; i++ {
		newData[i] = q.data[(q.head+i)%len(q.data)]
	}

	q.data = newData
	q.head = 0
	q.tail = q.size
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
