package dialect

// SQLite implements SQLite dialect.
type SQLite struct{}

func (SQLite) AppendPlaceholder(dst []byte, _ int) []byte {
	return append(dst, '?')
}

func (SQLite) AppendQuotedIdentifier(dst []byte, id string) []byte {
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

func (SQLite) Name() string { return "sqlite" }

func (SQLite) SupportsReturning() bool { return true }

func (SQLite) MaxParams() int { return 0 }
