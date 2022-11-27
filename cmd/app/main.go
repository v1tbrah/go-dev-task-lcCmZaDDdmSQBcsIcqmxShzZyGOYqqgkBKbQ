package main

import (
	"log"
	"os"

	"go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/app"
)

func main() {

	newApp, err := app.New()
	if err != nil {
		log.Print("creating app: " + err.Error())
		os.Exit(1)
	}

	if err = newApp.Run(); err != nil {
		log.Print("running app: " + err.Error())
		if !newApp.RMQConnIsClosed() {
			if err = newApp.CloseRMQConn(); err != nil {
				log.Print("closing rmq conn: " + err.Error())
			}
		}
		os.Exit(1)
	}

	if !newApp.RMQConnIsClosed() {
		if err = newApp.CloseRMQConn(); err != nil {
			log.Print("closing rmq conn: " + err.Error())
		}
	}

}
