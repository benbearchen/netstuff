package s

import (
	"sync"
)

type Processor struct {
	ints  *Ints
	mutex sync.RWMutex
}

func NewProcessor() *Processor {
	return new(Processor)
}

func (proc *Processor) Do(c ServerContext) {
	req := c.Get()
	reply := proc.reply(req)
	if len(reply) > 0 {
		for _, v := range reply {
			c.Reply(v)
		}
	}
}

func (proc *Processor) reply(req *ReqInts) []*ReplyInts {
	ints := proc.getInts()
	if ints == nil {
		return []*ReplyInts{&ReplyInts{0, 0, 0, make([]int32, 0)}}
	}

	var selected []int32 = nil
	if req.Offset >= 0 && req.Count > 0 && req.Slice > 0 {
		selected = ints.GetSlice(int(req.Offset), int(req.Count))
	}

	if selected == nil {
		return []*ReplyInts{&ReplyInts{int32(ints.Count()), 0, 0, make([]int32, 0)}}
	}

	slices := make([]*ReplyInts, 0, len(selected)/int(req.Slice)+1)
	for i := 0; i < len(selected); {
		c := int(req.Slice)
		if c > len(selected)-i {
			c = len(selected) - i
		}
		offset := int(req.Offset) + i
		reply := &ReplyInts{int32(ints.Count()), int32(offset), int32(c), selected[i : i+c]}
		slices = append(slices, reply)

		i += c
	}

	return slices
}

func (proc *Processor) getInts() *Ints {
	proc.mutex.RLock()
	defer proc.mutex.RUnlock()
	return proc.ints
}

func (proc *Processor) CreateInts(n int) {
	ints := NewInts()
	ints.Gen(n)

	proc.mutex.Lock()
	defer proc.mutex.Unlock()
	proc.ints = ints
}
