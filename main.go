package main

import (
	"log"

	"github.com/ds-vologdin/ticket-worker/storage/memory"
)

func main() {
	log.Println("start server")
	storage := memory.NewMock()
}
