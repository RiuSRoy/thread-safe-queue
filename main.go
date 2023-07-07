package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type ConcurrentQueue struct {
	queue []int32
	mu    sync.Mutex
}

func (q *ConcurrentQueue) Enqueue(item int32) {
	q.mu.Lock()
	defer q.mu.Unlock()             // anything post defer will be run at the end of the function
	q.queue = append(q.queue, item) // This particlular operation is not thread safe, so wrap it in a mutex
}

func (q *ConcurrentQueue) Dequeue() int32 {
	if len(q.queue) == 0 {
		panic("cannot dequeue from an empty queue")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item
}

var wgE sync.WaitGroup
var wgF sync.WaitGroup

func main() {

	q1 := ConcurrentQueue{
		queue: make([]int32, 0),
	}

	for i := 0; i < 100000; i++ {
		wgE.Add(1) // everytime I spin up a new thread, I write this command
		go func() {
			q1.Enqueue(rand.Int31())
			wgE.Done()
		}()
	}

	wgE.Wait()                 // Will make sure we wait for all go routines to complete
	fmt.Println(len(q1.queue)) // <1M because lot many threads are overwriting at the same location of the queue if not using lock on critical section

	for i := 0; i < 100000; i++ {
		wgF.Add(1)
		go func() {
			q1.Dequeue()
			wgF.Done()
		}()
	}

	wgF.Wait()
	fmt.Println(len(q1.queue))
}
