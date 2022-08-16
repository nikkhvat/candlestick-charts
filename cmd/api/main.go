package main

import (
	"forex/pkg/config"
	"forex/server"
	"log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)
	conf := config.GetConfig()

	app := server.NewApp()

	if err := app.Run(conf.Port); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
