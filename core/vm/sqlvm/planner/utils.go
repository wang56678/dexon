package planner

import "github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"

func bytesEq(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func findTableIdxByName(tables schema.Schema, name []byte) (
	schema.TableRef, bool) {

	for i, table := range tables {
		if bytesEq(name, table.Name) {
			return schema.TableRef(i), true
		}
	}
	return 0, false
}

func findColumnIdxByName(table *schema.Table, name []byte) (
	schema.ColumnRef, bool) {

	for i, c := range table.Columns {
		if bytesEq(name, c.Name) {
			return schema.ColumnRef(i), true
		}
	}
	return 0, false
}

// ColumnSet is a sorted slice of column idxs.
type ColumnSet []*schema.ColumnDescriptor

func compareTableRef(a, b schema.TableRef) int {
	switch {
	case a > b:
		return 1
	case a == b:
		return 0
	default:
		return -1
	}
}

func compareColumnRef(a, b schema.ColumnRef) int {
	switch {
	case a > b:
		return 1
	case a == b:
		return 0
	default:
		return -1
	}
}

func compareColumn(a, b *schema.ColumnDescriptor) int {
	comp := compareTableRef(a.Table, b.Table)
	switch comp {
	case 0:
		return compareColumnRef(a.Column, b.Column)
	default:
		return comp
	}
}

// Join creates a new set which is the union of c and other.
func (c ColumnSet) Join(other ColumnSet) ColumnSet {
	ret := make([]*schema.ColumnDescriptor, 0, len(c)+len(other))
	i, j := 0, 0
	for i != len(c) && j != len(other) {
		comp := compareColumn(c[i], other[j])
		switch comp {
		case 0:
			ret = append(ret, c[i])
			i++
			j++
		case 1:
			ret = append(ret, other[j])
			j++
		case -1:
			ret = append(ret, c[i])
			i++
		}
	}
	for i != len(c) {
		ret = append(ret, c[i])
		i++
	}
	for j != len(other) {
		ret = append(ret, other[j])
		j++
	}
	return ret
}

// Equal compares the two sets.
func (c ColumnSet) Equal(other ColumnSet) bool {
	if len(c) != len(other) {
		return false
	}
	for i := range c {
		if compareColumn(c[i], other[i]) != 0 {
			return false
		}
	}
	return true
}

// IsDisjoint checks if the two sets are disjoint.
func (c ColumnSet) IsDisjoint(other ColumnSet) bool {
	i, j := 0, 0
	for i != len(c) && j != len(other) {
		comp := compareColumn(c[i], other[j])
		switch comp {
		case 0:
			return false
		case 1:
			j++
		case -1:
			i++
		}
	}
	return true
}

// Contains checks if other is a subset of c.
func (c ColumnSet) Contains(other ColumnSet) bool {
	i, j := 0, 0
	for i != len(c) && j != len(other) {
		comp := compareColumn(c[i], other[j])
		switch comp {
		case 1:
			// Found some item not in c.
			return false
		case 0:
			j++
		}
		i++
	}
	return j == len(other)
}
