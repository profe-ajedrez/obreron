package obreron

import (
	"bytes"
	"slices"
	"strings"
	"unsafe"
)

type InsertStament struct {
	*stament
	withSelect bool
}

func Insert() *InsertStament {
	i := &InsertStament{pool.New().(*stament), false}
	i.add(insertS, "INSERT", "")

	return i
}

func (in *InsertStament) Ignore() *InsertStament {
	in.add(insertS, "IGNORE", "")
	return in
}

func (in *InsertStament) Into(table string) *InsertStament {
	in.add(insertS, "INTO", table)
	return in
}

func (in *InsertStament) Col(col string, p ...any) *InsertStament {

	if !in.firstCol {
		in.Clause(",", "")
		in.add(colsS, col, "")
		in.add(insP, "", strings.Repeat(", ?", len(p)), p...)
	} else {
		in.firstCol = true
		in.add(colsS, "(", col)
		pp := strings.Repeat(" ?,", len(p))
		spp := unsafe.Slice(unsafe.StringData(pp), len(pp))
		spp = spp[:len(spp)-1]
		pp = *(*string)(unsafe.Pointer(&spp))
		in.add(insP, "", pp, p...)
	}

	return in
}

func (in *InsertStament) ColIf(cond bool, col string, p ...any) *InsertStament {
	if cond {
		return in.Col(col, p...)
	}
	return in
}

func (in *InsertStament) ColSelect(col string, expr *SelectStament) *InsertStament {
	if !in.firstCol {
		in.Clause(",", "")
		in.add(colsS, col, "")
		q, pp := expr.Build()
		in.add(insP, "", q, pp...)
	} else {
		in.firstCol = true
		in.add(colsS, "(", col)
		q, pp := expr.Build()
		in.add(insP, "", q, pp...)
		in.withSelect = true
	}
	return in
}

func (in *InsertStament) ColSelectIf(cond bool, col string, expr *SelectStament) *InsertStament {
	if cond {
		return in.ColSelect(col, expr)
	}
	return in
}

func (in *InsertStament) Clause(clause, expr string, p ...any) *InsertStament {
	in.add(in.lastPos, clause, expr, p...)
	return in
}

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

	dest := make([]any, len(in.p))
	first := 0

	i := 0
	for i < len(in.s) && in.s[i].sType != insP {
		k := i
		j := 0

		if in.s[i].sType == insP {
			break
		}

		for k < len(in.s) && in.s[k].sType == in.s[i].sType {
			if j > 0 && j < len(in.s)-1 && in.s[i].sType != colsS {
				b.WriteString(" ")
			}

			b.Write(buf[in.s[k].start : in.s[k].start+in.s[k].length])

			if in.s[k].pIndex > -1 {
				copy(dest[first:], in.p[in.s[k].pIndex:in.s[k].pIndex+in.s[k].pCount])
				first += in.s[k].pCount
			}

			k++
			j++
		}

		i = k - 1

		if i < len(in.s)-1 {
			b.WriteString(" ")
		}
		i++
	}

	if in.withSelect {
		b.WriteString(") ")
	} else {
		b.WriteString(") VALUES (")
	}

	for i < len(in.s) {
		k := i
		j := 0

		for k < len(in.s) && in.s[k].sType == in.s[i].sType {
			if j > 0 && j < len(in.s)-1 && in.s[i].sType != colsS {
				b.WriteString(" ")
			}

			b.Write(buf[in.s[k].start : in.s[k].start+in.s[k].length])

			if in.s[k].pIndex > -1 {
				copy(dest[first:], in.p[in.s[k].pIndex:in.s[k].pIndex+in.s[k].pCount])
				first += in.s[k].pCount
			}

			k++
			j++
		}

		i = k - 1

		if i < len(in.s)-1 {
			b.WriteString(" ")
		}

		i++
	}

	if !in.withSelect {
		b.WriteString(" )")
	}

	ss := b.Bytes()
	return *(*string)(unsafe.Pointer(&ss)), dest
}

func (in *InsertStament) Close() {
	Close(in.stament)
}
