package controllers

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GeneratePassword(raw string) string {
	rand.Seed(time.Now().UnixNano())
	salt := []byte(randSeq(10))

	var ptime uint32 = 1
	var memory uint32 = 64 * 1024
	var threads uint8 = 4
	var length uint32 = 64

	hash := argon2.IDKey([]byte(raw), salt, ptime, memory, threads, length)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	return fmt.Sprintf(format, argon2.Version, memory, ptime, threads, b64Salt, b64Hash)
}

func ComparePasswords(password, hash string) bool {
	parts := strings.Split(hash, "$")

	var ptime uint32
	var memory uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &ptime, &threads)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	length := uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, ptime, memory, threads, length)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}
