package obreron

// SelectStm is a SELECT query builder.
//
// A SelectStm starts with "SELECT" and accumulates clauses such as columns, FROM, JOIN,
// WHERE, GROUP BY, HAVING, ORDER BY, LIMIT and OFFSET.
type SelectStm struct {
	*stament
}

// CloseSelect releases the resources used by s and returns them to the internal pool.
//
// s must not be used after calling CloseSelect.
//
// This is equivalent to calling s.Close().
func CloseSelect(s *SelectStm) {
	CloseStament(s.stament)
}

// Select constructs a new SelectStm initialized with the SELECT keyword.
//
// Example:
//
//	st := Select().Col("a1, a2, a3").From("client")
//	defer st.Close()
//
//	query, args := st.Build()
//	rows, err := db.Query(query, args...)
func Select() *SelectStm {
	st, ok := pool.Get().(*stament)

	if !ok {
		st = &stament{}
	}

	s := &SelectStm{
		st,
	}

	s.add(selectS, "SELECT", "")

	return s
}

// Close releases the resources used by the statement and returns them to the internal pool.
//
// Do not use st after calling Close.
func (st *SelectStm) Close() {
	CloseStament(st.stament)
}

// Col adds a column expression to the SELECT clause.
//
// - On the first call, expr is appended directly after SELECT.
// - On subsequent calls, expr is appended preceded by a comma.
//
// expr is a raw SQL fragment (e.g. "name", "COUNT(1) AS n", "? AS max_credit").
// Any parameters passed via p are appended to the args list in order.
//
// Example:
//
//	st := Select().
//		Col("name, mail").
//		Col("? AS max_credit", 1000000).
//		From("client")
//	defer st.Close()
//	q, args := st.Build()
func (st *SelectStm) Col(expr string, p ...any) *SelectStm {
	if !st.firstCol {
		st.add(colsS, ",", expr, p...)
		return st
	}

	st.add(colsS, "", expr, p...)
	st.firstCol = false

	return st
}

// ColIf adds a column expression to the SELECT clause only when cond is true.
//
// Example:
//
//	addMaxCredit := true
//	st := Select().
//		Col("name, mail").
//		ColIf(addMaxCredit, "? AS max_credit", 1000000).
//		From("client")
//	defer st.Close()
func (st *SelectStm) ColIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		if !st.firstCol {
			st.add(colsS, ",", expr, p...)
			return st
		}

		st.firstCol = false
		st.add(colsS, "", expr, p...)
	}

	return st
}

// From sets the source table for the SELECT statement.
//
// source is a raw SQL fragment. Typical values are "client" or "client c".
//
// Example:
//
//	st := Select().Col("*").From("client")
//	defer st.Close()
func (st *SelectStm) From(source string) *SelectStm {
	st.add(fromS, "FROM", source)
	return st
}

// Join adds an INNER JOIN clause to the query.
//
// expr is a raw SQL fragment. You can include the ON predicate in the same fragment:
//
//	st := Select().Col("*").From("client c").
//		Join("addresses a ON a.client_id = c.client_id")
//
// Or you can build the ON predicate with On/And/Or (useful when binding args):
//
//	st := Select().Col("*").From("client c").
//		Join("addresses a").
//		On("a.client_id = c.client_id").
//		And("c.status = ?", 0)
func (st *SelectStm) Join(expr string, p ...any) *SelectStm {
	st.add(joinS, "JOIN", expr, p...)
	return st
}

// JoinIf adds an INNER JOIN clause to the query only when cond is true.
func (st *SelectStm) JoinIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "JOIN", expr, p...)
	}

	return st
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (st *SelectStm) LeftJoin(expr string, p ...any) *SelectStm {
	st.add(joinS, "LEFT JOIN", expr, p...)
	return st
}

// LeftJoinIf adds a LEFT JOIN clause to the query only when cond is true.
func (st *SelectStm) LeftJoinIf(cond bool, join string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "LEFT JOIN", join, p...)
	}

	return st
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (st *SelectStm) RightJoin(expr string, p ...any) *SelectStm {
	st.add(joinS, "RIGHT JOIN", expr, p...)
	return st
}

// RightJoinIf adds a RIGHT JOIN clause to the query only when cond is true.
func (st *SelectStm) RightJoinIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "RIGHT JOIN", expr, p...)
	}

	return st
}

// OuterJoin adds an OUTER JOIN clause to the query.
//
// Note: MySQL commonly uses LEFT/RIGHT JOIN; FULL OUTER JOIN is not supported in MySQL.
// This method simply emits "OUTER JOIN" as written.
func (st *SelectStm) OuterJoin(expr string, p ...any) *SelectStm {
	st.add(joinS, "OUTER JOIN", expr, p...)
	return st
}

// OuterJoinIf adds an OUTER JOIN clause to the query only when cond is true.
func (st *SelectStm) OuterJoinIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "OUTER JOIN", expr, p...)
	}

	return st
}

// On adds an ON predicate. It is typically used immediately after a JOIN clause.
//
// Example:
//
//	st := Select().Col("*").From("client c").
//		Join("addresses a").
//		On("a.client_id = c.client_id")
func (st *SelectStm) On(on string, p ...any) *SelectStm {
	st.clause("ON", on, p...)
	return st
}

// OnIf adds an ON predicate only when cond is true.
func (st *SelectStm) OnIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("ON", expr, p...)
	}

	return st
}

// Where adds a WHERE predicate, or appends to an existing WHERE chain.
//
// Example:
//
//	st := Select().Col("*").From("client").
//		Where("status = ?", 1)
func (st *SelectStm) Where(cond string, p ...any) *SelectStm {
	st.where(cond, p...)
	return st
}

// And appends a predicate joined with AND.
//
// It can be used after WHERE, ON, and HAVING, depending on the last clause invoked.
func (st *SelectStm) And(expr string, p ...any) *SelectStm {
	st.clause("AND", expr, p...)
	return st
}

// AndIf appends a predicate joined with AND only when cond is true.
func (st *SelectStm) AndIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("AND", expr, p...)
	}

	return st
}

// Y appends an "AND" connector without an expression.
//
// This is a convenience helper when you want to add a connector and then build a
// structured clause (like InArgs) right after it.
//
// Example:
//
//	st := Select().
//		Col("*").
//		From("client").
//		Where("country = ?", "CL").
//		Y().
//		InArgs("status", 1, 2, 3, 4)
//
// Produces:
//
//	SELECT * FROM client WHERE country = ? AND status IN (?, ?, ?, ?)
func (st *SelectStm) Y() *SelectStm {
	st.clause("AND", "")
	return st
}

// Or appends a predicate joined with OR.
func (st *SelectStm) Or(expr string, p ...any) *SelectStm {
	st.clause("OR", expr, p...)
	return st
}

// OrIf appends a predicate joined with OR only when cond is true.
func (st *SelectStm) OrIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("OR", expr, p...)
	}

	return st
}

// Like appends a LIKE clause to the last emitted predicate.
//
// Example:
//
//	st := Select().
//		Col("a1, a2, a3").
//		From("client").
//		Where("1 = 1").
//		And("city").
//		Like("'%ago%'")
//
// Note: If you call Like() immediately after Select() without a prior predicate,
// you'll produce invalid SQL such as "SELECT LIKE ...".
func (st *SelectStm) Like(expr string, p ...any) *SelectStm {
	st.clause("LIKE", expr, p...)
	return st
}

// LikeIf appends a LIKE clause only when cond is true.
func (st *SelectStm) LikeIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("LIKE", expr, p...)
	}

	return st
}

// In appends an IN clause to the last emitted predicate.
//
// Example:
//
//	st := Select().
//		Col("a1, a2, a3").
//		From("client").
//		Where("1 = 1").
//		And("city").
//		In("'Nagoya'", "'Tokio'", "'Parral'")
//
// Note: In() treats expr as raw SQL; prefer InArgs() if you want parameter binding.
func (st *SelectStm) In(expr string, p ...any) *SelectStm {
	st.clause("IN (", expr+")", p...)
	return st
}

// InArgs appends an IN clause and automatically emits positional parameters.
//
// Example:
//
//	st := Select().
//		Col("*").
//		From("client").
//		Where("country = ?", "CL").
//		And("status").
//		InArgs("status", 1, 2, 3)
//
// If p is empty, the generated clause will never match any row (e.g. "... IN (NULL)").
func (st *SelectStm) InArgs(value string, p ...any) *SelectStm {
	st.inArgs(value, p...)
	return st
}

// GroupBy adds a GROUP BY clause.
//
// Multiple calls append group keys separated by commas.
func (st *SelectStm) GroupBy(grp string, p ...any) *SelectStm {
	if !st.grouped {
		st.add(groupS, "GROUP BY", grp, p...)
		st.grouped = true
	} else {
		st.add(groupS, ",", grp, p...)
	}

	return st
}

// Having adds a HAVING clause.
//
// Example:
//
//	st := Select().
//		Col("a1, COUNT(1) AS how_many").
//		From("client").
//		GroupBy("a1").
//		Having("how_many > ?", 100)
func (st *SelectStm) Having(hav string, p ...any) *SelectStm {
	st.add(havingS, "HAVING", hav, p...)
	return st
}

// OrderBy adds an ORDER BY clause.
func (st *SelectStm) OrderBy(expr string, p ...any) *SelectStm {
	st.add(orderS, "ORDER BY", expr, p...)
	return st
}

// Limit adds a LIMIT clause (MySQL-style) using a positional parameter.
func (st *SelectStm) Limit(limit int) *SelectStm {
	st.add(limitS, "LIMIT", "?", limit)
	return st
}

// Offset adds an OFFSET clause (MySQL-style) using a positional parameter.
func (st *SelectStm) Offset(off int) *SelectStm {
	st.add(offsetS, "OFFSET", "?", off)
	return st
}

// Clause inserts a custom clause at the current position (the position of the last clause added).
//
// This is useful for vendor-specific hints.
//
// Example:
//
//	st := Select().
//		Clause("SQL_NO_CACHE", "").
//		Col("a1, a2, a3").
//		From("client")
func (st *SelectStm) Clause(clause, expr string, p ...any) *SelectStm {
	st.add(st.lastPos, clause, expr, p...)
	return st
}

// ClauseIf inserts a custom clause at the current position only when cond is true.
func (st *SelectStm) ClauseIf(cond bool, clause, expr string, p ...any) *SelectStm {
	if cond {
		st.add(st.lastPos, clause, expr, p...)
	}

	return st
}
