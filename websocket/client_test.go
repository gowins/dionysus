package websocket

import (
	"net"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	Convey("Test Reconnect", t, func() {
		t.Parallel()
		addr := "127.0.0.1:9001"
		url := "ws://" + addr
		option := ClientOption{
			Url: url,
		}
		c, err := NewClient(option)
		So(err, ShouldNotEqual, nil)
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		c, err = NewClient(option)
		if err != nil {
			panic(err)
		}
		go c.Start()
		addr1 := c.GetLocalAddr()
		sc := <-ch
		s := NewServerConn(sc)
		So(addr1, ShouldEqual, s.GetRemoteAddr())
		StopServer(stop)
		c.CloseLocalConn()
		s.CloseLocalConn()
	})
}

func TestClient(t *testing.T) {
	Convey("Test Reconnect", t, func() {
		t.Parallel()
		addr := "127.0.0.1:9002"
		url := "ws://" + addr
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		option := ClientOption{
			Url: url,
		}
		c, err := NewClient(option)
		if err != nil {
			panic(err)
		}
		addr1 := c.GetLocalAddr()
		onOpen := false
		onMsg := false
		onClose := false
		onError := false
		c.AddEventListener(EventOpen, func(e *EventContext) {
			onOpen = true
		})
		c.AddEventListener(EventMessage, func(e *EventContext) {
			onMsg = true
		})
		c.AddEventListener(EventClose, func(e *EventContext) {
			onClose = true
			c, _ = NewClient(option)
			c.AddEventListener(EventError, func(e *EventContext) {
				onError = true
			})
			c.Start()
		})
		go c.Start()
		sc := <-ch
		s := NewServerConn(sc)
		time.Sleep(100 * time.Millisecond)
		So(onOpen, ShouldEqual, true)
		s.WriteText([]byte("hello"))
		time.Sleep(100 * time.Millisecond)
		So(onMsg, ShouldEqual, true)
		// Server send close to test OnClose
		s.Close("test")
		sc1 := <-ch
		s1 := NewServerConn(sc1)
		time.Sleep(500 * time.Millisecond)
		So(onClose, ShouldEqual, true)
		So(addr1, ShouldNotEqual, c.GetLocalAddr())
		// Stop server to test OnError
		s1.CloseLocalConn()
		time.Sleep(1 * time.Second)
		So(onError, ShouldEqual, true)
		c.CloseLocalConn()
		StopServer(stop)
	})
}

func TestClientIsDataQ(t *testing.T) {
	Convey("Tes t TestClientIsDataQ", t, func() {
		t.Parallel()
		addr := "127.0.0.1:9003"
		url := "ws://" + addr
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		option := ClientOption{
			Url: url,
		}
		c, err := NewClient(option)
		if err != nil {
			panic(err)
		}
		msgArr := []string{"hello", "hello1", "hello2"}
		expectMsgArr := make([]string, 0)
		c.AddEventListener(EventMessage, func(e *EventContext) {
			expectMsgArr = append(expectMsgArr, string(e.Msg))
		})
		sc := <-ch
		s := NewServerConn(sc)
		go c.Start()
		for _, v := range msgArr {
			s.WriteText([]byte(v))
		}
		time.Sleep(100 * time.Millisecond)
		So(msgArr, ShouldResemble, expectMsgArr)
		s.CloseLocalConn()
		StopServer(stop)
	})
}

func TestClientShutdown(t *testing.T) {
	Convey("Test TestClientShutdown", t, func() {
		t.Parallel()
		addr := "127.0.0.1:9004"
		url := "ws://" + addr
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		option := ClientOption{
			Url: url,
		}
		c, err := NewClient(option)
		if err != nil {
			panic(err)
		}
		msgArr := []string{"hello", "hello1", "hello2"}
		expectMsgArr := make([]string, 0)
		sc := <-ch
		s := NewWSConn(sc)
		c.AddEventListener(EventMessage, func(e *EventContext) {
			expectMsgArr = append(expectMsgArr, string(e.Msg))
		})
		go c.Start()
		for _, v := range msgArr {
			s.WriteText([]byte(v))
		}
		time.Sleep(time.Millisecond * 100)
		So(expectMsgArr, ShouldResemble, msgArr)
		go c.GraceShutDown()
		time.Sleep(time.Millisecond * 100)
		So(c.isDataQ.Length(), ShouldEqual, 0)
	})
}
