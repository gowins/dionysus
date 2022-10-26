package cmd

import (
	"github.com/gowins/dionysus/healthy"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/cobra"
	"os"
)

type healthCmd struct {
	cmd *cobra.Command
}

const (
	HealthProjectName = "healthx"
	HealthStatus      = "HEALTH_STATUS"
	statusOpen        = "open"
	statusClose       = "close"
)

func NewHealthCmd(use string) *healthCmd {
	if use == "" {
		use = HealthProjectName
	}
	h := &healthCmd{
		cmd: &cobra.Command{Use: use},
	}
	h.cmd.Short = "Check service healthy status"
	h.cmd.Long = "Healthy command will exited with code:0 and msg:success when service is health."
	h.cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		log.Setup(log.SetProjectName(HealthProjectName), log.WithWriter(os.Stdout))
		return nil
	}
	h.cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	h.cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	h.cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	return h
}

func GetHttpLivenessCmd() *cobra.Command {
	url := defaultWebServerAddr + healthy.HealthGroupPath + healthy.HealthLivenessPath
	return NewHealthCmd(healthy.HealthLiveness).GetHttpCheckCmd(url)
}

func GetHttpReadinessCmd() *cobra.Command {
	url := defaultWebServerAddr + healthy.HealthGroupPath + healthy.HealthReadinessPath
	return NewHealthCmd(healthy.HealthReadiness).GetHttpCheckCmd(url)
}

func GetHttpStartupCmd() *cobra.Command {
	url := defaultWebServerAddr + healthy.HealthGroupPath + healthy.HealthStartupPath
	return NewHealthCmd(healthy.HealthStartup).GetHttpCheckCmd(url)
}

func GetCtlLivenessCmd() *cobra.Command {
	return NewHealthCmd(healthy.HealthLiveness).GetCtlCheckCmd()
}

func GetCtlReadinessCmd() *cobra.Command {
	return NewHealthCmd(healthy.HealthReadiness).GetCtlCheckCmd()
}

func GetCtlStartupCmd() *cobra.Command {
	return NewHealthCmd(healthy.HealthStartup).GetCtlCheckCmd()
}

func GetGrpcLivenessCmd() *cobra.Command {
	return NewHealthCmd(healthy.HealthLiveness).GetGrpcCheckCmd()
}

func GetGrpcReadinessCmd() *cobra.Command {
	return NewHealthCmd(healthy.HealthReadiness).GetGrpcCheckCmd()
}

func GetGrpcStartupCmd() *cobra.Command {
	return NewHealthCmd(healthy.HealthStartup).GetGrpcCheckCmd()
}

func (h *healthCmd) GetCtlCheckCmd() *cobra.Command {
	h.cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := healthy.CheckCtlHealthyStat(h.cmd.Use); err != nil {
			log.Fatal(err)
		}
	}
	return h.cmd
}

func (h *healthCmd) GetHttpCheckCmd(url string) *cobra.Command {
	h.cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := healthy.CheckHttpHealthyStat(url, h.cmd.Use); err != nil {
			log.Fatal(err)
		}
	}
	return h.cmd
}

func (h *healthCmd) GetGrpcCheckCmd() *cobra.Command {
	h.cmd.Run = func(cmd *cobra.Command, args []string) {
		if status := os.Getenv(HealthStatus); status != "" {
			if status == statusOpen {
				if err := healthy.SetGrpcHealthyOpen(defaultGrpcAddr, h.cmd.Use); err != nil {
					log.Fatal(err)
				}
				return
			}
			if status == statusClose {
				if err := healthy.SetGrpcHealthyClose(defaultGrpcAddr, h.cmd.Use); err != nil {
					log.Fatal(err)
				}
				return
			}
		}
		if err := healthy.CheckGrpcHealthy(defaultGrpcAddr, h.cmd.Use); err != nil {
			log.Fatal(err)
		}
	}
	return h.cmd
}
