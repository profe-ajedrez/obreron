package obreron

import (
	"errors"
	"strings"
	"testing"
)

func TestBuildError_ErrorsIs(t *testing.T) {
	err := &BuildError{Op: "FROM", Dialect: "postgres", Err: ErrEmptyFrom}
	if !errors.Is(err, ErrEmptyFrom) {
		t.Fatalf("expected errors.Is to match ErrEmptyFrom")
	}
	if errors.Is(err, ErrNoCols) {
		t.Fatalf("did not expect errors.Is to match ErrNoCols")
	}
}

func TestBuildError_ErrorsAs(t *testing.T) {
	err := &BuildError{Op: "WHERE", Dialect: "mysql", Err: ErrPlaceholderMismatch}
	var be *BuildError
	if !errors.As(err, &be) {
		t.Fatalf("expected errors.As to match *BuildError")
	}
	if be.Op != "WHERE" || be.Dialect != "mysql" {
		t.Fatalf("unexpected BuildError fields: %+v", be)
	}
}

func TestBuildError_ErrorFormat(t *testing.T) {
	err := &BuildError{Op: "UPDATE", Dialect: "sqlite", Err: ErrMissingSet}
	s := err.Error()

	// Must include op and dialect context.
	if !strings.Contains(s, "UPDATE") {
		t.Fatalf("expected Error() to include op; got %q", s)
	}
	if !strings.Contains(s, "sqlite") {
		t.Fatalf("expected Error() to include dialect; got %q", s)
	}
	// Must include underlying sentinel.
	if !strings.Contains(s, ErrMissingSet.Error()) {
		t.Fatalf("expected Error() to include underlying error; got %q", s)
	}
}
