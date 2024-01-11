package basis

type Tx interface {
	Commit() error
	Rollback() error
}
