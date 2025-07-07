package components

import (
	"crypto"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

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

	data := []byte("撒旦法sadfasfd")
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

	data = []byte(strings.Repeat("顿发的sadfaasdfafsfd阿斯asdfadsdfd", 120))
	secretData, err = rsa.Encrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	pData, err = rsa.Decrypt(secretData)
	if err != nil {
		t.Fatal(err)
	}

	if string(pData) != string(data) {
		t.Fatalf("want %s, got %s", string(data), string(pData))
	}

	t.Logf("data:%s, secretData:%s, decryptData:%s", string(data), hex.EncodeToString(secretData), string(pData))
}

func TestRsa_DecryptOAEP(t *testing.T) {
	privateKey, publicKey := CreatePkcs1Keys(2048)
	t.Logf("privateKey:%s publicKey:%s", privateKey, publicKey)
	rsa, err := NewRsa(publicKey, privateKey)
	if err != nil {
		t.Fatalf("want nil, got %s", err.Error())
	}

	data := []byte("撒旦法sadfasfd")
	secretData, err := rsa.EncryptOAEP(data, sha256.New(), nil)
	if err != nil {
		t.Fatal(err)
	}
	pData, err := rsa.DecryptOAEP(secretData, sha256.New(), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%s, secretData:%s, decryptData:%s", string(data), hex.EncodeToString(secretData), string(pData))

	privateKey, publicKey = CreatePkcs8Keys(2048)
	rsa, _ = NewRsa(publicKey, privateKey)

	data = []byte(strings.Repeat("阿斯顿发的sadfaasdfasdfasdfsfd", 128))
	secretData, err = rsa.EncryptOAEP(data, sha1.New(), nil)
	if err != nil {
		t.Fatal(err)
	}
	pData, err = rsa.DecryptOAEP(secretData, sha1.New(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(pData) {
		t.Fatalf("want %s, got %s", string(data), string(pData))
	}

	t.Logf("data:%s, secretData:%s, decryptData:%s", string(data), hex.EncodeToString(secretData), string(pData))
}

func TestRsa_Verify(t *testing.T) {
	privateKey, publicKey := CreatePkcs1Keys(2048)
	rsa, _ := NewRsa(publicKey, privateKey)
	data := []byte("我撒旦法sadfaasdfasdfasdfsfd")
	sign, err := rsa.Sign(data, crypto.SHA1)
	if err != nil {
		t.Fatal(err)
	}

	if !rsa.Verify(data, sign, crypto.SHA1) {
		t.Fatal("want true got false")
	}
}
