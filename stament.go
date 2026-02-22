package obreron

import (
	"math"
	"sync"

	"github.com/profe-ajedrez/obreron/v3/dialect"
)

// NOTE: stament (intentionally misspelled for continuity with v2 internals)
// is the internal core state shared by all statement builders.
//
// Builders are NOT concurrency-safe. One goroutine should use one builder.
// A sync.Pool makes concurrent allocations of separate builders efficient.

// segment describes a single SQL fragment stored inside stament.buf.
//
// Layout is intentionally compact (16 bytes on amd64) to improve cache locality
// during sorting and building.
//
// Invariants:
//   - If pCount == 0 then pIndex must be -1.
//   - pIndex is the base index into stament.params.
//   - start/length address a slice window within stament.buf.
//
// The actual assembly order is controlled by sType (stable-sorted).
//
// Total: 16 bytes.
type segment struct {
	start  uint32 // offset in buf (max 4GiB)
	length uint32 // fragment length (max 4GiB)
	pIndex int32  // base index into params; -1 when pCount == 0
	pCount uint16 // number of params used by this fragment (max 65535)
	sType  uint8  // clause/type ordering key
	_      uint8  // explicit padding / reserved
}

func newSegment(start, length uint32, sType uint8, pIndex int32, pCount uint16) segment {
	if pCount == 0 {
		pIndex = -1
	}
	return segment{
		start:  start,
		length: length,
		pIndex: pIndex,
		pCount: pCount,
		sType:  sType,
	}
}

// placeholderMarker is written into stament.buf during scanning to mark a
// positional parameter placeholder. It must never appear in generated SQL.
const placeholderMarker byte = 0x1F

// stament is the internal accumulation buffer for a single SQL statement.
//
// Invariants:
//   - If err is non-nil, all mutating operations should become no-ops.
//   - dialect may be nil only while the stament is in the pool (after reset).
type stament struct {
	buf       []byte
	segs      []segment
	params    []any
	dialect   dialect.Dialect
	err       error
	lastSType uint8
	flags     uint8
}

func (st *stament) build() (string, []any, error) {
	buf, args, err := st.buildInto(nil, nil)
	if err != nil {
		return "", nil, err
	}
	return string(buf), args, nil
}

// stamentPool reuses internal buffers to reduce GC pressure.
var stamentPool = sync.Pool{New: func() any {
	const (
		bufGrow    = 256
		segsGrow   = 12
		paramsGrow = 8
	)
	return &stament{
		buf:    make([]byte, 0, bufGrow),
		segs:   make([]segment, 0, segsGrow),
		params: make([]any, 0, paramsGrow),
	}
}}

func acquireStament(d dialect.Dialect) *stament {
	st, ok := stamentPool.Get().(*stament)

	if !ok {
		return &stament{dialect: d}
	}

	// The pool may return a previously used stament; we always set dialect.
	st.dialect = d
	return st
}

func releaseStament(st *stament) {
	st.reset()
	stamentPool.Put(st)
}

func (st *stament) reset() {
	st.buf = st.buf[:0]
	st.segs = st.segs[:0]
	st.params = st.params[:0]
	st.err = nil
	st.flags = 0
	st.lastSType = 0
	// Explicitly clear dialect to avoid retaining references.
	st.dialect = nil
}

func (st *stament) setErr(op string, err error) {
	if st.err != nil {
		return
	}
	dialectName := ""
	if st.dialect != nil {
		dialectName = st.dialect.Name()
	}
	st.err = &BuildError{Op: op, Dialect: dialectName, Err: err}
}

// addFragment scans raw for placeholders, appends raw (with markers) to buf,
// appends args to params, and records a segment.
//
// It follows the spec rules for ??/???:
//   - ?? becomes a literal '?'
//   - ? becomes placeholderMarker (0x1F)
//   - ??? is interpreted as ?? + ? (greedy-left)
//
// On error, addFragment rolls back any buffer growth and does not mutate segs
// or params.
func (st *stament) addFragment(op string, sType uint8, raw string, args ...any) {
	if st.err != nil {
		return
	}

	origLen := uint32(len(st.buf))
	if origLen > math.MaxUint16 {
		st.buf = st.buf[:origLen]
		st.setErr(op, ErrTooManyParams)
		return
	}
	start := origLen

	out, placeholders, scanErr := scanPlaceholdersInto(st.buf, raw)
	if scanErr != nil {
		st.buf = st.buf[:origLen]
		st.setErr(op, scanErr)
		return
	}

	if placeholders != len(args) {
		st.buf = st.buf[:origLen]
		st.setErr(op, ErrPlaceholderMismatch)
		return
	}
	if placeholders > int(^uint16(0)) {
		st.buf = st.buf[:origLen]
		st.setErr(op, ErrTooManyParams)
		return
	}

	// Commit buffer changes.
	st.buf = out
	length := uint32(len(st.buf)) - origLen

	if placeholders > math.MaxUint16 {
		st.buf = st.buf[:origLen]
		st.setErr(op, ErrTooManyParams)
		return
	}
	pCount := uint16(placeholders)
	pIndex := int32(-1)
	if pCount > 0 {
		pIndex = int32(len(st.params))
	}

	// Append args.
	if pCount > 0 {
		st.params = append(st.params, args...)
	}

	st.segs = append(st.segs, newSegment(start, length, sType, pIndex, pCount))
}

// scanPlaceholdersInto appends raw into dst while applying the placeholder rules.
// It returns the updated dst, the number of placeholders found, and an error.
func scanPlaceholdersInto(dst []byte, raw string) ([]byte, int, error) {
	placeholders := 0
	for i := 0; i < len(raw); {
		b := raw[i]
		if b == placeholderMarker {
			return dst, placeholders, ErrReservedMarker
		}
		if b == '?' {
			// Greedy-left: consume ?? as a literal '?'
			if i+1 < len(raw) && raw[i+1] == '?' {
				dst = append(dst, '?')
				i += 2
				continue
			}
			dst = append(dst, placeholderMarker)
			placeholders++
			i++
			continue
		}
		dst = append(dst, b)
		i++
	}
	return dst, placeholders, nil
}
