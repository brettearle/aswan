package timer

import (
	"fmt"
	"strconv"
	"syscall"
	"time"
)

type SetTimer struct {
	job   string
	timer *time.Timer
}

func NewSetTimer(j string, t *time.Timer) *SetTimer {
	return &SetTimer{
		job:   j,
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
	//maybe syscall aswan to start some other process or a instance of a timer server style
	syscall.Exec("ls", []string{"ls"}, []string{"ls"})
	workStr := strconv.FormatInt(int64(t.Work), 10)
	duration, err := time.ParseDuration(workStr + "m")
	if err != nil {
		fmt.Printf("Problem converting %v", err)
	}

	timer := time.NewTimer(duration)
	res := NewSetTimer("work", timer)
	t.CurrentTimers = append(t.CurrentTimers, res)
	fmt.Printf("Sleep is done: %v", t)
}
func (t *TermPomoTimer) KillAllTimers() {
	if len(t.CurrentTimers) > 1 {
		for _, ti := range t.CurrentTimers {
			ti.timer.Stop()
		}
	}
}

func (t *TermPomoTimer) Stop() {

}

func (t *TermPomoTimer) Display() {}

func NewTermPomoTimer(work int, rest int, rounds int, lngRest int) *TermPomoTimer {
	return &TermPomoTimer{
		Work:     work,
		Rest:     rest,
		Rounds:   rounds,
		LongRest: lngRest,
	}
}
