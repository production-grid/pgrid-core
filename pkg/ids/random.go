package ids

import (
	crand "crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	mrand "math/rand"
	"strings"
)

const (
	activationCodeLength = 6
	activationChars      = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
)

// Rand is a cryptographically secure mrand.Rand which uses OS entropy as a
// source.
var Rand = mrand.New(&cryptoSrc{})

/*
NewSecureID generates a random 16 byte ID in Base32 encoding.
*/
func NewSecureID() string {

	rawBytes := RandomBytes(16)
	return strings.Replace(base32.StdEncoding.EncodeToString(rawBytes), "=", "", -1)

}

/*
NewDigitCode generates a random code with digits.
*/
func NewDigitCode(length int) string {
	const digits = "1234567890"

	return randomString(length, digits)
}

/*
NewActivationCode generates a random 6 character activation code.
*/
func NewActivationCode() string {
	return randomString(activationCodeLength, activationChars)
}

/*
ValidateActivationCode determines whether or not the code contains only valid
activation code characters.
*/
func ValidateActivationCode(code string) bool {
	if len(code) != activationCodeLength {
		return false
	}

	codeRunes := []rune(code)

	for _, codeRune := range codeRunes {
		if !strings.ContainsRune(code, codeRune) {
			return false
		}
	}

	return true
}

/*
NewHmacKey generates a random key for use in computing HMAC's.
*/
func NewHmacKey() string {

	rawBytes := RandomBytes(32)
	return hex.EncodeToString(rawBytes)

}

/*
RandomBytes returns a random slice of bytes of the given length.
*/
func RandomBytes(length int) []byte {
	result := make([]byte, length)
	crand.Read(result)

	return result
}

// randomString returns a cryptographically secure random string given a pool
// of characters and a length.
func randomString(length int, pool string) string {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = pool[Rand.Intn(len(pool))]
	}

	return string(result)
}

// cryptoSrc implements the mrand.Source interface.
type cryptoSrc struct{}

func (s *cryptoSrc) Seed(seed int64) { /*Results are not deterministic, so seeding is no-op.*/
}

func (s *cryptoSrc) Uint64() (value uint64) {
	binary.Read(crand.Reader, binary.BigEndian, &value)
	return
}

func (s *cryptoSrc) Int63() int64 {
	return int64(s.Uint64() & ^(uint64(1 << 63)))
}
