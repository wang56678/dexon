package planner

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dexon-foundation/dexon/core/vm/sqlvm/schema"
)

type PlannerUtilsTestSuite struct{ suite.Suite }

func makeColumnSet(cols []uint8) ColumnSet {
	var ret ColumnSet = make([]*schema.ColumnDescriptor, len(cols))
	for i := range ret {
		ret[i] = &schema.ColumnDescriptor{
			Table:  0,
			Column: schema.ColumnRef(cols[i]),
		}
	}
	return ret
}

func (s *PlannerUtilsTestSuite) TestColumnSet() {
	{
		// Join.
		var columns, expected ColumnSet
		columns = makeColumnSet([]uint8{1, 3, 5}).Join(
			makeColumnSet([]uint8{0, 1, 2, 4, 6}))
		expected = makeColumnSet([]uint8{0, 1, 2, 3, 4, 5, 6})
		s.Require().Equal(expected, columns)
		columns = makeColumnSet([]uint8{1, 3, 5}).Join(
			makeColumnSet([]uint8{3, 5}))
		expected = makeColumnSet([]uint8{1, 3, 5})
		s.Require().Equal(expected, columns)
		columns = ColumnSet{}.Join(makeColumnSet([]uint8{0}))
		expected = makeColumnSet([]uint8{0})
		s.Require().Equal(expected, columns)
		columns = makeColumnSet([]uint8{1}).Join(
			makeColumnSet([]uint8{1, 3}))
		expected = makeColumnSet([]uint8{1, 3})
		s.Require().Equal(expected, columns)
		columns = makeColumnSet([]uint8{5}).Join(makeColumnSet([]uint8{1, 3}))
		expected = makeColumnSet([]uint8{1, 3, 5})
		s.Require().Equal(expected, columns)
	}
	{
		// Equal.
		var equal bool
		// True cases.
		equal = ColumnSet{}.Equal(ColumnSet{})
		s.Require().True(equal)
		equal = makeColumnSet([]uint8{1, 2}).Equal(makeColumnSet([]uint8{1, 2}))
		s.Require().True(equal)
		// False cases.
		equal = ColumnSet{}.Equal(makeColumnSet([]uint8{1, 2}))
		s.Require().False(equal)
		equal = makeColumnSet([]uint8{1, 2}).Equal(ColumnSet{})
		s.Require().False(equal)
		equal = makeColumnSet([]uint8{2}).Equal(makeColumnSet([]uint8{1}))
		s.Require().False(equal)
		equal = makeColumnSet([]uint8{2}).Equal(makeColumnSet([]uint8{1, 3}))
		s.Require().False(equal)
		equal = makeColumnSet([]uint8{1, 3}).Equal(makeColumnSet([]uint8{2}))
		s.Require().False(equal)
	}
	{
		// Contains.
		var contains bool
		// True cases.
		contains = ColumnSet{}.Contains(ColumnSet{})
		s.Require().True(contains)
		contains = makeColumnSet([]uint8{1, 2}).Contains(ColumnSet{})
		s.Require().True(contains)
		contains = makeColumnSet([]uint8{1, 2}).Contains(
			makeColumnSet([]uint8{2}))
		s.Require().True(contains)
		contains = makeColumnSet([]uint8{1, 2}).Contains(
			makeColumnSet([]uint8{1}))
		s.Require().True(contains)
		contains = makeColumnSet([]uint8{1, 2, 3}).Contains(
			makeColumnSet([]uint8{1, 2}))
		s.Require().True(contains)
		// False cases.
		contains = makeColumnSet([]uint8{1}).Contains(makeColumnSet([]uint8{2}))
		s.Require().False(contains)
		contains = makeColumnSet([]uint8{2}).Contains(
			makeColumnSet([]uint8{1, 2}))
		s.Require().False(contains)
		contains = makeColumnSet([]uint8{1}).Contains(
			makeColumnSet([]uint8{1, 2}))
		s.Require().False(contains)
		contains = makeColumnSet([]uint8{1, 3, 5}).Contains(
			makeColumnSet([]uint8{4}))
		s.Require().False(contains)
	}
	{
		// IsDisjoint.
		var disjoin bool
		// True cases.
		disjoin = ColumnSet{}.IsDisjoint(ColumnSet{})
		s.Require().True(disjoin)
		disjoin = ColumnSet{}.IsDisjoint(makeColumnSet([]uint8{1}))
		s.Require().True(disjoin)
		disjoin = makeColumnSet([]uint8{1}).IsDisjoint(ColumnSet{})
		s.Require().True(disjoin)
		disjoin = makeColumnSet([]uint8{1}).IsDisjoint(
			makeColumnSet([]uint8{2}))
		s.Require().True(disjoin)
		// False cases.
		disjoin = makeColumnSet([]uint8{1, 2}).IsDisjoint(
			makeColumnSet([]uint8{2}))
		s.Require().False(disjoin)
		disjoin = makeColumnSet([]uint8{1, 2}).IsDisjoint(
			makeColumnSet([]uint8{1}))
		s.Require().False(disjoin)
		disjoin = makeColumnSet([]uint8{1, 2}).IsDisjoint(
			makeColumnSet([]uint8{0, 2}))
		s.Require().False(disjoin)
		disjoin = makeColumnSet([]uint8{1, 7}).IsDisjoint(
			makeColumnSet([]uint8{5, 6, 7}))
		s.Require().False(disjoin)
	}
}

func TestPlannerUtils(t *testing.T) {
	suite.Run(t, new(PlannerUtilsTestSuite))
}
