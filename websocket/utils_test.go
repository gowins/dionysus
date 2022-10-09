package websocket

import (
	"bytes"
	"io"
	"testing"

	"github.com/gobwas/ws/wsutil"

	"github.com/gobwas/ws"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReadData(t *testing.T) {
	Convey("Test ReadData", t, func() {
		s := ws.StateClientSide
		r := &wsutil.Reader{
			Source: bytes.NewReader(nil),
			State:  s,
		}
		exceptedErr := io.EOF
		_, _, err := readData(r, nil)

		So(err, ShouldEqual, exceptedErr)

		s = ws.StateServerSide
		f := ws.NewFrame(ws.OpText, true, nil)
		var rbuf bytes.Buffer
		ws.WriteFrame(&rbuf, f)
		r = &wsutil.Reader{
			Source: &rbuf,
			State:  s,
		}
		_, _, err = readData(r, nil)
		So(err, ShouldEqual, ws.ErrProtocolMaskRequired)
	})
}

func TestReadDataFragmented(t *testing.T) {
	Convey("Test ReadData Fragmented frame", t, func() {
		s := ws.StateServerSide
		var buf bytes.Buffer
		writeFirstFragmentedFrame(&buf, ws.OpText, []byte("fragment1"), s)
		writeContinueFragmentedFrame(&buf, []byte(","), s)
		writeLastFragmentedFrame(&buf, []byte("fragment2"), s)
		r := &wsutil.Reader{
			SkipHeaderCheck: true,
			Source:          bytes.NewReader(buf.Bytes()),
			State:           ws.StateClientSide,
		}
		expectedOp := ws.OpText
		expectedMsg := "fragment1,fragment2"
		msg, op, err := readData(r, nil)
		So(op, ShouldEqual, expectedOp)
		So(string(msg), ShouldEqual, expectedMsg)
		So(err, ShouldEqual, nil)
	})
}

func TestReadDataUnexpectedEOF(t *testing.T) {
	Convey("Test ReadData Only header written", t, func() {
		s := ws.StateClientSide
		var buf bytes.Buffer
		f := ws.NewTextFrame([]byte("this part will be lost"))
		if err := ws.WriteHeader(&buf, f.Header); err != nil {
			panic(err)
		}
		r := &wsutil.Reader{
			SkipHeaderCheck: true,
			Source:          bytes.NewReader(buf.Bytes()),
			State:           s,
		}
		expectedErr := io.ErrUnexpectedEOF
		_, _, err := readData(r, nil)
		So(err, ShouldEqual, expectedErr)
	})
}

func TestReadDataControlFrame(t *testing.T) {
	// If ReadData get control frame, the control frame will be handled by WSControlHandler,
	// then ReadData will continue read next frame, then get a EOF error
	Convey("Test ReadData Control Frame", t, func() {
		s := ws.StateClientSide
		var buf bytes.Buffer
		f := ws.NewFrame(ws.OpPing, true, nil)
		if err := ws.WriteFrame(&buf, f); err != nil {
			panic(err)
		}
		ch := wsutil.ControlFrameHandler(bytes.NewBuffer(buf.Bytes()), s)
		r := &wsutil.Reader{
			SkipHeaderCheck: true,
			Source:          bytes.NewReader(buf.Bytes()),
			State:           s,
			OnIntermediate:  ch,
		}
		msg, op, err := readData(r, nil)
		So(err, ShouldEqual, io.EOF)
		So(op, ShouldEqual, 0)
		So(msg, ShouldEqual, nil)
	})
}

func TestWriteData(t *testing.T) {
	Convey("Test writeData", t, func() {
		var buf bytes.Buffer
		err := writeData(&buf, ws.OpText, []byte("hello"), ws.StateClientSide)
		if err != nil {
			t.Errorf("writeData error")
		}
		f, err := ws.ReadFrame(&buf)
		So(f.Header.OpCode, ShouldEqual, ws.OpText)
		So(f.Header.Masked, ShouldEqual, true)
		So(f.Header.Length, ShouldEqual, 5)
		So(f.Header.Fin, ShouldEqual, true)
		msg := "hello"

		for i := 0; i < 100; i++ {
			msg = msg + "hello"

		}
		err = writeData(&buf, ws.OpPing, []byte(msg), ws.StateServerSide)
		So(err, ShouldEqual, ws.ErrProtocolControlPayloadOverflow)
	})
}

func TestWriteFragmentedData(t *testing.T) {
	Convey("Test WriteFragmentedData", t, func() {
		FrameSize = 10
		var buf bytes.Buffer
		var buf1 bytes.Buffer
		msg := []byte("helloworldhelloworldhelloworldhelloworldhelloworldhelloworldhelloworldaaaa")
		writeData(&buf, ws.OpText, msg, ws.StateServerSide)
		writeData(&buf1, ws.OpText, msg, ws.StateServerSide)

		i := 0
		textCount := 0
		continueCount := 0
		for {
			f, err := ws.ReadFrame(&buf)
			if err != nil {
				break
			}
			i++
			if f.Header.OpCode == ws.OpText {
				textCount++
			}
			if f.Header.OpCode == ws.OpContinuation {
				continueCount++
			}
		}
		So(textCount, ShouldEqual, 1)
		So(continueCount, ShouldEqual, 7)
		So(i, ShouldEqual, 8)

		r := &wsutil.Reader{
			SkipHeaderCheck: true,
			Source:          bytes.NewReader(buf1.Bytes()),
			State:           ws.StateClientSide,
		}
		ret, op, _ := readData(r, nil)
		So(op, ShouldEqual, ws.OpText)
		So(string(msg), ShouldEqual, string(ret))
	})
}

func TestWriteFirstFragmentedFrame(t *testing.T) {
	Convey("Test writeData", t, func() {
		var buf bytes.Buffer
		err := writeFirstFragmentedFrame(&buf, ws.OpBinary, []byte("hello"), ws.StateClientSide)
		if err != nil {
			t.Errorf("writeFirstFragmentedFrame error")
		}
		f, err := ws.ReadFrame(&buf)
		So(f.Header.Fin, ShouldEqual, false)
		So(f.Header.Masked, ShouldEqual, true)
		So(f.Header.OpCode, ShouldEqual, ws.OpBinary)
		So(f.Header.Length, ShouldEqual, len("hello"))
	})
}

func TestWriteContinueFragmentedFrame(t *testing.T) {
	Convey("Test writeData", t, func() {
		var buf bytes.Buffer
		err := writeContinueFragmentedFrame(&buf, []byte("hello"), ws.StateClientSide)
		if err != nil {
			t.Errorf("writeContinueFragmentedFrame error")
		}
		f, err := ws.ReadFrame(&buf)
		So(f.Header.Fin, ShouldEqual, false)
		So(f.Header.Masked, ShouldEqual, true)
		So(f.Header.OpCode, ShouldEqual, 0)
		So(f.Header.Length, ShouldEqual, len("hello"))
	})
}

func TestWriteLastFragmentedFrame(t *testing.T) {
	Convey("Test writeData", t, func() {
		var buf bytes.Buffer
		err := writeLastFragmentedFrame(&buf, []byte("hello"), ws.StateServerSide)
		if err != nil {
			t.Errorf("writeLastFragmentedFrame error")
		}
		f, err := ws.ReadFrame(&buf)
		So(f.Header.Fin, ShouldEqual, true)
		So(f.Header.Masked, ShouldEqual, false)
		So(f.Header.OpCode, ShouldEqual, 0)
		So(f.Header.Length, ShouldEqual, len("hello"))
	})
}

func TestWriteFrame(t *testing.T) {
	Convey("Test WriteFrame", t, func() {
		var buf bytes.Buffer
		f := ws.NewFrame(ws.OpText, true, nil)
		err := writeFrame(&buf, f, ws.StateServerSide)
		if err != nil {
			t.Errorf("writeFrame error")
		}
		f, err = ws.ReadFrame(&buf)
		So(f.Header.Masked, ShouldEqual, false)

		err = writeFrame(&buf, f, ws.StateClientSide)
		if err != nil {
			t.Errorf("writeFrame error")
		}
		f, err = ws.ReadFrame(&buf)
		So(f.Header.Masked, ShouldEqual, true)
	})
}

func TestClose(t *testing.T) {
	Convey("Test close", t, func() {
		var buf bytes.Buffer
		reason := "close"
		closeNormal(&buf, reason, ws.StateServerSide)
		f, _ := ws.ReadFrame(&buf)
		So(f.Header.OpCode, ShouldEqual, ws.OpClose)
		So(f.Header.Length, ShouldEqual, 2+len(reason))
	})
}
