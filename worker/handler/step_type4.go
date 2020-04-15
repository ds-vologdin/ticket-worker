package handler

import (
	"context"
	"log"
	"time"

	tk "github.com/ds-vologdin/ticket-worker/ticket"
)

type HandlerStepType4 struct{}

func (h *HandlerStepType4) Run(ctx context.Context, ticket tk.Ticket, step tk.TicketStep) Result {
	log.Printf("HandlerStepType4.Run: %v", step)
	time.Sleep(100 * time.Millisecond)
	return Result{Status: tk.Success, Details: "{\"answer\": \"OK\"}"}
}
