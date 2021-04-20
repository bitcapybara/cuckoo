package schedule

import "time"

type Cron struct {
	Expr string
	Zone time.Location
}

func (c Cron) Next(time time.Time) time.Time {
	panic("implement me")
}
