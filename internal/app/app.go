package app

import (
	"fmt"

	"github.com/artofey/sysmon/internal/statcollector"
)

type App struct {
	sc *statcollector.StatCollector
}

func (a *App) Run() {
	c := make(chan *statcollector.StatCollector)
	go statcollector.StartColecting(c)
	for sc := range c {
		fmt.Println(*sc.Lavg, *sc.Lcpu)
	}
}

func New() *App {
	app := App{
		sc: statcollector.New(),
	}
	return &app
}
