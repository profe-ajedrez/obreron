package dialect

// Dialect describes SQL syntax differences that matter to Obreron.
//
// Implementations MUST be immutable and concurrency-safe by construction.
//
// This interface is considered a contract (spec Rev 2.0).
type Dialect interface {
	// AppendPlaceholder appends the positional placeholder for argument n (0-based).
	// MySQL/SQLite: always '?' (ignores n)
	// PostgreSQL: '$1', '$2', ... (n+1)
	AppendPlaceholder(dst []byte, n int) []byte

	// AppendQuotedIdentifier appends an escaped/quoted SQL identifier.
	AppendQuotedIdentifier(dst []byte, id string) []byte

	// Name returns a short dialect name for errors/logging.
	Name() string

	// SupportsReturning indicates whether RETURNING is supported.
	SupportsReturning() bool

	// MaxParams returns a max positional-parameter limit (0 = unknown / do not validate).
	MaxParams() int
}
