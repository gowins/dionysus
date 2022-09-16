package grpool

import (
	"fmt"
	"github.com/gowins/dionysus/grpool"
)

func grpoolDemo() error {
	err := grpool.Submit(func() {
		panic("test panic")
		fmt.Println(123)
	})
	fmt.Println(666888)

	return err
}
