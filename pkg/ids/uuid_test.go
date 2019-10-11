package ids

import (
	"fmt"
	"strings"
	"testing"
)

func TestActivationCode(t *testing.T) {

	code := NewActivationCode()

	if code == "" {
		t.Error("Nothing returned")
	}

	fmt.Println("Activation Code:", code)

	if len(code) != 6 {
		t.Error("Wrong length for activation code")
	}

}

func TestSecureID(t *testing.T) {

	id := NewSecureID()

	fmt.Println("Secure ID:", id, len(id))

	if strings.Contains(id, "=") {
		t.Error("Id contains padding runes")
	}

	if len(id) != 26 {
		t.Error("Invalid Base32 Encoded UUID")
	}

}

func TestUUID(t *testing.T) {

	id := NewUUID()

	fmt.Println("UUID:", id, len(id))

	if len(id) != 36 {
		t.Error("Invalid UUID")
	}

	if id == "00000000-0000-0000-0000-000000000000" {
		t.Error("empty uuid")
	}

}

func TestBase32(t *testing.T) {

	id := NewBase32UUID()

	fmt.Println("Base 32:", id, len(id))

	if strings.Contains(id, "=") {
		t.Error("Id contains padding runes")
	}

	if len(id) != 26 {
		t.Error("Invalid Base32 Encoded UUID")
	}

}

func TestBase64(t *testing.T) {

	id := NewBase64UUID()

	fmt.Println("Base 64:", id, len(id))

	if strings.Contains(id, "=") {
		t.Error("Id contains padding runes")
	}

	if len(id) != 22 {
		t.Error("Invalid Base64 Encoded UUID")
	}

}
