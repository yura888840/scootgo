package application

import (
	"context"
	"time"

	"github.com/yura888840/scootgo/internal/timer/domain"
)

type StartTimerCommand struct {
	UserID         string
	ScooterID      string
	IdempotencyKey string
}

type TimerRepository interface {
	CreateOnce(
		ctx context.Context,
		userID string,
		idempotencyKey string,
		timer domain.Timer,
	) (domain.Timer, error)
}

type IDGenerator interface {
	NewID() string
}

type Clock interface {
	Now() time.Time
}

type StartTimer struct {
	repository  TimerRepository
	idGenerator IDGenerator
	clock       Clock
}

func NewStartTimer(repository TimerRepository, idGenerator IDGenerator, clock Clock) *StartTimer {
	return &StartTimer{
		repository:  repository,
		idGenerator: idGenerator,
		clock:       clock,
	}
}

func (useCase *StartTimer) Execute(ctx context.Context, command StartTimerCommand) (domain.Timer, error) {
	timer := domain.Timer{
		ID:        useCase.idGenerator.NewID(),
		UserID:    command.UserID,
		ScooterID: command.ScooterID,
		StartedAt: useCase.clock.Now(),
	}

	return useCase.repository.CreateOnce(ctx, command.UserID, command.IdempotencyKey, timer)
}
