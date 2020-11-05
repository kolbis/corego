package utils_test

import (
	"testing"

	"github.com/kolbis/corego/utils"
)

func TestBuildWithPrefix(t *testing.T) {
	want := "prefix-name-a-b-c"
	k := utils.NewKeys()
	is := k.Build("pRefiX", "nAme", "A", "B", "C")

	if want != is {
		t.Fail()
	}
}

func TestBuildWithoutPrefix(t *testing.T) {
	want := "name-a-b-c"
	k := utils.NewKeys()
	is := k.Build("", "nAme", "A", "B", "C")

	if want != is {
		t.Fail()
	}
}
