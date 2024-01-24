package elasticsearch

import (
	"bytes"

	"github.com/grpc-boot/base/v2/utils"
)

type BulkMarshal interface {
	Marshal() []byte
}

type BulkItem struct {
	list []BulkMarshal
}

func NewBulkItem(items ...BulkMarshal) *BulkItem {
	bi := BulkItem{}
	if len(items) == 0 {
		bi.list = []BulkMarshal{}
	} else {
		bi.list = items
	}

	return &bi
}

func (bi *BulkItem) Marshal() []byte {
	if len(bi.list) == 0 {
		return nil
	}

	var buf = bytes.NewBuffer(nil)
	for _, bulk := range bi.list {
		buf.Write(bulk.Marshal())
	}
	return buf.Bytes()
}

func (bi *BulkItem) WithIndex(index, id string, property Property) *BulkItem {
	item := BulkIndex{}
	item.Index = index
	item.Id = id
	item.property = property
	bi.list = append(bi.list, &item)
	return bi
}

func (bi *BulkItem) WithCreate(index, id string, property Property) *BulkItem {
	item := BulkCreate{}
	item.Index = index
	item.Id = id
	item.property = property
	bi.list = append(bi.list, &item)
	return bi
}

func (bi *BulkItem) WithUpdate(index, id string, property Updater) *BulkItem {
	item := BulkUpdate{}
	item.Index = index
	item.Id = id
	item.property = property
	bi.list = append(bi.list, &item)
	return bi
}

func (bi *BulkItem) WithDelete(index, id string) *BulkItem {
	item := BulkDelete{}
	item.Index = index
	item.Id = id
	bi.list = append(bi.list, &item)
	return bi
}

type bulkCondition struct {
	Index string `json:"_index"`
	Id    string `json:"_id"`
}

type BulkIndex struct {
	bulkCondition
	property Property
}

func (bi *BulkIndex) Marshal() []byte {
	var buf = bytes.NewBuffer(nil)
	condition, _ := utils.JsonMarshal(bi.bulkCondition)
	propertied, _ := utils.JsonMarshal(bi.property)
	buf.Grow(9 + len(condition) + len(propertied) + 3)
	buf.WriteString(`{"index":`)
	buf.Write(condition)
	buf.WriteString("}\n")
	buf.Write(propertied)
	buf.WriteByte('\n')
	return buf.Bytes()
}

type BulkCreate struct {
	bulkCondition
	property Property
}

func (bc *BulkCreate) Marshal() []byte {
	var buf = bytes.NewBuffer(nil)
	condition, _ := utils.JsonMarshal(bc.bulkCondition)
	propertied, _ := utils.JsonMarshal(bc.property)
	buf.Grow(10 + len(condition) + len(propertied) + 3)
	buf.WriteString(`{"create":`)
	buf.Write(condition)
	buf.WriteString("}\n")
	buf.Write(propertied)
	buf.WriteByte('\n')
	return buf.Bytes()
}

type BulkUpdate struct {
	bulkCondition
	property Updater
}

type Updater struct {
	Doc         Property `json:"doc,omitempty"`
	Script      Property `json:"script,omitempty"`
	DocAsUpsert bool     `json:"doc_as_upsert,omitempty"`
	Source      bool     `json:"_source,omitempty"`
}

func (bu *BulkUpdate) Marshal() []byte {
	var buf = bytes.NewBuffer(nil)
	condition, _ := utils.JsonMarshal(bu.bulkCondition)
	properties, _ := utils.JsonMarshal(bu.property)
	buf.Grow(10 + len(condition) + len(properties) + 3)
	buf.WriteString(`{"update":`)
	buf.Write(condition)
	buf.WriteString("}\n")
	buf.Write(properties)
	buf.WriteByte('\n')
	return buf.Bytes()
}

type BulkDelete struct {
	bulkCondition
}

func (bc *BulkDelete) Marshal() []byte {
	var buf = bytes.NewBuffer(nil)
	condition, _ := utils.JsonMarshal(bc.bulkCondition)
	buf.Grow(10 + len(condition) + 2)
	buf.WriteString(`{"delete":`)
	buf.Write(condition)
	buf.WriteString("}\n")
	return buf.Bytes()
}
