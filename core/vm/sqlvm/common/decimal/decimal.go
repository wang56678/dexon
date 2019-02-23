package decimal

import "github.com/dexon-foundation/decimal"

// Shared vars.
var (
	False = decimal.New(0, 0)
	True  = decimal.New(1, 0)

	Int64Max = decimal.New(1, 63).Sub(decimal.One)
	Int64Min = decimal.New(1, 63).Neg()

	UInt16Max = decimal.New(1, 16).Sub(decimal.One)
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
