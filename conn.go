package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn, s *Store) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Split(input, " ")

		if len(parts) < 2 || (parts[0] == "PUBLISH" && len(parts) < 3) {
			fmt.Fprintln(conn, "ERROR")
			continue
		}

		command := parts[0]
		args := parts[1:]

		response := s.handleCommand(command, args, conn)

		fmt.Fprintln(conn, response)
	}
}
