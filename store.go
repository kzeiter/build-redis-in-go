package main

import (
	"fmt"
	"strconv"
)

type store struct {
	data map[string]interface{}
	list map[string][]string
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

func (s *store) lPush(key string, value string) int {
	if s.list[key] == nil {
		s.list[key] = make([]string, 0)
	}

	s.list[key] = append([]string{value}, s.list[key]...)

	return len(s.list[key])
}

func (s *store) rPush(key string, value string) int {
	if s.list[key] == nil {
		s.list[key] = make([]string, 0)
	}

	s.list[key] = append(s.list[key], value)

	return len(s.list[key])
}

func (s *store) lPop(key string) string {
	if s.list[key] == nil {
		return ""
	}

	value := s.list[key][0]
	s.list[key] = s.list[key][1:]

	return value
}

func (s *store) rPop(key string) string {
	if s.list[key] == nil {
		return ""
	}

	value := s.list[key][len(s.list[key])-1]
	s.list[key] = s.list[key][:len(s.list[key])-1]

	return value
}

func (s *store) lLen(key string) int {
	if s.list[key] == nil {
		return 0
	}

	return len(s.list[key])
}

func (s *store) lIndex(key string, index int) string {
	if s.list[key] == nil || index >= len(s.list[key]) || index < 0 {
		return ""
	}

	return s.list[key][index]
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
		return fmt.Sprintf("%v", s.getNumber(args[0]))
	case "LPUSH":
		return fmt.Sprintf("%v", s.lPush(args[0], args[1]))
	case "RPUSH":
		return fmt.Sprintf("%v", s.rPush(args[0], args[1]))
	case "LPOP":
		return s.lPop(args[0])
	case "RPOP":
		return s.rPop(args[0])
	case "LLEN":
		return fmt.Sprintf("%v", s.lLen(args[0]))
	case "LINDEX":
		index, _ := strconv.Atoi(args[1])
		return s.lIndex(args[0], index)
	default:
		return "Unknown command"
	}
}
