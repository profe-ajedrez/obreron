package dialect

// Dialect describes SQL syntax differences that matter to Obreron.
//
// Implementations MUST be immutable and concurrency-safe by construction.
//
// This interface is considered a contract (spec Rev 2.0).
type Dialect interface {
	// AppendPlaceholder escribe el marcador posicional para el argumento n (base-0).
	// MySQL/SQLite: escribe '?'. PostgreSQL: escribe '$1', '$2', etc.
	// No aloca si dst tiene capacidad suficiente.
	AppendPlaceholder(dst []byte, n int) []byte

	// AppendQuotedIdentifier escribe el identificador con el escape del dialecto.
	// MySQL: `id`. PostgreSQL/SQLite: "id".
	// Escapa internamente: " → "" (Postgres/SQLite), ` → `` (MySQL).
	// El caller debe validar que id no está vacío antes de invocar.
	AppendQuotedIdentifier(dst []byte, id string) []byte

	// Name retorna el nombre del dialecto para logs y mensajes de error.
	Name() string

	// SupportsReturning indica si el dialecto soporta la cláusula RETURNING.
	SupportsReturning() bool

	// MaxParams retorna el límite de parámetros posicionales del driver.
	// 0 significa desconocido: el builder no valida el límite.
	// MySQL: 65535. PostgreSQL: 65535. SQLite: 999 (default de compilación).
	MaxParams() int
}
