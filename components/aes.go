package components

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"github.com/grpc-boot/base/v3/utils"
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
	return NewAesWithBytes(utils.String2Bytes(key), utils.String2Bytes(iv))
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
	data := pkcs7Padding(plain, a.block.BlockSize())

	secretData = make([]byte, len(data))
	cipher.NewCBCEncrypter(a.block, a.iv).CryptBlocks(secretData, data)

	return
}

// CbcDecrypt cbc解密
func (a *Aes) CbcDecrypt(secretData []byte) (data []byte, err error) {
	if len(secretData) == 0 || len(secretData)%a.block.BlockSize() != 0 {
		err = ErrAesDecrypt
		return
	}

	data = make([]byte, len(secretData))
	cipher.NewCBCDecrypter(a.block, a.iv).CryptBlocks(data, secretData)

	return pkcs7UnPadding(data)
}

func pkcs7Padding(src []byte, blockSize int) []byte {
	padChar := byte(blockSize - len(src)%blockSize)
	return append(src, bytes.Repeat([]byte{padChar}, int(padChar))...)
}

func pkcs7UnPadding(src []byte) ([]byte, error) {
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
