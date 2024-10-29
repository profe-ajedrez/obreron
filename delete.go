package obreron

type DeleteStm struct {
	*stament
}

func Delete() *DeleteStm {
	d := &DeleteStm{
		stament: pool.New().(*stament),
	}

	d.add(deleteS, "DELETE", "")

	return d
}

func (dst *DeleteStm) From(source string) *DeleteStm {
	dst.add(fromS, "FROM", source)
	return dst
}

func (dst *DeleteStm) Where(cond string, p ...any) *DeleteStm {
	dst.where(cond, p...)
	return dst
}

func (dst *DeleteStm) Y() *DeleteStm {
	dst.clause("AND", "")
	return dst
}

func (dst *DeleteStm) And(expr string, p ...any) *DeleteStm {
	dst.clause("AND", expr, p...)
	return dst
}

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

func (dst *DeleteStm) Like(expr string, p ...any) *DeleteStm {
	dst.clause("LIKE", expr, p...)
	return dst
}

func (dst *DeleteStm) LikeIf(cond bool, expr string, p ...any) *DeleteStm {
	if cond {
		dst.Like(expr, p...)
	}
	return dst
}

func (dst *DeleteStm) In(value, expr string, p ...any) *DeleteStm {
	dst.clause(value+" IN ("+expr+")", "", p...)
	return dst
}

// InArgs adds an In clause to the stament automatically setting the positional parameters of the query based on the
// passed parameters
func (up *DeleteStm) InArgs(value string, p ...any) *DeleteStm {
	up.stament.inArgs(value, p...)
	return up
}

func (dst *DeleteStm) Close() {
	closeStament(dst.stament)
}

func (dst *DeleteStm) OrderBy(expr string, p ...any) *DeleteStm {
	dst.add(limitS, "ORDER BY", expr, p...)
	return dst
}

func (dst *DeleteStm) Limit(limit int) *DeleteStm {
	dst.add(limitS, "LIMIT", "?", limit)
	return dst
}

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
