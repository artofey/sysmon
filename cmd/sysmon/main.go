package main

import (
	"log"

	"github.com/artofey/sysmon/internal/app"
)

func errHandle(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	app := app.New()
	app.Run()
}
