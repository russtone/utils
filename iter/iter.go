package iter

type Iterator interface {
	Next(*string) bool
	Reset()
	Count() uint64
	Close() error
}
