package orm

import (
	"fmt"
	"os"
	"strings"

	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/mysql"
	"github.com/grpc-boot/base/v2/utils"
)

func GenerateCodeWithMysql(f *basis.Flag) {
	var (
		templateFile  = f.TemplateFile()
		outDir        = f.OutDir()
		opt           = mysql.Flag2Options(f)
		templateModel = basis.DefaultModelTemplate()
	)

	if templateFile != "" {
		if exists, _ := utils.FileExists(templateFile); !exists {
			utils.RedFatal("file: %s not exists", templateFile)
		}

		data, err := os.ReadFile(templateFile)
		if err != nil {
			utils.RedFatal("read file: %s failed with error:%v", templateFile, err)
		}
		templateModel = utils.Bytes2String(data)
	}

	err := utils.MkDir(outDir, 0766)
	if err != nil {
		utils.RedFatal("create dir: %s failed with error:%v", outDir, err)
	}

	db, err := mysql.NewDb(opt)
	if err != nil {
		utils.RedFatal("init db failed with error: %v", err)
	}

	tables, err := db.ShowTables("")
	if err != nil {
		utils.RedFatal("show tables with error: %v", err)
	}

	if len(tables) == 0 {
		utils.RedFatal("no tables")
	}

	for index, table := range tables {
		utils.Green("【%d】: %s", index, table)
	}

	var (
		index int
		read  = true
	)

	for read {
		utils.Green("please enter the index of table：")
		_, err = fmt.Scanln(&index)
		if err != nil {
			utils.Red("read failed with error: %v", err)
			continue
		}

		if index >= len(tables) || index < 0 {
			utils.Red("invalid index: %d", index)
			continue
		}

		read = false
	}

	t, err := db.LoadTableSchema(tables[index])
	if err != nil {
		utils.RedFatal("read table schema failed with error: %v\n", err)
	}

	var (
		code    = t.GenerateCode(templateModel, f.PkgName())
		outFile = fmt.Sprintf("%s/%s.go", strings.TrimSuffix(outDir, "/"), tables[index])
		file    *os.File
	)

	file, err = os.Create(outFile)
	if err != nil {
		utils.RedFatal("create file: %s failed with error: %v", outFile, err)
	}

	defer file.Close()

	_, err = file.WriteString(code)
	if err != nil {
		utils.RedFatal("write code failed with error: %v", err)
	}

	utils.Green("The code has been written to the file：%s", outFile)
}
