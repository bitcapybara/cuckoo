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
	ParseOption  ParseOption
	FixedDelay   time.Duration
	FixedRate    time.Duration
	InitialDelay time.Duration
}

type JobId int64

type ExecutorId int64

type Job struct {
	Id           JobId
	ExecutorId   ExecutorId
	Comment      string
	Path         string
	ScheduleRule ScheduleRule
	RouteType    RouteType
	Timeout      time.Duration
	Router       RouteType
}
