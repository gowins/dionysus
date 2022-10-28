package healthy

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	pb "github.com/gowins/dionysus/healthy/proto"
	"github.com/gowins/dionysus/log"
)

var healthHttpClient = http.Client{Timeout: 30 * time.Second}

const (
	HealthGroupPath      = "/healthx"
	HealthLivenessPath   = "/liveness"
	HealthReadinessPath  = "/readiness"
	HealthStartupPath    = "/startup"
	HealthLiveness       = "liveness"
	HealthReadiness      = "readiness"
	HealthStartup        = "startup"
	CheckInterval        = time.Second * 5
	CheckIntervalTimeOut = CheckInterval * 3
	HealthStatus         = "HEALTH_STATUS"
	StatusOpen           = "open"
	StatusClose          = "close"
)

var (
	HealthyFile       = filepath.Join(os.TempDir(), "dio.healthy")
	livenessCheckers  = []Checker{}
	readinessCheckers = []Checker{}
	startupCheckers   = []Checker{}
)

type Checker func() error

type Health struct {
}

func New() *Health {
	return &Health{}
}

func RegReadinessCheckers(checkers ...Checker) {
	readinessCheckers = append(readinessCheckers, checkers...)
}

func RegLivenessCheckers(checkers ...Checker) {
	livenessCheckers = append(livenessCheckers, checkers...)
}

func RegStartupCheckers(checkers ...Checker) {
	startupCheckers = append(startupCheckers, checkers...)
}

func (h *Health) FileObserve(d time.Duration) error {
	if d <= 0 {
		return fmt.Errorf("duration should >= 0, got: %s ", d)
	}

	f, err := os.Create(HealthyFile)
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

	// ctl only support healthLiveness
	return CheckerFuncRun(HealthLiveness)
}

func CheckCtlHealthyStat(healthType string) error {
	if err := CheckerFuncRun(healthType); err != nil {
		return err
	}
	f, err := os.Open(HealthyFile)
	if err != nil {
		return fmt.Errorf("Open file err: %w ", err)
	}

	b, err := ioutil.ReadAll(f)
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

func CheckerFuncRun(checkType string) error {
	switch checkType {
	case HealthLiveness:
		for i, fn := range livenessCheckers {
			if e := fn(); e != nil {
				return fmt.Errorf("The %d'th checker err: %s ", i+1, e)
			}
		}
	case HealthReadiness:
		for i, fn := range readinessCheckers {
			if e := fn(); e != nil {
				return fmt.Errorf("The %d'th checker err: %s ", i+1, e)
			}
		}
	case HealthStartup:
		for i, fn := range startupCheckers {
			if e := fn(); e != nil {
				return fmt.Errorf("The %d'th checker err: %s ", i+1, e)
			}
		}
	}
	return nil
}

func CheckHttpHealthyStat(url string, checkType string) error {
	if err := CheckerFuncRun(checkType); err != nil {
		return err
	}
	url = "http://127.0.0.1" + url
	res, err := healthHttpClient.Get(url)
	if err != nil || res == nil {
		return fmt.Errorf("health check url %v, error %v", url, err)
	}
	if res == nil {
		return fmt.Errorf("health check url %v res nil", url)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("health check url %v, res statusCode %v", url, res.StatusCode)
	}
	return nil
}

func SetHttpHealthyOpen(url string, checkType string) error {
	return SetHttpHealthy(url, checkType, StatusOpen)
}

func SetHttpHealthyClose(url string, checkType string) error {
	return SetHttpHealthy(url, checkType, StatusClose)
}

func SetHttpHealthy(url string, checkType string, status string) error {
	if err := CheckerFuncRun(checkType); err != nil {
		return err
	}
	url = "http://127.0.0.1" + url + "/" + status
	res, err := healthHttpClient.Post(url, "application/json", nil)
	if err != nil || res == nil {
		return fmt.Errorf("health check url %v, error %v", url, err)
	}
	if res == nil {
		return fmt.Errorf("health check url %v res nil", url)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("health check url %v, res statusCode %v", url, res.StatusCode)
	}
	return nil
}

func CheckGrpcHealthy(grpcAddr, checkType string) error {
	grpcAddr = "127.0.0.1" + grpcAddr
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeoutCtx, grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewHealthServiceClient(conn)
	switch checkType {
	case HealthLiveness:
		return checkGrpcLiveness(c)
	case HealthReadiness:
		return checkGrpcReadiness(c)
	case HealthStartup:
		return checkGrpcStartup(c)
	}
	return nil
}

func SetGrpcHealthyOpen(grpcAddr string, checkType string) error {
	return SetGrpcHealthy(grpcAddr, checkType, true)
}

func SetGrpcHealthyClose(grpcAddr string, checkType string) error {
	return SetGrpcHealthy(grpcAddr, checkType, false)
}

func SetGrpcHealthy(grpcAddr string, checkType string, status bool) error {
	grpcAddr = "127.0.0.1" + grpcAddr
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeoutCtx, grpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewHealthServiceClient(conn)
	switch checkType {
	case HealthLiveness:
		return setGrpcLiveness(c, status)
	case HealthReadiness:
		return setGrpcReadiness(c, status)
	case HealthStartup:
		return setGrpcStartup(c, status)
	}
	return nil
}

func checkGrpcLiveness(c pb.HealthServiceClient) error {
	r, err := c.HealthLiveness(context.TODO(), &pb.HealthyRequest{})
	if err != nil {
		return err
	}
	if r.Response != HealthLiveness {
		return fmt.Errorf("get message: %v, want healthchecked", r.Response)
	}
	return nil
}

func checkGrpcReadiness(c pb.HealthServiceClient) error {
	r, err := c.HealthReadiness(context.TODO(), &pb.HealthyRequest{})
	if err != nil {
		return err
	}
	if r.Response != HealthReadiness {
		return fmt.Errorf("get message: %v, want healthchecked", r.Response)
	}
	return nil
}

func checkGrpcStartup(c pb.HealthServiceClient) error {
	r, err := c.HealthStartup(context.TODO(), &pb.HealthyRequest{})
	if err != nil {
		return err
	}
	if r.Response != HealthStartup {
		return fmt.Errorf("get message: %v, want healthchecked", r.Response)
	}
	return nil
}

func setGrpcLiveness(c pb.HealthServiceClient, status bool) error {
	_, err := c.HealthLivenessSet(context.TODO(), &pb.HealthyStatus{Status: status})
	return err
}

func setGrpcReadiness(c pb.HealthServiceClient, status bool) error {
	_, err := c.HealthReadinessSet(context.TODO(), &pb.HealthyStatus{Status: status})
	return err
}

func setGrpcStartup(c pb.HealthServiceClient, status bool) error {
	_, err := c.HealthStartupSet(context.TODO(), &pb.HealthyStatus{Status: status})
	return err
}

type GrpcHealthyServer struct {
	livenessStatus  bool
	readinessStatus bool
	startupStatus   bool
}

var grpcHealthyServer = &GrpcHealthyServer{
	livenessStatus:  true,
	readinessStatus: true,
	startupStatus:   true,
}

func GetGrpcHealthyServer() *GrpcHealthyServer {
	return grpcHealthyServer
}

func (h *GrpcHealthyServer) HealthLiveness(context.Context, *pb.HealthyRequest) (*pb.HealthyResponse, error) {
	if err := CheckerFuncRun(HealthLiveness); err != nil {
		return &pb.HealthyResponse{}, err
	}
	if !h.livenessStatus {
		return &pb.HealthyResponse{}, fmt.Errorf("livenessStatus is closed")
	}
	return &pb.HealthyResponse{
		Response: HealthLiveness,
	}, nil
}

func (h *GrpcHealthyServer) HealthLivenessSet(ctx context.Context, req *pb.HealthyStatus) (*pb.HealthyResponse, error) {
	h.livenessStatus = req.Status
	return &pb.HealthyResponse{
		Response: fmt.Sprintf("grpc liveness is set %v", req.Status),
	}, nil
}

func (h *GrpcHealthyServer) HealthReadiness(context.Context, *pb.HealthyRequest) (*pb.HealthyResponse, error) {
	if err := CheckerFuncRun(HealthReadiness); err != nil {
		return &pb.HealthyResponse{}, err
	}
	if !h.readinessStatus {
		return &pb.HealthyResponse{}, fmt.Errorf("readinessStatus is closed")
	}
	return &pb.HealthyResponse{
		Response: HealthReadiness,
	}, nil
}

func (h *GrpcHealthyServer) HealthReadinessSet(ctx context.Context, req *pb.HealthyStatus) (*pb.HealthyResponse, error) {
	h.readinessStatus = req.Status
	return &pb.HealthyResponse{
		Response: fmt.Sprintf("grpc readiness is set %v", req.Status),
	}, nil
}

func (h *GrpcHealthyServer) HealthStartup(context.Context, *pb.HealthyRequest) (*pb.HealthyResponse, error) {
	if err := CheckerFuncRun(HealthStartup); err != nil {
		return &pb.HealthyResponse{}, err
	}
	if !h.startupStatus {
		return &pb.HealthyResponse{}, fmt.Errorf("startupStatus is closed")
	}
	return &pb.HealthyResponse{
		Response: HealthStartup,
	}, nil
}

func (h *GrpcHealthyServer) HealthStartupSet(ctx context.Context, req *pb.HealthyStatus) (*pb.HealthyResponse, error) {
	h.startupStatus = req.Status
	return &pb.HealthyResponse{
		Response: fmt.Sprintf("grpc startup is set %v", req.Status),
	}, nil
}
