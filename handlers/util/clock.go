package util

import "time"

type Clock interface {
	Now() time.Time
}

func NewRealClock() RealClock {
	return RealClock{}
}

type RealClock struct{}

func (f RealClock) Now() time.Time { return time.Now() }

type FixedClock struct {
	time time.Time
}

func NewFixedClock(time time.Time) FixedClock {
	return FixedClock{
		time,
	}
}
func (f FixedClock) Now() time.Time { return f.time }
