package base

import (
	"math/rand"
	"strings"
)

const (
	seed = `ACQP12fsd3F4689qab12wVGczSjDF123SDfkynm32SDb123vxF12lEsS90DFm34SDF348z1230HDW7sdfj239sfdS8DF4SDFasdfjSDFekSDF230de0SDerpljmnF9sDFSDp23outF59mnSDqazxcS89HJvb34nhjLKHdfse6wrKyuipp0987G654A321vxbcLKhyweBNHTEw12332sdfGJTRsdfswf1296dACQghopMNKsdfwe123asd08AS231`
)

func randSeed(length int) string {
	middle := len(seed) / 2
	if length > middle {
		var strBuilder strings.Builder
		strBuilder.WriteString(randSeed(middle))
		rest := length - middle
		for i := 0; i < rest; i++ {
			strBuilder.WriteByte(seed[rand.Intn(len(seed))])
		}
		return strBuilder.String()
	}

	start := rand.Intn(len(seed))
	if len(seed)-start < length {
		start -= middle
	}

	return seed[start : start+length]
}

func RandBytes(length int) []byte {
	if length < 1 {
		return nil
	}

	data := make([]byte, 0, length)

	if length > 16 {
		rest := length - 16
		if rand.Intn(0xffff)%2 == 0 {
			data = append(data, RandBytes(16)...)
			data = append(data, randSeed(rest)...)
			return data
		}

		data = append(data, randSeed(rest)...)
		data = append(data, RandBytes(16)...)
		return data
	}

	if length == 1 {
		return append(data, randSeed(length)...)
	}

	suffix := strings.Repeat("0", length-1)
	min := Hex2Int64("1" + suffix)
	max := Hex2Int64("f" + suffix)

	return append(data, Uint64ToHex(uint64(min)+uint64(rand.Int63n(max)))...)
}
