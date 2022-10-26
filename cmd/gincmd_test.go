package cmd

import (
	"fmt"
	"github.com/gowins/dionysus/healthy"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/pflag"
	"net/http"
	"os"
	"testing"
	"time"
)

var healthHttpClient = http.Client{Timeout: 3 * time.Second}

func TestNewGinCommand(t *testing.T) {
	log.Setup(log.SetProjectName("projectName"), log.WithWriter(os.Stdout))
	gcmd := NewGinCommand()
	gcmd.RegFlagSet(&pflag.FlagSet{})
	gcmd.Flags()
	gcmd.GetCmd()
	gcmd.GetShutdownFunc()
	go func() {
		gcmd.registerHealth()
		gcmd.startServer()
		time.Sleep(time.Second * 10)
		gcmd.stopServer()
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
