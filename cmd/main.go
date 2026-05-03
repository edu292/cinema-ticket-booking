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

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	mux.HandleFunc("GET /movies", listMovies)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)

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
