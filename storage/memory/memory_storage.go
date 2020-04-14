package memory

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/ds-vologdin/ticket-worker/ticket"
	tk "github.com/ds-vologdin/ticket-worker/ticket"

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
	Tickets     map[tk.TicketID]tk.Ticket
	TicketSteps map[tk.TicketID][]tk.TicketStep
	OrdersSteps map[tk.TicketType][]tk.Step
	TicketTypes []tk.TicketType
	Steps       []tk.Step
	mx          sync.Mutex
}

func (s *TicketStorageMemory) Init() error {
	s.Tickets = make(map[tk.TicketID]tk.Ticket)
	s.TicketSteps = make(map[tk.TicketID][]tk.TicketStep)
	s.OrdersSteps = make(map[tk.TicketType][]tk.Step)
	s.TicketTypes = make([]tk.TicketType, 0)
	s.Steps = make([]tk.Step, 0)
	return nil
}

func (s *TicketStorageMemory) GetAllTickets() ([]tk.Ticket, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	tickets := make([]tk.Ticket, 0, len(s.Tickets))
	for _, ticket := range s.Tickets {
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (s *TicketStorageMemory) GetActiveTickets() ([]tk.Ticket, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	tickets := make([]tk.Ticket, 0, len(s.Tickets))
	for _, ticket := range s.Tickets {
		if ticket.Status != tk.Success {
			tickets = append(tickets, ticket)
		}
	}
	return tickets, nil
}

func (s *TicketStorageMemory) GetTicket(id ticket.TicketID) (tk.Ticket, []tk.TicketStep, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	ticket, ok := s.Tickets[id]
	if !ok {
		return ticket, nil, fmt.Errorf("Ticket with id %v is not found", id)
	}
	return ticket, s.TicketSteps[id], nil
}

func (s *TicketStorageMemory) GetOrderStepsByTicketType(ticketType tk.TicketType) ([]tk.Step, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	corretType := s.isCorrectTicketType(ticketType)
	if !corretType {
		return nil, fmt.Errorf("Incorrect ticket type")
	}
	return s.OrdersSteps[ticketType], nil
}

func (s *TicketStorageMemory) AddTicketAndInitSteps(tk tk.Ticket) error {
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

func (s *TicketStorageMemory) GetOnePendingTicketForAutoProcessing() (*tk.Ticket, []tk.TicketStep, error) {

	getStepByType := func(stepType tk.StepType) *tk.Step {
		for _, step := range s.Steps {
			if step.Type == stepType {
				return &step
			}
		}
		return nil
	}
	isNextStepAuto := func(ticketSteps []tk.TicketStep) bool {
		for _, step := range ticketSteps {
			if step.Status == tk.Pending {
				stepInfo := getStepByType(step.StepType)
				return stepInfo.Auto
			}
		}
		return false
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	for _, ticket := range s.Tickets {
		if ticket.Status == tk.Pending {
			ticketSteps := s.TicketSteps[ticket.ID]
			if isNextStepAuto(ticketSteps) {
				return &ticket, ticketSteps, nil
			}
		}
	}
	return nil, nil, nil
}

func (s *TicketStorageMemory) GetOnePendingTicketForManualProcessing() (*tk.Ticket, []tk.TicketStep, error) {

	getStepByType := func(stepType tk.StepType) *tk.Step {
		for _, step := range s.Steps {
			if step.Type == stepType {
				return &step
			}
		}
		return nil
	}
	isNextStepAuto := func(ticketSteps []tk.TicketStep) bool {
		for _, step := range ticketSteps {
			if step.Status == tk.Pending {
				stepInfo := getStepByType(step.StepType)
				return stepInfo.Auto
			}
		}
		return false
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	for _, ticket := range s.Tickets {
		if ticket.Status == tk.Pending {
			ticketSteps := s.TicketSteps[ticket.ID]
			if !isNextStepAuto(ticketSteps) {
				return &ticket, ticketSteps, nil
			}
		}
	}
	return nil, nil, nil
}

func (s *TicketStorageMemory) GetTicketStepsByTicketID(ticketID tk.TicketID) ([]tk.TicketStep, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	ticketSteps, ok := s.TicketSteps[ticketID]
	if !ok {
		return nil, fmt.Errorf("Ticket id is not found")
	}
	return ticketSteps, nil
}

func (s *TicketStorageMemory) MarkStepAsProcessed(ticketID tk.TicketID, serial int32) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	ticketSteps, ok := s.TicketSteps[ticketID]
	if !ok {
		return fmt.Errorf("TicketSteps is not found (ticketID: %v)", ticketID)
	}

	index, ticketStep, err := s.findTicketStep(ticketSteps, serial)
	if err != nil {
		return err
	}
	if ticketStep.Status != tk.Pending {
		return fmt.Errorf("Step status is %v", ticketStep.Status)
	}
	ticketStep.Status = tk.Pending
	ticketStep.Started = time.Now()

	ticketSteps[index] = *ticketStep
	return nil
}

func (s *TicketStorageMemory) SaveStepResult(ticketID tk.TicketID, serial int32, status tk.StepStatus, details string) error {
	if status != tk.Success && status != tk.Failed {
		return fmt.Errorf("Status %v is not support for save results", status)
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	ticketSteps, ok := s.TicketSteps[ticketID]
	if !ok {
		return fmt.Errorf("TicketSteps is not found (ticketID: %v)", ticketID)
	}

	index, ticketStep, err := s.findTicketStep(ticketSteps, serial)
	if err != nil {
		return err
	}
	ticketStep.Status = status
	ticketStep.Stoped = time.Now()
	ticketStep.Details = details
	ticketSteps[index] = *ticketStep

	isLastStep := len(ticketSteps) == index+1

	if isLastStep && status == tk.Success {
		ticket := s.Tickets[ticketID]
		ticket.Status = tk.Success
		ticket.Closed = time.Now()
		s.Tickets[ticketID] = ticket
	} else {
		steps := s.TicketSteps[ticketID]
		next := index + 1
		step := steps[next]
		step.Status = tk.Pending
		steps[next] = step
	}
	return nil
}

// TicketStorageMemory helpers
// You do not use mutex s.mx

func (s *TicketStorageMemory) isCorrectTicketType(ticketType tk.TicketType) bool {
	var corretType bool
	for _, tkType := range s.TicketTypes {
		if tkType == ticketType {
			corretType = true
			break
		}
	}
	return corretType
}

func (s *TicketStorageMemory) initTicketSteps(id tk.TicketID) error {
	getStatus := func(i int) tk.StepStatus {
		if i == 0 {
			return tk.Pending
		}
		return tk.Created
	}
	ticket := s.Tickets[id]
	steps := s.OrdersSteps[ticket.Type]
	ticketSteps := make([]tk.TicketStep, 0, len(steps))
	for i, step := range steps {
		nonce, err := uuid.NewRandom()
		if err != nil {
			log.Printf("Error uuid.NewRandom(): %v", err)
		}
		ticketStep := tk.TicketStep{
			TicketID:   id,
			StepType:   step.Type,
			SerialNumb: int32(i * 10),
			Created:    time.Now(),
			Status:     getStatus(i),
			Nonce:      tk.NonceType(nonce),
		}
		ticketSteps = append(ticketSteps, ticketStep)
	}
	s.TicketSteps[id] = ticketSteps
	return nil
}

func (s *TicketStorageMemory) findTicketStep(ticketSteps []tk.TicketStep, serial int32) (int, *tk.TicketStep, error) {
	for index, ticketStep := range ticketSteps {
		if ticketStep.SerialNumb == serial {
			return index, &ticketStep, nil
		}
	}
	return 0, nil, fmt.Errorf("TicketSteps is not found (serial: %v)", serial)
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
	storage.initMockTickets(100)
	return storage
}

func (s *TicketStorageMemory) initMockOrdersSteps() {
	s.mx.Lock()
	defer s.mx.Unlock()
	var order = []tk.Step{s.Steps[0], s.Steps[1], s.Steps[2], s.Steps[3]}
	s.OrdersSteps[TicketType1] = order
	order = []tk.Step{s.Steps[0], s.Steps[1], s.Steps[3]}
	s.OrdersSteps[TicketType2] = order
	order = []tk.Step{s.Steps[0], s.Steps[3]}
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

func getNewStep(stepType tk.StepType, auto bool) tk.Step {
	return tk.Step{Type: stepType, Auto: auto, Description: fmt.Sprintf("mock-%v", stepType)}
}

func getNewTicket(ticketType tk.TicketType) tk.Ticket {
	id := getNewTicketID()
	ticket := tk.Ticket{ID: id, Type: ticketType, Status: tk.Pending, Created: time.Now()}
	return ticket
}

func getNewTicketID() tk.TicketID {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("getNewTicketID error: %v", err)
	}
	return tk.TicketID(id)
}
