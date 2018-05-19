package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/skaji/go-server-worker/runner"
)

func main() {
	r := runner.NewRunner(10)
	log.Println("main start")
	r.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	r.Stop()
	log.Println("main stop")
}
