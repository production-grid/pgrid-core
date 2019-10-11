package money

import (
	"database/sql/driver"
	"errors"
	"math"
	"strconv"
	"strings"
)

// constants for use in great business success
const (
	BpsPerWhole       = 10000
	PctPointsPerWhole = 100
)

// BasisPoints models percentages at a hyper crazy level of detail.
type BasisPoints struct {
	Numerator   int
	Denominator int
}

//Float64 returns the ratio as a float 64.
func (bps *BasisPoints) Float64() float64 {
	return float64(bps.Numerator) / float64(bps.Denominator)
}

// ValidateBasisPointsSequential validates that the given basis points are sequential
func ValidateBasisPointsSequential(values ...*BasisPoints) error {

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

// NewBasisPoints generates a simple
func NewBasisPoints(bps int) *BasisPoints {

	return &BasisPoints{
		Numerator:   bps,
		Denominator: BpsPerWhole,
	}

}

// NewFractionalBasisPoints is used for jon compliant fractional basis points.
func NewFractionalBasisPoints(numerator int, denominator int) *BasisPoints {

	return &BasisPoints{
		Numerator:   numerator,
		Denominator: denominator,
	}

}

// ParseBasisPointsDecimal parses a bps ratio formatted as a decimal coefficient
func ParseBasisPointsDecimal(decString string) (*BasisPoints, error) {

	decString = strings.TrimRight(decString, "0")

	tokens := strings.Split(decString, ".")

	inputDecPlaces := len(tokens[1])
	minDecPlaces := int(math.Log10(BpsPerWhole))

	if inputDecPlaces <= minDecPlaces {
		//add padding right zeros
		paddedString := tokens[1] + strings.Repeat("0", minDecPlaces-inputDecPlaces)
		numer, err := strconv.Atoi(paddedString)
		if err != nil {
			return nil, err
		}
		result := BasisPoints{
			Numerator:   numer,
			Denominator: BpsPerWhole,
		}
		return &result, nil
	}

	denom := int(math.Pow10(inputDecPlaces))
	numer, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, err
	}
	result := BasisPoints{
		Numerator:   numer,
		Denominator: denom,
	}
	return &result, nil

}

// ParsePercentage assumes the input string is a percentage, even if the symbol is missing.
func ParsePercentage(pctString string) (*BasisPoints, error) {
	if len(strings.TrimSpace(pctString)) == 0 {
		return &BasisPoints{}, nil
	}

	if strings.HasPrefix(pctString, "%") {
		return ParseBasisPoints(pctString)
	}

	return ParseBasisPoints(pctString + "%")

}

// ParseBasisPoints can parse percentages or basis points
func ParseBasisPoints(bpsString string) (*BasisPoints, error) {

	if len(strings.TrimSpace(bpsString)) == 0 {
		return &BasisPoints{}, nil
	}

	//if the string is given as a percentage, translate to a bps string
	if strings.Contains(bpsString, "%") {
		bpsString = strings.Replace(bpsString, "%", "", 1)
		dotIndex := strings.Index(bpsString, ".")
		if dotIndex < 0 {
			bpsString = bpsString + "00"
		} else {
			bpsString = strings.Replace(bpsString, ".", "", 1)
			bpsDotIndex := dotIndex + 2
			if len(bpsString) <= bpsDotIndex {
				bpsString = bpsString + strings.Repeat("0", bpsDotIndex-len(bpsString))
			} else {
				bpsString = bpsString[:bpsDotIndex] + "." + bpsString[bpsDotIndex:]
			}
		}
	}

	//if we have a fractional component it gets tricky
	if strings.Contains(bpsString, ".") {
		tokens := strings.Split(bpsString, ".")
		wholeBps, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, err
		}
		fracVal, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, err
		}

		//set some starting assumptions
		decPlaces := len(tokens[1])
		denom := BpsPerWhole
		bps := wholeBps
		//if decPlaces is greater than zero, we have fractional basis points
		//and have to increse the denominator
		if decPlaces > 0 {
			denom = denom * int(math.Pow10(decPlaces))
			bps = bps * int(math.Pow10(decPlaces))
		}
		//if our denominator is greater than 10,000, we can simply the process
		//since there's no need to pad the result with right hand zeros
		if denom > BpsPerWhole {
			//this is a short circuit option because we know we don't have to right
			//pad the numerator
			baseString := strings.Replace(bpsString, ".", "", -1)
			numer, err := strconv.Atoi(baseString)
			if err != nil {
				return nil, err
			}
			result := BasisPoints{
				Numerator:   numer,
				Denominator: denom,
			}
			return &result, nil

		}

		//add the fractional component to the whole component
		bps += fracVal
		result := BasisPoints{
			Numerator:   bps,
			Denominator: denom,
		}
		return &result, nil

	}
	//if we made it here there are no fractional basis points or percentage points
	bps, err := strconv.Atoi(bpsString)
	if err != nil {
		return nil, err
	}
	result := BasisPoints{
		Numerator:   bps,
		Denominator: BpsPerWhole,
	}
	return &result, nil

}

//IsZero returns true if this is zero.
func (bps *BasisPoints) IsZero() bool {
	return bps.Numerator == 0
}

/*
Value implements the valuer interface in order to support sql serialization.
*/
func (bps BasisPoints) Value() (driver.Value, error) {

	return bps.DecimalString(), nil

}

/*
Scan implements the scan interface in order to support sql serialization.
*/
func (bps *BasisPoints) Scan(src interface{}) error {

	switch dbVal := src.(type) {
	case []uint8:
		floatString := string(dbVal)
		workingCopy, err := ParseBasisPointsDecimal(floatString)
		if err != nil {
			return err
		}
		bps.Numerator = workingCopy.Numerator
		bps.Denominator = workingCopy.Denominator

	}

	return nil

}

//DecimalString returns a normal decimal
func (bps *BasisPoints) DecimalString() string {

	if bps.Numerator == 0 {
		return "0"
	}

	sb := strings.Builder{}
	sb.WriteString("0.")

	numer := strconv.Itoa(bps.Numerator)

	expectedLen := int(math.Log10(float64(bps.Denominator)))

	leftPadding := expectedLen - len(numer)

	if leftPadding > 0 {
		sb.WriteString(strings.Repeat("0", leftPadding))
	}

	//strip right padded zeros
	numer = strings.TrimRight(numer, "0")

	sb.WriteString(numer)

	return sb.String()

}

func (bps *BasisPoints) String() string {
	if bps.Denominator > BpsPerWhole {
		wholePoints := bps.WholeBasisPoints()
		fracPoints := bps.FractionalBasisPoints()
		return strconv.Itoa(wholePoints) + "." + strconv.Itoa(fracPoints)
	}
	return strconv.Itoa(bps.Numerator)
}

//LTE returns true if the argument is less than or equal to the receiver
func (bps *BasisPoints) LTE(arg *BasisPoints) bool {

	rcvNum := bps.Numerator
	argNum := arg.Numerator

	if arg.Denominator > bps.Denominator {
		rcvNum *= (arg.Denominator / bps.Denominator)
	} else if arg.Denominator < bps.Denominator {
		argNum *= (bps.Denominator / arg.Denominator)
	}

	return argNum >= rcvNum

}

//LT returns true if the argument is less than the receiver
func (bps *BasisPoints) LT(arg *BasisPoints) bool {

	rcvNum := bps.Numerator
	argNum := arg.Numerator

	if arg.Denominator > bps.Denominator {
		rcvNum *= (arg.Denominator / bps.Denominator)
	} else if arg.Denominator < bps.Denominator {
		argNum *= (bps.Denominator / arg.Denominator)
	}

	return argNum > rcvNum

}

//GT returns true if the reciever is greater than the argument
func (bps *BasisPoints) GT(arg *BasisPoints) bool {

	rcvNum := bps.Numerator
	argNum := arg.Numerator

	if arg.Denominator > bps.Denominator {
		rcvNum *= (arg.Denominator / bps.Denominator)
	} else if arg.Denominator < bps.Denominator {
		argNum *= (bps.Denominator / arg.Denominator)
	}

	return argNum < rcvNum

}

//GTE returns true if the argument is greater than or equal to the receiver
func (bps *BasisPoints) GTE(arg *BasisPoints) bool {

	rcvNum := bps.Numerator
	argNum := arg.Numerator

	if arg.Denominator > bps.Denominator {
		rcvNum *= (arg.Denominator / bps.Denominator)
	} else if arg.Denominator < bps.Denominator {
		argNum *= (bps.Denominator / arg.Denominator)
	}

	return argNum <= rcvNum

}

// WholeBasisPoints returns whole basis points
func (bps *BasisPoints) WholeBasisPoints() int {

	denom := float64(bps.Denominator)
	numer := float64(bps.Numerator)
	return int(math.Floor(numer * BpsPerWhole / denom))

}

// FractionalBasisPoints returns fractional basis points
func (bps *BasisPoints) FractionalBasisPoints() int {

	denom := float64(bps.Denominator)
	numer := float64(bps.Numerator)

	extPlaces := denom / BpsPerWhole
	remainder := int(math.Mod(numer, extPlaces))

	return remainder

}

//PercentString returns the value as a percentage.
func (bps *BasisPoints) PercentString() string {

	wholePoints := float64(bps.WholeBasisPoints())

	percentagePoints := int(math.Floor(wholePoints / 100))

	fractionalPoints := int(math.Mod(wholePoints, 100))

	sb := strings.Builder{}

	sb.WriteString(strconv.Itoa(percentagePoints))

	if fractionalPoints > 0 {
		sb.WriteString(".")
		fracString := strconv.Itoa(fractionalPoints)
		if len(fracString) == 1 {
			sb.WriteString("0")
		}
		if bps.Denominator == BpsPerWhole {
			fracString = strings.TrimRight(fracString, "0")
		}
		sb.WriteString(fracString)
	} else if bps.FractionalBasisPoints() > 0 {
		//handle edge case where there are zero fractional percentage points, but fractional basis points
		sb.WriteString(".00")
	}
	fracBasisPoints := bps.FractionalBasisPoints()
	if fracBasisPoints > 0 {
		sb.WriteString(strconv.Itoa(fracBasisPoints))
	}
	sb.WriteString("%")

	return sb.String()
}
