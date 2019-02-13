package planner

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PlannerUtilsTestSuite struct{ suite.Suite }

func (s *PlannerUtilsTestSuite) TestColumnSet() {
	{
		// Join.
		var columns, expected ColumnSet
		columns = ColumnSet{1, 3, 5}.Join(ColumnSet{0, 1, 2, 4, 6})
		expected = ColumnSet{0, 1, 2, 3, 4, 5, 6}
		s.Require().Equal(expected, columns)
		columns = ColumnSet{1, 3, 5}.Join(ColumnSet{3, 5})
		expected = ColumnSet{1, 3, 5}
		s.Require().Equal(expected, columns)
		columns = ColumnSet{}.Join(ColumnSet{0})
		expected = ColumnSet{0}
		s.Require().Equal(expected, columns)
		columns = ColumnSet{1}.Join(ColumnSet{1, 3})
		expected = ColumnSet{1, 3}
		s.Require().Equal(expected, columns)
		columns = ColumnSet{5}.Join(ColumnSet{1, 3})
		expected = ColumnSet{1, 3, 5}
		s.Require().Equal(expected, columns)
	}
	{
		// Equal.
		var equal bool
		// True cases.
		equal = ColumnSet{}.Equal(ColumnSet{})
		s.Require().True(equal)
		equal = ColumnSet{1, 2}.Equal(ColumnSet{1, 2})
		s.Require().True(equal)
		// False cases.
		equal = ColumnSet{}.Equal(ColumnSet{1, 2})
		s.Require().False(equal)
		equal = ColumnSet{1, 2}.Equal(ColumnSet{})
		s.Require().False(equal)
		equal = ColumnSet{2}.Equal(ColumnSet{1})
		s.Require().False(equal)
		equal = ColumnSet{2}.Equal(ColumnSet{1, 3})
		s.Require().False(equal)
		equal = ColumnSet{1, 3}.Equal(ColumnSet{2})
		s.Require().False(equal)
	}
	{
		// Contains.
		var contains bool
		// True cases.
		contains = ColumnSet{}.Contains(ColumnSet{})
		s.Require().True(contains)
		contains = ColumnSet{1, 2}.Contains(ColumnSet{})
		s.Require().True(contains)
		contains = ColumnSet{1, 2}.Contains(ColumnSet{2})
		s.Require().True(contains)
		contains = ColumnSet{1, 2}.Contains(ColumnSet{1})
		s.Require().True(contains)
		contains = ColumnSet{1, 2, 3}.Contains(ColumnSet{1, 2})
		s.Require().True(contains)
		// False cases.
		contains = ColumnSet{1}.Contains(ColumnSet{2})
		s.Require().False(contains)
		contains = ColumnSet{2}.Contains(ColumnSet{1, 2})
		s.Require().False(contains)
		contains = ColumnSet{1}.Contains(ColumnSet{1, 2})
		s.Require().False(contains)
		contains = ColumnSet{1, 3, 5}.Contains(ColumnSet{4})
		s.Require().False(contains)
	}
	{
		// IsDisjoint.
		var disjoin bool
		// True cases.
		disjoin = ColumnSet{}.IsDisjoint(ColumnSet{})
		s.Require().True(disjoin)
		disjoin = ColumnSet{}.IsDisjoint(ColumnSet{1})
		s.Require().True(disjoin)
		disjoin = ColumnSet{1}.IsDisjoint(ColumnSet{})
		s.Require().True(disjoin)
		disjoin = ColumnSet{1}.IsDisjoint(ColumnSet{2})
		s.Require().True(disjoin)
		// False cases.
		disjoin = ColumnSet{1, 2}.IsDisjoint(ColumnSet{2})
		s.Require().False(disjoin)
		disjoin = ColumnSet{1, 2}.IsDisjoint(ColumnSet{1})
		s.Require().False(disjoin)
		disjoin = ColumnSet{1, 2}.IsDisjoint(ColumnSet{0, 2})
		s.Require().False(disjoin)
		disjoin = ColumnSet{1, 7}.IsDisjoint(ColumnSet{5, 6, 7})
		s.Require().False(disjoin)
	}
}

func TestPlannerUtils(t *testing.T) {
	suite.Run(t, new(PlannerUtilsTestSuite))
}
