package booking

import "sync"

type ConcurrentMemoryStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentMemoryStore() *ConcurrentMemoryStore {
	return &ConcurrentMemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurrentMemoryStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyTaken
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentMemoryStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()
	var bookings []Booking
	for _, value := range s.bookings {
		if value.MovieID == movieID {
			bookings = append(bookings, value)
		}
	}

	return bookings
}
