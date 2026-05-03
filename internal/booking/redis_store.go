package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func parseSession(val string) (Booking, error) {
	var data Booking
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}

	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil
}

func (s *RedisStore) Book(b Booking) (Booking, error) {
	session, err := s.hold(b)
	if err != nil {
		return Booking{}, err
	}

	return session, nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	var sessions []Booking

	ctx := context.Background()

	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		session, err := parseSession(val)
		if err != nil {
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions
}

func seatKey(b Booking) string {
	return fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := seatKey(b)

	b.ID = id
	val, _ := json.Marshal(b)

	err := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	}).Err()
	if err != nil {
		return Booking{}, ErrSeatAlreadyTaken
	}

	s.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL)

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    StatusHeld,
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}

func (s *RedisStore) getBookingFromSession(ctx context.Context, sessionID, userID string) (Booking, error) {
	key, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return Booking{}, err
	}

	val, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return Booking{}, err
	}

	booking, err := parseSession(val)
	if err != nil {
		return Booking{}, err
	}

	if booking.UserID != userID {
		log.Fatalf("sessioUser: %s, userID: %s", booking.UserID, userID)
		return Booking{}, ErrForbidden
	}

	return booking, nil
}

func (s *RedisStore) ConfirmSession(sessionID, userID string) error {
	ctx := context.Background()
	booking, err := s.getBookingFromSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}

	booking.Status = StatusConfimed

	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	key := seatKey(booking)
	err = s.rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStore) ReleaseSession(sessionID, userID string) error {
	ctx := context.Background()
	booking, err := s.getBookingFromSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}

	err = s.rdb.Del(ctx, seatKey(booking)).Err()
	if err != nil {
		return err
	}

	return nil
}
