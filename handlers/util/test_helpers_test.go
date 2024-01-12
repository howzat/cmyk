package util

import (
	"testing"
	"time"
)

func Test_TestLifespan(t *testing.T) {

	fixedDate := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	oneDay := 24 * time.Hour

	type args struct {
		lifespan Lifespan
		clock    Clock
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Short Lifespans should return 1 day in the future",
			args: args{
				lifespan: Short,
				clock:    NewFixedClock(fixedDate),
			},
			want: fixedDate.Add(oneDay).Unix(),
		}, {
			name: "When no Lifespans is given return Short",
			args: args{
				lifespan: -1,
				clock:    NewFixedClock(fixedDate),
			},
			want: fixedDate.Add(oneDay).Unix(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TestLifespan(tt.args.lifespan, tt.args.clock.Now()); got != tt.want {
				t.Errorf("TestLifespan() = %v, want %v", got, tt.want)
			}
		})
	}
}
