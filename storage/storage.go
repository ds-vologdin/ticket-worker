package storage

import (
	"github.com/ds-vologdin/ticket-worker/ticket"
)

type TicketStorage interface {
	Init() error

	// manage tickets
	AddTicketAndInitSteps(ticket.Ticket) error
	GetOrderStepsByTicketType(ticketType ticket.TicketType) ([]ticket.Step, error)

	// for back-office
	GetAllTickets() ([]ticket.Ticket, error)
	GetActiveTickets() ([]ticket.Ticket, error)
	GetTicket(id ticket.TicketID) (ticket.Ticket, []ticket.TicketStep, error)
	// for test
	GetOnePendingTicketForManualProcessing() (*ticket.Ticket, []ticket.TicketStep, error)

	// for worker
	GetOnePendingTicketForAutoProcessing() (*ticket.Ticket, []ticket.TicketStep, error)
	MarkStepAsProcessed(ticketID ticket.TicketID, serial int32) error
	SaveStepResult(ticketID ticket.TicketID, serial int32, status ticket.StepStatus, details string) error
}
