package kewpie

import (
	"testing"
)

// BenchmarkEnqueue measures the performance of the Enqueue operation.
func BenchmarkEnqueue(b *testing.B) {
	queue := NewQueue[int]()
	for n := 0; n < b.N; n++ {
		queue.Enqueue(n)
	}
}

// BenchmarkDequeue measures the performance of the Dequeue operation.
// This fills the queue first to ensure dequeue has work to do.
func BenchmarkDequeue(b *testing.B) {
	queue := NewQueue[int]()
	// Pre-fill the queue with a known quantity of elements to dequeue
	for n := 0; n < 10000; n++ {
		queue.Enqueue(n)
	}
	b.ResetTimer() // Start timing after setup

	for n := 0; n < b.N; n++ {
		queue.Dequeue()
	}
}
