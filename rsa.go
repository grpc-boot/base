package base

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"strings"
)

type Rsa struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRsa(publicKey, privateKey string) (r *Rsa, err error) {
	if privateKey != "" && strings.Index(privateKey, "BEGIN RSA") > 0 {
		return NewRsaWithPkcs1(publicKey, privateKey)
	}
	return NewRsaWithPkcs8(publicKey, privateKey)
}

func NewRsaWithPkcs8(publicKey, privateKey string) (r *Rsa, err error) {
	var (
		pubKey *rsa.PublicKey
		priKey *rsa.PrivateKey
	)

	if privateKey != "" {
		block, _ := pem.Decode([]byte(privateKey))
		var pKey interface{}
		pKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priKey = pKey.(*rsa.PrivateKey)
	}

	if publicKey != "" {
		block, _ := pem.Decode([]byte(publicKey))
		var pKey interface{}
		pKey, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pubKey = pKey.(*rsa.PublicKey)
	}

	return &Rsa{
		privateKey: priKey,
		publicKey:  pubKey,
	}, nil
}

func NewRsaWithPkcs1(publicKey, privateKey string) (r *Rsa, err error) {
	var (
		pubKey *rsa.PublicKey
		priKey *rsa.PrivateKey
	)

	if privateKey != "" {
		block, _ := pem.Decode([]byte(privateKey))
		var pKey interface{}
		pKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priKey = pKey.(*rsa.PrivateKey)
	}

	if publicKey != "" {
		block, _ := pem.Decode([]byte(publicKey))
		var pKey interface{}
		pKey, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		pubKey = pKey.(*rsa.PublicKey)
	}

	return &Rsa{
		privateKey: priKey,
		publicKey:  pubKey,
	}, nil
}

/**
 * 加密
 */
func (r *Rsa) Encrypt(data []byte) ([]byte, error) {
	blockLength := r.publicKey.N.BitLen()/8 - 11
	if len(data) <= blockLength {
		return rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, data)
	}

	buffer := bytes.NewBufferString("")

	pages := len(data) / blockLength

	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(data) {
				continue
			}
			end = len(data)
		}

		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, data[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 解密
 */
func (r *Rsa) Decrypt(secretData []byte) ([]byte, error) {
	if r.publicKey == nil {
		return rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, secretData)
	}

	blockLength := r.publicKey.N.BitLen() / 8
	if len(secretData) <= blockLength {
		return rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, secretData)
	}

	buffer := bytes.NewBufferString("")

	pages := len(secretData) / blockLength
	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(secretData) {
				continue
			}
			end = len(secretData)
		}

		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, secretData[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 签名
 */
func (r *Rsa) Sign(data []byte, algorithmSign crypto.Hash) ([]byte, error) {
	hash := algorithmSign.New()
	hash.Write(data)
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, algorithmSign, hash.Sum(nil))
	if err != nil {
		return nil, err
	}
	return sign, err
}

/**
 * 验签
 */
func (r *Rsa) Verify(data []byte, sign []byte, algorithmSign crypto.Hash) bool {
	h := algorithmSign.New()
	h.Write(data)
	return rsa.VerifyPKCS1v15(r.publicKey, algorithmSign, h.Sum(nil), sign) == nil
}

/**
 * 生成pkcs1格式公钥私钥
 */
func CreatePkcs1Keys(keyLength int) (privateKey, publicKey string) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
	}))

	derPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}

	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}))
	return
}

/**
 * 生成pkcs8格式公钥私钥
 */
func CreatePkcs8Keys(keyLength int) (privateKey, publicKey string) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: MarshalPkcs8PrivateKey(rsaPrivateKey),
	}))

	derPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}

	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}))
	return
}

func MarshalPkcs8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}
