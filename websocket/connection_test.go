package websocket

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gobwas/ws/wsutil"

	"github.com/gobwas/ws"

	. "github.com/smartystreets/goconvey/convey"
)

var stop = make(chan struct{})

func TServer(ch chan net.Conn, stop chan struct{}, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}
			u := ws.Upgrader{}

			_, err = u.Upgrade(conn)
			if err != nil {
				break
			}
			ch <- conn
		}
	}()

	select {
	case <-stop:
		fmt.Println("stop server :", addr)
		ln.Close()
		return nil
	}
}

func StopServer(stop chan struct{}) {
	stop <- struct{}{}
}

func TestReadWrite(t *testing.T) {
	Convey("Test Read", t, func() {
		t.Parallel()
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		addr := "127.0.0.1:9101"
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		url := "ws://" + addr
		conn, _, _, err := defaultDialer.Dial(context.Background(), url)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		c := NewClientConn(conn)
		sc := <-ch
		s := NewServerConn(sc)
		c.Write(OpText, []byte("hello"))
		s.readData(s.customHandle)
		item, shutdown := s.isDataQ.Get()
		time.Sleep(time.Millisecond * 100)
		So(s.readErr, ShouldEqual, nil)
		So(shutdown, ShouldEqual, false)
		So(string(item.Message), ShouldEqual, "hello")
		StopServer(stop)
		s.CloseLocalConn()
		s.CloseLocalConn()
	})
}

func TestWriteText(t *testing.T) {
	Convey("Test WriteText", t, func() {
		t.Parallel()
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		addr := "127.0.0.1:9102"
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		url := "ws://" + addr
		conn, _, _, err := defaultDialer.Dial(context.Background(), url)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		c := NewClientConn(conn)
		sc := <-ch
		s := NewServerConn(sc)
		c.WriteText([]byte("hello"))
		s.readData(s.customHandle)
		item, shutdown := s.isDataQ.Get()
		time.Sleep(time.Millisecond * 100)
		So(s.readErr, ShouldEqual, nil)
		So(shutdown, ShouldEqual, false)
		So(string(item.Message), ShouldEqual, "hello")
		StopServer(stop)
		s.CloseLocalConn()
		s.CloseLocalConn()
	})
}

func TestWriteBinary(t *testing.T) {
	Convey("Test WriteBinary", t, func() {
		t.Parallel()
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		addr := "127.0.0.1:9103"
		url := "ws://" + addr
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		conn, _, _, err := defaultDialer.Dial(context.Background(), url)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		c := NewClientConn(conn)
		sc := <-ch
		s := NewServerConn(sc)
		c.WriteBinary([]byte("hello"))
		s.readData(s.customHandle)
		item, shutdown := s.isDataQ.Get()
		time.Sleep(time.Millisecond * 100)
		So(s.readErr, ShouldEqual, nil)
		So(shutdown, ShouldEqual, false)
		So(string(item.Message), ShouldEqual, "hello")
		StopServer(stop)
		c.CloseLocalConn()
		s.CloseLocalConn()
	})
}

func TestConnClose(t *testing.T) {
	Convey("Test WriteBinary", t, func() {
		t.Parallel()
		ch := make(chan net.Conn, 1)
		stop := make(chan struct{})
		addr := "127.0.0.1:9104"
		url := "ws://" + addr
		go TServer(ch, stop, addr)
		time.Sleep(100 * time.Millisecond)
		conn, _, _, err := defaultDialer.Dial(context.Background(), url)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		c := NewClientConn(conn)
		sc := <-ch
		s := NewServerConn(sc)
		s.Close("hello")
		c.readData(c.customHandle)
		item, shutdown := c.isDataQ.Get()
		time.Sleep(time.Millisecond * 100)
		So(IsCloseErr(c.readErr), ShouldEqual, true)
		So(item, ShouldEqual, nil)
		So(shutdown, ShouldEqual, true)
		StopServer(stop)
		c.CloseLocalConn()
		s.CloseLocalConn()
	})
}

func TestCloseContext(t *testing.T) {
	Convey("Test close EventContext", t, func() {
		s := NewServerConn(nil)
		closeCode := ws.StatusNormalClosure
		closeErr := wsutil.ClosedError{Code: closeCode}
		e := s.GetErrorEventContext(closeErr)
		So(e.GetCloseCode(), ShouldEqual, closeCode)
	})
}
