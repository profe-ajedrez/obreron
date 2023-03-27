# Obreron

Sql Builder escrito en Go.

---

## Dialectos soportados

- [x] Mysql  
- [ ] Postgresql  
  

---

## ¿Por que?

Porqué escribir y usar SQL en go es aburrido y trabajoso. Sobre todo en el caso de tener que construirlas dinámicamente.

Hay varios proyectos similares, pero no encontramos ninguno que cumpliera con el requisito de ser liviano, o agregaban también una referencia a la conexión u otras dependencias.


## Instalación

Puede instalar obreron usando `go get`

```bash
$ go get github.com/profe-ajedrez/obreron
```



## Uso


Con obreron es fácil construir consultas. Por ejemplo el siguiente código 


```go
b = obreron.NewMaryBuilder()

q := b.Select(
    "id",
    "name",
    "mail",
    b.Quote("columna con espacios en el nombre"),
).From("users", "u").String()
```

produce la siguiente consulta en la variable `q`: 

```
SELECT id,name,mail, `columna con espacios en el nombre` FROM users u 
```

Pero para eso no necesitamos un Sql Builder. Estos brillan cuando debemos construir consultas dinámicas, como para el caso de armar filtros.

```go
b := NewMaryBuilder()

	// options es un struct que guarda las opciones para armar el filtro
	options := struct {
		useName     bool
		useFullName bool
		useAddress  bool
		status      int8
		limit       int64
	}{
		useName:     true, // agregar columna user_name
		useFullName: true, // agregar columna useFullName
		useAddress:  false, // NO agregar columna address
		status:      0, // para saber como filtrar users por su status
		limit:       25, // limite para la consulta
	}

    // Los campos user_id, user_mail y user_type serán agregados a la consulta
	b.Select(
		"user_id", "user_mail", "user_type",
	).From("users", "u").Where().Limit(options.limit)

    // Los campos user_name, user_fullname y user_address se agregaran a la consulta solo si la condición pasada es verdadera
	b.AddColumnIf(options.useName, "user_name", "").AddColumnIf(options.useFullName, "user_fullname", "")
	b.AddColumnIf(options.useAddress, "user_address", "")

    // Agregamos una condición en la que decimos que agregue el filtro por user_status si a opción correspondiente es > -1
	b.AndParamIf(options.status > -1, "user_status", "=", options.status)

    // tras construir la consulta, q contendrá la query y p los parámetros registrados para su uso
	q, p := b.Build()

```

El ejemplo anterior construye en `q` la consulta `SELECT user_id,user_mail,user_type,user_name,user_fullname FROM users u  WHERE 1=1  AND user_status = ? LIMIT 25 ` y deja en `p`  un `interface {}(int8) 0`


## Benchmarks

Presentamos el siguiente benchmark

```go

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

```

```bash
go test -benchmem -run=^$ -bench ^BenchmarkSlBuilder$ -count=5 -benchtime=5s  
goos: linux
goarch: amd64
pkg: github.com/profe-ajedrez/obreron
cpu: Intel(R) Core(TM) i7-10700 CPU @ 2.90GHz
BenchmarkSlBuilder-16             744958             10261 ns/op           13816 B/op        122 allocs/op
BenchmarkSlBuilder-16             712387              9246 ns/op           13817 B/op        122 allocs/op
BenchmarkSlBuilder-16             641104              8842 ns/op           13817 B/op        122 allocs/op
BenchmarkSlBuilder-16             695001              9042 ns/op           13817 B/op        122 allocs/op
BenchmarkSlBuilder-16             742594              9070 ns/op           13816 B/op        122 allocs/op
PASS
ok      github.com/profe-ajedrez/obreron        33.357s
```


## Static checks


### goreportcard-cli

```bash
~/go/bin/goreportcard-cli   
Grade .......... A+ 100.0%
Files .................. 4
Issues ................. 3
gofmt ............... 100%
go_vet .............. 100%
gocyclo ............. 100%
ineffassign ......... 100%
license ............. 100%
misspell ............. 25%
```

