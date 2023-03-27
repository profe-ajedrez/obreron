package obreron

import "fmt"

var (
	_ Dialect = Mysql{}
)

type Dialect interface {
	Quote(v interface{}) string
	ParamMark() string
	OpenEnclose() string
	CloseEnclose() string
}

type Mysql struct{}

func (m Mysql) Quote(v interface{}) string {
	return fmt.Sprintf("`%v`", v)
}

func (m Mysql) ParamMark() string {
	return "?"
}

func (m Mysql) OpenEnclose() string {
	return "("
}

func (m Mysql) CloseEnclose() string {
	return ")"
}
