package storage

import "github.com/ds-vologdin/ticket-worker/ticket"

type TicketStorage interface {
	Init() error

	GetOrderStepsByTicketType(ticketType ticket.TicketType) ([]ticket.Step, error)
	// InitTicketSteps(ticket.TicketID) error
	AddTicketAndInitSteps(ticket.Ticket) error
	// GetTickets() ([]ticket.Ticket, error)
	// GetTicket(id ticket.TicketID) (ticket.Ticket, error)
	GetTicketStepsByTicketID(ticketID ticket.TicketID) ([]ticket.TicketStep, error)

	// AddTicketWithEmptySteps(tkt ticket.Ticket) error
	// UpdateTicket(tkt ticket.Ticket) error

	// AddTicketStep(step ticket.TicketStep) error
}
