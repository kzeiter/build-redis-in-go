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

func (s *store) del(key string) {
	delete(s.data, key)
}

func (s *store) handleCommand(command string, args []string) string {
	switch command {
	case "SET":
		s.set(args[0], args[1])
		return "OK"
	case "GET":
		return s.get(args[0])
	case "DEL":
		s.del(args[0])
		return "OK"
	default:
		return "ERROR"
	}
}
