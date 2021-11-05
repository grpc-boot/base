package base

type CanHash interface {
	HashCode() (hashValue uint32)
}
