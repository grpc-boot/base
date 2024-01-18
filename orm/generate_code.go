package orm

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/mysql"
	"github.com/grpc-boot/base/v2/orm/sqlite"
	"github.com/grpc-boot/base/v2/utils"
)

var (
	generatorMap = map[string]func(f *basis.Flag) (basis.Generator, error){
		"mysql":  mysql.Flag2Generator,
		"sqlite": sqlite.Flag2Generator,
	}
)

func Flag2Generator(f *basis.Flag) (basis.Generator, error) {
	gen, exists := generatorMap[f.DriveName()]
	if !exists {
		return nil, errors.New("unsupported driver Name")
	}

	return gen(f)
}

func GenerateCode(driveName, tableName, pkgName, outDir, templateFile string, generator basis.Generator) {
	t, err := generator.LoadTableSchema(tableName)
	if err != nil {
		utils.RedFatal("read table schema failed with error: %v\n", err)
	}

	var templateModel = basis.DefaultModelTemplate()
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

	err = utils.MkDir(outDir, 0766)
	if err != nil {
		utils.RedFatal("create dir: %s failed with error:%v", outDir, err)
	}

	var (
		code    = t.GenerateCode(driveName, templateModel, pkgName)
		outFile = fmt.Sprintf("%s/%s.go", strings.TrimSuffix(outDir, "/"), tableName)
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

func GenerateCodeWithFlag(f *basis.Flag) {
	generator, err := Flag2Generator(f)
	if err != nil {
		utils.RedFatal("init generator failed with error: %v", err)
	}

	tables, err := generator.ShowTables("")
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
		index    int
		inputStr string
		read     = true
	)

	for read {
		utils.Green("please enter the index of table：")
		_, err = fmt.Scanln(&inputStr)
		if err != nil {
			utils.Red("read failed with error: %v", err)
			continue
		}

		index, err = strconv.Atoi(inputStr)
		if err != nil {
			utils.Red("parse index failed with error: %v", err)
			continue
		}

		if index >= len(tables) || index < 0 {
			utils.Red("invalid index: %d", index)
			continue
		}

		read = false
	}

	GenerateCode(f.DriveName(), tables[index], f.PkgName(), f.OutDir(), f.TemplateFile(), generator)
}
