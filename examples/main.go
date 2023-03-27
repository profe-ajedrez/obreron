package main

import "github.com/profe-ajedrez/obreron"

func main() {
	bl := obreron.NewMaryBuilder()

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
		obreron.NewMaryBuilder().Select("*").From("documento", "vdt").Where().AndParam("vdt.id_documento", "=", 126), "vdt", "vdt.id_documento = dvdt.id_documento",
	).Where().AndParam("vp.id_variante IS NOT NULL", "", nil)

	bl.AndParam("d.id_sucursal", "=", 126)
	bl.AndParamIf(filterByClassification, "p.clasificacion", "!=", 3).GroupBy("dd.id_detalle_desp")

	bl2 := obreron.NewMaryBuilder()

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

	_, _ = bl2.Build()
}
