package main

import "fmt"

func main() {
	s := store{data: make(map[string]string)}
	s.set("foo", "bar")
	fmt.Println(s.get("foo"))
}
