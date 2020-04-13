package memory

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ds-vologdin/ticket-worker/ticket"

	"github.com/google/uuid"
)

const (
	// Types of tickets
	TicketType1 = 1
	TicketType2 = 2
	TicketType3 = 3

	// Types of steps
	StepType1 = 1
	StepType2 = 2
	StepType3 = 3
	StepType4 = 4
)

type KeyTicketStep struct {
	TicketID   ticket.TicketID
	SerialNumb int32
}

type TicketStorageMemory struct {
	Tickets     map[ticket.TicketID]ticket.Ticket
	TicketSteps map[ticket.TicketID]ticket.TicketStep
	OrdersSteps map[ticket.TicketType][]ticket.Step
	TicketTypes []ticket.TicketType
	Steps       []ticket.Step
}

func (s *TicketStorageMemory) Init() error {
	s.Tickets = make(map[ticket.TicketID]ticket.Ticket)
	s.TicketSteps = make(map[ticket.TicketID]ticket.TicketStep)
	s.OrdersSteps = make(map[ticket.TicketType][]ticket.Step)
	s.TicketTypes = make([]ticket.TicketType, 0)
	s.Steps = make([]ticket.Step, 0)
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

// MOCK DATA

func NewMock() *TicketStorageMemory {

	var storage = &TicketStorageMemory{}
	err := storage.Init()
	if err != nil {
		log.Printf("Init TicketStorageMemory error: %v", err)
		return nil
	}
	storage.initMockTicketType()
	storage.initMockSteps()
	storage.initMockTickets(20)
	storage.initMockOrdersSteps()
	return storage
}

func (s *TicketStorageMemory) initMockOrdersSteps() {
	var order = []ticket.Step{s.Steps[0], s.Steps[1], s.Steps[2], s.Steps[3]}
	s.OrdersSteps[TicketType1] = order
	order = []ticket.Step{s.Steps[0], s.Steps[1], s.Steps[3]}
	s.OrdersSteps[TicketType2] = order
	order = []ticket.Step{s.Steps[0], s.Steps[3]}
	s.OrdersSteps[TicketType3] = order
}

func (s *TicketStorageMemory) initMockSteps() {
	s.Steps = append(s.Steps, getNewStep(StepType1, true))
	s.Steps = append(s.Steps, getNewStep(StepType2, true))
	s.Steps = append(s.Steps, getNewStep(StepType3, false))
	s.Steps = append(s.Steps, getNewStep(StepType4, true))
}

func (s *TicketStorageMemory) initMockTicketType() {
	s.TicketTypes = append(s.TicketTypes, TicketType1)
	s.TicketTypes = append(s.TicketTypes, TicketType2)
	s.TicketTypes = append(s.TicketTypes, TicketType3)
}

func (s *TicketStorageMemory) initMockTickets(count int) {
	for i := 0; i < count; i++ {
		ticketType := s.TicketTypes[rand.Intn(len(s.TicketTypes))]
		ticket := getNewTicket(ticketType)
		s.Tickets[ticket.ID] = ticket
	}
}

func getNewStep(stepType ticket.StepType, auto bool) ticket.Step {
	return ticket.Step{Type: stepType, Auto: auto, Description: fmt.Sprintf("mock-%v", stepType)}
}

func getNewTicket(ticketType ticket.TicketType) ticket.Ticket {
	id := getNewTicketID()
	ticket := ticket.Ticket{ID: id, Type: ticketType, Status: ticket.Pending, Created: time.Now()}
	return ticket
}

func getNewTicketID() ticket.TicketID {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("getNewTicketID error: %v", err)
	}
	return ticket.TicketID(id)
}
