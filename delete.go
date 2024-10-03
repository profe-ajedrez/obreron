package obreron

type DeleteStament struct {
	*stament
}

func Delete() *DeleteStament {
	d := &DeleteStament{
		stament: pool.New().(*stament),
	}

	d.add(deleteS, "DELETE", "")

	return d
}

func (dst *DeleteStament) From(source string) *DeleteStament {
	dst.add(fromS, "FROM", source)
	return dst
}

func (dst *DeleteStament) Where(cond string, p ...any) *DeleteStament {
	dst.where(cond, p...)
	return dst
}

func (dst *DeleteStament) Y() *DeleteStament {
	dst.clause("AND", "")
	return dst
}

func (dst *DeleteStament) And(expr string, p ...any) *DeleteStament {
	dst.clause("AND", expr, p...)
	return dst
}

func (dst *DeleteStament) AndIf(cond bool, expr string, p ...any) *DeleteStament {
	if cond {
		dst.clause("AND", expr, p...)
	}
	return dst
}

func (dst *DeleteStament) Or(expr string, p ...any) *DeleteStament {
	dst.clause("OR", expr, p...)
	return dst
}

func (dst *DeleteStament) OrIf(cond bool, expr string, p ...any) *DeleteStament {
	if cond {
		dst.clause("OR", expr, p...)
	}
	return dst
}

func (dst *DeleteStament) Like(expr string, p ...any) *DeleteStament {
	dst.clause("LIKE", expr, p...)
	return dst
}

func (dst *DeleteStament) LikeIf(cond bool, expr string, p ...any) *DeleteStament {
	if cond {
		dst.clause("LIKE", expr, p...)
	}
	return dst
}

func (dst *DeleteStament) In(value, expr string, p ...any) *DeleteStament {
	dst.clause(value+" IN ("+expr+")", "", p...)
	return dst
}

func (dst *DeleteStament) Close() {
	Close(dst.stament)
}

func (dst *DeleteStament) OrderBy(expr string, p ...any) *DeleteStament {
	dst.add(limitS, "ORDER BY", expr, p...)
	return dst
}

func (dst *DeleteStament) Limit(limit string, p ...any) *DeleteStament {
	dst.add(limitS, "LIMIT", limit, p...)
	return dst
}

func (dst *DeleteStament) Clause(clause, expr string, p ...any) *DeleteStament {
	dst.add(dst.lastPos, clause, expr, p...)
	return dst
}

func (dst *DeleteStament) ClauseIf(cond bool, clause, expr string, p ...any) *DeleteStament {
	if cond {
		dst.add(dst.lastPos, clause, expr, p...)
	}
	return dst
}
