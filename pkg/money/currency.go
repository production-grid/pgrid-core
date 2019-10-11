package money

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//RoundingMode indicates what technique to use for rounding values that fall out of the precision range.
type RoundingMode byte

//RoundingMode constants.
const (
	RoundUp      RoundingMode = iota //always round up
	RoundDown                        //always round down
	RoundNearest                     //round >= .5 up and < .5 down
)

/*
Currency models the integer and decimal components of currency along with currency code.
*/
type Currency struct {
	Integer            int
	Decimal            int
	DecimalDenominator int
	Negative           bool
	CurrencyCode       string
}

/*
CurrencyDefinition models locale specific characteristics of a currency.
*/
type CurrencyDefinition struct {
	DecimalPlaces       int
	DecimalSeparator    string
	CurrencySymbol      string
	ThousandsSeparator  string
	CurrencySymbolAfter bool
}

// Currency code constants.
const (
	CurrencyCodeBCP = "BCP"
	CurrencyCodeBTC = "BTC"
	CurrencyCodeCAD = "CAD"
	CurrencyCodeEUR = "EUR"
	CurrencyCodeGBP = "GBP"
	CurrencyCodeUSD = "USD"

	CurrencyCodeDefault = CurrencyCodeUSD
)

// CurrencyDefinitionMap includes configuration for various currencies.
var CurrencyDefinitionMap = map[string]CurrencyDefinition{
	CurrencyCodeBCP: CurrencyDefinition{DecimalPlaces: 0, DecimalSeparator: "."},
	CurrencyCodeBTC: CurrencyDefinition{DecimalPlaces: 8, DecimalSeparator: "."},
	CurrencyCodeCAD: CurrencyDefinition{DecimalPlaces: 2, DecimalSeparator: ".", ThousandsSeparator: ",", CurrencySymbol: "$"},
	CurrencyCodeUSD: CurrencyDefinition{DecimalPlaces: 2, DecimalSeparator: ".", ThousandsSeparator: ",", CurrencySymbol: "$"},
	CurrencyCodeGBP: CurrencyDefinition{DecimalPlaces: 2, DecimalSeparator: ".", ThousandsSeparator: ",", CurrencySymbol: "£"},
	CurrencyCodeEUR: CurrencyDefinition{DecimalPlaces: 2, DecimalSeparator: ",", ThousandsSeparator: ".", CurrencySymbol: "€", CurrencySymbolAfter: true},
}

// CurrencyDefinition returns the currency type definition. If it's empty, it
// will return the default.
func (base Currency) CurrencyDefinition() CurrencyDefinition {
	var code string

	if base.CurrencyCode == "" {
		code = CurrencyCodeDefault

	} else {
		code = base.CurrencyCode
	}

	def, ok := CurrencyDefinitionMap[code]
	if !ok {
		panic("unsupported currency: " + code)
	}

	return def
}

/*
Value implements the valuer interface in order to support sql serialization.
*/
func (base Currency) Value() (driver.Value, error) {

	if base.CurrencyCode == "" {
		return "0", nil
	}

	results := strconv.Itoa(base.Integer)

	def := base.CurrencyDefinition()

	results += string(def.DecimalSeparator)

	rawDecimal := strconv.Itoa(base.Decimal)

	decimalPlaces := int(math.Log10(float64(base.DecimalDenominator)))

	if len(rawDecimal) < decimalPlaces {
		rawDecimal = strings.Repeat("0", decimalPlaces-len(rawDecimal)) + rawDecimal
	}

	results += rawDecimal

	if base.Negative && !base.IsZero() {
		results = "-" + results
	}

	return results, nil

}

/*
Scan implements the scan interface in order to support sql serialization.
*/
func (base *Currency) Scan(src interface{}) error {

	switch dbVal := src.(type) {
	case []uint8:
		floatString := string(dbVal)
		tmpCur := ParseCurrency(CurrencyCodeUSD, floatString)
		base.Integer = tmpCur.Integer
		base.CurrencyCode = tmpCur.CurrencyCode
		base.Decimal = tmpCur.Decimal
		base.DecimalDenominator = tmpCur.DecimalDenominator
		base.Negative = tmpCur.Negative
	}

	return nil

}

/*
SetCurrencyCode sets the currency and computes the correct
decimal for the currency's number of decimal places.
*/
func (base *Currency) SetCurrencyCode(currencyCode string) *Currency {

	base.CurrencyCode = currencyCode
	return base.ComputeDecimal()

}

/*
ComputeDecimal computes the decimal component for a given currency.
*/
func (base *Currency) ComputeDecimal() *Currency {

	def := base.CurrencyDefinition()

	targetDenominator := int(math.Pow10(def.DecimalPlaces))

	if base.DecimalDenominator == 0 {
		return base
	}

	if targetDenominator < base.DecimalDenominator {
		return base
	}

	base.Decimal = base.Decimal * targetDenominator / base.DecimalDenominator

	return base

}

/*
FromFloat64 takes a float value and converts it to Currency.
*/
func FromFloat64(currencyCode string, value float64) Currency {
	return *ParseCurrency(currencyCode, fmt.Sprintf("%f", value))
}

func (base *Currency) decimalMultiplier() int {

	def := base.CurrencyDefinition()

	return int(math.Pow10(def.DecimalPlaces))

}

func (base *Currency) assertCurrencyMatch(operand *Currency) {

	if base.IsZero() || operand.IsZero() {
		return
	}

	if base.CurrencyCode != operand.CurrencyCode {
		panic(base.CurrencyCode + " and " + operand.CurrencyCode + " are different currencies")
	}

}

// normalizeCurrency converts zero amounts to match the currency code of
// non-zero amounts for operations on multiple currencies.
func (base *Currency) normalizeCurrency(operand *Currency) {
	base.assertCurrencyMatch(operand)

	if base.IsZero() {
		if !operand.IsZero() {
			base.CurrencyCode = operand.CurrencyCode
			if base.DecimalDenominator == 0 {
				def := operand.CurrencyDefinition()
				base.DecimalDenominator = int(math.Pow10(def.DecimalPlaces))
			}
		}
	} else {
		if operand.IsZero() {
			operand.CurrencyCode = base.CurrencyCode
			def := operand.CurrencyDefinition()
			operand.DecimalDenominator = int(math.Pow10(def.DecimalPlaces))
		}
	}

	if !base.IsZero() && base.DecimalDenominator == 0 {
		def := operand.CurrencyDefinition()
		base.DecimalDenominator = int(math.Pow10(def.DecimalPlaces))
	}

	if !operand.IsZero() && operand.DecimalDenominator == 0 {
		def := operand.CurrencyDefinition()
		operand.DecimalDenominator = int(math.Pow10(def.DecimalPlaces))
	}

}

//toScaledInts converts the given operands to integers in the same space
func toScaledInts(op1, op2 *Currency) (r1 int, r2 int, denom int) {

	if op1.DecimalDenominator == 0 {
		op1.normalize()
	}
	if op2.DecimalDenominator == 0 {
		op2.normalize()
	}

	denom = op1.DecimalDenominator
	if denom < op2.DecimalDenominator {
		denom = op2.DecimalDenominator
	}

	r1 = op1.toScaledInt()
	r2 = op2.toScaledInt()

	if denom > op1.DecimalDenominator {
		opPlaces := math.Log10(float64(op1.DecimalDenominator))
		denomPlaces := math.Log10(float64(denom))
		diff := int(denomPlaces - opPlaces)
		r1 *= int(math.Pow10(diff))

	}

	if denom > op2.DecimalDenominator {
		opPlaces := math.Log10(float64(op2.DecimalDenominator))
		denomPlaces := math.Log10(float64(denom))
		diff := int(denomPlaces - opPlaces)
		r2 *= int(math.Pow10(diff))
	}

	return r1, r2, denom
}

/*
Add adds the operand to base and returns the result.
*/
func (base *Currency) Add(operand *Currency) *Currency {

	baseInt, opInt, denom := toScaledInts(base, operand)

	sumScaled := baseInt + opInt

	currency := base.CurrencyCode
	if currency == "" {
		currency = operand.CurrencyCode
	}

	return fromScaledInt(currency, sumScaled, denom)

}

/*
Subtract subtracts the operand from base and returns the result.
*/
func (base *Currency) Subtract(operand *Currency) *Currency {

	baseInt, opInt, denom := toScaledInts(base, operand)

	diffScaled := baseInt - opInt

	currency := base.CurrencyCode
	if currency == "" {
		currency = operand.CurrencyCode
	}

	return fromScaledInt(currency, diffScaled, denom)

}

/*
IsZero returns true if currency value is zero.
*/
func (base *Currency) IsZero() bool {

	if base.Integer != 0 {
		return false
	}

	if base.Decimal != 0 {
		return false
	}

	return true

}

/*
Equals returns true if the operand is numerically equal to base.
*/
func (base *Currency) Equals(operand *Currency) bool {

	if base.IsZero() && operand.IsZero() {
		return true
	}

	if base.CurrencyCode != operand.CurrencyCode {
		return false
	}

	if base.Integer != operand.Integer {
		return false
	}

	if base.Decimal != operand.Decimal {
		return false
	}

	return true

}

/*
NotEquals returns true if the operand is not equal to base.
*/
func (base *Currency) NotEquals(operand *Currency) bool {
	return !base.Equals(operand)
}

/*
MultBPS multiplies currency by the given percentage or basis points.
*/
func (base *Currency) MultBPS(bps *BasisPoints) *Currency {

	return base.MultFloat(bps.Float64(), RoundUp)

}

/*
Pct multiplies itself by the given percentage and returns the result.
*/
func (base *Currency) Pct(operand float64) *Currency {

	decimalMultiplier := base.decimalMultiplier()

	result := Currency{CurrencyCode: base.CurrencyCode}

	baseFloat := float64(base.Integer) + float64(base.Decimal)/float64(decimalMultiplier)

	resultFloat := baseFloat * (operand / float64(decimalMultiplier))

	resultInt := int(resultFloat * float64(decimalMultiplier))

	result.Integer = resultInt / decimalMultiplier
	result.Decimal = resultInt % decimalMultiplier

	return &result
}

func intAbs(val int) int {
	if val < 0 {
		return 0 - val
	}
	return val
}

//Round rounds a given currency to the given level of precision
func (base *Currency) Round(decimalPlaces int, roundMode RoundingMode) *Currency {

	floatValue := base.ToFloat()

	denom := math.Pow10(decimalPlaces)

	switch roundMode {
	case RoundDown:
		floatValue = math.Floor(floatValue*denom) / denom
	case RoundUp:
		floatValue = math.Ceil(floatValue*denom) / denom
	case RoundNearest:
		floatValue = math.Round(floatValue*denom) / denom
	}

	return New(base.CurrencyCode, fmt.Sprintf("%f", floatValue))

}

func (base *Currency) toScaledInt() int {
	absResult := base.Decimal + (intAbs(base.Integer) * base.DecimalDenominator)
	if base.Negative {
		absResult = -absResult
	}
	return absResult
}

func fromScaledInt(currencyCode string, scaledInt int, decimalDenominator int) *Currency {

	result := Currency{
		CurrencyCode: currencyCode,
	}

	if decimalDenominator > 0 {
		scaledAbs := intAbs(scaledInt)
		result.Decimal = scaledInt % decimalDenominator
		if result.Decimal < 0 {
			result.Decimal = -result.Decimal
		}
		result.Integer = (scaledAbs - result.Decimal) / decimalDenominator
		result.DecimalDenominator = decimalDenominator
	} else {
		result.Integer = scaledInt
	}

	if scaledInt < 0 {
		result.Negative = true
	}

	return &result

}

/*
Mult multiplies itself by the operand and returns the result.
*/
func (base *Currency) Mult(operand *Currency) *Currency {

	base.normalizeCurrency(operand)

	scaledBase := base.toScaledInt()
	scaledOp := operand.toScaledInt()

	scaledResult := scaledBase * scaledOp / base.DecimalDenominator

	return fromScaledInt(base.CurrencyCode, scaledResult, base.DecimalDenominator)

}

func (base *Currency) normalize() {

	if base.DecimalDenominator == 0 {
		def := base.CurrencyDefinition()
		base.DecimalDenominator = int(math.Pow10(def.DecimalPlaces))
	}

}

/*
MultInt multiplies the receiver by an integer.
*/
func (base *Currency) MultInt(operand int) *Currency {

	if base.IsZero() {
		return base
	}

	base.normalize()

	scaledBase := base.toScaledInt()
	scaledResult := scaledBase * operand

	return fromScaledInt(base.CurrencyCode, scaledResult, base.DecimalDenominator)

}

// MultFloat multiplies the receiver by a floating point number.
func (base *Currency) MultFloat(operand float64, mode RoundingMode) *Currency {
	if base.IsZero() {
		return base
	}

	base.normalize()
	scaledBase := base.toScaledInt()

	scaledResult := float64(scaledBase) * operand

	var roundedResult int
	switch mode {
	case RoundUp:
		roundedResult = int(math.Ceil(scaledResult))
	case RoundNearest:
		roundedResult = int(math.Round(scaledResult))
	}

	return fromScaledInt(base.CurrencyCode, roundedResult, base.DecimalDenominator)
}

// ValidateCurrencySequential validates that the given basis points are sequential
func ValidateCurrencySequential(values ...*Currency) error {

	for i := 0; i < len(values)-1; i++ {
		lt := values[i]
		gt := values[i+1]
		if lt.IsZero() {
			continue
		}
		if gt.IsZero() {
			continue
		}
		if lt.GT(gt) {
			return errors.New("values out of sequence")
		}
	}

	return nil

}

/*
LT executes a less than comparison.
*/
func (base *Currency) LT(operand *Currency) bool {

	base.normalizeCurrency(operand)

	if base.Integer != operand.Integer {
		return base.Integer < operand.Integer
	}

	return base.Decimal < operand.Decimal

}

/*
LTE executes a less than comparison.
*/
func (base *Currency) LTE(operand *Currency) bool {

	base.normalizeCurrency(operand)

	if base.Equals(operand) {
		return true
	}

	return base.LT(operand)

}

/*
GTE executes a less than comparison.
*/
func (base *Currency) GTE(operand *Currency) bool {

	base.normalizeCurrency(operand)

	if base.Equals(operand) {
		return true
	}

	return base.GT(operand)

}

/*
GT executes a greater than comparison.
*/
func (base *Currency) GT(operand *Currency) bool {

	base.normalizeCurrency(operand)

	if base.Integer != operand.Integer {
		return base.Integer > operand.Integer
	}

	return base.Decimal > operand.Decimal

}

/*
DivideInt divides itself by the operand and returns the result.
*/
func (base *Currency) DivideInt(operand int, mode RoundingMode) *Currency {

	base.normalize()

	baseScaled := base.toScaledInt()

	absOperand := operand

	if operand < 0 {
		absOperand = -operand
	}

	remainder := baseScaled % absOperand

	resultScaled := (baseScaled - remainder) / absOperand

	switch mode {
	case RoundUp:
		if remainder > 0 {
			resultScaled++
		}
	case RoundNearest:
		if (absOperand - remainder) < remainder {
			resultScaled++
		}
	}

	result := fromScaledInt(base.CurrencyCode, resultScaled, base.DecimalDenominator)

	if base.Negative && (operand < 0) {
		result.Negative = false
	} else if base.Negative || (operand < 0) {
		result.Negative = true
	}

	return result

}

/*
Divide divides itself by the operand and returns the result.
*/
func (base *Currency) Divide(operand *Currency) *Currency {

	base.normalizeCurrency(operand)
	decimalMultiplier := base.decimalMultiplier()

	result := Currency{CurrencyCode: base.CurrencyCode}

	baseFloat := float64(base.Integer) + float64(base.Decimal)/float64(decimalMultiplier)
	opFloat := float64(operand.Integer) + float64(operand.Decimal)/float64(decimalMultiplier)

	resultFloat := baseFloat / opFloat

	resultInt := int(resultFloat * float64(decimalMultiplier))

	result.Integer = resultInt / decimalMultiplier
	result.Decimal = resultInt % decimalMultiplier

	return &result

}

/*
ToFloat returns the float version of the currency.
*/
func (base *Currency) ToFloat() float64 {

	result := float64(base.Integer)

	if base.DecimalDenominator == 0 {
		result += float64(base.Decimal) / float64(base.decimalMultiplier())
	} else {
		result += float64(base.Decimal) / float64(base.DecimalDenominator)
	}

	if base.Negative {
		result = -result
	}

	return result

}

/*
Negate flips the sign of the integer value.
*/
func (base *Currency) Negate() *Currency {

	result := *base
	result.Negative = !result.Negative
	return &result

}

/*
New creates a new currency from a string.
*/
func New(currencyCode string, value string) *Currency {

	return ParseCurrency(currencyCode, value)

}

/*
Zero returns a zero value constant.
*/
func Zero() *Currency {
	return &Currency{}
}

/*
ZeroDollars returns a zero value constant.
*/
func ZeroDollars() *Currency {
	return &Currency{CurrencyCode: CurrencyCodeUSD, DecimalDenominator: 100}
}

/*
One returns a one value constant.
*/
func One(currency string) *Currency {
	return &Currency{CurrencyCode: currency, Integer: 1}
}

/*
OneDollar dollar value constant.
*/
func OneDollar() *Currency {
	return &Currency{CurrencyCode: CurrencyCodeUSD, Integer: 1}
}

/*
String returns the default string representation of the currency.
*/
func (base *Currency) String() string {
	return base.FormatCurrency()
}

// Localize returns a localized human-readable string form of the currency.
func (base *Currency) Localize() string {
	def := base.CurrencyDefinition()

	var b strings.Builder

	if def.CurrencySymbolAfter {
		b.WriteString(base.FormatCurrencyRounded())
		b.WriteString(def.CurrencySymbol)
	} else {
		b.WriteString(def.CurrencySymbol)
		b.WriteString(base.FormatCurrencyRounded())
	}

	return b.String()
}

/*
FormatCurrency turns the given currency into a string.
*/
func (base *Currency) FormatCurrency() string {

	if base.CurrencyCode == "" {
		return "0"
	}

	lang := "en-US"
	if base.CurrencyCode == CurrencyCodeEUR {
		lang = "it"
	}

	var b strings.Builder

	p := message.NewPrinter(language.Make(lang))

	if base.Negative && !base.IsZero() {
		b.WriteRune('-')
	}

	b.WriteString(p.Sprintf("%d", base.Integer))

	def := base.CurrencyDefinition()

	b.WriteString(def.DecimalSeparator)

	rawDecimal := strconv.Itoa(base.Decimal)

	if base.DecimalDenominator != 0 {
		decPlaces := int(math.Log10(float64(base.DecimalDenominator)))

		if len(rawDecimal) < decPlaces {
			b.WriteString(strings.Repeat("0", decPlaces-len(rawDecimal)))
		}
	} else if def.DecimalPlaces > len(rawDecimal) {
		b.WriteString(strings.Repeat("0", def.DecimalPlaces-len(rawDecimal)))
	}

	b.WriteString(rawDecimal)

	return b.String()

}

// FormatCurrencyRounded converts the currency to a string and rounds it to the
// currency's precision.
func (base *Currency) FormatCurrencyRounded() string {
	if base.CurrencyCode == "" {
		return "0"
	}

	lang := "en-US"
	if base.CurrencyCode == CurrencyCodeEUR {
		lang = "it"
	}
	p := message.NewPrinter(language.Make(lang))

	float := base.ToFloat()

	return p.Sprintf("%.*f", base.CurrencyDefinition().DecimalPlaces, float)
}

func isZeroString(value string) bool {
	if value == "" {
		return true
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	if f == 0 {
		return true
	}
	return false
}

/*
ParseCurrency parses a string encoded currency for the given currency code.
*/
func ParseCurrency(currencyCode string, value string) *Currency {

	if isZeroString(value) {
		return Zero()
	}

	def, ok := CurrencyDefinitionMap[currencyCode]
	if !ok {
		panic(currencyCode + " is not a supported currency")
	}

	if strings.Contains(value, def.DecimalSeparator) {
		value = strings.TrimRight(value, "0")
	}

	currency := Currency{CurrencyCode: currencyCode}

	if strings.HasPrefix(value, "-") {
		currency.Negative = true
		value = value[1:]
	}

	tokens := strings.Split(value, def.DecimalSeparator)

	for i := range tokens {
		tokens[i] = stripSymbols(tokens[i])
	}

	if len(tokens) > 0 {
		currency.Integer, _ = strconv.Atoi(tokens[0])
		if len(tokens) > 1 {
			padded := padDecimal(tokens[1], def)
			currency.Decimal, _ = strconv.Atoi(padded)
			if len(tokens[1]) > def.DecimalPlaces {
				currency.DecimalDenominator = int(math.Pow10(len(tokens[1])))
			} else {
				currency.DecimalDenominator = int(math.Pow10(def.DecimalPlaces))
			}
		}
	} else {
		currency.DecimalDenominator = int(math.Pow10(def.DecimalPlaces))
	}

	if currency.IsZero() {
		currency.Negative = false
	}

	return &currency

}

func stripSymbols(numeric string) string {
	filter := func(r rune) rune {
		if r < '0' || r > '9' {
			// Non-numeric ascii character
			return -1
		}

		return r
	}

	return strings.Map(filter, numeric)
}

// ParseNumericCode parses a numeric encoded amount to currency.
func ParseNumericCode(currencyCode string, valueCode string) *Currency {
	def := CurrencyDefinitionMap[currencyCode]
	var err error

	var integer, decimal int
	if len(valueCode) > def.DecimalPlaces {
		decimal, err = strconv.Atoi(valueCode[len(valueCode)-def.DecimalPlaces : len(valueCode)])
		if err != nil {
			panic(err)
		}

		integer, err = strconv.Atoi(valueCode[:len(valueCode)-def.DecimalPlaces])
		if err != nil {
			panic(err)
		}
	} else {
		decimal, err = strconv.Atoi(valueCode)
		if err != nil {
			panic(err)
		}
	}

	return &Currency{
		CurrencyCode: CurrencyCodeUSD,
		Integer:      integer,
		Decimal:      decimal,
	}
}

// padDecimal zero pads a decimal to the correct length for the currency.
// use for parsing only
func padDecimal(raw string, def CurrencyDefinition) string {
	var b strings.Builder

	b.WriteString(raw)

	if len(raw) < def.DecimalPlaces {
		b.WriteString(strings.Repeat("0", def.DecimalPlaces-len(raw)))
	}

	return b.String()
}

//UnmarshalJSON implements the json.Unmarshaler interface.
func (base *Currency) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		var raw string
		if err := json.Unmarshal(b, &raw); err != nil {
			return err
		}

		*base = *ParseCurrency(CurrencyCodeUSD, raw)
		return nil
	}

	var raw float64
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	*base = FromFloat64(CurrencyCodeUSD, raw)

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (base Currency) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer

	b.WriteString("\"")
	b.WriteString(base.String())
	b.WriteString("\"")

	return b.Bytes(), nil
}
