package orm

import "flag"

func GenerateCodeWithMysql() string {
	var (
		h string
		P string
		p string
		d string
		c string
	)

	flag.StringVar(&h, "h", "127.0.0.1", "-h: host地址")
	flag.StringVar(&P, "P", "3306", "-P: 端口号")
	flag.StringVar(&p, "p", "", "-p: 密码")
	flag.StringVar(&d, "d", "", "-d: dbname")
	flag.StringVar(&c, "c", "utf8", "-c: charset")
	flag.Parse()
	return ""
}
