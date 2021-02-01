package app

import (
	"fmt"
	"log"

	"github.com/artofey/sysmon/internal/pb"
	"github.com/artofey/sysmon/internal/server"
	"github.com/artofey/sysmon/internal/statcollector"
)

func ErrHandle(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type App struct {
	Stats []*pb.StatSnapshot
}

func (a *App) Run() {
	c := make(chan *pb.StatSnapshot)
	go statcollector.StartColecting(c)

	GRPCServer, err := server.New(9000)
	ErrHandle(err)

	go func() {
		err := server.Start(GRPCServer)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}()

	for sc := range c {
		a.Stats = append(a.Stats, sc)
		fmt.Println(sc.Lavg.Load1, sc.Lavg.Load5, sc.Lavg.Load15)
		fmt.Println(sc.Lcpu.User, sc.Lcpu.System, sc.Lcpu.Idle)
	}
}

func New() *App {
	app := App{}
	return &app
}
