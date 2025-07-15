package utils

func FeistelEncrypt64(num, key uint64) uint64 {
	const rounds = 4
	var l, r = uint32(num >> 32), uint32(num)
	for i := 0; i < rounds; i++ {
		newL := r
		newR := l ^ FeistelRound32(r, key, uint32(i))
		l, r = newL, newR
	}
	return (uint64(l) << 32) | uint64(r)
}

func FeistelDecrypt64(num, key uint64) uint64 {
	const rounds = 4
	var l, r = uint32(num >> 32), uint32(num)
	for i := rounds - 1; i >= 0; i-- {
		newR := l
		newL := r ^ FeistelRound32(l, key, uint32(i))
		l, r = newL, newR
	}
	return (uint64(l) << 32) | uint64(r)
}

func FeistelRound32(r uint32, key uint64, round uint32) uint32 {
	r ^= uint32(key>>32) ^ uint32(key)
	r = ((r >> 3) | (r << 29)) + round*0x9e3779b9
	r ^= r >> 16
	return r
}

func FeistelEncrypt32(v, salt uint32) uint32 {
	const rounds = 3
	l, r := v>>16, v&0xFFFF
	for i := 0; i < rounds; i++ {
		l, r = r, l^FeistelRound16(r, salt, uint32(i))
	}
	return (l << 16) | r
}

func FeistelDecrypt32(v, salt uint32) uint32 {
	const rounds = 3
	l, r := v>>16, v&0xFFFF
	for i := rounds - 1; i >= 0; i-- {
		l, r = r^FeistelRound16(l, salt, uint32(i)), l
	}
	return (l << 16) | r
}

func FeistelRound16(r, salt, round uint32) uint32 {
	return ((r ^ salt ^ round) * 0x5bd1e995) & 0xFFFF
}
