package runtime

var jumpTable = [256]OpFunction{
	// 0x10
	ADD: opAdd,
	MUL: opMul,
	SUB: opSub,
	DIV: opDiv,
	MOD: opMod,

	// 0x20
	LT:    opLt,
	GT:    opGt,
	EQ:    opEq,
	AND:   opAnd,
	OR:    opOr,
	NOT:   opNot,
	UNION: opUnion,
	INTXN: opIntxn,
	LIKE:  opLike,

	// 0x40
	ZIP:    opZip,
	FIELD:  opField,
	PRUNE:  opPrune,
	SORT:   opSort,
	FILTER: opFilter,
	CAST:   opCast,

	// 0x60
	LOAD: opLoad,
}
