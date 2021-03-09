package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/artofey/sysmon/pkg/loadavg"
	"github.com/artofey/sysmon/pkg/loadcpu"
	"github.com/artofey/sysmon/pkg/server"
	"github.com/artofey/sysmon/pkg/statcollector"
	"github.com/artofey/sysmon/pkg/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	storage := storage.NewRingStorage(36000)
	collector := statcollector.NewStatCollector(storage, makeParsers())

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

	server.Shutdown()
	cancel()
	log.Print("App Down")
}

func makeParsers() []statcollector.Parser {
	return []statcollector.Parser{
		loadavg.NewParser(),
		loadcpu.NewParser(),
	}
}
