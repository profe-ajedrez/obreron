package obreron

import (
	"bytes"
	"slices"
	"sync"
	"unsafe"
)

var pool = &sync.Pool{
	New: func() interface{} {
		return &stament{grouped: false, firstCol: true}
	},
}

func closeStament(st *stament) {
	for i := range st.p {
		st.p[i] = nil
	}

	st.p = st.p[:0]
	st.s = st.s[:0]
	st.lastPos = 0
	st.grouped = false
	st.firstCol = true
	st.whereAdded = false
	st.buff.Reset()

	pool.Put(st)
}

type segment struct {
	start, length, pCount, pIndex, sType int
}

type stament struct {
	s          []segment
	p          []any
	buff       bytes.Buffer
	lastPos    int
	whereAdded bool
	grouped    bool
	firstCol   bool
}

func (st *stament) clause(clause, expr string, p ...any) {
	st.add(st.lastPos, clause, expr, p...)
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

	ss := b.Bytes()
	return *(*string)(unsafe.Pointer(&ss)), dest
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
		segments := make([]segment, len(st.s), cap(st.s)*2)
		copy(segments, st.s)
		st.s = segments
	}

	l := len(clause)
	_, _ = st.buff.WriteString(clause)

	if expr != "" {
		if l > 0 {
			_, _ = st.buff.WriteString(" ")
			l += 1
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
		//st.p = insertAt(st.p, p, len(st.p))
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
