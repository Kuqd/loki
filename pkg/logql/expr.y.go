// Code generated by goyacc -p expr -o pkg/logql/expr.y.go pkg/logql/expr.y. DO NOT EDIT.

//line pkg/logql/expr.y:2
package logql

import __yyfmt__ "fmt"

//line pkg/logql/expr.y:2

import (
	"github.com/prometheus/prometheus/pkg/labels"
	"time"
)

//line pkg/logql/expr.y:10
type exprSymType struct {
	yys                   int
	Expr                  Expr
	Filter                labels.MatchType
	Grouping              *grouping
	Labels                []string
	LogExpr               LogSelectorExpr
	LogRangeExpr          *logRange
	Matcher               *labels.Matcher
	Matchers              []*labels.Matcher
	RangeAggregationExpr  SampleExpr
	RangeOp               string
	Selector              []*labels.Matcher
	VectorAggregationExpr SampleExpr
	MetricExpr            SampleExpr
	VectorOp              string
	BinOpExpr             SampleExpr
	binOp                 string
	str                   string
	duration              time.Duration
	LiteralExpr           *literalExpr
}

const IDENTIFIER = 57346
const STRING = 57347
const NUMBER = 57348
const DURATION = 57349
const MATCHERS = 57350
const LABELS = 57351
const EQ = 57352
const NEQ = 57353
const RE = 57354
const NRE = 57355
const OPEN_BRACE = 57356
const CLOSE_BRACE = 57357
const OPEN_BRACKET = 57358
const CLOSE_BRACKET = 57359
const COMMA = 57360
const DOT = 57361
const PIPE_MATCH = 57362
const PIPE_EXACT = 57363
const OPEN_PARENTHESIS = 57364
const CLOSE_PARENTHESIS = 57365
const BY = 57366
const WITHOUT = 57367
const COUNT_OVER_TIME = 57368
const RATE = 57369
const SUM = 57370
const AVG = 57371
const MAX = 57372
const MIN = 57373
const COUNT = 57374
const STDDEV = 57375
const STDVAR = 57376
const BOTTOMK = 57377
const TOPK = 57378
const OR = 57379
const AND = 57380
const UNLESS = 57381
const ADD = 57382
const SUB = 57383
const MUL = 57384
const DIV = 57385
const MOD = 57386
const POW = 57387

var exprToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENTIFIER",
	"STRING",
	"NUMBER",
	"DURATION",
	"MATCHERS",
	"LABELS",
	"EQ",
	"NEQ",
	"RE",
	"NRE",
	"OPEN_BRACE",
	"CLOSE_BRACE",
	"OPEN_BRACKET",
	"CLOSE_BRACKET",
	"COMMA",
	"DOT",
	"PIPE_MATCH",
	"PIPE_EXACT",
	"OPEN_PARENTHESIS",
	"CLOSE_PARENTHESIS",
	"BY",
	"WITHOUT",
	"COUNT_OVER_TIME",
	"RATE",
	"SUM",
	"AVG",
	"MAX",
	"MIN",
	"COUNT",
	"STDDEV",
	"STDVAR",
	"BOTTOMK",
	"TOPK",
	"OR",
	"AND",
	"UNLESS",
	"ADD",
	"SUB",
	"MUL",
	"DIV",
	"MOD",
	"POW",
}
var exprStatenames = [...]string{}

const exprEofCode = 1
const exprErrCode = 2
const exprInitialStackSize = 16

//line pkg/logql/expr.y:182

//line yacctab:1
var exprExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 3,
	1, 2,
	23, 2,
	37, 2,
	38, 2,
	39, 2,
	40, 2,
	41, 2,
	42, 2,
	43, 2,
	44, 2,
	45, 2,
	-2, 0,
}

const exprPrivate = 57344

const exprLast = 200

var exprAct = [...]int{

	50, 4, 37, 98, 3, 46, 75, 14, 36, 107,
	49, 43, 51, 52, 109, 11, 31, 32, 33, 34,
	35, 36, 110, 6, 51, 52, 106, 17, 18, 19,
	20, 22, 23, 21, 24, 25, 26, 27, 79, 107,
	95, 15, 16, 96, 108, 83, 11, 33, 34, 35,
	36, 78, 82, 76, 6, 81, 48, 70, 17, 18,
	19, 20, 22, 23, 21, 24, 25, 26, 27, 11,
	54, 11, 15, 16, 53, 84, 88, 77, 89, 6,
	87, 86, 93, 97, 94, 80, 2, 100, 29, 30,
	31, 32, 33, 34, 35, 36, 104, 89, 105, 28,
	29, 30, 31, 32, 33, 34, 35, 36, 85, 102,
	111, 112, 101, 68, 55, 56, 57, 58, 59, 60,
	61, 62, 63, 38, 99, 10, 67, 90, 92, 69,
	65, 47, 42, 64, 41, 45, 42, 47, 41, 9,
	90, 39, 40, 13, 66, 39, 40, 8, 103, 42,
	5, 41, 12, 38, 71, 72, 73, 74, 39, 40,
	7, 91, 42, 44, 41, 1, 0, 0, 38, 0,
	0, 39, 40, 92, 66, 0, 0, 42, 0, 41,
	38, 0, 0, 0, 0, 0, 39, 40, 0, 42,
	0, 41, 0, 0, 0, 0, 0, 0, 39, 40,
}
var exprPact = [...]int{

	1, -1000, 62, 178, -1000, -1000, 57, -1000, -1000, -1000,
	-1000, 133, 34, -12, -1000, 68, 64, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 128, -1000, -1000,
	-1000, -1000, -1000, 151, 111, 42, -1000, 144, 55, 32,
	33, 30, 23, -1000, -1000, 50, -24, -24, 5, 5,
	-37, -37, -37, -37, -1000, -1000, -1000, -1000, -1000, 127,
	-1000, 103, 76, 75, 71, 138, 166, 55, 17, 25,
	62, 1, 120, 120, -1000, -1000, -1000, -1000, -1000, 107,
	-1000, -1000, -1000, 121, 125, 0, 1, 3, 21, -1000,
	-9, -1000, -1000, -1000, -1000, -1, -1000, 106, -1000, -1000,
	0, -1000, -1000,
}
var exprPgo = [...]int{

	0, 165, 85, 2, 0, 3, 4, 1, 6, 5,
	163, 160, 152, 150, 147, 143, 139, 125,
}
var exprR1 = [...]int{

	0, 1, 2, 2, 7, 7, 7, 7, 6, 6,
	6, 6, 6, 8, 8, 8, 8, 8, 11, 14,
	14, 14, 14, 14, 3, 3, 3, 3, 13, 13,
	13, 10, 10, 9, 9, 9, 9, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 17, 17, 17, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 12, 12,
	5, 5, 4, 4,
}
var exprR2 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	3, 3, 2, 2, 3, 3, 3, 2, 4, 4,
	5, 5, 6, 7, 1, 1, 1, 1, 3, 3,
	3, 1, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 1, 2, 2, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 3, 4, 4,
}
var exprChk = [...]int{

	-1000, -1, -2, -6, -7, -13, 22, -11, -14, -16,
	-17, 14, -12, -15, 6, 40, 41, 26, 27, 28,
	29, 32, 30, 31, 33, 34, 35, 36, 37, 38,
	39, 40, 41, 42, 43, 44, 45, -3, 2, 20,
	21, 13, 11, -6, -10, 2, -9, 4, 22, 22,
	-4, 24, 25, 6, 6, -2, -2, -2, -2, -2,
	-2, -2, -2, -2, 5, 2, 23, 15, 2, 18,
	15, 10, 11, 12, 13, -8, -6, 22, -7, 6,
	-2, 22, 22, 22, -9, 5, 5, 5, 5, -3,
	2, 23, 7, -6, -8, 23, 18, -7, -5, 4,
	-5, 5, 2, 23, -4, -7, 23, 18, 23, 23,
	23, 4, -4,
}
var exprDef = [...]int{

	0, -2, 1, -2, 3, 8, 0, 4, 5, 6,
	7, 0, 0, 0, 46, 0, 0, 58, 59, 49,
	50, 51, 52, 53, 54, 55, 56, 57, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 12, 24,
	25, 26, 27, 0, 0, 0, 31, 0, 0, 0,
	0, 0, 0, 47, 48, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 9, 11, 10, 28, 29, 0,
	30, 0, 0, 0, 0, 0, 0, 0, 3, 46,
	0, 0, 0, 0, 32, 33, 34, 35, 36, 0,
	17, 18, 13, 0, 0, 19, 0, 3, 0, 60,
	0, 14, 16, 15, 21, 3, 20, 0, 62, 63,
	22, 61, 23,
}
var exprTok1 = [...]int{

	1,
}
var exprTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45,
}
var exprTok3 = [...]int{
	0,
}

var exprErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	exprDebug        = 0
	exprErrorVerbose = false
)

type exprLexer interface {
	Lex(lval *exprSymType) int
	Error(s string)
}

type exprParser interface {
	Parse(exprLexer) int
	Lookahead() int
}

type exprParserImpl struct {
	lval  exprSymType
	stack [exprInitialStackSize]exprSymType
	char  int
}

func (p *exprParserImpl) Lookahead() int {
	return p.char
}

func exprNewParser() exprParser {
	return &exprParserImpl{}
}

const exprFlag = -1000

func exprTokname(c int) string {
	if c >= 1 && c-1 < len(exprToknames) {
		if exprToknames[c-1] != "" {
			return exprToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func exprStatname(s int) string {
	if s >= 0 && s < len(exprStatenames) {
		if exprStatenames[s] != "" {
			return exprStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func exprErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !exprErrorVerbose {
		return "syntax error"
	}

	for _, e := range exprErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + exprTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := exprPact[state]
	for tok := TOKSTART; tok-1 < len(exprToknames); tok++ {
		if n := base + tok; n >= 0 && n < exprLast && exprChk[exprAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if exprDef[state] == -2 {
		i := 0
		for exprExca[i] != -1 || exprExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; exprExca[i] >= 0; i += 2 {
			tok := exprExca[i]
			if tok < TOKSTART || exprExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if exprExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += exprTokname(tok)
	}
	return res
}

func exprlex1(lex exprLexer, lval *exprSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = exprTok1[0]
		goto out
	}
	if char < len(exprTok1) {
		token = exprTok1[char]
		goto out
	}
	if char >= exprPrivate {
		if char < exprPrivate+len(exprTok2) {
			token = exprTok2[char-exprPrivate]
			goto out
		}
	}
	for i := 0; i < len(exprTok3); i += 2 {
		token = exprTok3[i+0]
		if token == char {
			token = exprTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = exprTok2[1] /* unknown char */
	}
	if exprDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", exprTokname(token), uint(char))
	}
	return char, token
}

func exprParse(exprlex exprLexer) int {
	return exprNewParser().Parse(exprlex)
}

func (exprrcvr *exprParserImpl) Parse(exprlex exprLexer) int {
	var exprn int
	var exprVAL exprSymType
	var exprDollar []exprSymType
	_ = exprDollar // silence set and not used
	exprS := exprrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	exprstate := 0
	exprrcvr.char = -1
	exprtoken := -1 // exprrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		exprstate = -1
		exprrcvr.char = -1
		exprtoken = -1
	}()
	exprp := -1
	goto exprstack

ret0:
	return 0

ret1:
	return 1

exprstack:
	/* put a state and value onto the stack */
	if exprDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", exprTokname(exprtoken), exprStatname(exprstate))
	}

	exprp++
	if exprp >= len(exprS) {
		nyys := make([]exprSymType, len(exprS)*2)
		copy(nyys, exprS)
		exprS = nyys
	}
	exprS[exprp] = exprVAL
	exprS[exprp].yys = exprstate

exprnewstate:
	exprn = exprPact[exprstate]
	if exprn <= exprFlag {
		goto exprdefault /* simple state */
	}
	if exprrcvr.char < 0 {
		exprrcvr.char, exprtoken = exprlex1(exprlex, &exprrcvr.lval)
	}
	exprn += exprtoken
	if exprn < 0 || exprn >= exprLast {
		goto exprdefault
	}
	exprn = exprAct[exprn]
	if exprChk[exprn] == exprtoken { /* valid shift */
		exprrcvr.char = -1
		exprtoken = -1
		exprVAL = exprrcvr.lval
		exprstate = exprn
		if Errflag > 0 {
			Errflag--
		}
		goto exprstack
	}

exprdefault:
	/* default state action */
	exprn = exprDef[exprstate]
	if exprn == -2 {
		if exprrcvr.char < 0 {
			exprrcvr.char, exprtoken = exprlex1(exprlex, &exprrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if exprExca[xi+0] == -1 && exprExca[xi+1] == exprstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			exprn = exprExca[xi+0]
			if exprn < 0 || exprn == exprtoken {
				break
			}
		}
		exprn = exprExca[xi+1]
		if exprn < 0 {
			goto ret0
		}
	}
	if exprn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			exprlex.Error(exprErrorMessage(exprstate, exprtoken))
			Nerrs++
			if exprDebug >= 1 {
				__yyfmt__.Printf("%s", exprStatname(exprstate))
				__yyfmt__.Printf(" saw %s\n", exprTokname(exprtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for exprp >= 0 {
				exprn = exprPact[exprS[exprp].yys] + exprErrCode
				if exprn >= 0 && exprn < exprLast {
					exprstate = exprAct[exprn] /* simulate a shift of "error" */
					if exprChk[exprstate] == exprErrCode {
						goto exprstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if exprDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", exprS[exprp].yys)
				}
				exprp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if exprDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", exprTokname(exprtoken))
			}
			if exprtoken == exprEofCode {
				goto ret1
			}
			exprrcvr.char = -1
			exprtoken = -1
			goto exprnewstate /* try again in the same state */
		}
	}

	/* reduction by production exprn */
	if exprDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", exprn, exprStatname(exprstate))
	}

	exprnt := exprn
	exprpt := exprp
	_ = exprpt // guard against "declared and not used"

	exprp -= exprR2[exprn]
	// exprp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if exprp+1 >= len(exprS) {
		nyys := make([]exprSymType, len(exprS)*2)
		copy(nyys, exprS)
		exprS = nyys
	}
	exprVAL = exprS[exprp+1]

	/* consult goto table to find next state */
	exprn = exprR1[exprn]
	exprg := exprPgo[exprn]
	exprj := exprg + exprS[exprp].yys + 1

	if exprj >= exprLast {
		exprstate = exprAct[exprg]
	} else {
		exprstate = exprAct[exprj]
		if exprChk[exprstate] != -exprn {
			exprstate = exprAct[exprg]
		}
	}
	// dummy call; replaced with literal code
	switch exprnt {

	case 1:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:65
		{
			exprlex.(*lexer).expr = exprDollar[1].Expr
		}
	case 2:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:68
		{
			exprVAL.Expr = exprDollar[1].LogExpr
		}
	case 3:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:69
		{
			exprVAL.Expr = exprDollar[1].MetricExpr
		}
	case 4:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:73
		{
			exprVAL.MetricExpr = exprDollar[1].RangeAggregationExpr
		}
	case 5:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:74
		{
			exprVAL.MetricExpr = exprDollar[1].VectorAggregationExpr
		}
	case 6:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:75
		{
			exprVAL.MetricExpr = exprDollar[1].BinOpExpr
		}
	case 7:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:76
		{
			exprVAL.MetricExpr = exprDollar[1].LiteralExpr
		}
	case 8:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:80
		{
			exprVAL.LogExpr = newMatcherExpr(exprDollar[1].Selector)
		}
	case 9:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:81
		{
			exprVAL.LogExpr = NewFilterExpr(exprDollar[1].LogExpr, exprDollar[2].Filter, exprDollar[3].str)
		}
	case 10:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:82
		{
			exprVAL.LogExpr = exprDollar[2].LogExpr
		}
	case 13:
		exprDollar = exprS[exprpt-2 : exprpt+1]
//line pkg/logql/expr.y:88
		{
			exprVAL.LogRangeExpr = newLogRange(exprDollar[1].LogExpr, exprDollar[2].duration)
		}
	case 14:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:89
		{
			exprVAL.LogRangeExpr = addFilterToLogRangeExpr(exprDollar[1].LogRangeExpr, exprDollar[2].Filter, exprDollar[3].str)
		}
	case 15:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:90
		{
			exprVAL.LogRangeExpr = exprDollar[2].LogRangeExpr
		}
	case 18:
		exprDollar = exprS[exprpt-4 : exprpt+1]
//line pkg/logql/expr.y:95
		{
			exprVAL.RangeAggregationExpr = newRangeAggregationExpr(exprDollar[3].LogRangeExpr, exprDollar[1].RangeOp)
		}
	case 19:
		exprDollar = exprS[exprpt-4 : exprpt+1]
//line pkg/logql/expr.y:99
		{
			exprVAL.VectorAggregationExpr = mustNewVectorAggregationExpr(exprDollar[3].MetricExpr, exprDollar[1].VectorOp, nil, nil)
		}
	case 20:
		exprDollar = exprS[exprpt-5 : exprpt+1]
//line pkg/logql/expr.y:100
		{
			exprVAL.VectorAggregationExpr = mustNewVectorAggregationExpr(exprDollar[4].MetricExpr, exprDollar[1].VectorOp, exprDollar[2].Grouping, nil)
		}
	case 21:
		exprDollar = exprS[exprpt-5 : exprpt+1]
//line pkg/logql/expr.y:101
		{
			exprVAL.VectorAggregationExpr = mustNewVectorAggregationExpr(exprDollar[3].MetricExpr, exprDollar[1].VectorOp, exprDollar[5].Grouping, nil)
		}
	case 22:
		exprDollar = exprS[exprpt-6 : exprpt+1]
//line pkg/logql/expr.y:103
		{
			exprVAL.VectorAggregationExpr = mustNewVectorAggregationExpr(exprDollar[5].MetricExpr, exprDollar[1].VectorOp, nil, &exprDollar[3].str)
		}
	case 23:
		exprDollar = exprS[exprpt-7 : exprpt+1]
//line pkg/logql/expr.y:104
		{
			exprVAL.VectorAggregationExpr = mustNewVectorAggregationExpr(exprDollar[5].MetricExpr, exprDollar[1].VectorOp, exprDollar[7].Grouping, &exprDollar[3].str)
		}
	case 24:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:108
		{
			exprVAL.Filter = labels.MatchRegexp
		}
	case 25:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:109
		{
			exprVAL.Filter = labels.MatchEqual
		}
	case 26:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:110
		{
			exprVAL.Filter = labels.MatchNotRegexp
		}
	case 27:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:111
		{
			exprVAL.Filter = labels.MatchNotEqual
		}
	case 28:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:115
		{
			exprVAL.Selector = exprDollar[2].Matchers
		}
	case 29:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:116
		{
			exprVAL.Selector = exprDollar[2].Matchers
		}
	case 30:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:117
		{
		}
	case 31:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:121
		{
			exprVAL.Matchers = []*labels.Matcher{exprDollar[1].Matcher}
		}
	case 32:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:122
		{
			exprVAL.Matchers = append(exprDollar[1].Matchers, exprDollar[3].Matcher)
		}
	case 33:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:126
		{
			exprVAL.Matcher = mustNewMatcher(labels.MatchEqual, exprDollar[1].str, exprDollar[3].str)
		}
	case 34:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:127
		{
			exprVAL.Matcher = mustNewMatcher(labels.MatchNotEqual, exprDollar[1].str, exprDollar[3].str)
		}
	case 35:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:128
		{
			exprVAL.Matcher = mustNewMatcher(labels.MatchRegexp, exprDollar[1].str, exprDollar[3].str)
		}
	case 36:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:129
		{
			exprVAL.Matcher = mustNewMatcher(labels.MatchNotRegexp, exprDollar[1].str, exprDollar[3].str)
		}
	case 37:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:138
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("or", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 38:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:139
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("and", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 39:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:140
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("unless", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 40:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:141
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("+", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 41:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:142
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("-", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 42:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:143
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("*", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 43:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:144
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("/", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 44:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:145
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("%", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 45:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:146
		{
			exprVAL.BinOpExpr = mustNewBinOpExpr("^", exprDollar[1].Expr, exprDollar[3].Expr)
		}
	case 46:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:150
		{
			exprVAL.LiteralExpr = mustNewLiteralExpr(exprDollar[1].str, false)
		}
	case 47:
		exprDollar = exprS[exprpt-2 : exprpt+1]
//line pkg/logql/expr.y:151
		{
			exprVAL.LiteralExpr = mustNewLiteralExpr(exprDollar[2].str, false)
		}
	case 48:
		exprDollar = exprS[exprpt-2 : exprpt+1]
//line pkg/logql/expr.y:152
		{
			exprVAL.LiteralExpr = mustNewLiteralExpr(exprDollar[2].str, true)
		}
	case 49:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:156
		{
			exprVAL.VectorOp = OpTypeSum
		}
	case 50:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:157
		{
			exprVAL.VectorOp = OpTypeAvg
		}
	case 51:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:158
		{
			exprVAL.VectorOp = OpTypeCount
		}
	case 52:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:159
		{
			exprVAL.VectorOp = OpTypeMax
		}
	case 53:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:160
		{
			exprVAL.VectorOp = OpTypeMin
		}
	case 54:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:161
		{
			exprVAL.VectorOp = OpTypeStddev
		}
	case 55:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:162
		{
			exprVAL.VectorOp = OpTypeStdvar
		}
	case 56:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:163
		{
			exprVAL.VectorOp = OpTypeBottomK
		}
	case 57:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:164
		{
			exprVAL.VectorOp = OpTypeTopK
		}
	case 58:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:168
		{
			exprVAL.RangeOp = OpTypeCountOverTime
		}
	case 59:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:169
		{
			exprVAL.RangeOp = OpTypeRate
		}
	case 60:
		exprDollar = exprS[exprpt-1 : exprpt+1]
//line pkg/logql/expr.y:174
		{
			exprVAL.Labels = []string{exprDollar[1].str}
		}
	case 61:
		exprDollar = exprS[exprpt-3 : exprpt+1]
//line pkg/logql/expr.y:175
		{
			exprVAL.Labels = append(exprDollar[1].Labels, exprDollar[3].str)
		}
	case 62:
		exprDollar = exprS[exprpt-4 : exprpt+1]
//line pkg/logql/expr.y:179
		{
			exprVAL.Grouping = &grouping{without: false, groups: exprDollar[3].Labels}
		}
	case 63:
		exprDollar = exprS[exprpt-4 : exprpt+1]
//line pkg/logql/expr.y:180
		{
			exprVAL.Grouping = &grouping{without: true, groups: exprDollar[3].Labels}
		}
	}
	goto exprstack /* stack new state and value */
}
