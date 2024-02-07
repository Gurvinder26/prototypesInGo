package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type ConcurrentQueue struct {
	queue []int
	mu    sync.Mutex
}

func (q *ConcurrentQueue) Enqueue(num int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, num)
}

func (q *ConcurrentQueue) Dequeue() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		panic("No items present in the queue")
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item
}

func (q *ConcurrentQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}

func main() {

	queue := ConcurrentQueue{queue: make([]int, 0)}

	var wgE = sync.WaitGroup{}
	var wgD = sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		wgE.Add(1)
		go func() {
			queue.Enqueue(rand.Int())
			wgE.Done()
		}()
	}

	for i := 0; i < 10000; i++ {
		wgD.Add(1)
		go func() {
			queue.Dequeue()
			wgD.Done()
		}()
	}

	wgE.Wait()
	wgD.Wait()
	fmt.Println(queue.Size())
}
