package components

import (
	"crypto"
	"encoding/hex"
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

	data = []byte("阿斯顿发的sadfaasdfasdfasdfsfd")
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
	data := []byte("我撒旦法sadfaasdfasdfasdfsfd")
	sign, err := rsa.Sign(data, crypto.SHA1)
	if err != nil {
		t.Fatal(err)
	}

	if !rsa.Verify(data, sign, crypto.SHA1) {
		t.Fatal("want true got false")
	}
}
