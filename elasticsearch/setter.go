package elasticsearch

import (
	"bytes"
)

type Setter map[string]any

func (s Setter) ToBody() (setter []byte, err error) {
	if len(s) == 0 {
		return nil, ErrSetterEmpty
	}

	var buf bytes.Buffer
	for field, _ := range s {
		if buf.Len() > 0 {
			buf.WriteByte(';')
		}
		buf.WriteString("ctx._source.")
		buf.WriteString(field)
		buf.WriteString(`=params.`)
		buf.WriteString(field)
	}

	return Body{
		"script": Property{
			"source": buf.String(),
			"lang":   "painless",
			"params": s,
		},
	}.Marshal()
}
