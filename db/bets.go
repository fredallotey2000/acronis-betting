package db

import (
	"fmt"
	"sync"

	"acronis/models"
)

type BetsDB struct {
	placedBets sync.Map
}

// NewBets creates a new empty bet service
func NewBet() *BetsDB {
	return &BetsDB{}
}

// Find bet for a given id, if one exists
func (o *BetsDB) Find(id string) (models.Bet, error) {
	po, ok := o.placedBets.Load(id)
	if !ok {
		return models.Bet{}, fmt.Errorf("no bet found for %s bet id", id)
	}

	return toBet(po), nil
}

// Upsert creates or updates a bet in the bets DB
func (o *BetsDB) Upsert(bet models.Bet) {
	o.placedBets.Store(bet.Id, bet)
}

// toBet attempts to convert an interface to a bet
// panics if this not possible
func toBet(po interface{}) models.Bet {
	bet, ok := po.(models.Bet)
	if !ok {
		panic(fmt.Errorf("error casting %v to bet", po))
	}
	return bet
}

// FindAll returns all  bets in the system
func (p *BetsDB) FindAll() []models.Bet {
	var allBets []models.Bet
	p.placedBets.Range(func(_, value interface{}) bool {
		allBets = append(allBets, toBet(value))
		return true
	})

	// sort.Slice(allBets, func(i, j int) bool {
	// 	return allBets[i].Id < allBets[j].Id
	// })
	return allBets
}
