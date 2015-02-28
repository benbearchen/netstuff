package s

import (
	"math/rand"
	"time"
)

type Ints struct {
	Ints []int32
}

func NewInts() *Ints {
	i := new(Ints)
	i.Ints = make([]int32, 0)
	return i
}

func genSpace(s *rand.Rand, max int32, spaces []int32) {
	if len(spaces) == 0 || max <= 0 {
		return
	} else if len(spaces) == 1 {
		spaces[0] = s.Int31n(max)
	} else {
		p1 := s.Int31n(max)
		c1 := len(spaces) / 2
		if p1 > 0 {
			genSpace(s, p1, spaces[:c1])
		}

		if p1 < max {
			genSpace(s, max-p1, spaces[c1:len(spaces)])
		}
	}
}

func gen(s *rand.Rand, min, max int32, count int) []int32 {
	spaces := make([]int32, count-1)
	genSpace(s, max-min-int32(count), spaces)
	ints := make([]int32, count)
	ints[0] = min
	for i := 1; i < count; i++ {
		ints[i] = ints[i-1] + 1 + spaces[i-1]
	}

	return ints
}

func shuffle(s *rand.Rand, ints []int32) {
	src := make([]int32, len(ints))
	copy(src, ints)
	for i, j := range s.Perm(len(ints)) {
		ints[i] = src[j]
	}
}

func (i *Ints) Gen(count int) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := int32(1)
	max := int32(0x7FFFFFFF)
	if count < 0x7FFFFFFF/10 {
		for {
			a := s.Int31n(0x7FFFFFFF)
			b := s.Int31n(a)
			if a-b < int32(count*2) {
				continue
			}

			min = b
			max = a
			break
		}
	}

	ints := gen(s, min, max, count)
	shuffle(s, ints)
	i.Ints = ints
}

func (i *Ints) Count() int {
	return len(i.Ints)
}

func (i *Ints) GetSlice(offset, count int) []int32 {
	c := i.Count()
	if offset >= c {
		return nil
	}

	if offset+count > c {
		count = c - offset
	}

	return i.Ints[offset : offset+count]
}
