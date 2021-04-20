package schedule

import "time"

type FixedDelay struct {
	Delay uint64
	Init uint64
}

func (f FixedDelay) Next(time time.Time) time.Time {
	panic("implement me")
}
