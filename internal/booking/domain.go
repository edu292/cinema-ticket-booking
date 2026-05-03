package booking

import (
	"errors"
	"time"
)

var ErrSeatAlreadyTaken = errors.New("seat already taken.")

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
}
