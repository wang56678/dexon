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
	s.requireParseNoError(`select .0`)
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
	s.requireParseNoError(`select * from abc where a like a + 1 escape '*'`)
	s.requireParseNoError(`select * from abc where a like a + 1 escape a`)

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
	s.requireParseNoError(`create table "~!@#$%^&*()" ( a int32 references b ( a ) , b bytes primary key, c address not null default 1 + 1 )`)

	// Test create index.
	s.requireParseNoError(`create unique index a on a (a)`)
	s.requireParseNoError(`create index "~!@#$%^&*()" on ㄅ ( a , b )`)
	s.requireParseNoError(`create index ㄅㄆㄇ on 👍 ( 🌍 , 💯 )`)
}

func (s *ParserTestSuite) TestParseRules() {
	s.requireParseNoError(`
		SELECT
			C1,
			*,
			SUM(*),
			COUNT(*) + 1,
			*,
			NOT A >= B,
			NULL IS NULL,
			C2 OR C3 AND TRUE OR FALSE,
			C4 NOT IN (C5, 849, 2899 - C6),
			C7 + C8 IN (C9, 5566, 9487 * C10),
			C10 IS NULL,
			C11 IS NOT NULL,
			C12 LIKE 'dek_s%n',
			C13 || C14 NOT LIKE 'cob%h__d%',
			C15 LIKE 'dek_s\\%n' ESCAPE '\\',
			C16 <= C17 + 45,
			C18 >= C19 - 54,
			C20 <> 46 * C21,
			C22 != 64 / C22,
			C23 < C24 % C25,
			C26 > C27 / (C28 + C29),
			C30 = C31 * (C32 - C33),
			C34 || C35 || 'vm' || 'sql',
			C36 + C37 - C38 * C39 / C40 % C41,
			C42 - - C43 + + (C44) * -C45 ++ C46 / -C47,
			C48 + CAST(C49 % C50 AS INT88) - TSAC(),
			F(C51) * "F"(C52, "C53") + "!"('\U0010FFFF', '\x11\x23\xfd'),
			0x845 - 0x6ea - 0xbf,
			00244 - 1.56 + 24. - .34,
			1.2e1 - 2.04e-5 + -4.53e+10,
			1e1 + 1.e1 - .1e1,
			-1e1 + -1.e1 - -.1e1,
			0.0 + 0e0 - 0.e0 + .0e0 * 0.,
			-0.0 + -0e0 - -0.e0 + -.0e0 * -0.,
			'normal' || x'8e7a' || hex'abcdef' || C54
			FROM T
			WHERE W
			GROUP BY
				1,
				C1,
				C2 + 2,
				C3 - C4,
				C5 AND C6
			ORDER BY
				1,
				2 ASC,
				C1 DESC,
				C2 NULLS FIRST,
				C3 + C4 NULLS LAST,
				C5 * (C6 - C7) ASC NULLS FIRST,
				C8 || C9 || 'dexon' DESC NULLS LAST
			LIMIT 218 OFFSET 2019;

		UPDATE T
			SET
				C1 = C1 = C2 OR C3 <> C4,
				C2 = C2 IS NOT NULL,
				C3 = DEFAULT
			WHERE W;

		DELETE FROM T WHERE W;

		INSERT INTO T DEFAULT VALUES;
		INSERT INTO T VALUES (V1, V2, V3, V4, V5);
		INSERT INTO T VALUES (DEFAULT, DEFAULT, DEFAULT, DEFAULT, DEFAULT);
		INSERT INTO T (C1) VALUES (V1), (DEFAULT);
		INSERT INTO T (C1, C2, C3)
			VALUES (V1, V2, V3 + V4), (V5 IS NULL, DEFAULT, NULL);

		CREATE TABLE T (
			C1 UINT64 PRIMARY KEY AUTOINCREMENT,
			C2 ADDRESS REFERENCES U (D) NOT NULL,
			C3 UINT256 DEFAULT 3 * 2 + 1,
			C4 BYTES5 DEFAULT 'hello',
			C5 INT24 UNIQUE NOT NULL,
			C6 BYTES
		);

		CREATE TABLE T (
			C1 INT224,
			C2 UINT168,
			C3 FIXED72X0,
			C4 UFIXED80X80,
			C5 BYTES32,
			C6 BYTES1,
			C7 BYTE,
			C8 BYTES,
			C9 ADDRESS,
			C10 BOOL
		);

		CREATE INDEX I ON T (C1);
		CREATE UNIQUE INDEX I ON T (C2, C3);
	`)
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
