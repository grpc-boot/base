package kind

import "errors"

var (
	ErrIndex = errors.New("index out of range")
)

const (
	StateIndexMax = 30
)

type State int32

func StatusFromValueSlice(list []int32) State {
	if len(list) < 1 {
		return 0
	}

	var st State
	for _, value := range list {
		st = st | State(value)
	}

	return st
}

func StateFromSlice(list []uint8) (State, error) {
	var st State

	if err := st.Add(list...); err != nil {
		return 0, err
	}

	return st, nil
}

func (s *State) Add(indexList ...uint8) error {
	if len(indexList) < 1 {
		return nil
	}

	for _, index := range indexList {
		if index > StateIndexMax {
			return ErrIndex
		}

		*s = *s | (1 << index)
	}

	return nil
}

func (s *State) Remove(index uint8) {
	if index > StateIndexMax {
		return
	}

	*s = *s & ^(1 << index)
}

func (s *State) RemoveAll() {
	*s = 0
}

func (s *State) Has(index uint8) bool {
	if index > StateIndexMax {
		return false
	}

	return (*s & (1 << index)) > 0
}

func (s *State) Slice() []uint8 {
	if *s == 0 {
		return nil
	}

	var (
		st = make([]uint8, 0, StateIndexMax)
		i  uint8
	)

	for ; i <= StateIndexMax; i++ {
		if s.Has(i) {
			st = append(st, i)
		}
	}

	return st
}

func (s *State) ValueSlice() []int32 {
	if *s == 0 {
		return nil
	}

	var (
		st = make([]int32, 0, StateIndexMax)
		i  uint8
	)

	for ; i <= StateIndexMax; i++ {
		if s.Has(i) {
			st = append(st, 1<<i)
		}
	}

	return st
}

func (s *State) Merge(val State) {
	*s = *s | val
}

func (s *State) UnionSet(val State) []uint8 {
	if val == 0 || *s == 0 {
		return nil
	}

	var (
		st = make([]uint8, 0, StateIndexMax)
		i  uint8
	)

	for ; i <= StateIndexMax; i++ {
		if s.Has(i) || val.Has(i) {
			st = append(st, i)
		}
	}

	return st
}

func (s *State) Intersection(val State) []uint8 {
	if val == 0 || *s == 0 {
		return nil
	}

	var (
		st = make([]uint8, 0, StateIndexMax)
		i  uint8
	)

	for ; i <= StateIndexMax; i++ {
		if s.Has(i) && val.Has(i) {
			st = append(st, i)
		}
	}

	return st
}
