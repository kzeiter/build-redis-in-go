package main

import (
	"net"
)

func main() {
	s := &store{
		data: make(map[string]interface{}),
		list: make(map[string][]string),
		sets: make(map[string]map[string]bool),
	}

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go HandleConnection(conn, s)
	}
}
