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
	name           string
	expected       string
	expectedParams []any
	tc             *SelectStament
}) {

	cases = []struct {
		name           string
		expected       string
		expectedParams []any
		tc             *SelectStament
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
	}

	return cases
}

func deleteTestCases() []struct {
	name           string
	expected       string
	expectedParams []any
	tc             *DeleteStament
} {
	return []struct {
		name           string
		expected       string
		expectedParams []any
		tc             *DeleteStament
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
			expected:       "DELETE FROM client WHERE client_id = 100",
			expectedParams: nil,
			tc:             Delete().From("client").Where("client_id = 100"),
		},
		{
			name:           "del where conditions",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?)",
			expectedParams: []any{"'G03'"},
			tc: Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "'G03'"),
		},
		{
			name:           "del where conditions limit",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) LIMIT ?",
			expectedParams: []any{"'G03'", 100},
			tc: Delete().From("client").
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "'G03'").
				Limit("?", 100),
		},
		{
			name:           "del where conditions limit -- shuffled",
			expected:       "DELETE FROM client WHERE client_id = 100 AND estado_cliente = 0 AND regime_cliente IN ('G01','G02', ?) LIMIT ?",
			expectedParams: []any{"'G03'", 100},
			tc: Delete().From("client").
				Limit("?", 100).
				Where("client_id = 100").
				And("estado_cliente = 0").
				Y().In("regime_cliente", "'G01','G02', ?", "'G03'"),
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
			expected:       "DELETE FROM client WHERE client_id = 100 ORDER BY ciudad LIMIT 10",
			expectedParams: nil,
			tc:             Delete().From("client").Where("client_id = 100").OrderBy("ciudad").Limit("10"),
		},
	}
}

func updateTestCases() (tcs []struct {
	name           string
	expected       string
	expectedParams []any
	tc             *UpdateStament
}) {
	tcs = append(tcs, []struct {
		name           string
		expected       string
		expectedParams []any
		tc             *UpdateStament
	}{
		{
			name:           "update simple",
			expected:       "UPDATE client SET status = 0",
			expectedParams: nil,
			tc:             Update("client").Set("status = 0"),
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
			tc:             Update("client").Set("status = 0").Where("status = ?", 1).OrderBy("ciudad").Limit("?", 10),
		},
		{
			name:           "update where and order limit",
			expected:       "UPDATE client SET status = 0 WHERE status = ? AND country = ? ORDER BY ciudad LIMIT ?",
			expectedParams: []any{1, "'CL'", 10},
			tc:             Update("client").Set("status = 0").Where("status = ?", 1).And("country = ?", "'CL'").OrderBy("ciudad").Limit("?", 10),
		},
		{
			name:           "update select",
			expected:       "UPDATE items ,( SELECT id, retail / wholesale AS markup, quantity FROM items ) discounted SET items.retail = items.retail * 0.9 WHERE discounted.markup >= 1.3 AND discounted.quantity < 100 AND items.id = discounted.id",
			expectedParams: nil,
			tc: Update("items").
				ColSelect(Select().Col("id, retail / wholesale AS markup, quantity").From("items"), "discounted").
				Set("items.retail = items.retail * 0.9").
				Where("discounted.markup >= 1.3").
				And("discounted.quantity < 100").
				And("items.id = discounted.id"),
		},
		{
			name:           "update join",
			expected:       "UPDATE business AS b JOIN business_geocode AS g ON b.business_id = g.business_id SET b.mapx = g.latitude, b.mapy = g.longitude WHERE (b.mapx = '' or b.mapx = 0) AND g.latitude > 0",
			expectedParams: nil,
			tc: Update("business AS b").
				Join("business_geocode AS g").On("b.business_id = g.business_id").
				Set("b.mapx = g.latitude, b.mapy = g.longitude").
				Where("(b.mapx = '' or b.mapx = 0)").And("g.latitude > 0"),
		},
	}...)
	return tcs
}

func insertTestCases() (tcs []struct {
	name           string
	expected       string
	expectedParams []any
	tc             *InsertStament
}) {
	tcs = append(tcs, []struct {
		name           string
		expected       string
		expectedParams []any
		tc             *InsertStament
	}{
		{
			name:           "simple insert",
			expected:       "INSERT INTO client ( name, value ) VALUES ( ?, ? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'"},
			tc: Insert().Into("client").
				Col("name, value", "'some name'", "'somemail@mail.net'"),
		},
		{
			name:           "simple insert params",
			expected:       "INSERT INTO client ( name, value ) VALUES ( ?, ? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'"},
			tc:             Insert().Into("client").Col("name, value", "'some name'", "'somemail@mail.net'"),
		},
		{
			name:           "simple insert params shuffled",
			expected:       "INSERT INTO client ( name, value ) VALUES ( ?, ? )",
			expectedParams: []any{"'some name'", "'somemail@mail.net'"},
			tc:             Insert().Col("name, value", "'some name'", "'somemail@mail.net'").Into("client"),
		},
		{
			name:           "simple insert params select",
			expected:       "INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2",
			expectedParams: nil,
			tc: Insert().
				Into("courses").
				ColSelect("name, location, gid", Select().Col("name, location, 1").From("courses").Where("cid = 2")),
		},
	}...)

	return tcs
}

// INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2 --- Got
// INSERT INTO courses ( name, location, gid ) SELECT name, location, 1 FROM courses WHERE cid = 2
