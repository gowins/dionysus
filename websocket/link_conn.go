package websocket

import (
	"fmt"
	"sync"

	utilerrors "github.com/gowins/dionysus/errors"
	"github.com/gowins/dionysus/grpool"
)

var linkGrPool, _ = grpool.NewPool(300000)

type LinkWSConn struct {
	Wsconn *ServerConn
	Next   *LinkWSConn
}

type Link struct {
	lock   sync.Mutex
	length int
	Head   *LinkWSConn
}

func NewLink() *Link {
	return &Link{sync.Mutex{}, 0, nil}
}

func (link *Link) Push(wsconn *ServerConn) {
	link.lock.Lock()
	defer link.lock.Unlock()
	node := &LinkWSConn{Wsconn: wsconn}
	node.Next = link.Head
	link.Head = node
	link.length++
}

func (link *Link) Pop() (*ServerConn, bool) {
	link.lock.Lock()
	defer link.lock.Unlock()
	if link.Head == nil {
		return nil, false
	}
	wsconn := link.Head.Wsconn
	link.Head = link.Head.Next
	link.length--
	return wsconn, true
}

func (link *Link) Length() int {
	link.lock.Lock()
	defer link.lock.Unlock()
	return link.length
}

func (link *Link) SendMessage(code OpCode, msg []byte) error {
	length := link.Length()
	if length == 0 {
		return nil
	}
	var errlist []error
	errChan := make(chan error, length)
	// TODO make sure suitable value
	i := GetGoroutineCounts(length)
	for j := 0; j < i; j++ {
		err := linkGrPool.Submit(func() {
			sendMessage(code, msg, errChan, link)
		})
		if err != nil {
			wLog.Errorf("grpool Submit sendMessage error: %v", err)
			errChan <- fmt.Errorf("grpool Submit sendMessage error: %v", err)
		}
	}
	errCount := 0
	for err := range errChan {
		errCount++
		if err != nil {
			errlist = append(errlist, err)
		}
		if errCount == i {
			break
		}
	}
	close(errChan)
	return utilerrors.NewAggregate(errlist)
}

func (link *Link) SendPush(msg []byte, code OpCode) {
	if link.Length() == 0 {
		return
	}
	var wg sync.WaitGroup
	// TODO make sure suitable value
	i := GetGoroutineCounts(link.Length())
	for j := 0; j < i; j++ {
		wg.Add(1)
		err := linkGrPool.Submit(func() {
			sendPush(msg, link, code, &wg)
		})
		if err != nil {
			wLog.Errorf("grpool Submit sendPush error: %v", err)
			wg.Done()
		}
	}
	wg.Wait()
}

func sendPush(msg []byte, linkConn *Link, code OpCode, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		wsConn, ok := linkConn.Pop()
		if !ok {
			break
		}
		if wsConn.OnPush != nil {
			eventContext := NewEventContext(wsConn.Conn, msg, code)
			wsConn.OnPush(eventContext)
		}
	}
}

func sendMessage(code OpCode, msg []byte, errorChan chan error, linkConn *Link) {
	var errlist []error
	defer func() {
		aggregateError := utilerrors.NewAggregate(errlist)
		errorChan <- aggregateError
	}()
	for {
		wsConn, ok := linkConn.Pop()
		if !ok {
			break
		}
		err := wsConn.Write(code, msg)
		if err != nil {
			errlist = append(errlist, err)
		}
		if code == OpClose {
			wsConn.CloseLocalConn()
		}
	}
}

func GetGoroutineCounts(length int) int {
	switch {
	case length <= 10:
		return length
	case 10 < length && length <= 500:
		return 10
	case 500 < length && length < 5000:
		return length / 50
	default:
		return 100
	}
}
