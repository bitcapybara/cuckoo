package schedule

import "time"

type Schedule interface {
	Next(time.Time) time.Time
}
