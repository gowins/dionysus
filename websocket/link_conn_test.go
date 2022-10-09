package websocket

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestLink(t *testing.T) {
	link := NewLink()
	if link.Length() != 0 {
		t.Errorf("want length is 0 when a link is created")
		return
	}
	conn, ok := link.Pop()
	if conn != nil || ok {
		t.Errorf("want conn is nil and ok is false when a link is created")
		return
	}
	testConns := []*ServerConn{{Path: "test1"}, {Path: "test2"}, {Path: "test3"}, {Path: "test4"}, {Path: "test5"}}
	for i, testConn := range testConns {
		link.Push(testConn)
		if link.Length() != i+1 {
			t.Errorf("want length: %v, get length: %v", i+1, link.length)
			return
		}
	}
	wantPaths := []string{"test5", "test4", "test3", "test2", "test1"}
	for i, path := range wantPaths {
		conn, ok := link.Pop()
		if !ok || conn.Path != path {
			t.Errorf("want ok is true and want path: %v, get path: %v", path, conn.Path)
			return
		}
		if link.Length() != len(wantPaths)-1-i {
			t.Errorf("want length: %v, get length: %v", len(wantPaths)-1-i, link.length)
		}
	}
}

type testNetConn struct {
	ConnName string
}

func (testConn testNetConn) Read(b []byte) (n int, err error) { return 0, nil }
func (testConn testNetConn) Write(b []byte) (n int, err error) {
	return 10, fmt.Errorf("%v", testConn.ConnName)
}
func (testConn testNetConn) Close() error                       { return nil }
func (testConn testNetConn) LocalAddr() net.Addr                { return testAddr{} }
func (testConn testNetConn) RemoteAddr() net.Addr               { return testAddr{} }
func (testConn testNetConn) SetDeadline(t time.Time) error      { return nil }
func (testConn testNetConn) SetReadDeadline(t time.Time) error  { return nil }
func (testConn testNetConn) SetWriteDeadline(t time.Time) error { return nil }

type testAddr struct {
}

func (testAddr testAddr) Network() string {
	return "test"
}

func (testAddr testAddr) String() string {
	return "test"
}

func TestLink_SendMessage(t *testing.T) {
	tests := []struct {
		name      string
		connCount int
	}{
		{
			name:      "0 conn",
			connCount: 0,
		},
		{
			name:      "1 conn",
			connCount: 1,
		},
		{
			name:      "10 conns",
			connCount: 10,
		},
		{
			name:      "100 conns",
			connCount: 100,
		},
		{
			name:      "1000 conns",
			connCount: 1000,
		},
		{
			name:      "10000 conns",
			connCount: 10000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLink := NewLink()
			for i := 0; i < tt.connCount; i++ {
				c := testNetConn{
					ConnName: fmt.Sprintf("Name%v", i),
				}
				wsconn := NewWSConn(c)
				testLink.Push(wsconn)
			}
			err := testLink.SendMessage(OpText, []byte("error"))
			if tt.connCount == 0 || tt.connCount > 2000 {
				return
			}
			if err == nil {
				t.Errorf("want err not nil")
				return
			}
			for i := 0; i < tt.connCount; i++ {
				if !strings.Contains(err.Error(), fmt.Sprintf("Name%v", i)) {
					t.Errorf("want get Name%v in error: %v", i, err)
					return
				}
			}
		})
	}
}

func TestGetGoroutineCounts(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   int
	}{
		{
			length: 0,
			want:   0,
		},
		{
			length: 3,
			want:   3,
		},
		{
			length: 13,
			want:   10,
		},
		{
			length: 400,
			want:   10,
		},
		{
			length: 1000,
			want:   20,
		},
		{
			length: 2000,
			want:   40,
		},
		{
			length: 10000,
			want:   100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGoroutineCounts(tt.length); got != tt.want {
				t.Errorf("GetGoroutineCounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
