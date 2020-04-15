package handler

import (
	"context"

	tk "github.com/ds-vologdin/ticket-worker/ticket"
)

type Result struct {
	Status  tk.StepStatus
	Details string
}

type Handler interface {
	Run(ctx context.Context, ticket tk.Ticket, step tk.TicketStep) Result
}
