package obreron

// sType constants are the ordering key for SQL clause segments.
//
// Within a single stament, segments are stable-sorted by sType before
// assembly. Insertion order is preserved for segments with equal sType,
// so multiple WHERE conditions added in order stay in order.
//
// Each stament belongs to exactly one statement type, so numeric values
// may be reused across statement types (SELECT vs INSERT vs UPDATE vs
// DELETE are never mixed in the same stament).
//
// Gaps between values are intentional: they allow future clauses to be
// inserted without renumbering the entire table.
const (
	// sTypeKeyword is the primary keyword fragment.
	// SELECT / INSERT INTO table / UPDATE table / DELETE FROM table.
	sTypeKeyword uint8 = 10

	// sTypeColumnList is the projected column list.
	// SELECT: "col1, col2, …".
	// INSERT: "(col1, col2, …)".
	sTypeColumnList uint8 = 20

	// sTypeTarget covers the target-table or values clause.
	// SELECT / DELETE: "FROM <table>"
	// INSERT:          "VALUES (?, …)".
	sTypeTarget uint8 = 30

	// sTypeJoin covers all JOIN variants in SELECT.
	// sTypeSet covers SET assignments in UPDATE.
	// Both occupy the same slot (40) because they never appear together.
	sTypeJoin uint8 = 40
	sTypeSet  uint8 = 40

	// sTypeWhere covers WHERE and any AND/OR continuations.
	sTypeWhere uint8 = 50

	// sTypeGroupBy covers GROUP BY.
	sTypeGroupBy uint8 = 60

	// sTypeHaving covers HAVING.
	sTypeHaving uint8 = 70

	// sTypeOrderBy covers ORDER BY.
	sTypeOrderBy uint8 = 80

	// sTypeLimit covers LIMIT (stored as a literal, no placeholder).
	sTypeLimit uint8 = 90

	// sTypeOffset covers OFFSET (stored as a literal, no placeholder).
	sTypeOffset uint8 = 95

	// sTypeReturning covers RETURNING, gated by SupportsReturning().
	sTypeReturning uint8 = 200
)

// Structural flags stored in stament.flags.
// Each builder sets the relevant bit when a required clause is added.
// Build() reads these flags to enforce completeness invariants.
const (
	flagHasFrom   uint8 = 1 << 0 // FROM / target table was set
	flagHasSet    uint8 = 1 << 1 // At least one SET assignment (UPDATE)
	flagHasCols   uint8 = 1 << 2 // At least one column declared (INSERT / SELECT explicit cols)
	flagHasValues uint8 = 1 << 3 // VALUES clause was added (INSERT)
)
