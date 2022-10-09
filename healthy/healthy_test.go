package healthy

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/gowins/dionysus/healthy/proto"

	"google.golang.org/grpc"

	"github.com/gobwas/ws"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/gowins/dionysus/algs"
)

func TestHealthy(t *testing.T) {
	Convey("Test health check", t, func() {
		f, err := os.CreateTemp(os.TempDir(), "test.Health")
		So(err, ShouldBeNil)
		HealthyFile = f.Name()

		So(CheckHealthyStat(), ShouldBeError)

		Convey("Start write stat and try again", func() {
			go WriteHealthyStat()
			time.Sleep(CheckInterval + time.Second)
			So(CheckHealthyStat(), ShouldBeNil)
			time.Sleep(CheckInterval)
			So(CheckHealthyStat(), ShouldBeNil)
		})
	})
}

func TestHealth_Stat(t *testing.T) {
	Convey("Testing Stat", t, func() {
		Convey("normal case", func() {
			h := New()
			h.RegChecker(func() error { return nil })
			h.RegChecker(func() error { return nil })

			cs := []Checker{
				func() error { return nil },
				func() error { return nil },
				func() error { return nil },
			}
			h.RegChecker(cs...)
			So(h.Stat(), ShouldBeNil)
		})

		Convey("error case", func() {
			h := New()
			h.RegChecker(func() error { return nil })
			h.RegChecker(func() error { return errors.New("want err") })
			h.RegChecker(func() error { return nil })
			So(h.Stat(), ShouldResemble, errors.New("The 2'th checker err: want err "))
		})

		Convey("panic case", func() {
			h := New()
			h.RegChecker(func() error { return nil })
			h.RegChecker(func() error {
				panic("want panic")
			})
			h.RegChecker(func() error { return nil })

			So(h.Stat(), ShouldResemble, errors.New("Panic at Stat: want panic"))
		})

	})
}

func TestHealth_FileObserve(t *testing.T) {
	path := filepath.Join(os.TempDir(), algs.RandStr(10, false))
	Convey("Testing Stat", t, func() {
		Convey("Test negative duration", func() {
			h := New()
			So(h.FileObserve(-1*time.Second, path), ShouldBeError)
		})

		Convey("Test err path", func() {
			h := New()
			So(h.FileObserve(time.Second, ""), ShouldBeError)
		})

		Convey("Test Check", func() {
			HealthyFile = path
			So(CheckHealthyStat(), ShouldBeError)

			h := New()
			So(h.FileObserve(time.Second, path), ShouldBeNil)
			time.Sleep(time.Millisecond)
			So(CheckHealthyStat(), ShouldBeNil)
		})
	})
}

func TestGetWSHealthURL(t *testing.T) {
	Convey("Testing GetWSHealthURL", t, func() {
		// if not write
		os.Remove(WSPortFile)
		err := GetWSHealthURL()
		So(err, ShouldNotEqual, nil)
		So(WSHealthUrl, ShouldEqual, "ws://127.0.0.1:9999"+WSHealthPath)

		//write
		WritePortFile(":8777", WSPortFile)
		err = GetWSHealthURL()
		So(err, ShouldEqual, nil)
		So(WSHealthUrl, ShouldEqual, "ws://127.0.0.1:8777"+WSHealthPath)
		os.Remove(WSPortFile)
	})
}

func TestWriteWSPortFile(t *testing.T) {
	Convey("Test WritePortFile", t, func() {
		// normal port
		os.Remove(WSPortFile)
		err := WritePortFile(":6777", WSPortFile)
		So(err, ShouldEqual, nil)
		f, err := os.Open(WSPortFile)
		if err != nil {
			t.Error("Can not open port file")
		}

		port, err := ioutil.ReadAll(f)
		if err != nil {
			t.Error("Can not read port file")
		}
		So(string(port), ShouldEqual, "6777")

		err = WritePortFile("6888", WSPortFile)
		f, err = os.Open(WSPortFile)
		if err != nil {
			t.Error("Can not open port file")
		}

		port, err = ioutil.ReadAll(f)
		if err != nil {
			t.Error("Can not read port file")
		}
		So(string(port), ShouldEqual, "6888")
		os.Remove(WSPortFile)
	})
}

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
		ln.Close()
		return nil
	}
}

func StopServer(stop chan struct{}) {
	stop <- struct{}{}
}

func TestCheckWSHealthy(t *testing.T) {
	Convey("Test CheckWSHealthy", t, func() {
		t.Parallel()
		os.Remove(WSPortFile)
		err := CheckWSHealthy()
		So(err, ShouldNotEqual, nil)

		conns := make(chan net.Conn, 1)
		stop := make(chan struct{})
		go TServer(conns, stop, ":7778")
		WritePortFile(":7778", WSPortFile)
		GetWSHealthURL()
		time.Sleep(100 * time.Millisecond)
		err = CheckWSHealthy()
		So(err, ShouldEqual, nil)
		sc := <-conns
		sc.Close()
		StopServer(stop)
	})
}

func TestGrpcHealthy(t *testing.T) {
	err := WritePortFile(":19303", GrpcPortFile)
	if err != nil {
		t.Errorf("want err nil, get err: %v", err)
		return
	}
	err = SetGrpcPort()
	if err != nil {
		t.Errorf("want err nil, get err: %v", err)
		return
	}
	testGrpcServerStart()
	err = CheckGrpcHealthy()
	if err != nil {
		t.Errorf("want err nil, get err: %v", err)
		return
	}
}

func testGrpcServerStart() {
	testServer := grpc.NewServer()
	pb.RegisterHealthServer(testServer, &HealthGrpc{})
	go func() {
		lis, err := net.Listen("tcp", GrpcAddr)
		if err != nil {
			log.Fatalf("[error] net listen: addr %v, err: %v ", GrpcAddr, err)
		}
		if err := testServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	time.Sleep(time.Second * 1)
}
