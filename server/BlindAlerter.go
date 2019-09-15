package server

import (
	"fmt"
	"os"
	"time"
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type BlindAlerterFunc func(duration time.Duration, amount int)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

// 标准输入实现
func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		_, _ = fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}

// SpyBlindAlerter 用于测试
type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}
