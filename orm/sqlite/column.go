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
	d string

	_type          string
	_primaryKey    bool
	_autoIncrement bool
	_unsigned      bool
	_size          int
	_scale         int
}

func (c *column) format() {
	c._autoIncrement = strings.Contains(c.t, "autoincrement")
	c._primaryKey = strings.Contains(c.t, "primary key")
	c._unsigned = strings.Contains(c.t, "unsigned")

	matches := typeRegexp.FindStringSubmatch(c.t)
	if len(matches) > 1 {
		c._type = matches[1]
		if len(matches) > 2 {
			items := strings.Split(matches[2], ",")
			c._size, _ = strconv.Atoi(items[0])
			if len(items) > 1 {
				c._scale, _ = strconv.Atoi(items[1])
			}
		}
	}
}

func (c *column) Unsigned() bool {
	return c._unsigned
}

func (c *column) Comment() string {
	return ""
}

func (c *column) Field() string {
	return c.f
}

func (c *column) Name() string {
	return utils.BigCamelByChar(c.f, '_')
}

func (c *column) CanNull() bool {
	return !strings.Contains(c.t, "not null")
}

func (c *column) Size() int {
	return c._size
}

func (c *column) Scale() int {
	return c._scale
}

func (c *column) Extra() string {
	return ""
}

func (c *column) IsPrimaryKey() bool {
	return c._primaryKey
}

func (c *column) AutoIncrement() bool {
	return c._autoIncrement
}

func (c *column) GoType() string {
	switch c._type {
	case "tinyint", "int8", "int2":
		if c._unsigned {
			return "uint8"
		}
		return "int8"
	case "smallint", "mediumint", "int":
		if c._unsigned {
			return "uint32"
		}
		return "int32"
	case "bigint", "integer":
		if c._unsigned {
			return "uint64"
		}
		return "int64"
	case "float", "double", "real", "decimal", "numeric":
		return "float64"
	case "BLOB":
		return "[]byte"
	default:
		return "string"
	}
}

func (c *column) Key() string {
	if c._primaryKey {
		return "PRI"
	}
	return ""
}

func (c *column) Default() string {
	return c.d
}

func (c *column) Enums() []string {
	return nil
}
