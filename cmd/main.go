package main

import (
	"log"
	"net/http"

	"cinema-ticket-booking/internal/adapters/redis"
	"cinema-ticket-booking/internal/utils"

	"cinema-ticket-booking/internal/booking"
)

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

func main() {
	mux := http.NewServeMux()

	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	mux.Handle("/", http.FileServer(http.Dir("static")))

	mux.HandleFunc("GET /movies", listMovies)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)

	sessionsMux := http.NewServeMux()
	mux.Handle("/sessions/", http.StripPrefix("/sessions", utils.AuthMW(sessionsMux)))
	sessionsMux.HandleFunc("PUT /{sessionID}/confirm", bookingHandler.ConfirmSession)
	sessionsMux.HandleFunc("DELETE /{sessionID}", bookingHandler.ReleaseSession)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}
