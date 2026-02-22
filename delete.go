package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

type DeleteStm struct {
	st *stament
}

func newDelete(d dialect.Dialect) *DeleteStm { return &DeleteStm{st: acquireStament(d)} }

// Build builds the SQL and args.
func (ds *DeleteStm) Build() (string, []any, error) {
	return ds.st.build()
}

// BuildInto builds into caller-owned buffers.
func (ds *DeleteStm) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return ds.st.buildInto(buf, args)
}

func (ds *DeleteStm) MustBuild() (string, []any) {
	s, a, err := ds.Build()
	if err != nil {
		panic(err)
	}
	return s, a
}
