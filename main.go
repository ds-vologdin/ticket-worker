package main

import (
	"fmt"
	"log"

	"github.com/ds-vologdin/ticket-worker/storage/memory"
)

func main() {
	log.Println("start server")
	storage := memory.NewMock()
	for ticketID := range storage.Tickets {
		ticket := storage.Tickets[ticketID]
		fmt.Println(ticket)
		order := storage.OrdersSteps[ticket.Type]
		fmt.Printf("steps order: %v\n", order)
		ticketSteps, err := storage.GetTicketStepsByTicketID(ticket.ID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("ticket steps: %v\n", ticketSteps)
	}
}
