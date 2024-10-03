// Package obreron provides a simple, fast and cheap query builder
package obreron

// SelectStm is a select stament
type SelectStm struct {
	*stament
}

// Select Returns a select stament
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

func (st *SelectStm) Col(expr string, p ...any) *SelectStm {
	if !st.firstCol {
		st.add(colsS, ",", expr, p...)
		return st
	}

	st.add(colsS, "", expr, p...)
	st.firstCol = false
	return st
}

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

func (st *SelectStm) From(source string) *SelectStm {
	st.add(fromS, "FROM", source)
	return st
}

func (st *SelectStm) Join(expr string, p ...any) *SelectStm {
	st.add(joinS, "JOIN", expr, p...)
	return st
}

func (st *SelectStm) JoinIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.add(joinS, "JOIN", expr, p...)
	}
	return st
}

func (st *SelectStm) LeftJoin(expr string, p ...any) *SelectStm {
	st.add(joinS, "LEFT JOIN", expr, p...)
	return st
}

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

func (st *SelectStm) Where(cond string, p ...any) *SelectStm {
	st.where(cond, p...)

	return st
}

func (st *SelectStm) And(expr string, p ...any) *SelectStm {
	st.clause("AND", expr, p...)
	return st
}

func (st *SelectStm) AndIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("AND", expr, p...)
	}
	return st
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

func (st *SelectStm) Like(expr string, p ...any) *SelectStm {
	st.clause("LIKE", expr, p...)
	return st
}

func (st *SelectStm) LikeIf(cond bool, expr string, p ...any) *SelectStm {
	if cond {
		st.clause("LIKE", expr, p...)
	}
	return st
}

func (st *SelectStm) In(expr string, p ...any) *SelectStm {
	st.clause("IN (", expr+")", p...)
	return st
}

func (st *SelectStm) GroupBy(grp string, p ...any) *SelectStm {
	if !st.grouped {
		st.add(groupS, "GROUP BY", grp, p...)
		st.grouped = true
	} else {
		st.add(groupS, ",", grp, p...)
	}

	return st
}

func (st *SelectStm) Having(hav string, p ...any) *SelectStm {
	st.add(havingS, "HAVING", hav, p...)
	return st
}

func (st *SelectStm) OrderBy(expr string, p ...any) *SelectStm {
	st.add(limitS, "ORDER BY", expr, p...)
	return st
}

func (st *SelectStm) Limit(limit int) *SelectStm {
	st.add(limitS, "LIMIT", "?", limit)
	return st
}

func (st *SelectStm) Offset(off int) *SelectStm {
	st.add(offsetS, "OFFSET", "?", off)
	return st
}

func (st *SelectStm) Clause(clause, expr string, p ...any) *SelectStm {
	st.add(st.lastPos, clause, expr, p...)
	return st
}

func (st *SelectStm) ClauseIf(cond bool, clause, expr string, p ...any) *SelectStm {
	if cond {
		st.add(st.lastPos, clause, expr, p...)
	}
	return st
}
