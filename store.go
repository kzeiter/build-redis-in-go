package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type store struct {
	data        map[string]interface{}
	list        map[string][]string
	sets        map[string]map[string]bool
	subscribers map[string][]client
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
func (s *store) sadd(key string, value string) bool {
	if s.sets[key] == nil {
		s.sets[key] = make(map[string]bool)
	}

	if s.sets[key][value] {
		return false
	}

	s.sets[key][value] = true

	return true
}

func (s *store) srem(key string, value string) bool {
	if s.sets[key] == nil {
		return false
	}

	if !s.sets[key][value] {
		return false
	}

	delete(s.sets[key], value)

	return true
}

func (s *store) smembers(key string) []string {
	if s.sets[key] == nil {
		return []string{}
	}

	members := make([]string, 0, len(s.sets[key]))

	for member := range s.sets[key] {
		members = append(members, member)
	}

	return members
}

func (s *store) sismember(key string, value string) bool {
	if s.sets[key] == nil {
		return false
	}

	return s.sets[key][value]
}

func (s *store) subscribe(channel string, conn net.Conn) string {
	client := client{conn: conn}
	s.subscribers[channel] = append(s.subscribers[channel], client)
	return "OK"
}

func (s *store) publish(channel string, message string) {
	subscribers, ok := s.subscribers[channel]
	if ok {
		for _, subscriber := range subscribers {
			fmt.Fprintf(subscriber.conn, "+%s\n", message)
		}
	}
}

func (s *store) handleCommand(command string, args []string, conn net.Conn) string {
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
	case "SADD":
		return fmt.Sprintf("%v", s.sadd(args[0], args[1]))
	case "SREM":
		return fmt.Sprintf("%v", s.srem(args[0], args[1]))
	case "SMEMBERS":
		members := s.smembers(args[0])
		result := ""

		for _, member := range members {
			result += fmt.Sprintf("%v ", member)
		}

		return strings.TrimSpace(result)
	case "SISMEMBER":
		return fmt.Sprintf("%v", s.sismember(args[0], args[1]))
	case "SUBSCRIBE":
		return s.subscribe(args[0], conn)
	case "PUBLISH":
		s.publish(args[0], args[1])
		return "OK"
	default:
		return "Unknown command"
	}
}
