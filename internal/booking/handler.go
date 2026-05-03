package booking

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cinema-ticket-booking/internal/utils"
)

type handler struct {
	svc *Service
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

func NewHandler(svc *Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seat := seatInfo{
			SeatID: b.SeatID,
			UserID: b.UserID,
		}
		switch b.Status {
		case StatusHeld:
			seat.Booked = true
		case StatusConfimed:
			seat.Confirmed = true
		}

		seats = append(seats, seat)
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

func (h *handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	userID, _ := r.Context().Value(utils.ContextUserID).(string)

	err := h.svc.ReleaseSession(sessionID, userID)
	if err != nil {
		log.Fatal(err)
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	userID, _ := r.Context().Value(utils.ContextUserID).(string)

	err := h.svc.ConfirmSession(sessionID, userID)
	if err != nil {
		log.Fatal(err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
