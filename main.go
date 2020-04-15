package main

import (
	"context"
	"log"

	"github.com/ds-vologdin/ticket-worker/manual"
	"github.com/ds-vologdin/ticket-worker/server"
	"github.com/ds-vologdin/ticket-worker/worker"
)

const WorkerCount = 20

func main() {
	log.Println("start server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < WorkerCount; i++ {
		worker.StartTicketWorker(ctx)
		log.Println("start worker", i)
	}
	manual.StartTicketWorker(ctx)
	server.StartServer()
}
