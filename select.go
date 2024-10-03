package obreron

type SelectStament struct {
	*stament
}

func Select() *SelectStament {
	s := &SelectStament{
		pool.Get().(*stament),
	}

	s.add(selectS, "SELECT", "")
	return s
}

func (st *SelectStament) Close() {
	Close(st.stament)
}

func (st *SelectStament) Col(expr string, p ...any) *SelectStament {
	if !st.firstCol {
		st.add(colsS, ",", expr, p...)
		return st
	}

	st.add(colsS, "", expr, p...)
	st.firstCol = false
	return st
}

func (st *SelectStament) ColIf(cond bool, expr string, p ...any) *SelectStament {
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

func (st *SelectStament) From(source string) *SelectStament {
	st.add(fromS, "FROM", source)
	return st
}

func (st *SelectStament) Join(expr string, p ...any) *SelectStament {
	st.add(joinS, "JOIN", expr, p...)
	return st
}

func (st *SelectStament) JoinIf(cond bool, expr string, p ...any) *SelectStament {
	if cond {
		st.add(joinS, "JOIN", expr, p...)
	}
	return st
}

func (st *SelectStament) LeftJoin(expr string, p ...any) *SelectStament {
	st.add(joinS, "LEFT JOIN", expr, p...)
	return st
}

func (st *SelectStament) LeftJoinIf(cond bool, join string, p ...any) *SelectStament {
	if cond {
		st.add(joinS, "LEFT JOIN", join, p...)
	}
	return st
}

func (st *SelectStament) RightJoin(expr string, p ...any) *SelectStament {
	st.add(joinS, "RIGHT JOIN", expr, p...)
	return st
}

func (st *SelectStament) RightJoinIf(cond bool, expr string, p ...any) *SelectStament {
	if cond {
		st.add(joinS, "RIGHT JOIN", expr, p...)
	}
	return st
}

func (st *SelectStament) OuterJoin(expr string, p ...any) *SelectStament {
	st.add(joinS, "OUTER JOIN", expr, p...)
	return st
}

func (st *SelectStament) OuterJoinIf(cond bool, expr string, p ...any) *SelectStament {
	if cond {
		st.add(joinS, "OUTER JOIN", expr, p...)
	}
	return st
}

func (st *SelectStament) On(on string, p ...any) *SelectStament {
	st.clause("ON", on, p...)
	return st
}

func (st *SelectStament) OnIf(cond bool, expr string, p ...any) *SelectStament {
	if cond {
		st.clause("ON", expr, p...)
	}
	return st
}

func (st *SelectStament) Where(cond string, p ...any) *SelectStament {
	st.where(cond, p...)

	return st
}

func (st *SelectStament) And(expr string, p ...any) *SelectStament {
	st.clause("AND", expr, p...)
	return st
}

func (st *SelectStament) AndIf(cond bool, expr string, p ...any) *SelectStament {
	if cond {
		st.clause("AND", expr, p...)
	}
	return st
}

func (st *SelectStament) Or(expr string, p ...any) *SelectStament {
	st.clause("OR", expr, p...)
	return st
}

func (st *SelectStament) OrIf(cond bool, expr string, p ...any) *SelectStament {
	if cond {
		st.clause("OR", expr, p...)
	}
	return st
}

func (st *SelectStament) Like(expr string, p ...any) *SelectStament {
	st.clause("LIKE", expr, p...)
	return st
}

func (st *SelectStament) LikeIf(cond bool, expr string, p ...any) *SelectStament {
	if cond {
		st.clause("LIKE", expr, p...)
	}
	return st
}

func (st *SelectStament) In(expr string, p ...any) *SelectStament {
	st.clause("IN (", expr+")", p...)
	return st
}

func (st *SelectStament) GroupBy(grp string, p ...any) *SelectStament {
	if !st.grouped {
		st.add(groupS, "GROUP BY", grp, p...)
		st.grouped = true
	} else {
		st.add(groupS, ",", grp, p...)
	}

	return st
}

func (st *SelectStament) Having(hav string, p ...any) *SelectStament {
	st.add(havingS, "HAVING", hav, p...)
	return st
}

func (st *SelectStament) OrderBy(expr string, p ...any) *SelectStament {
	st.add(limitS, "ORDER BY", expr, p...)
	return st
}

func (st *SelectStament) Limit(limit string, p ...any) *SelectStament {
	st.add(limitS, "LIMIT", limit, p...)
	return st
}

func (st *SelectStament) Offset(off string, p ...any) *SelectStament {
	st.add(offsetS, "OFFSET", off, p...)
	return st
}

func (st *SelectStament) Clause(clause, expr string, p ...any) *SelectStament {
	st.add(st.lastPos, clause, expr, p...)
	return st
}

func (st *SelectStament) ClauseIf(cond bool, clause, expr string, p ...any) *SelectStament {
	if cond {
		st.add(st.lastPos, clause, expr, p...)
	}
	return st
}
