# SQLVM Design Doc

   * [SQLVM Design Doc](#sqlvm-design-doc)
      * [**Objective**](#objective)
      * [**Background**](#background)
      * [**Detailed design**](#detailed-design)
         * [Supported SQL subset grammar](#supported-sql-subset-grammar)
         * [Built-in functions](#built-in-functions)
         * [Cross-VM interaction (Internal transactions)](#cross-vm-interaction-internal-transactions)
         * [Permission control](#permission-control)
         * [Limitation](#limitation)
         * [ABI](#abi)
         * [Auto-generated functions for Database (WIP)](#auto-generated-functions-for-database-wip)
         * [Auto-generated functions for Tables (WIP)](#auto-generated-functions-for-tables-wip)
         * [Supported field types](#supported-field-types)
            * [Aliases](#aliases)
            * [Type ID allocation (2 bytes)](#type-id-allocation-2-bytes)
            * [Solidity type conversion](#solidity-type-conversion)
            * [Solidity integer conversion rules](#solidity-integer-conversion-rules)
            * [SQLVM type conversion](#sqlvm-type-conversion)
            * [Types of literals](#types-of-literals)
            * [Cast format](#cast-format)
         * [Unsupported field types](#unsupported-field-types)
         * [Table to key-value mapping](#table-to-key-value-mapping)
         * [Database metadata](#database-metadata)
         * [Table metadata](#table-metadata)
            * [Primary index &amp; data](#primary-index--data)
               * [List of keys](#list-of-keys)
               * [Actual data](#actual-data)
            * [Other indices](#other-indices)
               * [List of keys](#list-of-keys-1)
               * [Actual data](#actual-data-1)
            * [Sequence(auto increment)](#sequenceauto-increment)
            * [NULL value](#null-value)
            * [Compound index](#compound-index)
            * [Contract owner](#contract-owner)
            * [Table writer](#table-writer)
               * [List of keys](#list-of-keys-2)
               * [Actual data](#actual-data-2)
         * [Query planning](#query-planning)
            * [(TBD) Query planning algorithms](#tbd-query-planning-algorithms)
         * [(TBD) Pricing model](#tbd-pricing-model)
         * [Difference with SQL spec](#difference-with-sql-spec)
         * [Instruction set](#instruction-set)
            * [Base structs](#base-structs)
            * [Codes](#codes)
               * [Examples](#examples)
               * [Storage ops details](#storage-ops-details)
                  * [INSERT](#insert)
                  * [UPDATE](#update)
                  * [DELETE](#delete)
      * [**Corner cases**](#corner-cases)
      * [**Miscellaneous**](#miscellaneous)
      * [**TODO list**](#todo-list)
      * [**Authors**](#authors)

<!--- Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)  -->

## **Objective**

The objective is to provide the programmers with a friendly SQL interface which they are familiar with.  SQLVM is implemented alongside with EVM, such that DEXON DApp developer can use SQL syntax to handle their data. It would support a subset of SQL operations.

## **Background**

1. How solidity store data. [https://solidity.readthedocs.io/en/v0.4.24/miscellaneous.html](https://solidity.readthedocs.io/en/v0.4.24/miscellaneous.html)

1. How to map SQL syntax to key-value [https://www.cockroachlabs.com/blog/sql-in-cockroachdb-mapping-table-data-to-key-value-storage/](https://www.cockroachlabs.com/blog/sql-in-cockroachdb-mapping-table-data-to-key-value-storage/)

## **Detailed design**

### Supported SQL subset grammar

Core function reference: [https://www.sqlite.org/lang_corefunc.html](https://www.sqlite.org/lang_corefunc.html)

Aggregation function reference: [https://www.sqlite.org/lang_aggfunc.html](https://www.sqlite.org/lang_aggfunc.html)

Golang Planner reference: [https://gitlab.com/cznic/ql/blob/master/plan.go](https://gitlab.com/cznic/ql/blob/master/plan.go)
```sql
CREATE TABLE {name} {column-spec}
CREATE INDEX {name} ON {table} {columns}
SELECT {items} FROM {table} WHERE {condition} ORDER BY {expr}
UPDATE {table} SET {items} WHERE {condition}
DELETE FROM {table} WHERE {condition}
INSERT INTO {table} {columns} VALUES {values}
INSERT INTO {table} DEFAULT VALUES
```

### Built-in functions

```
BLOCKNUMBER() uint256
NOW() / BLOCKTIMESTAMP() uint256
BLOCKHASH(uint256) bytes32
BLOCKCOINBASE() address
BLOCKGASLIMIT() uint256
MSGSENDER() address
MSGDATA() bytes
TXORIGIN() address
RAND() uint256
BITAND(a, b T), T ∈ {uintX, intX, bytesX}
BITOR(a, b T), T ∈ {uintX, intX, bytesX}
BITNOT(a T), T ∈ {uintX, intX, bytesX}
BITXOR(a, b T), T ∈ {uintX, intX, bytesX}
OCTET_LENGTH(a T) N, T ∈ {bytes, bytesX}, N ∈ {uintX, intX}
SUBSTRING(str FROM pos FOR len) bytes, str ∈ {bytes}, pos, len ∈ {uintX, intX}
```

### Cross-VM interaction (Internal transactions)

To support cross-VM interaction, all VM must confront the [Solidity ABI specification](https://solidity.readthedocs.io/en/v0.5.2/abi-spec.html).

The transaction data sent to the database contract will be just plain SQL expressions. The returned rows will be encoded as a 2-level array, where the outer layer being the row data, and the inner layer being the data fields. We can consider modifying solidity so it could deserialize structures from function returns (or wait for the ABIv2 encoder). ([PoC](https://gist.github.com/aitjcize/51aeb7023febaa259216e6f5069a6e63))

### Permission control

By default, there is a writer list for each table, and the owner can edit the lists.
The owner can transfer the entire database to other owners.

### Limitation

* Max table number: 256
* Max column number: 256
* Max foreign key number: 256
* Max SELECT fields number: 65536
* Max records number: 2⁶⁴
* Max index number in a table: 256
* Max compound index field number: 256

### ABI

The ABI of the SQL contract will look just like an Ethereum ABI, except that all of the methods are auto-generated when the contract (schema) is deployed.

### Auto-generated functions for Database (WIP)

(Requires the [Solidity ABI v2 encoder](https://blog.ricmoo.com/solidity-abiv2-a-foray-into-the-experimental-a6afd3d47185))
```
function tables() public view returns ([]string)
function owner() public view returns (address)
function transferOwnership(address newOwner) onlyContractOwner public
function query(string statement) public view returns (uint256 error, bytes data)
function exec(string statement) public returns (uint256 error, uint256 affectedRows)
```

### Auto-generated functions for Tables (WIP)
```
function grantWriter(string table, address writer) onlyContractOwner public
function revokeWriter(string table, address writer) onlyContractOwner public
function fields(string table) public view returns ([]string)
function indices(string table) public view returns ([]string)
```

### Supported field types
**int{X}**, **uint{X}**, **bytes{X}**, **bytes**, **bool**, **address**

notes: \
**int{X}**: X=8..256, in step of 8. i.e. **int8**, **int16**, … **int256** \
**uint{X}**: X=8..256, in step of 8. i.e. **uint8**, **uint16**, … **uint256** \
**bytes{X}**: X=1..32. i.e. **bytes1**, **bytes2** … **bytes32**

#### Aliases
**boolean** = **bool** \
**byte** = **bytes1**

#### Type ID allocation (2 bytes)
|          | 1st byte | 2nd byte  | Number of types |
|-         |-         | -         | -               |
| unknown  | 0x00     | 0x00      | 1               |
| null     | 0x01     | 0x00      | 1               |
| *        | 0x01     | 0x01      | 1               |
| default  | 0x01     | 0x02      | 1               |
| bool     | 0x02     | 0x00      | 1               |
| address  | 0x03     | 0x00      | 1               |
| int      | 0x04     | 0x00~0x1f | 32              |
| uint     | 0x05     | 0x00~0x1f | 32              |
| bytes{X} | 0x06     | 0x00~0x1f | 32              |
| bytes    | 0x07     | 0x00      | 1               |

#### Solidity type conversion
| From/To | int                  | bytes{X}                   | bytes | bool    | address           |
| -       | -                    | -                          | -     | -       | -                 |
| int     | see the rule         | same size big-endian       | ✘     | ✘       | use integer rules |
| bytes{X}| same size big-endian | pad/crop zero on the right | ✘     | ✘       | only bytes20      |
| bytes   | ✘                    | ✘                          | ✔     | ✘       | ✘                 |
| bool    | ✘                    | ✘                          | ✘     | ✔       | ✘                 |
| address | use integer rules    | only bytes20               | ✘     | ✘       | ✔                 |


#### Solidity integer conversion rules
**int{X}** → **int{Y}**, X < Y: sign extend \
**int{X}** → **uint{Y}**, X < Y: sign extend \
**uint{X}** → **int{Y}**, X < Y: zero extend \
**uint{X}** → **uint{Y}**, X < Y: zero extend

**int{X}** → **int{Y}**, X > Y: discard higher bits \
**int{X}** → **uint{Y}**, X > Y: discard higher bits \
**uint{X}** → **int{Y}**, X > Y: discard higher bits \
**uint{X}** → **uint{Y}**, X > Y: discard higher bits

**address** has the same rule as **uint160**.

#### SQLVM type conversion
| From/To | int                  | bytes{X}                   | bytes | bool               | address           |
| -       | -                    | -                          | -     | -                  | -                 |
| int     | see the rule         | same size big-endian       | ✘     | true if source ≠ 0 | use integer rules |
| bytes{X}| same size big-endian | pad/crop zero on the right | ✔     | ✘                  | only bytes20      |
| bytes   | ✘                    | pad/crop zero on the right | ✔     | ✘                  | ✘                 |
| bool    | ✔                    | ✘                          | ✘     | ✔                  | ✘                 |
| address | use integer rules    | only bytes20               | ✘     | ✘                  | ✔                 |

#### Types of literals
```
TRUE / FALSE: bool
123 / 1. / 1.2e3 (decimal numbers): int / uint / fixed / ufixed (default to int256 / fixed128x18)
1.2 / .2 / 1.23e1 (decimal numbers): fixed / ufixed (default to fixed128x18)
0x12abc (hexadecimal numbers): int / uint (default to uint256)
‘Hello’ (normal strings): bytes / bytes{X} (default bytes)
hex’abcd1234’ (hexadecimal strings): bytes / bytes{X} (default bytes)
```

#### Cast format
```sql
CAST(expression AS type)
```

### Unsupported field types
```
fixed{M}x{N} / ufixed{M}x{N}: M must be divisible by 8 and goes from 8 to 256. N must be
between 0 and 80, inclusive.
```

### Table to key-value mapping

Similar to the design used in CockroachDB:

[https://www.cockroachlabs.com/blog/sql-in-cockroachdb-mapping-table-data-to-key-value-storage/](https://www.cockroachlabs.com/blog/sql-in-cockroachdb-mapping-table-data-to-key-value-storage/)

Data mapping including struct/list is being encoded using the solidity state layout: [https://solidity.readthedocs.io/en/v0.4.24/miscellaneous.html](https://solidity.readthedocs.io/en/v0.4.24/miscellaneous.html)

The key path is encoded using RLP, then pass through Keccak256. We define
```
PathKey(path) = Keccak256(RLPEncode(path))
```

### Database metadata

Metadata will be stored in `contract.Code`.
The first 4 bytes are reserved to note the SQLVM version,
and the rest will store the schema in RLP encoding.

### Table metadata

```golang
// Schema defines sqlvm schema struct.
type Schema []Table

// Table defiens sqlvm table struct.
type Table struct {
	Name    []byte
	Columns []Column
	Indices []Index
}

// IndexAttr defines bit flags for describing index attribute.
type IndexAttr uint16

const (
	// IndexAttrUnique indicates whether an index is unique.
	IndexAttrUnique IndexAttr = 1 << iota
	// IndexAttrReferenced indicates whether an index is referenced by columns
	// with foreign key constraints. This attribute cannot be specified by
	// users. It is computed automatically during contract creation.
	IndexAttrReferenced
)

// TableRef defines the type for table index in Schema.
type TableRef uint8

// ColumnRef defines the type for column index in Table.Columns.
type ColumnRef uint8

// IndexRef defines the type for array index of Column.Indices.
type IndexRef uint8

// SequenceRef defines the type for sequence index in Table.
type SequenceRef uint8

// IndexAttr defines bit flags for describing index attribute.
type IndexAttr uint16

// Index defines sqlvm index struct.
type Index struct {
	Name    []byte
	Attr    IndexAttr
	Columns []ColumnRef // Columns must be sorted in ascending order.
}

type ColumnAttr uint16

const (
	// ColumnAttrPrimaryKey is a no-op. Primary key constraints are converted
	// to a unique index during contract creation.
	ColumnAttrPrimaryKey ColumnAttr = 1 << iota
	// ColumnAttrNotNull is a no-op. We have not supported NULL values so all
	// columns are implicitly non-null.
	ColumnAttrNotNull
	// ColumnAttrUnique is a no-op. Unique constraints are converted to unique
	// indices during contract creation.
	ColumnAttrUnique
	// ColumnAttrHasDefault indicates whether a column has a default value. The
	// default value does not affect the starting value of AUTOINCREMENT.
	ColumnAttrHasDefault
	// ColumnAttrHasForeignKey indicates whether a column references a column
	// on a different table.
	ColumnAttrHasForeignKey
	// ColumnAttrHasSequence indicates whether a column is declared with
	// AUTOINCREMENT. It is only valid on integer fields.
	ColumnAttrHasSequence
)

// ColumnDescriptor identifies a column in a schema by array indices.
type ColumnDescriptor struct {
	Table  TableRef
	Column ColumnRef
}

// Column defines sqlvm index struct.
type Column struct {
	Name        []byte
	Type        ast.DataType
	Attr        ColumnAttr
	ForeignKeys []ColumnDescriptor
	Sequence    SequenceRef
	Default     []byte // serialized default value
}

// All metadata are encoded with RLPEncode.
CODE_SECTION = RLPEncode(Schema)
```

#### Primary index & data

##### List of keys

Key
```
PathKey(["tables", uint8(table_idx), "primary"])
```

Value
```golang
type PrimaryKeys struct {
	// Slot 1: Header.
	LastRowID   uint64
	RowCount    uint64
	<unused>    uint64
	<unused>    uint64
	// End of header.
	Slot2       uint256 // Bitmap, record 0 ~ 255.
	Slot3       uint256 // Bitmap, record 256 ~ 511.
	Slot4       uint256 // Bitmap, record 512 ~ 767.
    ...
}
```

Bitmap:

*   Use `uint64` to store primary index.
*   Pagination: 1 slot = 1 page
*   Reserve 1st slot for some information:
    *   Used slot number (0-63 bit)
    *   Record count (64-127 bit)
*   Not reuse bits
    * Only last bit will be reused.


##### Actual data
Key
```
PathKey(["tables", uint8(table_idx), "primary", uint64({row_id})])
```

Value
```
PackingEncode([column_1, column_2, column_3, ...])
```
`PackingEncode` will pack data if they can be put in one slot. It stores data slot by slot.
Dynamic bytes follow Solidity ABI encoding.


#### Other indices
##### List of keys
Key
```
PathKey(["tables", uint8(table_idx), "indices", uint8(index_idx)])
```

Value
```golang
type PossibleValues struct {
	// Slot 1.
	Length          uint64
	<unused>        uint64
	<unused>        uint64
	<unused>        uint64
	// Slot 2.
	ValueAddr1      uint256
	// Slot 3.
	ValueAddr2      uint256
	...
}
```
##### Actual data
Key
```
PathKey(["tables", uint8(table_idx), "indices", uint8(index_idx), field_1, field_2, field_3, ...])
```
Value
```golang
type IndexKeys struct {
	// Slot 1.
	Length                      uint64
	IndexToPossibleValuesOffset uint64
	ForeignKeyReferenceCount    uint64
	<unused>                    uint64
	// Slot 2.
	RowID1                      uint64
	RowID2                      uint64
	RowID3                      uint64
	RowID4                      uint64
	// Slot 3.
	RowID5                      uint64
	...
}
```

#### Sequence(auto increment)
Key
```
PathKey(["tables", uint8(table_idx), "sequence", uint8(sequence_idx)])
```

Value
```
uint64
```

#### NULL value
Do not support NULL value.

#### Compound index
Only (a, b, c) will be indexed.
(a), (b) or (a, b) won't be indexed automatically.

#### Contract owner
Key
```
PathKey(["owner"])
```

Value
```
address
```

#### Table writer

##### List of keys

Key
```
PathKey(["tables", uint8(table_idx), "writers"])
```

Value
```golang
type PossibleValues struct {
	// Slot 1.
	Length          uint64
	<unused>        uint64
	<unused>        uint64
	<unused>        uint64
	// Slot 2.
	Addr1           uint160
	<unused>        12 bytes
	// Slot 3.
	Addr2           uint160
	<unused>        12 bytes
	...
}
```

##### Actual data
Key
```
PathKey(["tables", uint8(table_idx), "writers", "addr"])
```
Value
```golang
type TableWriter struct {
	// Slot 1.
	IndexToPossibleValuesOffset    uint64
	<unused>                       uint64
	<unused>                       uint64
	<unused>                       uint64
}
```


### Query planning

The design of Ethereum state uses secure trie. Keys are first passed through Keccak256 to be converted to the actual key. This makes caching impossible, as the output of Keccak256 is proved to be randomized. We may want to consider bypassing secure tree for SQLVM's contract storage.

Since we probably only going to support a subset of SQL2011 (no JOIN), the query planning should be relatively simple. For the query to be efficient, most of the query should have an index, else a table scan will be triggered and it will probably be too slow to ever be usable.

#### (TBD) Query planning algorithms

CockroachDB cost-based SQL optimizer reference:
[https://www.cockroachlabs.com/blog/building-cost-based-sql-optimizer/](https://www.cockroachlabs.com/blog/building-cost-based-sql-optimizer/)

Sqlite query planner reference:
[https://www.sqlite.org/queryplanner.html](https://www.sqlite.org/queryplanner.html)


### (TBD) Pricing model

Charge gas by operations.


### Difference with SQL spec

1. No '' escape in strings and no "" escape in identifiers.
1. \ is an escape character in strings and identifiers.
1. No implicit type conversion.
1. No comment.
1. No multi-line string such as
```
'abc'
'def'
```
however,
```
'abc' ||
'def'
```
is allowed.


### Instruction set

#### Base structs

```golang
// Operand would be array-based value associated with meta to describe type of
// array element.
type Operand struct {
	IsImmediate   bool
	Meta          []ast.DataType
	Data          []Tuple
	RegisterIndex uint
}

// Tuple is collection of Raw.
type Tuple []*Raw

// Raw consist of decimal and byte slice which represents the real value
// of basic operand unit. Only dynmaic bytes, fixed bytes and address will be
// stored in []byte, and the other data type will be stored in decimal.
type Raw struct {
	Value decimal.Decimal
	Bytes []byte
}
```

#### Codes
```golang
// For example:
// * [] represents the data array
// * () represents multiple value tuple
// * value without brackets represents the single value tuple

// 0x10 range - arithmetic ops. (array-based)
const (
	ADD = iota + 0x10
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
	CONCAT
	// CONCAT(s1, s2) = s3
	NEG
	// NEG(t1) = -t1
)

// 0x20 range - comparison ops.
const (
	LT = iota + 0x20
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
	// LIKE(t1, pattern[, escape]) res
	// Immediate case
	// res = Like(t1, '%abc%') =
	//       ['_abc_', '123'] like ['%abc%'] = [true, false]
	// Not immediate case
	// res = Like(t1, t2, t3) =
	//      ['_abc', '1%23'] like ['%abc%', '_\%2_'] escape ["", "\"] =
	//      [true, true]
)

// 0x30 range - pk/index/field meta ops
const (
	REPEATPK = iota + 0x30
	// REPEATPK(table_id) res
	// res = [id1, id2, id3, ...]
	// REPEATPK(table_id) = [1, 2, 3, 5, 6, 7, ...]
	REPEATIDX
	// Scan given index value(s)
	// REPEATIDX(table_id, name_idx, idxv) res
	// res = [id2, id4, id5, id6]
	// REPEATIDX(
	//  [table_id, name_idx],
	//  [val1, val3]
	// ) = [5, 6, 7, 10, 11, ... ]
	REPEATIDXV
	// Get possible values from index value meta
	// REPEATIDXV(table_id, name_idx) res
	// res = [val1, val2, val3, ...]
	// REPEATIDXV(
	//   [table_id, name_idx]
	// ) = ["alice", "bob", "foo", "bar", ... ]
)

// 0x40 range - format/output ops
const (
	ZIP = iota + 0x40
	// ZIP(tgr, new) = res
	// ZIP([f1, f2, f3], [c1, c2, c3]) = [(f1, c1), (f2, c2), (f3, c3)]
	// ZIP(
	//    [(f1, c1), (f2, c2), (f3, c3)],
	//    [(x1, y1), (x2, y2), (x3, y3)]
	// ) = [(f1, c1, x1, y1), (f2, c2, x2, y2), (f3, c3, x3, y3)]
	FILED
	// FIELD(src, fields) = res
	// FIELD(
	//    [(r1f0, r1f1, r1f2), (r2f0, r2f1, r2f2),...], [(1, 2)]
	// ) = [(r1f1, r1f2), (r2f1, r2f2), ...]
	PRUNE // in-place op
	// PRUNE(src, fields) = res
	// PRUNE(
	//     [(r1f0, r1f1, r1f2), (r2f0, r2f1, r2f2),...], [1]
	// ) = [(r1f0, r1f2), (r2f0, r2f2), ...]
	SORT // in-place op
	// SORT(src, (ascending bool, field idx))
	// SORT(src, [(ascending, field)] ) = res
	// SORT(
	//   [(a1, a2), (b1, a2), (a1, b2), ...],
	//   [(asc, 0), (desc, 1), (asc, 2)]
	// ) = [(a1, a2), (a1, b2), (b1, a2), ...]
	FILTER
	// FILTER(src, cond) = res
	// FILTER([1, 2, 3, 4, 5], [true, false, true, false, false]) = [1, 3]
	CAST
	// CAST(t1, type) t2
	CUT
	// in-place op
	// CUT(src, range) = res
	// CUT(src, (start[, end]))
	// CUT(
	//      [(r1f1, r1f2, r1f3, r1f4),
	//      (r2f1, r2f2, r2f3, r2f4), ...], [(1, 2)]
	// ) = [(r1f1, r1f4), (r2f1, r2f4), ...]
	RANGE // in-place op
	// RANGE(src, range) = res
	// RANGE[src, (offset[, limit])]
	// RANGE([r1, r2, r3, r4, r5, r6], [(2, 3)]) = [r3, r4, r5]
)

// 0x50 range - function ops
const (
	FUNC = iota + 0x50
	// FUNC(t1, func id[, args...])) = res
)

// 0x60 range - storage ops
const (
	INSERT = iota + 0x60
	// INSERT(table_id, fields, values...) = res
	// the number of values depends on the length of fields
	// res not important
	// INSERT(
	//    table_id,
	//    [0, 2, 5],
	//    [field0],
	//    [field2],
	//    [Immediate value],
	// ) = _
	UPDATE
	// UPDATE(table_id, ids, fields, values...) = res
	// the number of values depends on the length of fields
	// res = uint64(affected row)
	// UPDATE(
	//    table_id,
	//    [0, 99],
	//    [1, 2, 3],
	//    [updated_field0-1, updated_field99-1],
	//    [updated_field0-2, updated_field99-2],
	//    [Immediate Value],
	// ) = uint64(2)
	LOAD
	// LOAD(table_id, ids, fields) = res
	// LOAD(
	//   table_id,
	//   [55, 66],
	//   [1, 2, 3],
	// ) = [
	//   (field55-1, field55-2, field55-3),
	//   (field66-1, field66-2, field66-3),
	// ]
	DELETE
	// DELETE(table_id, ids) = res
	// res = uint64(affected row)
	// DELETE(
	//    table_id,
	//    [1, 2, 3, 100],
	// ) = uint64(4)
)
```

##### Examples
*   SELECT * FROM a
    *   REPEATPK(a) t1
    *   LOAD(a, t1, [0,1,2,3,4,5]) t2
    *   ZIP(t2)
*   SELECT f1, f3, f5 FROM a WHERE f1=val (index hit)
    *   REPEATIDXV(a_f1_idx) t1
    *   EQ(t1, val) t3
    *   FILTER(t1, t3) t4
    *   REPEATIDX(a_f1_idx, t4) t5
    *   LOAD(a, t5, [1,3,5]) t6
    *   ZIP(t6)
*   SELECT f1, f3, f5 FROM a WHERE f1 > val (index scan)
    *   REPEATIDXV(a_f1_idx) t1
    *   GT(t1, val) t2
    *   FILTER(t1, t2) t3
    *   REPEATIDX(a_f1_idx, t3) t4
    *   LOAD(a, t4, [1, 3, 5]) t5
    *   ZIP(t5)
*   SELECT f1 > f2, f5 FROM a WHERE f3 > 10 (field operation)
    *   REPEATPK(a) t1
    *   LOAD(a, t1, [3]) t2
    *   GT(t2, 10) t4
    *   FILTER(t1, t4) t5
    *   LOAD(a, t5, [1, 2, 5]) t2
    *   FIELD(t2, [1]) t3
    *   FIELD(t2, [2]) t4
    *   GT(t3, t4) t5
    *   FIELD(t2, [5]) t9
    *   ZIP(t5, t9) t10
*   DELETE FROM a WHERE f2 = 'c' (idx hit)
    *   REPEATIDX(a_f2_idx, ['c']) t1
    *   DELETE(a, t1)
*   UPDATE a SET f1 = 'new'
    *   REPEATPK(a) t1
    *   UPDATE(a, t1, [1], ['new'])
*   INSERT a VALUES (f1, f2, f3), (g1, g2, g3);
    *   INSERT(a, [1,2,3], (f1, f2, f3))
    *   INSERT(a, [1,2,3], (g1, g2, g3))
*   INSERT INTO a (c0,c3,c5) VALUES (NOW(), MSGSENDER(), RAND());
    *   FUNC('NOW') t1          // t1 = [(block_height)]
    *   FUNC('MEGSENDER') t2    // t2 = [(address     )]
    *   FUNC('RAND') t3         // t3 = [(rand_num    )]
    *   ZIP(t1, t2) t1          // t1 = [(block_height, address)]
    *   ZIP(t1, t3) t1          // t1 = [(block_height, address, rand_num)]
    *   INSERT(a, [0,3,5], t1)

##### Storage ops details

###### INSERT

Row steps:
1. Acquire IDs from generator
1. Add new ID to pk
1. Update auto increment field if in need
1. Update default field if in need
1. Check foreign key exists
	1. Increase target reference count
1. Check index
	1. Create empty index if not existing
	1. If conflict with other unique index, return error
1. Update index
1. Insert data
1. Commit


###### UPDATE

Row steps:
1. Get old data by ID
1. Iterate IDs
	1. Check foreign key exists
		1. Increase target reference count
		1. Decrease old target reference count
	1. Update index
		1. Add new index
		1. Remove old index
    1. Update data
1. Checking unique index conflict. If conflict, return error
1. Commit

###### DELETE

Row steps:
1. Get old data by IDs
1. Iterate IDs
    1. Iterate indices
        1. If index contains more than 1 keys, remove from list
        1. If index contains only key
            1. If reference count not zero, return error
            1. Delete key and meta value
    1. Check foreign key exists
        1. Decrease target reference count
1. Commit

## **Corner cases**

This section contains some corner cases found in design stage. Can be used for test cases in integration test.

1. `select 1 from table where a > 1;`  (has problem in determining row number in ZIP op)
1. `select 1, a, 1 from table;` (similar case as 1.)
1. `select 1 where random() > 1;` (no table in select and the condition is not constant)
1. `select random() from table order by 1;` (support of column reference)
1. `select random() from table order by - -1;` (it's column reference, not expression)
1. `select random() from table order by random();` (in sqlite, it is equivalent to 'order by 1')
   1. those two random() are different.
1. `insert into table (a) values (random()+1) (1) (random());` (expression in insert)
1. `select * from table where column > random();` (we cannot use index on column as random is in condition)


## **Miscellaneous**

* If there is constraint conflict, abort transaction.
* (TBD) There is a default return limit for avoiding out of gas.
* There is a flag to decide to enable flow check or not.
* It can not transfer token to SQL contract address.
* Do not support expression default value.
* AUTO INCR use overall max value, can't auto increment with foreign key, only int/uint can use auto increment, dominate over 'default' setting `
* Unique check (in memory check, i.e. rehearsal in memory)`
* Dynamic Bytes support:
    * `OCTET_LENGTH`
    * `SUBSTRING`
* Use bytes to store fixed bytes data
* Fixed bytes supports:
    * `OCTET_LENGTH`
    * `fixBytes3Data || fixBytes5Data`
        * Return bytes8
        * Abort if larger than 32
    * CAST fixBytes to DynamciByte and Dynamicbyte to fixBytes
    * `fixBytes5Data [=|>|<] fixBytesData5`
    * `fixByte5Data bitwise op fixBytes5Data`
        *   `op: BITAND, BITOR, BITXOR, BITNOT`
* Fixed bytes do not supports:
    * `SUBSTRING`
    * `… WHERE fixBytesData LIKE XXX`
      * Do not support it. It should cast to dynamic bytes before LIKE
    * `SELECT fixBytes5Data [+|-|*|/] fixBytes5Data`
        * Solidity does not support arithmetic op for fixed bytes directly. User should cast to uint to do arithmetic op.

* Uint/Int support
    * `Bitwise op function`
* LIKE escape can only use bytes1.

## **TODO list**
*   Support(temporary save data) like join, subquery, streaming
*   Aggregation
    *   COUNT()
    *   DISTINCT
    *   SUM()
    *   MIN()/MAX()

## **Authors**
* [Wei-Ning Huang](https://github.com/aitjcize)
* [wmin0](https://github.com/wmin0)
* [lantw44](https://github.com/lantw44)
* [yenlinlai](https://github.com/yenlinlai)
* [Meng-Ying Yang](https://github.com/myyang)
* [Jhih-Ming Huang](https://github.com/jm-cobinhood)
