package iprange

import (
	"net"

	"github.com/russtone/utils/iter"
)

type iterator struct {
	rr  []IPRange
	idx int
	cur net.IP
}

func NewIterator(rr ...IPRange) iter.Iterator {
	return &iterator{
		rr:  rr,
		idx: 0,
		cur: nil,
	}
}

func (it *iterator) Next(out *string) bool {
	for i := it.idx; i < len(it.rr); i++ {
		it.cur = it.rr[i].next(it.cur)
		if it.cur != nil {
			break
		}

		it.idx++
	}

	if it.cur != nil {
		*out = it.cur.String()
		return true
	}

	return false
}

func (it *iterator) Count() uint64 {
	count := uint64(0)

	for _, r := range it.rr {
		count += r.Count()
	}

	return count
}

func (it *iterator) Reset() {
	it.cur = nil
	it.idx = 0
}

func (it *iterator) Close() error {
	return nil
}
