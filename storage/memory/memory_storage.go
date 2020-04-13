package memory

import (
	"github.com/ds-vologdin/ticket-worker/storage"
	"github.com/ds-vologdin/ticket-worker/ticket"
	"github.com/ds-vologdin/ticket-worker/ticket/storage"
)

type KeyTicketStep struct {
	TicketID   ticket.TicketID
	SerialNumb int32
}

type TicketStorageMemory struct {
	Tickets     map[ticket.TicketID]ticket.Ticket
	TicketSteps map[ticket.TicketID]ticket.TicketStep
	OrdersSteps map[ticket.TicketType][]ticket.Step
}

func (s *TicketStorageMemory) Init() error {
	s.Tickets = make(map[ticket.TicketID]ticket.Ticket)
	s.TicketSteps = make(map[ticket.TicketID]ticket.TicketStep)
	s.OrdersSteps = make(map[ticket.TicketType][]ticket.Step)
	return nil
}

// func GetStepsOrderByTicketType(ticketType TicketType) ([]StepType, error) {
// 	switch ticketType {
// 	case TicketType1:
// 		return []StepType{StepType1, StepType2, StepType3}, nil
// 	case TicketType2:
// 		return []StepType{StepType1, StepType4}, nil
// 	case TicketType3:
// 		return []StepType{StepType3, StepType4}, nil
// 	}
// 	return nil, nil
// }

func NewMock() *storage.TicketStorage {
	const (
		// Types of tickets
		TicketType1 = 1
		TicketType2 = 2
		TicketType3 = 3

		// Statuses of tickets
		Pending   = 1
		Processed = 2
		Failed    = 3
		Success   = 4

		// Types of steps
		StepType1 = 1
		StepType2 = 2
		StepType3 = 3
		StepType4 = 4
	)
	var storage TicketStorageMemory
	storage.Init()
	return &storage
}
