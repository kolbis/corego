package utils_test

import (
	"testing"
	"time"

	"github.com/kolbis/corego/utils"
)

func TestFromInt64ToString(t *testing.T) {
	var num int64 = 654321
	want := "654321"
	conv := utils.NewConvertor()

	is := conv.FromInt64ToString(num)

	if is != want {
		t.Fail()
	}
}

func TestFromStringToInt64(t *testing.T) {
	var num string = "654321"
	var want int64 = 654321
	conv := utils.NewConvertor()

	is := conv.FromStringToInt64(num)

	if is != want {
		t.Fail()
	}
}

func TestTimeToAndFromUnix(t *testing.T) {
	wantTime := time.Date(1977, 10, 16, 23, 0, 0, 0, time.UTC)
	conv := utils.NewConvertor()
	unix := conv.FromTimeToUnix(wantTime)
	isTime := conv.FromUnixToTime(unix)

	if isTime != wantTime {
		t.Fail()
	}
}
