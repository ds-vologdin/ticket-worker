package server

import (
	"github.com/ds-vologdin/ticket-worker/storage"
	"github.com/ds-vologdin/ticket-worker/storage/memory"
	"github.com/ds-vologdin/ticket-worker/ticket"
)

type serverTicket struct {
	Storage storage.TicketStorage
}

var ServerTicket serverTicket

func init() {
	ServerTicket.Storage = memory.NewMock()
}

// for back-office
func (s *serverTicket) AddTicket(ticket ticket.Ticket) error {
	err := s.Storage.AddTicketAndInitSteps(ticket)
	return err
}

// for worker
func (s *serverTicket) GetPendingSteps() {

}

func StartServer() {

}
