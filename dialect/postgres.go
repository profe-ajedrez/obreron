package dialect

import "strconv"

// Postgres implements PostgreSQL dialect.
type Postgres struct{}

func (Postgres) AppendPlaceholder(dst []byte, n int) []byte {
	const (
		pHolder = 10
	)
	dst = append(dst, '$')
	// n is 0-based; PostgreSQL placeholders are 1-based.
	return strconv.AppendInt(dst, int64(n+1), pHolder)
}

func (Postgres) AppendQuotedIdentifier(dst []byte, id string) []byte {
	dst = append(dst, '"')
	for i := range len(id) {
		b := id[i]
		if b == '"' {
			dst = append(dst, '"', '"')
			continue
		}
		dst = append(dst, b)
	}
	dst = append(dst, '"')
	return dst
}

func (Postgres) Name() string { return "postgres" }

func (Postgres) SupportsReturning() bool { return true }

// MaxParams returns a conservative positional-parameter limit for PostgreSQL
// bind parameters. The extended query protocol uses uint16 for parameter count
// (max 65535).
func (Postgres) MaxParams() int { return maxParamCount }

const (
	maxParamCount = 65_535
)
