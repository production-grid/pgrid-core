package ids

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

var (
	initializerSeeder sync.Once
	uuidChannel       chan uuid.UUID
)

/*
Identified is used for anything that has a string identifier.
*/
type Identified interface {
	Identifier() string
}

/*
NewUUID generates a string encoded UUID in the conventional format.
*/
func NewUUID() string {
	u := NewRawUUID()
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// doGenerateUuids creates an indefinite supply of V1 UUIDs and puts them into
// a blocking channel. This was intended to address UUID duplication.
// NOTE: This may no longer be necessary since replacing the
// package github.com/satori/go.uuid
// See: https://github.com/satori/go.uuid/issues/73
func doGenerateUuids() {
	for {
		u, err := uuid.NewV1()
		if err != nil {
			logging.Errorf("Failed to generate UUID: %v", err)
			continue
		}
		uuidChannel <- u
	}
}

/*
NewRawUUID generates a new raw UUID.
*/
func NewRawUUID() uuid.UUID {

	initializerSeeder.Do(func() {
		time.Sleep(1 * time.Second)
		uuidChannel = make(chan uuid.UUID, 1000)
		go doGenerateUuids()
	})

	return <-uuidChannel

}

/*
NewBase32UUID generates a new Base32 encoded UUID.
*/
func NewBase32UUID() string {

	u := NewRawUUID()
	rawBytes := u[0:]
	return strings.Replace(base32.StdEncoding.EncodeToString(rawBytes), "=", "", -1)

}

/*
NewBase64UUID generates a new Base64 encoded UUID.
*/
func NewBase64UUID() string {

	u := NewRawUUID()
	rawBytes := u[0:]
	return strings.Replace(base64.StdEncoding.EncodeToString(rawBytes), "=", "", -1)

}
