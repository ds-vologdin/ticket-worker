package ticket

import (
	"time"

	"github.com/google/uuid"
)

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

type TicketType int32
type TicketStatus int32
type StepType int32
type StepStatus int32
type TicketID uuid.UUID

type Ticket struct {
	ID          TicketID
	Type        TicketType
	Description string
	Created     time.Time
	Closed      time.Time
	Status      TicketStatus
}

type Step struct {
	Type        StepType
	Auto        bool
	Description string
}

// Создаются бэкендом при уведомлении со стороны воркера
type TicketStep struct {
	TicketID   TicketID
	StepType   StepType
	SerialNumb int32
	Nonce      uuid.UUID
	Status     StepStatus
	Start      time.Time
	Stop       time.Time
	Info       string
}

type TicketStepOrder struct {
	TicketType  TicketType
	Order       []Step
	Description string
}
