package websocket

import (
	"io"
	"io/ioutil"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var FrameSize = 1024

var defaultDialer = ws.Dialer{}

type CustomControlHandler func(hdr ws.Header, err error) error

func IsCloseErr(err error) bool {
	_, ok := err.(wsutil.ClosedError)
	return ok
}

func Read(c *Conn) ([]byte, byte, error) {
	return readData(c.Reader, c.customHandle)
}

func writeData(w io.Writer, op ws.OpCode, msg []byte, s ws.State) error {
	if op.IsControl() {
		if len(msg) > ws.MaxControlFramePayloadSize {
			return ws.ErrProtocolControlPayloadOverflow
		}
	}
	if op == ws.OpText || op == ws.OpBinary {
		if len(msg) > FrameSize {
			return writeFragmentedData(w, op, msg, s)
		}
	}
	return wsutil.WriteMessage(w, s, op, msg)
}

func writeFragmentedData(w io.Writer, op ws.OpCode, msg []byte, s ws.State) error {
	var i = 1
	var remaining = len(msg)
	for remaining > 0 {
		start := (i - 1) * FrameSize
		end := i * FrameSize
		if end >= len(msg) {
			end = len(msg)
		}

		frame := msg[start:end]
		if i == 1 {
			// first fragment frame
			err := writeFirstFragmentedFrame(w, op, frame, s)
			if err != nil {
				return err
			}
		} else {
			if remaining > FrameSize {
				// continue frame
				err := writeContinueFragmentedFrame(w, frame, s)
				if err != nil {
					return err
				}
			} else {
				// final frame
				err := writeLastFragmentedFrame(w, frame, s)
				if err != nil {
					return err
				}
				break
			}
		}
		remaining -= FrameSize
		i++
	}
	return nil
}

func writeFirstFragmentedFrame(w io.Writer, op ws.OpCode, msg []byte, s ws.State) error {
	f := ws.NewFrame(op, false, msg)
	return writeFrame(w, f, s)
}

func writeContinueFragmentedFrame(w io.Writer, msg []byte, s ws.State) error {
	f := ws.NewFrame(0, false, msg)
	return writeFrame(w, f, s)
}

func writeLastFragmentedFrame(w io.Writer, msg []byte, s ws.State) error {
	f := ws.NewFrame(0, true, msg)
	return writeFrame(w, f, s)
}

func writeFrame(w io.Writer, f ws.Frame, s ws.State) error {
	if s == ws.StateClientSide {
		f = ws.MaskFrameInPlace(f)
	}
	return ws.WriteFrame(w, f)
}

func readData(r *wsutil.Reader, h CustomControlHandler) ([]byte, byte, error) {
	controlHandler := r.OnIntermediate
	for {
		hdr, err := r.NextFrame()
		if err != nil {
			return nil, 0, err
		}
		if hdr.OpCode.IsControl() && controlHandler != nil {
			err := controlHandler(hdr, r)
			if h != nil {
				err = h(hdr, err)
			}
			if err != nil {
				return nil, 0, err
			}
			continue
		}
		bts, err := ioutil.ReadAll(r)
		return bts, byte(hdr.OpCode), err
	}
}

func writeClose(w io.Writer, statuscode ws.StatusCode, reason string, s ws.State) error {
	f := ws.NewCloseFrame(ws.NewCloseFrameBody(
		statuscode, reason,
	))
	return writeFrame(w, f, s)
}

func closeNormal(w io.Writer, reason string, s ws.State) error {
	return writeClose(w, ws.StatusNormalClosure, reason, s)
}
