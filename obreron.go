package obreron

import (
	"fmt"
	"strings"
)

var (
	// compile time check para garantizar que SqlBuilder cumple con la interface
	_ Builder = &SqlBuilder{}
)

// parseOpts enmascara el tipo bool para usar constantes para las opciones de parseo de parámetros
type parseOpts bool
type parseOptsMul int8

const (
	// Enclose indica poner un `(` antes del siguiente estamento
	Enclose = parseOptsMul(1)

	// NoEnclose indica NO poner un `(` antes del siguiente estamento
	NoEnclose = parseOptsMul(2)

	EncloseOnlyBuilders = parseOptsMul(3)

	// UseAs indica agregar una clausula `AS` antes del alias, si hubiera
	UseAs = parseOpts(true)

	// NoUseAs indica NO agregar una clausula `AS` antes del alias, si hubiera
	NoUseAs = parseOpts(false)

	// Quote indica si escapar con backticks el siguiente identificador
	Quote = parseOpts(true)

	// Quote indica NO escapar con backticks el siguiente identificador
	NoQuote = parseOpts(false)
)

// parsingOptions contiene opciones de parseo de elementos del query builder
type parsingOptions struct {
	AfterWork func(b Builder) error
	Alias     string
	Operator  string
	On        string
	Enclose   parseOptsMul
	UseAS     parseOpts
	Quote     parseOpts
}

// NewParsingOpts e Enclose, q Quote, u UseAs, a Alias, o Operator, on clausula ON
func NewParsingOpts(e parseOptsMul, q, u parseOpts, a string, o string, on string) parsingOptions {
	return parsingOptions{
		Enclose:  e,
		Alias:    a,
		Operator: o,
		UseAS:    u,
		Quote:    q,
		On:       on,
		AfterWork: func(b Builder) error {
			return nil
		},
	}
}

// Builder es una interface que representa una cosa capaz de agregar strings a un buffer y de construir algo con esos string agregados
type Builder interface {
	// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
	Len() int

	// Cap returns the capacity of the builder's underlying byte slice. It is the
	// total space allocated for the string being built and includes any bytes
	// already written.
	Cap() int

	// Grow grows b's capacity, if necessary, to guarantee space for
	// another n bytes. After Grow(n), at least n bytes can be written to b
	// without another allocation. If n is negative, Grow panics.
	Grow(n int)

	// Write appends the contents of p to b's buffer.
	// Write always returns len(p), nil.
	Write(p []byte) (int, error)

	// WriteByte appends the byte c to b's buffer.
	// The returned error is always nil.
	WriteByte(c byte) error

	// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
	// It returns the length of r and a nil error.
	WriteRune(r rune) (int, error)

	// WriteString appends the contents of s to b's buffer.
	// It returns the length of s and a nil error.
	WriteString(s string) (int, error)

	// String returns the accumulated string.
	String() string

	// Params retorna un slice con los paramétros que reemplazaran los marcadores en el buffer
	Params() []interface{}

	// Reset resets the Builder to be empty.
	Reset()

	// Build devuelve el string construido por obreron.Builder y los paramétros agregados para reemplazar a los marcadores en el buffer
	Build() (string, []interface{})

	// AddParam Agrega uno o mas parámetros para reemplazar los marcadores en el buffer
	AddParam(p ...interface{})

	// Dialect Retorna el dialecto usado por el builder
	Dialect() Dialect

	ResetParams()
}

// SqlBuilder engloba strings.Builder y le da el poder de tener dialecto y parámetros
type SqlBuilder struct {
	dialect Dialect
	strings.Builder
	params []interface{}
}

func newSqlBuilder(d Dialect) *SqlBuilder {
	return &SqlBuilder{
		Builder: strings.Builder{},
		params:  make([]interface{}, 0),
		dialect: d,
	}
}

func (sb *SqlBuilder) Build() (string, []interface{}) {
	return sb.String(), sb.params
}

func (sb *SqlBuilder) AddParam(p ...interface{}) {
	sb.params = append(sb.params, p...)
}

func (sb *SqlBuilder) Params() []interface{} {
	return sb.params
}

func (sb *SqlBuilder) Dialect() Dialect {
	return sb.dialect
}

func (sb *SqlBuilder) ResetParams() {
	sb.params = nil
}

// Select es el builder para consultas que tienen datos
type Select struct {
	columns *SqlBuilder
	joins   *SqlBuilder
	filter  *SqlBuilder
	order   *SqlBuilder
	source  *SqlBuilder
	group   *SqlBuilder
	having  *SqlBuilder

	*SqlBuilder

	limit  int64
	offset int64
}

// NewMaryBuilder devuelve un nuevo sql builder listo para trabajar
func NewMaryBuilder() *Select {

	d := Mysql{}
	s := Select{
		columns:    newSqlBuilder(d),
		joins:      newSqlBuilder(d),
		filter:     newSqlBuilder(d),
		order:      newSqlBuilder(d),
		source:     newSqlBuilder(d),
		group:      newSqlBuilder(d),
		having:     newSqlBuilder(d),
		SqlBuilder: newSqlBuilder(d),
		limit:      -1,
		offset:     -1,
	}
	return &s
}

func (s *Select) Params() []interface{} {
	sz := struct {
		size int
		lc   int
		ls   int
		lj   int
		lf   int
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
	s.limit = l
	return s
}

// Offset establece el offset de la consulta. Si este valor es -1 no se agregara la clausula OFFSET a la query construida
func (s *Select) Offset(o int64) *Select {
	s.offset = o
	return s
}

// Select define consultas para la consulta. Cada ve que se llama resetea el buffer de construción
func (s *Select) Select(cs ...interface{}) *Select {
	s.columns.Reset()
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS
	opt := NewParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, "", "", "")

	for _, c := range cs {
		parse(s, s.columns, c, nil, opt)
		s.columns.WriteString(",")
	}

	return s
}

// AddColum agrega una columna con su alias. la columna c puede ser string u otro SqlBuilder. Puede omitir el alias pasando un string vacio
func (s *Select) AddColumn(c interface{}, a string) *Select {
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS
	opt := NewParsingOpts(EncloseOnlyBuilders, NoQuote, UseAs, a, "", "")

	parse(s, s.columns, c, nil, opt)
	return s
}

// AddColumIf  agrega una columna con su alias si se cumple condición cond. la columna c puede ser string u otro SqlBuilder. Puede omitir el alias pasando un string vacio
func (s *Select) AddColumnIf(cond bool, c interface{}, a string) *Select {
	if cond {
		s.AddColumn(c, a)
	}
	return s
}

// From define el origen para obtener los datos de la consulta. Puede ser un string u otro sql builder
func (s *Select) From(source interface{}, a string) *Select {
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y usar clausula AS solo si se definio un alias
	opt := NewParsingOpts(EncloseOnlyBuilders, NoQuote, UseAs, a, "", "")
	s.source.WriteString(" FROM ")
	parse(s, s.source, source, nil, opt)
	return s
}

// Inner Agrega un inner join a la construcción de la query. El joinable c puede ser string o un SqlBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) Inner(c interface{}, a string, on string) *Select {
	return s.join(" INNER JOIN ", c, a, on)
}

// Inner Agrega un left join a la construcción de la query. El joinable c puede ser string o un SqlBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) Left(c interface{}, a string, on string) *Select {
	return s.join(" LEFT JOIN ", c, a, on)
}

// Inner Agrega un right join a la construcción de la query. El joinable c puede ser string o un SqlBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) Right(c interface{}, a string, on string) *Select {
	return s.join(" RIGHT JOIN ", c, a, on)
}

// Inner Agrega un right join a la construcción de la query si la condición `cond` es verdadera.
// El joinable c puede ser string o un SqlBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) RightIF(cond bool, c interface{}, a string, on string) *Select {
	if cond {
		return s.join(" RIGHT JOIN ", c, a, on)
	}
	return s
}

// Inner Agrega un inner join a la construcción de la query si la condición `cond` es verdadera.
// El joinable c puede ser string o un SqlBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) InnerIf(cond bool, c interface{}, a string, on string) *Select {
	if cond {
		return s.join(" INNER JOIN ", c, a, on)
	}
	return s
}

// Inner Agrega un left join a la construcción de la query si la condición `cond` es verdadera.
// El joinable c puede ser string o un SqlBuilder
// a es el alias, si no lo necesita puede pasarlo vacio.
// on contiene la condición para la clausula on, puede dejarla vacia para que no se agregue
func (s *Select) LeftIf(cond bool, c interface{}, a string, on string) *Select {
	if cond {
		return s.join(" LEFT JOIN ", c, a, on)
	}
	return s
}

// join es un mètodo helper privado que ayuda a la construcción de joines
func (s *Select) join(j string, c interface{}, a string, on string) *Select {
	s.joins.WriteString(j)
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS usando alias solo si se definio, pasando el on
	opt := NewParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, a, "", on)
	parse(s, s.joins, c, "", opt)
	return s
}

// GroupBy agrega la clausula GROUP BY al Sql Builder
func (s *Select) GroupBy(c string) *Select {
	s.group.WriteString(fmt.Sprintf(" GROUP BY %v ", c))
	return s
}

// Having agrega la clausula HAVING al Sql Builder
func (s *Select) Having(c string) *Select {
	s.group.WriteString(fmt.Sprintf(" HAVING %v ", c))
	return s
}

// OrderBy agrega la clausula ORDER BY al Sql Builder
func (s *Select) OrderBy(c string) *Select {
	s.order.WriteString(fmt.Sprintf(" ORDER BY %v ", c))
	return s
}

// Where inicializa la clausula where
func (s *Select) Where() *Select {
	s.filter.Reset()

	s.filter.WriteString(" WHERE 1=1 ")
	return s
}

// AndParam Agrega una condición usando conector AND
// c puede ser la condición como string o como un SqlBuilder
// op es el operador y param el parámetro de la condición
func (s *Select) AndParam(c interface{}, op string, param interface{}) *Select {
	s.filter.WriteString(" AND ")
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS ni  alias
	opt := NewParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, "", op, "")
	parse(s, s.filter, c, param, opt)

	return s
}

// AndParamIf Agrega una condición usando conector AND solo si cond es true
// c puede ser la condición como string o como un SqlBuilder
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
func (s *Select) AndId(cond bool, c string) *Select {
	if cond {
		s.AndParam(c, "", nil)
	}
	return s
}

// OrParam Agrega una condición usando conector AND
// c puede ser la condición como string o como un SqlBuilder
// op es el operador y param el parámetro de la condición
func (s *Select) OrParam(c interface{}, op string, param interface{}) *Select {
	s.filter.WriteString(" OR ")
	// Para este parseo cerrar entre parenstesis solo a los builders, no escapar y no usar clausula AS ni  alias
	opt := NewParsingOpts(EncloseOnlyBuilders, NoQuote, NoUseAs, "", op, "")
	parse(s, s.filter, c, param, opt)
	return s
}

func (s *Select) Quote(c interface{}) string {
	return s.SqlBuilder.dialect.Quote(c)
}

func (s *Select) OpenEnclose() string {
	return s.SqlBuilder.dialect.OpenEnclose()
}

func (s *Select) CloseEnclose() string {
	return s.SqlBuilder.dialect.CloseEnclose()
}

func (s *Select) String() string {
	s.WriteString("SELECT ")

	if s.columns.Len() > 1 {
		s.WriteString(s.columns.String()[0 : s.columns.Len()-1])
	}
	s.WriteString(s.source.String())
	s.WriteString(s.joins.String())

	if s.filter.Len() > 0 {
		s.WriteString(s.filter.String())
	}

	if s.group.Len() > 0 {
		s.WriteString(s.group.String())
	}

	if s.having.Len() > 0 {
		s.WriteString(s.having.String())
	}

	if s.order.Len() > 0 {
		s.WriteString(s.order.String())
	}

	if s.limit > -1 {
		s.WriteString(fmt.Sprintf(" LIMIT %v ", s.limit))
	}

	if s.offset > -1 {
		s.WriteString(fmt.Sprintf(" OFFSET %v ", s.offset))
	}

	return s.Builder.String()
}

func (s *Select) Build() (string, []interface{}) {
	return s.String(), s.Params()
}

// parse parsea elementos comunes de la consulta
func parse(master Builder, subject Builder, circunstance interface{}, parameter interface{}, opt parsingOptions) Builder {

	openHook(subject, circunstance, parameter, opt)

	switch circunstance.(type) {
	default:
		return subject
	case string:
		parseString(subject, circunstance, parameter, opt)
	case Builder:
		parseBuilder(master, subject, circunstance, parameter, opt)
	}

	closeHook(master, subject, circunstance, parameter, opt)

	return subject
}

// parseString parses a circunstance as string
func parseString(subject Builder, circunstance interface{}, parameter interface{}, opt parsingOptions) {
	sc := circunstance.(string)

	// if Quote try to escape builder result with backticks
	if opt.Quote {
		sc = subject.Dialect().Quote(sc)
	}

	subject.WriteString(sc)
}

// parseBuilder parses a circunstance as Builder
func parseBuilder(master Builder, subject Builder, circunstance interface{}, parameter interface{}, opt parsingOptions) {
	b := circunstance.(Builder)

	if opt.Enclose == EncloseOnlyBuilders {
		subject.WriteString(subject.Dialect().OpenEnclose())
	}

	// if Quote try to escape builder result wit backticks
	if opt.Quote {
		subject.WriteString(
			subject.Dialect().Quote(
				b.String(),
			),
		)
	} else {
		subject.WriteString(b.String())
	}

	if opt.Enclose == EncloseOnlyBuilders {
		subject.WriteString(subject.Dialect().CloseEnclose())
	}

	smt, ok := circunstance.(*Select)

	if ok {
		subject.AddParam(smt.Params()...)
	} else {
		subject.AddParam(b.Params()...)
	}
	// // carga los parámetros de la query en los params de master, normalmente *Select
	// master.AddParam(b.Params()...)
	// // después de cargarlos, resetee los paràmetros de la circunstancia, para evitar memory leaks.
	// b.ResetParams()
	// b.Reset()

}

// openHook concentra el proceso antes del parseo
func openHook(subject Builder, circunstance interface{}, parameter interface{}, opt parsingOptions) {
	// Este hack permite rodear de parentesis solo a circunstancias que sean Builders
	if opt.Enclose == Enclose {
		subject.WriteString(subject.Dialect().OpenEnclose())
	}
}

// closeHook concentra el proceso después del parseo
func closeHook(master Builder, subject Builder, circunstance interface{}, parameter interface{}, opt parsingOptions) {
	if opt.Enclose == Enclose {
		subject.WriteString(subject.Dialect().CloseEnclose())
	}

	if opt.Alias != "" {
		if opt.UseAS {
			subject.WriteString(fmt.Sprintf(" AS %v ", opt.Alias))
		} else {
			subject.WriteString(fmt.Sprintf(" %v ", opt.Alias))
		}
	}

	if opt.On != "" {
		subject.WriteString(fmt.Sprintf(" ON %v", opt.On))

		// este if es para agregar los posibles parametros en una clausula on
		if parameter != nil && parameter != "" {
			subject.AddParam(parameter)
		}
	}

	if opt.Operator != "" {
		subject.WriteString(fmt.Sprintf(" %v ", opt.Operator))
	}

	// // este if es para agregar los posibles parametros en una clausula WHERE
	// Cuidado! se debe mantener este orden para que funcione el hack de agregar la marca de parámetro junto con el paramètro.
	// justo después de agregar al operador
	if parameter != nil && (opt.On == "" && parameter != "") {
		subject.WriteString(subject.Dialect().ParamMark())
		subject.AddParam(parameter)
	}

	opt.AfterWork(subject)
}
