package obreron

import (
	"testing"
	"time"

	"github.com/pingcap/tidb/pkg/parser"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
)

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
	for k := 0; k <= 10000; k++ {
		for i, tc := range updateTestCases() {

			sql, p := tc.tc.Build()

			
			// Parser para verificar que el sql es correctamente construido
			_, _, err := par.ParseSQL(sql)

			if err != nil {
				t.Logf("[TEST CASE %d  %s] %v", i, tc.name, err)
				t.FailNow()
			}

			if tc.expected != "" {
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
			}

			CloseUpdate(tc.tc)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	for _, tc := range updateTestCases() {
		b.ResetTimer()
		b.Run(tc.name, func(b2 *testing.B) {
			for i := 0; i < b2.N; i++ {
				_, _ = tc.tc.Build()
				CloseUpdate(tc.tc)
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
	tc             *SelectStm
	name           string
	expected       string
	expectedParams []any
}) {

	cases = []struct {
		tc             *SelectStm
		name           string
		expected       string
		expectedParams []any
	}{
		/* SELECT TESTS */
		{
			name:           "columns - from",
			expected:       "SELECT a1, a2, a3 FROM client",
			expectedParams: nil,
			tc:             Select().Col("a1, a2, a3").From("client"),
		},
		{
			name:           "columns - from",
			expected:       "SELECT a1, a2, a3 FROM client",
			expectedParams: nil,
			tc:             Select().Col("a1, a2, a3").From("client"),
		},
		{
			name:           "columns params - from",
			expected:       "SELECT a1, a2, ? AS cien FROM client",
			expectedParams: []any{100},
			tc:             Select().Col("a1, a2, ? AS cien", 100).From("client"),
		},
		{
			name:           "columns params - columns params - from",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client",
			expectedParams: []any{10, 100},
			tc:             Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).From("client"),
		},
		{
			name:           "columns params - columns params - Col If - from",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client",
			expectedParams: []any{10, 100},
			tc:             Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client"),
		},
		{
			name:           "columns params - columns params - Col If - from - where",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client WHERE 1 = 1",
			expectedParams: []any{10, 100},
			tc:             Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client").Where("1 = 1"),
		},
		{
			name:           "columns params - Col If - columns params - Col If - from - join if -where",
			expected:       "SELECT SQL NO CACHE b2 AS b2, a1, a2, ? AS diez, a3, ? AS cien FROM client JOIN tablada t ON t.id=client.id LEFT JOIN tabla_b tb ON tb.id = client.id LEFT JOIN tabla_c tc ON tc.id = client.id RIGHT JOIN tabla_e te ON te.id = client.id AND 2 = 2 RIGHT JOIN tabla_f tf ON tf.id = client.id OUTER JOIN tabla_h th ON th.id = client.id AND 2 = 2 OUTER JOIN tabla_i ti ON ti.id = client.id WHERE 1 = 1 AND 3 = 3 AND colX LIKE '%chamullo%' OR 4 = 4 OR 5 = 5 GROUP BY colY , colZ HAVING colZ = 'abc' ORDER BY 1 ASC LIMIT ? OFFSET ?",
			expectedParams: []any{10, 100, 100, 200},
			tc: Select().Clause("SQL NO", "").ClauseIf(true, "CACHE", "").
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
			tc:             Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client").Where("a1 = ?", "'last name'"),
		},
		{
			name:           "columns params - columns params - Col If - from - where",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client WHERE a1 = ? AND a2 = ?",
			expectedParams: []any{10, 100, "'last name'", 1000.54},
			tc:             Select().Col("a1, a2, ? AS diez", 10).Col("a3, ? AS cien", 100).ColIf(false, "a4").From("client").Where("a1 = ?", "'last name'").And("a2 = ?", 1000.54),
		},
		{
			name:           "columns params - columns params - Col If - from - where shuffled",
			expected:       "SELECT a1, a2, ? AS diez, a3, ? AS cien FROM client WHERE a1 = ? AND a2 = ?",
			expectedParams: []any{10, 100, "'last name'", 1000.54},
			tc: Select().
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
			tc: Select().
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
			tc: Select().
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
			tc:             Select().Col("a1, a2, a3").From("client").Where("status").In("0, 1, 2, 3"),
		},
		{
			name:           "columns - where in",
			expected:       "SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND status IN ( 0, 1, 2, 3)",
			expectedParams: nil,
			tc:             Select().Col("a1, a2, a3").From("client").Where("1 = 1").And("status").In("0, 1, 2, 3"),
		},
		{
			name:           "columns - where in",
			expected:       "SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND status IN (?, ?, ?, ?)",
			expectedParams: []any{0, 1, 2, 3},
			tc:             Select().Col("a1, a2, a3").From("client").Where("1 = 1").Y().InArgs("status", 0, 1, 2, 3),
		},
		{
			name:           "columns - where like",
			expected:       "SELECT a1, a2, a3 FROM client WHERE 1 = 1 AND city LIKE '%ago%'",
			expectedParams: nil,
			tc:             Select().Col("a1, a2, a3").From("client").Where("1 = 1").And("city").Like("'%ago%'"),
		},
	}

	return cases
}

func deleteTestCases() []struct {
	tc             *DeleteStm
	name           string
	expected       string
	expectedParams []any
} {
	return []struct {
		tc             *DeleteStm
		name           string
		expected       string
		expectedParams []any
	}{
		/* DELETE TESTS */
		{
			name:           "simple del",
			expected:       "DELETE FROM client",
			expectedParams: nil,
			tc:             Delete().From("client"),
		},
		{
			name:           "simple del where",
			expected:       "DELETE FROM client WHERE client_id = 100 AND b = 3 OR 2 = 2 OR 3 = 3",
			expectedParams: nil,
			tc:             Delete().From("client").Where("client_id = 100").AndIf(true, "b = 3").OrIf(true, "2 = 2").Or("3 = 3"),
		},
		{
			name:           "del where conditions",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) AND a LIKE '%ago%' -- Comment\n",
			expectedParams: []any{"'G03'"},
			tc: Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "'G03'").And("a").LikeIf(true, "'%ago%'").ClauseIf(true, "-- Comment\n", ""),
		},
		{
			name:           "del where conditions",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN (?, ?, ?) AND a LIKE '%ago%' -- Comment\n",
			expectedParams: []any{"G01", "G02", "G03"},
			tc: Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().InArgs("regime_cliente", "G01", "G02", "G03").And("a").LikeIf(true, "'%ago%'").ClauseIf(true, "-- Comment\n", ""),
		},
		{
			name:           "del where conditions limit",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) LIMIT ?",
			expectedParams: []any{"'G03'", 100},
			tc: Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "'G03'").
				Limit(100),
		},
		{
			name:           "del where conditions limit -- shuffled",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) LIMIT ?",
			expectedParams: []any{"G03", 100},
			tc: Delete().From("client").
				Limit(100).
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "G03"),
		},
		{
			name:           "simple del where quick",
			expected:       "DELETE QUICK FROM client WHERE client_id = 100",
			expectedParams: nil,
			tc:             Delete().Clause("QUICK", "").From("client").Where("client_id = 100"),
		},
		{
			name:           "simple del where ignore",
			expected:       "DELETE IGNORE FROM client WHERE client_id = 100",
			expectedParams: nil,
			tc:             Delete().Clause("IGNORE", "").From("client").Where("client_id = 100"),
		},
		{
			name:           "simple del where partition",
			expected:       "DELETE PARTITION the_partition FROM client WHERE client_id = 100",
			expectedParams: nil,
			tc:             Delete().Clause("PARTITION", "the_partition").From("client").Where("client_id = 100"),
		},
		{
			name:           "simple del where order by limit",
			expected:       "DELETE FROM client WHERE client_id = 100 ORDER BY ciudad LIMIT ?",
			expectedParams: []any{10},
			tc:             Delete().From("client").Where("client_id = 100").OrderBy("ciudad").Limit(10),
		},
	}
}

func updateTestCases() (tcs []struct {
	tc             *UpdateStm
	name           string
	expected       string
	expectedParams []any
}) {
	tcs = append(tcs, []struct {
		tc             *UpdateStm
		name           string
		expected       string
		expectedParams []any
	}{

		{
			name:           "update simple",
			expected:       "UPDATE client SET status = 0",
			expectedParams: nil,
			tc:             Update("client").Set("status = 0"),
		},
		{
			name:           "update simple",
			expected:       "UPDATE client SET status = 0, name = ?",
			expectedParams: []any{"stitch"},
			tc:             Update("client").Set("status = 0").Set("name = ?", "stitch"),
		},
		{
			name:           "update simple",
			expected:       "UPDATE client SET status = 0 WHERE country = ? AND status IN (?, ?, ?, ?)",
			expectedParams: []any{"CL", 1, 2, 3, 4},
			tc:             Update("client").Set("status = 0").Where("country = ?", "CL").Y().In("status", "?, ?, ?, ?", 1, 2, 3, 4),
		},
		{
			name:           "update simple",
			expected:       "UPDATE client SET status = 0 WHERE country = ? AND status IN (?, ?, ?, ?)",
			expectedParams: []any{"CL", 1, 2, 3, 4},
			tc:             Update("client").Set("status = 0").Where("country = ?", "CL").Y().InArgs("status", 1, 2, 3, 4),
		},
		{
			name:           "update with IN with 1 param in args",
			expected:       "UPDATE client SET status = 0 WHERE country = ? AND status IN (?)",
			expectedParams: []any{"CL", 1},
			tc:             Update("client").Set("status = 0").Where("country = ?", "CL").Y().InArgs("status", 1),
		},
		{
			name:           "update where",
			expected:       "UPDATE client SET status = 0 WHERE status = ?",
			expectedParams: []any{1},
			tc:             Update("client").Set("status = 0").Where("status = ?", 1),
		},
		{
			name:           "update where order limit",
			expected:       "UPDATE client SET status = 0 WHERE status = ? ORDER BY ciudad LIMIT ?",
			expectedParams: []any{1, 10},
			tc:             Update("client").Set("status = 0").Where("status = ?", 1).OrderBy("ciudad").Limit(10),
		},
		{
			name:           "update where and order limit",
			expected:       "UPDATE client SET status = 0 WHERE status = ? AND country = ? ORDER BY ciudad LIMIT ?",
			expectedParams: []any{1, "'CL'", 10},
			tc:             Update("client").Set("status = 0").Where("status = ?", 1).And("country = ?", "'CL'").OrderBy("ciudad").Limit(10),
		},
		// UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9, a = 2, c = 3 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id AND regime_cliente IN ('G01','G02', ?) AND 2 = 2 OR 3 = 3 OR 4 = 4 AND colX LIKE '%ago%' AND colN LIKE '%oga%' AND colY (1, 2, 3) --- Got
		// UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9, a = 2, c = 3 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id AND regime_cliente IN ('G01','G02', ?) AND 2 = 2 OR 3 = 3 OR 4 = 4 AND colX LIKE '%ago%' AND colN LIKE '%oga%' AND colY LIKE (1, 2, 3)
		{
			name:           "update select",
			expected:       "UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9, a = 2, c = 3 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id AND regime_cliente IN ('G01','G02', ?) AND 2 = 2 OR 3 = 3 OR 4 = 4 AND colX LIKE '%ago%' AND colN LIKE '%oga%' AND colY IN (1, 2, 3)",
			expectedParams: []any{"'G03'"},
			tc: func() *UpdateStm {
			    s := Select().Col("id, retail / wholesale AS markup, quantity").From("items");
				defer CloseSelect(s)

			    return Update("items").
				ColSelectIf(true, s, "discounted").
				Set("items.retail = items.retail * 0.9").Set("a = 2").SetIf(true, "c = 3").
				Where("discounted.markup >= 1.3").
				And("discounted.quantity < 100").
				And("items.id = discounted.id").Y().In("regime_cliente", "'G01','G02', ?", "'G03'").AndIf(true, "2 = 2").Or("3 = 3").OrIf(true, "4 = 4").And("colX").Like("'%ago%'").AndIf(true, "colN").LikeIf(true, "'%oga%'").Y().In("colY", "1, 2, 3")

			}(),
		},
		// UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id --- Got
		// UPDATE items ,( SELECT , id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id
		{
			name:           "update join",
			expected:       "UPDATE business AS b JOIN business_geocode AS g ON b.business_id = g.business_id SET b.mapx = g.latitude, b.mapy = g.longitude WHERE (b.mapx = '' or b.mapx = 0) AND g.latitude > 0 AND 3 = 3",
			expectedParams: nil,
			tc: Update("business AS b").
				JoinIf(true, "business_geocode AS g").OnIf(true, "b.business_id = g.business_id").
				Set("b.mapx = g.latitude, b.mapy = g.longitude").
				Where("(b.mapx = '' or b.mapx = 0)").And("g.latitude > 0").ClauseIf(true, "AND", "3 = 3"),
		},
		{
			tc: Update("items").
				Set("items.retail = items.retail * 0.9").
				Set("a = 2").
				Where("discounted.markup >= 1.3").
				And("colX").
				Like("'%ago%'"),
			name:           "",
			expected:       "UPDATE items SET items.retail = items.retail * 0.9, a = 2 WHERE discounted.markup >= 1.3 AND colX LIKE '%ago%'",
			expectedParams: []any{},
		},
		{
			tc: specialCase_docupdater_UpdateTargetQuery(),
			name: "",
			expected: "",
			expectedParams: []any{},
		},
		{
			tc: specialCase_docupdater_UpdateResendQuery(),
			name: "",
			expected: "",
			expectedParams: []any{},
		},
		{
			tc: specialCase_docupdater_UpdateShippingQuery(),
			name: "",
			expected: "",
			expectedParams: []any{},
		},
	}...)
	return tcs
}

func insertTestCases() (tcs []struct {
	tc             *InsertStament
	name           string
	expected       string
	expectedParams []any
}) {
	tcs = append(tcs, []struct {
		tc             *InsertStament
		name           string
		expected       string
		expectedParams []any
	}{
		{
			name:           "simple insert",
			expected:       "INSERT IGNORE INTO client ( name, value ) VALUES ( ?,? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'"},
			tc: Insert().Ignore().Into("client").
				Col("name, value", "'some name'", "'somemail@mail.net'"),
		},
		{
			name:           "simple insert params",
			expected:       "INSERT INTO client ( name, value, data ) VALUES ( ?,?,? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'", "'some data'"},
			tc:             Insert().Into("client").Col("name", "'some name'").Col("value", "'somemail@mail.net'").ColIf(true, "data", "'some data'").ColIf(false, "info", 12),
		},
		{
			name:           "simple insert params shuffled",
			expected:       "INSERT INTO client ( name, value ) VALUES ( ?,? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'"},
			tc:             Insert().Col("name, value", "'some name'", "'somemail@mail.net'").Into("client"),
		},
		{
			name:           "simple insert params select",
			expected:       "INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2",
			expectedParams: nil,
			tc: Insert().
				Into("courses").
				ColSelectIf(true, "name, location, gid", Select().Col("name, location, 1").From("courses").Where("cid = 2")).
				ColSelectIf(false, "last_name, last_location, grid", Select().Col("last_name, last_location, 11").From("courses").Where("cid = 2")),
		},
	}...)

	return tcs
}

// INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2 --- Got
// INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2



type DocupdaterBody struct {
  ResourceID int
  StartID int
  EndID int
  UseExpirationDate int
  ClientID int
  OfficeID int
  SPS int
  RemitterID int
  UserID int
  LoggedUserID int
  DocumentTypeID int
  DocumentNumber int
  StartDateUTC time.Time
  EndDateUTC time.Time
  EndEpoch int
}

func (b *DocupdaterBody) ThereIsEndDate() bool {
	return b.EndEpoch > 0
}


func specialCase_docupdater_UpdateTargetQuery() *UpdateStm {
	body := DocupdaterBody{}

	// noCheckIds verifica si bypassear o no el siguiente filtro donde sea invocado,
	// esto según si esta definido un documento para actualizar, o un rango.
	noCheckIds := !(body.ResourceID > 0 || (body.StartID > 0 && body.EndID > 0))
	ob := Select().
		Col(`TRIM( TRAILING ',' FROM CONCAT(	
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
		WHERE _cc.id_cliente = _c.id_cliente),''))) AS destinatarios`).
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
		AndIf(noCheckIds && body.OfficeID > 0, "vds.id_sucursal = ?", body.OfficeID).
		AndIf(noCheckIds && body.OfficeID == 0 && body.LoggedUserID > 0 && body.SPS > 0, "vds.id_sucursal IN (SELECT id_sucursal FROM usuario_sucursal WHERE id_usuario = ?)", body.LoggedUserID).
		AndIf(body.ResourceID <= 0 && body.ClientID > 0, "vds.id_cliente = ?", body.ClientID).
		AndIf(noCheckIds && body.UseExpirationDate == 0, "vds.fecha_vencimiento_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(noCheckIds && body.UseExpirationDate == 1, "vds.fecha_emision_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(noCheckIds && body.UseExpirationDate == 0 && body.ThereIsEndDate(), "vds.fecha_emision_documento <= ?", body.EndDateUTC.Format("2006-01-02")).
		AndIf(noCheckIds && body.UseExpirationDate == 1 && body.ThereIsEndDate(), "vds.fecha_vencimiento_documento <= ?", body.EndDateUTC.Format("2006-01-02"))

	defer CloseSelect(ob)

	up := Update("vw_docs_search v").
		ColSelect(ob, "det").
		Set("v.destinatarios = det.destinatarios")

	if body.ResourceID > 0 {
		up.Where("v.id_venta_documento_tributario = det.id_venta_documento_tributario")
	} else {
		up.Where("v.id_venta_documento_tributario = ? ", body.ResourceID)
	}

	return up
}



func specialCase_docupdater_UpdateResendQuery() *UpdateStm {
	body := DocupdaterBody{}
	// noCheckIds verifica si bypassear o no el siguiente filtro donde sea invocado,
	// esto según si esta definido un documento para actualizar, o un rango.
	noCheckIds := !(body.ResourceID > 0 || (body.StartID > 0 && body.EndID > 0))

	ob := Select().
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
		AndIf(noCheckIds && body.OfficeID > 0, "vds.id_sucursal = ?", body.OfficeID).
		AndIf(noCheckIds && body.OfficeID == 0 && body.LoggedUserID > 0 && body.SPS > 0, "vds.id_sucursal IN (SELECT id_sucursal FROM usuario_sucursal WHERE id_usuario = ?)", body.LoggedUserID).
		AndIf(body.ResourceID <= 0 && body.ClientID > 0, "vds.id_cliente = ?", body.ClientID).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 0, "vds.fecha_vencimiento_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 1, "vds.fecha_emision_documento >= ?", body.StartDateUTC.Format("2006-01-02")).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 0 && body.ThereIsEndDate(), "vds.fecha_emision_documento <= ?", body.EndDateUTC.Format("2006-01-02")).
		AndIf(body.ResourceID <= 0 && body.UseExpirationDate == 1 && body.ThereIsEndDate(), "vds.fecha_vencimiento_documento <= ?", body.EndDateUTC.Format("2006-01-02")).
		And("e.tipo_envio = 0").
		GroupBy("vds.id_venta_documento_tributario")

	defer CloseSelect(ob)

	upd := Update("vw_docs_search v").
		ColSelect(ob, "det").
		Set("v.destinatarios_reenvio = det.destinatarios_reenvio").
		Where("v.id_venta_documento_tributario = det.id_venta_documento_tributario")
	return upd
}


func specialCase_docupdater_UpdateShippingQuery() *UpdateStm {
    body := DocupdaterBody{}
	ob := Select().
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

	defer CloseSelect(ob)

	upd := Update("vw_docs_search v").
		ColSelect(ob, "det").
		Set("v.es_despacho = det.es_despacho").
		Set("v.nombre_tipo_despacho = det.nombre_i18n_tipo").
		Set("v.id_sucursal_destino = det.id_sucursal_destino").
		Set("v.sucursal_destino = det.sucursal_destino").
		Set("v.recepcionado_destino = det.recepcionado_destino").
		Where("v.id_venta_documento_tributario = det.id_venta_documento_tributario")
	return upd
}