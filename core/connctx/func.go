package connctx

func getIndex(offset uint16) int {
	return int(offset / 8)
}

func getByte(offset uint16) byte {
	return byte(1 << (7 - (offset % 8)))
}

func realIndex(index int, length int) (realIndex int, err error) {
	if index > -1 {
		return index, nil
	}

	realIndex = length + index
	if realIndex < 0 {
		return 0, ErrIndexOutOfRange
	}

	return
}
