package obreron

// UpdateStm represents an update stament
type UpdateStm struct {
	snapQ string
	*stament
	snapP  []any
	closed bool
}

// Update returns an update stament
func Update(table string) *UpdateStm {
	st, ok := pool.Get().(*stament)
	if !ok {
		st = &stament{}
	}

	d := &UpdateStm{
		stament: st,
	}
	d.firstCol = true
	d.add(updateS, "UPDATE", table)

	return d
}

// Build returns sql query and its parameters
func (up *UpdateStm) Build() (string, []any) {
	// Si ya cerramos, devolvemos snapshot estable.
	if up.closed {
		return up.snapQ, append([]any(nil), up.snapP...)
	}

	// Build normal (usa el stament real)
	q, p := up.stament.Build()

	// Cachea el último build para que Close() no tenga que reconstruir.
	up.snapQ = q

	up.snapP = append([]any(nil), p...) // copia defensiva

	return q, p
}

// CloseUpdate resets and returns to the pool an update stament
func CloseUpdate(up *UpdateStm) {
	if up.closed {
		return
	}

	// Si todavía no hay snapshot, congélalo ahora.
	// (Esto cubre el caso malo: Close() ocurre antes de Build()).
	if up.snapQ == "" && len(up.snapP) == 0 {
		q, p := up.stament.Build()
		up.snapQ = q

		up.snapP = append([]any(nil), p...)
	}

	up.closed = true

	// Devuelve recursos al pool.
	CloseStament(up.stament)

	// Detacha el stament para evitar “mezclas” si alguien llama métodos luego.
	up.stament = nil
}

// ColSelect is a helper method which provides a way to build an update (select ...) stament
//
// # Example
//
//	upd := Update("items").
//		ColSelectIf(
//			true,
//			Select().
//			Col("id, retail / wholesale AS markup, quantity").
//			From("items"), "discounted"
//		).Set("items.retail = items.retail * 0.9").
//		Set("a = 2").
//		SetIf(true, "c = 3").
//		Where("discounted.markup >= 1.3").
//		And("discounted.quantity < 100").
//		And("items.id = discounted.id").
//
//	query, p := upd.Build() // builds UPDATE items, ( SELECT id, retail / wholesale AS markup, quantity FROM items) discounted SET a = 2, c = 3 WHERE 1 = 1 AND discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id
func (up *UpdateStm) ColSelect(col *SelectStm, alias string) *UpdateStm {
	up.Clause(",(", "")

	q, p := col.Build()
	up.Clause(q, "", p...)
	up.Clause(")", "")
	up.Clause(alias, "")

	return up
}

// ColSelectIf does the same work as [ColSelect] only when the cond parameter is true
func (up *UpdateStm) ColSelectIf(cond bool, col *SelectStm, alias string) *UpdateStm {
	if cond {
		up.ColSelect(col, alias)
	}

	return up
}

// Set adds set clause to the update stament
//
// # Examples
//
//	upd := Update("client").Set("status = 0").Where("status = ?", 1)
//	up2 := Update("client").Set("status = ?", 0).Where("status = ?", 1)
//	up3 := Update("client").Set("status = ?", 0).Set("name = ?", "stitch").Where("status = ?", 1)
func (up *UpdateStm) Set(expr string, p ...any) *UpdateStm {
	if !up.firstCol {
		up.Clause(", ", "")
		up.add(setS, "", expr, p...)
	} else {
		up.firstCol = false
		up.add(setS, "SET", expr, p...)
	}

	return up
}

// SetIf adds set clause to the update stament when the cond param is true
func (up *UpdateStm) SetIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.Set(expr, p...)
	}

	return up
}

// Where adds a where clause to the update stament
func (up *UpdateStm) Where(cond string, p ...any) *UpdateStm {
	up.where(cond, p...)
	return up
}

// Y adds an AND conector to the stament where is called. Its helpful when used with In()
//
// # Example
//
//	Update("client").Set("status = 0").Where("country = ?", "CL").Y().In("status", "", 1, 2, 3, 4)
//
// Produces: UPDATE client SET status = 0 WHERE country = ? AND status IN (?, ?, ?, ?)
func (up *UpdateStm) Y() *UpdateStm {
	up.clause("AND", "")
	return up
}

// And adds an AND conector with eventual parameters to the stament where is called
func (up *UpdateStm) And(expr string, p ...any) *UpdateStm {
	up.clause("AND", expr, p...)
	return up
}

// AndIf adds an AND conector with eventual parameters to the stament where is called, only when
// cond parameter is true
func (up *UpdateStm) AndIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.clause("AND", expr, p...)
	}

	return up
}

// Or adds an Or connector with eventual parameters to the stament where is called
func (up *UpdateStm) Or(expr string, p ...any) *UpdateStm {
	up.clause("OR", expr, p...)
	return up
}

// OrIf adds an Or connector with eventual parameters to the stament where is called only when cond parameter value is true
func (up *UpdateStm) OrIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.clause("OR", expr, p...)
	}

	return up
}

// Like adds a LIKE clause to the query after the last added clause
//
//	# Example
//
//	Update("items").
//	Set("items.retail = items.retail * 0.9").
//	Set("a = 2").
//	Where("discounted.markup >= 1.3").
//	And("colX").
//	Like("'%ago%'")
func (up *UpdateStm) Like(expr string, p ...any) *UpdateStm {
	up.clause("LIKE", expr, p...)
	return up
}

// LikeIf adds a LIKE clause to the query after the last added clause only when cond parameter value is true
//
//	# Example
//
//	Update("items").
//	Set("items.retail = items.retail * 0.9").
//	Set("a = 2").
//	Where("discounted.markup >= 1.3").
//	And("colX").
//	Like("'%ago%'")
func (up *UpdateStm) LikeIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.clause("LIKE", expr, p...)
	}

	return up
}

// In adds a IN clause to the query after the las clause added
//
// # Example
//
//	Update("client").
//	Set("status = 0").
//	Where("country = ?", "CL").
//	Y().In("status", "?, ?, ?, ?", 1, 2, 3, 4)
func (up *UpdateStm) In(value, expr string, p ...any) *UpdateStm {
	up.clause(value+" IN ("+expr+")", "", p...)
	return up
}

// InArgs adds an In clause to the stament automatically setting the positional parameters of the query based on the
// passed parameters
//
// # Example
//
//	Update("client").Set("status = 0").Where("country = ?", "CL").Y().InArgs("status", 1, 2, 3, 4)
//
// Produces: UPDATE client SET status = 0 WHERE country = ? AND status IN (?, ?, ?, ?)"
// If p is empty, the generated clause will never match any row.
func (up *UpdateStm) InArgs(value string, p ...any) *UpdateStm {
	up.inArgs(value, p...)
	return up
}

// Close frees up the resources used in the stament
func (up *UpdateStm) Close() {
	CloseUpdate(up)
}

// OrderBy adds an Order  clause
func (up *UpdateStm) OrderBy(expr string, p ...any) *UpdateStm {
	up.add(orderS, "ORDER BY", expr, p...)
	return up
}

// Limit adds a limit clause
func (up *UpdateStm) Limit(limit int) *UpdateStm {
	up.add(limitS, "LIMIT", "?", limit)
	return up
}

// Clause adds a RAW clause
func (up *UpdateStm) Clause(clause, expr string, p ...any) *UpdateStm {
	up.add(up.lastPos, clause, expr, p...)
	return up
}

// ClauseIf adds a RAW clause if cond is true
func (up *UpdateStm) ClauseIf(cond bool, clause, expr string, p ...any) *UpdateStm {
	if cond {
		up.Clause(clause, expr, p...)
	}

	return up
}

// Join adds a relation to the query in the form of an inner join
//
// # Example
//
//	Update("business AS b").
//	Join("business_geocode AS g").On("b.business_id = g.business_id").
//	Set("b.mapx = g.latitude, b.mapy = g.longitude").
//	Where("(b.mapx = '' or b.mapx = 0)").And("g.latitude > 0")
//
//	OUTPUT:
//	UPDATE business AS b JOIN business_geocode AS g ON b.business_id = g.business_id SET b.mapx = g.latitude, b.mapy = g.longitude WHERE (b.mapx = '' or b.mapx = 0) AND g.latitude > 0 AND 3 = 3
func (up *UpdateStm) Join(expr string, p ...any) *UpdateStm {
	up.add(updateS, "JOIN", expr, p...)
	return up
}

// JoinIf adds a join clause if cond is true
func (up *UpdateStm) JoinIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.Join(expr, p...)
	}

	return up
}

// On adds a On estament
func (up *UpdateStm) On(on string, p ...any) *UpdateStm {
	up.clause("ON", on, p...)
	return up
}

// OnIf adds a On stament if cond is true
func (up *UpdateStm) OnIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.On(expr, p...)
	}

	return up
}
