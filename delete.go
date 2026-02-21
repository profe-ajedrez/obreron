package obreron

import "bytes"

// DeleteStm represents a DELETE SQL statement builder.
//
// It provides a fluent API for constructing DELETE queries and collecting positional parameters
// compatible with database/sql.
//
// Resource lifecycle:
//
//	Builders borrow internal buffers from a pool. Call Close() when you're done with the builder.
//	Do not use the builder after Close().
type DeleteStm struct {
	*stament
}

// Delete initializes and returns a new DeleteStm instance.
//
// Example:
//
//	st := Delete().From("client").Where("client_id = ?", 100)
//	defer st.Close()
//	q, args := st.Build()
//	_, err := db.Exec(q, args...)
func Delete() *DeleteStm {
	// errcheck(check-type-assertions) requires checking the type assertion result.
	// In normal operation, the pool always stores *stament because pool.New returns *stament.
	raw := pool.Get()

	st, ok := raw.(*stament)
	if !ok || st == nil {
		st = &stament{
			grouped:    false,
			firstCol:   true,
			whereAdded: false,
			lastPos:    0,
			buff:       &bytes.Buffer{},
		}
	}

	d := &DeleteStm{
		stament: st,
	}

	d.add(deleteS, "DELETE", "")

	return d
}

// From sets the target table for the DELETE statement.
func (dst *DeleteStm) From(source string) *DeleteStm {
	dst.add(fromS, "FROM", source)
	return dst
}

// Where adds a condition to filter the query.
func (dst *DeleteStm) Where(cond string, p ...any) *DeleteStm {
	dst.where(cond, p...)
	return dst
}

// Y adds an AND connector without an expression.
//
// This is useful when building compound predicates, for example when chaining In()/InArgs().
func (dst *DeleteStm) Y() *DeleteStm {
	dst.clause("AND", "")
	return dst
}

// And adds a condition connected with an AND operator.
func (dst *DeleteStm) And(expr string, p ...any) *DeleteStm {
	dst.clause("AND", expr, p...)
	return dst
}

// AndIf adds an AND-connected condition only when cond is true.
func (dst *DeleteStm) AndIf(cond bool, expr string, p ...any) *DeleteStm {
	if cond {
		dst.clause("AND", expr, p...)
	}

	return dst
}

// Or adds a condition connected with an OR operator.
//
// Like And(), this can be used after WHERE and also after ON/HAVING depending on the last clause.
func (dst *DeleteStm) Or(expr string, p ...any) *DeleteStm {
	dst.clause("OR", expr, p...)
	return dst
}

// OrIf adds an OR-connected condition only when cond is true.
func (dst *DeleteStm) OrIf(cond bool, expr string, p ...any) *DeleteStm {
	if cond {
		dst.clause("OR", expr, p...)
	}

	return dst
}

// Like adds a LIKE clause after the last predicate.
func (dst *DeleteStm) Like(expr string, p ...any) *DeleteStm {
	dst.clause("LIKE", expr, p...)
	return dst
}

// LikeIf adds a LIKE clause only when cond is true.
func (dst *DeleteStm) LikeIf(cond bool, expr string, p ...any) *DeleteStm {
	if cond {
		dst.Like(expr, p...)
	}

	return dst
}

// In adds an IN clause after the last predicate.
//
// value is the left-hand side expression, expr is a raw SQL list (without the surrounding parentheses).
// Prefer InArgs() when you want parameter binding.
func (dst *DeleteStm) In(value, expr string, p ...any) *DeleteStm {
	dst.clause(value+" IN ("+expr+")", "", p...)
	return dst
}

// InArgs adds an IN clause and automatically emits positional parameters.
//
// If p is empty, the generated clause will never match any row (e.g. "... IN (NULL)").
func (dst *DeleteStm) InArgs(value string, p ...any) *DeleteStm {
	dst.inArgs(value, p...)
	return dst
}

// OrderBy adds an ORDER BY clause to the query.
func (dst *DeleteStm) OrderBy(expr string, p ...any) *DeleteStm {
	dst.add(orderS, "ORDER BY", expr, p...)
	return dst
}

// Limit adds a LIMIT clause to the query (MySQL-style) using a positional parameter.
func (dst *DeleteStm) Limit(limit int) *DeleteStm {
	dst.add(limitS, "LIMIT", "?", limit)
	return dst
}

// Clause adds a custom clause at the current statement position.
//
// This is an escape hatch for vendor-specific modifiers or features not covered by the API.
func (dst *DeleteStm) Clause(clause, expr string, p ...any) *DeleteStm {
	dst.add(dst.lastPos, clause, expr, p...)
	return dst
}

// ClauseIf adds a custom clause at the current statement position only when cond is true.
func (dst *DeleteStm) ClauseIf(cond bool, clause, expr string, p ...any) *DeleteStm {
	if cond {
		dst.Clause(clause, expr, p...)
	}

	return dst
}

// Close releases the builder resources back to the internal pool.
//
// Do not use dst after calling Close.
func (dst *DeleteStm) Close() {
	CloseStament(dst.stament)
}
