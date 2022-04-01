package db

import (
	"fmt"
	"sync"

	"acronis/models"
)

type EventsDB struct {
	events sync.Map
}

// NewProducts creates a new empty product DB
func NewEvent() (*EventsDB, error) {
	p := &EventsDB{}
	return p, nil
}

// Exists checks whether an event with a give id exists
func (p *EventsDB) Exists(id string) error {
	if _, ok := p.events.Load(id); !ok {
		return fmt.Errorf("no product found for id %s", id)
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

// Upsert creates or updates an event in the orders DB
func (p *EventsDB) Upsert(event models.Event) {
	p.events.Store(event.Id, event)
}

// FindAll returns all  event in the system
func (p *EventsDB) FindAll() []models.Event {
	var allProducts []models.Event
	p.events.Range(func(_, value interface{}) bool {
		allProducts = append(allProducts, toEvent(value))
		return true
	})

	// sort.Slice(allProducts, func(i, j int) bool {
	// 	return allProducts[i].Id < allProducts[j].Id
	// })
	return allProducts
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
