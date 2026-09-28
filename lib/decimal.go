package lib

import (
	"errors"
	"math"
	"strconv"

	flib "github.com/Files-com/files-sdk-go/v3/lib"
)

var errNotDecimal = errors.New("expected a decimal number such as 1.5 or 2e-3")

// ParseDecimalAttribute reads a string attribute that earlier provider
// versions parsed with strconv.ParseFloat before sending it as a number.
// Decimal text such as "1.5" or "2e-3" comes back unchanged as exact, so the
// API receives it without float rounding. Other spellings ParseFloat accepts,
// such as hexadecimal floats or digits separated by underscores, keep their
// former float64 value and rounding, with a nil exact. NaN, infinities and
// values ParseFloat rejects are errors.
func ParseDecimalAttribute(text string) (exact *string, legacy float64, err error) {
	if flib.IsDecimalText(text) {
		return &text, 0, nil
	}
	legacy, err = strconv.ParseFloat(text, 64)
	if err != nil || math.IsNaN(legacy) || math.IsInf(legacy, 0) {
		return nil, 0, errNotDecimal
	}
	return nil, legacy, nil
}
