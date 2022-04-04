package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ConfigureHandler configures the routes of this handler and binds handler functions to them
func ConfigureHandler(handler Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/").Handler(http.HandlerFunc(handler.Index))
	router.Methods("POST").Path("/events").
		Handler(http.HandlerFunc(handler.AddEvent))
	router.Methods("POST").Path("/bets").
		Handler(http.HandlerFunc(handler.AddBet))
	router.Methods("POST").Path("/totalprizes").
		Handler(http.HandlerFunc(handler.TotalPrizes))
	router.Methods("POST").Path("/endevent").Handler(http.HandlerFunc(handler.EndEvent))
	router.Methods("POST").Path("/checkwins").Handler(http.HandlerFunc(handler.CheckWins))
	router.Methods("GET").Path("/close").
		Handler(http.HandlerFunc(handler.CloseApp))

	return router
}
