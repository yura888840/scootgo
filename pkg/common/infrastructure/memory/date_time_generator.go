package memory

import "time"

type DateTimeGenerator struct{}

func (d DateTimeGenerator) GetCurrentDate() time.Time {
	currentTime := time.Now().UTC()

	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
}

func (d DateTimeGenerator) GetCurrentTime() time.Time {
	return time.Now().UTC()
}
