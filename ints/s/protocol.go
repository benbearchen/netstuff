package s

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ReqInts struct {
	Offset int32
	Count  int32
	Slice  int32
}

type ReplyInts struct {
	Total  int32
	Offset int32
	Count  int32
	Ints   []int32
}

func ParseReqInts(pkg []byte) (*ReqInts, error) {
	buf := bytes.NewBuffer(pkg)
	var offset, count, slice int32
	if binary.Read(buf, binary.BigEndian, &offset) != nil || binary.Read(buf, binary.BigEndian, &count) != nil || binary.Read(buf, binary.BigEndian, &slice) != nil {
		return nil, fmt.Errorf("ParseReqInts error")
	}

	req := &ReqInts{offset, count, slice}
	return req, nil
}

func (req *ReqInts) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 64*1024))
	binary.Write(buf, binary.BigEndian, req.Offset)
	binary.Write(buf, binary.BigEndian, req.Count)
	binary.Write(buf, binary.BigEndian, req.Slice)
	return buf.Bytes()
}

func ParseReplyInts(pkg []byte) (*ReplyInts, error) {
	buf := bytes.NewBuffer(pkg)
	var total, offset, count int32
	if binary.Read(buf, binary.BigEndian, &total) != nil || binary.Read(buf, binary.BigEndian, &offset) != nil || binary.Read(buf, binary.BigEndian, &count) != nil {
		return nil, fmt.Errorf("ParseReplyInts error")
	}

	ints := make([]int32, 0, count)
	for i := int32(0); i < count; i++ {
		var v int32
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			return nil, err
		} else {
			ints = append(ints, v)
		}
	}

	reply := &ReplyInts{total, offset, count, ints}
	return reply, nil
}

func (reply *ReplyInts) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 64*1024))
	binary.Write(buf, binary.BigEndian, reply.Total)
	binary.Write(buf, binary.BigEndian, reply.Offset)
	binary.Write(buf, binary.BigEndian, reply.Count)

	for _, v := range reply.Ints {
		binary.Write(buf, binary.BigEndian, v)
	}

	return buf.Bytes()
}
