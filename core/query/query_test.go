package query

import "testing"

func TestMysqlQuery_Sql(t *testing.T) {
	// SELECT * FROM `user` WHERE (`id` IN(?,?)) args1: [34 45]
	args1 := []interface{}{}
	query1 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(AndWhere(FieldMap{
			"`id`": {"IN", 34, 45},
		}))
	sql1 := query1.Sql(&args1)
	query1.Close()

	t.Logf("sql1: %s args1: %+v", sql1, args1)

	// sql2: SELECT * FROM `user` WHERE (`id` = ?) args2: [45]
	args2 := []interface{}{}
	query2 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(AndWhere(FieldMap{
			"`id`": {45},
		}))
	sql2 := query2.Sql(&args2)
	query2.Close()

	t.Logf("sql2: %s args2: %+v", sql2, args2)

	// SELECT * FROM `user` WHERE (`name` LIKE ? AND `id` = ? AND `isdel` = ?) args3: [ma% 45 0]
	args3 := []interface{}{}
	query3 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(AndWhere(FieldMap{
			"`id`":    {45},
			"`isdel`": {0},
			"`name`":  {"LIKE", "ma%"},
		}))
	sql3 := query3.Sql(&args3)
	query3.Close()

	t.Logf("sql3: %s args3: %+v", sql3, args3)

	// sql4: SELECT * FROM `user` WHERE (`id` = ? AND `isdel` = ?) OR (`id` = ? AND `isdel` = ?) LIMIT 10,10 args4: [45 0 46 0]
	args4 := []interface{}{}
	query4 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(AndWhere(FieldMap{
			"`id`":    {45},
			"`isdel`": {0},
		})).
		Or(AndCondition(FieldMap{
			"`id`":    {46},
			"`isdel`": {0},
		})).
		Offset(10).
		Limit(10)
	sql4 := query4.Sql(&args4)
	query4.Close()

	t.Logf("sql4: %s args4: %+v", sql4, args4)

	// sql5: SELECT COUNT(id) AS num,kind FROM `user` WHERE (`id` > ? AND `isdel` = ?) GROUP BY kind args5: [45 0]
	args5 := []interface{}{}
	query5 := Acquire4Mysql().
		Select("COUNT(id) AS num", "kind").
		From("`user`").
		Where(AndWhere(FieldMap{
			"`id`":    {">", 45},
			"`isdel`": {0},
		})).
		Group("kind")
	sql5 := query5.Sql(&args5)
	query5.Close()

	t.Logf("sql5: %s args5: %+v", sql5, args5)

	// SELECT COUNT(`id`) AS `num`,`kind` FROM `user` WHERE (`id` > ? AND `isdel` = ?) GROUP BY kind ORDER BY `num` DESC LIMIT 1,100 args6: [45 0]
	args6 := []interface{}{}
	query6 := Acquire4Mysql().
		Select("COUNT(`id`) AS `num`", "`kind`").
		From("`user`").
		Where(AndWhere(FieldMap{
			"`id`":    {">", 45},
			"`isdel`": {0},
		})).
		Group("kind").
		Order("`num` DESC").
		Offset(1).
		Limit(100)
	sql6 := query6.Sql(&args6)
	query6.Close()

	t.Logf("sql6: %s args6: %+v", sql6, args6)
}
