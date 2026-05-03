package booking

import (
	"errors"
	"time"
)

var (
	ErrSeatAlreadyTaken = errors.New("seat already taken.")
	ErrForbidden        = errors.New("operation on seat that is not own by user.")
)

type BookingStatus int

const (
	StatusHeld BookingStatus = iota
	StatusConfimed
)

type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    BookingStatus
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	ConfirmSession(sessionID, userID string) error
	ReleaseSession(sessionID, userID string) error
}
