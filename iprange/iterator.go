package iprange

import (
	"net"

	"github.com/russtone/utils/iter"
)

// iterator is an internal iterator state.
type iterator struct {
	// IP ranges
	rr []IterableRange

	// Index of current IP range.
	idx int

	// Current IP.
	cur net.IP
}

// NewIterator creates new iterator from IP ranges.
// Iterator implements iter.Iterator interface.
func NewIterator(rr ...IterableRange) iter.Iterator {
	return &iterator{
		rr:  rr,
		idx: 0,
		cur: nil,
	}
}

// Next returns next IP address in range as string.
// strings used instead of net.IP to implement iter.Iterator interface.
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

// Count returns total number of IP addreses in the iterator.
func (it *iterator) Count() uint64 {
	count := uint64(0)

	for _, r := range it.rr {
		count += r.Count()
	}

	return count
}

// Reset rewinds the iterator.
func (it *iterator) Reset() {
	it.cur = nil
	it.idx = 0
}

// Close do nothing.
// Required to satisfy iter.Iterator interface.
func (it *iterator) Close() error {
	return nil
}
