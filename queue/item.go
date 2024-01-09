package queue

import (
	"bytes"
	"strconv"
)

type Item struct {
	Id         string `json:"i"`
	Name       string `json:"n"`
	At         int64  `json:"a"`
	RetryCount int64  `json:"rc"`
}

func (i *Item) Member() string {
	var buf = bytes.NewBuffer(nil)

	at := strconv.FormatInt(i.At, 10)
	retryCount := strconv.FormatInt(i.RetryCount, 10)
	buf.Grow(len(i.Id) + len(i.Name) + len(at) + len(retryCount) + 6 + 7 + 6 + 6 + 1)
	buf.WriteString(`{"i":"`)
	buf.WriteString(i.Id)
	buf.WriteString(`","n":"`)
	buf.WriteString(i.Name)
	buf.WriteString(`","a":`)
	buf.WriteString(at)
	buf.WriteString(`,"rc":`)
	buf.WriteString(retryCount)
	buf.WriteByte('}')

	return buf.String()
}
