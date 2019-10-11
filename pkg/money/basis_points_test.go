package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequenceValidation(t *testing.T) {

	assert := assert.New(t)

	bps1, err := ParseBasisPoints("2.75%")
	assert.NoError(err)

	bps2, err := ParseBasisPoints("2.855%")
	assert.NoError(err)

	bps3, err := ParseBasisPoints("2.95%")
	assert.NoError(err)

	err = ValidateBasisPointsSequential(bps1, bps2, bps3)
	assert.NoError(err)

	err = ValidateBasisPointsSequential(bps2, bps1, bps3)
	assert.Error(err)

	err = ValidateBasisPointsSequential(bps1, bps3, bps2)
	assert.Error(err)

}

func TestPercentDoubleDecimalParsing(t *testing.T) {

	bps, err := ParseBasisPoints("2.75%")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(275, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal(275, bps.WholeBasisPoints())

	assert.Equal(0, bps.FractionalBasisPoints())
	assert.Equal("275", bps.String())
	assert.Equal("2.75%", bps.PercentString())
	assert.Equal("0.0275", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.0275")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal(275, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal("0.0275", bps.DecimalString())

}

func TestOneBasisPointParsing(t *testing.T) {

	bps, err := ParseBasisPoints("1")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(1, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal(1, bps.WholeBasisPoints())
	assert.Equal(0, bps.FractionalBasisPoints())
	assert.Equal("1", bps.String())
	assert.Equal("0.01%", bps.PercentString())
	assert.Equal("0.0001", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.0001000")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal(1, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal("0.0001", bps.DecimalString())

}

func TestSimpleBasisPointsParsing(t *testing.T) {

	bps, err := ParseBasisPoints("20")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(20, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal(20, bps.WholeBasisPoints())
	assert.Equal(0, bps.FractionalBasisPoints())
	assert.Equal("20", bps.String())
	assert.Equal("0.2%", bps.PercentString())
	assert.Equal("0.002", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.002")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal(20, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal("0.002", bps.DecimalString())

}

func TestSimplePercentageParsing(t *testing.T) {

	bps, err := ParseBasisPoints("2%")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(200, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal(200, bps.WholeBasisPoints())
	assert.Equal(0, bps.FractionalBasisPoints())
	assert.Equal("200", bps.String())
	assert.Equal("2%", bps.PercentString())
	assert.Equal("0.02", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.02")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.02", bps.DecimalString())

}

func TestFractionalPercentageBasisPointsParsing(t *testing.T) {

	bps, err := ParseBasisPoints("2.545%")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(2545, bps.Numerator)
	assert.Equal(100000, bps.Denominator)
	assert.Equal(254, bps.WholeBasisPoints())
	assert.Equal(5, bps.FractionalBasisPoints())
	assert.Equal("254.5", bps.String())
	assert.Equal("2.545%", bps.PercentString())
	assert.Equal("0.02545", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.02545")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.02545", bps.DecimalString())

}

func TestFractionalPercentageHalfBasisPointParsing(t *testing.T) {

	bps, err := ParseBasisPoints("0.005%")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(5, bps.Numerator)
	assert.Equal(100000, bps.Denominator)
	assert.Equal(0, bps.WholeBasisPoints())
	assert.Equal(5, bps.FractionalBasisPoints())
	assert.Equal("0.5", bps.String())
	assert.Equal("0.005%", bps.PercentString())
	assert.Equal("0.00005", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.00005")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.00005", bps.DecimalString())

}

func TestZeroTenthsPercentageParsingWithFractionalBasisPoints(t *testing.T) {

	bps, err := ParseBasisPoints("0.205%")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(205, bps.Numerator)
	assert.Equal(100000, bps.Denominator)
	assert.Equal(20, bps.WholeBasisPoints())
	assert.Equal(5, bps.FractionalBasisPoints())
	assert.Equal("20.5", bps.String())
	assert.Equal("0.205%", bps.PercentString())
	assert.Equal("0.00205", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.00205")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.00205", bps.DecimalString())

}

func TestTenthsPercentageParsing(t *testing.T) {

	bps, err := ParseBasisPoints("2.05%")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(205, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal(205, bps.WholeBasisPoints())
	assert.Equal(0, bps.FractionalBasisPoints())
	assert.Equal("205", bps.String())
	assert.Equal("2.05%", bps.PercentString())
	assert.Equal("0.0205", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.0205")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.0205", bps.DecimalString())

}

func TestFractionalPercentageParsing(t *testing.T) {

	bps, err := ParseBasisPoints("2.5%")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(250, bps.Numerator)
	assert.Equal(10000, bps.Denominator)
	assert.Equal(250, bps.WholeBasisPoints())
	assert.Equal(0, bps.FractionalBasisPoints())
	assert.Equal("250", bps.String())
	assert.Equal("2.5%", bps.PercentString())
	assert.Equal("0.025", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.025")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.025", bps.DecimalString())

}

func TestFractionalBasisPointsParsing(t *testing.T) {

	bps, err := ParseBasisPoints("3.5")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(35, bps.Numerator)
	assert.Equal(100000, bps.Denominator)
	assert.Equal(3, bps.WholeBasisPoints())
	assert.Equal(5, bps.FractionalBasisPoints())
	assert.Equal("3.5", bps.String())
	assert.Equal("0.035%", bps.PercentString())
	assert.Equal("0.00035", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.00035")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.00035", bps.DecimalString())

}

func TestTinyFractionalBasisPointsParsing(t *testing.T) {

	bps, err := ParseBasisPoints("3.55")

	assert := assert.New(t)
	assert.NoError(err)
	assert.NotNil(bps)

	assert.Equal(355, bps.Numerator)
	assert.Equal(1000000, bps.Denominator)
	assert.Equal(3, bps.WholeBasisPoints())
	assert.Equal(55, bps.FractionalBasisPoints())
	assert.Equal("3.55", bps.String())
	assert.Equal("0.0355%", bps.PercentString())
	assert.Equal("0.000355", bps.DecimalString())

	bps, err = ParseBasisPointsDecimal("0.000355")
	assert.NoError(err)
	assert.NotNil(bps)
	assert.Equal("0.000355", bps.DecimalString())

}
