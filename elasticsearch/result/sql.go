package result

import "github.com/grpc-boot/base/v2/orm/basis"

type Sql struct {
	Columns []Column `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Cursor  string   `json:"cursor"`
}

func (s *Sql) ToRecord() []basis.Record {
	if len(s.Columns) == 0 || len(s.Rows) == 0 {
		return make([]basis.Record, 0)
	}

	return make([]basis.Record, 0)
}

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
