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
	return zap.Duration("Dration", duration)
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
