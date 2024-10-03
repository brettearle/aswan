package timer

// TermPomoTimer is what is used to run pomodoro timer in terminal
type TermPomoTimer struct {
	Work int
	Rest int
	Rounds int
	LongRest int
}

func (t *TermPomoTimer) Start()   {}
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