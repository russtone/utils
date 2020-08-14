package iprange

import (
	"net"
)

type Iterator interface {
	Next(net.IP) bool
	Count() uint64
	Reset()
}

type iterator struct {
	rr  []IPRange
	idx int
	cur net.IP
}

func NewIterator(rr ...IPRange) Iterator {
	return &iterator{
		rr:  rr,
		idx: 0,
		cur: nil,
	}
}

func (it *iterator) Next(out net.IP) bool {
	if out.To16() != nil {
		out = out.To4()
	}

	for i := it.idx; i < len(it.rr); i++ {
		it.cur = it.rr[i].next(it.cur)
		if it.cur != nil {
			break
		}

		it.idx++
	}

	if it.cur != nil {
		copy(out, it.cur)
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
