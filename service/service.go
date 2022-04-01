package service

import (
	"fmt"

	"acronis/db"
	"acronis/models"
)

// service holds all the dependencies required for operations
type service struct {
	events       *db.EventsDB
	bets         *db.BetsDB
	incomingBets chan models.Bet
	done         chan struct{}
}

// IService is the interface we expose to outside packages
type IService interface {
	NewEvent(event models.Event) (*models.Event, error)
	PlaceBet(bet models.Bet) (*models.Bet, error)
	PrizePool() (float64, error)
	SetEventResults(event models.Event) (*models.Event, error)
	CheckWinAmount(event models.Event) (float64, error)
	Close()
}

// NewService creates a new Betting service with the correct database dependencies
func NewService() (IService, error) {
	done := make(chan struct{})
	e, err := db.NewEvent()
	if err != nil {
		return nil, err
	}
	o := service{
		events:       e,
		bets:         db.NewBet(),
		incomingBets: make(chan models.Bet),
		done:         done,
	}

	// start the bet processor
	go o.processIncomingBets()
	return &o, nil
}

// PlaceBet creates a new bet
func (r *service) NewEvent(incomingEvent models.Event) (*models.Event, error) {
	// add a new event
	event := models.NewEvent(incomingEvent)
	r.events.Upsert(event)
	return &event, nil
}

// PlaceBet creates a new bet
func (r *service) PlaceBet(incomingBet models.Bet) (*models.Bet, error) {
	bet := models.NewBet(incomingBet)
	// place the bet on the incoming bets channel
	select {
	case r.incomingBets <- bet:
		return &bet, nil
	case <-r.done:
		return nil, fmt.Errorf("betting app is closed, try again later")
	}
}

// PrizePool Compute total prices won
func (r *service) PrizePool() (float64, error) {
	pricePool := 0.00
	for _, bet := range r.bets.FindAll() {
		//fmt.Println(bet.EventId, "  ", bet.Status, " ", models.BetStatus_Won)
		if bet.Status == models.BetStatus_Won {

			event, err := r.events.Find(bet.EventId)
			fmt.Println(event.Name)
			if err != nil {
				return 0.0, err
			}
			pricePool += event.WinAmount
		}
	}
	return pricePool, nil
}

// SetEventResults sets an events results
func (r *service) SetEventResults(event models.Event) (*models.Event, error) {
	//update event with results
	//incoming events model will have Status to be closed and results of the event
	if event.Results == "" {
		return nil, fmt.Errorf("cannot end the event without a reuslt")
	} else if event.Status == "" || event.Status != "closed" {
		return nil, fmt.Errorf("cannot end the event without setting the status to closed")
	}
	r.events.Upsert(event)
	r.processWins(event)
	return &event, nil
}

// CheckWinAmount checkts the win amount for an event
func (r *service) processWins(event models.Event) {

	//find bets for this event and change thier status to won
	for _, bet := range r.bets.FindAll() {
		if bet.EventId == event.Id && bet.Status == models.BetStatus_Placed {
			if bet.Prediction == event.Results {
				bet.Status = models.BetStatus_Won
				r.bets.Upsert(bet)
			}
		}
	}
}

// CheckWinAmount checkts the win amount for an event
func (r *service) CheckWinAmount(event models.Event) (float64, error) {
	totalWinAmount := 0.00
	for _, bet := range r.bets.FindAll() {
		if bet.EventId == event.Id && bet.Status == models.BetStatus_Won {
			totalWinAmount += event.WinAmount
		}
	}
	return totalWinAmount, nil
}

// Close closes the betting app for incoming bets
func (r *service) Close() {
	close(r.done)
}

func (r *service) processIncomingBets() {
	fmt.Println("Bet processing started!")
	for {
		select {
		case bet := <-r.incomingBets:
			r.processBet(bet)
			fmt.Printf("Processing bet %s completed\n", bet.Id)
		case <-r.done:
			fmt.Println("Bet processing stopped!")
			return
		}
	}
}

// process bets does not add to db if the event is closed
func (r *service) processBet(bet models.Bet) {
	// check if the event is not closed
	event, err := r.events.Find(bet.EventId)
	if err != nil {
		return
	}
	if event.Status == models.EventStatus_Closed {
		return
	}
	r.bets.Upsert(bet)
}
