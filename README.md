![image](https://github.com/user-attachments/assets/c4a6b204-1aab-49a5-aeb4-950067a82e6d)


# Obreron

Fast and cheap sql builder.

[![Go Reference](https://pkg.go.dev/badge/github.com/profe-ajedrez/obreron/v2.svg)](https://pkg.go.dev/github.com/profe-ajedrez/obreron/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/profe-ajedrez/obreron/v2)](https://goreportcard.com/report/github.com/profe-ajedrez/obreron/v2)
[![Coverage Status](https://coveralls.io/repos/github/profe-ajedrez/obreron/badge.svg?branch=v2)](https://coveralls.io/github/profe-ajedrez/obreron?branch=v2)

## Supported Dialects

- [x] Mysql
- [ ] Postgresql (Work in progress)


## Why?

Because writing SQL in Go is boring.

## Instalation

Use `go get` to install v2

```bash
$ go get github.com/profe-ajedrez/obreron/v2
```

## Use

You could see the [examples](examples/) directory.

Import package

```go
import (
	v2 "github.com/profe-ajedrez/obreron/v2"
)
```

### Select

* Simple select

```go
// Produces SELECT a1, a2, a3 FROM client
query, _ := v2.Select().Col("a1, a2, a3").From("client").Build()
r, error := db.Query(query)
```

* Select/join/where/shuffled

```go
// Produces SELECT a1, a2, ? AS diez, colIf1, colIf2, ? AS zero, a3, ? AS cien FROM client c JOIN addresses a ON a.id_cliente = a.id_cliente JOIN phones p ON p.id_cliente = c.id_cliente JOIN mailes m ON m.id_cliente = m.id_cliente AND c.estado_cliente = ? LEFT JOIN left_joined lj ON lj.a1 = c.a1 WHERE a1 = ? AND a2 = ? AND a3 = 10 AND a16 = ?
// with params = []any{10, 0, 100, 0, "'last name'", 1000.54, 75}
query, params := v2.Select().
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
query, _ := v2.Select().
	Col("a1, a2, a3").
	ColIf(shouldAddName, "name")
	From("client").
	Build()

// Produces "SELECT a1, a2, a3 FROM client" when shouldAddName is false
// Produces "SELECT a1, a2, a3, name FROM client" when shouldAddName is true
```

This also can be applied to joins.

```go
query, _ := v2.Select().
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
query, _ := v2.Select().
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
query, params := v2.Select().
	Col("name, mail, ? AS max_credit", 1000000).
	From("client c").	
	Where("c.status = 0").And("country = ?", "CL").
    Limit("?", "100").Build()
```


### Delete

* Simple delete

```go
query, _ := v2.Delete().From("client").Build()
// Produces "DELETE FROM client"
```

* Simple del where

```go
query, _ := v2.Delete().From("client").Where("client_id = 100").Build()
// Produces "DELETE FROM client WHERE client_id = 100"
```

* Like with Select you can use parameters and conditionals with Delete

```go
query, params := v2.Delete().From("client").Where("client_id = ?", 100).Build()
// Produces "DELETE FROM client WHERE client_id = ?"
```

```go
query, params := v2.Delete().From("client").Where("1=1").AndIf(filterByClient, "client_id = ?", 100).Build()
// Produces "DELETE FROM client WHERE 1=1" when filterByClient is false
// Produces "DELETE FROM client WHERE 1=1 AND client_id = ?" when filterByClient is true
```


### Update

* Simple update

```go
query, _ := v2.Update("client").Set("status = 0").Build()
// Produces UPDATE client SET status = 0
```

* Update/where/order/limit

```go
query, _ := v2.Update("client").
	Set("status = 0").
	Where("status = ?", 1).
	OrderBy("ciudad").
	Limit(10).
	Build()
```

* You can use obreron to build an update/join query

```go
query, _ := v2.Update("business AS b").
Join("business_geocode AS g").On("b.business_id = g.business_id").
Set("b.mapx = g.latitude, b.mapy = g.longitude").
Where("(b.mapx = '' or b.mapx = 0)").
And("g.latitude > 0").
Build()

// Produces "UPDATE business AS b JOIN business_geocode AS g ON b.business_id = g.business_id SET b.mapx = g.latitude, b.mapy = g.longitude WHERE (b.mapx = '' or b.mapx = 0) AND g.latitude > 0"
```

* You can use obreron to build an update/select query

```go
query, _ := v2.Update("items").
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
query, params := Iv2.nsert().
	Into("client").
	Col("name, value", "'some name'", "'somemail@mail.net'").
    Build()

// Produces "INSERT INTO client ( name, value ) VALUES ( ?, ? )"
```

* insert select

```go
query, params := v2.Insert().
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
query, params := v2.Insert().Clause("IGNORE", "")
	Into("client").
	Col("name, value", "'some name'", "'somemail@mail.net'").
    Build()

// Produces "INSERT IGNORE INTO client ( name, value ) VALUES ( ?, ? )"
```

The `Clause` method always will inject the clause after the last uses building command
