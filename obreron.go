// Package obreron es un paquete que permite construir consultas sql dinámicamente
package obreron

import (
	"bytes"
	"fmt"
	"unsafe"
)

// parseOpts enmascara el tipo bool para usar constantes para las opciones de parseo de parámetros
type parseOpts bool
type parseOptsMul int8

const (
	// Enclose indica poner un `(` antes del siguiente estamento
	Enclose = parseOptsMul(1)

	// NoEnclose indica NO poner un `(` antes del siguiente estamento
	NoEnclose = parseOptsMul(2)

	// EncloseOnlyBuilders indica rodear de parentesis solo a SQLBuilders
	EncloseOnlyBuilders = parseOptsMul(3)

	// UseAs indica agregar una clausula `AS` antes del alias, si hubiera
	UseAs = parseOpts(true)

	// NoUseAs indica NO agregar una clausula `AS` antes del alias, si hubiera
	NoUseAs = parseOpts(false)

	// Quote indica si escapar con backticks el siguiente identificador
	Quote = parseOpts(true)

	// NoQuote indica NO escapar con backticks el siguiente identificador
	NoQuote = parseOpts(false)
)

// parsingOptions contiene opciones de parseo de elementos del query builder
type parsingOptions struct {
	AfterWork func(b *SQLBuilder) error
	Alias     string
	Operator  string
	On        string
	Enclose   parseOptsMul
	UseAS     parseOpts
	Quote     parseOpts
}

// newParsingOpts e Enclose, q Quote, u UseAs, a Alias, o Operator, on clausula ON
func newParsingOpts(e parseOptsMul, q, u parseOpts, a string, o string, on string) *parsingOptions {
	return &parsingOptions{
		Enclose:  e,
		Alias:    a,
		Operator: o,
		UseAS:    u,
		Quote:    q,
		On:       on,
		AfterWork: func(b *SQLBuilder) error {
			return nil
		},
	}
}

// SQLBuilder engloba bytes.Buffer y le da el poder de tener dialecto y parámetros
type SQLBuilder struct {
	dialect Dialect
	bytes.Buffer
	params []interface{}
}

func newSQLBuilder(d Dialect) *SQLBuilder {
	return &SQLBuilder{
		Buffer:  bytes.Buffer{},
		dialect: d,
	}
}

// Build construye la consulta devolviendo una tupla conteniendola en un string y los parámetros
// registrados para su uso
func (sb *SQLBuilder) Build() (string, []interface{}) {
	return *(*string)(unsafe.Pointer(&sb.Buffer)), sb.params
}

// AddParam agrega un paràmetro al SQLBuilder
func (sb *SQLBuilder) AddParam(p ...interface{}) {
	if len(p) > 0 {
		if len(sb.params) == 0 {
			sb.params = make([]interface{}, 0, len(p))
		}
		sb.params = append(sb.params, p...)
	}
}

// Params devuelve el slice de parámetros del SQLBuilder
func (sb *SQLBuilder) Params() []interface{} {
	return sb.params
}

// Dialect devuelve el dialecto de la consulta
func (sb *SQLBuilder) Dialect() Dialect {
	return sb.dialect
}

// ResetParams resetea el slices de parámetros
func (sb *SQLBuilder) ResetParams() {
	sb.params = make([]interface{}, 0)
}

// Select es el builder para consultas que tienen datos
type Select struct {
	columns *SQLBuilder
	joins   *SQLBuilder
	filter  *SQLBuilder
	order   *SQLBuilder
	source  *SQLBuilder
	group   *SQLBuilder
	having  *SQLBuilder

	*SQLBuilder

	limit  int64
	offset int64

	q string
}

// NewMaryBuilder devuelve un nuevo sql builder listo para trabajar
func NewMaryBuilder() *Select {

	d := Mysql{}
	s := Select{
		columns:    newSQLBuilder(d),
		joins:      newSQLBuilder(d),
		filter:     newSQLBuilder(d),
		order:      newSQLBuilder(d),
		source:     newSQLBuilder(d),
		group:      newSQLBuilder(d),
		having:     newSQLBuilder(d),
		SQLBuilder: newSQLBuilder(d),
		limit:      -1,
		offset:     -1,
	}
	return &s
}

func (s *Select) Reset() {
	s.columns.Reset()
	s.joins.Reset()
	s.filter.Reset()
	s.order.Reset()
	s.source.Reset()
	s.group.Reset()
	s.having.Reset()
	s.SQLBuilder.Reset()
	s.columns.Reset()
	s.joins.ResetParams()
	s.filter.ResetParams()
	s.order.ResetParams()
	s.source.ResetParams()
	s.group.ResetParams()
	s.having.ResetParams()
	s.SQLBuilder.ResetParams()
	s.limit = -1
	s.offset = -1
}

// Params devuelve los paramétros registrados para los componentes de la consulta
// en el orden esperado para las distintistas clausulas
func (s *Select) Params() []interface{} {
	sz := struct {
		// size cantidad total de paramétros a recibir
		size int
		// lc lc paràmetros de columns
		lc int
		ls int
		lj int
		lf int
	}{}

	if l := len(s.columns.params); l > 0 {
		sz.lc = l
		sz.size += l
	}

	if l := len(s.source.params); l > 0 {
		sz.ls = l
		sz.size += l
	}

	if l := len(s.joins.params); l > 0 {
		sz.lj = l
		sz.size += l
	}

	if l := len(s.filter.params); l > 0 {
		sz.lf = l
		sz.size += l
	}

	s.params = make([]interface{}, 0, sz.size)

	if sz.lc > 0 {
		s.params = append(s.params, s.columns.params...)
	}

	if sz.ls > 0 {
		s.params = append(s.params, s.source.params...)
	}

	if sz.lj > 0 {
		s.params = append(s.params, s.joins.params...)
	}

	if sz.lf > 0 {
		s.params = append(s.params, s.filter.params...)
	}
	return s.params
}

// Limit establece el limite de la consulta. Si este valor es -1 no se agregara la clausula OFFSET a la query construida
func (s *Select) Limit(l int64) *Select {
	s.q = ""
	s.limit = l
	return s
}

// Offset establece el offset de la consulta. Si este valor es -1 no se agregara la clausula OFFSET a la query construida
func (s *Select) Offset(o int64) *Select {
	s.q = ""
	s.offset = o
	return s
}

// Select define consultas para la consulta. Cada ve que se llama resetea el buffer de construción
func (s *Select) Select(cs ...interface{}) *Select {
	s.q = ""
	s.columns.Reset()
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS
	opt := newParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, "", "", "")

	for _, c := range cs {
		parse(s.columns, c, nil, opt)
		s.columns.WriteByte(44)
	}

	return s
}

// AddColumn agrega una columna con su alias. la columna c puede ser string u otro SQLBuilder. Puede omitir el alias pasando un string vacio
func (s *Select) AddColumn(c interface{}, a string) *Select {
	s.q = ""
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS
	parse(s.columns, c, nil, newParsingOpts(EncloseOnlyBuilders, NoQuote, UseAs, a, "", ""))
	s.columns.WriteByte(44)
	return s
}

// AddColumnIf agrega una columna con su alias si se cumple condición cond. la columna c puede ser string u otro SQLBuilder. Puede omitir el alias pasando un string vacio
func (s *Select) AddColumnIf(cond bool, c interface{}, a string) *Select {
	if cond {
		s.AddColumn(c, a)
	}
	return s
}

// From define el origen para obtener los datos de la consulta. Puede ser un string u otro sql builder
func (s *Select) From(source interface{}, a string) *Select {
	s.q = ""
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y usar clausula AS solo si se definio un alias
	s.source.WriteString(" FROM ")
	parse(s.source, source, nil, newParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, a, "", ""))
	return s
}

// Inner Agrega un inner join a la construcción de la query. El joinable c puede ser string o un SQLBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) Inner(c interface{}, a string, on string) *Select {
	return s.join(" INNER JOIN ", c, a, on)
}

// Left Agrega un left join a la construcción de la query. El joinable c puede ser string o un SQLBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) Left(c interface{}, a string, on string) *Select {
	return s.join(" LEFT JOIN ", c, a, on)
}

// Right Agrega un right join a la construcción de la query. El joinable c puede ser string o un SQLBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) Right(c interface{}, a string, on string) *Select {
	return s.join(" RIGHT JOIN ", c, a, on)
}

// RightIF Agrega un right join a la construcción de la query si la condición `cond` es verdadera.
// El joinable c puede ser string o un SQLBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) RightIF(cond bool, c interface{}, a string, on string) *Select {
	if cond {
		return s.join(" RIGHT JOIN ", c, a, on)
	}
	return s
}

// InnerIf Agrega un inner join a la construcción de la query si la condición `cond` es verdadera.
// El joinable c puede ser string o un SQLBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) InnerIf(cond bool, c interface{}, a string, on string) *Select {
	if cond {
		return s.join(" INNER JOIN ", c, a, on)
	}
	return s
}

// LeftIf Agrega un left join a la construcción de la query si la condición `cond` es verdadera.
// El joinable c puede ser string o un SQLBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) LeftIf(cond bool, c interface{}, a string, on string) *Select {
	if cond {
		return s.join(" LEFT JOIN ", c, a, on)
	}
	return s
}

// join es un método helper privado que ayuda a la construcción de joines
func (s *Select) join(j string, c interface{}, a string, on string) *Select {
	s.q = ""
	s.joins.WriteString(j)
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS usando alias solo si se definio, pasando el on
	parse(s.joins, c, "", newParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, a, "", on))
	return s
}

// GroupBy agrega la clausula GROUP BY al Sql Builder
func (s *Select) GroupBy(c string) *Select {
	s.q = ""
	s.group.WriteString(fmt.Sprintf(" GROUP BY %v ", c))
	return s
}

// Having agrega la clausula HAVING al Sql Builder
func (s *Select) Having(c string) *Select {
	s.q = ""
	s.group.WriteString(fmt.Sprintf(" HAVING %v ", c))
	return s
}

// OrderBy agrega la clausula ORDER BY al Sql Builder
func (s *Select) OrderBy(c string) *Select {
	s.q = ""
	s.order.WriteString(fmt.Sprintf(" ORDER BY %v ", c))
	return s
}

// Where inicializa la clausula where
func (s *Select) Where() *Select {
	s.q = ""
	s.filter.Reset()

	s.filter.WriteString(" WHERE 1=1 ")
	return s
}

// AndParam Agrega una condición usando conector AND
// c puede ser la condición como string o como un SQLBuilder
// op es el operador y param el parámetro de la condición
func (s *Select) AndParam(c interface{}, op string, param interface{}) *Select {
	s.q = ""
	s.filter.WriteString(" AND ")
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS ni  alias
	parse(s.filter, c, param, newParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, "", op, ""))

	return s
}

// AndParamIf Agrega una condición usando conector AND solo si cond es true
// c puede ser la condición como string o como un SQLBuilder
// op es el operador y param el parámetro de la condición
func (s *Select) AndParamIf(cond bool, c interface{}, op string, param interface{}) *Select {
	if cond {
		s.AndParam(c, op, param)
	}
	return s
}

// And Agrega una condición usando conector AND
// c es un strig conteniendo la condición completa
func (s *Select) And(c string) *Select {
	s.AndParam(c, "", nil)
	return s
}

// AndIf Agrega una condición usando conector AND solo si cond es true
// c es un strig conteniendo la condición completa
func (s *Select) AndIf(cond bool, c string) *Select {
	if cond {
		s.AndParam(c, "", nil)
	}
	return s
}

// OrParam Agrega una condición usando conector AND
// c puede ser la condición como string o como un SQLBuilder
// op es el operador y param el parámetro de la condición
func (s *Select) OrParam(c interface{}, op string, param interface{}) *Select {
	s.q = ""
	s.filter.WriteString(" OR ")
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS ni  alias
	parse(s.filter, c, param, newParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, "", op, ""))
	return s
}

// Quote escapa a su argumento según el dialecto de la consulta
func (s *Select) Quote(c interface{}) string {
	return s.SQLBuilder.dialect.Quote(c)
}

// OpenEnclose Agrega un abre parentesis ( la consulta
func (s *Select) OpenEnclose() string {
	return s.SQLBuilder.dialect.OpenEnclose()
}

// CloseEnclose Agrega un cierre de parentesis ) la consulta
func (s *Select) CloseEnclose() string {
	return s.SQLBuilder.dialect.CloseEnclose()
}

func (s *Select) String() string {
	if s.q != "" {
		return s.q
	}

	s.WriteString("SELECT ")

	if s.columns.Len() > 1 {
		s.Write(s.columns.Bytes()[0 : s.columns.Len()-1])
	}
	s.Write(s.source.Bytes())

	s.Write(s.joins.Bytes())

	if s.filter.Len() > 0 {
		s.Write(s.filter.Bytes())
	}

	if s.group.Len() > 0 {
		s.Write(s.group.Bytes())
	}

	if s.having.Len() > 0 {
		s.Write(s.having.Bytes())
	}

	if s.order.Len() > 0 {
		s.Write(s.order.Bytes())
	}

	if s.limit > -1 {
		s.WriteString(fmt.Sprintf(" LIMIT %v ", s.limit))
	}

	if s.offset > -1 {
		s.WriteString(fmt.Sprintf(" OFFSET %v ", s.offset))
	}

	s.q = *(*string)(unsafe.Pointer(&s.Buffer))

	return s.q
}

// Build construye la consulta devolviendo una tupla conteniendola en un string y los parámetros
// registrados para su uso
func (s *Select) Build() (string, []interface{}) {
	return s.String(), s.Params()
}

// parse parsea elementos comunes de la consulta
func parse(subject *SQLBuilder, circumstance interface{}, parameter interface{}, opt *parsingOptions) *SQLBuilder {

	openHook(subject, circumstance, parameter, opt)

	switch circumstance.(type) {
	default:
		return subject
	case string:
		parseString(subject, circumstance, parameter, opt)
	case *SQLBuilder:
		parseBuilder(subject, circumstance, parameter, opt)
	case *Select:
		parseSelect(subject, circumstance, parameter, opt)
	}

	closeHook(subject, circumstance, parameter, opt)

	return subject
}

// parseString parses a circumstance as string
func parseString(subject *SQLBuilder, circumstance interface{}, parameter interface{}, opt *parsingOptions) {
	sc := circumstance.(string)

	// if Quote try to escape builder result with backticks
	if opt.Quote {
		sc = subject.Dialect().Quote(sc)
	}

	_, _ = subject.WriteString(sc)
}

// parseBuilder parses a circumstance as  *SQLBuilder
func parseBuilder(subject *SQLBuilder, circumstance interface{}, parameter interface{}, opt *parsingOptions) {
	b := circumstance.(*SQLBuilder)

	if opt.Enclose == EncloseOnlyBuilders {
		_, _ = subject.WriteString(subject.Dialect().OpenEnclose())
	}

	// if Quote try to escape builder result wit backticks
	if opt.Quote {
		_, _ = subject.WriteString(
			subject.Dialect().Quote(
				b.String(),
			),
		)
	} else {
		_, _ = subject.Write(b.Bytes())
	}

	if opt.Enclose == EncloseOnlyBuilders {
		_, _ = subject.WriteString(subject.Dialect().CloseEnclose())
	}

	subject.AddParam(b.Params()...)

	// // carga los parámetros de la query en los params de master, normalmente *Select
	// master.AddParam(b.Params()...)
	// // después de cargarlos, resetee los paràmetros de la circunstancia, para evitar memory leaks.
	// b.ResetParams()
	// b.Reset()
}

func parseSelect(subject *SQLBuilder, circumstance interface{}, parameter interface{}, opt *parsingOptions) {
	smt := circumstance.(*Select)

	if opt.Enclose == EncloseOnlyBuilders {
		_, _ = subject.WriteString(subject.Dialect().OpenEnclose())
	}

	// if Quote try to escape builder result wit backticks
	if opt.Quote {
		_, _ = subject.WriteString(
			subject.Dialect().Quote(
				*(*string)(unsafe.Pointer(&smt.Buffer)),
			),
		)
	} else {
		_, _ = subject.WriteString(smt.String())
	}

	if opt.Enclose == EncloseOnlyBuilders {
		_, _ = subject.WriteString(subject.Dialect().CloseEnclose())
	}

	subject.AddParam(smt.Params()...)
}

// openHook concentra el proceso antes del parseo
func openHook(subject *SQLBuilder, circumstance interface{}, parameter interface{}, opt *parsingOptions) {
	// Este hack permite rodear de parentesis solo a circunstancias que sean Builders
	if opt.Enclose == Enclose {
		_, _ = subject.WriteString(subject.Dialect().OpenEnclose())
	}
}

// closeHook concentra el proceso después del parseo
func closeHook(subject *SQLBuilder, circumstance interface{}, parameter interface{}, opt *parsingOptions) {
	if opt.Enclose == Enclose {
		_, _ = subject.WriteString(subject.Dialect().CloseEnclose())
	}

	if opt.Alias != "" {
		if opt.UseAS {
			_, _ = subject.WriteString(fmt.Sprintf(" AS %v ", opt.Alias))
		} else {
			_, _ = subject.WriteString(fmt.Sprintf(" %v ", opt.Alias))
		}
	}

	if opt.On != "" {
		_, _ = subject.WriteString(fmt.Sprintf(" ON %v", opt.On))

		// este if es para agregar los posibles parametros en una clausula on
		if parameter != nil && parameter != "" {
			subject.AddParam(parameter)
		}
	}

	if opt.Operator != "" {
		_, _ = subject.WriteString(fmt.Sprintf(" %v ", opt.Operator))
	}

	// // este if es para agregar los posibles parametros en una clausula WHERE
	// Cuidado! se debe mantener este orden para que funcione el hack de agregar la marca de parámetro junto con el paramètro.
	// justo después de agregar al operador
	if parameter != nil && (opt.On == "" && parameter != "") {
		_, _ = subject.WriteString(subject.Dialect().ParamMark())
		subject.AddParam(parameter)
	}

	_ = opt.AfterWork(subject)
}
