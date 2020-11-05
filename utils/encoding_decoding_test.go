package utils_test

import (
	"testing"
	"time"

	"github.com/kolbis/corego/utils"
)

type innerStruct struct {
	ID int
}

type outerStruct struct {
	Name     string
	Inner    innerStruct
	Duration int64
}

const (
	tenSeconds time.Duration = time.Second * 10
)

func TestDecode(t *testing.T) {
	wantName := "guy kolbis"
	wantID := 555
	duration := tenSeconds.Milliseconds()
	input := map[string]interface{}{
		"Name":     wantName,
		"Duration": duration,
		"Inner": map[string]interface{}{
			"ID": wantID,
		},
	}

	var output outerStruct
	decoder := utils.NewDecoder()
	err := decoder.MapDecode(input, &output)

	if err != nil {
		t.Error(err)
	}

	if output.Name != wantName {
		t.Error("Name does not match")
	}

	if output.Inner.ID != wantID {
		t.Error("ID does not match")
	}
}
