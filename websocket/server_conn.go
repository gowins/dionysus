package websocket

import (
	"io/ioutil"
	"net"
	"time"
)

var KeepAliveInterval = 3 * time.Second
var LastPongTimeout = 3 * KeepAliveInterval

type ServerConn struct {
	*Conn
	Path        string
	IsClosed    bool
	ReadTimeout time.Duration
}

type AddrConns map[string]*ServerConn

func NewWSConn(conn net.Conn) *ServerConn {
	c := NewServerConn(conn)
	wsc := &ServerConn{
		Conn:     c,
		IsClosed: false,
	}
	return wsc
}

func (s *ServerConn) KeepAlive() {
	ticker := time.NewTicker(KeepAliveInterval)
	defer ticker.Stop()
	failedTimes := 0
	firstPing := true
	var checkTime time.Time
	for range ticker.C {
		if s.IsClosed {
			break
		}

		// err checks whether directly connected peer is lost
		err := s.Write(OpPing, nil)
		if err != nil {
			failedTimes++
		} else {
			if failedTimes > 0 {
				failedTimes--
			}
			if firstPing {
				firstPing = false
				checkTime = time.Now()
			}
		}

		if failedTimes >= 3 {
			break
		}

		// LastPong checks whether end peer  is lost
		// if we never receive Pong, check how long since first ping sent
		if !s.LastPong.IsZero() {
			checkTime = s.LastPong
		}
		if time.Since(checkTime) > LastPongTimeout {
			break
		}
	}
	s.IsClosed = true
	s.CloseLocalConn()
}

func (s *ServerConn) Handle() {
	if s.OnOpen != nil {
		ctx := NewEventContext(s.Conn, nil, OpText)
		s.OnOpen(ctx)
	}
	s.readData(s.customHandle)
	for {
		if s.readErr != nil && s.isDataQ.Length() == 0 {
			s.CloseLocalConn()
			s.IsClosed = true
			s.IsCloseErr(s.readErr)
			return
		}
		itemMsg, shutdown := s.isDataQ.Get()
		if shutdown {
			wLog.Infof("serverConn isDataQ is closed %v", shutdown)
		}
		if s.OnMessage != nil {
			if itemMsg != nil {
				ctx := NewEventContext(s.Conn, itemMsg.Message, OpText)
				s.OnMessage(ctx)
			}
		}
	}
}

func (s *ServerConn) readData(h CustomControlHandler) {
	err := linkGrPool.Submit(func() {
		for {
			if s.ReadTimeout > 0 {
				err := s.setReadDeadline(s.ReadTimeout)
				if err != nil {
					s.isDataQ.ShutDown()
					s.readErr = err
					return
				}
			}
			hdr, err := s.Reader.NextFrame()
			if err != nil {
				s.isDataQ.ShutDown()
				s.readErr = err
				return
			}
			if hdr.OpCode.IsControl() && s.Reader.OnIntermediate != nil {
				err := s.controlHandler(hdr)
				if h != nil {
					err = h(hdr, err)
				}
				if err != nil {
					s.isDataQ.ShutDown()
					s.readErr = err
					return
				}
				continue
			}
			bts, err := ioutil.ReadAll(s.Reader)
			if err != nil {
				s.isDataQ.ShutDown()
				s.readErr = err
				return
			}
			if insertFlag := s.isDataQ.Add(&WSMessage{Message: bts}); !insertFlag {
				wLog.Infof("insert serverConn isDataQ fail %v", insertFlag)
			}
		}
	})
	if err != nil {
		wLog.Errorf("server readData error:%v", err)
	}
}

func (s *ServerConn) IsCloseErr(e error) {
	if IsCloseErr(e) {
		if s.OnClose != nil {
			ctx := s.GetErrorEventContext(e)
			s.OnClose(ctx)
		}
	} else {
		if s.OnError != nil {
			ctx := NewEventContext(s.Conn, []byte(e.Error()), OpText)
			s.OnError(ctx)
		}
	}
}
