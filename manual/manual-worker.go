package manual

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ds-vologdin/ticket-worker/server"
	tk "github.com/ds-vologdin/ticket-worker/ticket"
)

const Frequency = 15 * time.Second

func runHandler(ctx context.Context) error {
	ticket, steps, err := server.ServerTicket.GetOnePendingTicketForManualProcessing()
	if err != nil {
		log.Printf("Worker: GetOnePendingTicket error (%v)", err)
		return err
	}
	if ticket == nil {
		log.Println("Worker: ticket is nil")
		return nil
	}
	if isStepsEmpty(steps) {
		log.Println("Worker: order of steps is empty")
		return nil
	}
	step := getCurrentStep(steps)
	if step == nil {
		err = fmt.Errorf("There is not Pending step")
		log.Printf("Worker: %v", err)
		return err
	}

	err = server.ServerTicket.MarkStepAsProcessed(step.TicketID, step.SerialNumb)
	if err != nil {
		log.Printf("Worker: MarkStepAsProcessed error (%v)", err)
		return err
	}

	err = server.ServerTicket.SaveStepResult(step.TicketID, step.SerialNumb, tk.Success, "manual processing")
	if err != nil {
		log.Printf("Worker: SaveStepResult error (%v)", err)
		return err
	}
	log.Printf("Manual worker: finished %v", step)
	return nil
}

func ticketWorker(ctx context.Context) {
	for {
		err := runHandler(ctx)
		if err != nil {
			log.Printf("Worker: handler error (%v)", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(Frequency):
			log.Println("Worker: new tour")
		}
	}
}

func StartTicketWorker(ctx context.Context) {
	go ticketWorker(ctx)
}

// helpers

func getCurrentStep(steps []tk.TicketStep) *tk.TicketStep {
	for _, step := range steps {
		if step.Status == tk.Pending {
			return &step
		}
	}
	return nil
}

func isStepsEmpty(steps []tk.TicketStep) bool {
	if steps == nil {
		return true
	}
	if len(steps) == 0 {
		return true
	}
	return false
}
