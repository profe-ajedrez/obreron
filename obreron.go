// Package obreron provides a small SQL query builder for Go.
//
// It is designed to be a "cheap" builder that produces:
//   - an SQL string
//   - a slice of positional parameters (args) compatible with database/sql
//
// obreron is NOT an ORM and does not execute queries. It only builds SQL.
//
// Security note:
//
//	This package does not escape or quote identifiers (tables/columns) for you.
//	Never pass untrusted input as table/column names or raw SQL fragments.
//
// Concurrency:
//
//	Builder instances are not safe for concurrent use. Create one builder per goroutine.
//
// Resource lifecycle:
//
//	Builders borrow internal buffers from a pool. Call Close() when you're done with a
//	builder instance to release resources back to the pool. Do not use a builder after Close().
package obreron

import (
	"bytes"
	"slices"
	"strings"
	"sync"
)

var pool = &sync.Pool{
	New: func() interface{} {
		return &stament{grouped: false, firstCol: true, whereAdded: false, lastPos: 0, buff: &bytes.Buffer{}}
	},
}

// CloseStament resets and closes a Select stament
func CloseStament(st *stament) {
	if st == nil {
		return
	}

	resetStament(st)
	st.buff.Reset()
	pool.Put(st)
}

func resetStament(st *stament) {
	// Importante: []any puede retener referencias a objetos grandes.
	// clear() suelta esas referencias sin perder la capacity del slice.
	clear(st.p)
	st.p = st.p[:0]
	// segment actualmente solo tiene ints,
	// pero clear + reslice mantiene el patrón y evita sorpresas si mañana cambia.
	clear(st.s)
	st.s = st.s[:0]
	st.lastPos = 0
	st.grouped = false
	st.firstCol = true
	st.whereAdded = false
}

type segment struct {
	start, length, pCount, pIndex, sType int
}

type stament struct {
	buff       *bytes.Buffer
	s          []segment
	p          []any
	lastPos    int
	whereAdded bool
	grouped    bool
	firstCol   bool
}

func (st *stament) clause(clause, expr string, p ...any) {
	st.add(st.lastPos, clause, expr, p...)
}

func (st *stament) inArgs(value string, p ...any) {
	if len(p) == 0 {
		// Empty IN list is invalid SQL in MySQL.
		// IN (NULL) produces NULL/unknown => treated as false in WHERE/ON/HAVING filters.
		st.clause(value+" IN (NULL)", "")
		return
	}

	if len(p) == 1 {
		st.clause(value+" IN (?)", "", p...)
		return
	}

	l := len(p)

	var builder strings.Builder

	const grow = 2
	builder.Grow(l * grow) // aprox: "?, ?, ?" => 1 + 3*(l-1)
	builder.WriteString("?")

	for i := 1; i < l; i++ {
		builder.WriteString(", ?")
	}

	st.clause(value+" IN ("+builder.String()+")", "", p...)
}

func (st *stament) where(cond string, p ...any) {
	if !st.whereAdded {
		st.add(whereS, "WHERE", cond, p...)
		st.whereAdded = true
	} else {
		st.add(whereS, "AND", cond, p...)
	}
}

// Build return the query as a string with the added parameters
func (st *stament) Build() (string, []any) {
	b := bytes.Buffer{}
	b.Grow(st.buff.Len())
	buf := st.buff.Bytes()

	slices.SortStableFunc(st.s, func(a, b segment) int {
		if a.sType < b.sType {
			return -1
		}

		if a.sType > b.sType {
			return +1
		}

		return 0
	})

	dest := orderQueryAndParams(st, &b, buf)

	return b.String(), dest
}

func orderQueryAndParams(st *stament, b *bytes.Buffer, buf []byte) []any {
	dest := make([]any, len(st.p))
	first := 0

	for i := 0; i < len(st.s); i++ {
		k := i
		j := 0

		for k < len(st.s) && st.s[k].sType == st.s[i].sType {
			if j > 0 && j < len(st.s)-1 && st.s[i].sType != colsS {
				b.WriteString(" ")
			}

			b.Write(buf[st.s[k].start : st.s[k].start+st.s[k].length])

			if st.s[k].pIndex > -1 {
				copy(dest[first:], st.p[st.s[k].pIndex:st.s[k].pIndex+st.s[k].pCount])
				first += st.s[k].pCount
			}

			k++
			j++
		}

		i = k - 1

		if i < len(st.s)-1 {
			b.WriteString(" ")
		}
	}

	return dest
}

func (st *stament) add(pos int, clause, expr string, p ...any) {
	pl := len(p)
	start := st.buff.Len()

	// Remember last clause added, to be used when method Clause is invoked
	st.lastPos = pos

	if cap(st.s) == len(st.s) {
		const capSize = 2

		segments := make([]segment, len(st.s), cap(st.s)*capSize)
		copy(segments, st.s)
		st.s = segments
	}

	l := len(clause)
	_, _ = st.buff.WriteString(clause)

	if expr != "" {
		if l > 0 {
			_, _ = st.buff.WriteString(" ")
			l++
		}

		_, _ = st.buff.WriteString(expr)
		l += len(expr)
	}

	st.s = append(st.s, segment{
		start:  start,
		length: l,
		pCount: pl,
		pIndex: -1,
		sType:  pos,
	})

	if pl > 0 {
		st.s[len(st.s)-1].pIndex = len(st.p)
		st.p = append(st.p, p...)
	}
}

const (
	selectS = 0
	deleteS = 0
	updateS = 0
	insertS = 0
	colsS   = 1
	setS    = 1
	fromS   = 2
	valueS  = 2
	joinS   = 3
	whereS  = 4
	groupS  = 5
	havingS = 6
	orderS  = 7
	limitS  = 8
	offsetS = 9
	insP    = 99
)
