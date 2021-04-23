package schedule

import (
	"errors"
	"github.com/bitcapybara/cuckoo/core"
	"github.com/emirpasic/gods/sets/treeset"
	"strings"
	"time"
)

const (
	second int = iota
	minute
	hour
	dayOfMonth
	month
	dayOfWeek
	year
	noSpecInt  = 98 // '?'
	allSpecInt = 99 // '*'
)

const (
	noSpec  = noSpecInt
	allSpec = allSpecInt
)

const maxYear = 9999

var monthMap = map[string]int{
	"JAN": 0,
	"FEB": 1,
	"MAR": 2,
	"APR": 3,
	"MAY": 4,
	"JUN": 5,
	"JUL": 6,
	"AUG": 7,
	"SEP": 8,
	"OCT": 9,
	"NOV": 10,
	"DEC": 11,
}

var dayMap = map[string]int{
	"SUN": 1,
	"MON": 2,
	"TUE": 3,
	"WED": 4,
	"THU": 5,
	"FRI": 6,
	"SAT": 7,
}

type scheduleCron struct {
	expression  string
	timezone    *time.Location
	seconds     *treeset.Set
	minutes     *treeset.Set
	hours       *treeset.Set
	daysOfMonth *treeset.Set
	months      *treeset.Set
	daysOfWeek  *treeset.Set
	years       *treeset.Set

	nthDayOfWeek  int
	lastDayOffset int

	lastDayOfWeek    bool
	lastDayOfMonth   bool
	nearestWeekday   bool
	expressionParsed bool
}

func newScheduleCron(rule core.ScheduleRule) *scheduleCron {
	if rule.CronExpr == "" {
		panic("cron 表达式不可为空！")
	}
	expression := strings.ToUpper(rule.CronExpr)
	timeZone := rule.TimeZone
	if timeZone == nil {
		timeZone = time.FixedZone(time.Now().Zone())
	}
	cron := &scheduleCron{
		expression: expression,
	}

	if buildErr := cron.buildExpression(expression); buildErr != nil {
		panic("表达式解析失败：" + buildErr.Error())
	}
	return cron
}

func (s *scheduleCron) Next(time time.Time) time.Time {
	panic("implement me")
}

func (s *scheduleCron) buildExpression(expression string) error {
	s.expressionParsed = true
	if s.seconds == nil {
		s.seconds = treeset.NewWithIntComparator()
	}
	if s.minutes == nil {
		s.minutes = treeset.NewWithIntComparator()
	}
	if s.hours == nil {
		s.hours = treeset.NewWithIntComparator()
	}
	if s.daysOfMonth == nil {
		s.daysOfMonth = treeset.NewWithIntComparator()
	}
	if s.months == nil {
		s.months = treeset.NewWithIntComparator()
	}
	if s.daysOfWeek == nil {
		s.daysOfWeek = treeset.NewWithIntComparator()
	}
	if s.years == nil {
		s.years = treeset.NewWithIntComparator()
	}

	exprOn := second
	exprSplit := strings.Split(expression, "\\s+")
	for i := 0; i < len(exprSplit) && exprOn <= year; exprOn++ {
		expr := exprSplit[i]

		if exprOn == dayOfMonth && strings.Contains(expr, "L") && len(expr) > 1 && strings.Contains(expr, ",") {
			return errors.New("support for specifying 'L' and 'LW' with other days of the month is not implemented")
		}
		if exprOn == dayOfWeek && strings.Contains(expr, "L") && len(expr) > 1 && strings.Contains(expr, ",") {
			return errors.New("support for specifying 'L' with other days of the week is not implemented")
		}
		if exprOn == dayOfWeek && strings.Contains(expr, "#") && strings.Contains(expr[strings.Index(expr, "#"):], "#") {
			return errors.New("support for specifying multiple 'nth' days is not implemented")
		}

		vSplit := strings.Split(expr, ",")
		for _, v := range vSplit {
			s.storeExpressionVals(v, exprOn)
		}
	}

	if exprOn <= dayOfWeek {
		return errors.New("unexpected end of expression")
	}

	if exprOn <= year {
		s.storeExpressionVals("*", year)
	}

	dow := s.getSet(dayOfWeek)
	dom := s.getSet(dayOfMonth)

	dayOfWSpec := !dow.Contains(noSpec)
	dayOfMSpec := !dom.Contains(noSpec)

	if (!dayOfMSpec || dayOfWSpec) && (!dayOfWSpec || dayOfMSpec) {
		return errors.New("support for specifying both a day-of-week AND a day-of-month parameter is not implemented")
	}

	return nil
}

func (s *scheduleCron) storeExpressionVals(str string, on int) {

}

func (s *scheduleCron) getSet(on int) *treeset.Set {
	switch on {
	case second:
		return s.seconds
	case minute:
		return s.minutes
	case hour:
		return s.hours
	case dayOfMonth:
		return s.daysOfMonth
	case month:
		return s.months
	case dayOfWeek:
		return s.daysOfWeek
	case year:
		return s.years
	default:
		return nil
	}
}
