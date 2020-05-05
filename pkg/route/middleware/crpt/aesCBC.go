package crpt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"

	"github.com/pkg/errors"
)

//AesEncryptCBC nil
type AesEncryptCBC struct {
}

//AesEncrypt encrypt
func AesEncrypt(origData []byte) ([]byte, error) {
	aesEnc := AesEncryptCBC{}
	key := []byte("1234567812345678")

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "AesEncrypt")
	}

	blockSize := block.BlockSize()
	origData = aesEnc.pkcs7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypt := make([]byte, len(origData))
	blockMode.CryptBlocks(crypt, origData)

	return crypt, nil
}

//AesDecrypt decrypt
func AesDecrypt(cryptKey []byte) ([]byte, error) {
	aesEnc := AesEncryptCBC{}
	key := []byte("1234567812345678")

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "AesDecrypt")
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cryptKey))
	blockMode.CryptBlocks(origData, cryptKey)
	origData = aesEnc.pkcs7UnPadding(origData)

	return origData, nil
}

func (a *AesEncryptCBC) pkcs7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (a *AesEncryptCBC) pkcs7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	return plantText[:(length - unPadding)]
}
