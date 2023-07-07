package query

import "strings"

type Values []interface{}

type Condition interface {
	Build() (sql string, args []interface{})
}

type AndCondition []Condition

type OrCondition []Condition

func (ac AndCondition) Build() (sql string, args []interface{}) {
	if len(ac) == 0 {
		return
	}

	var buffer strings.Builder

	args = make([]interface{}, 0)

	buffer.WriteByte('(')
	sql, buildArgs := ac[0].Build()
	buffer.WriteString(sql)
	args = append(args, buildArgs...)
	for index := 1; index < len(ac); index++ {
		sql, buildArgs = ac[index].Build()
		buffer.WriteString(" AND ")
		buffer.WriteString(sql)
		args = append(args, buildArgs...)

	}
	buffer.WriteByte(')')
	return buffer.String(), args
}

func (oc OrCondition) Build() (sql string, args []interface{}) {
	if len(oc) == 0 {
		return
	}

	var buffer strings.Builder

	args = make([]interface{}, 0)

	buffer.WriteByte('(')
	sql, buildArgs := oc[0].Build()
	buffer.WriteString(sql)
	args = append(args, buildArgs...)
	for index := 1; index < len(oc); index++ {
		sql, buildArgs = oc[index].Build()
		buffer.WriteString(" OR ")
		buffer.WriteString(sql)
		args = append(args, buildArgs...)

	}
	buffer.WriteByte(')')
	return buffer.String(), args
}

type Equal struct {
	Field string
	Value interface{}
}

func (e Equal) Build() (sql string, args []interface{}) {
	var buffer strings.Builder

	buffer.Grow(len(e.Field) + 2)

	buffer.WriteString(e.Field)
	buffer.WriteString("=?")

	return buffer.String(), []interface{}{e.Value}
}

type In struct {
	Field string `json:"field"`
	Value Values `json:"value"`
}

func (i In) Build() (sql string, args []interface{}) {
	if len(i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(i.Field) + 6 + (len(i.Value)-1)*2)
	args = make([]interface{}, len(i.Value), len(i.Value))

	buffer.WriteString(i.Field)
	buffer.WriteString(" IN(?")
	args[0] = i.Value[0]
	for index := 1; index < len(i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Between struct {
	Field string      `json:"field"`
	Start interface{} `json:"start"`
	End   interface{} `json:"end"`
}

func (b Between) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(b.Field) + 16)

	buffer.WriteString(b.Field)
	buffer.WriteString(" BETWEEN ? AND ?")

	return buffer.String(), []interface{}{b.Start, b.End}
}

type Contains struct {
	Field string `json:"field"`
	Words string `json:"words"`
}

func (c Contains) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(c.Field) + 7)

	buffer.WriteString(c.Field)
	buffer.WriteString(" LIKE ?")

	return buffer.String(), []interface{}{"%" + c.Words + "%"}
}

type BeginWith struct {
	Field string `json:"field"`
	Words string `json:"words"`
}

func (b BeginWith) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(b.Field) + 7)

	buffer.WriteString(b.Field)
	buffer.WriteString(" LIKE ?")

	return buffer.String(), []interface{}{b.Words + "%"}
}

type Gt struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (g Gt) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(g.Field) + 2)

	buffer.WriteString(g.Field)
	buffer.WriteString(">?")

	return buffer.String(), []interface{}{g.Value}
}

type Gte struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (g Gte) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(g.Field) + 3)

	buffer.WriteString(g.Field)
	buffer.WriteString(">=?")

	return buffer.String(), []interface{}{g.Value}
}

type Lt struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (l Lt) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(l.Field) + 2)

	buffer.WriteString(l.Field)
	buffer.WriteString("<?")

	return buffer.String(), []interface{}{l.Value}
}

type Lte struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (l Lte) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(l.Field) + 3)

	buffer.WriteString(l.Field)
	buffer.WriteString("<=?")

	return buffer.String(), []interface{}{l.Value}
}
