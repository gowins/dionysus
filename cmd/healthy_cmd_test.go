package cmd

import (
	"fmt"
	"github.com/gowins/dionysus/healthy"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

var _ = func() error {
	log.Setup(log.SetProjectName("projectName"), log.WithWriter(os.Stdout))
	return nil
}

func TestNewHealthCmd(t *testing.T) {
	NewHealthCmd("")
}

func TestGetHttpCmd(t *testing.T) {
	GetHttpLivenessCmd()
	GetHttpReadinessCmd()
	GetHttpStartupCmd()
}

func TestGetGrpcCmd(t *testing.T) {
	GetGrpcLivenessCmd()
	GetGrpcReadinessCmd()
	GetGrpcStartupCmd()
}

func TestGetCtlCmd(t *testing.T) {
	GetCtlLivenessCmd()
	GetCtlReadinessCmd()
	GetCtlStartupCmd()
}

func Test_healthCmd_GetCtlCheckCmd(t *testing.T) {
	cmd := NewHealthCmd(healthy.HealthLiveness).GetCtlCheckCmd()
	err := cmd.PersistentPreRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PreRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PostRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PersistentPostRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	cmd.Run(nil, nil)
}

func Test_healthCmd_GetGrpcCheckCmd(t *testing.T) {
	cmd := NewHealthCmd(healthy.HealthLiveness).GetGrpcCheckCmd()
	err := cmd.PersistentPreRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PreRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PostRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PersistentPostRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	if err := cmd.RunE(nil, nil); err != nil {
		t.Errorf("wann error nil get error %v", err)
	}
	healthy.RegLivenessCheckers(func() error {
		return fmt.Errorf("test healthy liveness error")
	})
	healthy.RegLivenessCheckers(func() error {
		return fmt.Errorf("test healthy liveness error")
	})
	os.Setenv(healthy.HealthStatus, healthy.StatusClose)
	cmd.Run(nil, nil)
	os.Setenv(healthy.HealthStatus, healthy.StatusOpen)
	cmd.Run(nil, nil)
}

func Test_healthCmd_GetHttpCheckCmd(t *testing.T) {
	cmd := NewHealthCmd(healthy.HealthLiveness).GetHttpCheckCmd("errurl")
	err := cmd.PersistentPreRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PreRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PostRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	err = cmd.PersistentPostRunE(nil, nil)
	if err != nil {
		t.Errorf("wann error nil get error %v", err)
		return
	}
	healthy.RegLivenessCheckers(func() error {
		return fmt.Errorf("test healthy liveness error")
	})
	os.Setenv(healthy.HealthStatus, healthy.StatusClose)
	cmd.Run(nil, nil)
	os.Setenv(healthy.HealthStatus, healthy.StatusOpen)
	cmd.Run(nil, nil)
}
