package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io"
	"math/rand"
	"strings"
	"time"
)

// MapToQuery will take in a map[string]string, and turn it into a query string
func MapToQuery(inMap map[string]string, start string) (res string) {
	if inMap == nil {
		return
	}
	var list []string
	for key, value := range inMap {
		list = append(list, key+"="+value)
	}
	res = strings.Join(list, "&")
	res = start + res
	return
}

// GenerateRandomString generates a random string for all kinds of purposes
// n - length of the generated string
// full - use the special characters or not.
func GenerateRandomString(n int, full bool) string {
	rand.Seed(time.Now().UnixNano())
	var LetterRunes = []rune{}
	if full {
		LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ{}!@#$^&*()~{}|><")
	} else {
		LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = LetterRunes[rand.Intn(len(LetterRunes))]
	}
	return string(b)
}

// Encrypt takes in a bytes slice and encrypts it using a key
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(crand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt reverses the Encrypt functionality.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Serialize takes in an object and serializes it into a bytes slice
func Serialize(obj interface{}) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(obj)
	if err != nil {
		return []byte{}, err
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(buffer.Len()))
	base64.StdEncoding.Encode(encoded, buffer.Bytes())
	return encoded, nil
}

// Unserialize reverses the Serialize functionality
func Unserialize(encodedBytes *[]byte, obj interface{}) error {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(*encodedBytes)))
	_, err := base64.StdEncoding.Decode(decoded, *encodedBytes)
	if err != nil {
		return err
	}
	buffer := bytes.Buffer{}
	buffer.Write(decoded)
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}
