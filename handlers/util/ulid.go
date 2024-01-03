package util

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func CurrentTimeAndULID(clock Clock) (time.Time, ulid.ULID, error) {
	currentTime := clock.Now()
	entropy := rand.New(rand.NewSource(currentTime.UnixNano()))
	ms := ulid.Timestamp(currentTime)
	id, err := ulid.New(ms, entropy)
	return currentTime, id, err
}
