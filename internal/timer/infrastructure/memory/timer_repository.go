package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/yura888840/scootgo/internal/timer/domain"
)

var ErrIdempotencyKeyConflict = errors.New("idempotency key already used for another timer")

type idempotencyScope struct {
	userID string
	key    string
}

type TimerRepository struct {
	mu          sync.Mutex
	timers      map[string]domain.Timer
	idempotency map[idempotencyScope]string
}

func NewTimerRepository() *TimerRepository {
	return &TimerRepository{
		timers:      make(map[string]domain.Timer),
		idempotency: make(map[idempotencyScope]string),
	}
}

func (repository *TimerRepository) CreateOnce(
	_ context.Context,
	userID string,
	idempotencyKey string,
	timer domain.Timer,
) (domain.Timer, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	scope := idempotencyScope{userID: userID, key: idempotencyKey}
	if timerID, exists := repository.idempotency[scope]; exists {
		storedTimer := repository.timers[timerID]
		if storedTimer.ScooterID != timer.ScooterID {
			return domain.Timer{}, ErrIdempotencyKeyConflict
		}

		return storedTimer, nil
	}

	repository.timers[timer.ID] = timer
	repository.idempotency[scope] = timer.ID

	return timer, nil
}

func (repository *TimerRepository) Count() int {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	return len(repository.timers)
}
