package memory

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
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

type TicketStorageMemory struct {
	Tickets     map[ticket.TicketID]ticket.Ticket
	TicketSteps map[ticket.TicketID][]ticket.TicketStep
	OrdersSteps map[ticket.TicketType][]ticket.Step
	TicketTypes []ticket.TicketType
	Steps       []ticket.Step
	mx          sync.Mutex
}

func (s *TicketStorageMemory) Init() error {
	s.Tickets = make(map[ticket.TicketID]ticket.Ticket)
	s.TicketSteps = make(map[ticket.TicketID][]ticket.TicketStep)
	s.OrdersSteps = make(map[ticket.TicketType][]ticket.Step)
	s.TicketTypes = make([]ticket.TicketType, 0)
	s.Steps = make([]ticket.Step, 0)
	return nil
}

func (s *TicketStorageMemory) GetOrderStepsByTicketType(ticketType ticket.TicketType) ([]ticket.Step, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	corretType := s.isCorrectTicketType(ticketType)
	if !corretType {
		return nil, fmt.Errorf("Incorrect ticket type")
	}
	return s.OrdersSteps[ticketType], nil
}

func (s *TicketStorageMemory) AddTicketAndInitSteps(tk ticket.Ticket) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if correct := s.isCorrectTicketType(tk.Type); !correct {
		return fmt.Errorf("Ticket type is incorrect")
	}
	s.Tickets[tk.ID] = tk
	err := s.initTicketSteps(tk.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TicketStorageMemory) GetTicketStepsByTicketID(ticketID ticket.TicketID) ([]ticket.TicketStep, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	ticketSteps, ok := s.TicketSteps[ticketID]
	if !ok {
		return nil, fmt.Errorf("Ticket id is not found")
	}
	return ticketSteps, nil
}

// TicketStorageMemory helpers
// You do not use mutex s.mx

func (s *TicketStorageMemory) isCorrectTicketType(ticketType ticket.TicketType) bool {
	var corretType bool
	for _, tkType := range s.TicketTypes {
		if tkType == ticketType {
			corretType = true
			break
		}
	}
	return corretType
}

func (s *TicketStorageMemory) initTicketSteps(id ticket.TicketID) error {
	currentTicket := s.Tickets[id]
	steps := s.OrdersSteps[currentTicket.Type]
	ticketSteps := make([]ticket.TicketStep, 0, len(steps))
	for i, step := range steps {
		ticketStep := ticket.TicketStep{
			TicketID:   id,
			StepType:   step.Type,
			SerialNumb: int32(i * 10),
			Created:    time.Now(),
		}
		ticketSteps = append(ticketSteps, ticketStep)
	}
	s.TicketSteps[id] = ticketSteps
	return nil
}

// New

func New() *TicketStorageMemory {
	var storage = &TicketStorageMemory{}
	err := storage.Init()
	if err != nil {
		log.Printf("TicketStorageMemory error: %v", err)
		return nil
	}
	return storage
}

// MOCK DATA

func NewMock() *TicketStorageMemory {
	storage := New()
	storage.initMockTicketType()
	storage.initMockSteps()
	storage.initMockOrdersSteps()
	storage.initMockTickets(200)
	return storage
}

func (s *TicketStorageMemory) initMockOrdersSteps() {
	s.mx.Lock()
	defer s.mx.Unlock()
	var order = []ticket.Step{s.Steps[0], s.Steps[1], s.Steps[2], s.Steps[3]}
	s.OrdersSteps[TicketType1] = order
	order = []ticket.Step{s.Steps[0], s.Steps[1], s.Steps[3]}
	s.OrdersSteps[TicketType2] = order
	order = []ticket.Step{s.Steps[0], s.Steps[3]}
	s.OrdersSteps[TicketType3] = order
}

func (s *TicketStorageMemory) initMockSteps() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Steps = append(s.Steps, getNewStep(StepType1, true))
	s.Steps = append(s.Steps, getNewStep(StepType2, true))
	s.Steps = append(s.Steps, getNewStep(StepType3, false))
	s.Steps = append(s.Steps, getNewStep(StepType4, true))
}

func (s *TicketStorageMemory) initMockTicketType() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.TicketTypes = append(s.TicketTypes, TicketType1)
	s.TicketTypes = append(s.TicketTypes, TicketType2)
	s.TicketTypes = append(s.TicketTypes, TicketType3)
}

func (s *TicketStorageMemory) initMockTickets(count int) {
	s.mx.Lock()
	defer s.mx.Unlock()
	for i := 0; i < count; i++ {
		ticketType := s.TicketTypes[rand.Intn(len(s.TicketTypes))]
		tk := getNewTicket(ticketType)
		s.Tickets[tk.ID] = tk
		err := s.initTicketSteps(tk.ID)
		if err != nil {
			log.Printf("initTicketSteps error: %v", err)
		}
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
