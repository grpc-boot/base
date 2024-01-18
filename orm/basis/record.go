package basis

import (
	"database/sql"
	"reflect"
	"strconv"

	"github.com/grpc-boot/base/v2/utils"
)

type Record map[string]string

func (r Record) Exists(key string) bool {
	_, exists := r[key]
	return exists
}

func (r Record) Bytes(key string) []byte {
	value, _ := r[key]
	return []byte(value)
}

func (r Record) String(key string) string {
	value, _ := r[key]
	return value
}

func (r Record) Bool(key string) bool {
	value, _ := r[key]
	if len(value) == 0 {
		return false
	}

	val, _ := strconv.ParseBool(value)
	return val
}

func (r Record) Float64(key string) float64 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseFloat(value, 64)
	return val
}

func (r Record) Float32(key string) float32 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseFloat(value, 32)
	return float32(val)
}

func (r Record) Int64(key string) int64 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseInt(value, 10, 64)
	return val
}

func (r Record) Uint64(key string) uint64 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseUint(value, 10, 64)
	return val
}

func (r Record) Int(key string) int {
	return int(r.Int64(key))
}

func (r Record) Uint(key string) uint {
	return uint(r.Uint64(key))
}

func (r Record) Int32(key string) int32 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseInt(value, 10, 32)
	return int32(val)
}

func (r Record) Uint32(key string) uint32 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseUint(value, 10, 32)
	return uint32(val)
}

func (r Record) Int16(key string) int16 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseInt(value, 10, 16)
	return int16(val)
}

func (r Record) Uint16(key string) uint16 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseUint(value, 10, 16)
	return uint16(val)
}

func (r Record) Int8(key string) int8 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseInt(value, 10, 8)
	return int8(val)
}

func (r Record) Uint8(key string) uint8 {
	value, _ := r[key]
	if len(value) == 0 {
		return 0
	}

	val, _ := strconv.ParseUint(value, 10, 8)
	return uint8(val)
}

func (r Record) Convert(out any) (err error) {
	mp := ParseMapping(out)

	value := reflect.ValueOf(out).Elem()
	for index, _ := range mp.names {
		switch mp.kinds[index] {
		case reflect.String:
			value.Field(index).SetString(r.String(mp.fields[index]))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value.Field(index).SetInt(r.Int64(mp.fields[index]))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value.Field(index).SetUint(r.Uint64(mp.fields[index]))
		case reflect.Float32, reflect.Float64:
			value.Field(index).SetFloat(r.Float64(mp.fields[index]))
		case reflect.Bool:
			value.Field(index).SetBool(r.Bool(mp.fields[index]))
		case reflect.Slice:
			value.Field(index).SetBytes(r.Bytes(mp.fields[index]))
		default:
			continue
		}
	}

	return nil
}

func ScanRecords(rows *sql.Rows) (records []Record, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	values := make([]any, len(fields), len(fields))
	for index, _ := range fields {
		values[index] = &[]byte{}
	}

	records = []Record{}
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		record := make(map[string]string, len(fields))
		for index, field := range fields {
			record[field] = utils.Bytes2String(*values[index].(*[]byte))
		}

		records = append(records, record)
	}
	return
}
