package main

import (
	"fmt"
	"log"

	"github.com/grpc-boot/base/v2/components"
	"github.com/grpc-boot/base/v2/utils"
)

func main() {
	utils.Green("from yaml:")
	i18n, err := components.NewI18nFromYaml("zh_CN", "msgs/")
	if err != nil {
		log.Fatal(err)
	}

	utils.Fuchsia(i18n.T("param_err", ""))
	utils.Red(i18n.T("param_err", "en"))
	utils.Yellow(i18n.T("param_err", "zh_CN"))

	fmt.Println()

	utils.Green("from json:")
	i18n, err = components.NewI18nFromJson("zh_CN", "msgs/")
	if err != nil {
		log.Fatal(err)
	}

	utils.Fuchsia(i18n.T("param_err", ""))
	utils.Red(i18n.T("param_err", "en"))
	utils.Yellow(i18n.T("param_err", "zh_CN"))
}
