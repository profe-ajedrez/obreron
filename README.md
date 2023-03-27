# Obreron

Sql Builder escrito en Go.

---

## Dialectos soportados

* Mysql

---


## Instalación

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
SELECT id,name,mail, `columna con espacios en el nombre` FROM users AS u 
```

