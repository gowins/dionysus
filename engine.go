package dionysus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gowins/dionysus/cmd"

	"github.com/spf13/cobra"
)

type Dio struct {
	cmd *cobra.Command
	// steps *step.Steps // 全局启动依赖项,暂时注释掉入口。

	// e             *gin.Engine
	// groups        []*CustomGroup
	// globalMiddles []func() gin.HandlerFunc
	// server        *http.Server
}

var dio *Dio

var (
	BuildTime = ""
	CommitID  = ""
	GoVersion = ""
)

func init() {
	dio = NewDio()
}

func NewDio() *Dio {
	d := &Dio{
		// root cmd
		cmd: &cobra.Command{},
		// gin cmd
	}

	return d
}

func Start(project string, cmds ...cmd.Commander) {
	// 0 set cmd use
	dio.cmd.Use = project
	if ex, err := os.Executable(); dio.cmd.Use == "" && err == nil {
		dio.cmd.Use = filepath.Base(ex)
	}

	// 1. global flags
	var pname, configPath string
	var version bool
	// var registryVal string
	dio.cmd.PersistentFlags().StringVarP(&configPath, "config", "c", os.Getenv("GAPI_CONFIG"), "deprecated: the config file path")
	dio.cmd.PersistentFlags().StringVarP(&pname, "name", "n", os.Getenv("GAPI_PROJECT_NAME"), "the project name")
	//dio.cmd.PersistentFlags().StringVar(&environment, "env", algs.FirstNotEmpty(os.Getenv(env.GetRunEnvKey()), env.Develop), "the project run mod available in [ product | gray| test | develop ]")
	dio.cmd.PersistentFlags().BoolVar(&version, "version", false, "Print build version")

	// 2. global pre run function
	dio.cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// //////////////////////////////////////////////////////////////// //
		// Project name
		if project != "" && pname == "" {
			pname = project
		}
		if pname == "" {
			return errors.New("the project name cannot be empty")
		}
		// //////////////////////////////////////////////////////////////// //

		// //////////////////////////////////////////////////////////////// //
		// Department name
		//if department != "" && dname == "" {
		//	dname = department
		//}
		//if dname == "" {
		//	return errors.New("the department name cannot be empty")
		//}

		return nil
	}

	dio.cmd.Run = func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println("Build Time:", BuildTime)
			fmt.Println("Go Version:", GoVersion)
			fmt.Println("Commit  ID:", CommitID)
		} else {
			_ = cmd.Help()
		}
	}

	// 3. append other cmds
	for _, c := range cmds {
		dio.cmd.AddCommand(c.GetCmd())
	}

	// TODO SIDE EFFECT
	// defer logger.Close()

	// 4. start
	if err := dio.cmd.Execute(); err != nil {
		panic(err)
	}
}
