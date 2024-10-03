package timer

import "testing"

func TestTimer(t *testing.T) {
	t.Run("Constructor should return TermPomoTimer", func(t *testing.T) {
		got := NewTermPomoTimer(15, 5, 8, 15)
		want := &TermPomoTimer{
			Work: 15,
			Rest: 5,
			Rounds: 8,
			LongRest: 15,
		}
		if got != want {
			t.Errorf("got %v wanted %v", got, want)
		}
	})	
}