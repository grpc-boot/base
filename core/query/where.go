package query

import "strings"

type whereItem struct {
	opt       string
	condition Condition
}

// Where Where对象
type Where interface {
	HasWhere() bool
	And(condition Condition) Where
	Or(condition Condition) Where
	Sql(args *[]interface{}) (sql string)
	Reset()
}

type where struct {
	items []whereItem
}

// AndWhere 实例化And条件Where
func AndWhere(fm FieldMap) Where {
	return &where{items: []whereItem{{
		condition: AndCondition(fm),
	}}}
}

// OrWhere 实例化Or条件Where
func OrWhere(fm FieldMap) Where {
	return &where{items: []whereItem{{
		condition: OrCondition(fm),
	}}}
}

// NewWhere 实例化Where
func NewWhere(condition Condition) Where {
	return &where{items: []whereItem{{
		opt:       And,
		condition: condition,
	}}}
}

// HasWhere 是否有where条件
func (w *where) HasWhere() bool {
	return len(w.items) > 0
}

// And 附加And条件
func (w *where) And(condition Condition) Where {
	w.items = append(w.items, whereItem{opt: And, condition: condition})
	return w
}

// Or 附加Or条件
func (w *where) Or(condition Condition) Where {
	w.items = append(w.items, whereItem{opt: Or, condition: condition})
	return w
}

// Sql 生成where sql
func (w *where) Sql(args *[]interface{}) (sql string) {
	if !w.HasWhere() {
		return
	}

	var (
		buf strings.Builder
	)

	for index, wc := range w.items {
		sqlStr := wc.condition.Sql(args)
		if sqlStr == "" {
			continue
		}

		if index == 0 {
			buf.WriteString(` WHERE `)
		} else {
			buf.WriteByte(' ')
			buf.WriteString(wc.opt)
			buf.WriteByte(' ')
		}
		buf.WriteString(sqlStr)
	}

	return buf.String()
}

// Reset ---
func (w *where) Reset() {
	if w.items == nil {
		w.items = []whereItem{}
		return
	}

	if w.HasWhere() {
		w.items = w.items[:0]
	}
}
