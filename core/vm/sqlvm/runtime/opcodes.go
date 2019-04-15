package runtime

// OpCode type
type OpCode byte

// 0x00 range - higher order ops.
const (
	NOP OpCode = iota
)

// 0x10 range - arithmetic ops. (array-based)
const (
	ADD OpCode = iota + 0x10
	MUL
	SUB
	DIV
	MOD
	CONCAT
	NEG
)

// 0x20 range - comparison ops.
const (
	LT OpCode = iota + 0x20
	GT
	EQ
	AND
	OR
	NOT
	UNION
	INTXN
	LIKE
)

// 0x30 range - pk/index/field meta ops
const (
	REPEATPK OpCode = iota + 0x30
	REPEATIDX
	REPEATIDXV
)

// 0x40 range - format/output ops
const (
	ZIP OpCode = iota + 0x40
	FIELD
	PRUNE
	SORT
	FILTER
	CAST
	CUT
	RANGE
)

// 0x50 range - function ops
const (
	FUNC = iota + 0x50
)

// 0x60 range - storage ops
const (
	INSERT = iota + 0x60
	UPDATE
	LOAD
	DELETE
)
