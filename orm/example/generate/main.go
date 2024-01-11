package main

import (
	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/orm"
	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/mysql"
	"github.com/grpc-boot/base/v2/utils"
)

var (
	f  *basis.Flag
	db *mysql.Db
)

// go run ./main.go -d {dbname} -p {your password} -o {out dir}
// go run ./main.go -h
func main() {
	f = basis.ParseFlag()
	f.Check()

	// generate code to dir
	orm.GenerateCodeWithMysql(f)

	// use model
	testModel()
}

func testModel() {
	err := logger.InitZapWithOption(logger.Option{
		Level:      -1,
		Path:       "./logs/",
		TickSecond: 5,
		MaxDays:    1,
	})
	if err != nil {
		utils.RedFatal("init logger failed with error: %v", err)
	}

	db, err = mysql.NewDb(mysql.Flag2Options(f))
	if err != nil {
		utils.RedFatal("init db failed with error: %v", err)
	}

	// for example
	/*ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	user, err := models.DefaultUsers.FindOne(ctx, db, models.Users{Id: 1})
	if err != nil {
		utils.Red("find one failed with error: %v", err)
	} else {
		utils.Green("find one: %+v", user)
	}*/
}
