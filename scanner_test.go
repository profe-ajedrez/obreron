package obreron

import (
    "errors"
    "testing"

    "github.com/profe-ajedrez/obreron/v3/dialect"
)

func TestScanPlaceholders_ConsumesDoubleQuestionAsLiteral(t *testing.T) {
    out, n, err := scanPlaceholdersInto(nil, "data ?? 'k'")
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }
    if n != 0 {
        t.Fatalf("expected 0 placeholders, got=%d", n)
    }
    if string(out) != "data ? 'k'" {
        t.Fatalf("unexpected output: %q", string(out))
    }
}

func TestScanPlaceholders_ReplacesSingleQuestionWithMarker(t *testing.T) {
    out, n, err := scanPlaceholdersInto(nil, "x = ? AND y = ?")
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }
    if n != 2 {
        t.Fatalf("expected 2 placeholders, got=%d", n)
    }
    // Expect markers in-place.
    want := []byte("x = ")
    want = append(want, placeholderMarker)
    want = append(want, []byte(" AND y = ")...)
    want = append(want, placeholderMarker)
    if string(out) != string(want) {
        t.Fatalf("unexpected output bytes: got=%v want=%v", out, want)
    }
}

func TestScanPlaceholders_TripleQuestionIsLiteralPlusMarker(t *testing.T) {
    out, n, err := scanPlaceholdersInto(nil, "a = ???")
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }
    if n != 1 {
        t.Fatalf("expected 1 placeholder, got=%d", n)
    }
    want := []byte("a = ?")
    want = append(want, placeholderMarker)
    if string(out) != string(want) {
        t.Fatalf("unexpected output: got=%v want=%v", out, want)
    }
}

func TestScanPlaceholders_EmptyIsOk(t *testing.T) {
    out, n, err := scanPlaceholdersInto(nil, "")
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }
    if n != 0 {
        t.Fatalf("expected 0 placeholders, got=%d", n)
    }
    if len(out) != 0 {
        t.Fatalf("expected empty out, got=%q", string(out))
    }
}

func TestScanPlaceholders_ReservedMarkerIsError(t *testing.T) {
    raw := "x = " + string([]byte{placeholderMarker})
    _, _, err := scanPlaceholdersInto(nil, raw)
    if !errors.Is(err, ErrReservedMarker) {
        t.Fatalf("expected ErrReservedMarker, got=%v", err)
    }
}

func TestAddFragment_PlaceholderMismatchDoesNotMutateState(t *testing.T) {
    st := acquireStament(dialect.MySQL{})
    defer releaseStament(st)

    st.addFragment("WHERE", 1, "x = ? AND y = ?", 1) // mismatch

    if st.err == nil {
        t.Fatalf("expected err")
    }
    if !errors.Is(st.err, ErrPlaceholderMismatch) {
        t.Fatalf("expected ErrPlaceholderMismatch, got=%v", st.err)
    }
    if len(st.buf) != 0 || len(st.segs) != 0 || len(st.params) != 0 {
        t.Fatalf("expected no mutation on error, got buf=%d segs=%d params=%d", len(st.buf), len(st.segs), len(st.params))
    }
}

func TestAddFragment_ReservedMarkerDoesNotMutateState(t *testing.T) {
    st := acquireStament(dialect.MySQL{})
    defer releaseStament(st)

    raw := "x = " + string([]byte{placeholderMarker})
    st.addFragment("WHERE", 1, raw, 1)

    if st.err == nil {
        t.Fatalf("expected err")
    }
    if !errors.Is(st.err, ErrReservedMarker) {
        t.Fatalf("expected ErrReservedMarker, got=%v", st.err)
    }
    if len(st.buf) != 0 || len(st.segs) != 0 || len(st.params) != 0 {
        t.Fatalf("expected no mutation on error, got buf=%d segs=%d params=%d", len(st.buf), len(st.segs), len(st.params))
    }
}
