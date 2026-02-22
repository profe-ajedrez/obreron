package obreron_test

import (
	"errors"
	"testing"

	"github.com/profe-ajedrez/obreron/v3"
	"github.com/profe-ajedrez/obreron/v3/dialect"
)

func TestStamentResetClearsAllFieldsIncludingDialect(t *testing.T) {
	st := obreron.AcquireStament(dialect.MySQL{})

	obreron.AppendBuff(st, 'a', 'b')
	obreron.AppendSegment(st, obreron.NewSegment(0, 2, 1, -1, 0))
	obreron.AppendParams(st, 1, 2)
	obreron.SetFlag(st, 0xFF)
	obreron.SetLastStype(st, 7)
	obreron.SetErr(st, "WHERE", obreron.ErrPlaceholderMismatch)
	obreron.ReleaseStament(st)

	if obreron.GetDialect(st) != nil {
		t.Fatalf("expected dialect=nil after reset, got=%T", obreron.GetDialect(st))
	}
	if errOb := obreron.GetErr(st); errOb != nil {
		t.Fatalf("expected err=nil after reset, got=%v", errOb)
	}
	if obreron.GetFlag(st) != 0 || obreron.GetLastStype(st) != 0 {
		t.Fatalf("expected flags/lastSType reset to 0, got flags=%d last=%d", obreron.GetFlag(st), obreron.GetLastStype(st))
	}
	if obreron.GetBuffLen(st) != 0 || obreron.GetSegsLength(st) != 0 || obreron.GetParamsLength(st) != 0 {
		t.Fatalf("expected slices cleared, got buf=%d segs=%d params=%d", obreron.GetBuffLen(st), obreron.GetSegsLength(st), obreron.GetParamsLength(st))
	}
}

func TestAcquireAlwaysAssignsDialect(t *testing.T) {
	st := obreron.AcquireStament(dialect.MySQL{})
	obreron.ReleaseStament(st)

	st2 := obreron.AcquireStament(dialect.Postgres{})
	defer obreron.ReleaseStament(st2)

	if obreron.GetDialect(st2) == nil {
		t.Fatalf("expected dialect assigned")
	}
	if obreron.GetDialect(st2).Name() != (dialect.Postgres{}).Name() {
		t.Fatalf("expected Postgres dialect, got=%s", obreron.GetDialect(st).Name())
	}
}

func TestSetErrFirstErrorWins(t *testing.T) {
	st := obreron.AcquireStament(dialect.MySQL{})
	defer obreron.ReleaseStament(st)

	obreron.SetErr(st, "WHERE", obreron.ErrPlaceholderMismatch)
	obreron.SetErr(st, "FROM", obreron.ErrEmptyFrom)

	if obreron.GetErr(st) == nil {
		t.Fatalf("expected err set")
	}
	if !errors.Is(obreron.GetErr(st), obreron.ErrPlaceholderMismatch) {
		t.Fatalf("expected first error to win (ErrPlaceholderMismatch), got=%v", obreron.GetErr(st))
	}

	var be *obreron.BuildError
	if !errors.As(obreron.GetErr(st), &be) {
		t.Fatalf("expected BuildError wrapper, got=%T", obreron.GetErr(st))
	}
	if be.Op != "WHERE" {
		t.Fatalf("expected Op=WHERE, got=%s", be.Op)
	}
}
