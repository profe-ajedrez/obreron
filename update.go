package obreron

type UpdateStament struct {
	*stament
}

func Update(table string) *UpdateStament {
	d := &UpdateStament{
		stament: pool.New().(*stament),
	}
	d.firstCol = true
	d.add(updateS, "UPDATE", table)
	return d
}

func (up *UpdateStament) ColSelect(col *SelectStm, alias string) *UpdateStament {

	up.Clause(",(", "")

	q, p := col.Build()
	up.Clause(q, "", p...)
	up.Clause(")", "")
	up.Clause(alias, "")

	return up
}

func (up *UpdateStament) ColSelectIf(cond bool, col *SelectStm, alias string) *UpdateStament {
	if cond {
		up.ColSelect(col, alias)
	}

	return up
}

func (up *UpdateStament) Set(expr string, p ...any) *UpdateStament {
	if !up.firstCol {
		up.Clause(", ", "")
		up.add(setS, "", expr, p...)
	} else {
		up.firstCol = false
		up.add(setS, "SET", expr, p...)
	}

	return up
}

func (up *UpdateStament) SetIf(cond bool, expr string, p ...any) *UpdateStament {
	if cond {
		up.Set(expr, p...)
	}

	return up
}

func (up *UpdateStament) Where(cond string, p ...any) *UpdateStament {
	up.where(cond, p...)
	return up
}

func (up *UpdateStament) Y() *UpdateStament {
	up.clause("AND", "")
	return up
}

func (up *UpdateStament) And(expr string, p ...any) *UpdateStament {
	up.clause("AND", expr, p...)
	return up
}

func (up *UpdateStament) AndIf(cond bool, expr string, p ...any) *UpdateStament {
	if cond {
		up.clause("AND", expr, p...)
	}
	return up
}

func (up *UpdateStament) Or(expr string, p ...any) *UpdateStament {
	up.clause("OR", expr, p...)
	return up
}

func (up *UpdateStament) OrIf(cond bool, expr string, p ...any) *UpdateStament {
	if cond {
		up.clause("OR", expr, p...)
	}
	return up
}

func (up *UpdateStament) Like(expr string, p ...any) *UpdateStament {
	up.clause("LIKE", expr, p...)
	return up
}

func (up *UpdateStament) LikeIf(cond bool, expr string, p ...any) *UpdateStament {
	if cond {
		up.clause("LIKE", expr, p...)
	}
	return up
}

func (up *UpdateStament) In(value, expr string, p ...any) *UpdateStament {
	up.clause(value+" IN ("+expr+")", "", p...)
	return up
}

func (up *UpdateStament) Close() {
	closeStament(up.stament)
}

func (up *UpdateStament) OrderBy(expr string, p ...any) *UpdateStament {
	up.add(limitS, "ORDER BY", expr, p...)
	return up
}

func (up *UpdateStament) Limit(limit int) *UpdateStament {
	up.add(limitS, "LIMIT", "?", limit)
	return up
}

func (up *UpdateStament) Clause(clause, expr string, p ...any) *UpdateStament {
	up.add(up.lastPos, clause, expr, p...)
	return up
}

func (up *UpdateStament) ClauseIf(cond bool, clause, expr string, p ...any) *UpdateStament {
	if cond {
		up.Clause(clause, expr, p...)
	}
	return up
}

func (up *UpdateStament) Join(expr string, p ...any) *UpdateStament {
	up.add(updateS, "JOIN", expr, p...)
	return up
}

func (up *UpdateStament) JoinIf(cond bool, expr string, p ...any) *UpdateStament {
	if cond {
		up.Join(expr, p...)

	}
	return up
}

func (up *UpdateStament) On(on string, p ...any) *UpdateStament {
	up.clause("ON", on, p...)
	return up
}

func (up *UpdateStament) OnIf(cond bool, expr string, p ...any) *UpdateStament {
	if cond {
		up.On(expr, p...)
	}
	return up
}
