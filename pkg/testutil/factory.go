package testutil

import (
	"math/rand"
	"time"

	"github.com/production-grid/pgrid-core/pkg/money"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "0123456789"

/*
GenerateTestString generates a random string of the given length.
*/
func GenerateTestString(length int) string {

	rand.Seed(int64(time.Now().Nanosecond()))

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)

}

/*
GenerateTestNumbers generates a random string containing only numbers.
*/
func GenerateTestNumbers(length int) string {

	rand.Seed(int64(time.Now().Nanosecond()))

	b := make([]byte, length)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}
	return string(b)

}

/*
GenerateTestPassword generates a random password guaranteed to meet PCI password
requirements.
*/
func GenerateTestPassword() string {
	return GenerateTestString(4) + GenerateTestNumbers(4)
}

/*
GenerateTestCurrency generates a test currency amount.
*/
func GenerateTestCurrency(max int) money.Currency {

	rand.Seed(int64(time.Now().Nanosecond()))

	result := money.Currency{CurrencyCode: "USD"}

	result.Integer = rand.Intn(max)
	result.Decimal = rand.Intn(100)
	result.DecimalDenominator = 100

	return result

}

/*
GenerateTestBasisPoints generates a test percentage in basis points.
*/
func GenerateTestBasisPoints(maxBps int) money.BasisPoints {

	rand.Seed(int64(time.Now().Nanosecond()))

	actualBps := rand.Intn(maxBps)

	return *money.NewBasisPoints(actualBps)

}

/*
GenerateTestEmail generates a random email address for testing.
*/
func GenerateTestEmail(unamelen int) string {

	return "test-" + GenerateTestString(unamelen) + "@blockchyp.com"

}
