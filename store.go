package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Store struct {
	data map[string]string
	list map[string][]string
	sets map[string]map[string]bool
	subs map[string][]client
	disk *diskStore
}

func (s *Store) set(key string, value string) string {
	s.data[key] = value
	s.disk.save(s.data)

	return value
}

func (s *Store) get(key string) string {
	return s.data[key]
}

func (s *Store) del(key string) {
	delete(s.data, key)
	s.disk.save(s.data)
}

func (s *Store) incr(key string) int {
	value, _ := strconv.Atoi(s.get(key))

	s.set(key, strconv.Itoa(value+1))

	return value + 1
}

func (s *Store) incrBy(key string, incr string) int {
	number, _ := strconv.Atoi(incr)

	value, _ := strconv.Atoi(s.get(key))

	s.set(key, strconv.Itoa(value+number))

	return value + number
}

func (s *Store) decr(key string) int {
	value, _ := strconv.Atoi(s.get(key))

	s.set(key, strconv.Itoa(value-1))

	return value - 1
}

func (s *Store) decrBy(key string, decr string) int {
	number, _ := strconv.Atoi(decr)

	value, _ := strconv.Atoi(s.get(key))

	s.set(key, strconv.Itoa(value-number))

	return value - number
}

func (s *Store) lPush(key string, value string) int {
	if s.list[key] == nil {
		s.list[key] = make([]string, 0)
	}

	s.list[key] = append([]string{value}, s.list[key]...)

	return len(s.list[key])
}

func (s *Store) rPush(key string, value string) int {
	if s.list[key] == nil {
		s.list[key] = make([]string, 0)
	}

	s.list[key] = append(s.list[key], value)

	return len(s.list[key])
}

func (s *Store) lPop(key string) string {
	if s.list[key] == nil {
		return ""
	}

	value := s.list[key][0]
	s.list[key] = s.list[key][1:]

	return value
}

func (s *Store) rPop(key string) string {
	if s.list[key] == nil {
		return ""
	}

	value := s.list[key][len(s.list[key])-1]
	s.list[key] = s.list[key][:len(s.list[key])-1]

	return value
}

func (s *Store) lLen(key string) int {
	if s.list[key] == nil {
		return 0
	}

	return len(s.list[key])
}

func (s *Store) lIndex(key string, index int) string {
	if s.list[key] == nil || index >= len(s.list[key]) || index < 0 {
		return ""
	}

	return s.list[key][index]
}
func (s *Store) sadd(key string, value string) bool {
	if s.sets[key] == nil {
		s.sets[key] = make(map[string]bool)
	}

	if s.sets[key][value] {
		return false
	}

	s.sets[key][value] = true

	return true
}

func (s *Store) srem(key string, value string) bool {
	if s.sets[key] == nil {
		return false
	}

	if !s.sets[key][value] {
		return false
	}

	delete(s.sets[key], value)

	return true
}

func (s *Store) smembers(key string) []string {
	if s.sets[key] == nil {
		return []string{}
	}

	members := make([]string, 0, len(s.sets[key]))

	for member := range s.sets[key] {
		members = append(members, member)
	}

	return members
}

func (s *Store) sismember(key string, value string) bool {
	if s.sets[key] == nil {
		return false
	}

	return s.sets[key][value]
}

func (s *Store) subscribe(channel string, conn net.Conn) string {
	client := client{conn: conn}
	s.subs[channel] = append(s.subs[channel], client)
	return "OK"
}

func (s *Store) publish(channel string, message string) {
	subs, ok := s.subs[channel]
	if ok {
		for _, subscriber := range subs {
			fmt.Fprintf(subscriber.conn, "+%s\n", message)
		}
	}
}

func (s *Store) handleCommand(command string, args []string, conn net.Conn) string {
	switch command {
	case "SET":
		return s.set(args[0], strings.Join(args[1:], " "))
	case "GET":
		return s.get(args[0])
	case "DEL":
		s.del(args[0])
		return "OK"
	case "INCR":
		return fmt.Sprintf("%v", s.incr(args[0]))
	case "INCRBY":
		return fmt.Sprintf("%v", s.incrBy(args[0], args[1]))
	case "DECR":
		return fmt.Sprintf("%v", s.decr(args[0]))
	case "DECRBY":
		return fmt.Sprintf("%v", s.decrBy(args[0], args[1]))
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

func NewStore(filename string) (*Store, error) {
	disk := &diskStore{filename: filename}

	store := &Store{
		data: make(map[string]string),
		list: make(map[string][]string),
		sets: make(map[string]map[string]bool),
		subs: make(map[string][]client),
		disk: disk,
	}

	data, _ := store.disk.load()
	if len(data) != 0 {
		store.data = data
	}

	return store, nil
}
