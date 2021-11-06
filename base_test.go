package base

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"hash/crc32"
	"testing"
	"time"
)

var (
	sm  ShardMap
	hs  HashSet
	btm Bitmap
)

type Data struct {
	CanHash

	id string
}

func (d *Data) HashCode() (hashValue uint32) {
	return crc32.ChecksumIEEE([]byte(d.id))
}

func init() {
	sm = NewShardMap()
	hs = NewHashSet(10)
	btm = NewBitmap(nil)
}

func TestHashValue(t *testing.T) {
	dd := &Data{
		id: "sfafd",
	}

	t.Log(dd.HashCode())
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
	privateKey, publicKey := CreateKeys(2048)
	t.Logf("privateKey:%s, publicKey:%s", privateKey, publicKey)
}

func TestCreatePkcs8Keys(t *testing.T) {
	privateKey, publicKey := CreatePkcs8Keys(2048)
	t.Logf("privateKey:%s, publicKey:%s", privateKey, publicKey)
}

func TestRsa_Decrypt(t *testing.T) {
	privateKey, publicKey := CreateKeys(2048)
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
	privateKey, publicKey := CreateKeys(2048)
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
