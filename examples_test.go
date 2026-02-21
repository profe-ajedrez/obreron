package obreron_test

import (
	"fmt"

	obreron "github.com/profe-ajedrez/obreron/v2"
)

func ExampleSelect_basic() {
	q, p := obreron.Select().
		Col("a1, a2, a3").
		From("client").
		Build()

	fmt.Println(q)
	fmt.Println(p)

	// Output:
	// SELECT a1, a2, a3 FROM client
	// []
}

func ExampleSelect_inArgsEmpty() {
	q, p := obreron.Select().
		Col("a1, a2, a3").
		From("client").
		Where("1 = 1").
		Y().
		InArgs("status").
		Build()

	fmt.Println(q)
	fmt.Println(p)

	// Output:
	// SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND status IN (NULL)
	// []
}

func ExampleInsert_ignore() {
	q, p := obreron.Insert().
		Ignore().
		Into("client").
		Col("name, value", "'some name'", "'somemail@mail.net'").
		Build()

	fmt.Println(q)
	fmt.Println(p)

	// Output:
	// INSERT IGNORE INTO client ( name, value ) VALUES ( ?,? )
	// ['some name' 'somemail@mail.net']
}

func ExampleUpdate_setMultiple() {
	q, p := obreron.Update("client").
		Set("status = 0").
		Set("name = ?", "stitch").
		Build()

	fmt.Println(q)
	fmt.Println(p)

	// Output:
	// UPDATE client SET status = 0, name = ?
	// [stitch]
}

func ExampleDelete_basic() {
	q, p := obreron.Delete().
		From("client").
		Build()

	fmt.Println(q)
	fmt.Println(p)

	// Output:
	// DELETE FROM client
	// []
}
