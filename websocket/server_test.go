package websocket

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gobwas/ws"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer_GetAllConns(t *testing.T) {
	route1 := AddrConns{"127.0.0.1:8888": newTestServerConn("test1"), "127.0.0.1:8889": newTestServerConn("test2")}
	route2 := AddrConns{"127.0.0.1:8886": newTestServerConn("test3"), "127.0.0.1:8887": newTestServerConn("test4"),
		"127.0.0.1:8885": newTestServerConn("test5")}
	route3 := AddrConns{}
	tests := []struct {
		ConnInfo    string
		testServers *Server
		wantWsConns int
	}{
		{
			testServers: &Server{},
			wantWsConns: 0,
		},
		{
			testServers: &Server{routeConns: map[string]AddrConns{"route1": route1}},
			wantWsConns: 2,
		},
		{
			testServers: &Server{routeConns: map[string]AddrConns{"route1": route1, "route2": route2}},
			wantWsConns: 5,
		},
		{
			testServers: &Server{routeConns: map[string]AddrConns{"route3": route3, "route2": route2}},
			wantWsConns: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.ConnInfo, func(t *testing.T) {
			if gotWsConns := tt.testServers.GetAllConns(); len(gotWsConns) != tt.wantWsConns {
				t.Errorf("GetAllConns() = %v, want %v", len(gotWsConns), tt.wantWsConns)
			}
		})
	}
}

func TestHandleMessages(t *testing.T) {
	var re string
	var lock sync.Mutex
	testServer := New("", ServerConfig{})
	go testServer.HandleMessages()
	msgQueue := testServer.GetMessageQueue()
	for i := 0; i < 66; i++ {
		sc := &ServerConn{
			Conn: NewServerConn(testConn{ConnConnInfo: "test"}),
		}
		sc.ConnInfo = fmt.Sprintf("ConnInfo%d", i)
		sc.OnPush = func(e *EventContext) {
			lock.Lock()
			re += string(e.Msg)
			lock.Unlock()
		}
		route := fmt.Sprintf("test%d", i%2)
		if addrConns, ok := testServer.routeConns[route]; ok {
			addrConns[sc.ConnInfo] = sc
		} else {
			testServer.routeConns[route] = AddrConns{sc.ConnInfo: sc}
		}
	}

	for i := 0; i < 66; i++ {
		msg := fmt.Sprintf("message%d", i)
		msgQueue.Add(&WSMessage{
			ServerConnName: fmt.Sprintf("ConnInfo%d", i),
			Message:        []byte(msg),
		})
	}
	time.Sleep(time.Millisecond * 300)
	var wantString string
	for i := 0; i < 66; i++ {
		wantMSG := fmt.Sprintf("message%d", i)
		wantString += wantMSG
	}
	if wantString != re {

		t.Errorf("want string: %v, got string: %v", wantString, re)
		return
	}
	re = ""
	wantString = ""
	msg := "haha"
	msgQueue.Add(&WSMessage{
		Route:   "all",
		Message: []byte(msg),
		OpCode:  OpText,
	})
	time.Sleep(time.Millisecond * 500)
	for i := 0; i < 66; i++ {
		wantString += "haha"
	}
	if re != wantString {
		t.Errorf("want string: %v, got string: %v", len(wantString), len(re))
		return
	}
	re = ""
	wantString = ""
	msgQueue.Add(&WSMessage{
		Route:   "test0",
		Message: []byte(msg),
		OpCode:  OpText,
	})
	time.Sleep(time.Millisecond * 500)
	for i := 0; i < 33; i++ {
		wantString += "haha"
	}
	if re != wantString {
		t.Errorf("want string: %v, got string: %v", len(wantString), len(re))
		return
	}
	testServer.Stop()
}

type testConn struct {
	ConnConnInfo string
}

func (testConn testConn) Read(b []byte) (n int, err error) { return 0, nil }
func (testConn testConn) Write(b []byte) (n int, err error) {
	return 10, fmt.Errorf("%v", testConn.ConnConnInfo)
}
func (testConn testConn) Close() error                       { return nil }
func (testConn testConn) LocalAddr() net.Addr                { return testAddress{} }
func (testConn testConn) RemoteAddr() net.Addr               { return testAddress{} }
func (testConn testConn) SetDeadline(t time.Time) error      { return nil }
func (testConn testConn) SetReadDeadline(t time.Time) error  { return nil }
func (testConn testConn) SetWriteDeadline(t time.Time) error { return nil }

type testAddress struct {
}

func (testAddr testAddress) Network() string {
	return "test"
}

func (testAddr testAddress) String() string {
	return "test"
}

func TestServer_GetRouteConns(t *testing.T) {
	route1 := AddrConns{"127.0.0.1:8888": newTestServerConn("test1"), "127.0.0.1:8889": newTestServerConn("test2")}
	route2 := AddrConns{"127.0.0.1:8886": newTestServerConn("test3"), "127.0.0.1:8887": newTestServerConn("test4"),
		"127.0.0.1:8885": newTestServerConn("test5")}
	route3 := AddrConns{}
	s := &Server{routeConns: map[string]AddrConns{"route1": route1, "route2": route2, "route3": route3}}
	tests := []struct {
		ConnInfo string
		path     string
		want     int
	}{
		{
			path: "notexist",
			want: 0,
		},
		{
			path: "route2",
			want: 3,
		},
		{
			path: "route1",
			want: 2,
		},
		{
			path: "route3",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.ConnInfo, func(t *testing.T) {
			if got := s.GetRouteConns(tt.path); len(got) != tt.want {
				t.Errorf("GetRouteConns() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestServer_SendGlobalMessage(t *testing.T) {
	testConns := make([]*Conn, 5)
	for i := 0; i < 5; i++ {
		test := fmt.Sprintf("test%d", i+1)
		testConn := NewServerConn(testConn{ConnConnInfo: test})
		testConn.ConnInfo = test
		testConns[i] = testConn
	}

	route1 := AddrConns{"127.0.0.1:8888": &ServerConn{Conn: testConns[0]},
		"127.0.0.1:8889": &ServerConn{Conn: testConns[1]}}
	route2 := AddrConns{"127.0.0.1:8886": &ServerConn{Conn: testConns[2]},
		"127.0.0.1:8887": &ServerConn{Conn: testConns[3]},
		"127.0.0.1:8885": &ServerConn{Conn: testConns[4]}}
	tests := []struct {
		name        string
		testServers *Server
		wantErr     []string
	}{
		{
			testServers: &Server{routeConns: map[string]AddrConns{"route1": route1}},
			wantErr:     []string{"test1", "test2"},
		},
		{
			testServers: &Server{routeConns: map[string]AddrConns{"route1": route1, "route2": route2}},
			wantErr:     []string{"test1", "test2", "test3", "test4", "test5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testServers.SendGlobalMessage(OpText, []byte("error"))
			if err == nil {
				t.Errorf("SendGlobalMessage() error is not nil")
				return
			}
			for _, v := range tt.wantErr {
				if !strings.Contains(err.Error(), v) {
					t.Errorf("want get %v in error: %v", v, err)
					return
				}
			}
		})
	}
}

func TestServer_SendRouteMessage(t *testing.T) {
	testConns := make([]*Conn, 5)
	for i := 0; i < 5; i++ {
		test := fmt.Sprintf("test%d", i+1)
		testConn := NewServerConn(testConn{ConnConnInfo: test})
		testConn.ConnInfo = test
		testConns[i] = testConn
	}

	route1 := AddrConns{"127.0.0.1:8888": &ServerConn{Conn: testConns[0]},
		"127.0.0.1:8889": &ServerConn{Conn: testConns[1]}}
	route2 := AddrConns{"127.0.0.1:8886": &ServerConn{Conn: testConns[2]},
		"127.0.0.1:8887": &ServerConn{Conn: testConns[3]},
		"127.0.0.1:8885": &ServerConn{Conn: testConns[4]}}
	testServers := Server{routeConns: map[string]AddrConns{"route1": route1, "route2": route2}}
	tests := []struct {
		name      string
		testRoute string
		wantErr   []string
	}{
		{
			testRoute: "route1",
			wantErr:   []string{"test1", "test2"},
		},
		{
			testRoute: "route2",
			wantErr:   []string{"test3", "test4", "test5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testServers.SendRouteMessage(tt.testRoute, OpText, []byte("error"))
			if err == nil {
				t.Errorf("SendGlobalMessage() error is not nil")
				return
			}
			for _, v := range tt.wantErr {
				if !strings.Contains(err.Error(), v) {
					t.Errorf("want get %v in error: %v", v, err)
					return
				}
			}
		})
	}
}

func TestStart(t *testing.T) {
	Convey("Test start", t, func() {
		t.Parallel()
		end := false
		snil := New(":xxxx", ServerConfig{})
		snil.Start()
		end = true
		So(end, ShouldEqual, true)
		end = false
		serr := New(":9900", ServerConfig{Cert: "/tmp/xxx", Key: "/tmp/xxx"})
		serr.Start()
		end = true
		So(end, ShouldEqual, true)
		s := New("0.0.0.0:9901", ServerConfig{})
		clientAddr := ""
		recStr := ""
		errStr := ""
		var closeCode uint16
		s.RegisterEventHandler("/test", EventOpen, func(e *EventContext) {
			clientAddr = e.Conn.GetRemoteAddr()
		})
		s.RegisterEventHandler("/test", EventMessage, func(e *EventContext) {
			recStr = string(e.Msg)
		})
		s.RegisterEventHandler("/test", EventClose, func(e *EventContext) {
			closeCode = e.GetCloseCode()
		})
		s.RegisterEventHandler("/test", EventError, func(e *EventContext) {
			errStr = string(e.Msg)
		})
		s.RegisterOnheader(func(connId string, key, value []byte) (conninfo string, err error) {
			if string(key) == "Conn-Info" {
				conninfo = string(value)
			}
			return
		})
		go s.Start()
		time.Sleep(100 * time.Millisecond)
		So(snil.eventHandlers, ShouldNotEqual, nil)
		url := "ws://127.0.0.1:9901/test"
		option := ClientOption{
			Url:  url,
			Head: map[string][]string{"Conn-Info": {"c1"}},
		}
		c, err := NewClient(option)
		time.Sleep(100 * time.Millisecond)
		So(err, ShouldEqual, nil)
		So(clientAddr, ShouldEqual, c.GetLocalAddr())

		c.WriteText([]byte("hello"))
		time.Sleep(100 * time.Millisecond)
		So(recStr, ShouldEqual, "hello")

		So(len(s.routeConns["/test"]), ShouldEqual, 1)
		option2 := ClientOption{
			Url:  url,
			Head: map[string][]string{"Conn-Info": {"c2"}},
		}
		c2, err := NewClient(option2)
		time.Sleep(100 * time.Millisecond)
		So(len(s.routeConns["/test"]), ShouldEqual, 2)

		c2.Close("c1 close")
		time.Sleep(200 * time.Millisecond)
		So(closeCode, ShouldEqual, CloseNormal)

		So(len(s.routeConns["/test"]), ShouldEqual, 1)
		fmt.Println("errStr", errStr)
		c.Close("")
		time.Sleep(200 * time.Millisecond)
		So(len(s.routeConns["/test"]), ShouldEqual, 0)
	})
}

func TestServer_BuildConnGrpool(t *testing.T) {
	s := New(":9983", ServerConfig{})
	tests := []struct {
		name     string
		maxConn  int
		wantConn int
	}{
		{
			maxConn:  defaultMaxConn,
			wantConn: defaultMaxConn,
		},
		{
			maxConn:  1024,
			wantConn: defaultMinConn,
		},
		{
			maxConn:  100000,
			wantConn: 100000,
		},
		{
			maxConn:  -1,
			wantConn: defaultMinConn,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.BuildConnGrpool(tt.maxConn)
			if s.connGrpool.Cap() != tt.wantConn {
				t.Errorf("want: %v, get: %v", tt.wantConn, s.connGrpool.Cap())
				return
			}
		})
	}
}

func TestServerUpgrade(t *testing.T) {
	Convey("Test TestServerUpgrade", t, func() {
		t.Parallel()
		addr := "127.0.0.1:9301"
		server := New(addr, ServerConfig{})

		type ConnInfo struct {
			token             string
			conninfo_header   string
			conninfo_protocol string
		}
		cis := map[string]*ConnInfo{}
		server.RegisterOnheader(func(connID string, key, value []byte) (string, error) {
			if string(key) == "Conn-Info" {
				if cis[connID] == nil {
					cis[connID] = &ConnInfo{
						conninfo_header: string(value),
					}
				} else {
					cis[connID].conninfo_header = string(value)
				}
				return string(value), nil
			}
			return "", nil
		})
		server.RegisterOnProtocol(func(connID string, bytes []byte) (string, bool) {
			if cis[connID] == nil {
				cis[connID] = &ConnInfo{
					conninfo_protocol: string(bytes),
				}
			} else {
				cis[connID].conninfo_protocol = string(bytes)
			}
			return string(bytes), true
		})
		server.RegisterOnBeforeUpgrade(func(connID string) (string, error) {
			conninfo := ""
			if cis[connID] == nil {
				return "", nil
			} else {
				if cis[connID].conninfo_protocol != "" {
					cis[connID].token = cis[connID].conninfo_protocol
				} else if cis[connID].conninfo_header != "" {
					cis[connID].token = cis[connID].conninfo_header
				}
				conninfo = cis[connID].token
			}
			return conninfo, nil
		})
		go server.Start()
		time.Sleep(time.Second)
		dialer := ws.Dialer{
			Protocols: []string{"ConnInfo-protocol-1"},
			Header:    ws.HandshakeHeaderHTTP(map[string][]string{"Conn-Info": {"ConnInfo-header-1"}}),
		}
		conn1, _, _, err := dialer.Dial(context.Background(), "ws://"+addr)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		dialer = ws.Dialer{
			Protocols: []string{"ConnInfo-protocol-2"},
			Header:    ws.HandshakeHeaderHTTP(map[string][]string{"Conn-Info": {"ConnInfo-header-2"}}),
		}
		conn2, _, _, err := dialer.Dial(context.Background(), "ws://"+addr)
		if err != nil {
			t.Errorf("Connect failed.")
		}
		So(err, ShouldEqual, nil)
		for _, v := range cis {
			So(v.token, ShouldEqual, v.conninfo_protocol)
		}
		conn1.Close()
		conn2.Close()
		server.Stop()
	})

}

func TestServerReadTimeout(t *testing.T) {
	Convey("Test TestServerUpgrade", t, func() {
		t.Parallel()
		var msg string
		// server without read timeout
		addr := "127.0.0.1:9302"
		server := New(addr, ServerConfig{})
		server.RegisterEventHandler("/test", EventMessage, func(e *EventContext) {
			msg = string(e.Msg)
		})
		go server.Start()
		time.Sleep(100 * time.Millisecond)
		client, err := NewClient(ClientOption{Url: "ws://" + addr + "/test"})
		if err != nil {
			t.Errorf("Connect failed.")
		}
		go client.Start()
		client.WriteText([]byte("hello"))
		time.Sleep(100 * time.Millisecond)
		So(msg, ShouldEqual, "hello")
		So(len(server.routeConns["/test"]), ShouldEqual, 1)
		time.Sleep(2 * time.Second)
		So(len(server.routeConns["/test"]), ShouldEqual, 1)
		client.Close("")
		server.Stop()

		// server with timeout
		addr1 := "127.0.0.1:9303"
		server1 := New(addr1, ServerConfig{})
		server1.RegisterEventHandler("/test1", EventMessage, func(e *EventContext) {
			msg = string(e.Msg)
		})
		server1.SetReadTimeout(time.Second)
		go server1.Start()
		time.Sleep(100 * time.Millisecond)
		client1, err1 := NewClient(ClientOption{Url: "ws://" + addr1 + "/test1"})
		if err1 != nil {
			t.Errorf("Connect failed.")
		}
		go client1.Start()
		client1.WriteText([]byte("hello1"))
		time.Sleep(100 * time.Millisecond)
		So(msg, ShouldEqual, "hello1")
		So(len(server1.routeConns["/test1"]), ShouldEqual, 1)
		time.Sleep(time.Second)
		// connection read timeout
		So(len(server1.routeConns["/test1"]), ShouldEqual, 0)
		client1.Close("")
		server1.Stop()
	})
}

func newTestServerConn(connInfo string) *ServerConn {
	return &ServerConn{Conn: &Conn{ConnInfo: connInfo}}
}
