package websocket

import (
	"sync"

	"gopkg.in/eapache/queue.v1"
)

type MessageQueue struct {
	queue *queue.Queue
	cond  *sync.Cond

	shuttingDown bool
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		queue: queue.New(),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}

func (q *MessageQueue) Length() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.queue.Length()
}

// return false when MessageQueue is shuttingDown
func (q *MessageQueue) Add(elem *WSMessage) bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.shuttingDown {
		return false
	}
	q.queue.Add(elem)
	q.cond.Signal()
	return true
}

func (q *MessageQueue) Get() (msg *WSMessage, shutdown bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for q.queue.Length() == 0 && !q.shuttingDown {
		q.cond.Wait()
	}
	if q.queue.Length() == 0 {
		// We must be shutting down.
		return nil, true
	}
	elem := q.queue.Remove()
	if msg, ok := elem.(*WSMessage); ok {
		return msg, false
	}
	return nil, false
}

// ShutDown will cause q to ignore all new items added to it. As soon as the
// worker goroutines have drained the existing items in the queue, they will be
// instructed to exit.
func (q *MessageQueue) ShutDown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.shuttingDown = true
	q.cond.Broadcast()
}

func (q *MessageQueue) ShuttingDown() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	return q.shuttingDown
}
