package connctx

import "sort"

type setlist []zsetItem

type zsetItem struct {
	member string
	score  float64
}

func (sl setlist) Len() int {
	return len(sl)
}

func (sl setlist) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

func (sl setlist) Less(i, j int) bool {
	return sl[i].score > sl[j].score
}

func (sl setlist) indexOf(member string) (index int) {
	for index, _ = range sl {
		if sl[index].member == member {
			return index
		}
	}

	return -1
}

type zset struct {
	sl setlist
}

func (z *zset) card() int {
	return len(z.sl)
}

func (z *zset) count(min, max float64) (num int) {
	if min > max {
		return
	}

	index := sort.Search(z.card(), func(i int) bool {
		return z.sl[i].score >= min && z.sl[i].score <= max
	})

	if index >= z.card() {
		return
	}

	for ; index < z.card(); index++ {
		if z.sl[index].score >= min && z.sl[index].score <= max {
			num++
		} else {
			break
		}
	}

	return
}

func (z *zset) addMap(items ...zsetItem) (newNum int) {
	if len(items) < 1 {
		return
	}

	if z.sl == nil {
		z.sl = make([]zsetItem, 0, len(items))
	}

	for _, item := range items {
		index := z.sl.indexOf(item.member)
		if index == -1 {
			z.sl = append(z.sl, item)
			newNum++
		} else {
			z.sl[index].score = item.score
		}
	}

	sort.Sort(z.sl)
	return
}

func (z *zset) incrby(value float64, member string) float64 {
	if z.card() < 1 {
		z.addMap(zsetItem{member: member, score: value})
		return value
	}

	index := z.sl.indexOf(member)
	if index < 0 {
		z.addMap(zsetItem{member: member, score: value})
		return value
	}

	z.sl[index].score += value

	return z.sl[index].score
}

func (z *zset) delMember(member string) (hasDel bool) {
	if z.sl[0].member == member {
		hasDel = true

		if z.card() == 1 {
			z.sl = z.sl[:0]
			return
		}

		z.sl = z.sl[1:]
		return
	}

	if z.sl[z.card()-1].member == member {
		hasDel = true
		z.sl = z.sl[:z.card()-1]
		return
	}

	index := 1
	for ; index < z.card()-1; index++ {
		if z.sl[index].member == member {
			copy(z.sl[index:], z.sl[index+1:])
			z.sl = z.sl[:z.card()-1]

			hasDel = true
			break
		}
	}

	return
}

func (z *zset) rem(members ...string) (delNum int) {
	if z.card() < 1 || len(members) < 1 {
		return
	}

	for _, member := range members {
		if z.delMember(member) {
			delNum++
		}
	}

	return
}

func (z *zset) rank(member string) (rank int, exists bool) {
	if z.card() < 1 {
		return
	}

	index := z.sl.indexOf(member)
	if index > -1 {
		exists = true
		rank = z.card() - index - 1
	}

	return
}

func (z *zset) revrank(member string) (rank int, exists bool) {
	if z.card() < 1 {
		return
	}

	index := z.sl.indexOf(member)
	if index > -1 {
		exists = true
		rank = index
	}

	return
}

func (z *zset) score(member string) (score float64, exists bool) {
	if z.card() < 1 {
		return
	}

	index := z.sl.indexOf(member)
	if index > -1 {
		exists = true
		score = z.sl[index].score
	}

	return
}

func (z *zset) revrange(start, stop int) (items []zsetItem) {
	if z.card() < 1 || start > stop {
		return
	}

	realStart, err := realIndex(start, z.card())
	if err != nil {
		return
	}

	realStop, err := realIndex(stop, z.card())
	if err != nil {
		realStop = z.card() - 1
	}

	if realStart > realStop {
		return
	}

	itemCount := realStop - realStart + 1
	items = make([]zsetItem, itemCount, itemCount)
	return
}
