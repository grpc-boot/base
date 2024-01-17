package sqlite

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/grpc-boot/base/v2/utils"
)

var (
	typeRegexp = regexp.MustCompile(`^(\w+)(?:\(([^\)]+)\))?`)
)

type column struct {
	f string
	t string
	n string
	k string
	d string
	e string
	c string

	_type          string
	_primaryKey    bool
	_autoIncrement bool
	_unsigned      bool
	_size          int
	_scale         int
	_enumValues    []string
}

func (c *column) format() {
	c._autoIncrement = strings.Contains(c.e, "auto_increment")
	c._primaryKey = strings.Contains(c.k, "PRI")
	c._unsigned = strings.Contains(c.t, "unsigned")

	matches := typeRegexp.FindStringSubmatch(c.t)
	if len(matches) > 1 {
		c._type = strings.ToLower(matches[1])

		if len(matches) > 2 {
			items := strings.Split(matches[2], ",")
			if c._type == "enum" {
				c._enumValues = make([]string, len(items))
				for i := 0; i < len(items); i++ {
					c._enumValues[i] = strings.Trim(items[i], "'")
				}
			} else {
				c._size, _ = strconv.Atoi(items[0])
				if len(items) > 1 {
					c._scale, _ = strconv.Atoi(items[1])
				}
			}
		}
	}
}

func (c *column) Unsigned() bool {
	return c._unsigned
}

func (c *column) Comment() string {
	return c.c
}

func (c *column) Field() string {
	return c.f
}

func (c *column) Name() string {
	return utils.BigCamelByChar(c.f, '_')
}

func (c *column) CanNull() bool {
	return c.n == "YES"
}

func (c *column) Size() int {
	return c._size
}

func (c *column) Scale() int {
	return c._scale
}

func (c *column) Extra() string {
	return c.e
}

func (c *column) IsPrimaryKey() bool {
	return c._primaryKey
}

func (c *column) AutoIncrement() bool {
	return c._autoIncrement
}

func (c *column) GoType() string {
	switch c._type {
	case "bit":
		if c._size > 32 {
			return "uint64"
		}
		return "uint32"
	case "tinyint":
		if c._unsigned {
			return "uint8"
		}
		return "int8"
	case "smallint", "mediumint", "int":
		if c._unsigned {
			return "uint32"
		}
		return "int32"
	case "bigint":
		if c._unsigned {
			return "uint64"
		}
		return "int64"
	case "float", "double", "real", "decimal", "numeric":
		return "float64"
	case "tinytext", "mediumtext", "longtext", "text", "varchar", "char", "enum", "json", "datetime", "year", "date", "time", "timestamp":
		return "string"
	case "blob", "longblob", "varbinary":
		return "[]byte"
	default:
		return ""
	}
}

func (c *column) Key() string {
	return c.k
}

func (c *column) Default() string {
	return c.d
}

func (c *column) Enums() []string {
	return c._enumValues
}
