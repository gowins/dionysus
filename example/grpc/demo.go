package main

import (
	"fmt"
	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
)

func main() {
	gcmd := cmd.NewGrpcCommand()
	d := dionysus.NewDio()
	if err := d.DioStart("grpcdemo", gcmd); err != nil {
		fmt.Printf("dio start error %v\n", err)
	}
}
