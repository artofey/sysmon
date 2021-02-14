package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/artofey/sysmon/pkg/server"
	"github.com/artofey/sysmon/pkg/statcollector"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	collector := statcollector.NewStatCollector()
	go collector.StartCollecting(ctx)

	server := server.NewServer(collector)

	go func() {
		if err := server.Run(50051); err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("App Started")

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
	<-exitCh

	log.Print("App Down")

	server.Shutdown()
	cancel()
}
