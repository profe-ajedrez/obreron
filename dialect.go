package obreron

import "fmt"

var (
	_ Dialect = Mysql{}
)

// Dialect representa el dialecto de la consulta
type Dialect interface {
	Quote(v interface{}) string
	ParamMark() string
	OpenEnclose() string
	CloseEnclose() string
}

// Mysql es un dialecto que permite construir consultas para mysql y mariadb
type Mysql struct{}

// Quote escapa a su argumento según el dialecto de la consulta
func (m Mysql) Quote(v interface{}) string {
	return fmt.Sprintf("`%v`", v)
}

// ParamMark devuelve una marca de parámetro posicional
func (m Mysql) ParamMark() string {
	return "?"
}

// OpenEnclose Agrega un abre parentesis ( la consulta
func (m Mysql) OpenEnclose() string {
	return "("
}

// CloseEnclose Agrega un cierre de parentesis ) la consulta
func (m Mysql) CloseEnclose() string {
	return ")"
}
