package main

import (
	"net"
)

func main() {
	s, _ := NewStore("data.out")

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
