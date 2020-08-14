package iter

import "github.com/russtone/utils/errors"

type combined struct {
	items []Iterator
	idx   int
}

func Combine(iterators ...Iterator) Iterator {
	return &combined{iterators, 0}
}

func (it *combined) Reset() {
	for _, item := range it.items {
		item.Reset()
	}
}

func (it *combined) Next(dest *string) bool {
	for i := it.idx; i < len(it.items); i++ {
		if it.items[i].Next(dest) {
			return true
		}
	}

	return false
}

func (it *combined) Count() uint64 {
	count := uint64(0)
	for _, item := range it.items {
		count += item.Count()
	}
	return count
}

func (it *combined) Close() error {
	ee := make(errors.Errors, 0)

	for _, item := range it.items {
		if err := item.Close(); err != nil {
			ee = append(ee, err)
		}
	}

	if len(ee) > 0 {
		return ee
	}

	return nil
}
