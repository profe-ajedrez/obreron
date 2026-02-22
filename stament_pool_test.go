package obreron

import (
    "errors"
    "testing"

    "github.com/profe-ajedrez/obreron/v3/dialect"
)

func TestStamentResetClearsAllFieldsIncludingDialect(t *testing.T) {
    st := acquireStament(dialect.MySQL{})
    // Mutate state.
    st.buf = append(st.buf, 'a', 'b')
    st.segs = append(st.segs, newSegment(0, 2, 1, -1, 0))
    st.params = append(st.params, 1, 2)
    st.flags = 0xFF
    st.lastSType = 7
    st.setErr("WHERE", ErrPlaceholderMismatch)

    releaseStament(st)

    if st.dialect != nil {
        t.Fatalf("expected dialect=nil after reset, got=%T", st.dialect)
    }
    if st.err != nil {
        t.Fatalf("expected err=nil after reset, got=%v", st.err)
    }
    if st.flags != 0 || st.lastSType != 0 {
        t.Fatalf("expected flags/lastSType reset to 0, got flags=%d last=%d", st.flags, st.lastSType)
    }
    if len(st.buf) != 0 || len(st.segs) != 0 || len(st.params) != 0 {
        t.Fatalf("expected slices cleared, got buf=%d segs=%d params=%d", len(st.buf), len(st.segs), len(st.params))
    }
}

func TestAcquireAlwaysAssignsDialect(t *testing.T) {
    st := acquireStament(dialect.MySQL{})
    releaseStament(st)

    st2 := acquireStament(dialect.Postgres{})
    defer releaseStament(st2)

    if st2.dialect == nil {
        t.Fatalf("expected dialect assigned")
    }
    if st2.dialect.Name() != (dialect.Postgres{}).Name() {
        t.Fatalf("expected Postgres dialect, got=%s", st2.dialect.Name())
    }
}

func TestSetErrFirstErrorWins(t *testing.T) {
    st := acquireStament(dialect.MySQL{})
    defer releaseStament(st)

    st.setErr("WHERE", ErrPlaceholderMismatch)
    st.setErr("FROM", ErrEmptyFrom)

    if st.err == nil {
        t.Fatalf("expected err set")
    }
    if !errors.Is(st.err, ErrPlaceholderMismatch) {
        t.Fatalf("expected first error to win (ErrPlaceholderMismatch), got=%v", st.err)
    }
    var be *BuildError
    if !errors.As(st.err, &be) {
        t.Fatalf("expected BuildError wrapper, got=%T", st.err)
    }
    if be.Op != "WHERE" {
        t.Fatalf("expected Op=WHERE, got=%s", be.Op)
    }
}
