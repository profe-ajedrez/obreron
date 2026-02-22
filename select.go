package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

type SelectStm struct {
	d dialect.Dialect // placeholder until core is implemented
}

func newSelect(d dialect.Dialect) *SelectStm {
	return &SelectStm{d: d}
}

// Build builds the SQL and args.
// Stub: implemented in later issues.
func (st *SelectStm) Build() (string, []any, error) { return "", nil, nil }

// MustBuild panics if Build returns an error.
func (st *SelectStm) MustBuild() (string, []any) {
	s, a, err := st.Build()
	if err != nil {
		panic(err)
	}
	return s, a
}

// BuildInto builds into caller-owned buffers.
// Stub: implemented in later issues.
func (st *SelectStm) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return buf, args, nil
}
