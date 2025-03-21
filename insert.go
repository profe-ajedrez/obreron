package obreron

import (
	"bytes"
	"slices"
	"strings"
	"unsafe"
)

// InsertStament represents an insert stament
type InsertStament struct {
	*stament
	withSelect bool
}

// Insert Returns an insert stament
//
// # Example
//
// ins := Insert().Into("client").Col("name", "'some name'").Col("value", "'somemail@mail.net'").ColIf(true, "data", "'some data'").ColIf(false, "info", 12)
//
// query, p := ins.Build()
//
// r, err := db.Exec(q, p...)
func Insert() *InsertStament {
	i := &InsertStament{pool.Get().(*stament), false}
	i.add(insertS, "INSERT", "")

	i.firstCol = true

	return i
}

// Ignore adds Ignore clause to the insert stament
//
// # Example
//
// ins := Insert().Ignore().Into("client").Col("name, value", "'some name'", "'somemail@mail.net'")
//
// query, p := ins.Build()
//
// r, err := db.Exec(q, p...)
func (in *InsertStament) Ignore() *InsertStament {
	in.add(insertS, "IGNORE", "")
	return in
}

// Into adds into clause to the insert stament
func (in *InsertStament) Into(table string) *InsertStament {
	in.add(insertS, "INTO", table)
	return in
}

// Col adds columns and values to the insert clause
//
// # Example
//
// ins := insInsert().Col("name, value", "'some name'", "'somemail@mail.net'").Into("client")
//
// query, p := ins.Build()
//
// r, err := db.Exec(q, p...)
func (in *InsertStament) Col(col string, p ...any) *InsertStament {

	if !in.firstCol {
		in.Clause(",", "")
		in.add(colsS, col, "")
		in.add(insP, "", strings.Repeat("?", len(p)), p...)
	} else {
		in.firstCol = false
		in.add(colsS, "(", col)
		pp := ""

		if l := len(p); l >= 1 {
			pp = strings.Repeat("?,", len(p)-1)
		}

		pp += "?"

		spp := unsafe.Slice(unsafe.StringData(pp), len(pp))
		spp = spp[:]
		pp = *(*string)(unsafe.Pointer(&spp))
		in.add(insP, "", pp, p...)
	}

	return in
}

// ColIf adds columns and values to the insert clause when the cond parameter is true
//
// # Example
//
// ins := insInsert().ColIf(true, "name, value", "'some name'", "'somemail@mail.net'").Into("client")
//
// query, p := ins.Build()
//
// r, err := db.Exec(q, p...)
func (in *InsertStament) ColIf(cond bool, col string, p ...any) *InsertStament {
	if cond {
		return in.Col(col, p...)
	}
	return in
}

// ColSelect is a helper method used to build insert select... staments
//
// # Example
//
// ins := Insert().Into("courses").ColSelectIf(true, "name, location, gid", Select().Col("name, location, 1").From("courses").Where("cid = 2")).ColSelectIf(false, "last_name, last_location, grid", Select().Col("last_name, last_location, 11").From("courses").Where("cid = 2"))
//
// query, p := ins.Build()
//
// r, err := db.Exec(q, p...)
//
// Produces: INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2
func (in *InsertStament) ColSelect(col string, expr *SelectStm) *InsertStament {
	in.firstCol = true
	in.add(colsS, "(", col)
	q, pp := expr.Build()
	in.add(insP, "", q, pp...)
	in.withSelect = true
	return in
}

func (in *InsertStament) ColSelectIf(cond bool, col string, expr *SelectStm) *InsertStament {
	if cond {
		return in.ColSelect(col, expr)
	}
	return in
}

func (in *InsertStament) Clause(clause, expr string, p ...any) *InsertStament {
	in.add(in.lastPos, clause, expr, p...)
	return in
}

// Build returns the query and the parameters as to be used by *sql.DB.query or *sql.DB.Exec
func (in *InsertStament) Build() (string, []any) {
	b := bytes.Buffer{}

	if in.withSelect {
		b.Grow(in.buff.Len() + 2)
	} else {
		b.Grow(in.buff.Len() + 10)
	}

	buf := in.buff.Bytes()

	slices.SortStableFunc(in.s, func(a, b segment) int {

		if a.sType < b.sType {
			return -1
		}

		if a.sType > b.sType {
			return +1
		}

		if a.sType == insP {
			return -1
		}

		return 0
	})

	i := posClauses(in, &b, buf)

	if in.withSelect {
		b.WriteString(") ")
	} else {
		b.WriteString(") VALUES ( ")
	}

	posParams(i, in, &b, buf)

	if !in.withSelect {
		b.WriteString(" )")
	}

	ss := b.Bytes()
	return *(*string)(unsafe.Pointer(&ss)), in.p
}

func posClauses(in *InsertStament, b *bytes.Buffer, buf []byte) int {
	i := 0
	for i < len(in.s) && in.s[i].sType != insP {
		k := i
		j := 0

		for k < len(in.s) && in.s[k].sType == in.s[i].sType {
			if j > 0 && j < len(in.s)-1 {
				if in.s[i].sType != colsS {
					b.WriteString(" ")
				} else {
					b.WriteString(", ")
				}
			}

			b.Write(buf[in.s[k].start : in.s[k].start+in.s[k].length])

			k++
			j++
		}

		i = k - 1

		if i < len(in.s)-1 {
			b.WriteString(" ")
		}
		i++
	}
	return i
}

func posParams(i int, in *InsertStament, b *bytes.Buffer, buf []byte) {
	for i < len(in.s) {
		k := i

		for k < len(in.s) && in.s[k].sType == in.s[i].sType {
			b.Write(buf[in.s[k].start : in.s[k].start+in.s[k].length])
			k++
		}

		i = k - 1

		i++
	}
}

// Close free resources used by the stament
func (in *InsertStament) Close() {
	CloseStament(in.stament)
}
