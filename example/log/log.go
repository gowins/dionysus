package main

import "github.com/gowins/dionysus/log"

func main() {
	log.Setup(log.SetProjectName("Test"))
	log.Debug("haha")
}
