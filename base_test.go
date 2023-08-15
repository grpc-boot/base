package base

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"hash/crc32"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/grpc-boot/base/core/zaplogger"

	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
)

var (
	sm  ShardMap
	hs  HashSet
	btm Bitmap
)

var (
	hashGroup = &Group{}
	hostList  = []string{
		"192.168.1.135:3551:v0",
		"192.168.1.135:3551:v1",
		"192.168.1.135:3551:v2",
		"192.168.1.135:3551:v3",
		"192.168.1.135:3551:v4",
		"192.168.1.135:3552:v0",
		"192.168.1.135:3552:v1",
		"192.168.1.135:3552:v2",
		"192.168.1.135:3552:v3",
		"192.168.1.135:3552:v4",
		"192.168.1.135:3553:v0",
		"192.168.1.135:3553:v1",
		"192.168.1.135:3553:v2",
		"192.168.1.135:3553:v3",
		"192.168.1.135:3553:v4",
		"192.168.1.135:3554:v0",
		"192.168.1.135:3554:v1",
		"192.168.1.135:3554:v2",
		"192.168.1.135:3554:v3",
		"192.168.1.135:3554:v4",
		"192.168.1.135:3555:v0",
		"192.168.1.135:3555:v1",
		"192.168.1.135:3555:v2",
		"192.168.1.135:3555:v3",
		"192.168.1.135:3555:v4",
		"192.168.1.135:3556:v0",
		"192.168.1.135:3556:v1",
		"192.168.1.135:3556:v2",
		"192.168.1.135:3556:v3",
		"192.168.1.135:3556:v4",
		"192.168.1.135:3557:v0",
		"192.168.1.135:3557:v1",
		"192.168.1.135:3557:v2",
		"192.168.1.135:3557:v3",
		"192.168.1.135:3557:v4",
		"192.168.1.135:3558:v0",
		"192.168.1.135:3558:v1",
		"192.168.1.135:3558:v2",
		"192.168.1.135:3558:v3",
		"192.168.1.135:3558:v4",
		"192.168.1.135:3559:v0",
		"192.168.1.135:3559:v1",
		"192.168.1.135:3559:v2",
		"192.168.1.135:3559:v3",
		"192.168.1.135:3559:v4",
	}
)

type Data struct {
	CanHash

	id string
}

func (d *Data) HashCode() (hashValue uint32) {
	return crc32.ChecksumIEEE([]byte(d.id))
}

type Group struct {
	ring HashRing
}

func init() {
	rand.Seed(time.Now().UnixNano())

	sm = NewShardMap()
	hs = NewHashSet(10)
	btm = NewBitmap(nil)
}

func TestFileExists(t *testing.T) {
	caseList := []string{
		"func.go",
		"function.go",
	}

	for _, fileName := range caseList {
		exist := FileExists(fileName)
		t.Logf("%s: %v", fileName, exist)
	}
}

func TestHashValue(t *testing.T) {
	dd := &Data{
		id: "sfafd",
	}

	t.Log(dd.HashCode())

	t.Logf("123 hash:%d\n", HashValue(123))

	t.Logf("123123123123.123 hash:%d\n", HashValue(123123123123.123))
	t.Logf("map hash:%d\n", HashValue(map[string]string{"sdfc": "sdf"}))
}

func TestMap(t *testing.T) {
	d := Data{
		id: "cc",
	}

	keyValue := map[interface{}]interface{}{
		"user": map[string]interface{}{
			"id":   15,
			"name": "ddadf",
		},
		"listLength": 34,
		"key":        "value",
		d:            55,
	}

	for key, value := range keyValue {
		sm.Set(key, value)
	}

	if int64(len(keyValue)) != sm.Length() {
		t.Fatalf("want %d, got %d", len(keyValue), sm.Length())
	}

	val, exists := sm.Get("user")
	if !exists {
		t.Fatalf("want true, got %t", exists)
	}

	if _, ok := val.(map[string]interface{}); !ok {
		t.Fatalf("want true, got %t", ok)
	}

	val, exists = sm.Get(d)
	if !exists {
		t.Fatalf("want true, got %t", exists)
	}

	if val != 55 {
		t.Fatal("want true, got false")
	}

	sm.Delete("key")

	if exists = sm.Exists("key"); exists {
		t.Fatalf("want false, got %t", exists)
	}

	if int64(len(keyValue)-1) != sm.Length() {
		t.Fatalf("want %d, got %d", len(keyValue)-1, sm.Length())
	}
}

// BenchmarkMap_SetParallel-4       3692262               330 ns/op              49 B/op          2 allocs/op
func BenchmarkMap_SetParallel(b *testing.B) {
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := time.Now().UnixNano()
			sm.Set(c, c)
		}
	})
}

// BenchmarkMap_Set-4               2544157               641 ns/op             143 B/op          2 allocs/op
func BenchmarkMap_Set(b *testing.B) {
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		c := time.Now().UnixNano()
		sm.Set(c, c)
	}
}

// BenchmarkMap_GetParallel-4      14161170                95.0 ns/op             8 B/op          1 allocs/op
func BenchmarkMap_GetParallel(b *testing.B) {
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sm.Get(time.Now().UnixNano())
		}
	})
}

// BenchmarkMap_Get-4               5327881               225 ns/op               8 B/op          1 allocs/op
func BenchmarkMap_Get(b *testing.B) {
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		sm.Get(time.Now().UnixNano())
	}
}

func TestHashSet_Length(t *testing.T) {
	length := hs.Length()
	if length != 0 {
		t.Fatalf("length want 0, got %d", length)
	}

	newNum := hs.Add("1", 1000)
	if newNum != 2 {
		t.Fatalf("newNum want 2, got %d", newNum)
	}

	length = hs.Length()
	if length != 2 {
		t.Fatalf("length want 2, got %d", length)
	}

	newNum = hs.Add(1, 1000)
	if newNum != 1 {
		t.Fatalf("newNum want 1, got %d", newNum)
	}

	exists := hs.Exists(1000)
	if !exists {
		t.Fatal("want true, got false")
	}

	exists = hs.Exists("1000")
	if exists {
		t.Fatal("want false, got true")
	}
}

func TestBitmap_AddTag(t *testing.T) {
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.AddTag(56)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.AddTag(8)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.AddTag(15)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.DelTag(15)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())

	btm.DelTag(0)
	t.Logf("hasBit:%t bitCount:%d binary:%s", btm.HasBit(), btm.BitCount(), btm.SprinfBinary())
}

func TestAes_CbcEncrypt(t *testing.T) {
	aes128, err := NewAes("GS317MrfMqnvFcEHIbPTuQ==", "1cijdjijnji89iju")
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("sdf312e3213")
	secretData := aes128.CbcEncrypt(data)
	plain, err := aes128.CbcDecrypt(secretData)
	if err != nil {
		t.Fatalf("%swant nil, got %s", plain, err)
	}

	if !bytes.Equal(data, plain) {
		t.Fatal("want true, got false")
	}

	t.Log(hex.EncodeToString(secretData))

	aes192, err := NewAes("DFSDFi8987ysdfasadf3a23f", "ji12jnchFgDejffe")
	secretData = aes192.CbcEncrypt(data)
	plain, err = aes192.CbcDecrypt(secretData)
	if !bytes.Equal(data, plain) {
		t.Fatal("want true, got false")
	}

	t.Log(hex.EncodeToString(secretData))

	aes256, err := NewAes("DFSDFi8987ysdfasadf3a23fJUDH7yuG", "uh12jnchFgDejffe")
	secretData = aes256.CbcEncrypt(data)
	plain, err = aes256.CbcDecrypt(secretData)
	if !bytes.Equal(data, plain) {
		t.Fatal("want true, got false")
	}

	t.Log(hex.EncodeToString(secretData))
}

func TestCreateKeys(t *testing.T) {
	privateKey, publicKey := CreatePkcs1Keys(2048)
	t.Logf("privateKey:%s, publicKey:%s", privateKey, publicKey)
}

func TestCreatePkcs8Keys(t *testing.T) {
	privateKey, publicKey := CreatePkcs8Keys(2048)
	t.Logf("privateKey:%s, publicKey:%s", privateKey, publicKey)
}

func TestRsa_Decrypt(t *testing.T) {
	privateKey, publicKey := CreatePkcs1Keys(2048)
	t.Logf("privateKey:%s publicKey:%s", privateKey, publicKey)
	rsa, err := NewRsa(publicKey, privateKey)
	if err != nil {
		t.Fatalf("want nil, got %s", err.Error())
	}

	data := []byte("sadfasfd")
	secretData, err := rsa.Encrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	pData, err := rsa.Decrypt(secretData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%s, secretData:%s, decryptData:%s", string(data), hex.EncodeToString(secretData), string(pData))

	privateKey, publicKey = CreatePkcs8Keys(2048)
	rsa, _ = NewRsa(publicKey, privateKey)

	data = []byte("sadfaasdfasdfasdfsfd")
	secretData, err = rsa.Encrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	pData, err = rsa.Decrypt(secretData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%s, secretData:%s, decryptData:%s", string(data), hex.EncodeToString(secretData), string(pData))
}

func TestRsa_Verify(t *testing.T) {
	privateKey, publicKey := CreatePkcs1Keys(2048)
	rsa, _ := NewRsa(publicKey, privateKey)
	data := []byte("sadfaasdfasdfasdfsfd")
	sign, err := rsa.Sign(data, crypto.SHA1)
	if err != nil {
		t.Fatal(err)
	}

	if !rsa.Verify(data, sign, crypto.SHA1) {
		t.Fatal("want true got false")
	}
}

func TestSnowFlake_Id(t *testing.T) {
	begin, _ := time.ParseInLocation("2006-01-02", `2023-01-01`, time.Local)
	sf, err := NewSFByIp(ModeWait, begin.Unix())
	if err != nil {
		t.Fatal(err)
	}

	id, _ := sf.Id()
	t.Log(id)
	ts, machineId, index := sf.Info(id)

	if index != 1 {
		t.Fatalf("want 1, got %d", index)
	}

	t.Logf("milliTimestamp:%d machine:%d, index:%d", ts, machineId, index)

	id, _ = sf.Id()
	t.Log(id)
	ts, machineId, index = sf.Info(id)

	if index != 2 {
		t.Fatalf("want 1, got %d", index)
	}

	t.Logf("milliTimestamp:%d machine:%d, index:%d", ts, machineId, index)
}

func TestNewSFByMachineFunc(t *testing.T) {
	os.Setenv("MNUM", "128")

	begin, _ := time.ParseInLocation("2006-01-02", `2021-01-01`, time.Local)
	sf, err := NewSFByMachineFunc(ModeWait, GetMachineIdByEnv("MNUM"), begin.Unix())
	if err != nil {
		t.Fatal(err)
	}

	id, _ := sf.Id()
	t.Log(id)
	ts, machineId, index := sf.Info(id)

	if index != 1 {
		t.Fatalf("want 1, got %d", index)
	}

	if machineId != 128 {
		t.Fatalf("want 128, got %d", machineId)
	}

	t.Logf("milliTimestamp:%d machine:%d, index:%d", ts, machineId, index)
}

func TestHashRing_Get(t *testing.T) {
	serverList := make([]CanHash, 0, len(hostList))

	for _, server := range hostList {
		serverList = append(serverList, &Data{
			id: server,
		})
	}

	hashGroup.ring = NewHashRing(serverList...)

	for end := 1 << 25; end > 0; end-- {
		_, err := hashGroup.ring.Get(strconv.Itoa(end))
		if err != nil {
			t.Fatal(err)
		}
	}

	hashGroup.ring.Range(func(index int, server CanHash, hitCount uint64) (handled bool) {
		t.Logf("index:%d, server.id:%s, hitCount:%d", index, server.(*Data).id, hitCount)
		return
	})
}

// go test -bench=. -benchmem -v
// BenchmarkHashRing_GetIndex-4    24148118                65.9 ns/op            16 B/op          1 allocs/op
func BenchmarkHashRing_Get(b *testing.B) {
	serverList := make([]CanHash, 0, len(hostList))

	for _, server := range hostList {
		serverList = append(serverList, &Data{
			id: server,
		})
	}

	hashGroup.ring = NewHashRing(serverList...)
	var val atomic.Uint64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := hashGroup.ring.Get([]byte(strconv.FormatUint(val.Add(1), 10)))
			if err != nil {
				b.Fatal(err.Error())
			}
		}
	})
}

func TestBigCamels(t *testing.T) {
	caseList := map[string]string{
		"user-info": "UserInfo",
		"user-":     "User",
		"User-info": "UserInfo",
		"user-Info": "UserInfo",
		"user":      "User",
	}

	for c, r := range caseList {
		if r != BigCamels('-', c) {
			t.Fatal("want true, got false")
		}
	}
}

func TestSmallCamels(t *testing.T) {
	caseList := map[string]string{
		"user-info": "userInfo",
		"user-":     "user",
		"User-Info": "userInfo",
		"user-Info": "userInfo",
		"user":      "user",
	}

	for c, r := range caseList {
		if r != SmallCamels('-', c) {
			t.Fatal("want true, got false")
		}
	}
}

func TestContext_Next(t *testing.T) {
	ctx := AcquireCtx([]func(ctx *Context){
		func(ctx *Context) {
			current := time.Now().Unix()
			t.Logf("first run:%d\n", current)
			ctx.Set("cc", current)

			ctx.Next()

			t.Logf("first done\n")
		},
		func(ctx *Context) {
			current := time.Now().Unix()
			t.Logf("second run:%d\n", current)

			ccValue, _ := ctx.GetInt64("cc")

			t.Logf("got cc:%d\n", ccValue)
			ctx.Set("dd", current)

			ctx.Abort()

			t.Logf("second done\n")
		},
		func(ctx *Context) {
			current := time.Now().Unix()
			t.Logf("third run:%d\n", current)

			ddValue, _ := ctx.GetInt64("dd")

			t.Logf("got dd:%d\n", ddValue)
			ctx.Set("ee", current)

			ctx.Next()

			t.Logf("third done\n")
		},
	})

	ctx.Next()
	ctx.Close()
}

func BenchmarkBigCamels(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BigCamels('-', "user-info-api")
	}
}

func BenchmarkSmallCamels(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SmallCamels('-', "user-info-api")
	}
}

func JsonParamProvider() [][]byte {
	return [][]byte{
		[]byte(`{}`),
		[]byte(`{"id": 1, "name": "masco"}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6, "sex":null}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6, "sex":null, "tags":[1, 2, 3]}`),
		[]byte(`{"id": 2, "name": "masco", "score": 55.6, "sex":null, "tags":[1, 2, 3], "hobby":["basketball","football"]}`),
	}
}

func TestUnmarshalJsonParam(t *testing.T) {
	cL := JsonParamProvider()

	for _, data := range cL {
		p, err := UnmarshalJsonParam(data)
		if err != nil {
			t.Fatalf("want nil, got %+v", err)
		}
		t.Logf("%+v", p)
	}
}

func TestJsonParam_Get(t *testing.T) {
	cL := JsonParamProvider()

	for _, data := range cL {
		p, _ := UnmarshalJsonParam(data)
		t.Logf("id:%d name:%s score:%10.2f sex:%v tags:%v hobby:%v", p.Int64("id"), p.String("name"), p.Float64("score"), p.Int("sex"), p.Uint32Slice("tags"), p.StringSlice("hobby"))
	}
}

func TestStatusWithCode(t *testing.T) {
	stsOk := StatusWithCode(OK)
	defer stsOk.Close()

	t.Logf("status: %s", stsOk.JsonMarshal())

	stsArg := StatusWithCode(CodeInvalidArgument)
	defer stsArg.Close()

	t.Logf("status: %s", stsArg.JsonMarshal())

	stsPermiss := StatusWithCode(CodePermissionDenied)
	stsPermiss.WithMsg("FORBIDDEN")
	defer stsPermiss.Close()

	t.Logf("status: %s", stsPermiss.JsonMarshal())
}

func TestStatusWithJsonUnmarshal(t *testing.T) {
	sts, err := StatusWithJsonUnmarshal([]byte(`{"code": 0, "msg": "OK", "data":{}}`))
	if err != nil {
		t.Fatalf("unmarshal error:%s", err.Error())
	}
	defer sts.Close()

	t.Logf("status: %s", sts.JsonMarshal())

	stsCancel, err := StatusWithJsonUnmarshal([]byte(`{"code": 1, "msg": "CANCELLED", "data:{}}`))
	if err == nil {
		t.Fatalf("want error, got nil")
	}

	if stsCancel != nil {
		t.Fatalf("want nil, got:%s", stsCancel.JsonMarshal())
	}
}

func TestStatus_Is(t *testing.T) {
	sts := StatusWithCode(CodeCanceled)
	defer sts.Close()
	if !sts.Is(CodeCanceled) {
		t.Fatalf("want true, got false")
	}

	if sts.Is(OK) {
		t.Fatalf("want false, got true")
	}
}

func TestStatus_Error(t *testing.T) {
	sts := StatusWithCode(CodeNotFound)
	defer sts.Close()

	t.Logf("error format: %s", sts.Error().Error())

	stsOk := StatusWithCode(OK)
	defer stsOk.Close()

	if stsOk.Error() != nil {
		t.Fatalf("want nil, got %+v", stsOk.Error())
	}
}

func TestZapError(t *testing.T) {
	opt := zaplogger.Option{
		Level: int8(zapcore.InfoLevel),
		Path:  "./logs",
	}

	err := InitZapWithOption(opt)
	if err != nil {
		t.Fatalf("want nil, got %s\n", err.Error())
	}

	res := IsLevel(zapcore.InfoLevel)
	if !res {
		t.Fatalf("want true, got %t\n", res)
	}

	res = IsLevel(zapcore.DebugLevel)
	if res {
		t.Fatalf("want false, got %t\n", res)
	}

	res = IsLevel(zapcore.ErrorLevel)
	if !res {
		t.Fatalf("want true, got %t\n", res)
	}

	Debug("dddddd", zaplogger.Mid("asdfasdffasd"))
	Info("iiiiiiiii", zaplogger.Event("Test"))
	Warn("wwwwww", zaplogger.UpdateAt(time.Now().Unix()))
	Error("eeeeeeee", zaplogger.Value([]interface{}{123123, "safasf"}))
}

func TestInt64ToHex(t *testing.T) {
	t.Logf(Int64ToHex(0))
	t.Logf(Int64ToHex(math.MaxInt16))
	t.Logf(Int64ToHex(math.MaxUint16))
}

func TestHex2Int64(t *testing.T) {
	t.Log(Hex2Int64("00000f"))
}

func TestInt64ToHexWithPad(t *testing.T) {
	t.Logf(Int64ToHexWithPad(0, 5))
	t.Logf(Int64ToHexWithPad(math.MaxUint8, 1))
	t.Logf(Int64ToHexWithPad(math.MaxUint8, 6))
	t.Logf(Int64ToHexWithPad(math.MaxInt16, 7))
	t.Logf(Int64ToHexWithPad(math.MaxUint16, 8))
}

var (
	dAes atomic.Value
)

func defaultAes() *Aes {
	aes, ok := dAes.Load().(*Aes)
	if ok {
		return aes
	}

	aes, _ = NewAes(`i8jdi8jdkfjujui1`, `yhbDCFRE67hbgfde`)
	dAes.Store(dAes)
	return aes
}

func TestV1_Unpack(t *testing.T) {
	transKey := `FR4rjdi8jdkfjujui1yhbDhbgfdeCE67`
	aes := defaultAes()
	v1, err := NewV1(aes, aes.CbcEncrypt([]byte(transKey)))
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	data := v1.Pack(&Package{
		Id:   0x1001,
		Name: "login",
		Param: JsonParam{
			"t": "v",
		},
	})

	t.Logf("pack--%s", data)

	pkg, err := v1.Unpack(data)
	t.Logf("%+v", pkg)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}
}

// BenchmarkV1_Pack-8   	 1066840	      1356 ns/op
func BenchmarkV1_Pack(b *testing.B) {
	transKey := `FR4rjdi8jdkfjujui1yhbDhbgfdeCE67`
	aes := defaultAes()
	v1, err := NewV1(aes, aes.CbcEncrypt([]byte(transKey)))
	if err != nil {
		b.Fatalf("want nil, got %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Pack(&Package{
			Id:   0x1001,
			Name: "login",
			Param: JsonParam{
				"t": "v",
			},
		})
	}
}

// BenchmarkV1_Unpack-8   	 1261839	       1149 ns/op
func BenchmarkV1_Unpack(b *testing.B) {
	transKey := `FR4rjdi8jdkfjujui1yhbDhbgfdeCE67`
	aes := defaultAes()
	v1, err := NewV1(aes, aes.CbcEncrypt([]byte(transKey)))
	if err != nil {
		b.Fatalf("want nil, got %s", err)
	}

	data := v1.Pack(&Package{
		Id:   0x1001,
		Name: "login",
		Param: JsonParam{
			"t": "v",
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = v1.Unpack(data); err != nil {
			b.Fatalf("want nil, got %s", err)
		}
	}
}

// BenchmarkV1_PackUnpack-8   	 1378827	       1010 ns/op
func BenchmarkV1_PackUnpack(b *testing.B) {
	transKey := `FR4rjdi8jdkfjujui1yhbDhbgfdeCE67`
	aes := defaultAes()
	v1, err := NewV1(aes, aes.CbcEncrypt([]byte(transKey)))
	if err != nil {
		b.Fatalf("want nil, got %s", err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := v1.Pack(&Package{
				Id:   0x1001,
				Name: "login",
				Param: JsonParam{
					"t": time.Now().Unix(),
				},
			})

			if _, err = v1.Unpack(data); err != nil {
				b.Fatalf("want nil, got %s", err)
			}
		}
	})
}

func TestV2_Unpack(t *testing.T) {
	transKey := `i1yhbDhbgfdeCE67`
	aes := defaultAes()
	protocol, err := NewV2(aes, aes.CbcEncrypt([]byte(transKey)))
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	explainData, err := aes.CbcDecrypt(protocol.ResponseKey())
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}
	t.Logf("iv:%s %s", protocol.ResponseKey(), explainData)

	data := protocol.Pack(&Package{
		Id:   0x1001,
		Name: "login",
		Param: JsonParam{
			"t":       "v",
			"current": time.Now().UnixNano(),
		},
	})

	t.Logf("pack--%s", data)

	pkg, err := protocol.Unpack(data)
	t.Logf("%+v", pkg)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}
}

func TestAccept_Accept(t *testing.T) {
	aes := defaultAes()
	accept := NewAccept(aes, LevelJson)

	protoV0, err := accept.Accept(LevelJson, nil)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	v1Key := []byte(`FR4#$%i8jdkfjujui1yhbDhbgfdeCE67`)

	protoV1, err := accept.Accept(LevelV1, aes.CbcEncrypt(v1Key))
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	v2Key := []byte(`FR4#$%i8jdkf@&ju`)
	protoV2, err := accept.Accept(LevelV2, aes.CbcEncrypt(v2Key))
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	pkg := &Package{
		Id:   EventLogin,
		Name: "login",
		Param: JsonParam{
			"currentTime": time.Now().Unix(),
		},
	}

	v0Pack := protoV0.Pack(pkg)
	v1Pack := protoV1.Pack(pkg)
	v2Pack := protoV2.Pack(pkg)

	v0Pkg, err := protoV0.Unpack(v0Pack)
	t.Logf("v0: %+v", v0Pkg)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	v1Pkg, err := protoV1.Unpack(v1Pack)
	t.Logf("v1: %+v", v1Pkg)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	v2Pkg, err := protoV2.Unpack(v2Pack)
	t.Logf("v2: %+v", v2Pkg)
	if err != nil {
		t.Fatalf("want nil, got %s", err)
	}

	t.Logf("v2-iv:%s", protoV2.ResponseKey())

	accept = NewAccept(aes, LevelV1)
	protoV0, err = accept.Accept(LevelJson, nil)
	if err == nil {
		t.Fatalf("want nil, got %s", err)
	}
}

func TestRandSeed(t *testing.T) {
	for i := 1; i < 128; i++ {
		str := randSeed(i)
		t.Logf(str)
		if len(str) != i {
			t.Fatalf("want %d, got %d", i, len(str))
		}
	}
}

func BenchmarkRandSeed(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			length := rand.Intn(len(seed) / 2)
			str := randSeed(length)
			if len(str) != length {
				b.Fatalf("want %d, got %d", length, len(str))
			}
		}
	})
}

func BenchmarkRandSeedBig(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			length := rand.Intn(1024)
			str := randSeed(length)
			if len(str) != length {
				b.Fatalf("want %d, got %d", length, len(str))
			}
		}
	})
}

func TestRandStr(t *testing.T) {
	for i := 1; i < 1024; i++ {
		rBytes := RandBytes(i)
		t.Logf("%s", rBytes)
		if len(rBytes) != i {
			t.Fatalf("want %d, got %d", i, len(rBytes))
		}
	}
}

func BenchmarkRandStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		length := 1 + rand.Intn(16)
		rBytes := RandBytes(length)
		if len(rBytes) != length {
			b.Fatalf("want %d, got %d", length, len(rBytes))
		}
	}
}

func TestState_Remove(t *testing.T) {
	var st State
	st.Add(0)
	if st != 1 {
		t.Fatalf("want 1, got %d", st)
	}

	st.Remove(0)
	if st != 0 {
		t.Fatalf("want 0, got %d", st)
	}

	for i := 0; i < StateIndexMax; i++ {
		var index = uint8(rand.Intn(StateIndexMax + 1))
		st.Add(index)
		if !st.Has(index) {
			t.Fatal("want true, got false")
		}

		st.Remove(index)

		if st.Has(index) {
			t.Fatal("want false, got true")
		}

		st.Add(index)
	}

	t.Logf("slice: %+v\nvalueSlice:%+v", st.Slice(), st.ValueSlice())

	st.RemoveAll()
	if st != 0 {
		t.Fatalf("want 0, got %d", st)
	}
}

func TestState_Merge(t *testing.T) {
	var A State
	var B State

	A.Add(0, 1, 4, 7, 8, 30)
	B.Add(2)

	A.Merge(B)

	if !A.Has(2) {
		t.Fatal("want true, got false")
	}

	t.Logf("union: %+v intersection:%+v", A.UnionSet(B), A.Intersection(B))

	B.Add(0, 12, 7, 8, 30)
	t.Logf("union: %+v intersection:%+v", A.UnionSet(B), A.Intersection(B))

	A.Merge(B)
	t.Logf("A: %+v B:%+v", A.Slice(), B.Slice())
}

func TestState_ValueSlice(t *testing.T) {
	var indexList = make([]uint8, 10, 10)
	for i := 0; i < 10; i++ {
		indexList[i] = uint8(rand.Intn(StateIndexMax))
	}

	st1, _ := StateFromSlice(indexList)
	st2 := StatusFromValueSlice(st1.ValueSlice())

	for _, index := range indexList {
		if !st1.Has(index) {
			t.Fatal("want true, got false")
		}

		if !st2.Has(index) {
			t.Fatal("want true, got false")
		}
	}

	t.Logf("st1:%+v, st2:%+v", st1.Slice(), st2.Slice())
}

func TestLocalIp(t *testing.T) {
	localIp, err := LocalIp()
	if err != nil {
		t.Fatalf("want nil, got %s", err.Error())
	}

	t.Logf("local ip:%s", localIp)
}

func TestLong2Ip(t *testing.T) {
	for index := 0; index < 100; index++ {
		longIp := rand.Uint32()
		ipStr := Long2Ip(longIp)
		uint32Ip, err := Ip2Long(ipStr)
		if err != nil {
			t.Fatalf("want nil, got %s", err.Error())
		}

		if uint32Ip != longIp {
			t.Fatalf("want equal, got %t", false)
		}

		t.Logf("longIp:%d stringIp:%s", longIp, ipStr)
	}
}

func BenchmarkLong2Ip(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Long2Ip(rand.Uint32())
		}
	})
}

func BenchmarkIp2Long(b *testing.B) {
	var (
		size = 500 * 10000
		list = make([]string, size)
	)

	for index := 0; index < size; index++ {
		list[index] = Long2Ip(rand.Uint32())
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := Ip2Long(list[rand.Intn(size)]); err != nil {
				b.Fatalf("want nil, got %v", err)
			}
		}
	})
}
