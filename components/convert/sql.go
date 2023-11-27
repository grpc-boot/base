package convert

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/grpc-boot/base/v2/internal"
	"github.com/grpc-boot/base/v2/utils"
)

func SqlRows2RowList(rows *sql.Rows) (list []Row, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return
	}

	defer rows.Close()

	values := make([]any, len(fields), len(fields))
	for index, _ := range fields {
		values[index] = &[]byte{}
	}

	list = []Row{}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			list = nil
			return
		}

		row := make(Row, len(fields))
		for index, field := range fields {
			row[field] = internal.Bytes2String(*values[index].(*[]byte))
		}

		list = append(list, row)
	}

	return
}

func SqlRows2StructList(rows *sql.Rows, out any, tagName string) (list []any, err error) {
	var (
		mapFields   map[string]cacheMapType
		t           = reflect.TypeOf(out)
		val, exists = _mapCache.Load(t)
	)

	if !exists {
		if t.Kind() != reflect.Ptr {
			err = fmt.Errorf("obj must be a pointer")
			return
		}

		if t.Elem().Kind() != reflect.Struct {
			err = fmt.Errorf("obj must be a pointer to a struct")
			return
		}

		sliceFields := parseType(t.Elem(), tagName)
		mapFields = slice2Map(sliceFields)

		_mapCache.Store(t, mapFields)
		_cache.Store(t, sliceFields)
	} else {
		mapFields, _ = val.(map[string]cacheMapType)
	}

	fields, err := rows.Columns()
	if err != nil {
		return
	}

	defer rows.Close()

	values := make([]any, len(fields), len(fields))
	for index, _ := range fields {
		values[index] = &[]byte{}
	}

	list = []any{}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			list = nil
			return
		}

		value := reflect.New(t.Elem())

		for index, field := range fields {
			if ct, ok := mapFields[field]; ok {
				v := *values[index].(*[]byte)

				switch ct.kind {
				case reflect.String:
					value.Field(ct.index).SetString(internal.Bytes2String(v))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value.Field(ct.index).SetInt(utils.Bytes2Int64(v))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					value.Field(ct.index).SetUint(utils.Bytes2Uint64(v))
				case reflect.Float32, reflect.Float64:
					value.Field(ct.index).SetFloat(utils.Bytes2Float64(v))
				case reflect.Bool:
					value.Field(ct.index).SetBool(utils.Bytes2Bool(v))
				default:
					continue
				}
			}
		}

		list = append(list, value.Interface())
	}

	return
}
