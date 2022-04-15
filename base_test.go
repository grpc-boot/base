package base

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"hash/crc32"
	"os"
	"strconv"
	"testing"
	"time"

	"go.uber.org/atomic"
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
		t.Fatal(err)
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
	rsa, _ := NewRsa(publicKey, privateKey)

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
	begin, _ := time.ParseInLocation("2006-01-02", `2021-01-01`, time.Local)
	sf, err := NewSFByIp(ModeWait, begin.UnixNano()/1e6)
	if err != nil {
		t.Fatal(err)
	}

	id, _ := sf.Id(1)
	t.Log(id)
	ts, machineId, logicId, index := sf.Info(id)
	if logicId != 1 {
		t.Fatalf("want 1, got %d", logicId)
	}

	if index != 1 {
		t.Fatalf("want 1, got %d", index)
	}
	t.Log(ts, machineId, logicId, index)
}

func TestNewSFByMachineFunc(t *testing.T) {
	os.Setenv("MNUM", "128")

	begin, _ := time.ParseInLocation("2006-01-02", `2021-01-01`, time.Local)
	sf, err := NewSFByMachineFunc(ModeWait, GetMachineIdByEnv("MNUM"), begin.UnixNano()/1e6)
	if err != nil {
		t.Fatal(err)
	}

	id, _ := sf.Id(3)
	t.Log(id)
	ts, machineId, logicId, index := sf.Info(id)
	if logicId != 3 {
		t.Fatalf("want 3, got %d", logicId)
	}

	if index != 1 {
		t.Fatalf("want 1, got %d", index)
	}

	if machineId != 128 {
		t.Fatalf("want 128, got %d", machineId)
	}

	t.Log(ts, machineId, logicId, index)
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
