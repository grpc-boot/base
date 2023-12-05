package query

import (
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/query/condition"
)

func TestInsert(t *testing.T) {
	current := time.Now().Unix()
	sql1, args1 := Insert("`user`", Columns{"name", "is_on", "created_at", "updated_at"}, []Row{
		{"mads", 1, current, current},
		{"asdf", 0, current, current},
	}, false)

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := Insert("`user`", Columns{"name", "is_on"}, []Row{
		{"asdf", 1},
	}, true)
	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestUpdate(t *testing.T) {
	sql1, args1 := Update("`user`", "is_on=1", condition.Equal{"id", 2})
	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := Update("`user`", "amount=amount+1", condition.In[int]{"id", kind.Slice[int]{45, 6, 7}})

	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestDelete(t *testing.T) {
	sql1, args1 := Delete("`user`", condition.Equal{"id", 1})

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := Delete("`user`", condition.In[uint8]{"id", kind.Slice[uint8]{45, 6, 7}})

	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestMysqlQuery_Sql(t *testing.T) {
	query1 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(condition.In[uint8]{"id", kind.Slice[uint8]{3, 5}})
	sql1, args1 := query1.Sql()
	query1.Close()

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	query2 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(condition.Equal{"id", 5})
	sql2, args2 := query2.Sql()
	query2.Close()

	t.Logf("sql2: %s with args:%+v", sql2, args2)

	query3 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(condition.And{
			condition.BeginWith{"`name`", "ma"},
			condition.Between{"id", 1, 7},
			condition.Equal{"is_on", 1},
		})
	sql3, args3 := query3.Sql()
	query3.Close()

	t.Logf("sql3: %s with args:%+v", sql3, args3)

	query4 := Acquire4Mysql().
		Select("*").
		From("`user`").
		Where(condition.Or{
			condition.And{
				condition.Contains{"`name`", "a"},
				condition.Gte{"id", 1},
				condition.Equal{"is_on", 1},
			},
			condition.And{
				condition.Lte{"id", 46},
				condition.Equal{"is_on", 0},
			},
		}).
		Offset(10).
		Limit(10)
	sql4, args4 := query4.Sql()
	query4.Close()

	t.Logf("sql4: %s with args:%+v", sql4, args4)

	query5 := Acquire4Mysql().
		Select("COUNT(id) AS num", "name").
		From("`user`").
		Where(condition.And{
			condition.Lt{"id", 5},
			condition.Equal{"is_on", 1},
		}).
		Group("name")
	sql5, args5 := query5.Sql()
	query5.Close()

	t.Logf("sql5: %s with args:%+v", sql5, args5)

	query6 := Acquire4Mysql().
		Select("COUNT(`id`) AS `num`", "`name`").
		From("`user`").
		Where(condition.And{
			condition.Gt{"id", 1},
			condition.Equal{"is_on", 1},
		}).
		Group("name").
		Order("`num` DESC").
		Offset(1).
		Limit(100)
	sql6, args6 := query6.Sql()
	query6.Close()

	t.Logf("sql6: %s with args:%+v", sql6, args6)

	query7 := Acquire4Mysql().
		Select("COUNT(`id`) AS `num`", "`name`").
		From("`user`").
		Where(condition.And{
			condition.Gt{"id", 1},
			condition.Equal{"is_on", 1},
		})
	sql7, args7 := query7.Count("id")
	query6.Close()

	t.Logf("sql7: %s with args:%+v", sql7, args7)
}
