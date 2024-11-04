package obreron

// DeleteStm represents a DELETE SQL statement builder.
// It provides a fluent interface for constructing DELETE queries.
type DeleteStm struct {
	*stament
}

func Delete() *DeleteStm {
	d := &DeleteStm{
		stament: pool.Get().(*stament),
	}

	d.add(deleteS, "DELETE", "")

	return d
}

// From sets the target table for the delete stament
//
// # Example
//
// s := Delete()
// From("client")
func (dst *DeleteStm) From(source string) *DeleteStm {
	dst.add(fromS, "FROM", source)
	return dst
}

// Where adds a condition to filter the query
func (dst *DeleteStm) Where(cond string, p ...any) *DeleteStm {
	dst.where(cond, p...)
	return dst
}

// Y adds an AND conector to the stament where is called. Its helpful when used with In()
func (dst *DeleteStm) Y() *DeleteStm {
	dst.clause("AND", "")
	return dst
}

// And adds a condition to the query connecting with an AND operator
func (dst *DeleteStm) And(expr string, p ...any) *DeleteStm {
	dst.clause("AND", expr, p...)
	return dst
}

// AndIf adds a condition to the query connecting with an AND operator only when cond parameter is true
func (dst *DeleteStm) AndIf(cond bool, expr string, p ...any) *DeleteStm {
	if cond {
		dst.clause("AND", expr, p...)
	}
	return dst
}

func (dst *DeleteStm) Or(expr string, p ...any) *DeleteStm {
	dst.clause("OR", expr, p...)
	return dst
}

func (dst *DeleteStm) OrIf(cond bool, expr string, p ...any) *DeleteStm {
	if cond {
		dst.clause("OR", expr, p...)
	}
	return dst
}

// Like adds a LIKE clause to the query after the las clause added
func (dst *DeleteStm) Like(expr string, p ...any) *DeleteStm {
	dst.clause("LIKE", expr, p...)
	return dst
}

// LikeIf adds a LIKE clause to the query after the las clause added, when cond is true
func (dst *DeleteStm) LikeIf(cond bool, expr string, p ...any) *DeleteStm {
	if cond {
		dst.Like(expr, p...)
	}
	return dst
}

// In adds a IN clause to the query after the las clause added
func (dst *DeleteStm) In(value, expr string, p ...any) *DeleteStm {
	dst.clause(value+" IN ("+expr+")", "", p...)
	return dst
}

// InArgs adds an IN clause to the statement with automatically generated positional parameters.
// Example:
//
//	Delete().From("users").Where("active = ?", true).InArgs("id", 1, 2, 3)
//
// Generates: DELETE FROM users WHERE active = ? AND id IN (?, ?, ?)
func (dst *DeleteStm) InArgs(value string, p ...any) *DeleteStm {
	dst.stament.inArgs(value, p...)
	return dst
}

// Close releases the statement back to the pool.
// After calling Close, the statement should not be used.
func (dst *DeleteStm) Close() {
	closeStament(dst.stament)
}

func (dst *DeleteStm) OrderBy(expr string, p ...any) *DeleteStm {
	dst.add(limitS, "ORDER BY", expr, p...)
	return dst
}

// Limit adds a LIMIT clause to the query
func (dst *DeleteStm) Limit(limit int) *DeleteStm {
	dst.add(limitS, "LIMIT", "?", limit)
	return dst
}

// Clause adds a custom clause to the query in the position were is invoked
func (dst *DeleteStm) Clause(clause, expr string, p ...any) *DeleteStm {
	dst.add(dst.lastPos, clause, expr, p...)
	return dst
}

func (dst *DeleteStm) ClauseIf(cond bool, clause, expr string, p ...any) *DeleteStm {
	if cond {
		dst.Clause(clause, expr, p...)
	}
	return dst
}
