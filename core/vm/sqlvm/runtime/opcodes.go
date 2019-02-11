package runtime

// OpCode type
type OpCode byte

// 0x10 range - arithmetic ops. (array-based)
const (
	ADD OpCode = iota + 0x10
	// ADD(t1, t2) res
	// res = t1 + t2 = [1, 2] + [2, 3] = [3, 5]
	MUL
	// MUL(t1, t2) res
	// res = t1 * t2 = [1, 2] * [2, 3] = [2, 6]
	SUB
	// SUB(t1, t2) res
	// res = t1 - t2 = [1, 2] - [2, 3] = [-1, -1]
	DIV
	// DIV(t1, t2) res
	// res = t1 / t2 = [1, 2] / [2, 3] = [0, 0]
	MOD
	// MOD(t1, t2) res
	// res = t1 % t2 = [1, 2] % [2, 3] = [1, 2]

)

// 0x20 range - comparison ops.
const (
	LT OpCode = iota + 0x20
	// LT(t1, t2) res
	// res = t1 < t2 = [1, 2] < [2, 3] = [true, true]
	GT
	// GT(t1, t2) res
	// res = t1 > t2 = [1, 2] > [2, 3] = [false, false]
	EQ
	// EQ(t1, t2) res
	// res = t1 == t2 = [1, 2] == [2, 3] = [false, false]
	AND
	// AND(t1, t2) res
	// res = t1 && t2  = [true, true] && [true, false] = [true, false]
	OR
	// OR(t1, t2) res
	// res = t1 || t2 = [false, false] || [true, false] = [true, false]
	NOT
	// NOT(t1) res
	// res = !t1 = ![false, true] = [true, false]
	UNION
	// UNION(t1, t2) res
	// res = t1 ∪ t2  = [1, 2] ∪ [2, 3] = [1, 2, 3]
	INTXN
	// INTXN(t1, t2) res
	// res = t1 ∩ t2 = [1, 2] ∩ [2, 3] = [2]
	LIKE
	// LIKE(t1, pattern) res
	// res = t1 like '%abc%' =
	//       ['_abc_', '123'] like '%abc%' = [true, false]
)

// 0x30 range - pk/index/field meta ops
const (
	REPEATPK OpCode = iota + 0x30
	// REPEATPK(pk) res    res = [id1, id2, id3, ...]
	// REPEATPK([tables, table_id, primary]) = [1, 2, 3, 5, 6, 7, ...]
	REPEATIDX
	// Scan given index value(s)
	// REPEATIDX(base, idxv) res    res = [id2, id4, id5, id6]
	// REPEATIDX(
	//     [tables, table_id, indices, name_idx],
	//     [val1, val3]
	// ) = [5, 6, 7, 10, 11, ... ]
	REPEATIDXV
	// Get possible values from index value meta
	// REPEATIDXV(idx) res    res = [val1, val2, val3, ...]
	// REPEATIDXV(
	//     [tables, table_id, indices, name_idx]
	// )  = ["alice", "bob", "foo", "bar", ... ]
)

// 0x40 range - format/output ops
const (
	ZIP OpCode = iota + 0x40
	// ZIP(tgr, new) = res
	// ZIP([f1, f2, f3], [c1, c2, c3]) = [(f1, c1), (f2, c2), (f3, c3)]
	// ZIP(
	//     [(f1, c1), (f2, c2), (f3, c3)],
	//     [(x1, (y1)), (x2, (y2)), (x3, (y3))]
	// ) = [(f1, c1, x1, (y1)), (f2, c2, x2, (y2)), (f3, c3, x2, (y2))]
	FIELD
	// (src, fields) = res
	// (
	//     [(r1f1, r1f2, r1f3), (r2f1, r2f2, r2f3),...], [2,3]
	// ) = [(r1f2, r1f3), (r2f2, r2f3), ...]
	SORT
	// SORT(src, [(field, order, null order)] ) = res
	// SORT(
	//     [(a1, a2, a3), (b1, a2, null), (a1, b2, null), ...],
	//     [(1, asc, null_first), (2, desc, null_last), (3, asc, null_last)]
	// ) = [(a1, a2, a3), (a1, b2, null), (b1, a2, null), ...]
	FILTER
	// FILTER(src, cond) = res
	// FILTER([1, 2, 3, 4, 5], [true, false, true, false, false]) = [1, 3]
	CAST
	// CAST(t1, types) t2
)

// 0x60 range - storage ops
const (
	STOREPK OpCode = iota + 0x60
	// STOREPK(pk, [uint256, uint256, ...])
	LOADPK
	// LOADPK(pk) array    array = [uint256, uint256, ...]
	STORE
	// STORE(base, ids, values, fields, idxes)
	// STORE(
	//     [tables, table_id, primary],
	//     [],
	//     [field1, field2, [field3.1, field3.2]],
	//     []
	// )
	// STORE(
	//     [tables, table_id, primary],
	//     [1, 2],
	//     [updated_field1, updated_field2],
	//     [1, 2]
	// )
	LOAD
	// LOAD(base, ids, fields) array
	// LOAD(
	//     [tables, table_id, primary], [1], []
	// ) = [(field1, field2, [field3.1, field3.2])]
	// LOAD(
	//     [tables, table_id, primary, [1], [1,3]
	// ) =[(field1, [field3.1, field3.2])］
	DELETE
	// DELETE(base, ids, idxes)
)
