package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

// UpdateStm is the builder for UPDATE statements.
//
// It is NOT concurrency-safe.
type UpdateStm struct {
	st *stament
}

func newUpdate(d dialect.Dialect) *UpdateStm {
	return &UpdateStm{st: acquireStament(d)}
}

// Build assembles the SQL and returns it with a defensive copy of args.
func (s *UpdateStm) Build() (string, []any, error) {
	buf, args, err := s.st.buildInto(nil, nil)
	if err != nil {
		return "", nil, err
	}
	return string(buf), args, nil
}

// MustBuild panics if Build returns an error.
func (s *UpdateStm) MustBuild() (string, []any) {
	sql, args, err := s.Build()
	if err != nil {
		panic(err)
	}
	return sql, args
}

// BuildInto assembles into caller-owned buffers to avoid allocation.
func (s *UpdateStm) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return s.st.buildInto(buf, args)
}
