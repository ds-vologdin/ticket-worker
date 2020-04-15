package server

import (
	"fmt"
	"log"
	"time"

	"github.com/ds-vologdin/ticket-worker/mock"
	"github.com/ds-vologdin/ticket-worker/storage"
	"github.com/ds-vologdin/ticket-worker/storage/memory"
	tk "github.com/ds-vologdin/ticket-worker/ticket"
)

type serverTicket struct {
	Storage storage.TicketStorage
}

var ServerTicket serverTicket

func init() {
	ServerTicket.Storage = memory.NewMock()
}

// ticket manage
func (s *serverTicket) AddTicketType1(details string) error {
	log.Printf("AddTicketType1: %v", details)
	ticketID, err := tk.NewTicketID()
	if err != nil {
		log.Printf("AddTicketType1: NewTicketID error (%v)", err)
		return err
	}
	ticket := tk.Ticket{
		ID:      ticketID,
		Type:    mock.TicketType1,
		Created: time.Now(),
		Status:  tk.Pending,
		Details: details,
	}
	return s.Storage.AddTicketAndInitSteps(ticket)
}

func (s *serverTicket) AddTicketType2(details string) error {
	log.Printf("AddTicketType2: %v", details)
	ticketID, err := tk.NewTicketID()
	if err != nil {
		log.Printf("AddTicketType2: NewTicketID error (%v)", err)
		return err
	}
	ticket := tk.Ticket{
		ID:      ticketID,
		Type:    mock.TicketType2,
		Created: time.Now(),
		Status:  tk.Pending,
		Details: details,
	}
	return s.Storage.AddTicketAndInitSteps(ticket)
}

// for back-office
func (s *serverTicket) GetAllTickets() ([]tk.Ticket, error) {
	return s.Storage.GetAllTickets()
}

func (s *serverTicket) GetActiveTickets() ([]tk.Ticket, error) {
	return s.Storage.GetActiveTickets()
}

func (s *serverTicket) GetTicket(id tk.TicketID) (tk.Ticket, []tk.TicketStep, error) {
	return s.Storage.GetTicket(id)
}

// for emulate manual processing
func (s *serverTicket) GetOnePendingTicketForManualProcessing() (*tk.Ticket, []tk.TicketStep, error) {
	return s.Storage.GetOnePendingTicketForManualProcessing()
}

// for worker
func (s *serverTicket) GetOnePendingTicket() (*tk.Ticket, []tk.TicketStep, error) {
	return s.Storage.GetOnePendingTicketForAutoProcessing()
}

func (s *serverTicket) MarkStepAsProcessed(ticketID tk.TicketID, serial int32) error {
	return s.Storage.MarkStepAsProcessed(ticketID, serial)
}

func (s *serverTicket) SaveStepResult(ticketID tk.TicketID, serial int32, status tk.StepStatus, details string) error {
	return s.Storage.SaveStepResult(ticketID, serial, status, details)
}

func StartServer() {
	const frequency = 5 * time.Second
	for {
		tickets, err := ServerTicket.GetActiveTickets()
		if err != nil {
			log.Printf("GetActiveTickets error: %v", err)
			time.Sleep(frequency)
			continue
		}
		log.Printf("count active tickets: %d", len(tickets))

		fmt.Println("=====================================================================")
		time.Sleep(frequency)
	}
}
