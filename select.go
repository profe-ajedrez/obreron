// Package obreron provides a simple, fast and cheap query builder
package obreron

// SelectStm is a select stament
type SelectStm struct {
	*stament
}

// Select Returns a select stament
//
// # Example
//
// query, _ := Select().Col("a1, a2, a3").From("client").Build()
// r, error := db.Query(q)
func Select() *SelectStm {
	s := &SelectStm{
		pool.Get().(*stament),
	}

	s.add(selectS, "SELECT", "")
	return s
}

// Close release the resources used by the stament
func (st *SelectStm) Close() {
	closeStament(st.stament)
}

// Col adds a column to the select stament.
//
// # Example
//
// s := Select()
// s.Col("name, mail").Col("? AS max_credit", 1000000).
// From("client")
func (st *SelectStm) Col(expr string, p ...any) *SelectStm {
	if !st.firstCol {
		st.add(colsS, ",", expr, p...)
		return st
	}

	st.add(colsS, "", expr, p...)
	st.firstCol = false
	return st
}

// ColIf adds a column to the select stament only when `cond` parameter is true.
//
// # Example
//
// addMaxCredit := true
//
// s := Select()
// s.Col("name, mail").ColIf(addMaxCredit, "? AS max_credit", 1000000).
// From("client")
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

// From sets the source table for the select stament
//
// # Example
//
// s := Select()
// s.Col("*").
// From("client")
func (st *SelectStm) From(source string) *SelectStm {
	st.add(fromS, "FROM", source)
	return st
}

// Join adds a relation to the query in the form of an inner join
//
// # Example
//
// s := Select().Col("*").From("client").
// Join("addresses a ON a.client_id = c.client_id")
//
// # Also On clause can be used along with connectors and parameters
//
// s := Select().Col("*").From("client").
// Join("addresses a").On("a.client_id = c.client_id").And("c.status = ?", 0)
func (st *SelectStm) Join(expr string, p ...any) *SelectStm {
	st.add(joinS, "JOIN", expr, p...)
	return st
}

// JoinIf adds a relation to the query in the form of an inner join only when the cond parameter is true
//
// # Example
//
// addJoin := true
// s := Select().Col("*").
// From("client").
// JoinIf(addJoin, "addresses a ON a.client_id = c.client_id")
//
// # Also OnIf clause can be used along with connectors and parameters
//
// s := Select().Col("*").
// From("client").
// JoinIf(aaddJoin, "addresses a").
// OnIf(addJoin, "a.client_id = c.client_id").And("c.status = ?", 0)
func (st *SelectStm) JoinIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "JOIN", expr, p...)
	}
	return st
}

// LeftJoin adds a relation to the query in the form of a left join
//
// # Example
//
// s := Select().Col("*").From("client").
// LeftJoin("addresses a ON a.client_id = c.client_id")
//
// # Also On clause can be used along with connectors and parameters
//
// s := Select().Col("*").From("client").
// LeftJoin("addresses a").On("a.client_id = c.client_id").And("c.status = ?", 0)
func (st *SelectStm) LeftJoin(expr string, p ...any) *SelectStm {
	st.add(joinS, "LEFT JOIN", expr, p...)
	return st
}

// LeftJoinIf adds a relation to the query in the form of a left join only when the cond parameter is true
//
// # Example
//
// addJoin := true
// s := Select().Col("*").
// From("client").
// LeftJoinIf(addJoin, "addresses a ON a.client_id = c.client_id")
//
// # Also OnIf clause can be used along with connectors and parameters
//
// s := Select().Col("*").
// From("client").
// LeftJoinIf(aaddJoin, "addresses a").
// OnIf(addJoin, "a.client_id = c.client_id").And("c.status = ?", 0)
func (st *SelectStm) LeftJoinIf(cond bool, join string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "LEFT JOIN", join, p...)
	}
	return st
}

func (st *SelectStm) RightJoin(expr string, p ...any) *SelectStm {
	st.add(joinS, "RIGHT JOIN", expr, p...)
	return st
}

func (st *SelectStm) RightJoinIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "RIGHT JOIN", expr, p...)
	}
	return st
}

func (st *SelectStm) OuterJoin(expr string, p ...any) *SelectStm {
	st.add(joinS, "OUTER JOIN", expr, p...)
	return st
}

func (st *SelectStm) OuterJoinIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "OUTER JOIN", expr, p...)
	}
	return st
}

func (st *SelectStm) On(on string, p ...any) *SelectStm {
	st.clause("ON", on, p...)
	return st
}

func (st *SelectStm) OnIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("ON", expr, p...)
	}
	return st
}

// Where adds a condition to filter the query
//
// # Example
//
// s := Select().Col("*").From("client").
// Where("status = ?", 1)
func (st *SelectStm) Where(cond string, p ...any) *SelectStm {
	st.where(cond, p...)

	return st
}

// And adds a condition to the query connecting with an AND operator
//
// # Example
//
// s := Select().Col("*").From("client").
// Where("status = ?", 1).And("country = ?", "CL")
//
// Also can be used in join and having clauses
func (st *SelectStm) And(expr string, p ...any) *SelectStm {
	st.clause("AND", expr, p...)
	return st
}

// AndIf adds a condition to the query connecting with an AND operator only when cond parameter is true
//
// # Example
//
// filterByCountry = true
// s := Select().Col("*").From("client").
// Where("status = ?", 1).AndIf("country = ?", "CL")
//
// Also can be used in join and having clauses
func (st *SelectStm) AndIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("AND", expr, p...)
	}
	return st
}

// Y adds an AND conector to the stament where is called. Its helpful when used with In()
//
// # Example
//
//	Update("client").Set("status = 0").Where("country = ?", "CL").Y().In("status", "", 1, 2, 3, 4)
//
// Produces: UPDATE client SET status = 0 WHERE country = ? AND status IN (?, ?, ?, ?)
func (up *SelectStm) Y() *SelectStm {
	up.clause("AND", "")
	return up
}

func (st *SelectStm) Or(expr string, p ...any) *SelectStm {
	st.clause("OR", expr, p...)
	return st
}

func (st *SelectStm) OrIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("OR", expr, p...)
	}
	return st
}

// Like adds a LIKE clause to the query after the las clause added
//
// # Example
//
// Select().Col("a1, a2, a3").From("client").Where("1 = 1").And("city").Like("'%ago%'")
//
// Observe that if you use it like Select().Like(..., will produce "SELECT LIKE"
func (st *SelectStm) Like(expr string, p ...any) *SelectStm {
	st.clause("LIKE", expr, p...)
	return st
}

// LikeIf adds a LIKE clause to the query after the las clause added, when cond is true
//
// # Example
//
// Select().Col("a1, a2, a3").From("client").Where("1 = 1").And("city").LikeIf(true, "'%ago%'")
func (st *SelectStm) LikeIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("LIKE", expr, p...)
	}
	return st
}

// In adds a IN clause to the query after the las clause added
//
// # Example
//
// Select().Col("a1, a2, a3").From("client").Where("1 = 1").And("city").In("'Nagoya'", "'Tokio", "'Parral'")
func (st *SelectStm) In(expr string, p ...any) *SelectStm {
	st.clause("IN (", expr+")", p...)
	return st
}

// InArgs adds an In clause to the stament automatically setting the positional parameters of the query based on the
// passed parameters
func (up *SelectStm) InArgs(value string, p ...any) *SelectStm {
	up.stament.inArgs(value, p...)
	return up
}

// GroupBy adds a GROUP BY clause to the query
//
// # Example
//
// Select().Col("a1, a2, a3").From("client").Where("1 = 1").GroupBy("a1")
func (st *SelectStm) GroupBy(grp string, p ...any) *SelectStm {
	if !st.grouped {
		st.add(groupS, "GROUP BY", grp, p...)
		st.grouped = true
	} else {
		st.add(groupS, ",", grp, p...)
	}

	return st
}

// Having adds a HAVING clause to the query
//
// # Example
//
// Select().Col("a1, a2, a3, COUNT(1) AS how_many").From("client").Where("1 = 1").GroupBy("a1").Having(how_many > 100)
func (st *SelectStm) Having(hav string, p ...any) *SelectStm {
	st.add(havingS, "HAVING", hav, p...)
	return st
}

// OrderBy adds an ORDER BY clause to the query
//
// # Example
//
// Select().Col("a1, a2, a3").From("client").Where("1 = 1").OrderBy("a1 ASC")
func (st *SelectStm) OrderBy(expr string, p ...any) *SelectStm {
	st.add(orderS, "ORDER BY", expr, p...)
	return st
}

// Limit adds a LIMIT clause to the query
//
// # Example
//
// Select().Col("a1, a2, a3").From("client").Where("1 = 1").Limit(100)
func (st *SelectStm) Limit(limit int) *SelectStm {
	st.add(limitS, "LIMIT", "?", limit)
	return st
}

func (st *SelectStm) Offset(off int) *SelectStm {
	st.add(offsetS, "OFFSET", "?", off)
	return st
}

// Clause adds a custom clause to the query in the position were is invoked
//
// # Example
//
// Select().Clause("SQL NO CACHE").Col("a1, a2, a3").From("client").Where("1 = 1")
func (st *SelectStm) Clause(clause, expr string, p ...any) *SelectStm {
	st.add(st.lastPos, clause, expr, p...)
	return st
}

// ClauseIf adds a custom clause to the query in the position were is invoked, whencond is true
//
// # Example
//
// Select().ClauseIf(true, "SQL NO CACHE").Col("a1, a2, a3").From("client").Where("1 = 1")
func (st *SelectStm) ClauseIf(cond bool, clause, expr string, p ...any) *SelectStm {
	if cond {
		st.add(st.lastPos, clause, expr, p...)
	}
	return st
}
