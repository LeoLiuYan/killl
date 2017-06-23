package ring

import (
	"github.com/errors"
	"killl/lib/log"
)

type Ring struct {
	rp   uint32
	wp   uint32
	num  uint32
	mask uint32
	data []Payload
}

type Payload struct {
	Name string
	ID   string
}

var (
	ErrRingEmpty = errors.New("ring buffer empty")
	ErrRingFull  = errors.New("ring buffer full")
)

func NewRing(num uint32) (r *Ring) {
	r = new(Ring)
	if num&(num-1) != 0 {
		for num&(num-1) != 0 {
			num &= num - 1
		}
		num = num << 1
	}
	r.num = num
	r.data = make([]Payload, num)
	r.mask = r.num - 1
	r.wp = 0
	r.rp = 0

	return
}

func (r *Ring) Set() (*Payload, error) {
	if r.wp-r.rp >= r.num {
		return nil, ErrRingFull
	}
	return &r.data[r.wp&r.mask], nil
}

func (r *Ring) SetA() {
	r.wp++
	log.Debugf("ring wp: %d, index: %d", r.wp, r.wp&r.mask)
}

func (r *Ring) Get() (*Payload, error) {
	if r.rp == r.wp {
		return nil, ErrRingEmpty
	}
	return &r.data[r.rp&r.mask], nil
}

func (r *Ring) GetA() {
	r.rp++
	log.Debugf("ring rp: %d, index; %d", r.rp, r.rp&r.mask)
}
