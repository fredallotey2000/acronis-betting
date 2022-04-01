package models

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventStatus_Open   EventStatus = "open"
	EventStatus_Closed EventStatus = "closed"
)

type Event struct {
	Id        string      `json:"id,omitempty"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Team1     string      `json:"team1"`
	Team2     string      `json:"team2"`
	Duration  int         `json:"duration"`
	Results   string      `json:"results,omitempty"`
	WinAmount float64     `json:"winAmount"`
	Status    EventStatus `json:"status,omitempty"`
	CreatedAt string      `json:"createdAt,omitempty"`
}

func NewEvent(event Event) Event {
	return Event{
		Id:        uuid.New().String(),
		Name:      event.Name,
		Type:      event.Type,
		Team1:     event.Team1,
		Team2:     event.Team2,
		Duration:  event.Duration,
		WinAmount: event.WinAmount,
		Status:    EventStatus_Open,
		CreatedAt: time.Now().Format(timeFormat),
	}
}
