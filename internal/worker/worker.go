package worker

import (
	"context"

	"github.com/google/uuid"
)

type Worker struct {
	workerID uuid.UUID
}

func NewWorker() *Worker {
	return &Worker{
		workerID: uuid.New(),
	}
}

func (w *Worker) Setup(ctx context.Context) error {

	return nil
}
