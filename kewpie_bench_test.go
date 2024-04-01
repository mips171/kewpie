package kewpie_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mips171/kewpie"
)

func benchmarkEnqueue(b *testing.B, queue *kewpie.Queue[int], size int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		b.StartTimer()
		for n := 0; n < size; n++ {
			queue.Enqueue(n)
		}
	}
}

func benchmarkDequeue(b *testing.B, queue *kewpie.Queue[int], size int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for n := 0; n < size; n++ {
			queue.Enqueue(n)
		}
		b.StartTimer()
		for n := 0; n < size; n++ {
			queue.Dequeue()
		}
	}
}

func BenchmarkDequeue1(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkDequeue(b, queue, 1)
}

func BenchmarkDequeue10(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkDequeue(b, queue, 10)
}

func BenchmarkDequeue100(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkDequeue(b, queue, 100)
}

func BenchmarkDequeue1000(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkDequeue(b, queue, 1000)
}

func BenchmarkDequeue10000(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkDequeue(b, queue, 10000)
}

func BenchmarkEnqueue1(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkEnqueue(b, queue, 1)
}

func BenchmarkEnqueue10(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkEnqueue(b, queue, 10)
}
func BenchmarkEnqueue100(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkEnqueue(b, queue, 100)
}

func BenchmarkEnqueue1000(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkEnqueue(b, queue, 1000)
}

func BenchmarkEnqueue10000(b *testing.B) {
	queue := kewpie.NewQueue[int]()
	benchmarkEnqueue(b, queue, 10000)
}

// BenchmarkEnqueueMessages benchmarks the enqueue operation with varying numbers of messages.
func BenchmarkEnqueueMessages(b *testing.B) {
	sizes := []int{1, 10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			queue := kewpie.NewQueue[Message]()
			rand.Seed(time.Now().UnixNano())
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < size; n++ {
					msg := generateRandomMessage()
					queue.Enqueue(msg)
				}
			}
		})
	}
}

// BenchmarkDequeueMessages benchmarks the dequeue operation with varying numbers of messages.
func BenchmarkDequeueMessages(b *testing.B) {
	sizes := []int{1, 10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			queue := kewpie.NewQueue[Message]() // Initialize the queue for Message structs
			rand.Seed(time.Now().UnixNano())    // Seed the random number generator

			// Pre-fill the queue with the required number of messages outside the b.N loop to ensure it doesn't affect benchmark timing
			for n := 0; n < size; n++ {
				msg := generateRandomMessage()
				queue.Enqueue(msg)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n := 0; n < size; n++ {
					_, err := queue.Dequeue()
					if err != nil {
						b.Fatalf("Dequeue failed: %v", err)
					}
				}

				// Re-fill the queue after timing to ensure it's ready for the next iteration
				if i < b.N-1 {
					for n := 0; n < size; n++ {
						queue.Enqueue(generateRandomMessage())
					}
				}
			}
		})
	}
}

func BenchmarkEnqueueBatch(b *testing.B) {
	batchSizes := []int{10, 100, 1000, 10000, 100000, 1000000}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize%d", batchSize), func(b *testing.B) {
			q := kewpie.NewQueue[Message]()

			batch := make([]Message, batchSize)
			for i := range batch {
				batch[i] = Message{ID: int64(i), Content: fmt.Sprintf("Message %d", i)}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				q.EnqueueBatch(batch)
			}
		})
	}
}

func BenchmarkDequeueBatch(b *testing.B) {
	batchSizes := []int{10, 100, 1000, 10000, 100000, 1000000}
	queueSize := 10000000

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize%d", batchSize), func(b *testing.B) {
			q := kewpie.NewQueue[Message]()

			// Pre-fill the queue efficiently
			preFillBatch := make([]Message, queueSize)
			for i := range preFillBatch {
				preFillBatch[i] = Message{ID: int64(i), Content: fmt.Sprintf("Message %d", i)}
			}
			q.EnqueueBatch(preFillBatch)

			dequeueIterations := queueSize / batchSize

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < dequeueIterations; j++ {
					q.DequeueBatch(batchSize)
				}

				b.StopTimer()

				q = kewpie.NewQueue[Message]()

				q.EnqueueBatch(preFillBatch)
				b.StartTimer()
			}
		})
	}
}

// Message represents a message in the queue with random data.
type Message struct {
	ID      int64
	Content string
	Time    time.Time
}

// generateRandomMessage generates a random Message instance.
func generateRandomMessage() Message {
	return Message{
		ID:      rand.Int63(),
		Content: generateRandomString(50),
		Time:    time.Now(),
	}
}

// generateRandomString generates a random string of a given length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
