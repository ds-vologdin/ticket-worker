package main

import (
	"log"

	"github.com/ds-vologdin/ticket-worker/server"
)

func main() {
	log.Println("start server")
	server.StartServer()
}
