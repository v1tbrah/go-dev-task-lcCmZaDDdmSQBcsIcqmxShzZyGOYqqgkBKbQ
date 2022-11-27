package main

import (
	"log"
	"os"

	"go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/logger"
)

func main() {

	newLogConsumer := logger.New()

	if err := newLogConsumer.Run(); err != nil {
		log.Print("error: " + err.Error())
		os.Exit(1)
	}

}
