package core

import (
	"time"
)

type ScheduleType uint8

const (
	Cron ScheduleType = iota
	FixedDelay
	FixedRate
)

type RouteType uint8

const (
	First RouteType = iota
	Last
	Round
	Random
)

type ScheduleRule struct {
	ScheduleType ScheduleType
	CronExpr     string
	TimeZone     *time.Location
	FixedDelay   int64
	FixedRate    int64
	InitialDelay int64
}

type JobId int64

type Job struct {
	Id           JobId
	Comment      string
	Path         string
	ScheduleRule ScheduleRule
	RouteType    RouteType
	Timeout      time.Duration
	router       RouteType
}
