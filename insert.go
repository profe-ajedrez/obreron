package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

type InsertStatement struct {
	st *stament
}

func newInsert(d dialect.Dialect) *InsertStatement { return &InsertStatement{st: acquireStament(d)} }

// Build builds the SQL and args.
func (in *InsertStatement) Build() (string, []any, error) {
	return in.st.build()
}

// BuildInto builds into caller-owned buffers.
func (in *InsertStatement) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return in.st.buildInto(buf, args)
}

func (in *InsertStatement) MustBuild() (string, []any) {
	s, a, err := in.Build()
	if err != nil {
		panic(err)
	}
	return s, a
}
