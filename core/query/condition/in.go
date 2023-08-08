package condition

import (
	"strings"
)

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

type StringIn struct {
	Field string   `json:"field"`
	Value []string `json:"value"`
}

func (si StringIn) Build() (sql string, args []interface{}) {
	if len(si.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(si.Field) + 6 + (len(si.Value)-1)*2)
	args = make([]interface{}, len(si.Value), len(si.Value))

	buffer.WriteString(si.Field)
	buffer.WriteString(" IN(?")
	args[0] = si.Value[0]
	for index := 1; index < len(si.Value); index++ {
		buffer.WriteString(",?")
		args[index] = si.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type IntIn struct {
	Field string `json:"field"`
	Value []int  `json:"value"`
}

func (ii IntIn) Build() (sql string, args []interface{}) {
	if len(ii.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(ii.Field) + 6 + (len(ii.Value)-1)*2)
	args = make([]interface{}, len(ii.Value), len(ii.Value))

	buffer.WriteString(ii.Field)
	buffer.WriteString(" IN(?")
	args[0] = ii.Value[0]
	for index := 1; index < len(ii.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ii.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type UintIn struct {
	Field string `json:"field"`
	Value []uint `json:"value"`
}

func (ui UintIn) Build() (sql string, args []interface{}) {
	if len(ui.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(ui.Field) + 6 + (len(ui.Value)-1)*2)
	args = make([]interface{}, len(ui.Value), len(ui.Value))

	buffer.WriteString(ui.Field)
	buffer.WriteString(" IN(?")
	args[0] = ui.Value[0]
	for index := 1; index < len(ui.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ui.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Int64In struct {
	Field string  `json:"field"`
	Value []int64 `json:"value"`
}

func (ii Int64In) Build() (sql string, args []interface{}) {
	if len(ii.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(ii.Field) + 6 + (len(ii.Value)-1)*2)
	args = make([]interface{}, len(ii.Value), len(ii.Value))

	buffer.WriteString(ii.Field)
	buffer.WriteString(" IN(?")
	args[0] = ii.Value[0]
	for index := 1; index < len(ii.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ii.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Uint64In struct {
	Field string   `json:"field"`
	Value []uint64 `json:"value"`
}

func (ui Uint64In) Build() (sql string, args []interface{}) {
	if len(ui.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(ui.Field) + 6 + (len(ui.Value)-1)*2)
	args = make([]interface{}, len(ui.Value), len(ui.Value))

	buffer.WriteString(ui.Field)
	buffer.WriteString(" IN(?")
	args[0] = ui.Value[0]
	for index := 1; index < len(ui.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ui.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Int32In struct {
	Field string  `json:"field"`
	Value []int32 `json:"value"`
}

func (i32i Int32In) Build() (sql string, args []interface{}) {
	if len(i32i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(i32i.Field) + 6 + (len(i32i.Value)-1)*2)
	args = make([]interface{}, len(i32i.Value), len(i32i.Value))

	buffer.WriteString(i32i.Field)
	buffer.WriteString(" IN(?")
	args[0] = i32i.Value[0]
	for index := 1; index < len(i32i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = i32i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Uint32In struct {
	Field string   `json:"field"`
	Value []uint32 `json:"value"`
}

func (u32i Uint32In) Build() (sql string, args []interface{}) {
	if len(u32i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(u32i.Field) + 6 + (len(u32i.Value)-1)*2)
	args = make([]interface{}, len(u32i.Value), len(u32i.Value))

	buffer.WriteString(u32i.Field)
	buffer.WriteString(" IN(?")
	args[0] = u32i.Value[0]
	for index := 1; index < len(u32i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = u32i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Int16In struct {
	Field string  `json:"field"`
	Value []int16 `json:"value"`
}

func (i16i Int16In) Build() (sql string, args []interface{}) {
	if len(i16i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(i16i.Field) + 6 + (len(i16i.Value)-1)*2)
	args = make([]interface{}, len(i16i.Value), len(i16i.Value))

	buffer.WriteString(i16i.Field)
	buffer.WriteString(" IN(?")
	args[0] = i16i.Value[0]
	for index := 1; index < len(i16i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = i16i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Uint16In struct {
	Field string   `json:"field"`
	Value []uint16 `json:"value"`
}

func (u16i Uint16In) Build() (sql string, args []interface{}) {
	if len(u16i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(u16i.Field) + 6 + (len(u16i.Value)-1)*2)
	args = make([]interface{}, len(u16i.Value), len(u16i.Value))

	buffer.WriteString(u16i.Field)
	buffer.WriteString(" IN(?")
	args[0] = u16i.Value[0]
	for index := 1; index < len(u16i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = u16i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Int8In struct {
	Field string `json:"field"`
	Value []int8 `json:"value"`
}

func (i8i Int8In) Build() (sql string, args []interface{}) {
	if len(i8i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(i8i.Field) + 6 + (len(i8i.Value)-1)*2)
	args = make([]interface{}, len(i8i.Value), len(i8i.Value))

	buffer.WriteString(i8i.Field)
	buffer.WriteString(" IN(?")
	args[0] = i8i.Value[0]
	for index := 1; index < len(i8i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = i8i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type Uint8In struct {
	Field string  `json:"field"`
	Value []uint8 `json:"value"`
}

func (u8i Uint8In) Build() (sql string, args []interface{}) {
	if len(u8i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(u8i.Field) + 6 + (len(u8i.Value)-1)*2)
	args = make([]interface{}, len(u8i.Value), len(u8i.Value))

	buffer.WriteString(u8i.Field)
	buffer.WriteString(" IN(?")
	args[0] = u8i.Value[0]
	for index := 1; index < len(u8i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = u8i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}
