package money

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntDivision(t *testing.T) {

	tests := []struct {
		name        string
		base        Currency
		operand     int
		expect      *Currency
		shouldPanic bool
	}{
		{
			name:    "SimpleDivision",
			base:    Currency{Integer: 3, Decimal: 33, CurrencyCode: "USD"},
			operand: 3,
			expect:  &Currency{Integer: 1, Decimal: 11, CurrencyCode: "USD"},
		},
		{
			name:    "CarryOverDivision",
			base:    Currency{Integer: 4, Decimal: 33, CurrencyCode: "USD"},
			operand: 3,
			expect:  &Currency{Integer: 1, Decimal: 44, CurrencyCode: "USD"},
		},
		{
			name:    "SmallNegativeDivision",
			base:    Currency{Integer: 0, Decimal: 12, Negative: true, CurrencyCode: "USD"},
			operand: 2,
			expect:  &Currency{Integer: 0, Decimal: 6, Negative: true, CurrencyCode: "USD"},
		},
		{
			name:    "SmallDoubleNegativeDivision",
			base:    Currency{Integer: 0, Decimal: 12, Negative: true, CurrencyCode: "USD"},
			operand: -2,
			expect:  &Currency{Integer: 0, Decimal: 6, Negative: false, CurrencyCode: "USD"},
		},
		{
			name:    "SmallNegativeIntegerDivision",
			base:    Currency{Integer: 0, Decimal: 12, Negative: false, CurrencyCode: "USD"},
			operand: -2,
			expect:  &Currency{Integer: 0, Decimal: 6, Negative: true, CurrencyCode: "USD"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() { test.base.DivideInt(test.operand, RoundNearest) })
				return
			}

			result := test.base.DivideInt(test.operand, RoundNearest)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}

}

func TestRounding(t *testing.T) {
	tests := []struct {
		name          string
		base          Currency
		decimalPlaces int
		roundingMode  RoundingMode
		expect        *Currency
	}{
		{
			name:          "RoundUp",
			base:          Currency{Integer: 5, Decimal: 505, DecimalDenominator: 1000, CurrencyCode: CurrencyCodeUSD},
			decimalPlaces: 2,
			roundingMode:  RoundUp,
			expect:        &Currency{Integer: 5, Decimal: 51, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:          "RoundNearest",
			base:          Currency{Integer: 5, Decimal: 505, DecimalDenominator: 1000, CurrencyCode: CurrencyCodeUSD},
			decimalPlaces: 2,
			roundingMode:  RoundNearest,
			expect:        &Currency{Integer: 5, Decimal: 51, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:          "RoundDown",
			base:          Currency{Integer: 5, Decimal: 505, DecimalDenominator: 1000, CurrencyCode: CurrencyCodeUSD},
			decimalPlaces: 2,
			roundingMode:  RoundDown,
			expect:        &Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:          "RoundDownSmallNegative",
			base:          Currency{Integer: 0, Decimal: 505, DecimalDenominator: 1000, Negative: true, CurrencyCode: CurrencyCodeUSD},
			decimalPlaces: 2,
			roundingMode:  RoundDown,
			expect:        &Currency{Integer: 0, Decimal: 51, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:          "RoundUpSmallNegative",
			base:          Currency{Integer: 0, Decimal: 505, DecimalDenominator: 1000, Negative: true, CurrencyCode: CurrencyCodeUSD},
			decimalPlaces: 2,
			roundingMode:  RoundUp,
			expect:        &Currency{Integer: 0, Decimal: 50, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			result := test.base.Round(test.decimalPlaces, test.roundingMode)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}

}

func TestAddition(t *testing.T) {
	tests := []struct {
		name        string
		base        Currency
		operand     Currency
		expect      *Currency
		shouldPanic bool
	}{
		{
			name:    "SimpleAddition",
			base:    Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 5, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 11, Decimal: 25, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "DecimalMagic",
			base:    Currency{Integer: 5, Decimal: 05, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 5, Decimal: 05, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 10, Decimal: 10, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroBase",
			base:    Currency{},
			operand: Currency{Integer: 5, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 5, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroBaseLeadingZero",
			base:    Currency{},
			operand: Currency{Integer: 5, Decimal: 05, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 5, Decimal: 05, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroOperand",
			base:    Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{},
			expect:  &Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "BigLeftDenom",
			base:    Currency{Integer: 5, Decimal: 50, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 5, Decimal: 12, DecimalDenominator: 1000, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 10, Decimal: 512, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "BigRightDenom",
			base:    Currency{Integer: 5, Decimal: 12, DecimalDenominator: 1000, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 5, Decimal: 50, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 10, Decimal: 512, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroZero",
			base:    Currency{},
			operand: Currency{},
			expect:  &Currency{},
		},
		{
			name:    "LeftSmallNegative",
			base:    Currency{Integer: 5, Decimal: 12, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 5, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "RightSmallNegative",
			base:    Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 5, Decimal: 12, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 5, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "DoubleSmallNegative",
			base:    Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 0, Decimal: 24, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() { test.base.Add(&test.operand) })
				return
			}

			result := test.base.Add(&test.operand)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}
}

func TestSubtraction(t *testing.T) {
	tests := []struct {
		name        string
		base        Currency
		operand     Currency
		expect      *Currency
		shouldPanic bool
	}{
		{
			name:    "SimpleSubtraction",
			base:    Currency{Integer: 5, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 0, Decimal: 25, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroBase",
			base:    Currency{},
			operand: Currency{Integer: 5, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: -5, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroOperand",
			base:    Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{},
			expect:  &Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroZero",
			base:    Currency{},
			operand: Currency{},
			expect:  &Currency{},
		},
		{
			name:    "LeftSmallNegative",
			base:    Currency{Integer: 5, Decimal: 12, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 5, Decimal: 24, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "RightSmallNegative",
			base:    Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 5, Decimal: 12, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 5, Decimal: 24, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "DoubleSmallNegative",
			base:    Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 0, Decimal: 12, Negative: true, DecimalDenominator: 100, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 0, Decimal: 0, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() { test.base.Subtract(&test.operand) })
				return
			}

			result := test.base.Subtract(&test.operand)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}
}

func TestFractionalPennyMultiplication(t *testing.T) {

	result := ZeroDollars()

	base := Currency{}
	base.Integer = 0
	base.CurrencyCode = "USD"
	base.Decimal = 18
	base.DecimalDenominator = 10000

	result = result.Add(base.MultInt(3113))

	assert.Equal(t, "5.6034", result.FormatCurrency())

}

func TestIntMultiplication(t *testing.T) {
	tests := []struct {
		name        string
		base        Currency
		operand     int
		expect      *Currency
		shouldPanic bool
	}{
		{
			name:    "SimpleMultiplication",
			base:    Currency{Integer: 2, Decimal: 55, CurrencyCode: CurrencyCodeUSD},
			operand: 3,
			expect:  &Currency{Integer: 7, Decimal: 65, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroBase",
			base:    Currency{CurrencyCode: CurrencyCodeUSD},
			operand: 5,
			expect:  &Currency{Integer: 0, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroOperand",
			base:    Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
			operand: 0,
			expect:  &Currency{Integer: 0, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "NinetyFiveTimesThree",
			base:    Currency{Integer: 0, Decimal: 95, CurrencyCode: CurrencyCodeUSD},
			operand: 3,
			expect:  &Currency{Integer: 2, Decimal: 85, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "NegativeMultiplication",
			base:    Currency{Integer: 1, Decimal: 25, Negative: true, CurrencyCode: CurrencyCodeUSD},
			operand: 3,
			expect:  &Currency{Integer: 3, Decimal: 75, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "FracPenny",
			base:    *ParseCurrency("USD", "0.0018"),
			operand: 3113,
			expect:  &Currency{Integer: 5, Decimal: 6034, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SmallNegativeBase",
			base:    *ParseCurrency("USD", "-0.10"),
			operand: 5,
			expect:  &Currency{Integer: 0, Decimal: 50, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SmallNegativeOperand",
			base:    *ParseCurrency("USD", "0.10"),
			operand: -5,
			expect:  &Currency{Integer: 0, Decimal: 50, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SmallDoubleNegative",
			base:    *ParseCurrency("USD", "-0.10"),
			operand: -5,
			expect:  &Currency{Integer: 0, Decimal: 50, Negative: false, CurrencyCode: CurrencyCodeUSD},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() { test.base.MultInt(test.operand) })
				return
			}

			result := test.base.MultInt(test.operand)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}
}

func TestMultiplication(t *testing.T) {
	tests := []struct {
		name        string
		base        Currency
		operand     Currency
		expect      *Currency
		shouldPanic bool
	}{
		{
			name:    "SimpleBase10Multiplication",
			base:    Currency{Integer: 10, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 2, Decimal: 00, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 20, Decimal: 00, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SimpleMultiplication",
			base:    Currency{Integer: 2, Decimal: 55, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 3, Decimal: 33, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 8, Decimal: 49, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroBase",
			base:    Currency{},
			operand: Currency{Integer: 5, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 0, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "ZeroOperand",
			base:    Currency{Integer: 5, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{},
			expect:  &Currency{Integer: 0, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:        "ZeroZero",
			base:        Currency{},
			operand:     Currency{},
			shouldPanic: true,
		},
		{
			name:    "SmallNegativeBase",
			base:    Currency{Integer: 0, Decimal: 75, Negative: true, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 0, Decimal: 25, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 0, Decimal: 18, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SmallNegativeOperand",
			base:    Currency{Integer: 0, Decimal: 75, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 0, Decimal: 25, Negative: true, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 0, Decimal: 18, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SmallDoubleNegative",
			base:    Currency{Integer: 0, Decimal: 75, Negative: true, CurrencyCode: CurrencyCodeUSD},
			operand: Currency{Integer: 0, Decimal: 25, Negative: true, CurrencyCode: CurrencyCodeUSD},
			expect:  &Currency{Integer: 0, Decimal: 18, CurrencyCode: CurrencyCodeUSD},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() { test.base.Mult(&test.operand) })
				return
			}

			result := test.base.Mult(&test.operand)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}
}

func TestFloatConversions(t *testing.T) {
	tests := []struct {
		name  string
		input float64

		expectInteger  int
		expectDecimal  int
		expectNegative bool
	}{
		{
			name:          "IntegerAndDecimal",
			input:         23.23,
			expectInteger: 23,
			expectDecimal: 23,
		},
		{
			name:          "OneSignificantDecimal",
			input:         0.5,
			expectInteger: 0,
			expectDecimal: 50,
		},
		{
			name:          "LeadingZeros",
			input:         0.05,
			expectInteger: 0,
			expectDecimal: 5,
		},
		{
			name:          "LeadingZeroArbitraryPrecision",
			input:         0.051,
			expectInteger: 0,
			expectDecimal: 51,
		},
		{
			name:           "SmallNegative",
			input:          -0.51,
			expectInteger:  0,
			expectDecimal:  51,
			expectNegative: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Input: %v", test.input)

			result := FromFloat64(CurrencyCodeUSD, test.input)

			t.Logf("Output: %s", result.Localize())

			assert.Equal(t, test.expectInteger, result.Integer)
			assert.Equal(t, test.expectDecimal, result.Decimal)

			assert.Equal(t, test.input, result.ToFloat())
		})
	}
}

func TestLocalize(t *testing.T) {
	tests := []struct {
		name  string
		input string

		expect string
	}{
		{
			name:   "IntegerAndDecimal",
			input:  "23.23",
			expect: "$23.23",
		},
		{
			name:   "OneSignificantDecimal",
			input:  "1",
			expect: "$1.00",
		},
		{
			name:   "PrecisionOverflowRoundUp",
			input:  "1.0099",
			expect: "$1.01",
		},
		{
			name:   "PrecisionOverflowRoundDown",
			input:  "1.0001",
			expect: "$1.00",
		},
		{
			name:   "PrecisionOverflowNoZerosRoundDown",
			input:  "1.111111",
			expect: "$1.11",
		},
		{
			name:   "PrecisionOverflowNoZerosRoundUp",
			input:  "1.999999",
			expect: "$2.00",
		},
		{
			name:   "SmallNegative",
			input:  "-0.12",
			expect: "$-0.12",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Input: %v", test.input)

			result := ParseCurrency(CurrencyCodeUSD, test.input)

			t.Logf("Output: %s", result.Localize())

			assert.Equal(t, test.expect, result.Localize())
		})
	}
}

func TestParseCurrency(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		output         string
		expectInteger  int
		expectDecimal  int
		expectNegative bool
	}{
		{
			name:          "IntegerAndDecimal",
			input:         "23.23",
			output:        "23.23",
			expectInteger: 23,
			expectDecimal: 23,
		},
		{
			name:           "IntegerAndDecimal",
			input:          "-34.07",
			output:         "-34.07",
			expectInteger:  34,
			expectDecimal:  7,
			expectNegative: true,
		},
		{
			name:          "OneSignificantDecimal",
			input:         "0.5",
			output:        "0.50",
			expectInteger: 0,
			expectDecimal: 50,
		},
		{
			name:          "OneSignificantInteger",
			input:         "1",
			output:        "1.0",
			expectInteger: 1,
		},
		{
			name:          "ThousandsComma",
			input:         "10,000.00",
			output:        "10000.00",
			expectInteger: 10000,
		},
		{
			name:           "SmallNegative",
			input:          "-0.11",
			output:         "-0.11",
			expectInteger:  0,
			expectDecimal:  11,
			expectNegative: true,
		},
		{
			name:          "WithSymbol",
			input:         "$123.11",
			output:        "123.11",
			expectInteger: 123,
			expectDecimal: 11,
		},
		{
			name:           "WithSymbolNegative",
			input:          "-$123.11",
			output:         "-123.11",
			expectInteger:  123,
			expectDecimal:  11,
			expectNegative: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("Input: %v", test.input)

			result := ParseCurrency(CurrencyCodeUSD, test.input)

			t.Logf("Output: %s", result.String())

			value, err := result.Value()
			assert.NoError(t, err)
			assert.Equal(t, test.expectInteger, result.Integer)
			assert.Equal(t, test.expectDecimal, result.Decimal)
			assert.Equal(t, test.output, value.(string))
			assert.Equal(t, test.expectNegative, result.Negative)
		})
	}
}

func TestBasisPoints(t *testing.T) {

	assert := assert.New(t)

	o1 := Currency{Integer: 5, Decimal: 55, CurrencyCode: CurrencyCodeUSD}

	o2, err := ParseBasisPoints("20%")

	assert.NoError(err)

	o3 := o1.MultBPS(o2)

	fmt.Println(o3.FormatCurrency())

	if o3.Integer != 1 {
		t.Error("Integer result incorrect")
	}

	if o3.Decimal != 11 {
		t.Error("Decimal result incorrect")
	}

}

func TestPercentage(t *testing.T) {

	o1 := Currency{Integer: 5, Decimal: 55, CurrencyCode: CurrencyCodeUSD}

	o3 := o1.Pct(20)

	fmt.Println(o3.FormatCurrency())

	if o3.Integer != 1 {
		t.Error("Integer result incorrect")
	}

	if o3.Decimal != 11 {
		t.Error("Decimal result incorrect")
	}

}

func TestParsing(t *testing.T) {

	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "0.035",
			expect: "0.035",
		},
		{
			input:  "0.05",
			expect: "0.05",
		},
		{
			input:  "0.50",
			expect: "0.50",
		},
		{
			input:  "5.55",
			expect: "5.55",
		},
		{
			input:  "5.5589",
			expect: "5.5589",
		},
		{
			input:  "34.07",
			expect: "34.07",
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := ParseCurrency(CurrencyCodeUSD, test.input)

			assert.Equal(t, test.expect, result.FormatCurrency())
		})
	}

}

type marshallerTestStruct struct {
	Money Currency `json:"money"`
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "SingleIntegerAsString",
			input:  "{\"money\":\"1\"}",
			expect: "1.00",
		},
		{
			name:   "Float",
			input:  "{\"money\":1.00}",
			expect: "1.00",
		},
		{
			name:   "Int",
			input:  "{\"money\":1}",
			expect: "1.00",
		},
		{
			name:   "FloatAsString",
			input:  "{\"money\":\"1.00\"}",
			expect: "1.00",
		},
		{
			name:   "NegativeFloatAsString",
			input:  "{\"money\":\"-1.00\"}",
			expect: "-1.00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result = &marshallerTestStruct{}
			if err := json.Unmarshal([]byte(test.input), result); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expect, result.Money.String())
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  marshallerTestStruct
		expect string
	}{
		{
			name: "Positive",
			input: marshallerTestStruct{
				Money: Currency{
					Integer:      1,
					CurrencyCode: CurrencyCodeUSD,
				},
			},
			expect: "{\"money\":\"1.00\"}",
		},
		{
			name: "Negative",
			input: marshallerTestStruct{
				Money: Currency{
					Integer:      -1,
					CurrencyCode: CurrencyCodeUSD,
				},
			},
			expect: "{\"money\":\"-1.00\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := json.Marshal(test.input)
			assert.NoError(t, err)

			assert.Equal(t, test.expect, string(result))
		})
	}
}

func TestMultBPS(t *testing.T) {
	tests := []struct {
		name    string
		base    *Currency
		operand *BasisPoints
		expect  *Currency
	}{
		{
			name:    "RoundUp",
			base:    &Currency{Integer: 1888, Decimal: 9, CurrencyCode: CurrencyCodeUSD},
			operand: NewBasisPoints(154),
			expect:  &Currency{Integer: 29, Decimal: 8, CurrencyCode: CurrencyCodeUSD},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.base.MultBPS(test.operand)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}

}

func TestMultFloat(t *testing.T) {
	tests := []struct {
		name    string
		base    *Currency
		operand float64
		round   RoundingMode
		expect  *Currency
	}{
		{
			name:    "RoundNumber",
			base:    &Currency{Integer: 10, Decimal: 50, CurrencyCode: CurrencyCodeUSD},
			operand: 2,
			round:   RoundUp,
			expect:  &Currency{Integer: 21, Decimal: 00, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "RoundNearest",
			base:    &Currency{Integer: 2, Decimal: 54, CurrencyCode: CurrencyCodeUSD},
			operand: 1.6666666666666666,
			round:   RoundNearest,
			expect:  &Currency{Integer: 4, Decimal: 23, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "RoundUp",
			base:    &Currency{Integer: 2, Decimal: 54, CurrencyCode: CurrencyCodeUSD},
			operand: 1.6666666666666666,
			round:   RoundUp,
			expect:  &Currency{Integer: 4, Decimal: 24, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "Zero",
			base:    &Currency{},
			operand: 1.00,
			round:   RoundUp,
			expect:  Zero(),
		},
		{
			name:    "SmallNegativeBase",
			base:    &Currency{Integer: 2, Decimal: 54, Negative: true, CurrencyCode: CurrencyCodeUSD},
			operand: 1.6666666666666666,
			round:   RoundUp,
			expect:  &Currency{Integer: 4, Decimal: 23, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SmallNegativeOperand",
			base:    &Currency{Integer: 2, Decimal: 54, CurrencyCode: CurrencyCodeUSD},
			operand: -1.6666666666666666,
			round:   RoundUp,
			expect:  &Currency{Integer: 4, Decimal: 23, Negative: true, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:    "SmallDoubleNegative",
			base:    &Currency{Integer: 2, Decimal: 54, Negative: true, CurrencyCode: CurrencyCodeUSD},
			operand: -1.6666666666666666,
			round:   RoundUp,
			expect:  &Currency{Integer: 4, Decimal: 24, CurrencyCode: CurrencyCodeUSD},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.base.MultFloat(test.operand, test.round)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}
}

func TestParseNumericCode(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect *Currency
		panics bool
	}{
		{
			name:   "CentsOnly",
			input:  "000000000010",
			expect: &Currency{Integer: 0, Decimal: 10, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:   "OneDollar",
			input:  "000000000100",
			expect: &Currency{Integer: 1, Decimal: 0, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:   "OneThousandDollarsNinetyNineCents",
			input:  "000000100099",
			expect: &Currency{Integer: 1000, Decimal: 99, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:   "ShortInput",
			input:  "98",
			expect: &Currency{Integer: 0, Decimal: 98, CurrencyCode: CurrencyCodeUSD},
		},
		{
			name:   "BadInput",
			input:  "BADINPUT",
			panics: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.panics {
				assert.Panics(t, func() {
					ParseNumericCode(CurrencyCodeUSD, test.input)
				})
				return
			}

			result := ParseNumericCode(CurrencyCodeUSD, test.input)

			assert.Equal(t, test.expect.String(), result.String())
		})
	}
}

func BenchmarkParseCurrency(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ParseCurrency(CurrencyCodeUSD, "123,456.78")
	}
}
