package crypto

import (
	"fmt"
	"testing"
	"time"
)

func TestPasswordHashing(t *testing.T) {

	password := "Gobo"

	hash := ComputeDoubleHash(password)

	if hash.InnerSalt == "" {
		t.Error("No double hash inner salt created")
	}

	fmt.Println("Double Hash Salt:", hash.InnerSalt)

	if hash.OuterHash == "" {
		t.Error("No double hash outer hash created")
	}

	fmt.Println("Double Hash Hash:", hash.OuterHash)

	if !CompareDoubleHash("Gobo", hash) {
		t.Error("Passwords don't match")
	}

	if CompareDoubleHash("Cookies", hash) {
		t.Error("Passwords match and they shouldn't")
	}

}

func TestPasswordHashingTimeDelay(t *testing.T) {

	start := time.Now()

	SimulateDoubleHashTimeDelay()

	end := time.Now()
	elapsed := end.Sub(start)

	if elapsed < 1000 {
		t.Error("Insufficient time delay")
	}

}

func TestScrypt(t *testing.T) {

	password := "Gobo"

	hash := computeScryptHash(password)

	if hash.Salt == "" {
		t.Error("Salt not computed")
	}

	fmt.Println("Scrypt Salt:", hash.Salt)

	if hash.Hash == "" {
		t.Error("Hash not computed")
	}

	fmt.Println("Scrypt Hash:", hash.Hash)

	if !compareScryptHash("Gobo", hash) {
		t.Error("Passwords don't match")
	}

	if compareScryptHash("Cookies", hash) {
		t.Error("Passwords match and they shouldn't")
	}

}

func TestBcrypt(t *testing.T) {

	password := "Gobo"

	hash := computeBcryptHash(password)

	if hash == "" {
		t.Error("Hash not computed")
	}

	fmt.Println("Bcrypt Hash:", hash)

	if !compareBcryptHash("Gobo", hash) {
		t.Error("Passwords don't match")
	}

	if compareBcryptHash("Cookies", hash) {
		t.Error("Passwords match and they shouldn't")
	}

}
