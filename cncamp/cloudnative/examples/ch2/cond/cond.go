package main

import (
	"fmt"
	"sync"
	"time"
)

type (
	Queue struct {
		queue []string
		cond  *sync.Cond
	}
)

func (q *Queue) Enqueue(item string) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.queue = append(q.queue, item)
	fmt.Printf("en-queuing: putting item %s, notify all\n", item)
	q.cond.Broadcast()
}

func (q *Queue) Dequeue() string {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if len(q.queue) == 0 {
		fmt.Println("de-queuing: no data, wait")
		q.cond.Wait()
	}

	ret := q.queue[0]
	q.queue = q.queue[1:]
	fmt.Printf("de-queuing: getting item %s\n", ret)
	return ret
}

func main() {
	q := Queue{
		queue: []string{},
		cond:  sync.NewCond(&sync.Mutex{}),
	}

	// en-queue
	go func() {
		for {
			q.Enqueue("en")
			time.Sleep(2 * time.Second)
		}
	}()

	// de-queue
	for {
		q.Dequeue()
		time.Sleep(time.Second)
	}
}
