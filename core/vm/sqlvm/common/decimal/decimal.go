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
