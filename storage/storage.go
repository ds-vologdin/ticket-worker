package storage

type TicketStorage interface {
	Init() error

	// GetOrderStepsByTicketType(ticketType ticket.TicketType) ([]ticket.StepType, error)
	// GetTickets() ([]ticket.Ticket, error)
	// GetTicket(id ticket.TicketID) (ticket.Ticket, error)
	// GetStepsByTicketID(ticketID ticket.TicketID) ([]ticket.Step, error)

	// AddTicketWithEmptySteps(tkt ticket.Ticket) error
	// UpdateTicket(tkt ticket.Ticket) error

	// AddTicketStep(step ticket.TicketStep) error
}
