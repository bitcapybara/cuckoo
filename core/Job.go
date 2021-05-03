package core

import (
	"time"
)

type ScheduleType uint8

const (
	Cron ScheduleType = iota
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
	Initial      time.Duration
	Duration     time.Duration
}

type JobId string

type Job struct {
	Id           JobId
	Group        string
	Path         string
	ScheduleRule ScheduleRule
	Timeout      time.Duration
	Router       RouteType
	Enable       bool
}
