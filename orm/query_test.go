package orm

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/condition"
	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/mysql"
)

var (
	db *mysql.Db
)

func init() {
	opts := mysql.DefaultMysqlOption()
	opts.Password = "12345678"
	opts.UserName = "root"
	opts.Host = "127.0.0.1"
	opts.Port = 3306
	opts.DbName = "users"

	var err error
	db, err = mysql.NewDb(opts)
	if err != nil {
		panic(err)
	}
}

func TestMysql_GenerateCode(t *testing.T) {
	tables, err := db.ShowTables("")
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	for _, table := range tables {
		tab, err := db.LoadTableSchema(table)
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}

		f, err := os.Create(fmt.Sprintf("./models/%s.go", table))
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}

		f.WriteString(tab.GenerateCode("mysql", basis.DefaultModelTemplate(), "models"))
		f.Close()
	}
}

func TestInsert(t *testing.T) {
	current := time.Now().Unix()
	sql1, args1 := mysql.Insert("`user`", basis.Columns{"name", "is_on", "created_at", "updated_at"}, []basis.Row{
		{"mads", 1, current, current},
		{"asdf", 0, current, current},
	}, false)

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := mysql.Insert("`user`", basis.Columns{"name", "is_on"}, []basis.Row{
		{"asdf", 1},
	}, true)
	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestUpdate(t *testing.T) {
	sql1, args1 := mysql.Update("`user`", "is_on=1", condition.Equal{"id", 2})
	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := mysql.Update("`user`", "amount=amount+1", condition.In[int]{"id", kind.Slice[int]{45, 6, 7}})

	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestDelete(t *testing.T) {
	sql1, args1 := mysql.Delete("`user`", condition.Equal{"id", 1})

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	sql2, args2 := mysql.Delete("`user`", condition.In[uint8]{"id", kind.Slice[uint8]{45, 6, 7}})

	t.Logf("sql2: %s with args:%+v", sql2, args2)
}

func TestMysqlQuery_Sql(t *testing.T) {
	query1 := mysql.AcquireQuery().
		Select("*").
		From("`user`").
		Where(condition.In[uint8]{"id", kind.Slice[uint8]{3, 5}})
	sql1, args1 := query1.Sql()
	query1.Close()

	t.Logf("sql1: %s with args:%+v", sql1, args1)

	query2 := mysql.AcquireQuery().
		Select("*").
		From("`user`").
		Where(condition.Equal{"id", 5})
	sql2, args2 := query2.Sql()
	query2.Close()

	t.Logf("sql2: %s with args:%+v", sql2, args2)

	query3 := mysql.AcquireQuery().
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

	query4 := mysql.AcquireQuery().
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

	query5 := mysql.AcquireQuery().
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

	query6 := mysql.AcquireQuery().
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

	query7 := mysql.AcquireQuery().
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
