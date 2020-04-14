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
	Created   = "Created"
	Pending   = "Pending"
	Processed = "Processed"
	Failed    = "Failed"
	Success   = "Success"

	// Types of steps
	StepType1 = 1
	StepType2 = 2
	StepType3 = 3
	StepType4 = 4
)

type TicketType int32
type TicketStatus string
type StepType int32
type StepStatus string
type TicketID uuid.UUID
type NonceType uuid.UUID

type Ticket struct {
	ID          TicketID
	Type        TicketType
	Description string
	Created     time.Time
	Closed      time.Time
	Status      TicketStatus
	Details     string // json
}

type Step struct {
	Type        StepType
	Auto        bool
	Description string
}

// Создаются бэкендом при уведомлении со стороны воркера
// key: (TicketID, SerialNumb)
type TicketStep struct {
	TicketID   TicketID
	SerialNumb int32
	StepType   StepType
	Nonce      NonceType
	Status     StepStatus
	Created    time.Time
	Started    time.Time
	Stoped     time.Time
	Details    string // json
}

type TicketStepOrder struct {
	TicketType  TicketType
	Order       []Step
	Description string
}

func NewTicketID() (TicketID, error) {
	id, err := uuid.NewRandom()
	return TicketID(id), err
}
