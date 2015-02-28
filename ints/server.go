package main

import (
	"github.com/benbearchen/netstuff/ints/s"

	"fmt"
)

func main() {
	proc := s.NewProcessor()
	s := s.NewServer()
	port := 2015
	addr := fmt.Sprintf(":%d", port)
	s.Bind(addr, proc)
	fmt.Println("bind on: ", addr)
	s.Run()
}
