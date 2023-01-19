package connctx

var (
	setValue = struct{}{}
)

type set map[interface{}]struct{}

func (s set) exists(key interface{}) bool {
	_, exists := s[key]
	return exists
}

func (s set) add(keys ...interface{}) (newNum int) {
	for _, key := range keys {
		if !s.exists(key) {
			newNum++
			s[key] = setValue
		}
	}

	return
}

func (s set) card() int {
	return int(len(s))
}

func (s set) members() (keys []interface{}) {
	if s.card() < 1 {
		return
	}

	keys = make([]interface{}, s.card())

	index := 0

	for key, _ := range s {
		keys[index] = key
		index++
	}

	return
}

func (s set) isMember(key interface{}) bool {
	_, exists := s[key]
	return exists
}

func (s set) rem(keys ...interface{}) (delNum int) {
	for _, key := range keys {
		if _, exists := s[key]; exists {
			delete(s, key)
			delNum++
		}
	}

	return
}

func (s set) pop() (item interface{}) {
	if s.card() < 1 {
		return
	}

	for key, _ := range s {
		item = key
		delete(s, key)
		break
	}

	return
}

func (s set) randMember() (item interface{}) {
	if s.card() < 1 {
		return
	}

	for key, _ := range s {
		item = key
		break
	}

	return
}
