package schedule

import "time"

type FixedRate struct {
	Rate uint64
	Init uint64
}

func (f FixedRate) Next(t time.Time) time.Time {
	panic("implement me")
}
