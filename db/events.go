package db

import (
	"fmt"
	"sync"

	"acronis/models"
)

type EventsDB struct {
	events sync.Map
}

// NewEvent creates a new empty event DB
func NewEvent() (*EventsDB, error) {
	p := &EventsDB{}
	return p, nil
}

// Exists checks whether an event with a give id exists
func (p *EventsDB) Exists(id string) error {
	if _, ok := p.events.Load(id); !ok {
		return fmt.Errorf("no event found for id %s", id)
	}
	return nil
}

// Find returns a given event if one exists
func (p *EventsDB) Find(id string) (models.Event, error) {
	pp, ok := p.events.Load(id)
	if !ok {
		return models.Event{}, fmt.Errorf("no event found for id %s", id)
	}

	return toEvent(pp), nil
}

// Upsert creates or updates an event in the events DB
func (p *EventsDB) Upsert(event models.Event) {
	p.events.Store(event.Id, event)
}

// FindAll returns all  event in the system
func (p *EventsDB) FindAll() []models.Event {
	var allEvents []models.Event
	p.events.Range(func(_, value interface{}) bool {
		allEvents = append(allEvents, toEvent(value))
		return true
	})

	// sort.Slice(allEvents, func(i, j int) bool {
	// 	return allEvents[i].Id < allEvents[j].Id
	// })
	return allEvents
}

// toEvent attempts to convert an interface to an event
// panics if if this not possible
func toEvent(pp interface{}) models.Event {
	prod, ok := pp.(models.Event)
	if !ok {
		panic(fmt.Errorf("error casting %v to event", pp))
	}
	return prod
}
