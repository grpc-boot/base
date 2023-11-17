package kind

type Set[V comparable] interface {
	Add(value V) (isNew bool)
	Del(list ...V) (delNum int)
	Exists(value V) (exists bool)
	Length() (length int)
	List() (list Slice[V])
}

// set hash set，非线程安全
type set[V comparable] struct {
	data map[V]struct{}
}

func NewSet[V comparable](initSize uint) Set[V] {
	return &set[V]{
		data: make(map[V]struct{}, initSize),
	}
}

func (s *set[V]) Add(value V) (isNew bool) {
	_, exists := s.data[value]
	if !exists {
		s.data[value] = Empty
	}

	isNew = !exists

	return
}

func (s *set[V]) Del(list ...V) (delNum int) {
	if len(list) < 1 {
		return
	}

	for _, value := range list {
		_, exists := s.data[value]
		if exists {
			delete(s.data, value)
			delNum++
		}
	}

	return
}

func (s *set[V]) Exists(value V) (exists bool) {
	_, exists = s.data[value]
	return
}

func (s *set[V]) Length() (length int) {
	return len(s.data)
}

func (s *set[V]) List() (list Slice[V]) {
	list = make(Slice[V], len(s.data))

	if len(s.data) < 1 {
		return
	}

	i := 0

	for value, _ := range s.data {
		list[i] = value
		i++
	}

	return
}
