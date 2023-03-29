package obreron

import "testing"

func TestSimpleSQlBuild(t *testing.T) {
	b := NewMaryBuilder()

	q := b.Select(
		"id",
		b.Quote("columna_de_nombre_largo"),
	).From("table_a", "t").String()

	t.Log(q)

	expected := "SELECT id,`columna_de_nombre_largo` FROM table_a t "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}
}

func TestOtheSimpleSQlBuild(t *testing.T) {
	b := NewMaryBuilder()

	q := b.Select(
		"id",
		"name",
		"mail",
		b.Quote("columna con espacios en el nombre"),
	).From("users", "u").String()

	t.Log(q)

	expected := "SELECT id,name,mail,`columna con espacios en el nombre` FROM users u "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}
}

func TestSubQuerySQlBuild(t *testing.T) {
	b := NewMaryBuilder()

	b2 := NewMaryBuilder()

	b2.Select("SUM(monto_neto) AS monto_neto").From("table_c", "cc")

	q := b.Select(
		"id",
		b2,
	).From("table_a", "t").String()

	t.Log(q)

	expected := "SELECT id,(SELECT SUM(monto_neto) AS monto_neto FROM table_c AS cc ) FROM table_a t "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}
}

func TestSimpleSQlWithLimitBuild(t *testing.T) {
	b := NewMaryBuilder()

	q := b.Select(
		"id",
		b.Quote("columna_de_nombre_largo"),
	).From(
		"table_a",
		"t",
	).Limit(10).String()

	t.Log(q)

	expected := "SELECT id,`columna_de_nombre_largo` FROM table_a AS t  LIMIT 10 "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}
}

func TestAddingColumnsAfter(t *testing.T) {
	b := NewMaryBuilder()

	b.Select(
		"id",
		b.Quote("columna_de_nombre_largo"),
	).From(
		"table_a",
		"t",
	).Limit(10)

	b.AddColumn("`esta_columna_se_agregara_despues_de_formado_el_builder`", "otro_alias")
	b.AddColumnIf(false, "`esta_columna_no_se_agregara_al_builder`", "")

	q := b.String()
	t.Log(q)

	expected := "SELECT id,`columna_de_nombre_largo`,`esta_columna_se_agregara_despues_de_formado_el_builder` AS otro_alias FROM table_a t  LIMIT 10 "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}
}

func TestOrderBy(t *testing.T) {
	b := NewMaryBuilder()

	b.Select(
		"id",
		b.Quote("columna_de_nombre_largo"),
	).From(
		"table_a",
		"t",
	).Limit(10)

	b.AddColumn("`esta_columna_se_agregara_despues_de_formado_el_builder`", "otro_alias")
	b.AddColumnIf(false, "`esta_columna_no_se_agregara_al_builder`", "")

	q := b.OrderBy("1 ASC").String()
	t.Log(q)

	expected := "SELECT id,`columna_de_nombre_largo`,`esta_columna_se_agregara_despues_de_formado_el_builder` AS otro_alias FROM table_a  t  ORDER BY 1 ASC  LIMIT 10 "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}
}

func TestWhered(t *testing.T) {

	q, p := testWhere(t)
	t.Log(q)
	t.Log(p)

	expected := "SELECT user_id,user_mail,user_type,user_name,user_fullname FROM users u  WHERE 1=1  AND user_status = ? LIMIT 25 "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
		t.FailNow()
	}

	if len(p) != 1 {
		t.Log("expected  : 1 params\n")
		t.Logf("generated: %d", len(p))
		t.FailNow()
	}

	p0, ok := p[0].(int8)

	if !ok {
		t.Log("expected  : int param\n")
		t.FailNow()
	}

	if ok && p0 != 0 {
		t.Log("expected  : 0 as first param\n")
		t.Logf("generated: %d", p[0])
		t.FailNow()
	}

}

func TestJoin(t *testing.T) {
	q := testJoin(t)

	t.Log(q)

	expected := "SELECT id,(SELECT SUM(monto_neto) AS monto_neto FROM table_c AS cc ) FROM table_a t  WHERE 1=1  AND 'B' = 'B'"

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}
}

func TestHeavySlBuilder(t *testing.T) {
	q, p := heavyQueryBuild(t)
	t.Log(q)
	t.Log(p)

	expected := "SELECT created_at_ms,fecha_movimiento,nombre_producto,cantidad_entrada,cantidad_salida,stock,costo,codigo_variante_producto,codigo_barras,num_doc_tributario,nombre_tipo_documento,id_venta_documento_tributario,id_detalle_ingreso_stock,id_consumo_stock,usuario_movimiento,id_despacho,uso_documento,numero_serie  FROM (SELECT dd.created_at_ms,d.fecha_despacho AS fecha_movimiento,CONCAT(\n\t\tp.nombre_producto, ' ', vp.descripcion_variante\n\t\t) AS nombre_producto,0 AS cantidad_entrada,dd.cantidad_desp AS cantidad_salida,dd.stock AS stock,dd.costo AS costo,vp.codigo,vp.barras,CASE WHEN vdt.estado = 0 THEN IFNULL(vdt.num_doc, '') ELSE '' END AS num_doc,CASE WHEN vdt.estado = 0 THEN IFNULL(td.nombre_tipo, '') ELSE 'nulo' END AS nombre_tipo,CASE WHEN vdt.estado = 0 THEN IFNULL(\n\t\tvdt.id_documento, \n\t\t0\n\t\t) ELSE -99 END AS id_documento,IFNULL(dd.ids_detalle_ingreso, 0) AS id_detalle_ingreso,0 AS id_consumo,CONCAT(\n\t\tus.nombre_usuario, ' ', us.apellido_usuario\n\t\t) AS usuario_movimiento,d.id_despacho AS id_despacho,td.uso_documento AS uso_documento,dd.id_detalle_despacho,vdt.estado_documento,dd.numero_serie FROM detalle_desp dd  INNER JOIN cart_it ci  ON dd.id_cart_it=ci.id_cart_it INNER JOIN variante vp  ON ci.id_variante= vp.id_variante INNER JOIN producto p  ON vp.id_producto = p.id_producto INNER JOIN detalle_venta_documento_tributario dvdt  ON dd.id_detalle_despacho = dvdt.id_detalle_despacho INNER JOIN (SELECT * FROM documento vdt  WHERE 1=1  AND vdt.id_documento = ?) vdt  ON vdt.id_documento = dvdt.id_documento WHERE 1=1  AND vp.id_variante IS NOT NULL AND d.id_sucursal = ? AND p.clasificacion != ? GROUP BY dd.id_detalle_desp ) out_detail  GROUP BY id_detalle_desp  ORDER BY id_despacho ASC "

	if q != expected {
		t.Logf("expected : %s", expected)
		t.Logf("generated: %s", q)
	}

	if len(p) != 3 {
		t.Log("expected : 3 params\n")
		t.Logf("generated: %d\n", len(p))
	}

	if p[0] != 126 {
		t.Log("expected p[0] == 126 params\n")
		t.Logf("generated: %d\n", p[0])
	}

	if p[1] != 126 {
		t.Log("expected p[1] == 126 params\n")
		t.Logf("generated: %d\n", p[1])
	}

	if p[2] != 3 {
		t.Log("expected p[2] == 3 params\n")
		t.Logf("generated: %d\n", p[2])
	}
}

func BenchmarkSlBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		heavyQueryBuild(b)
	}

}

func heavyQueryBuild(b testing.TB) (string, []interface{}) {
	bl := NewMaryBuilder()

	filterByClassification := true

	bl.Select(
		"dd.created_at_ms",
		"d.fecha_despacho AS fecha_movimiento",
		`CONCAT(
		p.nombre_producto, ' ', vp.descripcion_variante
		) AS nombre_producto`,
		"0 AS cantidad_entrada",
		"dd.cantidad_desp AS cantidad_salida",
		"dd.stock AS stock",
		"dd.costo AS costo",
		"vp.codigo",
		"vp.barras",
		"CASE WHEN vdt.estado = 0 THEN IFNULL(vdt.num_doc, '') ELSE '' END AS num_doc",
		"CASE WHEN vdt.estado = 0 THEN IFNULL(td.nombre_tipo, '') ELSE 'nulo' END AS nombre_tipo",
		`CASE WHEN vdt.estado = 0 THEN IFNULL(
		vdt.id_documento, 
		0
		) ELSE -99 END AS id_documento`,
		"IFNULL(dd.ids_detalle_ingreso, 0) AS id_detalle_ingreso",
		"0 AS id_consumo",
		`CONCAT(
		us.nombre_usuario, ' ', us.apellido_usuario
		) AS usuario_movimiento`,
		"d.id_despacho AS id_despacho",
		"td.uso_documento AS uso_documento",
		"dd.id_detalle_despacho",
		"vdt.estado_documento",
		"dd.numero_serie",
	).From(
		"detalle_desp", "dd",
	).Inner(
		"cart_it", "ci", "dd.id_cart_it=ci.id_cart_it",
	).Inner(
		"variante", "vp", "ci.id_variante= vp.id_variante",
	).Inner(
		"producto", "p", "vp.id_producto = p.id_producto",
	).Inner(
		"detalle_venta_documento_tributario", "dvdt", "dd.id_detalle_despacho = dvdt.id_detalle_despacho",
	).Inner(
		NewMaryBuilder().Select("*").From("documento", "vdt").Where().AndParam("vdt.id_documento", "=", 126), "vdt", "vdt.id_documento = dvdt.id_documento",
	).Where().AndParam("vp.id_variante IS NOT NULL", "", nil)

	bl.AndParam("d.id_sucursal", "=", 126)
	bl.AndParamIf(filterByClassification, "p.clasificacion", "!=", 3).GroupBy("dd.id_detalle_desp")

	bl2 := NewMaryBuilder()

	bl2.Select().From(bl, "out_detail")

	bl2.AddColumn("created_at_ms", "")
	bl2.AddColumn("fecha_movimiento", "")
	bl2.AddColumn("nombre_producto", "")
	bl2.AddColumn("cantidad_entrada", "")
	bl2.AddColumn("cantidad_salida", "")
	bl2.AddColumn("stock", "")
	bl2.AddColumn("costo", "")
	bl2.AddColumn("codigo_variante_producto", "")
	bl2.AddColumn("codigo_barras", "")
	bl2.AddColumn("num_doc_tributario", "")
	bl2.AddColumn("nombre_tipo_documento", "")
	bl2.AddColumn("id_venta_documento_tributario", "")
	bl2.AddColumn("id_detalle_ingreso_stock", "")
	bl2.AddColumn("id_consumo_stock", "")
	bl2.AddColumn("usuario_movimiento", "")
	bl2.AddColumn("id_despacho", "")
	bl2.AddColumn("uso_documento", "")
	bl2.AddColumn("numero_serie ", "")

	bl2.GroupBy("id_detalle_desp")
	bl2.OrderBy("id_despacho ASC")

	q, params := bl2.Build()

	return q, params
}

func BenchmarkJoin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		testJoin(b)
	}
}

func testJoin(t testing.TB) string {
	b := NewMaryBuilder()

	b2 := NewMaryBuilder()

	b2.Select("SUM(monto_neto) AS monto_neto").From("table_c", "cc")

	return b.Select(
		"id",
		b2,
	).From("table_a", "t").Where().AndParam("'B' = 'B'", "", nil).String()
}

func testWhere(t testing.TB) (string, []interface{}) {
	b := NewMaryBuilder()

	// options es un struct que guarda las opciones para armar el filtro
	options := struct {
		useName     bool
		useFullName bool
		useAddress  bool
		status      int8
		limit       int64
	}{
		useName:     true,
		useFullName: true,
		useAddress:  false,
		status:      0,
		limit:       25,
	}

	b.Select(
		"user_id", "user_mail", "user_type",
	).From("users", "u").Where().Limit(options.limit)

	b.AddColumnIf(options.useName, "user_name", "").AddColumnIf(options.useFullName, "user_fullname", "")
	b.AddColumnIf(options.useAddress, "user_address", "")

	b.AndParamIf(options.status > -1, "user_status", "=", options.status)

	return b.Build()
}
