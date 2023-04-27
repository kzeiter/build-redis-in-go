package main

import (
	"fmt"
	"strconv"
)

type store struct {
	data map[string]interface{}
}

func (s *store) set(key string, value interface{}) {
	s.data[key] = value
}

func (s *store) get(key string) interface{} {
	return s.data[key]
}

func (s *store) del(key string) {
	delete(s.data, key)
}

func (s *store) setString(key, value string) {
	s.set(key, value)
}

func (s *store) getString(key string) string {
	value, _ := s.get(key).(string)
	return value
}

func (s *store) setNumber(key string, value string) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return
	}
	s.set(key, number)
}

func (s *store) getNumber(key string) int {
	value, _ := s.get(key).(int)
	return value
}

func (s *store) handleCommand(command string, args []string) string {
	switch command {
	case "SET":
		s.set(args[0], args[1])
		return "OK"
	case "GET":
		return fmt.Sprintf("%v", s.get(args[0]))
	case "DEL":
		s.del(args[0])
		return "OK"
	case "SETSTR":
		s.setString(args[0], args[1])
		return "OK"
	case "GETSTR":
		return s.getString(args[0])
	case "SETNUM":
		s.setNumber(args[0], args[1])
		return "OK"
	case "GETNUM":
		return fmt.Sprintf("%d", s.getNumber(args[0]))
	default:
		return "ERROR"
	}
}
