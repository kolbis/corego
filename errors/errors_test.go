package errors_test

import (
	"errors"
	"strings"
	"testing"

	tleerrors "github.com/kolbis/corego/errors"
)

func TestApplicationError(t *testing.T) {
	format := "error from %s %s"
	err := tleerrors.NewApplicationErrorf(format, "guy", "kolbis")
	isErrString := err.Error()
	wantErrString := "error from guy kolbis application error"

	if tleerrors.IsApplicationError(err) == false {
		t.Fail()
	}

	if wantErrString != isErrString {
		t.Errorf("wanted: '%s', got: '%s'", wantErrString, isErrString)
	}
}

func TestErrorWithAnnotation(t *testing.T) {
	msg := "original error from guy kolbis"
	err := tleerrors.New(msg)
	err = tleerrors.Annotate(err, "more information about the error")

	isErrString := err.Error()
	wantErrString := "more information about the error: original error from guy kolbis application error"

	if wantErrString != isErrString {
		t.Errorf("wanted: '%s', got: '%s'", wantErrString, isErrString)
	}
}

func TestErrorWithWrap(t *testing.T) {
	msg := "original error from guy kolbis"
	err := errors.New(msg)
	errForbidden := tleerrors.NewForbiddenError(err, "forbidden!")
	newerr := tleerrors.Wrap(err, errForbidden)

	if tleerrors.IsForbiddenError(newerr) == false {
		t.Fail()
	}

	isErrString := newerr.Error()
	wantErrString := "forbidden! forbidden error: original error from guy kolbis"

	if wantErrString != isErrString {
		t.Errorf("wanted: '%s', got: '%s'", wantErrString, isErrString)
	}
}

func TestNotFoundOriginatedFromCallerFile(t *testing.T) {
	msg := "this is an error!"
	err := tleerrors.NewNotFoundErrorf(msg)
	stack := tleerrors.ErrorStack(err)
	want := "errors_test.go"

	if strings.Contains(stack, want) == false {
		t.Fail()
	}
}

func TestNew(t *testing.T) {
	msg := "new error"
	err := tleerrors.New(msg)
	stack := tleerrors.ErrorStack(err)
	want := "errors_test.go"

	if strings.Contains(stack, want) == false {
		t.Fail()
	}
}
