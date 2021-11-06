package main

import (
	"fmt"
	"log"

	"github.com/grpc-boot/base"
)

func main() {
	base.Green("from yaml:")
	i18n, err := base.NewI18nFromYaml("zh_CN", "msgs/")
	if err != nil {
		log.Fatal(err)
	}

	base.Fuchsia(i18n.T("param_err", ""))
	base.Red(i18n.T("param_err", "en"))
	base.Yellow(i18n.T("param_err", "zh_CN"))

	fmt.Println()

	base.Green("from json:")
	i18n, err = base.NewI18nFromJson("zh_CN", "msgs/")
	if err != nil {
		log.Fatal(err)
	}

	base.Fuchsia(i18n.T("param_err", ""))
	base.Red(i18n.T("param_err", "en"))
	base.Yellow(i18n.T("param_err", "zh_CN"))
}
