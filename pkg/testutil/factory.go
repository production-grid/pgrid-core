package testutil

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/production-grid/pgrid-core/pkg/money"
	"github.com/stretchr/testify/assert"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "0123456789"

// PopulateTestData fills an entity with dummy data.
func PopulateTestData(entity interface{}) interface{} {

	el := reflect.ValueOf(entity).Elem()
	t := el.Type()

	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		column := fld.Tag.Get("col")
		if column == "" {
			continue
		}
		if column == "id" {
			continue
		}
		val := el.FieldByIndex(fld.Index)
		v := val.Interface()
		switch v.(type) {
		case string:
			val.SetString(GenerateTestString(16))
		case bool:
			val.SetBool(true)
		case time.Time:
			val.Set(reflect.ValueOf(time.Now()))
		case *time.Time:
			n := time.Now()
			val.Set(reflect.ValueOf(&n))
		}
	}

	return entity
}

// AssertEquivalent ensures two entities are populated with the same data.
func AssertEquivalent(t *testing.T, actual interface{}, expected interface{}) {

	assert := assert.New(t)

	valExpected := reflect.ValueOf(expected)
	valActual := reflect.ValueOf(actual)
	tp := reflect.ValueOf(expected).Elem().Type()

	for i := 0; i < tp.NumField(); i++ {
		fld := tp.Field(i)
		column := fld.Tag.Get("col")
		if column == "" {
			continue
		}
		if column == "id" {
			continue
		}
		fldExpected := valExpected.Elem().FieldByIndex(fld.Index).Addr().Interface()
		fldActual := valActual.Elem().FieldByIndex(fld.Index).Addr().Interface()
		switch fldExpected.(type) {
		case *string:
			e := fldExpected.(*string)
			a := fldActual.(*string)
			assert.Equal(a, e)
		case *bool:
			e := fldExpected.(*bool)
			a := fldActual.(*bool)
			assert.Equal(a, e)
		case *time.Time:
			e := fldExpected.(*time.Time)
			a := fldActual.(*time.Time)
			assert.True(a.Equal(*e))
		case **time.Time:
			e := fldExpected.(**time.Time)
			a := fldActual.(**time.Time)
			assert.True((*a).Equal(**e))
		}
	}

}

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
