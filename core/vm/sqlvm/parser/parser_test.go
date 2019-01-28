package parser

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/text/encoding/traditionalchinese"
)

type ParserTestSuite struct{ suite.Suite }

func (s *ParserTestSuite) requireParseNoError(sql string) {
	_, err := Parse([]byte(sql))
	s.Require().NoError(err)
}

func (s *ParserTestSuite) TestParse() {
	// Test stmt.
	s.requireParseNoError(``)
	s.requireParseNoError(`;`)
	s.requireParseNoError(`;;;select 1;;;;`)

	// Test expr.
	s.requireParseNoError(`select 1 + 2 * 3`)
	s.requireParseNoError(`select a(1 + 1)`)
	s.requireParseNoError(`select hEx'12'`)
	s.requireParseNoError(`select x'12'`)
	s.requireParseNoError(`select 0xABC`)
	s.requireParseNoError(`select true and false or true and false or true`)
	s.requireParseNoError(`SeLeCT '1' NoT LiKe '1';`)
	s.requireParseNoError(`select a in (1,2) is not null not in (true)`)
	s.requireParseNoError(`select count(*)`)
	s.requireParseNoError(`select cast(a as fixed65535X1)`)
	s.requireParseNoError(`select "now!" ( not a + b, aa( + 3 + .1 + 1. ) + - .3e-9  + 1.e-10, 'a' || 'b' and true )`)

	// Test where.
	s.requireParseNoError(`select * from abc where abc is null`)
	s.requireParseNoError(`select * from abc where abc is not null`)
	s.requireParseNoError(`select * from abc where abc in (1, 1)`)
	s.requireParseNoError(`select * from abc where abc not in (1, 1)`)
	s.requireParseNoError(`select * from abc where not true`)
	s.requireParseNoError(`select * from abc where a like a + 1`)

	// Test some logic expr and no from.
	s.requireParseNoError(`select 1 where a is not null = b`)
	s.requireParseNoError(`select 1 where null = null is null and true`)
	s.requireParseNoError(`select 1 where null is null = null`)
	s.requireParseNoError(`SELECT 1 + 2 WHERE 3 <> 4`)

	// Test order by.
	s.requireParseNoError(`select a=b+1 from a order by a desc`)
	s.requireParseNoError(`select 1 from a order by b + 1 desc`)
	s.requireParseNoError(`select 1 from a order by b + 1 nulls first`)
	s.requireParseNoError(`select 1 from a order by b + 1 desc nulls last`)

	// Test group by.
	s.requireParseNoError(`select 1 from a group by b + 1`)

	// Test insert.
	s.requireParseNoError(`insert into "abc"(a) values (f(a, b),g(a),h())`)
	s.requireParseNoError(`insert into "abc"(a) values (1,2,3), (f(a, b),g(a),h())`)
	s.requireParseNoError(`insert into a default values`)
	s.requireParseNoError(`insert into a values (default)`)

	// Test update.
	s.requireParseNoError(`update "~!@#$%^&*()" set b = default where a is null;`)
	s.requireParseNoError(`update "~!@#$%^&*()" set b = default, a = 123 where a is null;`)

	// Test delete.
	s.requireParseNoError(`delete from a where b is null`)

	// Test create table.
	s.requireParseNoError(`create table a (a int32 not null unique primary key default 0)`)
	s.requireParseNoError(`create table "~!@#$%^&*()" ( a int32 references b ( a ) , b string primary key, c address not null default 1 + 1 )`)

	// Test create index.
	s.requireParseNoError(`create unique index a on a (a)`)
	s.requireParseNoError(`create index "~!@#$%^&*()" on ㄅ ( a , b )`)
	s.requireParseNoError(`create index ㄅㄆㄇ on 👍 ( 🌍 , 💯 )`)
}

func (s *ParserTestSuite) TestParseInvalidUTF8() {
	query := `SELECT ㄅ FROM 東 WHERE — - ─ = ██`
	query, err := traditionalchinese.Big5.NewEncoder().String(query)
	s.Require().NoError(err)
	s.requireParseNoError(query)
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
