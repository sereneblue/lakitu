package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/alexedwards/argon2id"
	"github.com/shengdoushi/base58"
	"golang.org/x/crypto/argon2"
)

func decrypt(enc, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(hash(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(enc) < gcm.NonceSize() {
		return nil, errors.New("malformed encrypted data")
	}

	return gcm.Open(nil,
		enc[:gcm.NonceSize()],
		enc[gcm.NonceSize():],
		nil,
	)
}

func encrypt(text, key []byte) (string, error) {
	block, err := aes.NewCipher(hash(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	return string(gcm.Seal(nonce, nonce, text, nil)), nil
}

func hash(key []byte) []byte {
	h := sha256.New()
	h.Write(key)
	return h.Sum(nil)
}

func CreateHashWithSalt(password string, salt []byte, params *argon2id.Params) string {
	key := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Key)
	return hash
}

func Encrypt(str, key string) (string, error) {
	return encrypt([]byte(str), []byte(key))
}

func Decrypt(enc, key string) (string, error) {
	data, err := decrypt([]byte(enc), []byte(key))
	return string(data), err
}

func GenerateRandomKey() (string, error) {
	b := make([]byte, 128)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	alphabet := base58.NewAlphabet("ABCDEFGHJKLMNPQRSTUVWXYZ123456789abcdefghijkmnopqrstuvwxyz")

	return base58.Encode(b, alphabet), err
}
