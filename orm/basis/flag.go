package basis

import (
	"flag"
	"os"
	"strconv"

	"github.com/grpc-boot/base/v2/utils"
)

type Flag struct {
	host         string
	password     string
	userName     string
	dbName       string
	port         string
	charset      string
	outDir       string
	templateFile string
	packageName  string
}

func (f *Flag) Host() string {
	return f.host
}

func (f *Flag) Port() uint32 {
	port, _ := strconv.ParseUint(f.port, 10, 32)
	return uint32(port)
}

func (f *Flag) DbName() string {
	return f.dbName
}

func (f *Flag) UserName() string {
	return f.userName
}

func (f *Flag) Password() string {
	return f.password
}

func (f *Flag) Charset() string {
	return f.charset
}

func (f *Flag) OutDir() string {
	return f.outDir
}

func (f *Flag) TemplateFile() string {
	return f.templateFile
}

func (f *Flag) PkgName() string {
	return f.packageName
}

func (f *Flag) Check() {
	if f.dbName == "" {
		utils.Red("required: -d")
		flag.Usage()
		os.Exit(1)
	}

	if f.password == "" {
		utils.Red("required: -p")
		flag.Usage()
		os.Exit(1)
	}

	utils.Green("using\ndb: %s\nhost: %s\nport: %s\nusername: %s\ncharset: %s\ntemplate file: %s\npackage name: %s\nout dir: %s",
		f.dbName,
		f.host,
		f.port,
		f.userName,
		f.charset,
		f.templateFile,
		f.packageName,
		f.outDir,
	)
}

func ParseFlag() *Flag {
	f := &Flag{}

	flag.StringVar(&f.host, "h", "127.0.0.1", "-h: host")
	flag.StringVar(&f.port, "P", "3306", "-P: port")
	flag.StringVar(&f.userName, "u", "root", "-u: username")
	flag.StringVar(&f.password, "p", "", "-p: password")
	flag.StringVar(&f.dbName, "d", "", "-d: dbname")
	flag.StringVar(&f.charset, "c", "utf8", "-c: charset")
	flag.StringVar(&f.outDir, "o", "./", "-o: out path")
	flag.StringVar(&f.templateFile, "t", "", "-t: template file")
	flag.StringVar(&f.packageName, "g", "models", "-g: package name")
	flag.Parse()

	return f
}
