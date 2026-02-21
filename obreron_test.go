package obreron_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/pingcap/tidb/pkg/parser"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/profe-ajedrez/obreron/v2"
)

func TestUpdateBuildAfterCloseReturnsSnapshot(t *testing.T) {
	makeBad := func() *obreron.UpdateStm {
		ob := obreron.Select().
			Col("1").
			From("t").
			Where("1=1")
		defer obreron.CloseSelect(ob)

		up := obreron.Update("vw_docs_search v").
			ColSelect(ob, "det").
			Set("v.x = det.x").
			Where("v.id = det.id")

		defer obreron.CloseUpdate(up) // patrón malo
		return up
	}

	up := makeBad()

	// Churn: intenta forzar reuse del pool
	s := obreron.Select().Col("1").From("x")
	_, _ = s.Build()
	s.Close()

	q, _ := up.Build()
	if q == "" || !strings.HasPrefix(q, "UPDATE ") {
		t.Fatalf("expected stable UPDATE snapshot, got: %q", q)
	}
}

func TestUseAfterCloseSequential(t *testing.T) {
	// helper que retorna builder ya cerrado (patrón real a buscar)
	bad := func() *obreron.UpdateStm {
		up := obreron.Update("t").Set("a=1")
		defer obreron.CloseUpdate(up)

		return up
	}

	runtime.GOMAXPROCS(1)

	up := bad()

	// churn del pool: fuerza a que el mismo *stament sea reutilizado
	for i := 0; i < 5000; i++ {
		s := obreron.Select().Col("1").From("x")
		_, _ = s.Build()
		s.Close()
	}

	q, _ := up.Build()
	if q == "" {
		t.Fatalf("corrupted SQL: empty string (builder likely closed/recycled before Build)")
	}

	if !strings.HasPrefix(q, "UPDATE ") {
		t.Fatalf("corrupted SQL: %s", q)
	}
}

func TestCorruptionStress(t *testing.T) {
	t.Parallel()

	const goroutines = 32

	const iters = 2000

	errCh := make(chan error, goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			for i := 0; i < iters; i++ {
				// construye el update con subselect, como tu caso real
				ob := obreron.Select().
					Col("1").
					From("t").
					Where("1=1")
				up := obreron.Update("vw_docs_search v").
					ColSelect(ob, "det").
					Set("v.x = 1").
					Where("v.id = det.id")

				q, _ := up.Build()

				// quick invariant: un UPDATE debe empezar con UPDATE
				if !strings.HasPrefix(q, "UPDATE") {
					errCh <- fmt.Errorf("corrupt SQL: %s", q)
					return
				}

				obreron.CloseSelect(ob)
				obreron.CloseUpdate(up)
			}

			errCh <- nil
		}()
	}

	for i := 0; i < goroutines; i++ {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}
}

// TestInsertH3Regression verifica que obreron.InsertStament.Col genera placeholders
// correctos (con comas) cuando se pasan múltiples params en posiciones no-primeras.
//
// Cómo correr solo este test:
//
//	go test -run TestInsertH3Regression -v
//
// El test DEBE FALLAR antes de aplicar el fix en insert.go.
// El test DEBE PASAR después del fix.
func TestInsertH3Regression(t *testing.T) {
	for i, tc := range insertH3RegressionCases() {
		t.Run(tc.name, func(t *testing.T) {
			sql, p := tc.tc.Build()
			defer tc.tc.Close()

			if sql != tc.expected {
				t.Errorf(
					"[caso %d] SQL incorrecto\n  want: %s\n   got: %s",
					i, tc.expected, sql,
				)
			}

			if len(p) != len(tc.expectedParams) {
				t.Errorf(
					"[caso %d] longitud de params incorrecta: want %d got %d\n  params: %v",
					i, len(tc.expectedParams), len(p), p,
				)

				return // sin esto el siguiente loop puede panic
			}

			for k := range tc.expectedParams {
				if p[k] != tc.expectedParams[k] {
					t.Errorf(
						"[caso %d] params[%d]: want %v got %v",
						i, k, tc.expectedParams[k], p[k],
					)
				}
			}
		})
	}
}

// SELECT a1, a2, ? AS diez, colIf1, colIf2, ? AS zero, a3, ? AS cien FROM client c JOIN addresses a ON a.id_cliente = a.id_cliente JOIN phones p ON p.id_cliente = c.id_cliente JOIN mailes m ON m.id_cliente = m.id_cliente AND c.estado_cliente = ? LEFT JOIN left_joined lj ON lj.a1 = c.a1 WHERE a1 = ? AND a2 = ? AND a3 = 10 AND a16 = ? --- Got
// SELECT a1, a2, ? AS diez, colIf1, colIf2, ? AS zero, a3, ? AS cien FROM client c LEFT JOIN left_joined lj ON lj.a1 = c.a1 JOIN addresses a ON a.id_cliente = a.id_cliente JOIN phones p ON p.id_cliente = c.id_cliente JOIN mailes m ON m.id_cliente = m.id_cliente AND c.estado_cliente = ? WHERE a1 = ? AND a2 = ? AND a3 = 10 AND a16 = ?

func TestSelect(t *testing.T) {
	for i, tc := range selectTestCases() {
		sql, p := tc.tc.Build()

		if sql != tc.expected {
			t.Logf("[Test case %d %s] Failed! Expected %s --- Got %s", i, tc.name, tc.expected, sql)
			t.FailNow()
		}

		if len(p) != len(tc.expectedParams) {
			t.Logf("[Test case %d %s] Failed! Params Length Expected %d --- Got %d", i, tc.name, len(tc.expectedParams), len(p))
			t.FailNow()
		}

		for k := range tc.expectedParams {
			if p[k] != tc.expectedParams[k] {
				t.Logf("[Test case %d %s] Failed! Param[%d] Expected %v --- Got p[%d] = %v", i, tc.name, k, tc.expectedParams[k], k, p[k])
				t.FailNow()
			}
		}

		tc.tc.Close()
	}
}

func BenchmarkSelect(b *testing.B) {
	for _, tc := range selectTestCases() {
		b.ResetTimer()
		b.Run(tc.name, func(b2 *testing.B) {
			for i := 0; i < b2.N; i++ {
				_, _ = tc.tc.Build()
				tc.tc.Close()
			}
		})
	}
}

func TestDelete(t *testing.T) {
	for i, tc := range deleteTestCases() {
		sql, p := tc.tc.Build()

		if sql != tc.expected {
			t.Logf("[Test case %d %s] Failed! Expected %s --- Got %s", i, tc.name, tc.expected, sql)
			t.FailNow()
		}

		if len(p) != len(tc.expectedParams) {
			t.Logf("[Test case %d %s] Failed! Params Length Expected %d --- Got %d", i, tc.name, len(tc.expectedParams), len(p))
			t.FailNow()
		}

		for k := range tc.expectedParams {
			if p[k] != tc.expectedParams[k] {
				t.Logf("[Test case %d %s] Failed! Param[%d] Expected %v --- Got p[%d] = %v", i, tc.name, k, tc.expectedParams[k], k, p[k])
				t.FailNow()
			}
		}

		tc.tc.Close()
	}
}

func BenchmarkDelete(b *testing.B) {
	for _, tc := range deleteTestCases() {
		b.ResetTimer()
		b.Run(tc.name, func(b2 *testing.B) {
			for i := 0; i < b2.N; i++ {
				_, _ = tc.tc.Build()
				tc.tc.Close()
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	par := parser.New()

	for round := 0; round <= 1000; round++ {
		for i, tc := range updateTestCases() {
			upStm := tc.tc()
			sql, p := upStm.Build()

			// Parser para verificar que el sql es correctamente construido
			_, _, err := par.ParseSQL(sql)
			if err != nil {
				t.Logf("[TEST CASE %d ROUND %d %s] %v", i, round, tc.name, err)
				t.FailNow()
			}

			if tc.expected != "" {
				if sql != tc.expected {
					t.Logf("[Test case %d ROUND %d %s] Failed! Expected %s --- Got %s", i, round, tc.name, tc.expected, sql)
					t.FailNow()
				}

				if len(p) != len(tc.expectedParams) {
					t.Logf("[Test case %d ROUND %d %s] Failed! Params Length Expected %d --- Got %d", i, round, tc.name, len(tc.expectedParams), len(p))
					t.FailNow()
				}

				for pi := range tc.expectedParams {
					if p[pi] != tc.expectedParams[pi] {
						t.Logf("[Test case %d ROUND %d %s] Failed! Param[%d] Expected %v --- Got p[%d] = %v",
							i, round, tc.name, pi, tc.expectedParams[pi], pi, p[pi])
						t.FailNow()
					}
				}

				obreron.CloseUpdate(upStm)
			}
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	for _, tc := range updateTestCases() {
		b.ResetTimer()
		b.Run(tc.name, func(b2 *testing.B) {
			for i := 0; i < b2.N; i++ {
				upStm := tc.tc()
				_, _ = upStm.Build()
				obreron.CloseUpdate(upStm)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	for i, tc := range insertTestCases() {
		sql, p := tc.tc.Build()

		if sql != tc.expected {
			t.Logf("[Test case %d %s] Failed! Expected %s --- Got %s", i, tc.name, tc.expected, sql)
			t.FailNow()
		}

		if len(p) != len(tc.expectedParams) {
			t.Logf("[Test case %d %s] Failed! Params Length Expected %d --- Got %d", i, tc.name, len(tc.expectedParams), len(p))
			t.FailNow()
		}

		for k := range tc.expectedParams {
			if p[k] != tc.expectedParams[k] {
				t.Logf("[Test case %d %s] Failed! Param[%d] Expected %v --- Got p[%d] = %v", i, tc.name, k, tc.expectedParams[k], k, p[k])
				t.FailNow()
			}
		}

		tc.tc.Close()
	}
}

func BenchmarkInsert(b *testing.B) {
	for _, tc := range insertTestCases() {
		b.ResetTimer()
		b.Run(tc.name, func(b2 *testing.B) {
			for i := 0; i < b2.N; i++ {
				_, _ = tc.tc.Build()
				tc.tc.Close()
			}
		})
	}
}

func selectTestCases() (cases []struct {
	tc             *obreron.SelectStm
	name           string
	expected       string
	expectedParams []any
}) {
	cases = []struct {
		tc             *obreron.SelectStm
		name           string
		expected       string
		expectedParams []any
	}{
		/* SELECT TESTS */
		{
			name:           "columns - from",
			expected:       "SELECT a1, a2, a3 FROM client",
			expectedParams: nil,
			tc:             obreron.Select().Col("a1, a2, a3").From("client"),
		},
		{
			name:           "columns - from",
			expected:       "SELECT a1, a2, a3 FROM client",
			expectedParams: nil,
			tc:             obreron.Select().Col("a1, a2, a3").From("client"),
		},
		{
			name:           "columns params - from",
			expected:       "SELECT a1, a2, ? AS cien FROM client",
			expectedParams: []any{100},
			tc:             obreron.Select().Col("a1, a2, ? AS cien", 100).From("client"),
		},
		{
			name:           "columns params - columns params - from",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client",
			expectedParams: []any{10, 100},
			tc:             obreron.Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).From("client"),
		},
		{
			name:           "columns params - columns params - Col If - from",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client",
			expectedParams: []any{10, 100},
			tc:             obreron.Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client"),
		},
		{
			name:           "columns params - columns params - Col If - from - where",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client WHERE 1 = 1",
			expectedParams: []any{10, 100},
			tc:             obreron.Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client").Where("1 = 1"),
		},
		{
			name:           "columns params - Col If - columns params - Col If - from - join if -where",
			expected:       "SELECT SQL NO CACHE b2 AS b2, a1, a2, ? AS diez, a3, ? AS cien FROM client JOIN tablada t ON t.id=client.id LEFT JOIN tabla_b tb ON tb.id = client.id LEFT JOIN tabla_c tc ON tc.id = client.id RIGHT JOIN tabla_e te ON te.id = client.id AND 2 = 2 RIGHT JOIN tabla_f tf ON tf.id = client.id OUTER JOIN tabla_h th ON th.id = client.id AND 2 = 2 OUTER JOIN tabla_i ti ON ti.id = client.id WHERE 1 = 1 AND 3 = 3 AND colX LIKE '%chamullo%' OR 4 = 4 OR 5 = 5 GROUP BY colY , colZ HAVING colZ = 'abc' ORDER BY 1 ASC LIMIT ? OFFSET ?",
			expectedParams: []any{10, 100, 100, 200},
			tc: obreron.Select().Clause("SQL NO", "").ClauseIf(true, "CACHE", "").
				ColIf(true, "b2 AS b2").
				Col("a1, a2, ? AS diez", 10).
				Col("a3, ? AS cien", 100).
				ColIf(false, "a4").
				From("client").
				JoinIf(true, "tablada t ON t.id=client.id").
				LeftJoin("tabla_b tb ON tb.id = client.id").
				LeftJoinIf(true, "tabla_c tc").OnIf(true, "tc.id = client.id").
				LeftJoinIf(false, "tabla_d td ON td.id = client.id").
				RightJoin("tabla_e te ON te.id = client.id").AndIf(true, "2 = 2").
				RightJoinIf(true, "tabla_f tf ON tf.id = client.id").
				RightJoinIf(false, "tabla_g tg ON tg.id = client.id").
				OuterJoin("tabla_h th ON th.id = client.id").AndIf(true, "2 = 2").
				OuterJoinIf(true, "tabla_i ti ON ti.id = client.id").
				OuterJoinIf(false, "tabla_j tg ON tj.id = client.id").
				Where("1 = 1").AndIf(true, "3 = 3").And("colX").LikeIf(true, "'%chamullo%'").Or("4 = 4").OrIf(true, "5 = 5").
				GroupBy("colY").GroupBy("colZ").Having("colZ = 'abc'").
				OrderBy("1 ASC").
				Limit(100).
				Offset(200),
		},
		{
			name:           "columns params - columns params - Col If - from - where",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client WHERE a1 = ?",
			expectedParams: []any{10, 100, "'last name'"},
			tc:             obreron.Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client").Where("a1 = ?", "'last name'"),
		},
		{
			name:           "columns params - columns params - Col If - from - where",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client WHERE a1 = ? AND a2 = ?",
			expectedParams: []any{10, 100, "'last name'", 1000.54},
			tc:             obreron.Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client").Where("a1 = ?", "'last name'").And("a2 = ?", 1000.54),
		},
		{
			name:           "columns params - columns params - Col If - from - where shuffled",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client WHERE a1 = ? AND a2 = ?",
			expectedParams: []any{10, 100, "'last name'", 1000.54},
			tc: obreron.Select().
				Where("a1 = ?", "'last name'").
				And("a2 = ?", 1000.54).
				Col("a1, a2, ? AS diez", 10).
				From("client").
				Col("a3, ? AS cien", 100).
				ColIf(false, "a4"),
		},
		{
			name:           "complex query shuffled",
			expected:       `SELECT a1, a2, ? AS diez, colIf1, colIf2, ? AS zero, a3, ? AS cien FROM client c JOIN addresses a ON a.id_cliente = a.id_cliente JOIN phones p ON p.id_cliente = c.id_cliente JOIN mailes m ON m.id_cliente = m.id_cliente AND c.estado_cliente = ? LEFT JOIN left_joined lj ON lj.a1 = c.a1 WHERE a1 = ? AND a2 = ? AND a3 = 10 AND a16 = ?`,
			expectedParams: []any{10, 0, 100, 0, "'last name'", 1000.54, 75},
			tc: obreron.Select().
				Where("a1 = ?", "'last name'").
				Col("a1, a2, ? AS diez", 10).
				ColIf(true, `colIf1, colIf2, ? AS zero`, 0).
				Col("a3, ? AS cien", 100).
				ColIf(false, "a4").
				Where("a2 = ?", 1000.54).
				And("a3 = 10").And("a16 = ?", 75).
				AndIf(false, "will_not_be_shown = ?", 10).
				Join("addresses a ON a.id_cliente = a.id_cliente").
				Join("phones p").On("p.id_cliente = c.id_cliente").
				Join("mailes m").On("m.id_cliente = m.id_cliente").
				And("c.estado_cliente = ?", 0).
				JoinIf(false, "not_to_join ntj").OnIf(false, "ntj.will_not_be_shown = c.will_not_be_shown").
				LeftJoin("left_joined lj").On("lj.a1 = c.a1").
				From("client c"),
		},
		{
			name:           "complex query badly shuffled",
			expected:       `SELECT a1, a2, ? AS diez, colIf1, colIf2, ? AS zero, a3, ? AS cien FROM client c JOIN addresses a ON a.id_cliente = a.id_cliente JOIN phones p ON p.id_cliente = c.id_cliente JOIN mailes m ON m.id_cliente = m.id_cliente AND c.estado_cliente = ? LEFT JOIN left_joined lj ON lj.a1 = c.a1 WHERE a1 = ? AND a2 = ? AND a3 = 10 AND a16 = ?`,
			expectedParams: []any{10, 0, 100, 0, "'last name'", 1000.54, 75},
			tc: obreron.Select().
				Where("a1 = ?", "'last name'").
				Join("addresses a ON a.id_cliente = a.id_cliente").
				Join("phones p").On("p.id_cliente = c.id_cliente").
				Col("a1, a2, ? AS diez", 10).
				ColIf(true, `colIf1, colIf2, ? AS zero`, 0).
				Where("a2 = ?", 1000.54).
				And("a3 = 10").And("a16 = ?", 75).
				AndIf(false, "will_not_be_shown = ?", 10).
				Col("a3, ? AS cien", 100).
				ColIf(false, "a4").
				Join("mailes m").On("m.id_cliente = m.id_cliente").
				And("c.estado_cliente = ?", 0).
				JoinIf(false, "not_to_join ntj").OnIf(false, "ntj.will_not_be_shown = c.will_not_be_shown").
				LeftJoin("left_joined lj").On("lj.a1 = c.a1").
				From("client c"),
		},
		{
			name:           "columns - where in",
			expected:       "SELECT a1, a2, a3 FROM client WHERE status IN ( 0, 1, 2, 3)",
			expectedParams: nil,
			tc:             obreron.Select().Col("a1, a2, a3").From("client").Where("status").In("0, 1, 2, 3"),
		},
		{
			name:           "columns - where in",
			expected:       "SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND status IN ( 0, 1, 2, 3)",
			expectedParams: nil,
			tc:             obreron.Select().Col("a1, a2, a3").From("client").Where("1 = 1").And("status").In("0, 1, 2, 3"),
		},
		{
			name:           "columns - where in",
			expected:       "SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND status IN (?, ?, ?, ?)",
			expectedParams: []any{0, 1, 2, 3},
			tc:             obreron.Select().Col("a1, a2, a3").From("client").Where("1 = 1").Y().InArgs("status", 0, 1, 2, 3),
		},
		{
			name:           "columns - where like",
			expected:       "SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND city LIKE '%ago%'",
			expectedParams: nil,
			tc:             obreron.Select().Col("a1, a2, a3").From("client").Where("1 = 1").And("city").Like("'%ago%'"),
		},
		{
			name:           "columns - where in (InArgs empty)",
			expected:       "SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND status IN (NULL)",
			expectedParams: nil,
			tc:             obreron.Select().Col("a1, a2, a3").From("client").Where("1 = 1").Y().InArgs("status"),
		},
	}

	return cases
}

func deleteTestCases() []struct {
	tc             *obreron.DeleteStm
	name           string
	expected       string
	expectedParams []any
} {
	return []struct {
		tc             *obreron.DeleteStm
		name           string
		expected       string
		expectedParams []any
	}{
		/* DELETE TESTS */
		{
			name:           "simple del",
			expected:       "DELETE FROM client",
			expectedParams: nil,
			tc:             obreron.Delete().From("client"),
		},
		{
			name:           "simple del where",
			expected:       "DELETE FROM client WHERE client_id = 100 AND b = 3 OR 2 = 2 OR 3 = 3",
			expectedParams: nil,
			tc:             obreron.Delete().From("client").Where("client_id = 100").AndIf(true, "b = 3").OrIf(true, "2 = 2").Or("3 = 3"),
		},
		{
			name:           "del where conditions",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) AND a LIKE '%ago%' -- Comment\n",
			expectedParams: []any{"'G03'"},
			tc: obreron.Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "'G03'").And("a").LikeIf(true, "'%ago%'").ClauseIf(true, "-- Comment\n", ""),
		},
		{
			name:           "del where conditions",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN (?, ?, ?) AND a LIKE '%ago%' -- Comment\n",
			expectedParams: []any{"G01", "G02", "G03"},
			tc: obreron.Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().InArgs("regime_cliente", "G01", "G02", "G03").And("a").LikeIf(true, "'%ago%'").ClauseIf(true, "-- Comment\n", ""),
		},
		{
			name:           "del where conditions limit",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) LIMIT ?",
			expectedParams: []any{"'G03'", 100},
			tc: obreron.Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "'G03'").
				Limit(100),
		},
		{
			name:           "del where conditions limit -- shuffled",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) LIMIT ?",
			expectedParams: []any{"G03", 100},
			tc: obreron.Delete().From("client").
				Limit(100).
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "G03"),
		},
		{
			name:           "simple del where quick",
			expected:       "DELETE QUICK FROM client WHERE client_id = 100",
			expectedParams: nil,
			tc:             obreron.Delete().Clause("QUICK", "").From("client").Where("client_id = 100"),
		},
		{
			name:           "simple del where ignore",
			expected:       "DELETE IGNORE FROM client WHERE client_id = 100",
			expectedParams: nil,
			tc:             obreron.Delete().Clause("IGNORE", "").From("client").Where("client_id = 100"),
		},
		{
			name:           "simple del where partition",
			expected:       "DELETE PARTITION the_partition FROM client WHERE client_id = 100",
			expectedParams: nil,
			tc:             obreron.Delete().Clause("PARTITION", "the_partition").From("client").Where("client_id = 100"),
		},
		{
			name:           "simple del where order by limit",
			expected:       "DELETE FROM client WHERE client_id = 100 ORDER BY ciudad LIMIT ?",
			expectedParams: []any{10},
			tc:             obreron.Delete().From("client").Where("client_id = 100").OrderBy("ciudad").Limit(10),
		},
		{
			name:           "simple del where limit order by -- shuffled",
			expected:       "DELETE FROM client WHERE client_id = 100 ORDER BY ciudad LIMIT ?",
			expectedParams: []any{10},
			tc:             obreron.Delete().From("client").Where("client_id = 100").Limit(10).OrderBy("ciudad"),
		},
	}
}

func updateTestCases() (tcs []struct {
	tc             func() *obreron.UpdateStm
	name           string
	expected       string
	expectedParams []any
}) {
	tcs = append(tcs, []struct {
		tc             func() *obreron.UpdateStm
		name           string
		expected       string
		expectedParams []any
	}{
		{
			name:           "UPDATE client SET status = 0",
			expected:       "UPDATE client SET status = 0",
			expectedParams: nil,
			tc:             func() *obreron.UpdateStm { return obreron.Update("client").Set("status = 0") },
		},
		{
			name:           "UPDATE client SET status = 0, name = ?",
			expected:       "UPDATE client SET status = 0, name = ?",
			expectedParams: []any{"stitch"},
			tc:             func() *obreron.UpdateStm { return obreron.Update("client").Set("status = 0").Set("name = ?", "stitch") },
		},
		{
			name:           "UPDATE client SET status = 1 WHERE country = ? AND status IN (?, ?, ?, ?)",
			expected:       "UPDATE client SET status = 1 WHERE country = ? AND status IN (?, ?, ?, ?)",
			expectedParams: []any{"CL", 1, 2, 3, 4},
			tc: func() *obreron.UpdateStm {
				return obreron.Update("client").Set("status = 1").Where("country = ?", "CL").Y().In("status", "?, ?, ?, ?", 1, 2, 3, 4)
			},
		},
		{
			name:           "UPDATE client SET status = 2 WHERE country = ? AND status IN (?, ?, ?, ?) InArgs",
			expected:       "UPDATE client SET status = 2 WHERE country = ? AND status IN (?, ?, ?, ?)",
			expectedParams: []any{"CL", 1, 2, 3, 4},
			tc: func() *obreron.UpdateStm {
				return obreron.Update("client").Set("status = 2").Where("country = ?", "CL").Y().InArgs("status", 1, 2, 3, 4)
			},
		},
		{
			name:           "update with IN with 1 param in args",
			expected:       "UPDATE client SET status = 0 WHERE country = ? AND status IN (?)",
			expectedParams: []any{"CL", 1},
			tc: func() *obreron.UpdateStm {
				return obreron.Update("client").Set("status = 0").Where("country = ?", "CL").Y().InArgs("status", 1)
			},
		},
		{
			name:           "update where",
			expected:       "UPDATE client SET status = 0 WHERE status = ?",
			expectedParams: []any{1},
			tc:             func() *obreron.UpdateStm { return obreron.Update("client").Set("status = 0").Where("status = ?", 1) },
		},
		{
			name:           "update where order limit",
			expected:       "UPDATE client SET status = 0 WHERE status = ? ORDER BY ciudad LIMIT ?",
			expectedParams: []any{1, 10},
			tc: func() *obreron.UpdateStm {
				return obreron.Update("client").Set("status = 0").Where("status = ?", 1).OrderBy("ciudad").Limit(10)
			},
		},
		{
			name:           "UPDATE client SET status = 0 WHERE status = ? AND country = ? ORDER BY ciudad LIMIT ?",
			expected:       "UPDATE client SET status = 0 WHERE status = ? AND country = ? ORDER BY ciudad LIMIT ?",
			expectedParams: []any{1, "'CL'", 10},
			tc: func() *obreron.UpdateStm {
				return obreron.Update("client").Set("status = 0").Where("status = ?", 1).And("country = ?", "'CL'").OrderBy("ciudad").Limit(10)
			},
		},
		// UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9, a = 2, c = 3 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id AND regime_cliente IN ('G01','G02', ?) AND 2 = 2 OR 3 = 3 OR 4 = 4 AND colX LIKE '%ago%' AND colN LIKE '%oga%' AND colY (1, 2, 3) --- Got
		// UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9, a = 2, c = 3 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id AND regime_cliente IN ('G01','G02', ?) AND 2 = 2 OR 3 = 3 OR 4 = 4 AND colX LIKE '%ago%' AND colN LIKE '%oga%' AND colY LIKE (1, 2, 3)
		{
			name:           "update select",
			expected:       "UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9, a = 2, c = 3 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id AND regime_cliente IN ('G01','G02', ?) AND 2 = 2 OR 3 = 3 OR 4 = 4 AND colX LIKE '%ago%' AND colN LIKE '%oga%' AND colY IN (1, 2, 3)",
			expectedParams: []any{"'G03'"},
			tc: func() *obreron.UpdateStm {
				s := obreron.Select().Col("id, retail / wholesale AS markup, quantity").From("items")
				defer obreron.CloseSelect(s)

				return obreron.Update("items").
					ColSelectIf(true, s, "discounted").
					Set("items.retail = items.retail * 0.9").Set("a = 2").SetIf(true, "c = 3").
					Where("discounted.markup >= 1.3").
					And("discounted.quantity < 100").
					And("items.id = discounted.id").Y().In("regime_cliente", "'G01','G02', ?", "'G03'").AndIf(true, "2 = 2").Or("3 = 3").OrIf(true, "4 = 4").And("colX").Like("'%ago%'").AndIf(true, "colN").LikeIf(true, "'%oga%'").Y().In("colY", "1, 2, 3")
			},
		},
		// UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id --- Got
		// UPDATE items ,( SELECT , id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id
		{
			name:           "update join",
			expected:       "UPDATE business AS b JOIN business_geocode AS g ON b.business_id = g.business_id SET b.mapx = g.latitude, b.mapy = g.longitude WHERE (b.mapx = '' or b.mapx = 0) AND g.latitude > 0 AND 3 = 3",
			expectedParams: nil,
			tc: func() *obreron.UpdateStm {
				return obreron.Update("business AS b").
					JoinIf(true, "business_geocode AS g").OnIf(true, "b.business_id = g.business_id").
					Set("b.mapx = g.latitude, b.mapy = g.longitude").
					Where("(b.mapx = '' or b.mapx = 0)").And("g.latitude > 0").ClauseIf(true, "AND", "3 = 3")
			},
		},
		{
			tc: func() *obreron.UpdateStm {
				return obreron.Update("items").
					Set("items.retail = items.retail * 0.9").
					Set("a = 2").
					Where("discounted.markup >= 1.3").
					And("colX").
					Like("'%ago%'")
			},
			name:           "",
			expected:       "UPDATE items SET items.retail = items.retail * 0.9, a = 2 WHERE discounted.markup >= 1.3 AND colX LIKE '%ago%'",
			expectedParams: []any{},
		},
		{
			tc:             specialCaseUpdateTargetQuery,
			name:           "",
			expected:       "",
			expectedParams: []any{},
		},
		{
			tc:             specialCaseUpdateResendQuery,
			name:           "",
			expected:       "",
			expectedParams: []any{},
		},
		{
			tc:             specialCaseUpdateShippingQuery,
			name:           "",
			expected:       "",
			expectedParams: []any{},
		},
		{
			name:           "UPDATE client SET status = 0 WHERE country = ? AND status IN (NULL) InArgs empty",
			expected:       "UPDATE client SET status = 0 WHERE country = ? AND status IN (NULL)",
			expectedParams: []any{"CL"},
			tc: func() *obreron.UpdateStm {
				return obreron.Update("client").Set("status = 0").Where("country = ?", "CL").Y().InArgs("status")
			},
		},
	}...)

	return tcs
}

func insertTestCases() (tcs []struct {
	tc             *obreron.InsertStament
	name           string
	expected       string
	expectedParams []any
}) {
	tcs = append(tcs, []struct {
		tc             *obreron.InsertStament
		name           string
		expected       string
		expectedParams []any
	}{
		{
			name:           "simple insert",
			expected:       "INSERT IGNORE INTO client ( name, value ) VALUES ( ?,? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'"},
			tc: obreron.Insert().Ignore().Into("client").
				Col("name, value", "'some name'", "'somemail@mail.net'"),
		},
		{
			name:           "simple insert params",
			expected:       "INSERT INTO client ( name, value, data ) VALUES ( ?,?,? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'", "'some data'"},
			tc:             obreron.Insert().Into("client").Col("name", "'some name'").Col("value", "'somemail@mail.net'").ColIf(true, "data", "'some data'").ColIf(false, "info", 12),
		},
		{
			name:           "simple insert params shuffled",
			expected:       "INSERT INTO client ( name, value ) VALUES ( ?,? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'"},
			tc:             obreron.Insert().Col("name, value", "'some name'", "'somemail@mail.net'").Into("client"),
		},
		{
			name:           "simple insert params select",
			expected:       "INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2",
			expectedParams: nil,
			tc: obreron.Insert().
				Into("courses").
				ColSelectIf(true, "name, location, gid", obreron.Select().Col("name, location, 1").From("courses").Where("cid = 2")).
				ColSelectIf(false, "last_name, last_location, grid", obreron.Select().Col("last_name, last_location, 11").From("courses").Where("cid = 2")),
		},
		{
			name:           "insert params multi in second col",
			expected:       "INSERT INTO client ( name, value, data ) VALUES ( ?,?,? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'", "'some data'"},
			tc:             obreron.Insert().Into("client").Col("name", "'some name'").Col("value, data", "'somemail@mail.net'", "'some data'"),
		},
	}...)

	return tcs
}

func TestInsertBuildCopiesParams(t *testing.T) {
	ins := obreron.Insert().Into("client").Col("name, value", "'some name'", "'somemail@mail.net'")

	_, p := ins.Build()
	ins.Close()

	if len(p) != 2 {
		t.Fatalf("expected 2 params, got %d", len(p))
	}

	if p[0] != "'some name'" || p[1] != "'somemail@mail.net'" {
		t.Fatalf("params mutated after Close(): %#v", p)
	}
}

// insertH3RegressionCases contiene los casos de tabla para el bug H3:
// obreron.InsertStament.Col genera placeholders sin comas ("??") en columns no-primeras
// cuando se pasan más de un parámetro.
//
// Raíz del bug (insert.go):
//
//	// Rama non-firstCol (BUGGY):
//	in.add(insP, "", strings.Repeat("?", len(p)), p...)
//	//                 ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//	//                 Con len(p)=2 produce "??" en vez de "?,?"
//
//	// Rama firstCol (CORRECTO, para referencia):
//	pp = strings.Repeat("?,", len(p)-1)  // "?," para n-1 params
//	pp += "?"                             // añade el último sin coma
//
// Fix de una línea en insert.go:258:
//
//	// Antes:
//	in.add(insP, "", strings.Repeat("?", len(p)), p...)
//	// Después:
//	in.add(insP, "", strings.Repeat("?,", len(p)-1)+"?", p...)
//
// Cómo verificar que el test falla ANTES del fix y pasa DESPUÉS:
//
//	$ go test -run TestInsertH3Regression -v       # debe fallar
//	$ # aplicar el fix en insert.go
//	$ go test -run TestInsertH3Regression -v       # debe pasar
//	$ go test ./...                                # no debe romper nada más
func insertH3RegressionCases() []struct {
	tc             *obreron.InsertStament
	name           string
	expected       string
	expectedParams []any
} {
	return []struct {
		tc             *obreron.InsertStament
		name           string
		expected       string
		expectedParams []any
	}{
		// ------------------------------------------------------------------ //
		// Caso base de referencia — YA FUNCIONA (1 param en posición no-primera)
		// Incluido para documentar qué sí funciona y como guardia de regresión.
		// ------------------------------------------------------------------ //
		{
			name:           "H3-baseline: second Col with 1 param (works before fix)",
			expected:       "INSERT INTO t ( a, b ) VALUES ( ?,? )",
			expectedParams: []any{1, 2},
			tc:             obreron.Insert().Into("t").Col("a", 1).Col("b", 2),
		},

		// ------------------------------------------------------------------ //
		// Casos que FALLAN antes del fix (H3 regression)
		// ------------------------------------------------------------------ //

		// Caso mínimo reproducible: segundo Col con 2 params.
		// Bug: strings.Repeat("?", 2) = "??" → VALUES ( ?,?? )
		// Fix: strings.Repeat("?,", 1)+"?" = "?,?" → VALUES ( ?,?,? )
		{
			name:           "H3-regression: second Col with 2 params",
			expected:       "INSERT INTO t ( a, b, c ) VALUES ( ?,?,? )",
			expectedParams: []any{1, 2, 3},
			tc:             obreron.Insert().Into("t").Col("a", 1).Col("b, c", 2, 3),
		},

		// Segundo Col con 3 params.
		// Bug: "???" → VALUES ( ?,??? )
		// Fix: "?,?,?" → VALUES ( ?,?,?,? )
		{
			name:           "H3-regression: second Col with 3 params",
			expected:       "INSERT INTO t ( a, b, c, d ) VALUES ( ?,?,?,? )",
			expectedParams: []any{1, 2, 3, 4},
			tc:             obreron.Insert().Into("t").Col("a", 1).Col("b, c, d", 2, 3, 4),
		},

		// Ambos: primer Col Y segundo Col con múltiples params.
		// Verifica que la rama firstCol sigue correcta tras el fix.
		// Bug: VALUES ( ?,?,?? ) — el segundo bloque de placeholders está roto.
		// Fix: VALUES ( ?,?,?,? )
		{
			name:           "H3-regression: both first and second Col with 2 params each",
			expected:       "INSERT INTO t ( a, b, c, d ) VALUES ( ?,?,?,? )",
			expectedParams: []any{1, 2, 3, 4},
			tc:             obreron.Insert().Into("t").Col("a, b", 1, 2).Col("c, d", 3, 4),
		},

		// Tres Cols: primero (1 param), segundo (2 params), tercero (2 params).
		// Verifica que el bug afecta a TODAS las posiciones no-primeras con >1 param.
		// Bug: VALUES ( ?,?,??,?? )
		// Fix: VALUES ( ?,?,?,?,? )
		{
			name:           "H3-regression: third Col with multiple params also broken",
			expected:       "INSERT INTO t ( a, b, c, d, e ) VALUES ( ?,?,?,?,? )",
			expectedParams: []any{1, 2, 3, 4, 5},
			tc:             obreron.Insert().Into("t").Col("a", 1).Col("b, c", 2, 3).Col("d, e", 4, 5),
		},

		// ColIf(true, ...) en posición no-primera con múltiples params.
		// ColIf delega en Col, así que hereda el mismo bug.
		{
			name:           "H3-regression: ColIf(true) in non-first position with 2 params",
			expected:       "INSERT INTO t ( a, b, c ) VALUES ( ?,?,? )",
			expectedParams: []any{1, 2, 3},
			tc:             obreron.Insert().Into("t").Col("a", 1).ColIf(true, "b, c", 2, 3),
		},

		// ColIf(false, ...) NO debe añadir nada — guardia de regresión.
		{
			name:           "H3-baseline: ColIf(false) in non-first position does not add cols",
			expected:       "INSERT INTO t ( a ) VALUES ( ? )",
			expectedParams: []any{1},
			tc:             obreron.Insert().Into("t").Col("a", 1).ColIf(false, "b, c", 2, 3),
		},

		// Primer Col con múltiples params + ColIf(true) con múltiples params.
		// La rama firstCol ya es correcta; este caso aísla que el fix no la rompe.
		{
			name:           "H3-regression: first Col multi-param + ColIf(true) multi-param",
			expected:       "INSERT INTO t ( a, b, c, d ) VALUES ( ?,?,?,? )",
			expectedParams: []any{1, 2, 3, 4},
			tc:             obreron.Insert().Into("t").Col("a, b", 1, 2).ColIf(true, "c, d", 3, 4),
		},
	}
}

// INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2 --- Got
// INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2

type DocupdaterBody struct {
	StartDateUTC      time.Time
	EndDateUTC        time.Time
	SPS               int
	UseExpirationDate int
	ClientID          int
	OfficeID          int
	ResourceID        int
	RemitterID        int
	UserID            int
	LoggedUserID      int
	DocumentTypeID    int
	DocumentNumber    int
	EndID             int
	StartID           int
	EndEpoch          int
}

func (b *DocupdaterBody) ThereIsEndDate() bool {
	return b.EndEpoch > 0
}

func specialCaseUpdateTargetQuery() *obreron.UpdateStm {
	body := DocupdaterBody{}

	// noCheckIDs verifica si bypassear o no el siguiente filtro donde sea invocado,
	// esto según si esta definido un documento para actualizar, o un rango.
	noCheckIDs := body.ResourceID < 0 && (body.StartID > 0 && body.EndID > 0)
	ob := obreron.Select().
		Col(colStr).
		Col(`vds.id_venta_documento_tributario`).
		From("venta_documento_tributario vds").
		Clause(" STRAIGHT_JOIN ", "cliente _c ON vds.id_cliente = _c.id_cliente ").
		Where("1=1").
		AndIf(body.ResourceID > 0, "vds.id_venta_documento_tributario = ?", body.ResourceID).
		AndIf(body.ResourceID <= 0 && body.StartID > 0 && body.EndID > 0, "vds.id_venta_documento_tributario BETWEEN  ? AND ?", body.StartID, body.EndID).
		AndIf(body.ResourceID <= 0 && body.DocumentTypeID > 0, "vds.id_tipo_documento_tributario = ?", body.DocumentTypeID).
		AndIf(body.ResourceID <= 0 && body.DocumentNumber > 0, "vds.num_doc_tributario = ?", body.DocumentNumber).
		AndIf(body.ResourceID <= 0 && body.UserID > 0, "vds.id_usuario = ?", body.UserID).
		AndIf(body.ResourceID <= 0 && body.RemitterID > 0, "vds.id_emisor = ?", body.RemitterID).
		AndIf(noCheckIDs && body.OfficeID > 0, "vds.id_sucursal = ?", body.OfficeID).
		AndIf(noCheckIDs && body.OfficeID == 0 && body.LoggedUserID > 0 && body.SPS > 0, "vds.id_sucursal IN (SELECT id_sucursal FROM usuario_sucursal WHERE id_usuario = ?)", body.LoggedUserID).
		AndIf(body.ResourceID <= 0 && body.ClientID > 0, "vds.id_cliente = ?", body.ClientID).
		AndIf(noCheckIDs && body.UseExpirationDate == 0, "vds.fecha_vencimiento_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(noCheckIDs && body.UseExpirationDate == 1, "vds.fecha_emission_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(noCheckIDs && body.UseExpirationDate == 0 && body.ThereIsEndDate(), "vds.fecha_emission_documento <= ?", body.EndDateUTC.Format("2006-01-02")).
		AndIf(noCheckIDs && body.UseExpirationDate == 1 && body.ThereIsEndDate(), "vds.fecha_vencimiento_documento <= ?", body.EndDateUTC.Format("2006-01-02"))

	defer obreron.CloseSelect(ob)

	up := obreron.Update("vw_docs_search v").
		ColSelect(ob, "det").
		Set("v.destinatarios = det.destinatarios")

	defer obreron.CloseUpdate(up)

	if body.ResourceID == 0 {
		up.Where("v.id_venta_documento_tributario = det.id_venta_documento_tributario")
	} else {
		up.Where("v.id_venta_documento_tributario = ? ", body.ResourceID)
	}

	return up
}

func specialCaseUpdateResendQuery() *obreron.UpdateStm {
	body := DocupdaterBody{}
	// noCheckIDs verifica si bypassear o no el siguiente filtro donde sea invocado,
	// esto según si esta definido un documento para actualizar, o un rango.
	noCheckIDs := body.ResourceID < 0 && (body.StartID > 0 && body.EndID > 0)

	ob := obreron.Select().
		Col(`IFNULL(
CONVERT(
	GROUP_CONCAT(
			CONCAT(COALESCE(de.nombre_destinatario,''), ':', COALESCE(de.email_destinatario,''), ':', de.id_detalle_envio_documento)
		) USING latin1
	),
	''
) AS destinatarios_reenvio`).
		Col("vds.id_venta_documento_tributario").
		From("venta_documento_tributario vds").
		Join("envio_documento e ON vds.id_venta_documento_tributario = e.id_venta_documento_tributario").
		Join("detalle_envio_documento de ON e.id_envio_documento = de.id_envio_documento").
		Where("1=1").
		AndIf(body.ResourceID > 0, "vds.id_venta_documento_tributario = ?", body.ResourceID).
		AndIf(body.ResourceID <= 0 && body.StartID > 0 && body.EndID > 0, "vds.id_venta_documento_tributario BETWEEN ? AND ?", body.StartID, body.EndID).
		AndIf(body.ResourceID <= 0 && body.DocumentTypeID > 0, "vds.id_tipo_documento_tributario = ?", body.DocumentTypeID).
		AndIf(body.ResourceID <= 0 && body.DocumentNumber > 0, "vds.num_doc_tributario = ?", body.DocumentNumber).
		AndIf(body.ResourceID <= 0 && body.UserID > 0, "vds.id_usuario = ?", body.UserID).
		AndIf(body.ResourceID <= 0 && body.RemitterID > 0, "vds.id_emisor = ?", body.RemitterID).
		AndIf(noCheckIDs && body.OfficeID > 0, "vds.id_sucursal = ?", body.OfficeID).
		AndIf(noCheckIDs && body.OfficeID == 0 && body.LoggedUserID > 0 && body.SPS > 0, "vds.id_sucursal IN (SELECT id_sucursal FROM usuario_sucursal WHERE id_usuario = ?)", body.LoggedUserID).
		AndIf(body.ResourceID <= 0 && body.ClientID > 0, "vds.id_cliente = ?", body.ClientID).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 0, "vds.fecha_vencimiento_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 1, "vds.fecha_emission_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 0 && body.ThereIsEndDate(), "vds.fecha_emission_documento <= ?", body.EndDateUTC.Format("2006-01-02")).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 1 && body.ThereIsEndDate(), "vds.fecha_vencimiento_documento <= ?", body.EndDateUTC.Format("2006-01-02")).
		And("e.tipo_envio = 0").
		GroupBy("vds.id_venta_documento_tributario")

	defer obreron.CloseSelect(ob)

	upd := obreron.Update("vw_docs_search v").
		ColSelect(ob, "det").
		Set("v.destinatarios_reenvio = det.destinatarios_reenvio").
		Where("v.id_venta_documento_tributario = det.id_venta_documento_tributario")

	return upd
}

func specialCaseUpdateShippingQuery() *obreron.UpdateStm {
	body := DocupdaterBody{}
	ob := obreron.Select().
		Col("1 AS es_despacho").
		Col("td.nombre_i18n_tipo").
		Col("IFNULL(d.id_sucursal_destino, 0) AS id_sucursal_destino").
		Col("IFNULL(s.nombre_sucursal,'') AS sucursal_destino").
		Col("d.recepcionado_destino").
		Col("vdt.id_venta_documento_tributario").
		From("venta_documento_tributario vdt").
		Join("tipo_documento_tributario tdoc ON tdoc.id_tipo_documento_tributario = vdt.id_tipo_documento_tributario").
		Join("detalle_venta_documento_tributario dvdt ON vdt.id_venta_documento_tributario = dvdt.id_venta_documento_tributario").
		Join("detalle_despacho dd ON dvdt.id_detalle_despacho = dd.id_detalle_despacho").
		Join("despacho d ON dd.id_despacho = d.id_despacho").
		Join("tipo_despacho td ON d.id_tipo_despacho = td.id_tipo_despacho").
		LeftJoin("sucursal s ON d.id_sucursal_destino = s.id_sucursal").
		Where("1=1").
		AndIf(body.ResourceID > 0, "vdt.id_venta_documento_tributario = ?", body.ResourceID).
		AndIf(body.ResourceID <= 0 && body.StartID > 0, "vdt.id_venta_documento_tributario >= ?", body.StartID).
		AndIf(body.ResourceID <= 0 && body.EndID > 0, "vdt.id_venta_documento_tributario <= ?", body.EndID).
		And("uso_documento = 2").
		GroupBy("vdt.id_venta_documento_tributario")

	defer obreron.CloseSelect(ob)

	upd := obreron.Update("vw_docs_search v").
		ColSelect(ob, "det").
		Set("v.es_despacho = det.es_despacho").
		Set("v.nombre_tipo_despacho = det.nombre_i18n_tipo").
		Set("v.id_sucursal_destino = det.id_sucursal_destino").
		Set("v.sucursal_destino = det.sucursal_destino").
		Set("v.recepcionado_destino = det.recepcionado_destino").
		Where("v.id_venta_documento_tributario = det.id_venta_documento_tributario")

	return upd
}

const (
	colStr = `TRIM( TRAILING ',' FROM CONCAT(	
	CASE WHEN _c.email_cliente IS NULL OR LENGTH(TRIM(COALESCE(_c.email_cliente,''))) = 0 
	    THEN ''
		ELSE CONCAT(COALESCE(_c.nombre_cliente,''), ' ', COALESCE(_c.apellido_cliente,''), ':', COALESCE(_c.email_cliente,'')) 
	END, ',', 
	IFNULL((
	    SELECT CONCAT(GROUP_CONCAT(
		    CASE WHEN _cc.email_contacto IS NULL OR LENGTH(TRIM(COALESCE(_cc.email_contacto,''))) = 0 
			    THEN ''
				ELSE CONCAT(COALESCE(_cc.nombre_contacto,''), ' ', COALESCE(_cc.apellido_contacto,''), ':', CONCAT(_cc.email_contacto,'')) 
			END
		), ',')
		FROM contacto_cliente _cc
		WHERE _cc.id_cliente = _c.id_cliente),''))) AS destinatarios`
)
