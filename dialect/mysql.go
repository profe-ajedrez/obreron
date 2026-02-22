package dialect

// MySQL implements the MySQL/MariaDB dialect.
type MySQL struct{}

func (MySQL) AppendPlaceholder(dst []byte, _ int) []byte {
	return append(dst, '?')
}

func (MySQL) AppendQuotedIdentifier(dst []byte, id string) []byte {
	dst = append(dst, '`')
	for i := 0; i < len(id); i++ {
		b := id[i]
		if b == '`' {
			dst = append(dst, '`', '`')
			continue
		}
		dst = append(dst, b)
	}
	dst = append(dst, '`')
	return dst
}

func (MySQL) Name() string { return "mysql" }

func (MySQL) SupportsReturning() bool { return false }

// MaxParams returns a conservative positional-parameter limit for MySQL/MariaDB
// prepared statements. This is commonly 65535.
func (MySQL) MaxParams() int { return 65_535 }

// MariaDB is an alias of MySQL behavior for v3.0.
type MariaDB struct{ MySQL }

func (MariaDB) Name() string { return "mariadb" }
