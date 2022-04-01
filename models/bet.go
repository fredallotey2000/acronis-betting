package models

import (
	"time"

	"github.com/google/uuid"
)

type BetStatus string

const (
	BetStatus_Won      BetStatus = "Won"
	BetStatus_Lost     BetStatus = "Lost"
	BetStatus_Rejected BetStatus = "Rejected"
	BetStatus_Placed   BetStatus = "Placed"
)

const timeFormat = "2006-01-02 15:04:05.000"

type Bet struct {
	Id         string    `json:"id,omitempty"`
	EventId    string    `json:"eventId"`
	Prediction string    `json:"prediction"`
	Status     BetStatus `json:"status,omitempty"`
	Error      string    `json:"error,omitempty"`
	CreatedAt  string    `json:"createdAt,omitempty"`
}

func NewBet(bet Bet) Bet {
	return Bet{
		Id:         uuid.New().String(),
		EventId:    bet.EventId,
		Prediction: bet.Prediction,
		Status:     BetStatus_Placed,
		CreatedAt:  time.Now().Format(timeFormat),
	}
}

func (b *Bet) BetPlaced() {
	b.Status = BetStatus_Placed
}
