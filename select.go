package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

// SelectStm is the builder for SELECT statements.
//
// It is NOT concurrency-safe. One goroutine must own one builder instance.
// A [sync.Pool] can be used to efficiently share builders across goroutines:
// acquire → populate → build → release.
type SelectStm struct {
	st *stament
}

func newSelect(d dialect.Dialect) *SelectStm {
	return &SelectStm{st: acquireStament(d)}
}

// Build assembles the SQL statement and returns it together with a
// defensive copy of the argument slice.
//
// Signature: (string, []any, error) — stable contract per spec Rev 2.0.
//
// If any fluent method accumulated an error, Build returns it here.
// The returned args slice is always freshly allocated (nil when empty),
// so the caller may safely mutate it without affecting internal state.
func (s *SelectStm) Build() (string, []any, error) {
	buf, args, err := s.st.buildInto(nil, nil)
	if err != nil {
		return "", nil, err
	}
	return string(buf), args, nil
}

// MustBuild is a convenience wrapper that panics on error.
func (s *SelectStm) MustBuild() (string, []any) {
	sql, args, err := s.Build()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// BuildInto assembles into caller-owned buffers to avoid allocation.
// The caller is responsible for pre-slicing (buf[:0], args[:0]) between reuses.
func (s *SelectStm) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return s.st.buildInto(buf, args)
}
