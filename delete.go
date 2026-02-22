package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

type DeleteStm struct {
	d dialect.Dialect
}

func newDelete(d dialect.Dialect) *DeleteStm { return &DeleteStm{d: d} }

func (st *DeleteStm) Build() (string, []any, error) { return "", nil, nil }

func (st *DeleteStm) MustBuild() (string, []any) {
	s, a, err := st.Build()
	if err != nil {
		panic(err)
	}
	return s, a
}

func (st *DeleteStm) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return buf, args, nil
}
