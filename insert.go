package obreron

import (
	"bytes"
	"slices"
	"strings"
)

// InsertStament is an INSERT query builder.
//
// It supports:
//   - INSERT [IGNORE] INTO <table> (<cols...>) VALUES (<placeholders...>)
//   - INSERT [IGNORE] INTO <table> (<cols...>) <SELECT ...>
//
// The builder keeps SQL text and parameters separately, returning them via Build().
//
// Note: The type name is kept as-is for compatibility (typo preserved).
type InsertStament struct {
	*stament
	withSelect bool
}

// Insert constructs a new InsertStament initialized with INSERT.
//
// Example:
//
//	ins := Insert().
//		Into("client").
//		Col("name", "'some name'").
//		Col("value", "'somemail@mail.net'").
//		ColIf(true, "data", "'some data'")
//	defer ins.Close()
//
//	query, args := ins.Build()
//	_, err := db.Exec(query, args...)
func Insert() *InsertStament {
	in, ok := pool.Get().(*stament)

	if !ok {
		in = &stament{}
	}

	i := &InsertStament{in, false}
	i.add(insertS, "INSERT", "")

	i.firstCol = true

	return i
}

// Ignore adds the IGNORE modifier to the INSERT statement.
//
// MySQL will ignore certain errors (e.g., duplicate key) depending on server settings.
//
// Example:
//
//	ins := Insert().Ignore().
//		Into("client").
//		Col("name, value", "'some name'", "'somemail@mail.net'")
func (in *InsertStament) Ignore() *InsertStament {
	in.add(insertS, "IGNORE", "")

	return in
}

// Into sets the target table for the INSERT statement.
//
// table is a raw SQL fragment (e.g. "client" or "client c").
func (in *InsertStament) Into(table string) *InsertStament {
	in.add(insertS, "INTO", table)

	return in
}

// Col adds one or more columns and their corresponding values to the INSERT statement.
//
// The API accepts a comma-separated list of columns in col. The values are provided via p.
//
// Semantics:
//   - When len(p) > 0, the builder emits positional placeholders ("?") for each argument.
//   - When len(p) == 0, the builder emits DEFAULT for each column listed in col.
//     This makes it possible to explicitly rely on column defaults without causing runtime errors.
//
// Example (VALUES insert):
//
//	ins := Insert().
//		Into("client").
//		Col("name, value", "'some name'", "'somemail@mail.net'")
//
// Example (DEFAULT values):
//
//	ins := Insert().
//		Into("client").
//		Col("created_at") // emits DEFAULT
func (in *InsertStament) Col(col string, p ...any) *InsertStament {
	// If no params are provided, emit DEFAULT for each column.
	// This avoids generating placeholders without args and avoids panics on negative Repeat counts.
	if len(p) == 0 {
		def := defaultsForCols(col)

		if !in.firstCol {
			in.Clause(",", "")
			in.add(colsS, col, "")
			in.add(insP, "", def)

			return in
		}

		in.firstCol = false
		in.add(colsS, "(", col)
		in.add(insP, "", def)

		return in
	}

	ph := placeholders(len(p))

	if !in.firstCol {
		in.Clause(",", "")
		in.add(colsS, col, "")
		in.add(insP, "", ph, p...)

		return in
	}

	in.firstCol = false
	in.add(colsS, "(", col)
	in.add(insP, "", ph, p...)

	return in
}

// placeholders returns a comma-separated list of positional placeholders for n parameters.
//
// Examples:
//
//	n=1 -> "?"
//	n=3 -> "?,?,?"
func placeholders(n int) string {
	if n <= 1 {
		return "?"
	}

	var b strings.Builder
	// "?,?," ~ 2 chars per placeholder minus 1 comma.
	b.Grow(n*2 - 1)
	b.WriteByte('?')

	for i := 1; i < n; i++ {
		b.WriteString(",?")
	}

	return b.String()
}

// defaultsForCols returns a comma-separated DEFAULT list for a comma-separated column list.
//
// Example:
//
//	col="a, b, c" -> "DEFAULT,DEFAULT,DEFAULT"
func defaultsForCols(col string) string {
	n := 0
	parts := strings.Split(col, ",")

	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			n++
		}
	}

	if n <= 1 {
		return "DEFAULT"
	}

	var b strings.Builder
	// len("DEFAULT") == 7; add commas between entries.
	b.Grow(n*8 + (n - 1))
	b.WriteString("DEFAULT")

	for i := 1; i < n; i++ {
		b.WriteString(",DEFAULT")
	}

	return b.String()
}

// ColIf adds columns and values to the INSERT statement only when cond is true.
//
// It is a convenience wrapper over Col().
func (in *InsertStament) ColIf(cond bool, col string, p ...any) *InsertStament {
	if cond {
		return in.Col(col, p...)
	}

	return in
}

// ColSelect configures an INSERT ... SELECT statement.
//
// It emits:
//
//	INSERT INTO <table> (<col>) <select-query>
//
// Important:
//   - ColSelect resets the internal "first column" tracking for the insert column list.
//   - The provided SelectStm is built immediately; its args are appended to the INSERT args.
//
// Example:
//
//	ins := Insert().
//		Into("courses").
//		ColSelect("name, location, gid",
//			Select().Col("name, location, 1").From("courses").Where("cid = 2"),
//		)
//
// Produces:
//
//	INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2
func (in *InsertStament) ColSelect(col string, expr *SelectStm) *InsertStament {
	in.firstCol = true
	in.add(colsS, "(", col)

	q, pp := expr.Build()
	in.add(insP, "", q, pp...)
	in.withSelect = true

	return in
}

// ColSelectIf adds an INSERT ... SELECT clause only when cond is true.
func (in *InsertStament) ColSelectIf(cond bool, col string, expr *SelectStm) *InsertStament {
	if cond {
		return in.ColSelect(col, expr)
	}

	return in
}

// Clause appends a raw clause at the current statement position.
//
// This is an escape hatch for vendor-specific modifiers or features not covered by the API.
//
// Example:
//
//	Insert().Clause("HIGH_PRIORITY", "").Into("t").Col("a", 1)
func (in *InsertStament) Clause(clause, expr string, p ...any) *InsertStament {
	in.add(in.lastPos, clause, expr, p...)

	return in
}

// Build returns the SQL query and its positional parameters.
//
// The returned args slice is a copy of the internal parameter slice, so it remains valid
// even if you later call Close() on the builder.
//
// Usage:
//
//	ins := Insert().Into("client").Col("name", "bob")
//	defer ins.Close()
//	q, args := ins.Build()
//	_, err := db.Exec(q, args...)
func (in *InsertStament) Build() (string, []any) {
	var b bytes.Buffer

	// Heuristic growth to reduce reallocations. INSERT...SELECT and INSERT...VALUES differ slightly.
	if in.withSelect {
		const heuristic = 2
		b.Grow(in.buff.Len() + heuristic)
	} else {
		const heuristic = 10
		b.Grow(in.buff.Len() + heuristic)
	}

	buf := in.buff.Bytes()

	// Order segments by statement type so the final SQL is assembled correctly even if the user
	// calls builder methods in a flexible order.
	slices.SortStableFunc(in.s, func(a, b segment) int {
		if a.sType < b.sType {
			return -1
		}

		if a.sType > b.sType {
			return +1
		}

		// Ensure params (insP) are placed after the column list for INSERT.
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

	// Defensive copy: do not expose internal slice to callers.
	var dest []any
	if len(in.p) > 0 {
		dest = make([]any, len(in.p))
		copy(dest, in.p)
	}

	return b.String(), dest
}

// posClauses writes all non-parameter segments (keywords, INTO, column list, etc.) into b.
//
// It returns the index of the first parameter segment (insP) within in.s.
func posClauses(in *InsertStament, b *bytes.Buffer, buf []byte) int {
	i := 0
	for i < len(in.s) && in.s[i].sType != insP {
		k := i
		j := 0

		for k < len(in.s) && in.s[k].sType == in.s[i].sType {
			// Insert a separator between segments of the same type.
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

// posParams writes the parameter/value segments into b (placeholders, SELECT query, etc.).
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

// Close releases resources used by the statement and returns them to the internal pool.
//
// Do not use in after calling Close.
func (in *InsertStament) Close() {
	CloseStament(in.stament)
}
