package booking

import (
	"encoding/json"
	"net/http"
	"time"

	"cinema-ticket-booking/internal/utils"
)

type handler struct {
	svc *Service
}

type seatInfo struct {
	SeatID string `json:"seat_id"`
	UserID string `json:"user_id"`
	Booked bool   `json:"booked"`
}

func NewHandler(svc *Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID: b.SeatID,
			UserID: b.UserID,
			Booked: true,
		})
	}

	utils.WriteJSON(w, http.StatusOK, seats)
}

func (h *handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	type holdRequest struct {
		UserID string `json:"user_id"`
	}

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	data := Booking{
		UserID:  req.UserID,
		MovieID: movieID,
		SeatID:  seatID,
	}

	session, err := h.svc.Book(data)
	if err != nil {
		return
	}

	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movieID"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	utils.WriteJSON(w, http.StatusOK, holdResponse{
		SeatID:    seatID,
		MovieID:   session.MovieID,
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}
