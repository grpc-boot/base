package redistack

import (
	"context"
	"time"

	"github.com/redis/rueidis"
)

func StreamAdd(timeout time.Duration, client rueidis.Client, key, maxLen string, fv FieldValue) (id string, err error) {
	fieldValue := client.B().Xadd().Key(key).Maxlen().Almost().Threshold(maxLen).Id(`*`).FieldValue()
	for field, value := range fv {
		fieldValue = fieldValue.FieldValue(field, value)
	}

	return DoWithTimeout(timeout, client, fieldValue.Build()).ToString()
}

func StreamGroupCreate(timeout time.Duration, client rueidis.Client, key, group, id string) (isNew bool, err error) {
	cmd := client.B().XgroupCreate().Key(key).Group(group).Id(id).Mkstream().Build()
	res, err := DoWithTimeout(timeout, client, cmd).ToString()
	if err == nil {
		isNew = res == OK
		return
	}

	if ErrBusyGroup(err) {
		err = nil
		return
	}

	return
}

func StreamReadGroup(blockMilliSeconds int64, client rueidis.Client, key, group, consumer string, maxCount int64, id string, noAck bool) (msgMap map[string][]rueidis.XRangeEntry, err error) {
	readBlock := client.B().Xreadgroup().Group(group, consumer).Count(maxCount).Block(blockMilliSeconds)
	if noAck {
		return Do(context.Background(), client, readBlock.Noack().Streams().Key(key).Id(id).Build()).AsXRead()
	}

	return Do(context.Background(), client, readBlock.Streams().Key(key).Id(id).Build()).AsXRead()
}

func StreamAck(timeout time.Duration, client rueidis.Client, key, group string, idList ...string) (okCount int64, err error) {
	cmd := client.B().Xack().Key(key).Group(group).Id(idList...).Build()
	return DoWithTimeout(timeout, client, cmd).AsInt64()
}

func StreamGroupCreateConsumer(timeout time.Duration, client rueidis.Client, key, group, consumer string) (isNew bool, err error) {
	cmd := client.B().XgroupCreateconsumer().Key(key).Group(group).Consumer(consumer).Build()
	res, err := DoWithTimeout(timeout, client, cmd).AsInt64()
	isNew = res == 1
	return
}

func StreamLen(timeout time.Duration, client rueidis.Client, key string) (length int64, err error) {
	cmd := client.B().Xlen().Key(key).Build()
	return DoWithTimeout(timeout, client, cmd).AsInt64()
}

func StreamInfo(timeout time.Duration, client rueidis.Client, key string) (info InfoStream, err error) {
	cmd := client.B().XinfoStream().Key(key).Full().Build()
	msgMap, err := DoWithTimeout(timeout, client, cmd).AsMap()
	if err != nil {
		return
	}
	return ToInfoStream(msgMap), nil
}
