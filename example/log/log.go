package main

import (
	"github.com/gowins/dionysus/log"
	"github.com/gowins/dionysus/log/writer/rotate"
)

func main() {
	// rotate example
	cfg := rotate.NewWriterConfig()
	cfg.Dir = "."
	w, err := rotate.NewRotateLogger(cfg)
	if err != nil {
		panic(err)
	}
	log.Setup(log.SetProjectName("Test"), log.WithWriter(w))
	log.Debug("haha")
}
