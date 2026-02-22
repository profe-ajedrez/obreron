package obreron

import "github.com/profe-ajedrez/obreron/v3/dialect"

// export_test.go exposes internal symbols needed by the obreron_test package.
// This file is compiled ONLY during `go test`; it is invisible to importers.

// sType constants re-exported for black-box tests.
const (
	STypeKeyword    = sTypeKeyword
	STypeColumnList = sTypeColumnList
	STypeTarget     = sTypeTarget
	STypeJoin       = sTypeJoin
	STypeSet        = sTypeSet
	STypeWhere      = sTypeWhere
	STypeGroupBy    = sTypeGroupBy
	STypeHaving     = sTypeHaving
	STypeOrderBy    = sTypeOrderBy
	STypeLimit      = sTypeLimit
	STypeOffset     = sTypeOffset
	STypeReturning  = sTypeReturning
)

// AddFrag calls the internal addFragment on SelectStm's stament.
// It allows build_test.go to inject raw fragments without needing
// the fluent API (which is implemented in later issues).
func (s *SelectStm) AddFrag(op string, sType uint8, raw string, args ...any) *SelectStm {
	s.st.addFragment(op, sType, raw, args...)
	return s
}

// AddFrag for InsertStatement.
func (s *InsertStatement) AddFrag(op string, sType uint8, raw string, args ...any) *InsertStatement {
	s.st.addFragment(op, sType, raw, args...)
	return s
}

// AddFrag for UpdateStm.
func (s *UpdateStm) AddFrag(op string, sType uint8, raw string, args ...any) *UpdateStm {
	s.st.addFragment(op, sType, raw, args...)
	return s
}

// AddFrag for DeleteStm.
func (s *DeleteStm) AddFrag(op string, sType uint8, raw string, args ...any) *DeleteStm {
	s.st.addFragment(op, sType, raw, args...)
	return s
}

// ForceErr injects an error into SelectStm's stament (for MustBuild panic tests).
func (s *SelectStm) ForceErr(op string, err error) *SelectStm {
	s.st.setErr(op, err)
	return s
}

func ReleaseStament(st *stament) {
	releaseStament(st)
}

func AcquireStament(d dialect.Dialect) *stament {
	return acquireStament(d)
}

func AppendBuff(st *stament, strs ...byte) {
	st.buf = append(st.buf, strs...)
}

func AppendSegment(st *stament, segs ...segment) {
	st.segs = append(st.segs, segs...)
}

func AppendParams(st *stament, params ...any) {
	st.params = append(st.params, params...)
}

func NewSegment(start, length uint32, sType uint8, pIndex int32, pCount uint16) segment {
	return newSegment(start, length, sType, pIndex, pCount)
}

func SetFlag(st *stament, flags uint8) {
	st.flags = flags
}

func GetFlag(st *stament) uint8 {
	return st.flags
}

func SetLastStype(st *stament, sType uint8) {
	st.lastSType = sType
}

func GetLastStype(st *stament) uint8 {
	return st.lastSType
}

func SetErr(st *stament, op string, err error) {
	st.setErr(op, err)
}

func GetErr(st *stament) error {
	return st.err
}

func GetDialect(st *stament) dialect.Dialect {
	return st.dialect
}

func GetBuffLen(st *stament) int {
	return len(st.buf)
}

func GetSegsLength(st *stament) int {
	return len(st.segs)
}

func GetParamsLength(st *stament) int {
	return len(st.params)
}

func GetParamIndex(s segment) int32 {
	return s.pIndex
}

func GetParamsCount(s segment) uint16 {
	return s.pCount
}

func NewEmptySegment() segment {
	return segment{}
}
