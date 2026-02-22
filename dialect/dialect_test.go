package dialect

import (
	"testing"
)

func TestDialect_AppendPlaceholder(t *testing.T) {
	var (
		my MySQL
		pg Postgres
		sq SQLite
	)

	if got := string(my.AppendPlaceholder(nil, 0)); got != "?" {
		t.Fatalf("MySQL placeholder: got %q want %q", got, "?")
	}
	if got := string(my.AppendPlaceholder(nil, 999)); got != "?" {
		t.Fatalf("MySQL placeholder (ignores n): got %q want %q", got, "?")
	}

	if got := string(pg.AppendPlaceholder(nil, 0)); got != "$1" {
		t.Fatalf("Postgres placeholder n=0: got %q want %q", got, "$1")
	}
	if got := string(pg.AppendPlaceholder(nil, 1)); got != "$2" {
		t.Fatalf("Postgres placeholder n=1: got %q want %q", got, "$2")
	}

	if got := string(sq.AppendPlaceholder(nil, 0)); got != "?" {
		t.Fatalf("SQLite placeholder: got %q want %q", got, "?")
	}
}

func TestDialect_AppendQuotedIdentifier(t *testing.T) {
	var (
		my MySQL
		pg Postgres
		sq SQLite
	)

	// MySQL uses backticks and escapes ` as ``.
	if got := string(my.AppendQuotedIdentifier(nil, "my`table")); got != "`my``table`" {
		t.Fatalf("MySQL quoted identifier: got %q want %q", got, "`my``table`")
	}

	// Postgres/SQLite use double-quotes and escape " as "".
	if got := string(pg.AppendQuotedIdentifier(nil, "my\"table")); got != "\"my\"\"table\"" {
		t.Fatalf("Postgres quoted identifier: got %q want %q", got, "\"my\"\"table\"")
	}
	if got := string(sq.AppendQuotedIdentifier(nil, "my\"table")); got != "\"my\"\"table\"" {
		t.Fatalf("SQLite quoted identifier: got %q want %q", got, "\"my\"\"table\"")
	}
}

func TestDialect_Metadata(t *testing.T) {
	dialects := []Dialect{MySQL{}, MariaDB{}, Postgres{}, SQLite{}}
	for _, d := range dialects {
		if d.Name() == "" {
			t.Fatalf("Name must not be empty for %T", d)
		}
		// SupportsReturning is a hard contract for v3.0 target dialects.
		switch d.Name() {
		case "mysql", "mariadb":
			if d.SupportsReturning() {
				t.Fatalf("%s should not support RETURNING", d.Name())
			}
		case "postgres", "sqlite":
			if !d.SupportsReturning() {
				t.Fatalf("%s should support RETURNING", d.Name())
			}
		}
	}
}

func TestDialect_MaxParamsSemantics(t *testing.T) {
	// Contract: 0 means "unknown / do not validate".
	if got := (SQLite{}).MaxParams(); got != 0 {
		t.Fatalf("SQLite MaxParams: got %d want %d", got, 0)
	}

	// For MySQL/Postgres we expose a conservative known limit.
	if got := (MySQL{}).MaxParams(); got <= 0 {
		t.Fatalf("MySQL MaxParams must be > 0, got %d", got)
	}
	if got := (Postgres{}).MaxParams(); got <= 0 {
		t.Fatalf("Postgres MaxParams must be > 0, got %d", got)
	}
}
