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

	if len(ni.Value) == 1 {
		return NotEqual{Field: ni.Field, Value: ni.Value[0]}.Build()
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

type NotStringIn struct {
	Field string   `json:"field"`
	Value []string `json:"value"`
}

func (nsi NotStringIn) Build() (sql string, args []interface{}) {
	if len(nsi.Value) == 0 {
		return
	}

	if len(nsi.Value) == 1 {
		return NotEqual{Field: nsi.Field, Value: nsi.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(nsi.Field) + 10 + (len(nsi.Value)-1)*2)
	args = make([]interface{}, len(nsi.Value), len(nsi.Value))

	buffer.WriteString(nsi.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = nsi.Value[0]
	for index := 1; index < len(nsi.Value); index++ {
		buffer.WriteString(",?")
		args[index] = nsi.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotIntIn struct {
	Field string `json:"field"`
	Value []int  `json:"value"`
}

func (nii NotIntIn) Build() (sql string, args []interface{}) {
	if len(nii.Value) == 0 {
		return
	}

	if len(nii.Value) == 1 {
		return NotEqual{Field: nii.Field, Value: nii.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(nii.Field) + 10 + (len(nii.Value)-1)*2)
	args = make([]interface{}, len(nii.Value), len(nii.Value))

	buffer.WriteString(nii.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = nii.Value[0]
	for index := 1; index < len(nii.Value); index++ {
		buffer.WriteString(",?")
		args[index] = nii.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotUintIn struct {
	Field string `json:"field"`
	Value []uint `json:"value"`
}

func (nui NotUintIn) Build() (sql string, args []interface{}) {
	if len(nui.Value) == 0 {
		return
	}

	if len(nui.Value) == 1 {
		return NotEqual{Field: nui.Field, Value: nui.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(nui.Field) + 10 + (len(nui.Value)-1)*2)
	args = make([]interface{}, len(nui.Value), len(nui.Value))

	buffer.WriteString(nui.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = nui.Value[0]
	for index := 1; index < len(nui.Value); index++ {
		buffer.WriteString(",?")
		args[index] = nui.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotInt64In struct {
	Field string  `json:"field"`
	Value []int64 `json:"value"`
}

func (ni64i NotInt64In) Build() (sql string, args []interface{}) {
	if len(ni64i.Value) == 0 {
		return
	}

	if len(ni64i.Value) == 1 {
		return NotEqual{Field: ni64i.Field, Value: ni64i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(ni64i.Field) + 10 + (len(ni64i.Value)-1)*2)
	args = make([]interface{}, len(ni64i.Value), len(ni64i.Value))

	buffer.WriteString(ni64i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = ni64i.Value[0]
	for index := 1; index < len(ni64i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ni64i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotUint64In struct {
	Field string   `json:"field"`
	Value []uint64 `json:"value"`
}

func (nu64i NotUint64In) Build() (sql string, args []interface{}) {
	if len(nu64i.Value) == 0 {
		return
	}

	if len(nu64i.Value) == 1 {
		return NotEqual{Field: nu64i.Field, Value: nu64i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(nu64i.Field) + 10 + (len(nu64i.Value)-1)*2)
	args = make([]interface{}, len(nu64i.Value), len(nu64i.Value))

	buffer.WriteString(nu64i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = nu64i.Value[0]
	for index := 1; index < len(nu64i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = nu64i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotInt32In struct {
	Field string  `json:"field"`
	Value []int32 `json:"value"`
}

func (ni32i NotInt32In) Build() (sql string, args []interface{}) {
	if len(ni32i.Value) == 0 {
		return
	}

	if len(ni32i.Value) == 1 {
		return NotEqual{Field: ni32i.Field, Value: ni32i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(ni32i.Field) + 10 + (len(ni32i.Value)-1)*2)
	args = make([]interface{}, len(ni32i.Value), len(ni32i.Value))

	buffer.WriteString(ni32i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = ni32i.Value[0]
	for index := 1; index < len(ni32i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ni32i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotUint32In struct {
	Field string   `json:"field"`
	Value []uint32 `json:"value"`
}

func (nu32i NotUint32In) Build() (sql string, args []interface{}) {
	if len(nu32i.Value) == 0 {
		return
	}

	if len(nu32i.Value) == 1 {
		return NotEqual{Field: nu32i.Field, Value: nu32i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(nu32i.Field) + 10 + (len(nu32i.Value)-1)*2)
	args = make([]interface{}, len(nu32i.Value), len(nu32i.Value))

	buffer.WriteString(nu32i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = nu32i.Value[0]
	for index := 1; index < len(nu32i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = nu32i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotInt16In struct {
	Field string  `json:"field"`
	Value []int16 `json:"value"`
}

func (ni16i NotInt16In) Build() (sql string, args []interface{}) {
	if len(ni16i.Value) == 0 {
		return
	}

	if len(ni16i.Value) == 1 {
		return NotEqual{Field: ni16i.Field, Value: ni16i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(ni16i.Field) + 10 + (len(ni16i.Value)-1)*2)
	args = make([]interface{}, len(ni16i.Value), len(ni16i.Value))

	buffer.WriteString(ni16i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = ni16i.Value[0]
	for index := 1; index < len(ni16i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ni16i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotUint16In struct {
	Field string   `json:"field"`
	Value []uint16 `json:"value"`
}

func (nu16i NotUint16In) Build() (sql string, args []interface{}) {
	if len(nu16i.Value) == 0 {
		return
	}

	if len(nu16i.Value) == 1 {
		return NotEqual{Field: nu16i.Field, Value: nu16i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(nu16i.Field) + 10 + (len(nu16i.Value)-1)*2)
	args = make([]interface{}, len(nu16i.Value), len(nu16i.Value))

	buffer.WriteString(nu16i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = nu16i.Value[0]
	for index := 1; index < len(nu16i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = nu16i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotInt8In struct {
	Field string `json:"field"`
	Value []int8 `json:"value"`
}

func (ni8i NotInt8In) Build() (sql string, args []interface{}) {
	if len(ni8i.Value) == 0 {
		return
	}

	if len(ni8i.Value) == 1 {
		return NotEqual{Field: ni8i.Field, Value: ni8i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(ni8i.Field) + 10 + (len(ni8i.Value)-1)*2)
	args = make([]interface{}, len(ni8i.Value), len(ni8i.Value))

	buffer.WriteString(ni8i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = ni8i.Value[0]
	for index := 1; index < len(ni8i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ni8i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}

type NotUint8In struct {
	Field string  `json:"field"`
	Value []uint8 `json:"value"`
}

func (nu8i NotUint8In) Build() (sql string, args []interface{}) {
	if len(nu8i.Value) == 0 {
		return
	}

	if len(nu8i.Value) == 1 {
		return NotEqual{Field: nu8i.Field, Value: nu8i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(nu8i.Field) + 10 + (len(nu8i.Value)-1)*2)
	args = make([]interface{}, len(nu8i.Value), len(nu8i.Value))

	buffer.WriteString(nu8i.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = nu8i.Value[0]
	for index := 1; index < len(nu8i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = nu8i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}
