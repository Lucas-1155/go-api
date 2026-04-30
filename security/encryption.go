package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func getSecretKey() []byte {
	key := os.Getenv("API_SECRET_KEY")
	// Validação de segurança: AES-256 exige 32 bytes
	if len(key) != 32 {
		fmt.Printf("AVISO: A chave API_SECRET_KEY deve ter 32 bytes! Atual: %d\n", len(key))
	}
	return []byte(key)
}

func Encrypt(text string) (string, error) {
	block, _ := aes.NewCipher(getSecretKey())
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(cryptoText)
	block, _ := aes.NewCipher(getSecretKey())
	gcm, _ := cipher.NewGCM(block)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}
