package timer

type TimeTracker interface {
	Start()
	Stop()
	Reset()
	Set()
	Display()
}

func StartTimer(t TimeTracker) {
	t.Start()
}

type TermClock struct {

}

func (t *TermClock) Start(){}
func (t *TermClock) Stop(){}
func (t *TermClock) Reset(){}
func (t *TermClock) Set(){}
func (t *TermClock) Display(){}

func NewTermClock()*TermClock{
	return &TermClock{}
}