package base

import (
	"strconv"
	"strings"

	"github.com/grpc-boot/base/v2/utils"
)

type Table struct {
	name    string
	primary Column
	columns []Column
}

func NewTable(name string, columns []Column) *Table {
	table := &Table{
		name:    name,
		columns: columns,
	}

	for _, col := range table.columns {
		if col.IsPrimaryKey() {
			table.primary = col
			break
		}
	}

	return table
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) Columns() []Column {
	return t.columns
}

func (t *Table) Primary() Column {
	return t.primary
}

func (t *Table) StructName() string {
	return utils.BigCamelByChar(t.name, '_')
}

func (t *Table) getThis() string {
	items := strings.Split(t.name, "_")
	if len(items) == 1 {
		if len(items) > 1 {
			return strings.ToLower(items[0][:2])
		}
	}

	ts := make([]byte, 0, len(items))
	for _, b := range items {
		if b != "" {
			ts = append(ts, strings.ToLower(b[:1])[0])
		}
	}

	return utils.Bytes2String(ts)
}

func (t *Table) convertCode(column Column) string {
	switch column.GoType() {
	case "int64":
		return "row.Int64"
	case "uint64":
		return "row.Uint64"
	case "int32":
		return "row.Int32"
	case "uint32":
		return "row.Uint32"
	case "int8":
		return "row.Int8"
	case "uint8":
		return "row.Uint8"
	case "string":
		return "row.String"
	case "float64":
		return "row.Float64"
	default:
		return "row.Bytes"
	}
}

func (t *Table) GenerateCode(template, packageName string) string {
	var (
		code            = strings.ReplaceAll(template, "{package}", packageName)
		primaryField    = "id"
		primaryProperty = "Id"
	)

	if t.Primary() != nil {
		primaryField = t.primary.Field()
		primaryProperty = t.primary.Name()
	}

	code = strings.ReplaceAll(code, "{struct}", t.StructName())
	code = strings.ReplaceAll(code, "{this}", t.getThis())
	code = strings.ReplaceAll(code, "{table}", t.name)
	code = strings.ReplaceAll(code, "{primary}", primaryField)
	code = strings.ReplaceAll(code, "{Primary}", primaryProperty)
	code = strings.ReplaceAll(code, "{fields}", t.buildFields())
	code = strings.ReplaceAll(code, "{columns}", t.buildColumns(true))
	code = strings.ReplaceAll(code, "{rows}", t.buildRows(true))
	code = strings.ReplaceAll(code, "{assignment}", t.buildAssignment())

	return strings.ReplaceAll(code, "{smallCamel}", utils.SmallCamelByChar(t.name, '_'))
}

func (t *Table) buildColumns(skipPrimary bool) string {
	var buf strings.Builder

	for _, col := range t.columns {
		if skipPrimary && col.IsPrimaryKey() {
			continue
		}

		buf.WriteString("\n\t\t")
		buf.WriteString("fm.Field(\"")
		buf.WriteString(col.Name())
		buf.WriteString("\"),")
	}

	return buf.String()
}

func (t *Table) buildRows(skipPrimary bool) string {
	var buf strings.Builder

	for _, col := range t.columns {
		if skipPrimary && col.IsPrimaryKey() {
			continue
		}

		buf.WriteString("\n\t\t")
		buf.WriteString("{smallCamel}.")
		buf.WriteString(col.Name())
		buf.WriteString(",")
	}

	return buf.String()
}

func (t *Table) buildFields() string {
	var buf strings.Builder

	for _, col := range t.columns {
		buf.WriteString("\n\t")
		buf.WriteString(col.Name())
		buf.WriteByte('\t')
		buf.WriteString(col.GoType())
		buf.WriteByte('\t')
		buf.WriteString("`json:\"")
		buf.WriteString(utils.SmallCamelByChar(col.Field(), '_'))
		buf.WriteByte('"')
		buf.WriteString(" field:")
		buf.WriteString(strconv.Quote(col.Field()))
		buf.WriteString(" size:")
		buf.WriteString(strconv.Quote(strconv.Itoa(col.Size())))

		buf.WriteString(" label:")
		buf.WriteString(strconv.Quote(col.Comment()))

		if col.Default() != "" {
			buf.WriteString(" default:")
			buf.WriteString(strconv.Quote(col.Default()))
		}

		if len(col.Enums()) > 0 {
			buf.WriteString(" enums:")
			buf.WriteString(strconv.Quote(strings.Join(col.Enums(), ",")))
		}

		buf.WriteByte('`')
	}

	return buf.String()
}

func (t *Table) buildAssignment() string {
	var buf strings.Builder
	for _, col := range t.columns {
		buf.WriteString("\n\t\t\t")
		buf.WriteString(col.Name())
		buf.WriteByte(':')
		buf.WriteString(t.convertCode(col))
		buf.WriteString("(fm.Field(\"")
		buf.WriteString(col.Name())
		buf.WriteString("\")),")
	}

	return buf.String()
}
