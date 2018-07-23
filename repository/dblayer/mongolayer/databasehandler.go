package mongolayer

import (
	models "bookmyevents/lib"
)

// DatabaseHandler defines repository methods
type DatabaseHandler interface {
	AddEvent(models.Event) ([]byte, error)
	FindEvent([]byte) (models.Event, error)
	FindEventByName(string) (models.Event, error)
	FindAllAvailableEvents() ([]models.Event, error)
}
