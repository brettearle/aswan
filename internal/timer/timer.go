package timer

import (
	"fmt"
	"strconv"
	"syscall"
	"time"
)

// TermPomoTimer is what is used to run pomodoro timer in terminal
type TermPomoTimer struct {
	Work          int
	Rest          int
	Rounds        int
	LongRest      int
	CurrentTimers []*time.Timer
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
	t.CurrentTimers = append(t.CurrentTimers, timer)
	// time.Sleep(duration)
	fmt.Printf("Sleep is done: %v", t)
}
func (t *TermPomoTimer) Stop() {
}
func (t *TermPomoTimer) Reset()   {}
func (t *TermPomoTimer) Set()     {}
func (t *TermPomoTimer) Display() {}

func NewTermPomoTimer(work int, rest int, rounds int, lngRest int) *TermPomoTimer {
	return &TermPomoTimer{
		Work:     work,
		Rest:     rest,
		Rounds:   rounds,
		LongRest: lngRest,
	}
}
