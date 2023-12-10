package kind

type RedisValue interface {
	Number | ~bool | ~string | ~[]byte
}
