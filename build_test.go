package obreron_test

import (
	"errors"
	"testing"

	"github.com/profe-ajedrez/obreron/v3"
	"github.com/profe-ajedrez/obreron/v3/dialect"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// buildInto — assembler behaviour via .AddFrag() (defined in export_test.go)
// ---------------------------------------------------------------------------

func TestBuild_Empty(t *testing.T) {
	sql, args, err := obreron.Select(dialect.MySQL{}).Build()
	mustNoErr(t, err)
	if sql != "" {
		t.Fatalf("expected empty SQL, got %q", sql)
	}
	if len(args) != 0 {
		t.Fatalf("expected no args, got %v", args)
	}
}

func TestBuild_SingleFragmentNoPlaceholders(t *testing.T) {
	sql, _, err := obreron.Select(dialect.MySQL{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT 1").
		Build()
	mustNoErr(t, err)
	if sql != "SELECT 1" {
		t.Fatalf("got %q want %q", sql, "SELECT 1")
	}
}

func TestBuild_MultiFragmentSpaceSeparated(t *testing.T) {
	sql, _, err := obreron.Select(dialect.MySQL{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT id, name").
		AddFrag("FROM", obreron.STypeTarget, "FROM users").
		Build()
	mustNoErr(t, err)
	want := "SELECT id, name FROM users"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
}

func TestBuild_SortByType_AddedOutOfOrder(t *testing.T) {
	// Fragments added in wrong clause order; assembler must stable-sort them.
	sql, _, err := obreron.Select(dialect.MySQL{}).
		AddFrag("FROM", obreron.STypeTarget, "FROM users").   // sType=30
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT id"). // sType=10
		Build()
	mustNoErr(t, err)
	want := "SELECT id FROM users"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
}

func TestBuild_StableSortPreservesInsertionOrder(t *testing.T) {
	// Two WHERE fragments share the same sType; insertion order must be kept.
	sql, args, err := obreron.Select(dialect.MySQL{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT id").
		AddFrag("FROM", obreron.STypeTarget, "FROM t").
		AddFrag("WHERE", obreron.STypeWhere, "WHERE a = ?", 1).
		AddFrag("AND", obreron.STypeWhere, "AND b = ?", 2). // same sType
		Build()
	mustNoErr(t, err)
	want := "SELECT id FROM t WHERE a = ? AND b = ?"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != 2 {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestBuild_MySQLPlaceholderExpansion(t *testing.T) {
	sql, args, err := obreron.Select(dialect.MySQL{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT id").
		AddFrag("WHERE", obreron.STypeWhere, "WHERE x = ? AND y = ?", 10, 20).
		Build()
	mustNoErr(t, err)
	want := "SELECT id WHERE x = ? AND y = ?"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 2 || args[0] != 10 || args[1] != 20 {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestBuild_PostgresPlaceholderExpansion(t *testing.T) {
	sql, args, err := obreron.Select(dialect.Postgres{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT id").
		AddFrag("WHERE", obreron.STypeWhere, "WHERE x = ? AND y = ?", "foo", "bar").
		Build()
	mustNoErr(t, err)
	want := "SELECT id WHERE x = $1 AND y = $2"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 2 || args[0] != "foo" || args[1] != "bar" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestBuild_ArgsReorderedAfterSegmentSort(t *testing.T) {
	// WHERE added first but assembled last (sType 50 > 10/30).
	// Arg must still appear at $1 in Postgres output.
	sql, args, err := obreron.Select(dialect.Postgres{}).
		AddFrag("WHERE", obreron.STypeWhere, "WHERE id = ?", 99).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT id").
		AddFrag("FROM", obreron.STypeTarget, "FROM users").
		Build()
	mustNoErr(t, err)
	want := "SELECT id FROM users WHERE id = $1"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 1 || args[0] != 99 {
		t.Fatalf("unexpected args: %v (want [99])", args)
	}
}

func TestBuild_MultiSegmentArgsReorderWithMultipleParams(t *testing.T) {
	// ORDER BY (80) added first; WHERE (50) last.
	// After sort: SELECT(10) → FROM(30) → WHERE(50) → ORDER BY(80).
	sql, args, err := obreron.Select(dialect.Postgres{}).
		AddFrag("ORDER BY", obreron.STypeOrderBy, "ORDER BY name ASC").
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT id, name").
		AddFrag("FROM", obreron.STypeTarget, "FROM products").
		AddFrag("WHERE", obreron.STypeWhere, "WHERE price > ? AND active = ?", 100, true).
		Build()
	mustNoErr(t, err)
	want := "SELECT id, name FROM products WHERE price > $1 AND active = $2 ORDER BY name ASC"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 2 || args[0] != 100 || args[1] != true {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestBuild_LiteralDoubleQuestionmark(t *testing.T) {
	// ?? → literal '?', not a placeholder; no args consumed.
	sql, args, err := obreron.Select(dialect.MySQL{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT '??'").
		Build()
	mustNoErr(t, err)
	want := "SELECT '?'"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 0 {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestBuild_AccumulatedErrorShortCircuits(t *testing.T) {
	sql, args, err := obreron.Select(dialect.MySQL{}).
		ForceErr("test", obreron.ErrEmptyFrom).
		Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, obreron.ErrEmptyFrom) {
		t.Fatalf("expected ErrEmptyFrom, got %v", err)
	}
	if sql != "" || len(args) != 0 {
		t.Fatalf("expected empty sql/args on error, got sql=%q args=%v", sql, args)
	}
}

// ---------------------------------------------------------------------------
// BuildInto — caller-owned buffers
// ---------------------------------------------------------------------------

func TestSelectStm_BuildInto_ReusesBuffers(t *testing.T) {
	buf := make([]byte, 0, 64)
	args := make([]any, 0, 4)
	buf, args, err := obreron.Select(dialect.MySQL{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT 1").
		BuildInto(buf, args)
	mustNoErr(t, err)
	if string(buf) != "SELECT 1" {
		t.Fatalf("got %q", string(buf))
	}
	if len(args) != 0 {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestBuildInto_AppendsToExistingBuffers(t *testing.T) {
	existingBuf := []byte("PREFIX:")
	existingArgs := []any{"pre"}
	buf, args, err := obreron.Select(dialect.MySQL{}).
		AddFrag("SELECT", obreron.STypeKeyword, "SELECT 1").
		BuildInto(existingBuf, existingArgs)
	mustNoErr(t, err)
	if string(buf) != "PREFIX:SELECT 1" {
		t.Fatalf("got %q", string(buf))
	}
	if len(args) != 1 || args[0] != "pre" {
		t.Fatalf("expected pre-existing arg preserved: %v", args)
	}
}

// ---------------------------------------------------------------------------
// MustBuild
// ---------------------------------------------------------------------------

func TestMustBuild_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from MustBuild, got none")
		}
	}()
	obreron.Select(dialect.MySQL{}).
		ForceErr("test", obreron.ErrEmptyFrom).
		MustBuild()
}

// ---------------------------------------------------------------------------
// Insert / Update / Delete builders
// ---------------------------------------------------------------------------

func TestInsertStatement_Build(t *testing.T) {
	sql, args, err := obreron.Insert(dialect.Postgres{}).
		AddFrag("INSERT", obreron.STypeKeyword, "INSERT INTO users (name)").
		AddFrag("VALUES", obreron.STypeTarget, "VALUES (?)", "alice").
		Build()
	mustNoErr(t, err)
	want := "INSERT INTO users (name) VALUES ($1)"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 1 || args[0] != "alice" {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestUpdateStm_Build(t *testing.T) {
	sql, args, err := obreron.Update(dialect.MySQL{}).
		AddFrag("UPDATE", obreron.STypeKeyword, "UPDATE users").
		AddFrag("SET", obreron.STypeSet, "SET name = ?", "bob").
		AddFrag("WHERE", obreron.STypeWhere, "WHERE id = ?", 7).
		Build()
	mustNoErr(t, err)
	want := "UPDATE users SET name = ? WHERE id = ?"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 2 || args[0] != "bob" || args[1] != 7 {
		t.Fatalf("unexpected args: %v", args)
	}
}

func TestDeleteStm_Build(t *testing.T) {
	sql, args, err := obreron.Delete(dialect.Postgres{}).
		AddFrag("DELETE", obreron.STypeKeyword, "DELETE FROM users").
		AddFrag("WHERE", obreron.STypeWhere, "WHERE id = ?", 42).
		Build()
	mustNoErr(t, err)
	want := "DELETE FROM users WHERE id = $1"
	if sql != want {
		t.Fatalf("got %q want %q", sql, want)
	}
	if len(args) != 1 || args[0] != 42 {
		t.Fatalf("unexpected args: %v", args)
	}
}

// ---------------------------------------------------------------------------
// DB factory
// ---------------------------------------------------------------------------

func TestDB_FactoryCreatesBuilders(t *testing.T) {
	db := obreron.New(dialect.MySQL{})
	if db.Select() == nil {
		t.Fatal("Select() returned nil")
	}
	if db.Insert() == nil {
		t.Fatal("Insert() returned nil")
	}
	if db.Update() == nil {
		t.Fatal("Update() returned nil")
	}
	if db.Delete() == nil {
		t.Fatal("Delete() returned nil")
	}
}

func TestPackageLevel_ConstructorFunctions(t *testing.T) {
	d := dialect.Postgres{}
	if obreron.Select(d) == nil {
		t.Fatal("Select() returned nil")
	}
	if obreron.Insert(d) == nil {
		t.Fatal("Insert() returned nil")
	}
	if obreron.Update(d) == nil {
		t.Fatal("Update() returned nil")
	}
	if obreron.Delete(d) == nil {
		t.Fatal("Delete() returned nil")
	}
}
