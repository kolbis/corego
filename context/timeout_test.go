package context_test

import (
	"testing"
	"time"

	tlectx "github.com/kolbis/corego/context"
)

func TestNewTimeout(t *testing.T) {
	calc := tlectx.NewTimeoutCalculator()
	duration, _ := calc.NewTimeout()

	if duration != tlectx.MaxTimeout {
		t.Fail()
	}
}

func TestNextTimeout(t *testing.T) {
	calc := tlectx.NewTimeoutCalculator()
	duration1, deadline1 := calc.NewTimeout()
	duration2, deadline2 := calc.NextTimeout(duration1, deadline1)
	wantTimeout := time.Second * 13

	if duration2 >= wantTimeout {
		t.Fail()
	}

	if deadline2.After(deadline1) {
		t.Fail()
	}
}
