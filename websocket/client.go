package websocket

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/gobwas/ws"
)

type WSClient struct {
	*Conn
	Dialer ws.Dialer
}

type ClientOption struct {
	Url       string
	TLSConfig *tls.Config
	Head      map[string][]string
	Protocols []string
}

func NewClient(option ClientOption) (*WSClient, error) {
	dialer := ws.Dialer{Protocols: option.Protocols}
	if option.TLSConfig != nil {
		dialer.TLSConfig = option.TLSConfig
	}
	if len(option.Head) != 0 {
		dialer.Header = ws.HandshakeHeaderHTTP(option.Head)
	}
	conn, _, _, err := dialer.Dial(context.Background(), option.Url)

	if err != nil {
		return nil, err
	}
	c := NewClientConn(conn)
	cli := &WSClient{
		Conn:   c,
		Dialer: dialer,
	}
	return cli, nil
}

func (c *WSClient) Start() {
	if c.OnOpen != nil {
		ctx := NewEventContext(c.Conn, nil, OpText)
		c.OnOpen(ctx)
	}
	c.readData(c.customHandle)
	for {
		if c.readErr != nil && c.isDataQ.Length() == 0 {
			c.CloseLocalConn()
			c.IsCloseErr(c.readErr)
			return
		}
		itemMsg, shutdown := c.isDataQ.Get()
		if shutdown {
			wLog.Infof("clientConn isDataQ is closed %v", shutdown)
		}
		if c.OnMessage != nil {
			if itemMsg != nil {
				ctx := NewEventContext(c.Conn, itemMsg.Message, OpText)
				c.OnMessage(ctx)
			}
		}
	}
}

func (c *WSClient) IsCloseErr(e error) {
	if IsCloseErr(e) {
		if c.OnClose != nil {
			ctx := c.GetErrorEventContext(e)
			c.OnClose(ctx)
		}
	} else {
		if c.OnError != nil {
			ctx := NewEventContext(c.Conn, []byte(e.Error()), OpText)
			c.OnError(ctx)
		}
	}
}

// will shutDown by dataQueue work out
func (c *WSClient) GraceShutDown() {
	c.CloseLocalConn()
	c.isDataQ.ShutDown()
	for c.isDataQ.Length() != 0 {
		time.Sleep(time.Millisecond * 30)
	}
}
