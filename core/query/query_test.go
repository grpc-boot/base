package query

import "testing"

func TestInsert(t *testing.T) {
	sql1, args1 := Insert("`user`", Columns{"name", "is_on"}, []Row{
		{"mads", 1},
		{"asdf", 0},
	}, false)

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := Insert("`user`", Columns{"name", "is_on"}, []Row{
		{"asdf", 1},
	}, true)

	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestUpdate(t *testing.T) {
	sql1, args1 := Update("`user`", "is_on=1", Equal{"id", 1})

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := Update("`user`", "amount=amount+1", In{"id", Values{45, 6, 7}})

	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestDelete(t *testing.T) {
	sql1, args1 := Delete("`user`", Equal{"id", 1})

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := Delete("`user`", In{"id", Values{45, 6, 7}})

	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestMysqlQuery_Sql(t *testing.T) {
	query1 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(In{"id", Values{34, 45}})
	sql1, args1 := query1.Sql()
	query1.Close()

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	query2 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(Equal{"id", 45})
	sql2, args2 := query2.Sql()
	query2.Close()

	t.Logf("sql2: %s with args:%+v", sql2, args2)

	query3 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(AndCondition{
			BeginWith{"`name`", "ma"},
			Between{"id", 45, 67},
			Equal{"isdel", 0},
		})
	sql3, args3 := query3.Sql()
	query3.Close()

	t.Logf("sql3: %s with args:%+v", sql3, args3)

	query4 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(OrCondition{
			AndCondition{
				Contains{"`name`", "cc"},
				Gte{"id", 45},
				Equal{"isdel", 0},
			},
			AndCondition{
				Lte{"id", 46},
				Equal{"isdel", 0},
			},
		}).
		Offset(10).
		Limit(10)
	sql4, args4 := query4.Sql()
	query4.Close()

	t.Logf("sql4: %s with args:%+v", sql4, args4)

	query5 := Acquire4Mysql().
		Select("COUNT(id) AS num", "kind").
		From("`user`").
		Where(AndCondition{
			Lt{"id", 45},
			Equal{"isdel", 0},
		}).
		Group("kind")
	sql5, args5 := query5.Sql()
	query5.Close()

	t.Logf("sql5: %s with args:%+v", sql5, args5)

	query6 := Acquire4Mysql().
		Select("COUNT(`id`) AS `num`", "`kind`").
		From("`user`").
		Where(AndCondition{
			Gt{"id", 45},
			Equal{"isdel", 0},
		}).
		Group("kind").
		Order("`num` DESC").
		Offset(1).
		Limit(100)
	sql6, args6 := query6.Sql()
	query6.Close()

	t.Logf("sql6: %s with args:%+v", sql6, args6)
}
