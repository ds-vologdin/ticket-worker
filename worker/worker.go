package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ds-vologdin/ticket-worker/mock"
	"github.com/ds-vologdin/ticket-worker/server"
	tk "github.com/ds-vologdin/ticket-worker/ticket"
	"github.com/ds-vologdin/ticket-worker/worker/handler"
)

const Frequency = 1 * time.Second

func runHandler(ctx context.Context) error {
	ticket, steps, err := server.ServerTicket.GetOnePendingTicket()
	if err != nil {
		log.Printf("Worker: GetOnePendingTicket error (%v)", err)
		return err
	}
	if ticket == nil {
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

	handler := selectHandler(*step)
	result := handler.Run(ctx, *ticket, *step)

	err = server.ServerTicket.SaveStepResult(step.TicketID, step.SerialNumb, result.Status, result.Details)
	if err != nil {
		log.Printf("Worker: SaveStepResult error (%v)", err)
		return err
	}
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
			continue
		}
	}
}

func StartTicketWorker(ctx context.Context) {
	go ticketWorker(ctx)
}

// helpers

func selectHandler(step tk.TicketStep) handler.Handler {
	switch step.StepType {
	case mock.StepType1:
		return &handler.HandlerStepType1{}
	case mock.StepType2:
		return &handler.HandlerStepType2{}
	case mock.StepType3:
		return &handler.HandlerStepType3{}
	case mock.StepType4:
		return &handler.HandlerStepType4{}
	default:
		log.Printf("selectHandler: %v is an unknown step type", step.StepType)
	}
	return nil
}

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
