package cmd

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gowins/dionysus/healthy"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/pflag"
)

var healthHttpClient = http.Client{Timeout: 3 * time.Second}

func TestNewGinCommand(t *testing.T) {
	defer func() {
		_ = recover()
	}()
	log.Setup(log.SetProjectName("gin_cmd"), log.WithWriter(os.Stdout), log.WithOnFatal(&log.MockCheckWriteHook{}))
	gcmd := NewGinCommand()
	gcmd.RegFlagSet(&pflag.FlagSet{})
	gcmd.Flags()
	gcmd.GetCmd()
	gcmd.RegShutdownFunc(StopStep{
		StepName: "stopgccm",
		StopFn: func() {
			fmt.Printf("this is stop gin")
		},
	})
	stopfn := gcmd.GetShutdownFunc()
	go func() {
		gcmd.registerHealth()
		gcmd.startServer()
		time.Sleep(time.Second * 10)
	}()
	time.Sleep(time.Second * 3)
	setHttpHealthyStatOpen(healthy.HealthLiveness)
	if err := checkHttpHealthyStat(healthy.HealthLiveness); err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	setHttpHealthyStatOpen(healthy.HealthReadiness)
	if err := checkHttpHealthyStat(healthy.HealthReadiness); err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	setHttpHealthyStatOpen(healthy.HealthStartup)
	if err := checkHttpHealthyStat(healthy.HealthStartup); err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	setHttpHealthyStatClose(healthy.HealthLiveness)
	if err := checkHttpHealthyStat(healthy.HealthLiveness); err == nil {
		t.Errorf("want error not nil")
		return
	}
	setHttpHealthyStatClose(healthy.HealthReadiness)
	if err := checkHttpHealthyStat(healthy.HealthReadiness); err == nil {
		t.Errorf("want error not nil")
		return
	}
	setHttpHealthyStatClose(healthy.HealthStartup)
	if err := checkHttpHealthyStat(healthy.HealthStartup); err == nil {
		t.Errorf("want error not nil")
		return
	}
	time.Sleep(time.Second * 8)
	stopfn()
}

func checkHttpHealthyStat(checkType string) error {
	url := "http://127.0.0.1" + defaultWebServerAddr + healthy.HealthGroupPath + "/" + checkType
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

func setHttpHealthyStatOpen(checkType string) error {
	url := "http://127.0.0.1" + defaultWebServerAddr + healthy.HealthGroupPath + "/" + checkType + "/open"
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

func setHttpHealthyStatClose(checkType string) error {
	url := "http://127.0.0.1" + defaultWebServerAddr + healthy.HealthGroupPath + "/" + checkType + "/close"
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
