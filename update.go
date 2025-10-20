package obreron

// UpdateStm represents an update stament
type UpdateStm struct {
	*stament
}

// Update returns an update stament
//
// # Example
//
// upd := Update("client").Set("status = 0").Where("status = ?", 1)
//
// query, p := upd.Build() // builds UPDATE client SET status = 0 WHERE status = ?  and stores in p []any{1}
//
// r, err := db.Exec(query, p...)
func Update(table string) *UpdateStm {
	//d := upPool.Get().(*UpdateStm)
	//d.stament = pool.Get().(*stament)
	d := &UpdateStm{
		stament: pool.Get().(*stament),
	}
	d.firstCol = true
	d.add(updateS, "UPDATE", table)
	return d
}

func CloseUpdate(s *UpdateStm) {
	CloseStament(s.stament)
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

// Set adds set clause to the update stament when the cond param is true
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

// And adds an AND conector with eventual parameters to the stament where is called, only when
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
func (up *UpdateStm) InArgs(value string, p ...any) *UpdateStm {
	up.stament.inArgs(value, p...)
	return up
}

// Close frees up the resources used in the stament
func (up *UpdateStm) Close() {
	CloseStament(up.stament)
}

func (up *UpdateStm) OrderBy(expr string, p ...any) *UpdateStm {
	up.add(limitS, "ORDER BY", expr, p...)
	return up
}

func (up *UpdateStm) Limit(limit int) *UpdateStm {
	up.add(limitS, "LIMIT", "?", limit)
	return up
}

func (up *UpdateStm) Clause(clause, expr string, p ...any) *UpdateStm {
	up.add(up.lastPos, clause, expr, p...)
	return up
}

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

func (up *UpdateStm) JoinIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.Join(expr, p...)

	}
	return up
}

func (up *UpdateStm) On(on string, p ...any) *UpdateStm {
	up.clause("ON", on, p...)
	return up
}

func (up *UpdateStm) OnIf(cond bool, expr string, p ...any) *UpdateStm {
	if cond {
		up.On(expr, p...)
	}
	return up
}
