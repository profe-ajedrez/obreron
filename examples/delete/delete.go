package main

import (
	"fmt"

	v2 "github.com/profe-ajedrez/obreron/v2"
)

func main() {
	// Simple delete
	fmt.Println("Simple Delete")
	q, p := v2.Delete().From("client").Build()
	fmt.Printf("query: %s \nparams: %v\n\n", q, p)
	//	OUTPUT:
	//	query: DELETE FROM client
	//	params: []

	// Delete with conditions
	fmt.Println("Delete With conditions")
	q, p = v2.Delete().From("client").
		Where("client_id = 100").
		And("estado_cliente = 0").
		Y().In("regime_cliente", "'G01','G02', ?", "'G03'").
		And("a").
		Build()

	fmt.Printf("query: %s \nparams: %v\n\n", q, p)
	//	OUTPUT:
	//	query: DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) AND a
	//	params: ['G03']
}
