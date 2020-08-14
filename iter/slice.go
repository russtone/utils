package iter

type slice struct {
	items []string
	idx   int
}

func Slice(items []string) Iterator {
	return &slice{items, 0}
}

func (it *slice) Next(dest *string) bool {
	if it.idx < len(it.items) {
		*dest = it.items[it.idx]
		it.idx++
		return true
	}

	return false
}

func (it *slice) Reset() {
	it.idx = 0
}

func (it *slice) Count() uint64 {
	return uint64(len(it.items))
}

func (it *slice) Close() error {
	return nil
}
