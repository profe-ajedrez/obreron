# Obreron

Fast and cheap sql builder.


## Supported Dialects

- [x] Mysql
- [ ] Postgresql (Work in progress)


## Why?

Because writing SQL in Go is boring.

## How fast is it?

```bash
goos: linux
goarch: amd64
pkg: github.com/profe-ajedrez/obreron/v2
cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
BenchmarkSelect/columns_-_from-8                44610848                62.10 ns/op           24 B/op          0 allocs/op
BenchmarkSelect/columns_-_from#01-8             24653995                63.52 ns/op           25 B/op          0 allocs/op
BenchmarkSelect/columns_params_-_from-8         32579121                59.88 ns/op           16 B/op          0 allocs/op
BenchmarkSelect/columns_params_-_columns_params_-_from-8                28909456                71.47 ns/op           23 B/op          0 allocs/op
BenchmarkSelect/columns_params_-_columns_params_-_Col_If_-_from-8               35311484                64.07 ns/op           19 B/op          0 allocs/op
BenchmarkSelect/columns_params_-_columns_params_-_Col_If_-_from_-_where-8       36645351                69.15 ns/op           21 B/op          0 allocs/op
BenchmarkSelect/columns_params_-_columns_params_-_Col_If_-_from_-_where#01-8            28030076                71.77 ns/op           28 B/op          0 allocs/op
BenchmarkSelect/columns_params_-_columns_params_-_Col_If_-_from_-_where#02-8            33171466                62.83 ns/op           16 B/op          0 allocs/op
BenchmarkSelect/columns_params_-_columns_params_-_Col_If_-_from_-_where_shuffled-8      39448615                69.81 ns/op           20 B/op          0 allocs/op
BenchmarkSelect/complex_query_shuffled-8                                                25634710                75.16 ns/op           23 B/op          0 allocs/op
BenchmarkSelect/complex_query_badly_shuffled-8                                          45230116                68.14 ns/op           19 B/op          0 allocs/op
BenchmarkDelete/simple_del-8                                                            24681375                68.65 ns/op           21 B/op          0 allocs/op
BenchmarkDelete/simple_del_where-8                                                      23892154                68.72 ns/op           22 B/op          0 allocs/op
BenchmarkDelete/del_where_conditions-8                                                  23020927                73.73 ns/op           23 B/op          0 allocs/op
BenchmarkDelete/del_where_conditions_limit-8                                            23444376                73.24 ns/op           22 B/op          0 allocs/op
BenchmarkDelete/del_where_conditions_limit_--_shuffled-8                                20474908                75.58 ns/op           19 B/op          0 allocs/op
BenchmarkDelete/simple_del_where_quick-8                                                18463716                82.85 ns/op           22 B/op          0 allocs/op
BenchmarkDelete/simple_del_where_ignore-8                                               21899524                71.26 ns/op           24 B/op          0 allocs/op
BenchmarkDelete/simple_del_where_partition-8                                            22774772                76.24 ns/op           23 B/op          0 allocs/op
BenchmarkDelete/simple_del_where_order_by_limit-8                                       24599192                72.83 ns/op           24 B/op          0 allocs/op
BenchmarkUpdate/update_simple-8                                                         26563249                75.88 ns/op           20 B/op          0 allocs/op
BenchmarkUpdate/update_where-8                                                          23154291                73.58 ns/op           17 B/op          0 allocs/op
BenchmarkUpdate/update_where_order_limit-8                                              25949305                68.23 ns/op           20 B/op          0 allocs/op
BenchmarkUpdate/update_where_and_order_limit-8                                          44520046                79.09 ns/op           27 B/op          0 allocs/op
BenchmarkUpdate/update_select-8                                                         23995230                79.18 ns/op           23 B/op          0 allocs/op
BenchmarkUpdate/update_join-8                                                           22023397                80.67 ns/op           21 B/op          0 allocs/op
BenchmarkInsert/simple_insert-8                                                         14466757               145.0 ns/op            85 B/op          1 allocs/op
BenchmarkInsert/simple_insert_params-8                                                  14740395               134.9 ns/op            84 B/op          1 allocs/op
BenchmarkInsert/simple_insert_params_shuffled-8                                         12962823               121.0 ns/op            85 B/op          1 allocs/op
BenchmarkInsert/simple_insert_params_select-8                                           16434265               135.9 ns/op            84 B/op          1 allocs/op
```

## Instalation

Use `go get`

```bash
$ go get github.com/profe-ajedrez/obreron/v2
```

## Use

### Select

* Simple select

```go
// Produces SELECT a1, a2, a3 FROM client
query, _ := Select().Col("a1, a2, a3").From("client").Build()
r, error := db.Query(query)
```

* Select/join/where/shuffled

```go
// Produces SELECT a1, a2, ? AS diez, colIf1, colIf2, ? AS zero, a3, ? AS cien FROM client c JOIN addresses a ON a.id_cliente = a.id_cliente JOIN phones p ON p.id_cliente = c.id_cliente JOIN mailes m ON m.id_cliente = m.id_cliente AND c.estado_cliente = ? LEFT JOIN left_joined lj ON lj.a1 = c.a1 WHERE a1 = ? AND a2 = ? AND a3 = 10 AND a16 = ?
// with params = []any{10, 0, 100, 0, "'last name'", 1000.54, 75}
query, params := Select().
    Where("a1 = ?", "'last name'").
    Col("a1, a2, ? AS diez", 10).
    Col(`colIf1, colIf2, ? AS zero`, 0).
    Col("a3, ? AS cien", 100).    
    Where("a2 = ?", 1000.54).
    And("a3 = 10").And("a16 = ?", 75).
    Join("addresses a ON a.id_cliente = a.id_cliente").
    Join("phones p").On("p.id_cliente = c.id_cliente").
    Join("mailes m").On("m.id_cliente = m.id_cliente").
    And("c.estado_cliente = ?", 0).    
    LeftJoin("left_joined lj").On("lj.a1 = c.a1").
    From("client c").
    Build()

r, error := db.Query(query, params...)
```

Note that in this example we purposely shuffled the order of the clauses and yet the query was built correctly

* Conditional elements

Sometimes we need to check for a condition to build dynamic sql

This example adds the column `name` to the query only if the variable `shouldAddName` is true. 

```go
query, _ := Select().
	Col("a1, a2, a3").
	ColIf(shouldAddName, "name")
	From("client").
	Build()

// Produces "SELECT a1, a2, a3 FROM client" when shouldAddName is false
// Produces "SELECT a1, a2, a3, name FROM client" when shouldAddName is true
```

This also can be applied to joins.

```go
query, _ := Select().
	Col("*").
	From("client c").
	Join("addresses a").On("a.client_id = c.client_id").
    JoinIf(shouldGetPhones, "phones p ON p.client_id = c.client_id").
    Build()

// Produces "SELECT * FROM client c JOIN a.client_id = c.client_id" if shouldGetPhones is false
// Produces "SELECT * FROM client c JOIN a.client_id = c.client_id JOIN phones p ON p.client_id = c.client_id" " if shouldGetPhones is true
```

And boolean connectors

```go
query, _ := Select().
	Col("*").
	From("client c").	
	Where("c.status = 0").AndIf(shouldFilterByCountry, "country = 'CL'").
    Build()

// Produces "SELECT * FROM client c WHERE c.status = 0" when shouldFilterByCountry is false
// Produces "SELECT * FROM client c WHERE c.status = 0 AND country = 'CL'" when shouldFilterByCountry is true
```

* Params

You can add params to almost any clause

```go
query, params := Select().
	Col("name, mail, ? AS max_credit", 1000000).
	From("client c").	
	Where("c.status = 0").And("country = ?", "CL").
    Limit("?", "100").Build()
```


### Delete

* Simple delete

```go
query, _ := Delete().From("client").Build()
// Produces "DELETE FROM client"
```

* Simple del where

```go
query, _ := Delete().From("client").Where("client_id = 100").Build()
// Produces "DELETE FROM client WHERE client_id = 100"
```

* Like with Select you can use parameters and conditionals with Delete

```go
query, params := Delete().From("client").Where("client_id = ?", 100).Build()
// Produces "DELETE FROM client WHERE client_id = ?"
```

```go
query, params := Delete().From("client").Where("1=1").AndIf(filterByClient, "client_id = ?", 100).Build()
// Produces "DELETE FROM client WHERE 1=1" when filterByClient is false
// Produces "DELETE FROM client WHERE 1=1 AND client_id = ?" when filterByClient is true
```


### Update

* Simple update

```go
query, _ := Update("client").Set("status = 0").Build()
// Produces UPDATE client SET status = 0
```

* Update/where/order/limit

```go
query, _ := Update("client").
	Set("status = 0").
	Where("status = ?", 1).
	OrderBy("ciudad").
	Limit(10).
	Build()
```

* You can use obreron to build an update/join query

```go
query, _ := obreron.Update("business AS b").
Join("business_geocode AS g").On("b.business_id = g.business_id").
Set("b.mapx = g.latitude, b.mapy = g.longitude").
Where("(b.mapx = '' or b.mapx = 0)").
And("g.latitude > 0").
Build()

// Produces "UPDATE business AS b JOIN business_geocode AS g ON b.business_id = g.business_id SET b.mapx = g.latitude, b.mapy = g.longitude WHERE (b.mapx = '' or b.mapx = 0) AND g.latitude > 0"
```

* You can use obreron to build an update/select query

```go
query, _ := obreron.Update("items").
				ColSelect(Select().Col("id, retail / wholesale AS markup, quantity").From("items"), "discounted").
				Set("items.retail = items.retail * 0.9").
				Where("discounted.markup >= 1.3").
				And("discounted.quantity < 100").
				And("items.id = discounted.id").
	            Build()
// Produces UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id
```

## Insert

* Simple insert

```go
query, params := Insert().
	Into("client").
	Col("name, value", "'some name'", "'somemail@mail.net'").
    Build()

// Produces "INSERT INTO client ( name, value ) VALUES ( ?, ? )"
```

* insert select

```go
query, params := nsert().
    Into("courses").
    ColSelect("name, location, gid", 
		Select().
		Col("name, location, 1").
	    From("courses").
	    Where("cid = 2")
	).Build()

// Produces       "INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2"
```

## Other clauses

You can add others clauses using the `Clause` method 

```go
query, params := Insert().Clause("IGNORE", "")
	Into("client").
	Col("name, value", "'some name'", "'somemail@mail.net'").
    Build()

// Produces "INSERT IGNORE INTO client ( name, value ) VALUES ( ?, ? )"
```

The `Clause` method always will inject the clause after the last uses building command

