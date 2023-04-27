package main

type store struct {
	data map[string]string
}

func (s *store) set(key, value string) {
	s.data[key] = value
}

func (s *store) get(key string) string {
	return s.data[key]
}
