package obreron_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/profe-ajedrez/obreron/v3"
)

func TestBuildError_ErrorsIs(t *testing.T) {
	err := &obreron.BuildError{Op: "FROM", Dialect: "postgres", Err: obreron.ErrEmptyFrom}
	if !errors.Is(err, obreron.ErrEmptyFrom) {
		t.Fatalf("expected errors.Is to match obreron.ErrEmptyFrom")
	}
	if errors.Is(err, obreron.ErrNoCols) {
		t.Fatalf("did not expect errors.Is to match obreron.ErrNoCols")
	}
}

func TestBuildError_ErrorsAs(t *testing.T) {
	err := &obreron.BuildError{Op: "WHERE", Dialect: "mysql", Err: obreron.ErrPlaceholderMismatch}
	var be *obreron.BuildError
	if !errors.As(err, &be) {
		t.Fatalf("expected errors.As to match *obreron.BuildError")
	}
	if be.Op != "WHERE" || be.Dialect != "mysql" {
		t.Fatalf("unexpected obreron.BuildError fields: %+v", be)
	}
}

func TestBuildError_ErrorFormat(t *testing.T) {
	err := &obreron.BuildError{Op: "UPDATE", Dialect: "sqlite", Err: obreron.ErrMissingSet}
	s := err.Error()

	// Must include op and dialect context.
	if !strings.Contains(s, "UPDATE") {
		t.Fatalf("expected Error() to include op; got %q", s)
	}
	if !strings.Contains(s, "sqlite") {
		t.Fatalf("expected Error() to include dialect; got %q", s)
	}
	// Must include underlying sentinel.
	if !strings.Contains(s, obreron.ErrMissingSet.Error()) {
		t.Fatalf("expected Error() to include underlying error; got %q", s)
	}
}
