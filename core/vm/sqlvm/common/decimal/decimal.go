package decimal

import (
	"fmt"
	"math"

	"github.com/dexon-foundation/decimal"
)

// Shared vars.
var (
	False = decimal.New(0, 0)
	True  = decimal.New(1, 0)

	MaxInt64 = decimal.New(math.MaxInt64, 0)
	MinInt64 = decimal.New(math.MinInt64, 0)

	MaxUint16 = decimal.New(math.MaxUint16, 0)
	MaxUint64 = decimal.RequireFromString(fmt.Sprint(uint64(math.MaxUint64)))

	Dec257 = decimal.New(257, 0)
)

// Val2Bool convert value to boolean definition.
func Val2Bool(v decimal.Decimal) decimal.Decimal {
	if v.IsZero() {
		return False
	}
	return True
}

// IsTrue returns given value is decimal defined True in golang boolean.
func IsTrue(v decimal.Decimal) bool {
	return v.Cmp(True) == 0
}

// IsFalse returns given value is decimal defined False in golang boolean.
func IsFalse(v decimal.Decimal) bool {
	return v.Cmp(False) == 0
}
