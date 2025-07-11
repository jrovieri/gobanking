package main

import (
	"github.com/jrovieri/gobanking/app"
	"github.com/jrovieri/gobanking/logger"
)

func main() {

	logger.Info("Starting the application...")
	app.Start()

}
