package obreron

import "testing"

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
	for i, tc := range updateTestCases() {
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

func BenchmarkUpdate(b *testing.B) {
	for _, tc := range updateTestCases() {
		b.ResetTimer()
		b.Run(tc.name, func(b2 *testing.B) {
			for i := 0; i < b2.N; i++ {
				_, _ = tc.tc.Build()
				tc.tc.Close()
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
			tc: Update("items").
				ColSelectIf(true, Select().Col("id, retail / wholesale AS markup, quantity").From("items"), "discounted").
				Set("items.retail = items.retail * 0.9").Set("a = 2").SetIf(true, "c = 3").
				Where("discounted.markup >= 1.3").
				And("discounted.quantity < 100").
				And("items.id = discounted.id").Y().In("regime_cliente", "'G01','G02', ?", "'G03'").AndIf(true, "2 = 2").Or("3 = 3").OrIf(true, "4 = 4").And("colX").Like("'%ago%'").AndIf(true, "colN").LikeIf(true, "'%oga%'").Y().In("colY", "1, 2, 3"),
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
