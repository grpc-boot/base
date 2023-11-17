package components

func Max6() int64 {
	return defaultIdCode.Max6()
}

func Max8() int64 {
	return defaultIdCode.Max8()
}

func Code6(id int64) (code []byte, err *BError) {
	return defaultIdCode.Code6(id)
}

func Code6String(id int64) (code string, err *BError) {
	return defaultIdCode.Code6String(id)
}

func Code8(id int64) (code []byte, err *BError) {
	return defaultIdCode.Code8(id)
}

func Code8String(id int64) (code string, err *BError) {
	return defaultIdCode.Code8String(id)
}

func Code2Id(code []byte) (id int64, err *BError) {
	return defaultIdCode.Code2Id(code)
}

func CodeString2Id(code string) (id int64, err *BError) {
	return defaultIdCode.CodeString2Id(code)
}

func Code6To8(code6 []byte) (code8 []byte, err *BError) {
	return defaultIdCode.Code6To8(code6)
}

func CodeString6To8(code6 string) (code8 string, err *BError) {
	return defaultIdCode.CodeString6To8(code6)
}
