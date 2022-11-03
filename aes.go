package base

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	ErrInvalidPaddingChar = errors.New(`invalid padding char`)
	ErrAesDecrypt         = errors.New(`aes decrypt error`)
)

// Aes Aes加密
type Aes struct {
	iv    []byte
	block cipher.Block
}

// NewAes 实例化Aes
func NewAes(key, iv string) (a *Aes, err error) {
	var b cipher.Block
	b, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if iv == "" {
		iv = key[:b.BlockSize()]
	}

	a = &Aes{
		iv:    []byte(iv),
		block: b,
	}
	return
}

// NewAesWithBytes 实例化Aes
func NewAesWithBytes(key, iv []byte) (a *Aes, err error) {
	var b cipher.Block
	b, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	a = &Aes{
		iv:    iv,
		block: b,
	}
	return
}

// CbcEncrypt cbc加密
func (a *Aes) CbcEncrypt(plain []byte) (secretData []byte) {
	data := pkcs5Padding(plain, a.block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(a.block, a.iv)
	secretData = make([]byte, len(data))
	blockMode.CryptBlocks(secretData, data)
	return
}

// CbcDecrypt cbc解密
func (a *Aes) CbcDecrypt(secretData []byte) (data []byte, err error) {
	defer func() {
		er := recover()
		if er != nil {
			err = errors.New(er.(string))
			data = nil
		}
	}()

	blockMode := cipher.NewCBCDecrypter(a.block, a.iv)
	data = make([]byte, len(secretData))
	blockMode.CryptBlocks(data, secretData)
	return pkcs5UnPadding(data)
}

func pkcs5Padding(src []byte, blockSize int) []byte {
	padChar := byte(blockSize - len(src)%blockSize)
	return append(src, bytes.Repeat([]byte{padChar}, int(padChar))...)
}

func pkcs5UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, ErrAesDecrypt
	}

	padChar := int(src[length-1])
	if length-padChar < 0 {
		return nil, ErrInvalidPaddingChar
	}

	return src[:length-padChar], nil
}
