package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/skaji/go-server-worker/runner"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := runner.NewRunner(10)
	log.Println("main start")
	r.Start(ctx)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	cancel()
	r.Wait()
	log.Println("main stop")
}
