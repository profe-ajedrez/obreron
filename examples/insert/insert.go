package main

import (
	"fmt"

	v2 "github.com/profe-ajedrez/obreron/v2"
)

func main() {
	// Simple Insert
	fmt.Println("Simple Insert")
	q, p := v2.Insert().
		Into("client").
		Col("name, mail", "some name", "somemail@mail.net").
		Build()
	fmt.Printf("query: %s \nparams: %v\n\n", q, p)
	//	OUTPUT:
	//	qquery: INSERT INTO client ( name, mail ) VALUES ( ?,? )
	//	nparams: [some name somemail@mail.net]

	// Conditional insert
	fmt.Println("Conditional Insert")
	name := ""
	mail := "somemail@mail.net"

	q, p = v2.Insert().
		Into("client").
		ColIf(len(name) > 0, "name", name).
		ColIf(len(mail) > 0, "mail", mail).
		Build()
	fmt.Printf("query: %s \nparams: %v\n\n", q, p)
	//	OUTPUT:
	//	query: INSERT INTO client ( mail ) VALUES ( ? )
	//	params: [somemail@mail.net]

	// Insert Select
	fmt.Println("Insert Select")
	q, p = v2.Insert().
		Into("courses").
		ColSelect("name, location, gid",
			v2.Select().
				Col("name, location, 1").
				From("courses").
				Where("cid = 2"),
		).Build()
	fmt.Printf("query: %s \nparams: %v\n", q, p)
	//	OUTPUT
	//	query: INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2
	//	params: []
}
