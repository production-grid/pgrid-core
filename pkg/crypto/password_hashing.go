package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const bcryptDifficulty = 14

type scryptParams struct {
	Cost            int
	BlockSize       int // block size in bytes
	Parallelization int
	SaltLen         int
	DKLen           int
}

// ScryptHash models the properties of an scrypt hash.
type ScryptHash struct {
	Salt string
	Hash string
}

// DoubleHash models a double password hash.
type DoubleHash struct {
	OuterHash string
	InnerSalt string
}

var defaultParams = scryptParams{Cost: 16384, BlockSize: 8, Parallelization: 1, SaltLen: 16, DKLen: 32}

// ComputeDoubleHash computes the double password hash for a new password.
func ComputeDoubleHash(password string) DoubleHash {

	var result DoubleHash

	innerHash := computeScryptHash(password)

	result.InnerSalt = innerHash.Salt

	result.OuterHash = computeBcryptHash(innerHash.Hash)

	return result
}

// CompareDoubleHash determines whether or not the given hash matches the presented
// password.
func CompareDoubleHash(password string, hash DoubleHash) bool {

	saltBytes, _ := base64.StdEncoding.DecodeString(hash.InnerSalt)

	dk, _ := scrypt.Key([]byte(password), saltBytes, defaultParams.Cost, defaultParams.BlockSize, defaultParams.Parallelization, defaultParams.DKLen)

	presentedHash := base64.StdEncoding.EncodeToString(dk)

	return compareBcryptHash(presentedHash, hash.OuterHash)

}

/*
SimulateDoubleHashTimeDelay simulates the time delay that would be associated
with computing two hashes.  Used to prevent timing attacks from determining
which user id's are valid.
*/
func SimulateDoubleHashTimeDelay() {

	simulatedPassword := generateScryptSalt()

	saltBytes := generateScryptSalt()

	dk, _ := scrypt.Key(simulatedPassword, saltBytes, defaultParams.Cost, defaultParams.BlockSize, defaultParams.Parallelization, defaultParams.DKLen)

	presentedHash := base64.StdEncoding.EncodeToString(dk)

	computeBcryptHash(presentedHash)

}

func computeScryptHash(password string) ScryptHash {

	var hash ScryptHash

	saltBytes := generateScryptSalt()

	hash.Salt = base64.StdEncoding.EncodeToString(saltBytes)

	dk, _ := scrypt.Key([]byte(password), saltBytes, defaultParams.Cost, defaultParams.BlockSize, defaultParams.Parallelization, defaultParams.DKLen)

	hash.Hash = base64.StdEncoding.EncodeToString(dk)

	return hash

}

func compareScryptHash(password string, hash ScryptHash) bool {

	saltBytes, _ := base64.StdEncoding.DecodeString(hash.Salt)

	dk, _ := scrypt.Key([]byte(password), saltBytes, defaultParams.Cost, defaultParams.BlockSize, defaultParams.Parallelization, defaultParams.DKLen)

	presentedHash := base64.StdEncoding.EncodeToString(dk)

	return hash.Hash == presentedHash
}

func generateScryptSalt() []byte {

	b := make([]byte, defaultParams.SaltLen)
	rand.Read(b)
	return b

}

func computeBcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcryptDifficulty)
	return base64.StdEncoding.EncodeToString(bytes)
}

func compareBcryptHash(password string, hash string) bool {
	rawHash, _ := base64.StdEncoding.DecodeString(hash)
	err := bcrypt.CompareHashAndPassword(rawHash, []byte(password))
	return err == nil
}
