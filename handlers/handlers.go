package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"acronis/models"
	"acronis/service"
)

// handler holds all the dependencies required for server requests
type handler struct {
	service service.IService
	once    sync.Once
}

// Handler is the interface we expose to outside packages
type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
	AddEvent(w http.ResponseWriter, r *http.Request)
	AddBet(w http.ResponseWriter, r *http.Request)
	TotalPrizes(w http.ResponseWriter, r *http.Request)
	EndEvent(w http.ResponseWriter, r *http.Request)
	CheckWins(w http.ResponseWriter, r *http.Request)
	CloseApp(w http.ResponseWriter, r *http.Request)
}

func New() (Handler, error) {
	r, err := service.NewService()
	if err != nil {
		return nil, err
	}
	h := handler{
		service: r,
	}
	return &h, nil
}

// Index returns a simple welcome response for the homepage
func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	writeResponse(w, http.StatusOK, "Welcome to the betting App!", nil)
}

// AddEvent creates a new event with the given parameters
func (h *handler) AddEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid event body:%v", err))
		return
	}
	event_, err := h.service.NewEvent(event)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, err)
		return
	}
	writeResponse(w, http.StatusOK, event_, nil)
}

// AddBet creates a new bet with the given parameters
func (h *handler) AddBet(w http.ResponseWriter, r *http.Request) {
	var bet models.Bet
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&bet); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid bet body:%v", err))
		return
	}
	bet_, err := h.service.PlaceBet(bet)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, err)
		return
	}
	writeResponse(w, http.StatusOK, bet_, nil)
}

// TotalPrizes calculates the total prize to be given to winers
func (h *handler) TotalPrizes(w http.ResponseWriter, r *http.Request) {

	var event models.Event
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid event body:%v", err))
		return
	}

	total, err := h.service.PrizePool(event)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusNotFound, nil, err)
		return
	}
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, total, nil)
}

// EndEvent ends an event and processes winnngs
func (h *handler) EndEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid event body:%v", err))
		return
	}
	event_, err := h.service.SetEventResults(event)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, err)
		return
	}
	writeResponse(w, http.StatusOK, event_, nil)
}

// CheckWins checks the wins for an event
func (h *handler) CheckWins(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	// Read the request body
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid event body:%v", err))
		return
	}
	wins, err := h.service.CheckWins(event)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, nil, err)
		return
	}
	writeResponse(w, http.StatusOK, wins, nil)
}

// CloseApp closes the betting app for new bets
func (h *handler) CloseApp(w http.ResponseWriter, r *http.Request) {
	h.invokeClose()
	writeResponse(w, http.StatusOK, "The betting App is now closed!", nil)
}

func (h *handler) invokeClose() {
	h.once.Do(func() {
		h.service.Close()
	})
}
