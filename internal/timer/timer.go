package timer

import (
	"fmt"
	"strconv"
	"time"
)

type SetTimer struct {
	job   string
	start time.Time
	timer *time.Timer
}

func NewSetTimer(j string, t *time.Timer) *SetTimer {
	return &SetTimer{
		job:   j,
		start: time.Now(),
		timer: t,
	}
}

// TermPomoTimer is what is used to run pomodoro timer in terminal
type TermPomoTimer struct {
	Work          int
	Rest          int
	Rounds        int
	LongRest      int
	CurrentTimers []*SetTimer
}

func (t *TermPomoTimer) Start() {
	workStr := strconv.FormatInt(int64(t.Work), 10)
	duration, err := time.ParseDuration(workStr + "m")
	if err != nil {
		fmt.Printf("Problem converting %v", err)
	}

	timer := time.NewTimer(duration)

	//  TODO: make this not just work
	res := NewSetTimer("work", timer)
	t.CurrentTimers = append(t.CurrentTimers, res)
}

func (t *TermPomoTimer) KillAllTimers() {
	if len(t.CurrentTimers) > 1 {
		for _, ti := range t.CurrentTimers {
			ti.timer.Stop()
		}
	}
}

func (t *TermPomoTimer) Stop(jb string) {
	var tm *SetTimer
	for _, ct := range t.CurrentTimers {
		if ct.job == jb {
			tm = ct
		}
	}
	tm.timer.Stop()
}

func (t *TermPomoTimer) Display() {
	for _, ct := range t.CurrentTimers {
		diff := time.Since(ct.start)
		diffStr := diff.String()
		fmt.Printf("Time: %v", diffStr)
	}

}

func NewTermPomoTimer(work int, rest int, rounds int, lngRest int) *TermPomoTimer {
	return &TermPomoTimer{
		Work:     work,
		Rest:     rest,
		Rounds:   rounds,
		LongRest: lngRest,
	}
}
