package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"service/internal/pkg/config"
	"strings"
)

type PrivateAPIEncryption interface {
	Encrypt(key string) (string, error)
	Decrypt(name, publicKey string) (string, error)
}

func NewPrivateAPIEncryption(id string) PrivateAPIEncryption {
	return &privateAPIEncryption{
		id: id,
	}
}

type privateAPIEncryption struct {
	id string
}

func (ec *privateAPIEncryption) Encrypt(key string) (string, error) {
	rawKey := ec.padKey(fmt.Sprintf("%s%f", key, ec.id), 16)

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", fmt.Errorf("cipher error: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("IV error: %w", err)
	}

	plaintext, err := bcrypt.GenerateFromPassword([]byte(key), 10)
	if err != nil {
		return "", fmt.Errorf("bcrypt key error: %w", err)
	}

	plainBytes := ec.pkcs7Pad(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plainBytes)

	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)

	full := append(iv, []byte(ciphertextBase64)...)
	finalEncoded := base64.StdEncoding.EncodeToString(full)

	return "_open:" + finalEncoded, nil
}

func (ec *privateAPIEncryption) Decrypt(name, publicKey string) (string, error) {
	credential, ok := config.PrivateAPICredential[name]
	if !ok {
		return "", fmt.Errorf("private api credential not found")
	}

	if strings.HasPrefix(publicKey, "_open:") {
		publicKey = strings.TrimPrefix(publicKey, "_open:")
	}

	raw, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}
	if len(raw) < aes.BlockSize {
		return "", fmt.Errorf("data too short")
	}

	iv := raw[:aes.BlockSize]
	ciphertextBase64 := raw[aes.BlockSize:]

	ciphertext, err := base64.StdEncoding.DecodeString(string(ciphertextBase64))
	if err != nil {
		return "", fmt.Errorf("decode nested base64 error: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	rawKey := ec.padKey(fmt.Sprintf("%s%f", credential["key"].(string), ec.id), 16)
	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return "", fmt.Errorf("aes cipher error: %w", err)
	}

	decrypted := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, ciphertext)

	result, err := ec.pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpad error: %w", err)
	}

	return string(result), nil
}

func (ec *privateAPIEncryption) padKey(key string, length int) []byte {
	if len(key) >= length {
		return []byte(key[:length])
	}
	return append([]byte(key), bytes.Repeat([]byte("0"), length-len(key))...)
}

func (ec *privateAPIEncryption) pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data")
	}
	pad := int(data[length-1])
	if pad == 0 || pad > blockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:length-pad], nil
}

func (ec *privateAPIEncryption) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...)
}
