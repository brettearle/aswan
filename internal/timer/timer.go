package timer

import (
	"fmt"
	"strconv"
	"time"
)

// TermPomoTimer is what is used to run pomodoro timer in terminal
type TermPomoTimer struct {
	Work int
	Rest int
	Rounds int
	LongRest int
	CurrentTimers []*time.Timer
}

func (t *TermPomoTimer) Start() {
	//does it need timer parameters like system etc
	workStr := strconv.FormatInt(int64(t.Work), 10)
	duration, err := time.ParseDuration(workStr + "m")

	if err != nil {
		fmt.Printf("Problem converting %v", err)
	}
	timer := time.NewTimer(duration)
	t.CurrentTimers = append(t.CurrentTimers, timer)
	time.Sleep(duration)
	//is it work or rest or lrest
	//work starts at Work
	//rest starts at Rest
	//lrest starts at Longrest
	//set up timer
	//start process in background
	//bring forward once timer is up
}
func (t *TermPomoTimer) Stop()    {}
func (t *TermPomoTimer) Reset()   {}
func (t *TermPomoTimer) Set()     {}
func (t *TermPomoTimer) Display() {}

func NewTermPomoTimer(work int, rest int, rounds int, lngRest int ) *TermPomoTimer {
	return &TermPomoTimer{
		Work: work,
		Rest: rest,
		Rounds: rounds,
		LongRest: lngRest,
	}
}