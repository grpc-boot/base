package result

import (
	"errors"

	"github.com/grpc-boot/base/v2/http_client"
)

type Sql struct {
	Result

	Columns []Column `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Cursor  string   `json:"cursor"`
}

func (s *Sql) ToRecord() ([]Source, error) {
	if s.HasError() {
		return nil, errors.New(s.Error.Reason)
	}

	if len(s.Columns) == 0 || len(s.Rows) == 0 {
		return make([]Source, 0), nil
	}

	res := make([]Source, len(s.Rows))
	for index, row := range s.Rows {
		src := make(Source, len(row))
		for i, _ := range row {
			src[s.Columns[i].Name] = row[i]
		}
		res[index] = src
	}

	return res, nil
}

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func ToSql(resp *http_client.Response, err error) (*Sql, error) {
	if err != nil {
		return nil, err
	}

	res := &Sql{}
	res.Status = resp.GetStatus()
	err = resp.UnmarshalJson(res)
	if err == nil && res.HasError() {
		return nil, errors.New(res.Error.Reason)
	}
	return res, err
}
