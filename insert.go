package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

type InsertStatement struct {
	d dialect.Dialect
}

func newInsert(d dialect.Dialect) *InsertStatement { return &InsertStatement{d: d} }

func (st *InsertStatement) Build() (string, []any, error) { return "", nil, nil }

func (st *InsertStatement) MustBuild() (string, []any) {
	s, a, err := st.Build()
	if err != nil {
		panic(err)
	}
	return s, a
}

func (st *InsertStatement) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return buf, args, nil
}
