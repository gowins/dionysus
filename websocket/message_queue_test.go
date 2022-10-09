package websocket

import (
	"fmt"
	"testing"
	"time"
)

func TestMessageQueue(t *testing.T) {
	msgQueue := NewMessageQueue()
	if msgQueue.Length() != 0 {
		t.Errorf("want length 0, get %v", msgQueue.Length())
		return
	}
	go func() {
		for i := 0; i < 128; i++ {
			msgQueue.Add(&WSMessage{
				ServerConnName: fmt.Sprintf("ServerConnName%d", i),
			})
		}
	}()
	for i := 0; i < 128; i++ {
		wsMSG, stoped := msgQueue.Get()
		if wsMSG == nil || stoped {
			t.Errorf("want wsMSG is not nil, stop is false")
			return
		}
		if wsMSG.ServerConnName != fmt.Sprintf("ServerConnName%d", i) {
			t.Errorf("want ServerConnName%d", i)
			return
		}
	}
	go msgQueue.ShutDown()
	time.Sleep(time.Millisecond * 100)
	if !msgQueue.ShuttingDown() {
		t.Errorf("want ShuttingDown is true")
		return
	}
	ok := msgQueue.Add(&WSMessage{
		ServerConnName: "test",
	})
	if ok {
		t.Errorf("want msgQueue Add false")
		return
	}
	for {
		wsMSG, stoped := msgQueue.Get()
		if wsMSG != nil || !stoped {
			t.Errorf("want wsMSG is  nil, stop is true")
		} else {
			break
		}
	}
}
