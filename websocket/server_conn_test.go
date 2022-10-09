package websocket

import (
	"context"
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

/*func TestKeepAlive(t *testing.T) {
	Convey("Test KeepAlive", t, func() {
		t.Parallel()
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		addr := "127.0.0.1:9999"
		url := "ws://" + addr
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		conn, _, _, err := defaultDialer.Dial(context.Background(), url)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		c := NewClientConn(conn)
		sc := <-ch
		s := NewWSConn(sc)
		go s.KeepAlive()
		So(s.IsClosed, ShouldEqual, false)
		time.Sleep(4 * time.Second)
		f, _ := ws.ReadFrame(c.netConn)
		So(f.Header.OpCode, ShouldEqual, ws.OpPing)
		c.CloseLocalConn()
		time.Sleep(1 + 4*KeepAliveInterval)
		So(s.IsClosed, ShouldEqual, true)
		s.CloseLocalConn()

		// test if client to return pong and not close
		conn1, _, _, err := defaultDialer.Dial(context.Background(), url)
		sc1 := <-ch
		s1 := NewWSConn(sc1)
		go s1.KeepAlive()
		time.Sleep(1 + 4*KeepAliveInterval)
		So(s.IsClosed, ShouldEqual, true)
		conn1.Close()

		StopServer(stop)
	})
}*/

func TestHandle(t *testing.T) {
	Convey("Test Handle", t, func() {
		t.Parallel()
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		addr := "127.0.0.1:9202"
		url := "ws://" + addr
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		conn, _, _, err := defaultDialer.Dial(context.Background(), url)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		c := NewClientConn(conn)
		sc := <-ch
		s := NewWSConn(sc)
		onOpen := false
		onMsg := false
		onClose := false
		recStr := ""
		s.AddEventListener(EventOpen, func(e *EventContext) {
			onOpen = true
		})
		s.AddEventListener(EventMessage, func(e *EventContext) {
			onMsg = true
			recStr = string(e.Msg)
		})
		s.AddEventListener(EventClose, func(e *EventContext) {
			onClose = true
		})
		go s.Handle()
		time.Sleep(100 * time.Millisecond)
		So(onOpen, ShouldEqual, true)
		c.WriteText([]byte("hello"))
		time.Sleep(100 * time.Millisecond)
		So(onMsg, ShouldEqual, true)
		So(recStr, ShouldEqual, "hello")
		// client send close to server
		c.Close("test")
		time.Sleep(100 * time.Millisecond)
		So(onClose, ShouldEqual, true)
		s.CloseLocalConn()
		c.CloseLocalConn()
		StopServer(stop)
	})
}

func (s *ServerConn) testHandle() {
	if s.OnOpen != nil {
		ctx := NewEventContext(s.Conn, nil, OpText)
		s.OnOpen(ctx)
	}
	if s.ReadTimeout > 0 {
		err := s.setReadDeadline(s.ReadTimeout)
		if err != nil {
			s.CloseLocalConn()
			s.IsClosed = true
			s.IsCloseErr(err)
			return
		}
	}
	s.readData(s.customHandle)
	// concurrent serverConn
	time.Sleep(time.Second)
	for {
		if s.readErr != nil && s.isDataQ.Length() == 0 {
			s.CloseLocalConn()
			s.IsClosed = true
			s.IsCloseErr(s.readErr)
			return
		}
		if s.OnMessage != nil {
			itemMsg, _ := s.isDataQ.Get()
			if itemMsg != nil {
				ctx := NewEventContext(s.Conn, itemMsg.Message, OpText)
				s.OnMessage(ctx)
			}
		}
	}
}

func TestOnErrorHandle(t *testing.T) {
	Convey("Test OnErrorHandle", t, func() {
		t.Parallel()
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		addr := "127.0.0.1:9203"
		url := "ws://" + addr
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		conn, _, _, err := defaultDialer.Dial(context.Background(), url)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		c := NewClientConn(conn)
		sc := <-ch
		s := NewWSConn(sc)
		onError := false
		s.AddEventListener(EventError, func(e *EventContext) {
			onError = true
		})
		// client send close to server
		go s.testHandle()
		c.Close("test")
		time.Sleep(100 * time.Millisecond)
		// read a broken connection
		s.testHandle()
		So(onError, ShouldEqual, true)
		s.CloseLocalConn()
		c.CloseLocalConn()
		StopServer(stop)
	})
}
