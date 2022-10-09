package healthy

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "github.com/gowins/dionysus/healthy/proto"
	"github.com/gowins/dionysus/websocket"
)

const CheckInterval = time.Second * 5
const CheckIntervalTimeOut = CheckInterval * 3

var (
	HealthyFile = filepath.Join(os.TempDir(), "dio.healthy")

	// WSPortFile Variables for websocket health check
	WSPortFile = filepath.Join(os.TempDir(), "ws.port")
	// GrpcPortFile Variables for grpc health check
	GrpcPortFile = filepath.Join(os.TempDir(), "grpc.port")
	WSHealthPath = "/healthz"
	WSHealthUrl  = "ws://127.0.0.1:9999" + WSHealthPath
	GrpcAddr     = ":8080"
)

type Checker func() error

type Health struct {
	checkerChain []Checker
}

func New() *Health {
	return &Health{}
}

func (h *Health) RegChecker(c ...Checker) {
	h.checkerChain = append(h.checkerChain, c...)
}

func (h *Health) FileObserve(d time.Duration, path string) error {
	if d <= 0 {
		return fmt.Errorf("duration should >= 0, got: %s ", d)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Write Health file err: %w ", err)
	}

	go func() {
		t := time.NewTicker(d)
		for {
			if h.Stat() == nil {
				_, err := f.WriteAt([]byte(strconv.Itoa(int(time.Now().Unix()))), 0)
				if err != nil {
					log.Errorf("Write Health stat to file:%s error:%s", f.Name(), err)
				}
			}

			<-t.C
		}
	}()

	return nil
}

func (h *Health) Stat() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic at Stat: %s", r)
			log.Error(err)
		}
	}()

	for i, cc := range h.checkerChain {
		if e := cc(); e != nil {
			err = fmt.Errorf("The %d'th checker err: %s ", i+1, e)
			return
		}
	}

	return
}

func CheckHealthyStat() error {
	f, err := os.Open(HealthyFile)
	if err != nil {
		return fmt.Errorf("Open file err: %w ", err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("Read file err: %w ", err)
	}

	ti, err := strconv.Atoi(string(b))
	if err != nil {
		return fmt.Errorf("Parse time err: %w ", err)
	}

	t := time.Unix(int64(ti), 0)
	if time.Since(t) > CheckIntervalTimeOut {
		return fmt.Errorf("Service did not update healthy file more than %s. Last time updated at: %s ", CheckIntervalTimeOut.String(), t.String())
	}
	return nil
}

// Deprecated:: This func will delete after v1.0
func WriteHealthyStat() {
	f, err := os.Create(HealthyFile)
	if err != nil {
		log.Fatal(fmt.Errorf("Write healthy file err: %w ", err))
	}
	defer f.Close()

	t := time.NewTicker(CheckInterval)

	for now := range t.C {
		_, _ = f.WriteAt([]byte(strconv.Itoa(int(now.Unix()))), 0)
	}
}

func GetWSHealthURL() error {
	f, err := os.Open(WSPortFile)
	if err != nil {
		return err
	}

	port, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	WSHealthUrl = "ws://127.0.0.1:" + string(port) + WSHealthPath
	return nil
}

func WritePortFile(addr, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	addrs := strings.Split(addr, ":")
	port := addrs[len(addrs)-1]
	_, err = f.Write([]byte(port))
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func CheckWSHealthy() error {
	cli, err := websocket.NewClient(websocket.ClientOption{Url: WSHealthUrl})
	if err != nil {
		return err
	}
	cli.CloseLocalConn()
	return nil
}

func SetGrpcPort() error {
	f, err := os.Open(GrpcPortFile)
	if err != nil {
		return err
	}

	port, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	GrpcAddr = "127.0.0.1:" + string(port)
	return nil
}

func CheckGrpcHealthy() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeoutCtx, GrpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewHealthClient(conn)
	r, err := c.SayHealth(context.Background(), &pb.HealthRequest{Name: "health"})
	if err != nil {
		return err
	}
	if r.Message != "healthchecked" {
		return fmt.Errorf("get message: %v, want healthchecked", r.Message)
	}
	return nil
}

type HealthGrpc struct {
}

func (h *HealthGrpc) SayHealth(ctx context.Context, in *pb.HealthRequest) (*pb.HealthReply, error) {
	return &pb.HealthReply{Message: in.GetName() + "checked"}, nil
}
