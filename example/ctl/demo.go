package ctl

import (
	"context"
	"errors"
	"fmt"
	"time"

	base_frame "github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/pflag"
)

var (
	green  = string([]byte{27, 91, 51, 50, 109})
	white  = string([]byte{27, 91, 51, 55, 109})
	yellow = string([]byte{27, 91, 51, 51, 109})
	red    = string([]byte{27, 91, 51, 49, 109})
	blue   = string([]byte{27, 91, 51, 52, 109})
	// magenta = string([]byte{27, 91, 51, 53, 109})
	// cyan    = string([]byte{27, 91, 51, 54, 109})
	// color   = []string{green, white, yellow, red, blue, magenta, cyan}
)

func main() {
	// 0. 注册自己的flags
	subSet := &pflag.FlagSet{}
	subSet.StringP("runflag", "r", "default run value", "subSet as run flag")
	subSet.StringP("shutdownflag", "s", "default shutdown value", "subSet as shutdown flag")
	subSet.Int("sleep", 1, "sleep how many seconds")

	// 1. 新建对应cmd对象，绑定flags
	c := cmd.NewCtlCommand()
	c.RegFlagSet(subSet)

	// 2. 注册前置方法
	checkFlags := func() error {
		i, err := subSet.GetInt("sleep")
		if err != nil {
			return err
		}

		if i < 0 {
			return errors.New("Flag sleep can not be negtive ")
		}
		fmt.Println(red, "=========== pre run check flags =========== ", white)

		return nil
	}

	if err := c.RegPreRunFunc("CheckFlags", checkFlags); err != nil {
		panic(err)
	}

	// 3.注册绑定用户主要方法
	run := func(ctx context.Context) {
		rf, err := subSet.GetString("runflag")
		if err != nil {
			log.Error("error:", err)
		}

		// c.Flags() 包括了之前声明 subSet，和 subSet.GetInt("sleep") 效力相同
		s, err := c.Flags().GetInt("sleep")
		if err != nil {
			log.Error("error:", err)
		}

		userMainFunc(rf, s)
	}
	if err := c.RegRunFunc(run); err != nil {
		panic(err)
	}

	// 4.1 [可选] 注册shutdown方法，只有程序收到 os.Signal 退出的时候会执行
	shut := func(ctx context.Context) {
		fmt.Println(blue, "===========  Shutdown flag:", ctx.Value(cmd.CtxKey("shutdownflag")), " =========== ", white)
	}
	if err := c.RegShutdownFunc(shut); err != nil {
		panic(err)
	}

	// 5. 注册后置方法，所有的都会运行
	err2 := c.RegPostRunFunc("PostPrint2", func() error {
		fmt.Println(green, "=========== post 2 =========== ", white)
		return nil
	})
	if err2 != nil {
		panic(err2)
	}

	err1 := c.RegPostRunFunc("PostPrint1", func() error {
		fmt.Println(green, "=========== post 1 =========== ", white)
		return nil
	})
	if err1 != nil {
		panic(err1)
	}

	// 6. 启动，运行下面的命令试试看
	//
	// 6.1 ctl 模式帮助
	// 		go run example/ctl/main.go ctl -h
	//
	// 6.2 ctl 自己退出
	// 		go run example/ctl/main.go ctl --sleep 5 -r run_args
	//
	// 6.3 ctl 收到 kill 信号
	// 		go run example/ctl/main.go ctl --sleep 300 -r run_args -s shut_args
	// 		ps aux |grep 'main ctl'|grep -v "grep" |awk '{print $2}'|xargs kill
	//
	// 6.4 healthy 健康检查
	// 		go run example/ctl/main.go ctl --sleep 300 -r run_args -s shut_args
	// 		go run example/ctl/main.go healthy
	// 		ps aux |grep 'main ctl'|grep -v "grep" |awk '{print $2}'|xargs kill
	//
	// 		15s 之后再试试 healthy 命令
	// 		go run example/ctl/main.go healthy
	base_frame.Start("ctl", c)
}

// 这里模拟用户的主要逻辑在其他包内
func userMainFunc(rf string, s int) {
	fmt.Println(yellow, "=========== Run flag:", rf, " =========== ", white)
	fmt.Printf("%s =========== ctl will sleep %s ===========  %s \n", yellow, (time.Second * time.Duration(s)).String(), white)
	time.Sleep(time.Second * time.Duration(s))
}
