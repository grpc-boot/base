package main

import (
	"context"
	"github.com/grpc-boot/base/v2/orm/condition"
	"time"

	"github.com/grpc-boot/base/v2/orm"
	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/example/generate/models"
	"github.com/grpc-boot/base/v2/orm/sqlite"
	"github.com/grpc-boot/base/v2/utils"
)

var (
	f  *basis.Flag
	db basis.Generator
)

// go run ./main.go -d {dbname} -p {your password} -o {out dir}
// go run ./main.go -d user.db -dr sqlite -o ./models
// go run ./main.go -h
func main() {
	f = basis.ParseFlag()
	f.Check()

	if f.DriveName() == "sqlite" {
		createTable4Sqlite(f)
	}

	// generate code to dir
	orm.GenerateCodeWithFlag(f)

	testModel(f)
}

func createTable4Sqlite(f *basis.Flag) {
	g, err := sqlite.Flag2Generator(f)
	if err != nil {
		utils.RedFatal("init db failed with error:%v", err)
	}

	tableSql := `CREATE TABLE IF NOT EXISTS user (
    id INTEGER NOT NULL PRIMARY KEY autoincrement,
    user_name VARCHAR(64) NOT NULL DEFAULT '',
    head_img BLOB NOT NULL,
    created_at BIGINT unsigned DEFAULT 0,
    status INT8 DEFAULT 0,
    amount REAL NOT NULL,
    remark TEXT DEFAULT ''
)`

	_, err = g.ExecContext(context.Background(), tableSql)
	if err != nil {
		utils.RedFatal("create table failed with error:%v", err)
	}

	indexsSql, _ := sqlite.ShowIndexs("user")
	rows, err := g.QueryContext(context.Background(), indexsSql)
	if err != nil {
		utils.RedFatal("query index failed with error:%v", err)
	}

	records, err := basis.ScanRecords(rows)
	if err != nil {
		utils.RedFatal("scan index failed with error:%v", err)
	}

	utils.Green("scan index list: %v", records)

	hasIndex := false
	for _, record := range records {
		if record["name"] == "created_at" {
			hasIndex = true
			break
		}
	}

	if !hasIndex {
		indexSql, _ := sqlite.AddIndex("created_at", "user", true, "created_at")
		_, err = g.ExecContext(context.Background(), indexSql)
		if err != nil {
			utils.RedFatal("create index failed with error:%v", err)
		}
	}
}

func testModel(f *basis.Flag) {
	db, _ = orm.Flag2Generator(f)
	user := models.User{
		UserName:  "u1",
		HeadImg:   "https://sfasdf.png",
		CreatedAt: uint64(time.Now().Unix()),
		Status:    100,
		Remark:    "remark test",
		Amount:    10.1234,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id, err := user.Add(ctx, db, user)
	if err != nil {
		utils.RedFatal("insert failed with error: %s", err)
	}
	utils.Green("insert id: %d", id)

	query := user.Query().
		IndexBy("created_at").
		Where(condition.Gt{
			Field: "created_at",
			Value: 0,
		}).
		Group("id").
		Order("created_at DESC, id ASC").
		Limit(10)

	userList, err := user.FindByQuery(ctx, db, query)
	if err != nil {
		utils.RedFatal("query failed with error: %s", err)
	}
	utils.Green("user list: %v", userList)
}
