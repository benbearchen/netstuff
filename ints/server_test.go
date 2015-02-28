package main

import (
	"testing"
)

import (
	"github.com/benbearchen/netstuff/ints/s"

	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func TestMain(t *testing.T) {
	proc := s.NewProcessor()
	s2015 := s.NewServer()
	addr2015, _ := net.ResolveUDPAddr("udp", "127.0.0.1:2015")
	s2015.Bind(addr2015.String(), proc)
	fmt.Println("bind on: ", addr2015)
	s2015.Run()

	proc.CreateInts(1000000)

	addr, err := net.ResolveUDPAddr("udp", ":2016")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	var w sync.WaitGroup
	w.Add(2)

	ch := make(chan *s.ReplyInts)
	stop := false

	go func() {
		defer w.Done()
		defer close(ch)

		for !stop {
			buf := make([]byte, 64*1024)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				if e, ok := err.(net.Error); ok && e.Timeout() {
					continue
				} else {
					t.Errorf("recv err: %v", err)
				}
			}

			reply, err := s.ParseReplyInts(buf[:n])
			if err == nil {
				ch <- reply
			} else {
				t.Errorf("parse err: %v", err)
			}
		}
	}()

	go func() {
		defer w.Done()
		defer func() { stop = true }()

		var nums []int32 = nil
		next := 0

		for !stop && (nums == nil || next < len(nums)) {
			req := &s.ReqInts{int32(next), 300, 300}
			conn.WriteToUDP(req.Bytes(), addr2015)

			select {
			case reply, ok := <-ch:
				if !ok {
					stop = true
				} else {
					if nums == nil {
						nums = make([]int32, reply.Total)
					} else if int(reply.Total) != len(nums) {
						stop = true
						continue
					}

					if next == int(reply.Offset) && int(reply.Offset+reply.Count) <= len(nums) {
						copy(nums[int(reply.Offset):int32(reply.Offset+reply.Count)], reply.Ints)
						next += int(reply.Count)
					}
				}
			case <-time.NewTimer(time.Millisecond * 10).C:
			}
		}

		if nums != nil {
			var sum int64 = 0
			t := make([]int, 16)
			ints := make([]int, len(nums))
			for i, v := range nums {
				sum += int64(v)
				t[v % 16]++
				ints[i] = int(v)
			}

			fmt.Println("sum:", sum, "avg:", float64(sum) / float64(len(nums)))
			fmt.Println("t[%16]:", t)

			sort.Ints(ints)
			if len(ints) % 2 == 0 {
				mid := ints[len(ints)/2-1:len(ints)/2+1]
				fmt.Println("mid:", (mid[0] + mid[1]) / 2)
			} else {
				fmt.Println("mid:", ints[len(ints)/2])
			}
		}
	}()

	w.Wait()
}
