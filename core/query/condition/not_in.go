package condition

import "strings"

type NotIn struct {
	Field string `json:"field"`
	Value Values `json:"value"`
}

func (ni NotIn) Build() (sql string, args []interface{}) {
	if len(ni.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(ni.Field) + 10 + (len(ni.Value)-1)*2)
	args = make([]interface{}, len(ni.Value), len(ni.Value))

	buffer.WriteString(ni.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = ni.Value[0]
	for index := 1; index < len(ni.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ni.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}
