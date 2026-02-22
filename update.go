package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

type UpdateStm struct {
	d     dialect.Dialect
	table string
}

func newUpdate(d dialect.Dialect, table string) *UpdateStm { return &UpdateStm{d: d, table: table} }

func (st *UpdateStm) Build() (string, []any, error) { return "", nil, nil }

func (st *UpdateStm) MustBuild() (string, []any) {
	s, a, err := st.Build()
	if err != nil {
		panic(err)
	}
	return s, a
}

func (st *UpdateStm) BuildInto(buf []byte, args []any) ([]byte, []any, error) {
	return buf, args, nil
}
