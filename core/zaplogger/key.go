package zaplogger

import (
	"time"

	"go.uber.org/zap"
)

func Id(id int64) zap.Field {
	return zap.Int64("Id", id)
}

func UpdateAt(updatedAt int64) zap.Field {
	return zap.Int64("UpdatedAt", updatedAt)
}

func Ip(ip string) zap.Field {
	return zap.String("Ip", ip)
}

func TraceId(traceId string) zap.Field {
	return zap.String("TraceId", traceId)
}

func Mid(mid string) zap.Field {
	return zap.String("Mid", mid)
}

func Cmd(cmd string) zap.Field {
	return zap.String("Cmd", cmd)
}

func Event(name string) zap.Field {
	return zap.String("Event", name)
}

func Method(method string) zap.Field {
	return zap.String("Method", method)
}

func Path(path string) zap.Field {
	return zap.String("Path", path)
}

func Uri(uri string) zap.Field {
	return zap.String("Uri", uri)
}

func Duration(duration time.Duration) zap.Field {
	return zap.Duration("Duration", duration)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Args(args ...interface{}) zap.Field {
	return zap.Any("Args", args)
}

func Params(params ...interface{}) zap.Field {
	return zap.Any("Params", params)
}

func Value(value interface{}) zap.Field {
	return zap.Any("Value", value)
}

func Data(data []byte) zap.Field {
	return zap.ByteString("Data", data)
}

func Addr(addr string) zap.Field {
	return zap.String("Addr", addr)
}

func Offset(oft string) zap.Field {
	return zap.String("Offset", oft)
}

func Size(size int64) zap.Field {
	return zap.Int64("Size", size)
}

func Key(key string) zap.Field {
	return zap.String("Key", key)
}

func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Strings(key string, value []string) zap.Field {
	return zap.Strings(key, value)
}

func BytesString(key string, value []byte) zap.Field {
	return zap.ByteString(key, value)
}

func BytesStrings(key string, value [][]byte) zap.Field {
	return zap.ByteStrings(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Ints(key string, value []int) zap.Field {
	return zap.Ints(key, value)
}

func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

func Int64s(key string, value []int64) zap.Field {
	return zap.Int64s(key, value)
}

func Int32(key string, value int32) zap.Field {
	return zap.Int32(key, value)
}

func Int32s(key string, value []int32) zap.Field {
	return zap.Int32s(key, value)
}

func Int16(key string, value int16) zap.Field {
	return zap.Int16(key, value)
}

func Int16s(key string, value []int16) zap.Field {
	return zap.Int16s(key, value)
}

func Int8(key string, value int8) zap.Field {
	return zap.Int8(key, value)
}

func Uint(key string, value uint) zap.Field {
	return zap.Uint(key, value)
}

func Uints(key string, value []uint) zap.Field {
	return zap.Uints(key, value)
}

func Uint64(key string, value uint64) zap.Field {
	return zap.Uint64(key, value)
}

func Uint64s(key string, value []uint64) zap.Field {
	return zap.Uint64s(key, value)
}

func Uint32(key string, value uint32) zap.Field {
	return zap.Uint32(key, value)
}

func Uint32s(key string, value []uint32) zap.Field {
	return zap.Uint32s(key, value)
}

func Uint16(key string, value uint16) zap.Field {
	return zap.Uint16(key, value)
}

func Uint16s(key string, value []uint16) zap.Field {
	return zap.Uint16s(key, value)
}

func Uint8(key string, value uint8) zap.Field {
	return zap.Uint8(key, value)
}

func Uint8s(key string, value []uint8) zap.Field {
	return zap.Uint8s(key, value)
}

func Time(key string, value time.Time) zap.Field {
	return zap.Time(key, value)
}

func Times(key string, value []time.Time) zap.Field {
	return zap.Times(key, value)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Binary(key string, value []byte) zap.Field {
	return zap.Binary(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func Bools(key string, value []bool) zap.Field {
	return zap.Bools(key, value)
}
