package util

import (
	"strconv"
	"strings"

	"github.com/production-grid/pgrid-core/pkg/logging"
	"golang.org/x/text/message"
)

/*
FormatInteger returns the string representation of an integer
in localized format.
*/
func FormatInteger(value int) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprint(value)
}

// IntPointer returns a pointer to an int literal.
func IntPointer(value int) *int {

	return &value

}

/*
ParseInteger parses an integer.  Everybody Got that?
*/
func ParseInteger(rawString string) *int {
	rawString = strings.Replace(rawString, ",", "", -1)
	value, err := strconv.Atoi(rawString)
	if err != nil {
		logging.Warn(err)
		return nil
	}
	return &value
}

/*
FormatDecimal returns the string representation of a decimal
in localized format.
*/
func FormatDecimal(value float64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%.1f", value)

}

/*
FormatCurrencyWithSymbol returns the string representation of a decimal value
in localized currency format with a dollar sign.
*/
func FormatCurrencyWithSymbol(value float64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("$%.2f", value)

}

/*
FormatCurrency returns the string representation of a decimal value
in localized currency format.
*/
func FormatCurrency(value float64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%.2f", value)

}
