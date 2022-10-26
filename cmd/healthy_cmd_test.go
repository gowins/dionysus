package cmd

import (
	"github.com/gowins/dionysus/healthy"
	"testing"
)

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
	cmd.Run(nil, nil)
}
