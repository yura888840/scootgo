package memory_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/yura888840/scootgo/internal/timer/application"
	"github.com/yura888840/scootgo/internal/timer/infrastructure/memory"
)

type sequenceIDGenerator struct {
	next atomic.Uint64
}

func (generator *sequenceIDGenerator) NewID() string {
	return fmt.Sprintf("timer-%d", generator.next.Add(1))
}

type fixedClock struct {
	now time.Time
}

func (clock fixedClock) Now() time.Time {
	return clock.now
}

func TestStartTimerConcurrentExecution(t *testing.T) {
	tests := []struct {
		name           string
		commands       func() []application.StartTimerCommand
		expectedTimers int
		expectSameID   bool
	}{
		{
			name: "same idempotency key creates one timer",
			commands: func() []application.StartTimerCommand {
				commands := make([]application.StartTimerCommand, 20)
				for index := range commands {
					commands[index] = application.StartTimerCommand{
						UserID:         "user-1",
						ScooterID:      "scooter-1",
						IdempotencyKey: "request-1",
					}
				}

				return commands
			},
			expectedTimers: 1,
			expectSameID:   true,
		},
		{
			name: "different idempotency keys create different timers",
			commands: func() []application.StartTimerCommand {
				commands := make([]application.StartTimerCommand, 20)
				for index := range commands {
					commands[index] = application.StartTimerCommand{
						UserID:         "user-1",
						ScooterID:      "scooter-1",
						IdempotencyKey: fmt.Sprintf("request-%d", index),
					}
				}

				return commands
			},
			expectedTimers: 20,
			expectSameID:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := memory.NewTimerRepository()
			useCase := application.NewStartTimer(
				repository,
				&sequenceIDGenerator{},
				fixedClock{now: time.Date(2026, time.September, 20, 17, 0, 0, 0, time.UTC)},
			)
			commands := test.commands()
			results := make([]string, len(commands))
			errors := make([]error, len(commands))

			var waitGroup sync.WaitGroup
			waitGroup.Add(len(commands))
			for index, command := range commands {
				go func() {
					defer waitGroup.Done()

					timer, err := useCase.Execute(context.Background(), command)
					results[index] = timer.ID
					errors[index] = err
				}()
			}
			waitGroup.Wait()

			for _, err := range errors {
				require.NoError(t, err)
			}
			require.Equal(t, test.expectedTimers, repository.Count())

			uniqueIDs := make(map[string]struct{}, len(results))
			for _, timerID := range results {
				uniqueIDs[timerID] = struct{}{}
			}
			if test.expectSameID {
				require.Len(t, uniqueIDs, 1)
			} else {
				require.Len(t, uniqueIDs, len(commands))
			}
		})
	}
}
