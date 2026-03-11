// Code generated from PPLParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package generated // PPLParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type PPLParser struct {
	*antlr.BaseParser
}

var PPLParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func pplparserParserInit() {
	staticData := &PPLParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "'where'", "'fields'", "'stats'", "'sort'", "'head'", "'describe'",
		"'showdatasources'", "'explain'", "'chart'", "'timechart'", "'bin'",
		"'dedup'", "'top'", "'rare'", "'eval'", "'rename'", "'replace'", "'fillnull'",
		"'parse'", "'rex'", "'lookup'", "'append'", "'join'", "'table'", "'eventstats'",
		"'streamstats'", "'reverse'", "'flatten'", "'by'", "'as'", "'with'",
		"'output'", "'type'", "'value'", "'inner'", "'left'", "'right'", "'outer'",
		"'full'", "'span'", "'bins'", "'keepevents'", "'consecutive'", "'sortby'",
		"'limit'", "'countfield'", "'percentfield'", "'showcount'", "'showperc'",
		"'useother'", "'otherstr'", "'nullstr'", "'asc'", "'desc'", "", "",
		"", "'='", "", "'<'", "'<='", "'>'", "'>='", "", "", "'+'", "'-'", "'*'",
		"'/'", "'%'", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "'|'", "','", "'.'",
		"'('", "')'", "'['", "']'",
	}
	staticData.SymbolicNames = []string{
		"", "SEARCH", "WHERE", "FIELDS", "STATS", "SORT", "HEAD", "DESCRIBE",
		"SHOWDATASOURCES", "EXPLAIN", "CHART", "TIMECHART", "BIN", "DEDUP",
		"TOP", "RARE", "EVAL", "RENAME", "REPLACE", "FILLNULL", "PARSE", "REX",
		"LOOKUP", "APPEND", "JOIN", "TABLE", "EVENTSTATS", "STREAMSTATS", "REVERSE",
		"FLATTEN", "BY", "AS", "WITH", "OUTPUT", "TYPE", "VALUE", "INNER", "LEFT",
		"RIGHT", "OUTER", "FULL", "SPAN", "BINS", "KEEPEVENTS", "CONSECUTIVE",
		"SORTBY", "LIMIT", "COUNTFIELD", "PERCENTFIELD", "SHOWCOUNT", "SHOWPERC",
		"USEOTHER", "OTHERSTR", "NULLSTR", "ASC", "DESC", "AND", "OR", "NOT",
		"EQ", "NEQ", "LT", "LTE", "GT", "GTE", "LIKE", "IN", "PLUS", "MINUS",
		"STAR", "SLASH", "PERCENT", "TRUE", "FALSE", "NULL", "CASE", "WHEN",
		"THEN", "ELSE", "END", "COUNT", "SUM", "AVG", "MIN", "MAX", "DC", "DISTINCT_COUNT",
		"VAR", "VARP", "STDEV", "STDEVP", "PERCENTILE", "MEDIAN", "MODE", "EARLIEST",
		"LATEST", "VALUES", "RANGE", "DISTINCT", "PIPE", "COMMA", "DOT", "LPAREN",
		"RPAREN", "LBRACKET", "RBRACKET", "IDENTIFIER", "INTEGER", "DECIMAL",
		"STRING", "WS", "LINE_COMMENT", "BLOCK_COMMENT",
	}
	staticData.RuleNames = []string{
		"query", "searchQuery", "processingCommand", "metadataCommand", "searchCommand",
		"whereCommand", "fieldsCommand", "fieldList", "statsCommand", "aggregationList",
		"aggregation", "sortCommand", "sortFieldList", "sortField", "headCommand",
		"chartCommand", "chartOptions", "timechartCommand", "timechartOptions",
		"timeSpan", "binCommand", "binOptions", "dedupCommand", "dedupOptions",
		"topCommand", "topOptions", "rareCommand", "evalCommand", "evalAssignment",
		"renameCommand", "renameAssignment", "replaceCommand", "replaceMapping",
		"fillnullCommand", "fillnullAssignment", "parseCommand", "rexCommand",
		"lookupCommand", "lookupOutputList", "lookupOutputField", "appendCommand",
		"joinCommand", "joinType", "tableCommand", "eventstatsCommand", "streamstatsCommand",
		"streamstatsOptions", "reverseCommand", "flattenCommand", "booleanValue",
		"describeCommand", "showDatasourcesCommand", "explainCommand", "expression",
		"orExpression", "andExpression", "notExpression", "comparisonExpression",
		"additiveExpression", "multiplicativeExpression", "unaryExpression",
		"primaryExpression", "literal", "fieldReference", "functionCall", "aggregationFunction",
		"expressionList", "caseExpression", "whenClause",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 112, 704, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57, 7, 57,
		2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7, 62, 2,
		63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67, 2, 68,
		7, 68, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 3, 0,
		149, 8, 0, 1, 1, 1, 1, 1, 1, 5, 1, 154, 8, 1, 10, 1, 12, 1, 157, 9, 1,
		1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 2, 3, 2, 185, 8, 2, 1, 3, 1, 3, 3, 3, 189, 8, 3, 1, 4, 1, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 198, 8, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1,
		6, 1, 6, 1, 6, 1, 6, 3, 6, 208, 8, 6, 1, 7, 1, 7, 1, 7, 5, 7, 213, 8, 7,
		10, 7, 12, 7, 216, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 222, 8, 8, 1, 9,
		1, 9, 1, 9, 5, 9, 227, 8, 9, 10, 9, 12, 9, 230, 9, 9, 1, 10, 1, 10, 1,
		10, 3, 10, 235, 8, 10, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 5, 12,
		243, 8, 12, 10, 12, 12, 12, 246, 9, 12, 1, 13, 1, 13, 3, 13, 250, 8, 13,
		1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 259, 8, 15, 1,
		15, 5, 15, 262, 8, 15, 10, 15, 12, 15, 265, 9, 15, 1, 16, 1, 16, 1, 16,
		1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1,
		16, 1, 16, 3, 16, 282, 8, 16, 1, 17, 1, 17, 5, 17, 286, 8, 17, 10, 17,
		12, 17, 289, 9, 17, 1, 17, 1, 17, 1, 17, 3, 17, 294, 8, 17, 1, 18, 1, 18,
		1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 3,
		18, 308, 8, 18, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 5, 20, 316, 8,
		20, 10, 20, 12, 20, 319, 9, 20, 1, 21, 1, 21, 1, 21, 1, 21, 3, 21, 325,
		8, 21, 1, 21, 1, 21, 1, 21, 3, 21, 330, 8, 21, 1, 22, 1, 22, 3, 22, 334,
		8, 22, 1, 22, 1, 22, 5, 22, 338, 8, 22, 10, 22, 12, 22, 341, 9, 22, 1,
		23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 3, 23, 351, 8, 23,
		1, 24, 1, 24, 3, 24, 355, 8, 24, 1, 24, 1, 24, 1, 24, 3, 24, 360, 8, 24,
		1, 24, 5, 24, 363, 8, 24, 10, 24, 12, 24, 366, 9, 24, 1, 25, 1, 25, 1,
		25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25,
		1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 389, 8,
		25, 1, 26, 1, 26, 3, 26, 393, 8, 26, 1, 26, 1, 26, 1, 26, 3, 26, 398, 8,
		26, 1, 26, 5, 26, 401, 8, 26, 10, 26, 12, 26, 404, 9, 26, 1, 27, 1, 27,
		1, 27, 1, 27, 5, 27, 410, 8, 27, 10, 27, 12, 27, 413, 9, 27, 1, 28, 1,
		28, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 1, 29, 5, 29, 423, 8, 29, 10, 29,
		12, 29, 426, 9, 29, 1, 30, 1, 30, 1, 30, 1, 30, 1, 31, 1, 31, 1, 31, 1,
		31, 5, 31, 436, 8, 31, 10, 31, 12, 31, 439, 9, 31, 1, 31, 1, 31, 1, 31,
		1, 32, 1, 32, 1, 32, 1, 32, 1, 33, 1, 33, 1, 33, 1, 33, 1, 33, 1, 33, 3,
		33, 454, 8, 33, 1, 33, 1, 33, 1, 33, 1, 33, 5, 33, 460, 8, 33, 10, 33,
		12, 33, 463, 9, 33, 3, 33, 465, 8, 33, 1, 34, 1, 34, 1, 34, 1, 34, 1, 35,
		1, 35, 1, 35, 3, 35, 474, 8, 35, 1, 35, 1, 35, 1, 35, 1, 36, 1, 36, 1,
		36, 1, 36, 3, 36, 483, 8, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 1, 37,
		1, 37, 3, 37, 492, 8, 37, 1, 37, 1, 37, 1, 37, 1, 38, 1, 38, 1, 38, 5,
		38, 500, 8, 38, 10, 38, 12, 38, 503, 9, 38, 1, 39, 1, 39, 1, 39, 3, 39,
		508, 8, 39, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 41, 1, 41, 1, 41, 1,
		41, 3, 41, 519, 8, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 42, 1, 42,
		1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1, 44, 1, 44, 3, 44, 535, 8, 44, 1,
		45, 1, 45, 5, 45, 539, 8, 45, 10, 45, 12, 45, 542, 9, 45, 1, 45, 1, 45,
		1, 45, 3, 45, 547, 8, 45, 1, 46, 1, 46, 1, 46, 1, 46, 1, 46, 3, 46, 554,
		8, 46, 1, 47, 1, 47, 1, 48, 1, 48, 1, 48, 1, 49, 1, 49, 1, 50, 1, 50, 1,
		50, 1, 51, 1, 51, 1, 52, 1, 52, 1, 53, 1, 53, 1, 54, 1, 54, 1, 54, 5, 54,
		575, 8, 54, 10, 54, 12, 54, 578, 9, 54, 1, 55, 1, 55, 1, 55, 5, 55, 583,
		8, 55, 10, 55, 12, 55, 586, 9, 55, 1, 56, 1, 56, 1, 56, 3, 56, 591, 8,
		56, 1, 57, 1, 57, 1, 57, 1, 57, 1, 57, 1, 57, 1, 57, 1, 57, 1, 57, 1, 57,
		3, 57, 603, 8, 57, 1, 58, 1, 58, 1, 58, 5, 58, 608, 8, 58, 10, 58, 12,
		58, 611, 9, 58, 1, 59, 1, 59, 1, 59, 5, 59, 616, 8, 59, 10, 59, 12, 59,
		619, 9, 59, 1, 60, 1, 60, 1, 60, 3, 60, 624, 8, 60, 1, 61, 1, 61, 1, 61,
		1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 3, 61, 634, 8, 61, 1, 62, 1, 62, 1,
		63, 1, 63, 1, 63, 5, 63, 641, 8, 63, 10, 63, 12, 63, 644, 9, 63, 1, 63,
		1, 63, 1, 63, 1, 63, 3, 63, 650, 8, 63, 1, 64, 1, 64, 1, 64, 1, 64, 1,
		64, 1, 64, 3, 64, 658, 8, 64, 1, 64, 1, 64, 1, 64, 1, 64, 1, 64, 1, 64,
		1, 64, 1, 64, 1, 64, 1, 64, 3, 64, 670, 8, 64, 1, 64, 1, 64, 1, 64, 3,
		64, 675, 8, 64, 1, 65, 1, 65, 1, 66, 1, 66, 1, 66, 5, 66, 682, 8, 66, 10,
		66, 12, 66, 685, 9, 66, 1, 67, 1, 67, 4, 67, 689, 8, 67, 11, 67, 12, 67,
		690, 1, 67, 1, 67, 3, 67, 695, 8, 67, 1, 67, 1, 67, 1, 68, 1, 68, 1, 68,
		1, 68, 1, 68, 1, 68, 0, 0, 69, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22,
		24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58,
		60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94,
		96, 98, 100, 102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124,
		126, 128, 130, 132, 134, 136, 0, 8, 1, 0, 54, 55, 1, 0, 36, 40, 1, 0, 72,
		73, 1, 0, 59, 64, 1, 0, 67, 68, 1, 0, 69, 71, 2, 0, 72, 74, 107, 109, 1,
		0, 80, 97, 741, 0, 148, 1, 0, 0, 0, 2, 150, 1, 0, 0, 0, 4, 184, 1, 0, 0,
		0, 6, 188, 1, 0, 0, 0, 8, 197, 1, 0, 0, 0, 10, 199, 1, 0, 0, 0, 12, 207,
		1, 0, 0, 0, 14, 209, 1, 0, 0, 0, 16, 217, 1, 0, 0, 0, 18, 223, 1, 0, 0,
		0, 20, 231, 1, 0, 0, 0, 22, 236, 1, 0, 0, 0, 24, 239, 1, 0, 0, 0, 26, 247,
		1, 0, 0, 0, 28, 251, 1, 0, 0, 0, 30, 254, 1, 0, 0, 0, 32, 281, 1, 0, 0,
		0, 34, 283, 1, 0, 0, 0, 36, 307, 1, 0, 0, 0, 38, 309, 1, 0, 0, 0, 40, 312,
		1, 0, 0, 0, 42, 329, 1, 0, 0, 0, 44, 331, 1, 0, 0, 0, 46, 350, 1, 0, 0,
		0, 48, 352, 1, 0, 0, 0, 50, 388, 1, 0, 0, 0, 52, 390, 1, 0, 0, 0, 54, 405,
		1, 0, 0, 0, 56, 414, 1, 0, 0, 0, 58, 418, 1, 0, 0, 0, 60, 427, 1, 0, 0,
		0, 62, 431, 1, 0, 0, 0, 64, 443, 1, 0, 0, 0, 66, 464, 1, 0, 0, 0, 68, 466,
		1, 0, 0, 0, 70, 470, 1, 0, 0, 0, 72, 478, 1, 0, 0, 0, 74, 486, 1, 0, 0,
		0, 76, 496, 1, 0, 0, 0, 78, 504, 1, 0, 0, 0, 80, 509, 1, 0, 0, 0, 82, 514,
		1, 0, 0, 0, 84, 525, 1, 0, 0, 0, 86, 527, 1, 0, 0, 0, 88, 530, 1, 0, 0,
		0, 90, 536, 1, 0, 0, 0, 92, 548, 1, 0, 0, 0, 94, 555, 1, 0, 0, 0, 96, 557,
		1, 0, 0, 0, 98, 560, 1, 0, 0, 0, 100, 562, 1, 0, 0, 0, 102, 565, 1, 0,
		0, 0, 104, 567, 1, 0, 0, 0, 106, 569, 1, 0, 0, 0, 108, 571, 1, 0, 0, 0,
		110, 579, 1, 0, 0, 0, 112, 590, 1, 0, 0, 0, 114, 592, 1, 0, 0, 0, 116,
		604, 1, 0, 0, 0, 118, 612, 1, 0, 0, 0, 120, 623, 1, 0, 0, 0, 122, 633,
		1, 0, 0, 0, 124, 635, 1, 0, 0, 0, 126, 649, 1, 0, 0, 0, 128, 674, 1, 0,
		0, 0, 130, 676, 1, 0, 0, 0, 132, 678, 1, 0, 0, 0, 134, 686, 1, 0, 0, 0,
		136, 698, 1, 0, 0, 0, 138, 139, 3, 104, 52, 0, 139, 140, 3, 2, 1, 0, 140,
		141, 5, 0, 0, 1, 141, 149, 1, 0, 0, 0, 142, 143, 3, 2, 1, 0, 143, 144,
		5, 0, 0, 1, 144, 149, 1, 0, 0, 0, 145, 146, 3, 6, 3, 0, 146, 147, 5, 0,
		0, 1, 147, 149, 1, 0, 0, 0, 148, 138, 1, 0, 0, 0, 148, 142, 1, 0, 0, 0,
		148, 145, 1, 0, 0, 0, 149, 1, 1, 0, 0, 0, 150, 155, 3, 8, 4, 0, 151, 152,
		5, 99, 0, 0, 152, 154, 3, 4, 2, 0, 153, 151, 1, 0, 0, 0, 154, 157, 1, 0,
		0, 0, 155, 153, 1, 0, 0, 0, 155, 156, 1, 0, 0, 0, 156, 3, 1, 0, 0, 0, 157,
		155, 1, 0, 0, 0, 158, 185, 3, 10, 5, 0, 159, 185, 3, 12, 6, 0, 160, 185,
		3, 16, 8, 0, 161, 185, 3, 22, 11, 0, 162, 185, 3, 28, 14, 0, 163, 185,
		3, 30, 15, 0, 164, 185, 3, 34, 17, 0, 165, 185, 3, 40, 20, 0, 166, 185,
		3, 44, 22, 0, 167, 185, 3, 48, 24, 0, 168, 185, 3, 52, 26, 0, 169, 185,
		3, 54, 27, 0, 170, 185, 3, 58, 29, 0, 171, 185, 3, 62, 31, 0, 172, 185,
		3, 66, 33, 0, 173, 185, 3, 70, 35, 0, 174, 185, 3, 72, 36, 0, 175, 185,
		3, 74, 37, 0, 176, 185, 3, 80, 40, 0, 177, 185, 3, 82, 41, 0, 178, 185,
		3, 86, 43, 0, 179, 185, 3, 88, 44, 0, 180, 185, 3, 90, 45, 0, 181, 185,
		3, 94, 47, 0, 182, 185, 3, 96, 48, 0, 183, 185, 3, 66, 33, 0, 184, 158,
		1, 0, 0, 0, 184, 159, 1, 0, 0, 0, 184, 160, 1, 0, 0, 0, 184, 161, 1, 0,
		0, 0, 184, 162, 1, 0, 0, 0, 184, 163, 1, 0, 0, 0, 184, 164, 1, 0, 0, 0,
		184, 165, 1, 0, 0, 0, 184, 166, 1, 0, 0, 0, 184, 167, 1, 0, 0, 0, 184,
		168, 1, 0, 0, 0, 184, 169, 1, 0, 0, 0, 184, 170, 1, 0, 0, 0, 184, 171,
		1, 0, 0, 0, 184, 172, 1, 0, 0, 0, 184, 173, 1, 0, 0, 0, 184, 174, 1, 0,
		0, 0, 184, 175, 1, 0, 0, 0, 184, 176, 1, 0, 0, 0, 184, 177, 1, 0, 0, 0,
		184, 178, 1, 0, 0, 0, 184, 179, 1, 0, 0, 0, 184, 180, 1, 0, 0, 0, 184,
		181, 1, 0, 0, 0, 184, 182, 1, 0, 0, 0, 184, 183, 1, 0, 0, 0, 185, 5, 1,
		0, 0, 0, 186, 189, 3, 100, 50, 0, 187, 189, 3, 102, 51, 0, 188, 186, 1,
		0, 0, 0, 188, 187, 1, 0, 0, 0, 189, 7, 1, 0, 0, 0, 190, 191, 5, 1, 0, 0,
		191, 192, 5, 1, 0, 0, 192, 193, 5, 59, 0, 0, 193, 198, 5, 106, 0, 0, 194,
		195, 5, 1, 0, 0, 195, 196, 5, 59, 0, 0, 196, 198, 5, 106, 0, 0, 197, 190,
		1, 0, 0, 0, 197, 194, 1, 0, 0, 0, 198, 9, 1, 0, 0, 0, 199, 200, 5, 2, 0,
		0, 200, 201, 3, 106, 53, 0, 201, 11, 1, 0, 0, 0, 202, 203, 5, 3, 0, 0,
		203, 204, 5, 68, 0, 0, 204, 208, 3, 14, 7, 0, 205, 206, 5, 3, 0, 0, 206,
		208, 3, 14, 7, 0, 207, 202, 1, 0, 0, 0, 207, 205, 1, 0, 0, 0, 208, 13,
		1, 0, 0, 0, 209, 214, 3, 106, 53, 0, 210, 211, 5, 100, 0, 0, 211, 213,
		3, 106, 53, 0, 212, 210, 1, 0, 0, 0, 213, 216, 1, 0, 0, 0, 214, 212, 1,
		0, 0, 0, 214, 215, 1, 0, 0, 0, 215, 15, 1, 0, 0, 0, 216, 214, 1, 0, 0,
		0, 217, 218, 5, 4, 0, 0, 218, 221, 3, 18, 9, 0, 219, 220, 5, 30, 0, 0,
		220, 222, 3, 14, 7, 0, 221, 219, 1, 0, 0, 0, 221, 222, 1, 0, 0, 0, 222,
		17, 1, 0, 0, 0, 223, 228, 3, 20, 10, 0, 224, 225, 5, 100, 0, 0, 225, 227,
		3, 20, 10, 0, 226, 224, 1, 0, 0, 0, 227, 230, 1, 0, 0, 0, 228, 226, 1,
		0, 0, 0, 228, 229, 1, 0, 0, 0, 229, 19, 1, 0, 0, 0, 230, 228, 1, 0, 0,
		0, 231, 234, 3, 106, 53, 0, 232, 233, 5, 31, 0, 0, 233, 235, 5, 106, 0,
		0, 234, 232, 1, 0, 0, 0, 234, 235, 1, 0, 0, 0, 235, 21, 1, 0, 0, 0, 236,
		237, 5, 5, 0, 0, 237, 238, 3, 24, 12, 0, 238, 23, 1, 0, 0, 0, 239, 244,
		3, 26, 13, 0, 240, 241, 5, 100, 0, 0, 241, 243, 3, 26, 13, 0, 242, 240,
		1, 0, 0, 0, 243, 246, 1, 0, 0, 0, 244, 242, 1, 0, 0, 0, 244, 245, 1, 0,
		0, 0, 245, 25, 1, 0, 0, 0, 246, 244, 1, 0, 0, 0, 247, 249, 3, 106, 53,
		0, 248, 250, 7, 0, 0, 0, 249, 248, 1, 0, 0, 0, 249, 250, 1, 0, 0, 0, 250,
		27, 1, 0, 0, 0, 251, 252, 5, 6, 0, 0, 252, 253, 5, 107, 0, 0, 253, 29,
		1, 0, 0, 0, 254, 255, 5, 10, 0, 0, 255, 258, 3, 18, 9, 0, 256, 257, 5,
		30, 0, 0, 257, 259, 3, 14, 7, 0, 258, 256, 1, 0, 0, 0, 258, 259, 1, 0,
		0, 0, 259, 263, 1, 0, 0, 0, 260, 262, 3, 32, 16, 0, 261, 260, 1, 0, 0,
		0, 262, 265, 1, 0, 0, 0, 263, 261, 1, 0, 0, 0, 263, 264, 1, 0, 0, 0, 264,
		31, 1, 0, 0, 0, 265, 263, 1, 0, 0, 0, 266, 267, 5, 41, 0, 0, 267, 268,
		5, 59, 0, 0, 268, 282, 3, 38, 19, 0, 269, 270, 5, 46, 0, 0, 270, 271, 5,
		59, 0, 0, 271, 282, 5, 107, 0, 0, 272, 273, 5, 51, 0, 0, 273, 274, 5, 59,
		0, 0, 274, 282, 3, 98, 49, 0, 275, 276, 5, 52, 0, 0, 276, 277, 5, 59, 0,
		0, 277, 282, 5, 109, 0, 0, 278, 279, 5, 53, 0, 0, 279, 280, 5, 59, 0, 0,
		280, 282, 5, 109, 0, 0, 281, 266, 1, 0, 0, 0, 281, 269, 1, 0, 0, 0, 281,
		272, 1, 0, 0, 0, 281, 275, 1, 0, 0, 0, 281, 278, 1, 0, 0, 0, 282, 33, 1,
		0, 0, 0, 283, 287, 5, 11, 0, 0, 284, 286, 3, 36, 18, 0, 285, 284, 1, 0,
		0, 0, 286, 289, 1, 0, 0, 0, 287, 285, 1, 0, 0, 0, 287, 288, 1, 0, 0, 0,
		288, 290, 1, 0, 0, 0, 289, 287, 1, 0, 0, 0, 290, 293, 3, 18, 9, 0, 291,
		292, 5, 30, 0, 0, 292, 294, 3, 14, 7, 0, 293, 291, 1, 0, 0, 0, 293, 294,
		1, 0, 0, 0, 294, 35, 1, 0, 0, 0, 295, 296, 5, 41, 0, 0, 296, 297, 5, 59,
		0, 0, 297, 308, 3, 38, 19, 0, 298, 299, 5, 42, 0, 0, 299, 300, 5, 59, 0,
		0, 300, 308, 5, 107, 0, 0, 301, 302, 5, 46, 0, 0, 302, 303, 5, 59, 0, 0,
		303, 308, 5, 107, 0, 0, 304, 305, 5, 51, 0, 0, 305, 306, 5, 59, 0, 0, 306,
		308, 3, 98, 49, 0, 307, 295, 1, 0, 0, 0, 307, 298, 1, 0, 0, 0, 307, 301,
		1, 0, 0, 0, 307, 304, 1, 0, 0, 0, 308, 37, 1, 0, 0, 0, 309, 310, 5, 107,
		0, 0, 310, 311, 5, 106, 0, 0, 311, 39, 1, 0, 0, 0, 312, 313, 5, 12, 0,
		0, 313, 317, 3, 126, 63, 0, 314, 316, 3, 42, 21, 0, 315, 314, 1, 0, 0,
		0, 316, 319, 1, 0, 0, 0, 317, 315, 1, 0, 0, 0, 317, 318, 1, 0, 0, 0, 318,
		41, 1, 0, 0, 0, 319, 317, 1, 0, 0, 0, 320, 321, 5, 41, 0, 0, 321, 324,
		5, 59, 0, 0, 322, 325, 3, 38, 19, 0, 323, 325, 5, 106, 0, 0, 324, 322,
		1, 0, 0, 0, 324, 323, 1, 0, 0, 0, 325, 330, 1, 0, 0, 0, 326, 327, 5, 42,
		0, 0, 327, 328, 5, 59, 0, 0, 328, 330, 5, 107, 0, 0, 329, 320, 1, 0, 0,
		0, 329, 326, 1, 0, 0, 0, 330, 43, 1, 0, 0, 0, 331, 333, 5, 13, 0, 0, 332,
		334, 5, 107, 0, 0, 333, 332, 1, 0, 0, 0, 333, 334, 1, 0, 0, 0, 334, 335,
		1, 0, 0, 0, 335, 339, 3, 14, 7, 0, 336, 338, 3, 46, 23, 0, 337, 336, 1,
		0, 0, 0, 338, 341, 1, 0, 0, 0, 339, 337, 1, 0, 0, 0, 339, 340, 1, 0, 0,
		0, 340, 45, 1, 0, 0, 0, 341, 339, 1, 0, 0, 0, 342, 343, 5, 43, 0, 0, 343,
		344, 5, 59, 0, 0, 344, 351, 3, 98, 49, 0, 345, 346, 5, 44, 0, 0, 346, 347,
		5, 59, 0, 0, 347, 351, 3, 98, 49, 0, 348, 349, 5, 45, 0, 0, 349, 351, 3,
		24, 12, 0, 350, 342, 1, 0, 0, 0, 350, 345, 1, 0, 0, 0, 350, 348, 1, 0,
		0, 0, 351, 47, 1, 0, 0, 0, 352, 354, 5, 14, 0, 0, 353, 355, 5, 107, 0,
		0, 354, 353, 1, 0, 0, 0, 354, 355, 1, 0, 0, 0, 355, 356, 1, 0, 0, 0, 356,
		359, 3, 14, 7, 0, 357, 358, 5, 30, 0, 0, 358, 360, 3, 14, 7, 0, 359, 357,
		1, 0, 0, 0, 359, 360, 1, 0, 0, 0, 360, 364, 1, 0, 0, 0, 361, 363, 3, 50,
		25, 0, 362, 361, 1, 0, 0, 0, 363, 366, 1, 0, 0, 0, 364, 362, 1, 0, 0, 0,
		364, 365, 1, 0, 0, 0, 365, 49, 1, 0, 0, 0, 366, 364, 1, 0, 0, 0, 367, 368,
		5, 47, 0, 0, 368, 369, 5, 59, 0, 0, 369, 389, 5, 109, 0, 0, 370, 371, 5,
		48, 0, 0, 371, 372, 5, 59, 0, 0, 372, 389, 5, 109, 0, 0, 373, 374, 5, 49,
		0, 0, 374, 375, 5, 59, 0, 0, 375, 389, 3, 98, 49, 0, 376, 377, 5, 50, 0,
		0, 377, 378, 5, 59, 0, 0, 378, 389, 3, 98, 49, 0, 379, 380, 5, 46, 0, 0,
		380, 381, 5, 59, 0, 0, 381, 389, 5, 107, 0, 0, 382, 383, 5, 51, 0, 0, 383,
		384, 5, 59, 0, 0, 384, 389, 3, 98, 49, 0, 385, 386, 5, 52, 0, 0, 386, 387,
		5, 59, 0, 0, 387, 389, 5, 109, 0, 0, 388, 367, 1, 0, 0, 0, 388, 370, 1,
		0, 0, 0, 388, 373, 1, 0, 0, 0, 388, 376, 1, 0, 0, 0, 388, 379, 1, 0, 0,
		0, 388, 382, 1, 0, 0, 0, 388, 385, 1, 0, 0, 0, 389, 51, 1, 0, 0, 0, 390,
		392, 5, 15, 0, 0, 391, 393, 5, 107, 0, 0, 392, 391, 1, 0, 0, 0, 392, 393,
		1, 0, 0, 0, 393, 394, 1, 0, 0, 0, 394, 397, 3, 14, 7, 0, 395, 396, 5, 30,
		0, 0, 396, 398, 3, 14, 7, 0, 397, 395, 1, 0, 0, 0, 397, 398, 1, 0, 0, 0,
		398, 402, 1, 0, 0, 0, 399, 401, 3, 50, 25, 0, 400, 399, 1, 0, 0, 0, 401,
		404, 1, 0, 0, 0, 402, 400, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0, 403, 53, 1,
		0, 0, 0, 404, 402, 1, 0, 0, 0, 405, 406, 5, 16, 0, 0, 406, 411, 3, 56,
		28, 0, 407, 408, 5, 100, 0, 0, 408, 410, 3, 56, 28, 0, 409, 407, 1, 0,
		0, 0, 410, 413, 1, 0, 0, 0, 411, 409, 1, 0, 0, 0, 411, 412, 1, 0, 0, 0,
		412, 55, 1, 0, 0, 0, 413, 411, 1, 0, 0, 0, 414, 415, 5, 106, 0, 0, 415,
		416, 5, 59, 0, 0, 416, 417, 3, 106, 53, 0, 417, 57, 1, 0, 0, 0, 418, 419,
		5, 17, 0, 0, 419, 424, 3, 60, 30, 0, 420, 421, 5, 100, 0, 0, 421, 423,
		3, 60, 30, 0, 422, 420, 1, 0, 0, 0, 423, 426, 1, 0, 0, 0, 424, 422, 1,
		0, 0, 0, 424, 425, 1, 0, 0, 0, 425, 59, 1, 0, 0, 0, 426, 424, 1, 0, 0,
		0, 427, 428, 5, 106, 0, 0, 428, 429, 5, 31, 0, 0, 429, 430, 5, 106, 0,
		0, 430, 61, 1, 0, 0, 0, 431, 432, 5, 18, 0, 0, 432, 437, 3, 64, 32, 0,
		433, 434, 5, 100, 0, 0, 434, 436, 3, 64, 32, 0, 435, 433, 1, 0, 0, 0, 436,
		439, 1, 0, 0, 0, 437, 435, 1, 0, 0, 0, 437, 438, 1, 0, 0, 0, 438, 440,
		1, 0, 0, 0, 439, 437, 1, 0, 0, 0, 440, 441, 5, 66, 0, 0, 441, 442, 5, 106,
		0, 0, 442, 63, 1, 0, 0, 0, 443, 444, 3, 106, 53, 0, 444, 445, 5, 32, 0,
		0, 445, 446, 3, 106, 53, 0, 446, 65, 1, 0, 0, 0, 447, 448, 5, 19, 0, 0,
		448, 449, 5, 35, 0, 0, 449, 450, 5, 59, 0, 0, 450, 453, 3, 106, 53, 0,
		451, 452, 5, 3, 0, 0, 452, 454, 3, 14, 7, 0, 453, 451, 1, 0, 0, 0, 453,
		454, 1, 0, 0, 0, 454, 465, 1, 0, 0, 0, 455, 456, 5, 19, 0, 0, 456, 461,
		3, 68, 34, 0, 457, 458, 5, 100, 0, 0, 458, 460, 3, 68, 34, 0, 459, 457,
		1, 0, 0, 0, 460, 463, 1, 0, 0, 0, 461, 459, 1, 0, 0, 0, 461, 462, 1, 0,
		0, 0, 462, 465, 1, 0, 0, 0, 463, 461, 1, 0, 0, 0, 464, 447, 1, 0, 0, 0,
		464, 455, 1, 0, 0, 0, 465, 67, 1, 0, 0, 0, 466, 467, 5, 106, 0, 0, 467,
		468, 5, 59, 0, 0, 468, 469, 3, 106, 53, 0, 469, 69, 1, 0, 0, 0, 470, 473,
		5, 20, 0, 0, 471, 472, 5, 106, 0, 0, 472, 474, 5, 59, 0, 0, 473, 471, 1,
		0, 0, 0, 473, 474, 1, 0, 0, 0, 474, 475, 1, 0, 0, 0, 475, 476, 5, 106,
		0, 0, 476, 477, 5, 109, 0, 0, 477, 71, 1, 0, 0, 0, 478, 482, 5, 21, 0,
		0, 479, 480, 5, 106, 0, 0, 480, 481, 5, 59, 0, 0, 481, 483, 5, 106, 0,
		0, 482, 479, 1, 0, 0, 0, 482, 483, 1, 0, 0, 0, 483, 484, 1, 0, 0, 0, 484,
		485, 5, 109, 0, 0, 485, 73, 1, 0, 0, 0, 486, 487, 5, 22, 0, 0, 487, 488,
		5, 106, 0, 0, 488, 491, 5, 106, 0, 0, 489, 490, 5, 31, 0, 0, 490, 492,
		5, 106, 0, 0, 491, 489, 1, 0, 0, 0, 491, 492, 1, 0, 0, 0, 492, 493, 1,
		0, 0, 0, 493, 494, 5, 33, 0, 0, 494, 495, 3, 76, 38, 0, 495, 75, 1, 0,
		0, 0, 496, 501, 3, 78, 39, 0, 497, 498, 5, 100, 0, 0, 498, 500, 3, 78,
		39, 0, 499, 497, 1, 0, 0, 0, 500, 503, 1, 0, 0, 0, 501, 499, 1, 0, 0, 0,
		501, 502, 1, 0, 0, 0, 502, 77, 1, 0, 0, 0, 503, 501, 1, 0, 0, 0, 504, 507,
		5, 106, 0, 0, 505, 506, 5, 31, 0, 0, 506, 508, 5, 106, 0, 0, 507, 505,
		1, 0, 0, 0, 507, 508, 1, 0, 0, 0, 508, 79, 1, 0, 0, 0, 509, 510, 5, 23,
		0, 0, 510, 511, 5, 104, 0, 0, 511, 512, 3, 2, 1, 0, 512, 513, 5, 105, 0,
		0, 513, 81, 1, 0, 0, 0, 514, 518, 5, 24, 0, 0, 515, 516, 5, 34, 0, 0, 516,
		517, 5, 59, 0, 0, 517, 519, 3, 84, 42, 0, 518, 515, 1, 0, 0, 0, 518, 519,
		1, 0, 0, 0, 519, 520, 1, 0, 0, 0, 520, 521, 5, 106, 0, 0, 521, 522, 5,
		104, 0, 0, 522, 523, 3, 2, 1, 0, 523, 524, 5, 105, 0, 0, 524, 83, 1, 0,
		0, 0, 525, 526, 7, 1, 0, 0, 526, 85, 1, 0, 0, 0, 527, 528, 5, 25, 0, 0,
		528, 529, 3, 14, 7, 0, 529, 87, 1, 0, 0, 0, 530, 531, 5, 26, 0, 0, 531,
		534, 3, 18, 9, 0, 532, 533, 5, 30, 0, 0, 533, 535, 3, 14, 7, 0, 534, 532,
		1, 0, 0, 0, 534, 535, 1, 0, 0, 0, 535, 89, 1, 0, 0, 0, 536, 540, 5, 27,
		0, 0, 537, 539, 3, 92, 46, 0, 538, 537, 1, 0, 0, 0, 539, 542, 1, 0, 0,
		0, 540, 538, 1, 0, 0, 0, 540, 541, 1, 0, 0, 0, 541, 543, 1, 0, 0, 0, 542,
		540, 1, 0, 0, 0, 543, 546, 3, 18, 9, 0, 544, 545, 5, 30, 0, 0, 545, 547,
		3, 14, 7, 0, 546, 544, 1, 0, 0, 0, 546, 547, 1, 0, 0, 0, 547, 91, 1, 0,
		0, 0, 548, 549, 5, 106, 0, 0, 549, 553, 5, 59, 0, 0, 550, 554, 5, 107,
		0, 0, 551, 554, 5, 109, 0, 0, 552, 554, 3, 98, 49, 0, 553, 550, 1, 0, 0,
		0, 553, 551, 1, 0, 0, 0, 553, 552, 1, 0, 0, 0, 554, 93, 1, 0, 0, 0, 555,
		556, 5, 28, 0, 0, 556, 95, 1, 0, 0, 0, 557, 558, 5, 29, 0, 0, 558, 559,
		3, 126, 63, 0, 559, 97, 1, 0, 0, 0, 560, 561, 7, 2, 0, 0, 561, 99, 1, 0,
		0, 0, 562, 563, 5, 7, 0, 0, 563, 564, 5, 106, 0, 0, 564, 101, 1, 0, 0,
		0, 565, 566, 5, 8, 0, 0, 566, 103, 1, 0, 0, 0, 567, 568, 5, 9, 0, 0, 568,
		105, 1, 0, 0, 0, 569, 570, 3, 108, 54, 0, 570, 107, 1, 0, 0, 0, 571, 576,
		3, 110, 55, 0, 572, 573, 5, 57, 0, 0, 573, 575, 3, 110, 55, 0, 574, 572,
		1, 0, 0, 0, 575, 578, 1, 0, 0, 0, 576, 574, 1, 0, 0, 0, 576, 577, 1, 0,
		0, 0, 577, 109, 1, 0, 0, 0, 578, 576, 1, 0, 0, 0, 579, 584, 3, 112, 56,
		0, 580, 581, 5, 56, 0, 0, 581, 583, 3, 112, 56, 0, 582, 580, 1, 0, 0, 0,
		583, 586, 1, 0, 0, 0, 584, 582, 1, 0, 0, 0, 584, 585, 1, 0, 0, 0, 585,
		111, 1, 0, 0, 0, 586, 584, 1, 0, 0, 0, 587, 588, 5, 58, 0, 0, 588, 591,
		3, 112, 56, 0, 589, 591, 3, 114, 57, 0, 590, 587, 1, 0, 0, 0, 590, 589,
		1, 0, 0, 0, 591, 113, 1, 0, 0, 0, 592, 602, 3, 116, 58, 0, 593, 594, 7,
		3, 0, 0, 594, 603, 3, 116, 58, 0, 595, 596, 5, 65, 0, 0, 596, 603, 3, 116,
		58, 0, 597, 598, 5, 66, 0, 0, 598, 599, 5, 102, 0, 0, 599, 600, 3, 132,
		66, 0, 600, 601, 5, 103, 0, 0, 601, 603, 1, 0, 0, 0, 602, 593, 1, 0, 0,
		0, 602, 595, 1, 0, 0, 0, 602, 597, 1, 0, 0, 0, 602, 603, 1, 0, 0, 0, 603,
		115, 1, 0, 0, 0, 604, 609, 3, 118, 59, 0, 605, 606, 7, 4, 0, 0, 606, 608,
		3, 118, 59, 0, 607, 605, 1, 0, 0, 0, 608, 611, 1, 0, 0, 0, 609, 607, 1,
		0, 0, 0, 609, 610, 1, 0, 0, 0, 610, 117, 1, 0, 0, 0, 611, 609, 1, 0, 0,
		0, 612, 617, 3, 120, 60, 0, 613, 614, 7, 5, 0, 0, 614, 616, 3, 120, 60,
		0, 615, 613, 1, 0, 0, 0, 616, 619, 1, 0, 0, 0, 617, 615, 1, 0, 0, 0, 617,
		618, 1, 0, 0, 0, 618, 119, 1, 0, 0, 0, 619, 617, 1, 0, 0, 0, 620, 621,
		7, 4, 0, 0, 621, 624, 3, 120, 60, 0, 622, 624, 3, 122, 61, 0, 623, 620,
		1, 0, 0, 0, 623, 622, 1, 0, 0, 0, 624, 121, 1, 0, 0, 0, 625, 634, 3, 124,
		62, 0, 626, 634, 3, 126, 63, 0, 627, 634, 3, 128, 64, 0, 628, 634, 3, 134,
		67, 0, 629, 630, 5, 102, 0, 0, 630, 631, 3, 106, 53, 0, 631, 632, 5, 103,
		0, 0, 632, 634, 1, 0, 0, 0, 633, 625, 1, 0, 0, 0, 633, 626, 1, 0, 0, 0,
		633, 627, 1, 0, 0, 0, 633, 628, 1, 0, 0, 0, 633, 629, 1, 0, 0, 0, 634,
		123, 1, 0, 0, 0, 635, 636, 7, 6, 0, 0, 636, 125, 1, 0, 0, 0, 637, 642,
		5, 106, 0, 0, 638, 639, 5, 101, 0, 0, 639, 641, 5, 106, 0, 0, 640, 638,
		1, 0, 0, 0, 641, 644, 1, 0, 0, 0, 642, 640, 1, 0, 0, 0, 642, 643, 1, 0,
		0, 0, 643, 650, 1, 0, 0, 0, 644, 642, 1, 0, 0, 0, 645, 646, 5, 106, 0,
		0, 646, 647, 5, 104, 0, 0, 647, 648, 5, 107, 0, 0, 648, 650, 5, 105, 0,
		0, 649, 637, 1, 0, 0, 0, 649, 645, 1, 0, 0, 0, 650, 127, 1, 0, 0, 0, 651,
		652, 5, 106, 0, 0, 652, 653, 5, 102, 0, 0, 653, 675, 5, 103, 0, 0, 654,
		655, 5, 106, 0, 0, 655, 657, 5, 102, 0, 0, 656, 658, 5, 98, 0, 0, 657,
		656, 1, 0, 0, 0, 657, 658, 1, 0, 0, 0, 658, 659, 1, 0, 0, 0, 659, 660,
		3, 132, 66, 0, 660, 661, 5, 103, 0, 0, 661, 675, 1, 0, 0, 0, 662, 663,
		3, 130, 65, 0, 663, 664, 5, 102, 0, 0, 664, 665, 5, 103, 0, 0, 665, 675,
		1, 0, 0, 0, 666, 667, 3, 130, 65, 0, 667, 669, 5, 102, 0, 0, 668, 670,
		5, 98, 0, 0, 669, 668, 1, 0, 0, 0, 669, 670, 1, 0, 0, 0, 670, 671, 1, 0,
		0, 0, 671, 672, 3, 132, 66, 0, 672, 673, 5, 103, 0, 0, 673, 675, 1, 0,
		0, 0, 674, 651, 1, 0, 0, 0, 674, 654, 1, 0, 0, 0, 674, 662, 1, 0, 0, 0,
		674, 666, 1, 0, 0, 0, 675, 129, 1, 0, 0, 0, 676, 677, 7, 7, 0, 0, 677,
		131, 1, 0, 0, 0, 678, 683, 3, 106, 53, 0, 679, 680, 5, 100, 0, 0, 680,
		682, 3, 106, 53, 0, 681, 679, 1, 0, 0, 0, 682, 685, 1, 0, 0, 0, 683, 681,
		1, 0, 0, 0, 683, 684, 1, 0, 0, 0, 684, 133, 1, 0, 0, 0, 685, 683, 1, 0,
		0, 0, 686, 688, 5, 75, 0, 0, 687, 689, 3, 136, 68, 0, 688, 687, 1, 0, 0,
		0, 689, 690, 1, 0, 0, 0, 690, 688, 1, 0, 0, 0, 690, 691, 1, 0, 0, 0, 691,
		694, 1, 0, 0, 0, 692, 693, 5, 78, 0, 0, 693, 695, 3, 106, 53, 0, 694, 692,
		1, 0, 0, 0, 694, 695, 1, 0, 0, 0, 695, 696, 1, 0, 0, 0, 696, 697, 5, 79,
		0, 0, 697, 135, 1, 0, 0, 0, 698, 699, 5, 76, 0, 0, 699, 700, 3, 106, 53,
		0, 700, 701, 5, 77, 0, 0, 701, 702, 3, 106, 53, 0, 702, 137, 1, 0, 0, 0,
		63, 148, 155, 184, 188, 197, 207, 214, 221, 228, 234, 244, 249, 258, 263,
		281, 287, 293, 307, 317, 324, 329, 333, 339, 350, 354, 359, 364, 388, 392,
		397, 402, 411, 424, 437, 453, 461, 464, 473, 482, 491, 501, 507, 518, 534,
		540, 546, 553, 576, 584, 590, 602, 609, 617, 623, 633, 642, 649, 657, 669,
		674, 683, 690, 694,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// PPLParserInit initializes any static state used to implement PPLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewPPLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func PPLParserInit() {
	staticData := &PPLParserParserStaticData
	staticData.once.Do(pplparserParserInit)
}

// NewPPLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewPPLParser(input antlr.TokenStream) *PPLParser {
	PPLParserInit()
	this := new(PPLParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &PPLParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "PPLParser.g4"

	return this
}

// PPLParser tokens.
const (
	PPLParserEOF             = antlr.TokenEOF
	PPLParserSEARCH          = 1
	PPLParserWHERE           = 2
	PPLParserFIELDS          = 3
	PPLParserSTATS           = 4
	PPLParserSORT            = 5
	PPLParserHEAD            = 6
	PPLParserDESCRIBE        = 7
	PPLParserSHOWDATASOURCES = 8
	PPLParserEXPLAIN         = 9
	PPLParserCHART           = 10
	PPLParserTIMECHART       = 11
	PPLParserBIN             = 12
	PPLParserDEDUP           = 13
	PPLParserTOP             = 14
	PPLParserRARE            = 15
	PPLParserEVAL            = 16
	PPLParserRENAME          = 17
	PPLParserREPLACE         = 18
	PPLParserFILLNULL        = 19
	PPLParserPARSE           = 20
	PPLParserREX             = 21
	PPLParserLOOKUP          = 22
	PPLParserAPPEND          = 23
	PPLParserJOIN            = 24
	PPLParserTABLE           = 25
	PPLParserEVENTSTATS      = 26
	PPLParserSTREAMSTATS     = 27
	PPLParserREVERSE         = 28
	PPLParserFLATTEN         = 29
	PPLParserBY              = 30
	PPLParserAS              = 31
	PPLParserWITH            = 32
	PPLParserOUTPUT          = 33
	PPLParserTYPE            = 34
	PPLParserVALUE           = 35
	PPLParserINNER           = 36
	PPLParserLEFT            = 37
	PPLParserRIGHT           = 38
	PPLParserOUTER           = 39
	PPLParserFULL            = 40
	PPLParserSPAN            = 41
	PPLParserBINS            = 42
	PPLParserKEEPEVENTS      = 43
	PPLParserCONSECUTIVE     = 44
	PPLParserSORTBY          = 45
	PPLParserLIMIT           = 46
	PPLParserCOUNTFIELD      = 47
	PPLParserPERCENTFIELD    = 48
	PPLParserSHOWCOUNT       = 49
	PPLParserSHOWPERC        = 50
	PPLParserUSEOTHER        = 51
	PPLParserOTHERSTR        = 52
	PPLParserNULLSTR         = 53
	PPLParserASC             = 54
	PPLParserDESC            = 55
	PPLParserAND             = 56
	PPLParserOR              = 57
	PPLParserNOT             = 58
	PPLParserEQ              = 59
	PPLParserNEQ             = 60
	PPLParserLT              = 61
	PPLParserLTE             = 62
	PPLParserGT              = 63
	PPLParserGTE             = 64
	PPLParserLIKE            = 65
	PPLParserIN              = 66
	PPLParserPLUS            = 67
	PPLParserMINUS           = 68
	PPLParserSTAR            = 69
	PPLParserSLASH           = 70
	PPLParserPERCENT         = 71
	PPLParserTRUE            = 72
	PPLParserFALSE           = 73
	PPLParserNULL            = 74
	PPLParserCASE            = 75
	PPLParserWHEN            = 76
	PPLParserTHEN            = 77
	PPLParserELSE            = 78
	PPLParserEND             = 79
	PPLParserCOUNT           = 80
	PPLParserSUM             = 81
	PPLParserAVG             = 82
	PPLParserMIN             = 83
	PPLParserMAX             = 84
	PPLParserDC              = 85
	PPLParserDISTINCT_COUNT  = 86
	PPLParserVAR             = 87
	PPLParserVARP            = 88
	PPLParserSTDEV           = 89
	PPLParserSTDEVP          = 90
	PPLParserPERCENTILE      = 91
	PPLParserMEDIAN          = 92
	PPLParserMODE            = 93
	PPLParserEARLIEST        = 94
	PPLParserLATEST          = 95
	PPLParserVALUES          = 96
	PPLParserRANGE           = 97
	PPLParserDISTINCT        = 98
	PPLParserPIPE            = 99
	PPLParserCOMMA           = 100
	PPLParserDOT             = 101
	PPLParserLPAREN          = 102
	PPLParserRPAREN          = 103
	PPLParserLBRACKET        = 104
	PPLParserRBRACKET        = 105
	PPLParserIDENTIFIER      = 106
	PPLParserINTEGER         = 107
	PPLParserDECIMAL         = 108
	PPLParserSTRING          = 109
	PPLParserWS              = 110
	PPLParserLINE_COMMENT    = 111
	PPLParserBLOCK_COMMENT   = 112
)

// PPLParser rules.
const (
	PPLParserRULE_query                    = 0
	PPLParserRULE_searchQuery              = 1
	PPLParserRULE_processingCommand        = 2
	PPLParserRULE_metadataCommand          = 3
	PPLParserRULE_searchCommand            = 4
	PPLParserRULE_whereCommand             = 5
	PPLParserRULE_fieldsCommand            = 6
	PPLParserRULE_fieldList                = 7
	PPLParserRULE_statsCommand             = 8
	PPLParserRULE_aggregationList          = 9
	PPLParserRULE_aggregation              = 10
	PPLParserRULE_sortCommand              = 11
	PPLParserRULE_sortFieldList            = 12
	PPLParserRULE_sortField                = 13
	PPLParserRULE_headCommand              = 14
	PPLParserRULE_chartCommand             = 15
	PPLParserRULE_chartOptions             = 16
	PPLParserRULE_timechartCommand         = 17
	PPLParserRULE_timechartOptions         = 18
	PPLParserRULE_timeSpan                 = 19
	PPLParserRULE_binCommand               = 20
	PPLParserRULE_binOptions               = 21
	PPLParserRULE_dedupCommand             = 22
	PPLParserRULE_dedupOptions             = 23
	PPLParserRULE_topCommand               = 24
	PPLParserRULE_topOptions               = 25
	PPLParserRULE_rareCommand              = 26
	PPLParserRULE_evalCommand              = 27
	PPLParserRULE_evalAssignment           = 28
	PPLParserRULE_renameCommand            = 29
	PPLParserRULE_renameAssignment         = 30
	PPLParserRULE_replaceCommand           = 31
	PPLParserRULE_replaceMapping           = 32
	PPLParserRULE_fillnullCommand          = 33
	PPLParserRULE_fillnullAssignment       = 34
	PPLParserRULE_parseCommand             = 35
	PPLParserRULE_rexCommand               = 36
	PPLParserRULE_lookupCommand            = 37
	PPLParserRULE_lookupOutputList         = 38
	PPLParserRULE_lookupOutputField        = 39
	PPLParserRULE_appendCommand            = 40
	PPLParserRULE_joinCommand              = 41
	PPLParserRULE_joinType                 = 42
	PPLParserRULE_tableCommand             = 43
	PPLParserRULE_eventstatsCommand        = 44
	PPLParserRULE_streamstatsCommand       = 45
	PPLParserRULE_streamstatsOptions       = 46
	PPLParserRULE_reverseCommand           = 47
	PPLParserRULE_flattenCommand           = 48
	PPLParserRULE_booleanValue             = 49
	PPLParserRULE_describeCommand          = 50
	PPLParserRULE_showDatasourcesCommand   = 51
	PPLParserRULE_explainCommand           = 52
	PPLParserRULE_expression               = 53
	PPLParserRULE_orExpression             = 54
	PPLParserRULE_andExpression            = 55
	PPLParserRULE_notExpression            = 56
	PPLParserRULE_comparisonExpression     = 57
	PPLParserRULE_additiveExpression       = 58
	PPLParserRULE_multiplicativeExpression = 59
	PPLParserRULE_unaryExpression          = 60
	PPLParserRULE_primaryExpression        = 61
	PPLParserRULE_literal                  = 62
	PPLParserRULE_fieldReference           = 63
	PPLParserRULE_functionCall             = 64
	PPLParserRULE_aggregationFunction      = 65
	PPLParserRULE_expressionList           = 66
	PPLParserRULE_caseExpression           = 67
	PPLParserRULE_whenClause               = 68
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ExplainCommand() IExplainCommandContext
	SearchQuery() ISearchQueryContext
	EOF() antlr.TerminalNode
	MetadataCommand() IMetadataCommandContext

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) ExplainCommand() IExplainCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExplainCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExplainCommandContext)
}

func (s *QueryContext) SearchQuery() ISearchQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISearchQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISearchQueryContext)
}

func (s *QueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(PPLParserEOF, 0)
}

func (s *QueryContext) MetadataCommand() IMetadataCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMetadataCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMetadataCommandContext)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (s *QueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, PPLParserRULE_query)
	p.SetState(148)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserEXPLAIN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(138)
			p.ExplainCommand()
		}
		{
			p.SetState(139)
			p.SearchQuery()
		}
		{
			p.SetState(140)
			p.Match(PPLParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserSEARCH:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(142)
			p.SearchQuery()
		}
		{
			p.SetState(143)
			p.Match(PPLParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserDESCRIBE, PPLParserSHOWDATASOURCES:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(145)
			p.MetadataCommand()
		}
		{
			p.SetState(146)
			p.Match(PPLParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISearchQueryContext is an interface to support dynamic dispatch.
type ISearchQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SearchCommand() ISearchCommandContext
	AllPIPE() []antlr.TerminalNode
	PIPE(i int) antlr.TerminalNode
	AllProcessingCommand() []IProcessingCommandContext
	ProcessingCommand(i int) IProcessingCommandContext

	// IsSearchQueryContext differentiates from other interfaces.
	IsSearchQueryContext()
}

type SearchQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySearchQueryContext() *SearchQueryContext {
	var p = new(SearchQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_searchQuery
	return p
}

func InitEmptySearchQueryContext(p *SearchQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_searchQuery
}

func (*SearchQueryContext) IsSearchQueryContext() {}

func NewSearchQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SearchQueryContext {
	var p = new(SearchQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_searchQuery

	return p
}

func (s *SearchQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *SearchQueryContext) SearchCommand() ISearchCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISearchCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISearchCommandContext)
}

func (s *SearchQueryContext) AllPIPE() []antlr.TerminalNode {
	return s.GetTokens(PPLParserPIPE)
}

func (s *SearchQueryContext) PIPE(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserPIPE, i)
}

func (s *SearchQueryContext) AllProcessingCommand() []IProcessingCommandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IProcessingCommandContext); ok {
			len++
		}
	}

	tst := make([]IProcessingCommandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IProcessingCommandContext); ok {
			tst[i] = t.(IProcessingCommandContext)
			i++
		}
	}

	return tst
}

func (s *SearchQueryContext) ProcessingCommand(i int) IProcessingCommandContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IProcessingCommandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IProcessingCommandContext)
}

func (s *SearchQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SearchQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SearchQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterSearchQuery(s)
	}
}

func (s *SearchQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitSearchQuery(s)
	}
}

func (s *SearchQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitSearchQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) SearchQuery() (localctx ISearchQueryContext) {
	localctx = NewSearchQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PPLParserRULE_searchQuery)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(150)
		p.SearchCommand()
	}
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserPIPE {
		{
			p.SetState(151)
			p.Match(PPLParserPIPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(152)
			p.ProcessingCommand()
		}

		p.SetState(157)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IProcessingCommandContext is an interface to support dynamic dispatch.
type IProcessingCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WhereCommand() IWhereCommandContext
	FieldsCommand() IFieldsCommandContext
	StatsCommand() IStatsCommandContext
	SortCommand() ISortCommandContext
	HeadCommand() IHeadCommandContext
	ChartCommand() IChartCommandContext
	TimechartCommand() ITimechartCommandContext
	BinCommand() IBinCommandContext
	DedupCommand() IDedupCommandContext
	TopCommand() ITopCommandContext
	RareCommand() IRareCommandContext
	EvalCommand() IEvalCommandContext
	RenameCommand() IRenameCommandContext
	ReplaceCommand() IReplaceCommandContext
	FillnullCommand() IFillnullCommandContext
	ParseCommand() IParseCommandContext
	RexCommand() IRexCommandContext
	LookupCommand() ILookupCommandContext
	AppendCommand() IAppendCommandContext
	JoinCommand() IJoinCommandContext
	TableCommand() ITableCommandContext
	EventstatsCommand() IEventstatsCommandContext
	StreamstatsCommand() IStreamstatsCommandContext
	ReverseCommand() IReverseCommandContext
	FlattenCommand() IFlattenCommandContext

	// IsProcessingCommandContext differentiates from other interfaces.
	IsProcessingCommandContext()
}

type ProcessingCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProcessingCommandContext() *ProcessingCommandContext {
	var p = new(ProcessingCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_processingCommand
	return p
}

func InitEmptyProcessingCommandContext(p *ProcessingCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_processingCommand
}

func (*ProcessingCommandContext) IsProcessingCommandContext() {}

func NewProcessingCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProcessingCommandContext {
	var p = new(ProcessingCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_processingCommand

	return p
}

func (s *ProcessingCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ProcessingCommandContext) WhereCommand() IWhereCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereCommandContext)
}

func (s *ProcessingCommandContext) FieldsCommand() IFieldsCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldsCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldsCommandContext)
}

func (s *ProcessingCommandContext) StatsCommand() IStatsCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatsCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatsCommandContext)
}

func (s *ProcessingCommandContext) SortCommand() ISortCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortCommandContext)
}

func (s *ProcessingCommandContext) HeadCommand() IHeadCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHeadCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHeadCommandContext)
}

func (s *ProcessingCommandContext) ChartCommand() IChartCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IChartCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IChartCommandContext)
}

func (s *ProcessingCommandContext) TimechartCommand() ITimechartCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimechartCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimechartCommandContext)
}

func (s *ProcessingCommandContext) BinCommand() IBinCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBinCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBinCommandContext)
}

func (s *ProcessingCommandContext) DedupCommand() IDedupCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDedupCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDedupCommandContext)
}

func (s *ProcessingCommandContext) TopCommand() ITopCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITopCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITopCommandContext)
}

func (s *ProcessingCommandContext) RareCommand() IRareCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRareCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRareCommandContext)
}

func (s *ProcessingCommandContext) EvalCommand() IEvalCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEvalCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEvalCommandContext)
}

func (s *ProcessingCommandContext) RenameCommand() IRenameCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRenameCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRenameCommandContext)
}

func (s *ProcessingCommandContext) ReplaceCommand() IReplaceCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReplaceCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReplaceCommandContext)
}

func (s *ProcessingCommandContext) FillnullCommand() IFillnullCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFillnullCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFillnullCommandContext)
}

func (s *ProcessingCommandContext) ParseCommand() IParseCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParseCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParseCommandContext)
}

func (s *ProcessingCommandContext) RexCommand() IRexCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRexCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRexCommandContext)
}

func (s *ProcessingCommandContext) LookupCommand() ILookupCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILookupCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILookupCommandContext)
}

func (s *ProcessingCommandContext) AppendCommand() IAppendCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAppendCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAppendCommandContext)
}

func (s *ProcessingCommandContext) JoinCommand() IJoinCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoinCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoinCommandContext)
}

func (s *ProcessingCommandContext) TableCommand() ITableCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableCommandContext)
}

func (s *ProcessingCommandContext) EventstatsCommand() IEventstatsCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEventstatsCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEventstatsCommandContext)
}

func (s *ProcessingCommandContext) StreamstatsCommand() IStreamstatsCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStreamstatsCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStreamstatsCommandContext)
}

func (s *ProcessingCommandContext) ReverseCommand() IReverseCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReverseCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReverseCommandContext)
}

func (s *ProcessingCommandContext) FlattenCommand() IFlattenCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFlattenCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFlattenCommandContext)
}

func (s *ProcessingCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProcessingCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProcessingCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterProcessingCommand(s)
	}
}

func (s *ProcessingCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitProcessingCommand(s)
	}
}

func (s *ProcessingCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitProcessingCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ProcessingCommand() (localctx IProcessingCommandContext) {
	localctx = NewProcessingCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PPLParserRULE_processingCommand)
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(158)
			p.WhereCommand()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(159)
			p.FieldsCommand()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(160)
			p.StatsCommand()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(161)
			p.SortCommand()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(162)
			p.HeadCommand()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(163)
			p.ChartCommand()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(164)
			p.TimechartCommand()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(165)
			p.BinCommand()
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(166)
			p.DedupCommand()
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(167)
			p.TopCommand()
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(168)
			p.RareCommand()
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(169)
			p.EvalCommand()
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(170)
			p.RenameCommand()
		}

	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(171)
			p.ReplaceCommand()
		}

	case 15:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(172)
			p.FillnullCommand()
		}

	case 16:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(173)
			p.ParseCommand()
		}

	case 17:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(174)
			p.RexCommand()
		}

	case 18:
		p.EnterOuterAlt(localctx, 18)
		{
			p.SetState(175)
			p.LookupCommand()
		}

	case 19:
		p.EnterOuterAlt(localctx, 19)
		{
			p.SetState(176)
			p.AppendCommand()
		}

	case 20:
		p.EnterOuterAlt(localctx, 20)
		{
			p.SetState(177)
			p.JoinCommand()
		}

	case 21:
		p.EnterOuterAlt(localctx, 21)
		{
			p.SetState(178)
			p.TableCommand()
		}

	case 22:
		p.EnterOuterAlt(localctx, 22)
		{
			p.SetState(179)
			p.EventstatsCommand()
		}

	case 23:
		p.EnterOuterAlt(localctx, 23)
		{
			p.SetState(180)
			p.StreamstatsCommand()
		}

	case 24:
		p.EnterOuterAlt(localctx, 24)
		{
			p.SetState(181)
			p.ReverseCommand()
		}

	case 25:
		p.EnterOuterAlt(localctx, 25)
		{
			p.SetState(182)
			p.FlattenCommand()
		}

	case 26:
		p.EnterOuterAlt(localctx, 26)
		{
			p.SetState(183)
			p.FillnullCommand()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMetadataCommandContext is an interface to support dynamic dispatch.
type IMetadataCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DescribeCommand() IDescribeCommandContext
	ShowDatasourcesCommand() IShowDatasourcesCommandContext

	// IsMetadataCommandContext differentiates from other interfaces.
	IsMetadataCommandContext()
}

type MetadataCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetadataCommandContext() *MetadataCommandContext {
	var p = new(MetadataCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_metadataCommand
	return p
}

func InitEmptyMetadataCommandContext(p *MetadataCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_metadataCommand
}

func (*MetadataCommandContext) IsMetadataCommandContext() {}

func NewMetadataCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetadataCommandContext {
	var p = new(MetadataCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_metadataCommand

	return p
}

func (s *MetadataCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *MetadataCommandContext) DescribeCommand() IDescribeCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDescribeCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDescribeCommandContext)
}

func (s *MetadataCommandContext) ShowDatasourcesCommand() IShowDatasourcesCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowDatasourcesCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowDatasourcesCommandContext)
}

func (s *MetadataCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetadataCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MetadataCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterMetadataCommand(s)
	}
}

func (s *MetadataCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitMetadataCommand(s)
	}
}

func (s *MetadataCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitMetadataCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) MetadataCommand() (localctx IMetadataCommandContext) {
	localctx = NewMetadataCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PPLParserRULE_metadataCommand)
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserDESCRIBE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(186)
			p.DescribeCommand()
		}

	case PPLParserSHOWDATASOURCES:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(187)
			p.ShowDatasourcesCommand()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISearchCommandContext is an interface to support dynamic dispatch.
type ISearchCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSearchCommandContext differentiates from other interfaces.
	IsSearchCommandContext()
}

type SearchCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySearchCommandContext() *SearchCommandContext {
	var p = new(SearchCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_searchCommand
	return p
}

func InitEmptySearchCommandContext(p *SearchCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_searchCommand
}

func (*SearchCommandContext) IsSearchCommandContext() {}

func NewSearchCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SearchCommandContext {
	var p = new(SearchCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_searchCommand

	return p
}

func (s *SearchCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *SearchCommandContext) CopyAll(ctx *SearchCommandContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *SearchCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SearchCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SearchWithKeywordContext struct {
	SearchCommandContext
}

func NewSearchWithKeywordContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SearchWithKeywordContext {
	var p = new(SearchWithKeywordContext)

	InitEmptySearchCommandContext(&p.SearchCommandContext)
	p.parser = parser
	p.CopyAll(ctx.(*SearchCommandContext))

	return p
}

func (s *SearchWithKeywordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SearchWithKeywordContext) AllSEARCH() []antlr.TerminalNode {
	return s.GetTokens(PPLParserSEARCH)
}

func (s *SearchWithKeywordContext) SEARCH(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserSEARCH, i)
}

func (s *SearchWithKeywordContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *SearchWithKeywordContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *SearchWithKeywordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterSearchWithKeyword(s)
	}
}

func (s *SearchWithKeywordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitSearchWithKeyword(s)
	}
}

func (s *SearchWithKeywordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitSearchWithKeyword(s)

	default:
		return t.VisitChildren(s)
	}
}

type SearchWithSourceContext struct {
	SearchCommandContext
}

func NewSearchWithSourceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SearchWithSourceContext {
	var p = new(SearchWithSourceContext)

	InitEmptySearchCommandContext(&p.SearchCommandContext)
	p.parser = parser
	p.CopyAll(ctx.(*SearchCommandContext))

	return p
}

func (s *SearchWithSourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SearchWithSourceContext) SEARCH() antlr.TerminalNode {
	return s.GetToken(PPLParserSEARCH, 0)
}

func (s *SearchWithSourceContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *SearchWithSourceContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *SearchWithSourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterSearchWithSource(s)
	}
}

func (s *SearchWithSourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitSearchWithSource(s)
	}
}

func (s *SearchWithSourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitSearchWithSource(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) SearchCommand() (localctx ISearchCommandContext) {
	localctx = NewSearchCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, PPLParserRULE_searchCommand)
	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSearchWithKeywordContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(190)
			p.Match(PPLParserSEARCH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(191)
			p.Match(PPLParserSEARCH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(192)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(193)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewSearchWithSourceContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(194)
			p.Match(PPLParserSEARCH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(195)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(196)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWhereCommandContext is an interface to support dynamic dispatch.
type IWhereCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHERE() antlr.TerminalNode
	Expression() IExpressionContext

	// IsWhereCommandContext differentiates from other interfaces.
	IsWhereCommandContext()
}

type WhereCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereCommandContext() *WhereCommandContext {
	var p = new(WhereCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_whereCommand
	return p
}

func InitEmptyWhereCommandContext(p *WhereCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_whereCommand
}

func (*WhereCommandContext) IsWhereCommandContext() {}

func NewWhereCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereCommandContext {
	var p = new(WhereCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_whereCommand

	return p
}

func (s *WhereCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereCommandContext) WHERE() antlr.TerminalNode {
	return s.GetToken(PPLParserWHERE, 0)
}

func (s *WhereCommandContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *WhereCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhereCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterWhereCommand(s)
	}
}

func (s *WhereCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitWhereCommand(s)
	}
}

func (s *WhereCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitWhereCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) WhereCommand() (localctx IWhereCommandContext) {
	localctx = NewWhereCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, PPLParserRULE_whereCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(199)
		p.Match(PPLParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(200)
		p.Expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldsCommandContext is an interface to support dynamic dispatch.
type IFieldsCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsFieldsCommandContext differentiates from other interfaces.
	IsFieldsCommandContext()
}

type FieldsCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldsCommandContext() *FieldsCommandContext {
	var p = new(FieldsCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fieldsCommand
	return p
}

func InitEmptyFieldsCommandContext(p *FieldsCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fieldsCommand
}

func (*FieldsCommandContext) IsFieldsCommandContext() {}

func NewFieldsCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldsCommandContext {
	var p = new(FieldsCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_fieldsCommand

	return p
}

func (s *FieldsCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldsCommandContext) CopyAll(ctx *FieldsCommandContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *FieldsCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldsCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type FieldsExcludeContext struct {
	FieldsCommandContext
}

func NewFieldsExcludeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FieldsExcludeContext {
	var p = new(FieldsExcludeContext)

	InitEmptyFieldsCommandContext(&p.FieldsCommandContext)
	p.parser = parser
	p.CopyAll(ctx.(*FieldsCommandContext))

	return p
}

func (s *FieldsExcludeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldsExcludeContext) FIELDS() antlr.TerminalNode {
	return s.GetToken(PPLParserFIELDS, 0)
}

func (s *FieldsExcludeContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PPLParserMINUS, 0)
}

func (s *FieldsExcludeContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *FieldsExcludeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFieldsExclude(s)
	}
}

func (s *FieldsExcludeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFieldsExclude(s)
	}
}

func (s *FieldsExcludeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFieldsExclude(s)

	default:
		return t.VisitChildren(s)
	}
}

type FieldsIncludeContext struct {
	FieldsCommandContext
}

func NewFieldsIncludeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FieldsIncludeContext {
	var p = new(FieldsIncludeContext)

	InitEmptyFieldsCommandContext(&p.FieldsCommandContext)
	p.parser = parser
	p.CopyAll(ctx.(*FieldsCommandContext))

	return p
}

func (s *FieldsIncludeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldsIncludeContext) FIELDS() antlr.TerminalNode {
	return s.GetToken(PPLParserFIELDS, 0)
}

func (s *FieldsIncludeContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *FieldsIncludeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFieldsInclude(s)
	}
}

func (s *FieldsIncludeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFieldsInclude(s)
	}
}

func (s *FieldsIncludeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFieldsInclude(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) FieldsCommand() (localctx IFieldsCommandContext) {
	localctx = NewFieldsCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, PPLParserRULE_fieldsCommand)
	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFieldsExcludeContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(202)
			p.Match(PPLParserFIELDS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(203)
			p.Match(PPLParserMINUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(204)
			p.FieldList()
		}

	case 2:
		localctx = NewFieldsIncludeContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(205)
			p.Match(PPLParserFIELDS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(206)
			p.FieldList()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldListContext is an interface to support dynamic dispatch.
type IFieldListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFieldListContext differentiates from other interfaces.
	IsFieldListContext()
}

type FieldListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldListContext() *FieldListContext {
	var p = new(FieldListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fieldList
	return p
}

func InitEmptyFieldListContext(p *FieldListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fieldList
}

func (*FieldListContext) IsFieldListContext() {}

func NewFieldListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldListContext {
	var p = new(FieldListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_fieldList

	return p
}

func (s *FieldListContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldListContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *FieldListContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FieldListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *FieldListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *FieldListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFieldList(s)
	}
}

func (s *FieldListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFieldList(s)
	}
}

func (s *FieldListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFieldList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) FieldList() (localctx IFieldListContext) {
	localctx = NewFieldListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, PPLParserRULE_fieldList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(209)
		p.Expression()
	}
	p.SetState(214)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(210)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(211)
			p.Expression()
		}

		p.SetState(216)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStatsCommandContext is an interface to support dynamic dispatch.
type IStatsCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STATS() antlr.TerminalNode
	AggregationList() IAggregationListContext
	BY() antlr.TerminalNode
	FieldList() IFieldListContext

	// IsStatsCommandContext differentiates from other interfaces.
	IsStatsCommandContext()
}

type StatsCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatsCommandContext() *StatsCommandContext {
	var p = new(StatsCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_statsCommand
	return p
}

func InitEmptyStatsCommandContext(p *StatsCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_statsCommand
}

func (*StatsCommandContext) IsStatsCommandContext() {}

func NewStatsCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatsCommandContext {
	var p = new(StatsCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_statsCommand

	return p
}

func (s *StatsCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *StatsCommandContext) STATS() antlr.TerminalNode {
	return s.GetToken(PPLParserSTATS, 0)
}

func (s *StatsCommandContext) AggregationList() IAggregationListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationListContext)
}

func (s *StatsCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(PPLParserBY, 0)
}

func (s *StatsCommandContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *StatsCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatsCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatsCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterStatsCommand(s)
	}
}

func (s *StatsCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitStatsCommand(s)
	}
}

func (s *StatsCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitStatsCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) StatsCommand() (localctx IStatsCommandContext) {
	localctx = NewStatsCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, PPLParserRULE_statsCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(217)
		p.Match(PPLParserSTATS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(218)
		p.AggregationList()
	}
	p.SetState(221)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserBY {
		{
			p.SetState(219)
			p.Match(PPLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(220)
			p.FieldList()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAggregationListContext is an interface to support dynamic dispatch.
type IAggregationListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAggregation() []IAggregationContext
	Aggregation(i int) IAggregationContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsAggregationListContext differentiates from other interfaces.
	IsAggregationListContext()
}

type AggregationListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationListContext() *AggregationListContext {
	var p = new(AggregationListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_aggregationList
	return p
}

func InitEmptyAggregationListContext(p *AggregationListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_aggregationList
}

func (*AggregationListContext) IsAggregationListContext() {}

func NewAggregationListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationListContext {
	var p = new(AggregationListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_aggregationList

	return p
}

func (s *AggregationListContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationListContext) AllAggregation() []IAggregationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAggregationContext); ok {
			len++
		}
	}

	tst := make([]IAggregationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAggregationContext); ok {
			tst[i] = t.(IAggregationContext)
			i++
		}
	}

	return tst
}

func (s *AggregationListContext) Aggregation(i int) IAggregationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationContext)
}

func (s *AggregationListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *AggregationListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *AggregationListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AggregationListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAggregationList(s)
	}
}

func (s *AggregationListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAggregationList(s)
	}
}

func (s *AggregationListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAggregationList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) AggregationList() (localctx IAggregationListContext) {
	localctx = NewAggregationListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, PPLParserRULE_aggregationList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(223)
		p.Aggregation()
	}
	p.SetState(228)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(224)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(225)
			p.Aggregation()
		}

		p.SetState(230)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAggregationContext is an interface to support dynamic dispatch.
type IAggregationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	AS() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsAggregationContext differentiates from other interfaces.
	IsAggregationContext()
}

type AggregationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationContext() *AggregationContext {
	var p = new(AggregationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_aggregation
	return p
}

func InitEmptyAggregationContext(p *AggregationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_aggregation
}

func (*AggregationContext) IsAggregationContext() {}

func NewAggregationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationContext {
	var p = new(AggregationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_aggregation

	return p
}

func (s *AggregationContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AggregationContext) AS() antlr.TerminalNode {
	return s.GetToken(PPLParserAS, 0)
}

func (s *AggregationContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *AggregationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AggregationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAggregation(s)
	}
}

func (s *AggregationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAggregation(s)
	}
}

func (s *AggregationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAggregation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) Aggregation() (localctx IAggregationContext) {
	localctx = NewAggregationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, PPLParserRULE_aggregation)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(231)
		p.Expression()
	}
	p.SetState(234)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserAS {
		{
			p.SetState(232)
			p.Match(PPLParserAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(233)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISortCommandContext is an interface to support dynamic dispatch.
type ISortCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SORT() antlr.TerminalNode
	SortFieldList() ISortFieldListContext

	// IsSortCommandContext differentiates from other interfaces.
	IsSortCommandContext()
}

type SortCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortCommandContext() *SortCommandContext {
	var p = new(SortCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_sortCommand
	return p
}

func InitEmptySortCommandContext(p *SortCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_sortCommand
}

func (*SortCommandContext) IsSortCommandContext() {}

func NewSortCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortCommandContext {
	var p = new(SortCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_sortCommand

	return p
}

func (s *SortCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *SortCommandContext) SORT() antlr.TerminalNode {
	return s.GetToken(PPLParserSORT, 0)
}

func (s *SortCommandContext) SortFieldList() ISortFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortFieldListContext)
}

func (s *SortCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterSortCommand(s)
	}
}

func (s *SortCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitSortCommand(s)
	}
}

func (s *SortCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitSortCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) SortCommand() (localctx ISortCommandContext) {
	localctx = NewSortCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, PPLParserRULE_sortCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(236)
		p.Match(PPLParserSORT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(237)
		p.SortFieldList()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISortFieldListContext is an interface to support dynamic dispatch.
type ISortFieldListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSortField() []ISortFieldContext
	SortField(i int) ISortFieldContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsSortFieldListContext differentiates from other interfaces.
	IsSortFieldListContext()
}

type SortFieldListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortFieldListContext() *SortFieldListContext {
	var p = new(SortFieldListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_sortFieldList
	return p
}

func InitEmptySortFieldListContext(p *SortFieldListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_sortFieldList
}

func (*SortFieldListContext) IsSortFieldListContext() {}

func NewSortFieldListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortFieldListContext {
	var p = new(SortFieldListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_sortFieldList

	return p
}

func (s *SortFieldListContext) GetParser() antlr.Parser { return s.parser }

func (s *SortFieldListContext) AllSortField() []ISortFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISortFieldContext); ok {
			len++
		}
	}

	tst := make([]ISortFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISortFieldContext); ok {
			tst[i] = t.(ISortFieldContext)
			i++
		}
	}

	return tst
}

func (s *SortFieldListContext) SortField(i int) ISortFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortFieldContext)
}

func (s *SortFieldListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *SortFieldListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *SortFieldListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortFieldListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortFieldListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterSortFieldList(s)
	}
}

func (s *SortFieldListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitSortFieldList(s)
	}
}

func (s *SortFieldListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitSortFieldList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) SortFieldList() (localctx ISortFieldListContext) {
	localctx = NewSortFieldListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, PPLParserRULE_sortFieldList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(239)
		p.SortField()
	}
	p.SetState(244)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(240)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(241)
			p.SortField()
		}

		p.SetState(246)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISortFieldContext is an interface to support dynamic dispatch.
type ISortFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	ASC() antlr.TerminalNode
	DESC() antlr.TerminalNode

	// IsSortFieldContext differentiates from other interfaces.
	IsSortFieldContext()
}

type SortFieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortFieldContext() *SortFieldContext {
	var p = new(SortFieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_sortField
	return p
}

func InitEmptySortFieldContext(p *SortFieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_sortField
}

func (*SortFieldContext) IsSortFieldContext() {}

func NewSortFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortFieldContext {
	var p = new(SortFieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_sortField

	return p
}

func (s *SortFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *SortFieldContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SortFieldContext) ASC() antlr.TerminalNode {
	return s.GetToken(PPLParserASC, 0)
}

func (s *SortFieldContext) DESC() antlr.TerminalNode {
	return s.GetToken(PPLParserDESC, 0)
}

func (s *SortFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterSortField(s)
	}
}

func (s *SortFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitSortField(s)
	}
}

func (s *SortFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitSortField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) SortField() (localctx ISortFieldContext) {
	localctx = NewSortFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, PPLParserRULE_sortField)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(247)
		p.Expression()
	}
	p.SetState(249)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserASC || _la == PPLParserDESC {
		{
			p.SetState(248)
			_la = p.GetTokenStream().LA(1)

			if !(_la == PPLParserASC || _la == PPLParserDESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHeadCommandContext is an interface to support dynamic dispatch.
type IHeadCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HEAD() antlr.TerminalNode
	INTEGER() antlr.TerminalNode

	// IsHeadCommandContext differentiates from other interfaces.
	IsHeadCommandContext()
}

type HeadCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHeadCommandContext() *HeadCommandContext {
	var p = new(HeadCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_headCommand
	return p
}

func InitEmptyHeadCommandContext(p *HeadCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_headCommand
}

func (*HeadCommandContext) IsHeadCommandContext() {}

func NewHeadCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HeadCommandContext {
	var p = new(HeadCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_headCommand

	return p
}

func (s *HeadCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *HeadCommandContext) HEAD() antlr.TerminalNode {
	return s.GetToken(PPLParserHEAD, 0)
}

func (s *HeadCommandContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *HeadCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HeadCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HeadCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterHeadCommand(s)
	}
}

func (s *HeadCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitHeadCommand(s)
	}
}

func (s *HeadCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitHeadCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) HeadCommand() (localctx IHeadCommandContext) {
	localctx = NewHeadCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, PPLParserRULE_headCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(251)
		p.Match(PPLParserHEAD)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(252)
		p.Match(PPLParserINTEGER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IChartCommandContext is an interface to support dynamic dispatch.
type IChartCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CHART() antlr.TerminalNode
	AggregationList() IAggregationListContext
	BY() antlr.TerminalNode
	FieldList() IFieldListContext
	AllChartOptions() []IChartOptionsContext
	ChartOptions(i int) IChartOptionsContext

	// IsChartCommandContext differentiates from other interfaces.
	IsChartCommandContext()
}

type ChartCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyChartCommandContext() *ChartCommandContext {
	var p = new(ChartCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_chartCommand
	return p
}

func InitEmptyChartCommandContext(p *ChartCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_chartCommand
}

func (*ChartCommandContext) IsChartCommandContext() {}

func NewChartCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ChartCommandContext {
	var p = new(ChartCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_chartCommand

	return p
}

func (s *ChartCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ChartCommandContext) CHART() antlr.TerminalNode {
	return s.GetToken(PPLParserCHART, 0)
}

func (s *ChartCommandContext) AggregationList() IAggregationListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationListContext)
}

func (s *ChartCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(PPLParserBY, 0)
}

func (s *ChartCommandContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *ChartCommandContext) AllChartOptions() []IChartOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IChartOptionsContext); ok {
			len++
		}
	}

	tst := make([]IChartOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IChartOptionsContext); ok {
			tst[i] = t.(IChartOptionsContext)
			i++
		}
	}

	return tst
}

func (s *ChartCommandContext) ChartOptions(i int) IChartOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IChartOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IChartOptionsContext)
}

func (s *ChartCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ChartCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ChartCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterChartCommand(s)
	}
}

func (s *ChartCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitChartCommand(s)
	}
}

func (s *ChartCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitChartCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ChartCommand() (localctx IChartCommandContext) {
	localctx = NewChartCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, PPLParserRULE_chartCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(254)
		p.Match(PPLParserCHART)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(255)
		p.AggregationList()
	}
	p.SetState(258)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserBY {
		{
			p.SetState(256)
			p.Match(PPLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(257)
			p.FieldList()
		}

	}
	p.SetState(263)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&15835166463229952) != 0 {
		{
			p.SetState(260)
			p.ChartOptions()
		}

		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IChartOptionsContext is an interface to support dynamic dispatch.
type IChartOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SPAN() antlr.TerminalNode
	EQ() antlr.TerminalNode
	TimeSpan() ITimeSpanContext
	LIMIT() antlr.TerminalNode
	INTEGER() antlr.TerminalNode
	USEOTHER() antlr.TerminalNode
	BooleanValue() IBooleanValueContext
	OTHERSTR() antlr.TerminalNode
	STRING() antlr.TerminalNode
	NULLSTR() antlr.TerminalNode

	// IsChartOptionsContext differentiates from other interfaces.
	IsChartOptionsContext()
}

type ChartOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyChartOptionsContext() *ChartOptionsContext {
	var p = new(ChartOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_chartOptions
	return p
}

func InitEmptyChartOptionsContext(p *ChartOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_chartOptions
}

func (*ChartOptionsContext) IsChartOptionsContext() {}

func NewChartOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ChartOptionsContext {
	var p = new(ChartOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_chartOptions

	return p
}

func (s *ChartOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *ChartOptionsContext) SPAN() antlr.TerminalNode {
	return s.GetToken(PPLParserSPAN, 0)
}

func (s *ChartOptionsContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *ChartOptionsContext) TimeSpan() ITimeSpanContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeSpanContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeSpanContext)
}

func (s *ChartOptionsContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(PPLParserLIMIT, 0)
}

func (s *ChartOptionsContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *ChartOptionsContext) USEOTHER() antlr.TerminalNode {
	return s.GetToken(PPLParserUSEOTHER, 0)
}

func (s *ChartOptionsContext) BooleanValue() IBooleanValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanValueContext)
}

func (s *ChartOptionsContext) OTHERSTR() antlr.TerminalNode {
	return s.GetToken(PPLParserOTHERSTR, 0)
}

func (s *ChartOptionsContext) STRING() antlr.TerminalNode {
	return s.GetToken(PPLParserSTRING, 0)
}

func (s *ChartOptionsContext) NULLSTR() antlr.TerminalNode {
	return s.GetToken(PPLParserNULLSTR, 0)
}

func (s *ChartOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ChartOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ChartOptionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterChartOptions(s)
	}
}

func (s *ChartOptionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitChartOptions(s)
	}
}

func (s *ChartOptionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitChartOptions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ChartOptions() (localctx IChartOptionsContext) {
	localctx = NewChartOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, PPLParserRULE_chartOptions)
	p.SetState(281)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserSPAN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(266)
			p.Match(PPLParserSPAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(267)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(268)
			p.TimeSpan()
		}

	case PPLParserLIMIT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(269)
			p.Match(PPLParserLIMIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(270)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(271)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserUSEOTHER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(272)
			p.Match(PPLParserUSEOTHER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(273)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(274)
			p.BooleanValue()
		}

	case PPLParserOTHERSTR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(275)
			p.Match(PPLParserOTHERSTR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(276)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(277)
			p.Match(PPLParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserNULLSTR:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(278)
			p.Match(PPLParserNULLSTR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(279)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(280)
			p.Match(PPLParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimechartCommandContext is an interface to support dynamic dispatch.
type ITimechartCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TIMECHART() antlr.TerminalNode
	AggregationList() IAggregationListContext
	AllTimechartOptions() []ITimechartOptionsContext
	TimechartOptions(i int) ITimechartOptionsContext
	BY() antlr.TerminalNode
	FieldList() IFieldListContext

	// IsTimechartCommandContext differentiates from other interfaces.
	IsTimechartCommandContext()
}

type TimechartCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimechartCommandContext() *TimechartCommandContext {
	var p = new(TimechartCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_timechartCommand
	return p
}

func InitEmptyTimechartCommandContext(p *TimechartCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_timechartCommand
}

func (*TimechartCommandContext) IsTimechartCommandContext() {}

func NewTimechartCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimechartCommandContext {
	var p = new(TimechartCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_timechartCommand

	return p
}

func (s *TimechartCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *TimechartCommandContext) TIMECHART() antlr.TerminalNode {
	return s.GetToken(PPLParserTIMECHART, 0)
}

func (s *TimechartCommandContext) AggregationList() IAggregationListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationListContext)
}

func (s *TimechartCommandContext) AllTimechartOptions() []ITimechartOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITimechartOptionsContext); ok {
			len++
		}
	}

	tst := make([]ITimechartOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITimechartOptionsContext); ok {
			tst[i] = t.(ITimechartOptionsContext)
			i++
		}
	}

	return tst
}

func (s *TimechartCommandContext) TimechartOptions(i int) ITimechartOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimechartOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimechartOptionsContext)
}

func (s *TimechartCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(PPLParserBY, 0)
}

func (s *TimechartCommandContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *TimechartCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimechartCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimechartCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterTimechartCommand(s)
	}
}

func (s *TimechartCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitTimechartCommand(s)
	}
}

func (s *TimechartCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitTimechartCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) TimechartCommand() (localctx ITimechartCommandContext) {
	localctx = NewTimechartCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, PPLParserRULE_timechartCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(283)
		p.Match(PPLParserTIMECHART)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2328765627629568) != 0 {
		{
			p.SetState(284)
			p.TimechartOptions()
		}

		p.SetState(289)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(290)
		p.AggregationList()
	}
	p.SetState(293)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserBY {
		{
			p.SetState(291)
			p.Match(PPLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(292)
			p.FieldList()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimechartOptionsContext is an interface to support dynamic dispatch.
type ITimechartOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SPAN() antlr.TerminalNode
	EQ() antlr.TerminalNode
	TimeSpan() ITimeSpanContext
	BINS() antlr.TerminalNode
	INTEGER() antlr.TerminalNode
	LIMIT() antlr.TerminalNode
	USEOTHER() antlr.TerminalNode
	BooleanValue() IBooleanValueContext

	// IsTimechartOptionsContext differentiates from other interfaces.
	IsTimechartOptionsContext()
}

type TimechartOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimechartOptionsContext() *TimechartOptionsContext {
	var p = new(TimechartOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_timechartOptions
	return p
}

func InitEmptyTimechartOptionsContext(p *TimechartOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_timechartOptions
}

func (*TimechartOptionsContext) IsTimechartOptionsContext() {}

func NewTimechartOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimechartOptionsContext {
	var p = new(TimechartOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_timechartOptions

	return p
}

func (s *TimechartOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *TimechartOptionsContext) SPAN() antlr.TerminalNode {
	return s.GetToken(PPLParserSPAN, 0)
}

func (s *TimechartOptionsContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *TimechartOptionsContext) TimeSpan() ITimeSpanContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeSpanContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeSpanContext)
}

func (s *TimechartOptionsContext) BINS() antlr.TerminalNode {
	return s.GetToken(PPLParserBINS, 0)
}

func (s *TimechartOptionsContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *TimechartOptionsContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(PPLParserLIMIT, 0)
}

func (s *TimechartOptionsContext) USEOTHER() antlr.TerminalNode {
	return s.GetToken(PPLParserUSEOTHER, 0)
}

func (s *TimechartOptionsContext) BooleanValue() IBooleanValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanValueContext)
}

func (s *TimechartOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimechartOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimechartOptionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterTimechartOptions(s)
	}
}

func (s *TimechartOptionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitTimechartOptions(s)
	}
}

func (s *TimechartOptionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitTimechartOptions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) TimechartOptions() (localctx ITimechartOptionsContext) {
	localctx = NewTimechartOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, PPLParserRULE_timechartOptions)
	p.SetState(307)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserSPAN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(295)
			p.Match(PPLParserSPAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(296)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(297)
			p.TimeSpan()
		}

	case PPLParserBINS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(298)
			p.Match(PPLParserBINS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(299)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(300)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserLIMIT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(301)
			p.Match(PPLParserLIMIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(302)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(303)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserUSEOTHER:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(304)
			p.Match(PPLParserUSEOTHER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(305)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(306)
			p.BooleanValue()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITimeSpanContext is an interface to support dynamic dispatch.
type ITimeSpanContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsTimeSpanContext differentiates from other interfaces.
	IsTimeSpanContext()
}

type TimeSpanContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeSpanContext() *TimeSpanContext {
	var p = new(TimeSpanContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_timeSpan
	return p
}

func InitEmptyTimeSpanContext(p *TimeSpanContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_timeSpan
}

func (*TimeSpanContext) IsTimeSpanContext() {}

func NewTimeSpanContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeSpanContext {
	var p = new(TimeSpanContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_timeSpan

	return p
}

func (s *TimeSpanContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeSpanContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *TimeSpanContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *TimeSpanContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeSpanContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TimeSpanContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterTimeSpan(s)
	}
}

func (s *TimeSpanContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitTimeSpan(s)
	}
}

func (s *TimeSpanContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitTimeSpan(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) TimeSpan() (localctx ITimeSpanContext) {
	localctx = NewTimeSpanContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, PPLParserRULE_timeSpan)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(309)
		p.Match(PPLParserINTEGER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(310)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBinCommandContext is an interface to support dynamic dispatch.
type IBinCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BIN() antlr.TerminalNode
	FieldReference() IFieldReferenceContext
	AllBinOptions() []IBinOptionsContext
	BinOptions(i int) IBinOptionsContext

	// IsBinCommandContext differentiates from other interfaces.
	IsBinCommandContext()
}

type BinCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinCommandContext() *BinCommandContext {
	var p = new(BinCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_binCommand
	return p
}

func InitEmptyBinCommandContext(p *BinCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_binCommand
}

func (*BinCommandContext) IsBinCommandContext() {}

func NewBinCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinCommandContext {
	var p = new(BinCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_binCommand

	return p
}

func (s *BinCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *BinCommandContext) BIN() antlr.TerminalNode {
	return s.GetToken(PPLParserBIN, 0)
}

func (s *BinCommandContext) FieldReference() IFieldReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldReferenceContext)
}

func (s *BinCommandContext) AllBinOptions() []IBinOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBinOptionsContext); ok {
			len++
		}
	}

	tst := make([]IBinOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBinOptionsContext); ok {
			tst[i] = t.(IBinOptionsContext)
			i++
		}
	}

	return tst
}

func (s *BinCommandContext) BinOptions(i int) IBinOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBinOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBinOptionsContext)
}

func (s *BinCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BinCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterBinCommand(s)
	}
}

func (s *BinCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitBinCommand(s)
	}
}

func (s *BinCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitBinCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) BinCommand() (localctx IBinCommandContext) {
	localctx = NewBinCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, PPLParserRULE_binCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(312)
		p.Match(PPLParserBIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(313)
		p.FieldReference()
	}
	p.SetState(317)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserSPAN || _la == PPLParserBINS {
		{
			p.SetState(314)
			p.BinOptions()
		}

		p.SetState(319)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBinOptionsContext is an interface to support dynamic dispatch.
type IBinOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SPAN() antlr.TerminalNode
	EQ() antlr.TerminalNode
	TimeSpan() ITimeSpanContext
	IDENTIFIER() antlr.TerminalNode
	BINS() antlr.TerminalNode
	INTEGER() antlr.TerminalNode

	// IsBinOptionsContext differentiates from other interfaces.
	IsBinOptionsContext()
}

type BinOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinOptionsContext() *BinOptionsContext {
	var p = new(BinOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_binOptions
	return p
}

func InitEmptyBinOptionsContext(p *BinOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_binOptions
}

func (*BinOptionsContext) IsBinOptionsContext() {}

func NewBinOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinOptionsContext {
	var p = new(BinOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_binOptions

	return p
}

func (s *BinOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *BinOptionsContext) SPAN() antlr.TerminalNode {
	return s.GetToken(PPLParserSPAN, 0)
}

func (s *BinOptionsContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *BinOptionsContext) TimeSpan() ITimeSpanContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeSpanContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeSpanContext)
}

func (s *BinOptionsContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *BinOptionsContext) BINS() antlr.TerminalNode {
	return s.GetToken(PPLParserBINS, 0)
}

func (s *BinOptionsContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *BinOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BinOptionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterBinOptions(s)
	}
}

func (s *BinOptionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitBinOptions(s)
	}
}

func (s *BinOptionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitBinOptions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) BinOptions() (localctx IBinOptionsContext) {
	localctx = NewBinOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, PPLParserRULE_binOptions)
	p.SetState(329)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserSPAN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(320)
			p.Match(PPLParserSPAN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(321)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(324)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case PPLParserINTEGER:
			{
				p.SetState(322)
				p.TimeSpan()
			}

		case PPLParserIDENTIFIER:
			{
				p.SetState(323)
				p.Match(PPLParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

	case PPLParserBINS:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(326)
			p.Match(PPLParserBINS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(327)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(328)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDedupCommandContext is an interface to support dynamic dispatch.
type IDedupCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEDUP() antlr.TerminalNode
	FieldList() IFieldListContext
	INTEGER() antlr.TerminalNode
	AllDedupOptions() []IDedupOptionsContext
	DedupOptions(i int) IDedupOptionsContext

	// IsDedupCommandContext differentiates from other interfaces.
	IsDedupCommandContext()
}

type DedupCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDedupCommandContext() *DedupCommandContext {
	var p = new(DedupCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_dedupCommand
	return p
}

func InitEmptyDedupCommandContext(p *DedupCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_dedupCommand
}

func (*DedupCommandContext) IsDedupCommandContext() {}

func NewDedupCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DedupCommandContext {
	var p = new(DedupCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_dedupCommand

	return p
}

func (s *DedupCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *DedupCommandContext) DEDUP() antlr.TerminalNode {
	return s.GetToken(PPLParserDEDUP, 0)
}

func (s *DedupCommandContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *DedupCommandContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *DedupCommandContext) AllDedupOptions() []IDedupOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDedupOptionsContext); ok {
			len++
		}
	}

	tst := make([]IDedupOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDedupOptionsContext); ok {
			tst[i] = t.(IDedupOptionsContext)
			i++
		}
	}

	return tst
}

func (s *DedupCommandContext) DedupOptions(i int) IDedupOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDedupOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDedupOptionsContext)
}

func (s *DedupCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DedupCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DedupCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterDedupCommand(s)
	}
}

func (s *DedupCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitDedupCommand(s)
	}
}

func (s *DedupCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitDedupCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) DedupCommand() (localctx IDedupCommandContext) {
	localctx = NewDedupCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, PPLParserRULE_dedupCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(331)
		p.Match(PPLParserDEDUP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(333)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(332)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(335)
		p.FieldList()
	}
	p.SetState(339)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&61572651155456) != 0 {
		{
			p.SetState(336)
			p.DedupOptions()
		}

		p.SetState(341)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDedupOptionsContext is an interface to support dynamic dispatch.
type IDedupOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KEEPEVENTS() antlr.TerminalNode
	EQ() antlr.TerminalNode
	BooleanValue() IBooleanValueContext
	CONSECUTIVE() antlr.TerminalNode
	SORTBY() antlr.TerminalNode
	SortFieldList() ISortFieldListContext

	// IsDedupOptionsContext differentiates from other interfaces.
	IsDedupOptionsContext()
}

type DedupOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDedupOptionsContext() *DedupOptionsContext {
	var p = new(DedupOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_dedupOptions
	return p
}

func InitEmptyDedupOptionsContext(p *DedupOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_dedupOptions
}

func (*DedupOptionsContext) IsDedupOptionsContext() {}

func NewDedupOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DedupOptionsContext {
	var p = new(DedupOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_dedupOptions

	return p
}

func (s *DedupOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *DedupOptionsContext) KEEPEVENTS() antlr.TerminalNode {
	return s.GetToken(PPLParserKEEPEVENTS, 0)
}

func (s *DedupOptionsContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *DedupOptionsContext) BooleanValue() IBooleanValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanValueContext)
}

func (s *DedupOptionsContext) CONSECUTIVE() antlr.TerminalNode {
	return s.GetToken(PPLParserCONSECUTIVE, 0)
}

func (s *DedupOptionsContext) SORTBY() antlr.TerminalNode {
	return s.GetToken(PPLParserSORTBY, 0)
}

func (s *DedupOptionsContext) SortFieldList() ISortFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortFieldListContext)
}

func (s *DedupOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DedupOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DedupOptionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterDedupOptions(s)
	}
}

func (s *DedupOptionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitDedupOptions(s)
	}
}

func (s *DedupOptionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitDedupOptions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) DedupOptions() (localctx IDedupOptionsContext) {
	localctx = NewDedupOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, PPLParserRULE_dedupOptions)
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserKEEPEVENTS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(342)
			p.Match(PPLParserKEEPEVENTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(343)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(344)
			p.BooleanValue()
		}

	case PPLParserCONSECUTIVE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(345)
			p.Match(PPLParserCONSECUTIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(346)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(347)
			p.BooleanValue()
		}

	case PPLParserSORTBY:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(348)
			p.Match(PPLParserSORTBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(349)
			p.SortFieldList()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITopCommandContext is an interface to support dynamic dispatch.
type ITopCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TOP() antlr.TerminalNode
	AllFieldList() []IFieldListContext
	FieldList(i int) IFieldListContext
	INTEGER() antlr.TerminalNode
	BY() antlr.TerminalNode
	AllTopOptions() []ITopOptionsContext
	TopOptions(i int) ITopOptionsContext

	// IsTopCommandContext differentiates from other interfaces.
	IsTopCommandContext()
}

type TopCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTopCommandContext() *TopCommandContext {
	var p = new(TopCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_topCommand
	return p
}

func InitEmptyTopCommandContext(p *TopCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_topCommand
}

func (*TopCommandContext) IsTopCommandContext() {}

func NewTopCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopCommandContext {
	var p = new(TopCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_topCommand

	return p
}

func (s *TopCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *TopCommandContext) TOP() antlr.TerminalNode {
	return s.GetToken(PPLParserTOP, 0)
}

func (s *TopCommandContext) AllFieldList() []IFieldListContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldListContext); ok {
			len++
		}
	}

	tst := make([]IFieldListContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldListContext); ok {
			tst[i] = t.(IFieldListContext)
			i++
		}
	}

	return tst
}

func (s *TopCommandContext) FieldList(i int) IFieldListContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *TopCommandContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *TopCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(PPLParserBY, 0)
}

func (s *TopCommandContext) AllTopOptions() []ITopOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITopOptionsContext); ok {
			len++
		}
	}

	tst := make([]ITopOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITopOptionsContext); ok {
			tst[i] = t.(ITopOptionsContext)
			i++
		}
	}

	return tst
}

func (s *TopCommandContext) TopOptions(i int) ITopOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITopOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITopOptionsContext)
}

func (s *TopCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TopCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TopCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterTopCommand(s)
	}
}

func (s *TopCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitTopCommand(s)
	}
}

func (s *TopCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitTopCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) TopCommand() (localctx ITopCommandContext) {
	localctx = NewTopCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, PPLParserRULE_topCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(352)
		p.Match(PPLParserTOP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(354)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(353)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(356)
		p.FieldList()
	}
	p.SetState(359)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserBY {
		{
			p.SetState(357)
			p.Match(PPLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(358)
			p.FieldList()
		}

	}
	p.SetState(364)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&8936830510563328) != 0 {
		{
			p.SetState(361)
			p.TopOptions()
		}

		p.SetState(366)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITopOptionsContext is an interface to support dynamic dispatch.
type ITopOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COUNTFIELD() antlr.TerminalNode
	EQ() antlr.TerminalNode
	STRING() antlr.TerminalNode
	PERCENTFIELD() antlr.TerminalNode
	SHOWCOUNT() antlr.TerminalNode
	BooleanValue() IBooleanValueContext
	SHOWPERC() antlr.TerminalNode
	LIMIT() antlr.TerminalNode
	INTEGER() antlr.TerminalNode
	USEOTHER() antlr.TerminalNode
	OTHERSTR() antlr.TerminalNode

	// IsTopOptionsContext differentiates from other interfaces.
	IsTopOptionsContext()
}

type TopOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTopOptionsContext() *TopOptionsContext {
	var p = new(TopOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_topOptions
	return p
}

func InitEmptyTopOptionsContext(p *TopOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_topOptions
}

func (*TopOptionsContext) IsTopOptionsContext() {}

func NewTopOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopOptionsContext {
	var p = new(TopOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_topOptions

	return p
}

func (s *TopOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *TopOptionsContext) COUNTFIELD() antlr.TerminalNode {
	return s.GetToken(PPLParserCOUNTFIELD, 0)
}

func (s *TopOptionsContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *TopOptionsContext) STRING() antlr.TerminalNode {
	return s.GetToken(PPLParserSTRING, 0)
}

func (s *TopOptionsContext) PERCENTFIELD() antlr.TerminalNode {
	return s.GetToken(PPLParserPERCENTFIELD, 0)
}

func (s *TopOptionsContext) SHOWCOUNT() antlr.TerminalNode {
	return s.GetToken(PPLParserSHOWCOUNT, 0)
}

func (s *TopOptionsContext) BooleanValue() IBooleanValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanValueContext)
}

func (s *TopOptionsContext) SHOWPERC() antlr.TerminalNode {
	return s.GetToken(PPLParserSHOWPERC, 0)
}

func (s *TopOptionsContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(PPLParserLIMIT, 0)
}

func (s *TopOptionsContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *TopOptionsContext) USEOTHER() antlr.TerminalNode {
	return s.GetToken(PPLParserUSEOTHER, 0)
}

func (s *TopOptionsContext) OTHERSTR() antlr.TerminalNode {
	return s.GetToken(PPLParserOTHERSTR, 0)
}

func (s *TopOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TopOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TopOptionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterTopOptions(s)
	}
}

func (s *TopOptionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitTopOptions(s)
	}
}

func (s *TopOptionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitTopOptions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) TopOptions() (localctx ITopOptionsContext) {
	localctx = NewTopOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, PPLParserRULE_topOptions)
	p.SetState(388)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserCOUNTFIELD:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(367)
			p.Match(PPLParserCOUNTFIELD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(368)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(369)
			p.Match(PPLParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserPERCENTFIELD:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(370)
			p.Match(PPLParserPERCENTFIELD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(371)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(372)
			p.Match(PPLParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserSHOWCOUNT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(373)
			p.Match(PPLParserSHOWCOUNT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(374)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(375)
			p.BooleanValue()
		}

	case PPLParserSHOWPERC:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(376)
			p.Match(PPLParserSHOWPERC)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(377)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(378)
			p.BooleanValue()
		}

	case PPLParserLIMIT:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(379)
			p.Match(PPLParserLIMIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(380)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(381)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserUSEOTHER:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(382)
			p.Match(PPLParserUSEOTHER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(383)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(384)
			p.BooleanValue()
		}

	case PPLParserOTHERSTR:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(385)
			p.Match(PPLParserOTHERSTR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(386)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(387)
			p.Match(PPLParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRareCommandContext is an interface to support dynamic dispatch.
type IRareCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RARE() antlr.TerminalNode
	AllFieldList() []IFieldListContext
	FieldList(i int) IFieldListContext
	INTEGER() antlr.TerminalNode
	BY() antlr.TerminalNode
	AllTopOptions() []ITopOptionsContext
	TopOptions(i int) ITopOptionsContext

	// IsRareCommandContext differentiates from other interfaces.
	IsRareCommandContext()
}

type RareCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRareCommandContext() *RareCommandContext {
	var p = new(RareCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_rareCommand
	return p
}

func InitEmptyRareCommandContext(p *RareCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_rareCommand
}

func (*RareCommandContext) IsRareCommandContext() {}

func NewRareCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RareCommandContext {
	var p = new(RareCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_rareCommand

	return p
}

func (s *RareCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *RareCommandContext) RARE() antlr.TerminalNode {
	return s.GetToken(PPLParserRARE, 0)
}

func (s *RareCommandContext) AllFieldList() []IFieldListContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldListContext); ok {
			len++
		}
	}

	tst := make([]IFieldListContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldListContext); ok {
			tst[i] = t.(IFieldListContext)
			i++
		}
	}

	return tst
}

func (s *RareCommandContext) FieldList(i int) IFieldListContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *RareCommandContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *RareCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(PPLParserBY, 0)
}

func (s *RareCommandContext) AllTopOptions() []ITopOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITopOptionsContext); ok {
			len++
		}
	}

	tst := make([]ITopOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITopOptionsContext); ok {
			tst[i] = t.(ITopOptionsContext)
			i++
		}
	}

	return tst
}

func (s *RareCommandContext) TopOptions(i int) ITopOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITopOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITopOptionsContext)
}

func (s *RareCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RareCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RareCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterRareCommand(s)
	}
}

func (s *RareCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitRareCommand(s)
	}
}

func (s *RareCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitRareCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) RareCommand() (localctx IRareCommandContext) {
	localctx = NewRareCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, PPLParserRULE_rareCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(390)
		p.Match(PPLParserRARE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(392)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 28, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(391)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(394)
		p.FieldList()
	}
	p.SetState(397)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserBY {
		{
			p.SetState(395)
			p.Match(PPLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(396)
			p.FieldList()
		}

	}
	p.SetState(402)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&8936830510563328) != 0 {
		{
			p.SetState(399)
			p.TopOptions()
		}

		p.SetState(404)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEvalCommandContext is an interface to support dynamic dispatch.
type IEvalCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EVAL() antlr.TerminalNode
	AllEvalAssignment() []IEvalAssignmentContext
	EvalAssignment(i int) IEvalAssignmentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsEvalCommandContext differentiates from other interfaces.
	IsEvalCommandContext()
}

type EvalCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEvalCommandContext() *EvalCommandContext {
	var p = new(EvalCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_evalCommand
	return p
}

func InitEmptyEvalCommandContext(p *EvalCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_evalCommand
}

func (*EvalCommandContext) IsEvalCommandContext() {}

func NewEvalCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EvalCommandContext {
	var p = new(EvalCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_evalCommand

	return p
}

func (s *EvalCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *EvalCommandContext) EVAL() antlr.TerminalNode {
	return s.GetToken(PPLParserEVAL, 0)
}

func (s *EvalCommandContext) AllEvalAssignment() []IEvalAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IEvalAssignmentContext); ok {
			len++
		}
	}

	tst := make([]IEvalAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IEvalAssignmentContext); ok {
			tst[i] = t.(IEvalAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *EvalCommandContext) EvalAssignment(i int) IEvalAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEvalAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEvalAssignmentContext)
}

func (s *EvalCommandContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *EvalCommandContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *EvalCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EvalCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EvalCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterEvalCommand(s)
	}
}

func (s *EvalCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitEvalCommand(s)
	}
}

func (s *EvalCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitEvalCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) EvalCommand() (localctx IEvalCommandContext) {
	localctx = NewEvalCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, PPLParserRULE_evalCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(405)
		p.Match(PPLParserEVAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(406)
		p.EvalAssignment()
	}
	p.SetState(411)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(407)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(408)
			p.EvalAssignment()
		}

		p.SetState(413)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEvalAssignmentContext is an interface to support dynamic dispatch.
type IEvalAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	EQ() antlr.TerminalNode
	Expression() IExpressionContext

	// IsEvalAssignmentContext differentiates from other interfaces.
	IsEvalAssignmentContext()
}

type EvalAssignmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEvalAssignmentContext() *EvalAssignmentContext {
	var p = new(EvalAssignmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_evalAssignment
	return p
}

func InitEmptyEvalAssignmentContext(p *EvalAssignmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_evalAssignment
}

func (*EvalAssignmentContext) IsEvalAssignmentContext() {}

func NewEvalAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EvalAssignmentContext {
	var p = new(EvalAssignmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_evalAssignment

	return p
}

func (s *EvalAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *EvalAssignmentContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *EvalAssignmentContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *EvalAssignmentContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *EvalAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EvalAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EvalAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterEvalAssignment(s)
	}
}

func (s *EvalAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitEvalAssignment(s)
	}
}

func (s *EvalAssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitEvalAssignment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) EvalAssignment() (localctx IEvalAssignmentContext) {
	localctx = NewEvalAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, PPLParserRULE_evalAssignment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(414)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(415)
		p.Match(PPLParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(416)
		p.Expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRenameCommandContext is an interface to support dynamic dispatch.
type IRenameCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RENAME() antlr.TerminalNode
	AllRenameAssignment() []IRenameAssignmentContext
	RenameAssignment(i int) IRenameAssignmentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsRenameCommandContext differentiates from other interfaces.
	IsRenameCommandContext()
}

type RenameCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRenameCommandContext() *RenameCommandContext {
	var p = new(RenameCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_renameCommand
	return p
}

func InitEmptyRenameCommandContext(p *RenameCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_renameCommand
}

func (*RenameCommandContext) IsRenameCommandContext() {}

func NewRenameCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RenameCommandContext {
	var p = new(RenameCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_renameCommand

	return p
}

func (s *RenameCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *RenameCommandContext) RENAME() antlr.TerminalNode {
	return s.GetToken(PPLParserRENAME, 0)
}

func (s *RenameCommandContext) AllRenameAssignment() []IRenameAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRenameAssignmentContext); ok {
			len++
		}
	}

	tst := make([]IRenameAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRenameAssignmentContext); ok {
			tst[i] = t.(IRenameAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *RenameCommandContext) RenameAssignment(i int) IRenameAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRenameAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRenameAssignmentContext)
}

func (s *RenameCommandContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *RenameCommandContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *RenameCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RenameCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RenameCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterRenameCommand(s)
	}
}

func (s *RenameCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitRenameCommand(s)
	}
}

func (s *RenameCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitRenameCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) RenameCommand() (localctx IRenameCommandContext) {
	localctx = NewRenameCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, PPLParserRULE_renameCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(418)
		p.Match(PPLParserRENAME)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(419)
		p.RenameAssignment()
	}
	p.SetState(424)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(420)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(421)
			p.RenameAssignment()
		}

		p.SetState(426)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRenameAssignmentContext is an interface to support dynamic dispatch.
type IRenameAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AS() antlr.TerminalNode

	// IsRenameAssignmentContext differentiates from other interfaces.
	IsRenameAssignmentContext()
}

type RenameAssignmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRenameAssignmentContext() *RenameAssignmentContext {
	var p = new(RenameAssignmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_renameAssignment
	return p
}

func InitEmptyRenameAssignmentContext(p *RenameAssignmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_renameAssignment
}

func (*RenameAssignmentContext) IsRenameAssignmentContext() {}

func NewRenameAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RenameAssignmentContext {
	var p = new(RenameAssignmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_renameAssignment

	return p
}

func (s *RenameAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *RenameAssignmentContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(PPLParserIDENTIFIER)
}

func (s *RenameAssignmentContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, i)
}

func (s *RenameAssignmentContext) AS() antlr.TerminalNode {
	return s.GetToken(PPLParserAS, 0)
}

func (s *RenameAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RenameAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RenameAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterRenameAssignment(s)
	}
}

func (s *RenameAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitRenameAssignment(s)
	}
}

func (s *RenameAssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitRenameAssignment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) RenameAssignment() (localctx IRenameAssignmentContext) {
	localctx = NewRenameAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, PPLParserRULE_renameAssignment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(427)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(428)
		p.Match(PPLParserAS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(429)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IReplaceCommandContext is an interface to support dynamic dispatch.
type IReplaceCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REPLACE() antlr.TerminalNode
	AllReplaceMapping() []IReplaceMappingContext
	ReplaceMapping(i int) IReplaceMappingContext
	IN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsReplaceCommandContext differentiates from other interfaces.
	IsReplaceCommandContext()
}

type ReplaceCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReplaceCommandContext() *ReplaceCommandContext {
	var p = new(ReplaceCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_replaceCommand
	return p
}

func InitEmptyReplaceCommandContext(p *ReplaceCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_replaceCommand
}

func (*ReplaceCommandContext) IsReplaceCommandContext() {}

func NewReplaceCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReplaceCommandContext {
	var p = new(ReplaceCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_replaceCommand

	return p
}

func (s *ReplaceCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ReplaceCommandContext) REPLACE() antlr.TerminalNode {
	return s.GetToken(PPLParserREPLACE, 0)
}

func (s *ReplaceCommandContext) AllReplaceMapping() []IReplaceMappingContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IReplaceMappingContext); ok {
			len++
		}
	}

	tst := make([]IReplaceMappingContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IReplaceMappingContext); ok {
			tst[i] = t.(IReplaceMappingContext)
			i++
		}
	}

	return tst
}

func (s *ReplaceCommandContext) ReplaceMapping(i int) IReplaceMappingContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IReplaceMappingContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IReplaceMappingContext)
}

func (s *ReplaceCommandContext) IN() antlr.TerminalNode {
	return s.GetToken(PPLParserIN, 0)
}

func (s *ReplaceCommandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *ReplaceCommandContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *ReplaceCommandContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *ReplaceCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReplaceCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReplaceCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterReplaceCommand(s)
	}
}

func (s *ReplaceCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitReplaceCommand(s)
	}
}

func (s *ReplaceCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitReplaceCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ReplaceCommand() (localctx IReplaceCommandContext) {
	localctx = NewReplaceCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, PPLParserRULE_replaceCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(431)
		p.Match(PPLParserREPLACE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(432)
		p.ReplaceMapping()
	}
	p.SetState(437)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(433)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(434)
			p.ReplaceMapping()
		}

		p.SetState(439)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(440)
		p.Match(PPLParserIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(441)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IReplaceMappingContext is an interface to support dynamic dispatch.
type IReplaceMappingContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	WITH() antlr.TerminalNode

	// IsReplaceMappingContext differentiates from other interfaces.
	IsReplaceMappingContext()
}

type ReplaceMappingContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReplaceMappingContext() *ReplaceMappingContext {
	var p = new(ReplaceMappingContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_replaceMapping
	return p
}

func InitEmptyReplaceMappingContext(p *ReplaceMappingContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_replaceMapping
}

func (*ReplaceMappingContext) IsReplaceMappingContext() {}

func NewReplaceMappingContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReplaceMappingContext {
	var p = new(ReplaceMappingContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_replaceMapping

	return p
}

func (s *ReplaceMappingContext) GetParser() antlr.Parser { return s.parser }

func (s *ReplaceMappingContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ReplaceMappingContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ReplaceMappingContext) WITH() antlr.TerminalNode {
	return s.GetToken(PPLParserWITH, 0)
}

func (s *ReplaceMappingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReplaceMappingContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReplaceMappingContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterReplaceMapping(s)
	}
}

func (s *ReplaceMappingContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitReplaceMapping(s)
	}
}

func (s *ReplaceMappingContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitReplaceMapping(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ReplaceMapping() (localctx IReplaceMappingContext) {
	localctx = NewReplaceMappingContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, PPLParserRULE_replaceMapping)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(443)
		p.Expression()
	}
	{
		p.SetState(444)
		p.Match(PPLParserWITH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(445)
		p.Expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFillnullCommandContext is an interface to support dynamic dispatch.
type IFillnullCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsFillnullCommandContext differentiates from other interfaces.
	IsFillnullCommandContext()
}

type FillnullCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillnullCommandContext() *FillnullCommandContext {
	var p = new(FillnullCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fillnullCommand
	return p
}

func InitEmptyFillnullCommandContext(p *FillnullCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fillnullCommand
}

func (*FillnullCommandContext) IsFillnullCommandContext() {}

func NewFillnullCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillnullCommandContext {
	var p = new(FillnullCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_fillnullCommand

	return p
}

func (s *FillnullCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *FillnullCommandContext) CopyAll(ctx *FillnullCommandContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *FillnullCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillnullCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type FillnullWithAssignmentsContext struct {
	FillnullCommandContext
}

func NewFillnullWithAssignmentsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FillnullWithAssignmentsContext {
	var p = new(FillnullWithAssignmentsContext)

	InitEmptyFillnullCommandContext(&p.FillnullCommandContext)
	p.parser = parser
	p.CopyAll(ctx.(*FillnullCommandContext))

	return p
}

func (s *FillnullWithAssignmentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillnullWithAssignmentsContext) FILLNULL() antlr.TerminalNode {
	return s.GetToken(PPLParserFILLNULL, 0)
}

func (s *FillnullWithAssignmentsContext) AllFillnullAssignment() []IFillnullAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFillnullAssignmentContext); ok {
			len++
		}
	}

	tst := make([]IFillnullAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFillnullAssignmentContext); ok {
			tst[i] = t.(IFillnullAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *FillnullWithAssignmentsContext) FillnullAssignment(i int) IFillnullAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFillnullAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFillnullAssignmentContext)
}

func (s *FillnullWithAssignmentsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *FillnullWithAssignmentsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *FillnullWithAssignmentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFillnullWithAssignments(s)
	}
}

func (s *FillnullWithAssignmentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFillnullWithAssignments(s)
	}
}

func (s *FillnullWithAssignmentsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFillnullWithAssignments(s)

	default:
		return t.VisitChildren(s)
	}
}

type FillnullWithDefaultContext struct {
	FillnullCommandContext
}

func NewFillnullWithDefaultContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FillnullWithDefaultContext {
	var p = new(FillnullWithDefaultContext)

	InitEmptyFillnullCommandContext(&p.FillnullCommandContext)
	p.parser = parser
	p.CopyAll(ctx.(*FillnullCommandContext))

	return p
}

func (s *FillnullWithDefaultContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillnullWithDefaultContext) FILLNULL() antlr.TerminalNode {
	return s.GetToken(PPLParserFILLNULL, 0)
}

func (s *FillnullWithDefaultContext) VALUE() antlr.TerminalNode {
	return s.GetToken(PPLParserVALUE, 0)
}

func (s *FillnullWithDefaultContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *FillnullWithDefaultContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FillnullWithDefaultContext) FIELDS() antlr.TerminalNode {
	return s.GetToken(PPLParserFIELDS, 0)
}

func (s *FillnullWithDefaultContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *FillnullWithDefaultContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFillnullWithDefault(s)
	}
}

func (s *FillnullWithDefaultContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFillnullWithDefault(s)
	}
}

func (s *FillnullWithDefaultContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFillnullWithDefault(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) FillnullCommand() (localctx IFillnullCommandContext) {
	localctx = NewFillnullCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, PPLParserRULE_fillnullCommand)
	var _la int

	p.SetState(464)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 36, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFillnullWithDefaultContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(447)
			p.Match(PPLParserFILLNULL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(448)
			p.Match(PPLParserVALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(449)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(450)
			p.Expression()
		}
		p.SetState(453)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == PPLParserFIELDS {
			{
				p.SetState(451)
				p.Match(PPLParserFIELDS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(452)
				p.FieldList()
			}

		}

	case 2:
		localctx = NewFillnullWithAssignmentsContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(455)
			p.Match(PPLParserFILLNULL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(456)
			p.FillnullAssignment()
		}
		p.SetState(461)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == PPLParserCOMMA {
			{
				p.SetState(457)
				p.Match(PPLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(458)
				p.FillnullAssignment()
			}

			p.SetState(463)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFillnullAssignmentContext is an interface to support dynamic dispatch.
type IFillnullAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	EQ() antlr.TerminalNode
	Expression() IExpressionContext

	// IsFillnullAssignmentContext differentiates from other interfaces.
	IsFillnullAssignmentContext()
}

type FillnullAssignmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillnullAssignmentContext() *FillnullAssignmentContext {
	var p = new(FillnullAssignmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fillnullAssignment
	return p
}

func InitEmptyFillnullAssignmentContext(p *FillnullAssignmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fillnullAssignment
}

func (*FillnullAssignmentContext) IsFillnullAssignmentContext() {}

func NewFillnullAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillnullAssignmentContext {
	var p = new(FillnullAssignmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_fillnullAssignment

	return p
}

func (s *FillnullAssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *FillnullAssignmentContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *FillnullAssignmentContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *FillnullAssignmentContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FillnullAssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillnullAssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FillnullAssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFillnullAssignment(s)
	}
}

func (s *FillnullAssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFillnullAssignment(s)
	}
}

func (s *FillnullAssignmentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFillnullAssignment(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) FillnullAssignment() (localctx IFillnullAssignmentContext) {
	localctx = NewFillnullAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, PPLParserRULE_fillnullAssignment)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(466)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(467)
		p.Match(PPLParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(468)
		p.Expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParseCommandContext is an interface to support dynamic dispatch.
type IParseCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PARSE() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	STRING() antlr.TerminalNode
	EQ() antlr.TerminalNode

	// IsParseCommandContext differentiates from other interfaces.
	IsParseCommandContext()
}

type ParseCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParseCommandContext() *ParseCommandContext {
	var p = new(ParseCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_parseCommand
	return p
}

func InitEmptyParseCommandContext(p *ParseCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_parseCommand
}

func (*ParseCommandContext) IsParseCommandContext() {}

func NewParseCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseCommandContext {
	var p = new(ParseCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_parseCommand

	return p
}

func (s *ParseCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ParseCommandContext) PARSE() antlr.TerminalNode {
	return s.GetToken(PPLParserPARSE, 0)
}

func (s *ParseCommandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(PPLParserIDENTIFIER)
}

func (s *ParseCommandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, i)
}

func (s *ParseCommandContext) STRING() antlr.TerminalNode {
	return s.GetToken(PPLParserSTRING, 0)
}

func (s *ParseCommandContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *ParseCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParseCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParseCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterParseCommand(s)
	}
}

func (s *ParseCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitParseCommand(s)
	}
}

func (s *ParseCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitParseCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ParseCommand() (localctx IParseCommandContext) {
	localctx = NewParseCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, PPLParserRULE_parseCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(470)
		p.Match(PPLParserPARSE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(473)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(471)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(472)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(475)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(476)
		p.Match(PPLParserSTRING)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRexCommandContext is an interface to support dynamic dispatch.
type IRexCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REX() antlr.TerminalNode
	STRING() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	EQ() antlr.TerminalNode

	// IsRexCommandContext differentiates from other interfaces.
	IsRexCommandContext()
}

type RexCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRexCommandContext() *RexCommandContext {
	var p = new(RexCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_rexCommand
	return p
}

func InitEmptyRexCommandContext(p *RexCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_rexCommand
}

func (*RexCommandContext) IsRexCommandContext() {}

func NewRexCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RexCommandContext {
	var p = new(RexCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_rexCommand

	return p
}

func (s *RexCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *RexCommandContext) REX() antlr.TerminalNode {
	return s.GetToken(PPLParserREX, 0)
}

func (s *RexCommandContext) STRING() antlr.TerminalNode {
	return s.GetToken(PPLParserSTRING, 0)
}

func (s *RexCommandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(PPLParserIDENTIFIER)
}

func (s *RexCommandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, i)
}

func (s *RexCommandContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *RexCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RexCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RexCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterRexCommand(s)
	}
}

func (s *RexCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitRexCommand(s)
	}
}

func (s *RexCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitRexCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) RexCommand() (localctx IRexCommandContext) {
	localctx = NewRexCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, PPLParserRULE_rexCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(478)
		p.Match(PPLParserREX)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(482)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserIDENTIFIER {
		{
			p.SetState(479)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(480)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(481)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(484)
		p.Match(PPLParserSTRING)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILookupCommandContext is an interface to support dynamic dispatch.
type ILookupCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LOOKUP() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	OUTPUT() antlr.TerminalNode
	LookupOutputList() ILookupOutputListContext
	AS() antlr.TerminalNode

	// IsLookupCommandContext differentiates from other interfaces.
	IsLookupCommandContext()
}

type LookupCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLookupCommandContext() *LookupCommandContext {
	var p = new(LookupCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_lookupCommand
	return p
}

func InitEmptyLookupCommandContext(p *LookupCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_lookupCommand
}

func (*LookupCommandContext) IsLookupCommandContext() {}

func NewLookupCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LookupCommandContext {
	var p = new(LookupCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_lookupCommand

	return p
}

func (s *LookupCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *LookupCommandContext) LOOKUP() antlr.TerminalNode {
	return s.GetToken(PPLParserLOOKUP, 0)
}

func (s *LookupCommandContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(PPLParserIDENTIFIER)
}

func (s *LookupCommandContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, i)
}

func (s *LookupCommandContext) OUTPUT() antlr.TerminalNode {
	return s.GetToken(PPLParserOUTPUT, 0)
}

func (s *LookupCommandContext) LookupOutputList() ILookupOutputListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILookupOutputListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILookupOutputListContext)
}

func (s *LookupCommandContext) AS() antlr.TerminalNode {
	return s.GetToken(PPLParserAS, 0)
}

func (s *LookupCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LookupCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LookupCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterLookupCommand(s)
	}
}

func (s *LookupCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitLookupCommand(s)
	}
}

func (s *LookupCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitLookupCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) LookupCommand() (localctx ILookupCommandContext) {
	localctx = NewLookupCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, PPLParserRULE_lookupCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(486)
		p.Match(PPLParserLOOKUP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(487)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(488)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(491)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserAS {
		{
			p.SetState(489)
			p.Match(PPLParserAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(490)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(493)
		p.Match(PPLParserOUTPUT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(494)
		p.LookupOutputList()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILookupOutputListContext is an interface to support dynamic dispatch.
type ILookupOutputListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLookupOutputField() []ILookupOutputFieldContext
	LookupOutputField(i int) ILookupOutputFieldContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsLookupOutputListContext differentiates from other interfaces.
	IsLookupOutputListContext()
}

type LookupOutputListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLookupOutputListContext() *LookupOutputListContext {
	var p = new(LookupOutputListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_lookupOutputList
	return p
}

func InitEmptyLookupOutputListContext(p *LookupOutputListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_lookupOutputList
}

func (*LookupOutputListContext) IsLookupOutputListContext() {}

func NewLookupOutputListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LookupOutputListContext {
	var p = new(LookupOutputListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_lookupOutputList

	return p
}

func (s *LookupOutputListContext) GetParser() antlr.Parser { return s.parser }

func (s *LookupOutputListContext) AllLookupOutputField() []ILookupOutputFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILookupOutputFieldContext); ok {
			len++
		}
	}

	tst := make([]ILookupOutputFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILookupOutputFieldContext); ok {
			tst[i] = t.(ILookupOutputFieldContext)
			i++
		}
	}

	return tst
}

func (s *LookupOutputListContext) LookupOutputField(i int) ILookupOutputFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILookupOutputFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILookupOutputFieldContext)
}

func (s *LookupOutputListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *LookupOutputListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *LookupOutputListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LookupOutputListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LookupOutputListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterLookupOutputList(s)
	}
}

func (s *LookupOutputListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitLookupOutputList(s)
	}
}

func (s *LookupOutputListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitLookupOutputList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) LookupOutputList() (localctx ILookupOutputListContext) {
	localctx = NewLookupOutputListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, PPLParserRULE_lookupOutputList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(496)
		p.LookupOutputField()
	}
	p.SetState(501)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(497)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(498)
			p.LookupOutputField()
		}

		p.SetState(503)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILookupOutputFieldContext is an interface to support dynamic dispatch.
type ILookupOutputFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AS() antlr.TerminalNode

	// IsLookupOutputFieldContext differentiates from other interfaces.
	IsLookupOutputFieldContext()
}

type LookupOutputFieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLookupOutputFieldContext() *LookupOutputFieldContext {
	var p = new(LookupOutputFieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_lookupOutputField
	return p
}

func InitEmptyLookupOutputFieldContext(p *LookupOutputFieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_lookupOutputField
}

func (*LookupOutputFieldContext) IsLookupOutputFieldContext() {}

func NewLookupOutputFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LookupOutputFieldContext {
	var p = new(LookupOutputFieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_lookupOutputField

	return p
}

func (s *LookupOutputFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *LookupOutputFieldContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(PPLParserIDENTIFIER)
}

func (s *LookupOutputFieldContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, i)
}

func (s *LookupOutputFieldContext) AS() antlr.TerminalNode {
	return s.GetToken(PPLParserAS, 0)
}

func (s *LookupOutputFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LookupOutputFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LookupOutputFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterLookupOutputField(s)
	}
}

func (s *LookupOutputFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitLookupOutputField(s)
	}
}

func (s *LookupOutputFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitLookupOutputField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) LookupOutputField() (localctx ILookupOutputFieldContext) {
	localctx = NewLookupOutputFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, PPLParserRULE_lookupOutputField)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(504)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(507)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserAS {
		{
			p.SetState(505)
			p.Match(PPLParserAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(506)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAppendCommandContext is an interface to support dynamic dispatch.
type IAppendCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	APPEND() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	SearchQuery() ISearchQueryContext
	RBRACKET() antlr.TerminalNode

	// IsAppendCommandContext differentiates from other interfaces.
	IsAppendCommandContext()
}

type AppendCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAppendCommandContext() *AppendCommandContext {
	var p = new(AppendCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_appendCommand
	return p
}

func InitEmptyAppendCommandContext(p *AppendCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_appendCommand
}

func (*AppendCommandContext) IsAppendCommandContext() {}

func NewAppendCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AppendCommandContext {
	var p = new(AppendCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_appendCommand

	return p
}

func (s *AppendCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *AppendCommandContext) APPEND() antlr.TerminalNode {
	return s.GetToken(PPLParserAPPEND, 0)
}

func (s *AppendCommandContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(PPLParserLBRACKET, 0)
}

func (s *AppendCommandContext) SearchQuery() ISearchQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISearchQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISearchQueryContext)
}

func (s *AppendCommandContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PPLParserRBRACKET, 0)
}

func (s *AppendCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AppendCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AppendCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAppendCommand(s)
	}
}

func (s *AppendCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAppendCommand(s)
	}
}

func (s *AppendCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAppendCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) AppendCommand() (localctx IAppendCommandContext) {
	localctx = NewAppendCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, PPLParserRULE_appendCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(509)
		p.Match(PPLParserAPPEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(510)
		p.Match(PPLParserLBRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(511)
		p.SearchQuery()
	}
	{
		p.SetState(512)
		p.Match(PPLParserRBRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IJoinCommandContext is an interface to support dynamic dispatch.
type IJoinCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	JOIN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	SearchQuery() ISearchQueryContext
	RBRACKET() antlr.TerminalNode
	TYPE() antlr.TerminalNode
	EQ() antlr.TerminalNode
	JoinType() IJoinTypeContext

	// IsJoinCommandContext differentiates from other interfaces.
	IsJoinCommandContext()
}

type JoinCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJoinCommandContext() *JoinCommandContext {
	var p = new(JoinCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_joinCommand
	return p
}

func InitEmptyJoinCommandContext(p *JoinCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_joinCommand
}

func (*JoinCommandContext) IsJoinCommandContext() {}

func NewJoinCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JoinCommandContext {
	var p = new(JoinCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_joinCommand

	return p
}

func (s *JoinCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *JoinCommandContext) JOIN() antlr.TerminalNode {
	return s.GetToken(PPLParserJOIN, 0)
}

func (s *JoinCommandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *JoinCommandContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(PPLParserLBRACKET, 0)
}

func (s *JoinCommandContext) SearchQuery() ISearchQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISearchQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISearchQueryContext)
}

func (s *JoinCommandContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PPLParserRBRACKET, 0)
}

func (s *JoinCommandContext) TYPE() antlr.TerminalNode {
	return s.GetToken(PPLParserTYPE, 0)
}

func (s *JoinCommandContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *JoinCommandContext) JoinType() IJoinTypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoinTypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoinTypeContext)
}

func (s *JoinCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *JoinCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterJoinCommand(s)
	}
}

func (s *JoinCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitJoinCommand(s)
	}
}

func (s *JoinCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitJoinCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) JoinCommand() (localctx IJoinCommandContext) {
	localctx = NewJoinCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, PPLParserRULE_joinCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(514)
		p.Match(PPLParserJOIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(518)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserTYPE {
		{
			p.SetState(515)
			p.Match(PPLParserTYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(516)
			p.Match(PPLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(517)
			p.JoinType()
		}

	}
	{
		p.SetState(520)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(521)
		p.Match(PPLParserLBRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(522)
		p.SearchQuery()
	}
	{
		p.SetState(523)
		p.Match(PPLParserRBRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IJoinTypeContext is an interface to support dynamic dispatch.
type IJoinTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INNER() antlr.TerminalNode
	LEFT() antlr.TerminalNode
	RIGHT() antlr.TerminalNode
	OUTER() antlr.TerminalNode
	FULL() antlr.TerminalNode

	// IsJoinTypeContext differentiates from other interfaces.
	IsJoinTypeContext()
}

type JoinTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJoinTypeContext() *JoinTypeContext {
	var p = new(JoinTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_joinType
	return p
}

func InitEmptyJoinTypeContext(p *JoinTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_joinType
}

func (*JoinTypeContext) IsJoinTypeContext() {}

func NewJoinTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JoinTypeContext {
	var p = new(JoinTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_joinType

	return p
}

func (s *JoinTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *JoinTypeContext) INNER() antlr.TerminalNode {
	return s.GetToken(PPLParserINNER, 0)
}

func (s *JoinTypeContext) LEFT() antlr.TerminalNode {
	return s.GetToken(PPLParserLEFT, 0)
}

func (s *JoinTypeContext) RIGHT() antlr.TerminalNode {
	return s.GetToken(PPLParserRIGHT, 0)
}

func (s *JoinTypeContext) OUTER() antlr.TerminalNode {
	return s.GetToken(PPLParserOUTER, 0)
}

func (s *JoinTypeContext) FULL() antlr.TerminalNode {
	return s.GetToken(PPLParserFULL, 0)
}

func (s *JoinTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *JoinTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterJoinType(s)
	}
}

func (s *JoinTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitJoinType(s)
	}
}

func (s *JoinTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitJoinType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) JoinType() (localctx IJoinTypeContext) {
	localctx = NewJoinTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, PPLParserRULE_joinType)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(525)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2130303778816) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITableCommandContext is an interface to support dynamic dispatch.
type ITableCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TABLE() antlr.TerminalNode
	FieldList() IFieldListContext

	// IsTableCommandContext differentiates from other interfaces.
	IsTableCommandContext()
}

type TableCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTableCommandContext() *TableCommandContext {
	var p = new(TableCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_tableCommand
	return p
}

func InitEmptyTableCommandContext(p *TableCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_tableCommand
}

func (*TableCommandContext) IsTableCommandContext() {}

func NewTableCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TableCommandContext {
	var p = new(TableCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_tableCommand

	return p
}

func (s *TableCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *TableCommandContext) TABLE() antlr.TerminalNode {
	return s.GetToken(PPLParserTABLE, 0)
}

func (s *TableCommandContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *TableCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TableCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TableCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterTableCommand(s)
	}
}

func (s *TableCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitTableCommand(s)
	}
}

func (s *TableCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitTableCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) TableCommand() (localctx ITableCommandContext) {
	localctx = NewTableCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, PPLParserRULE_tableCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(527)
		p.Match(PPLParserTABLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(528)
		p.FieldList()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IEventstatsCommandContext is an interface to support dynamic dispatch.
type IEventstatsCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EVENTSTATS() antlr.TerminalNode
	AggregationList() IAggregationListContext
	BY() antlr.TerminalNode
	FieldList() IFieldListContext

	// IsEventstatsCommandContext differentiates from other interfaces.
	IsEventstatsCommandContext()
}

type EventstatsCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEventstatsCommandContext() *EventstatsCommandContext {
	var p = new(EventstatsCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_eventstatsCommand
	return p
}

func InitEmptyEventstatsCommandContext(p *EventstatsCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_eventstatsCommand
}

func (*EventstatsCommandContext) IsEventstatsCommandContext() {}

func NewEventstatsCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EventstatsCommandContext {
	var p = new(EventstatsCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_eventstatsCommand

	return p
}

func (s *EventstatsCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *EventstatsCommandContext) EVENTSTATS() antlr.TerminalNode {
	return s.GetToken(PPLParserEVENTSTATS, 0)
}

func (s *EventstatsCommandContext) AggregationList() IAggregationListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationListContext)
}

func (s *EventstatsCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(PPLParserBY, 0)
}

func (s *EventstatsCommandContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *EventstatsCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EventstatsCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EventstatsCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterEventstatsCommand(s)
	}
}

func (s *EventstatsCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitEventstatsCommand(s)
	}
}

func (s *EventstatsCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitEventstatsCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) EventstatsCommand() (localctx IEventstatsCommandContext) {
	localctx = NewEventstatsCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, PPLParserRULE_eventstatsCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(530)
		p.Match(PPLParserEVENTSTATS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(531)
		p.AggregationList()
	}
	p.SetState(534)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserBY {
		{
			p.SetState(532)
			p.Match(PPLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(533)
			p.FieldList()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStreamstatsCommandContext is an interface to support dynamic dispatch.
type IStreamstatsCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STREAMSTATS() antlr.TerminalNode
	AggregationList() IAggregationListContext
	AllStreamstatsOptions() []IStreamstatsOptionsContext
	StreamstatsOptions(i int) IStreamstatsOptionsContext
	BY() antlr.TerminalNode
	FieldList() IFieldListContext

	// IsStreamstatsCommandContext differentiates from other interfaces.
	IsStreamstatsCommandContext()
}

type StreamstatsCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStreamstatsCommandContext() *StreamstatsCommandContext {
	var p = new(StreamstatsCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_streamstatsCommand
	return p
}

func InitEmptyStreamstatsCommandContext(p *StreamstatsCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_streamstatsCommand
}

func (*StreamstatsCommandContext) IsStreamstatsCommandContext() {}

func NewStreamstatsCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StreamstatsCommandContext {
	var p = new(StreamstatsCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_streamstatsCommand

	return p
}

func (s *StreamstatsCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *StreamstatsCommandContext) STREAMSTATS() antlr.TerminalNode {
	return s.GetToken(PPLParserSTREAMSTATS, 0)
}

func (s *StreamstatsCommandContext) AggregationList() IAggregationListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationListContext)
}

func (s *StreamstatsCommandContext) AllStreamstatsOptions() []IStreamstatsOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStreamstatsOptionsContext); ok {
			len++
		}
	}

	tst := make([]IStreamstatsOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStreamstatsOptionsContext); ok {
			tst[i] = t.(IStreamstatsOptionsContext)
			i++
		}
	}

	return tst
}

func (s *StreamstatsCommandContext) StreamstatsOptions(i int) IStreamstatsOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStreamstatsOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStreamstatsOptionsContext)
}

func (s *StreamstatsCommandContext) BY() antlr.TerminalNode {
	return s.GetToken(PPLParserBY, 0)
}

func (s *StreamstatsCommandContext) FieldList() IFieldListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldListContext)
}

func (s *StreamstatsCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StreamstatsCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StreamstatsCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterStreamstatsCommand(s)
	}
}

func (s *StreamstatsCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitStreamstatsCommand(s)
	}
}

func (s *StreamstatsCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitStreamstatsCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) StreamstatsCommand() (localctx IStreamstatsCommandContext) {
	localctx = NewStreamstatsCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, PPLParserRULE_streamstatsCommand)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(536)
		p.Match(PPLParserSTREAMSTATS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(540)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(537)
				p.StreamstatsOptions()
			}

		}
		p.SetState(542)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	{
		p.SetState(543)
		p.AggregationList()
	}
	p.SetState(546)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserBY {
		{
			p.SetState(544)
			p.Match(PPLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(545)
			p.FieldList()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStreamstatsOptionsContext is an interface to support dynamic dispatch.
type IStreamstatsOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	EQ() antlr.TerminalNode
	INTEGER() antlr.TerminalNode
	STRING() antlr.TerminalNode
	BooleanValue() IBooleanValueContext

	// IsStreamstatsOptionsContext differentiates from other interfaces.
	IsStreamstatsOptionsContext()
}

type StreamstatsOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStreamstatsOptionsContext() *StreamstatsOptionsContext {
	var p = new(StreamstatsOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_streamstatsOptions
	return p
}

func InitEmptyStreamstatsOptionsContext(p *StreamstatsOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_streamstatsOptions
}

func (*StreamstatsOptionsContext) IsStreamstatsOptionsContext() {}

func NewStreamstatsOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StreamstatsOptionsContext {
	var p = new(StreamstatsOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_streamstatsOptions

	return p
}

func (s *StreamstatsOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *StreamstatsOptionsContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *StreamstatsOptionsContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *StreamstatsOptionsContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *StreamstatsOptionsContext) STRING() antlr.TerminalNode {
	return s.GetToken(PPLParserSTRING, 0)
}

func (s *StreamstatsOptionsContext) BooleanValue() IBooleanValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanValueContext)
}

func (s *StreamstatsOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StreamstatsOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StreamstatsOptionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterStreamstatsOptions(s)
	}
}

func (s *StreamstatsOptionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitStreamstatsOptions(s)
	}
}

func (s *StreamstatsOptionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitStreamstatsOptions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) StreamstatsOptions() (localctx IStreamstatsOptionsContext) {
	localctx = NewStreamstatsOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, PPLParserRULE_streamstatsOptions)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(548)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(549)
		p.Match(PPLParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(553)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserINTEGER:
		{
			p.SetState(550)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserSTRING:
		{
			p.SetState(551)
			p.Match(PPLParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PPLParserTRUE, PPLParserFALSE:
		{
			p.SetState(552)
			p.BooleanValue()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IReverseCommandContext is an interface to support dynamic dispatch.
type IReverseCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REVERSE() antlr.TerminalNode

	// IsReverseCommandContext differentiates from other interfaces.
	IsReverseCommandContext()
}

type ReverseCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReverseCommandContext() *ReverseCommandContext {
	var p = new(ReverseCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_reverseCommand
	return p
}

func InitEmptyReverseCommandContext(p *ReverseCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_reverseCommand
}

func (*ReverseCommandContext) IsReverseCommandContext() {}

func NewReverseCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReverseCommandContext {
	var p = new(ReverseCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_reverseCommand

	return p
}

func (s *ReverseCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ReverseCommandContext) REVERSE() antlr.TerminalNode {
	return s.GetToken(PPLParserREVERSE, 0)
}

func (s *ReverseCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReverseCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReverseCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterReverseCommand(s)
	}
}

func (s *ReverseCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitReverseCommand(s)
	}
}

func (s *ReverseCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitReverseCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ReverseCommand() (localctx IReverseCommandContext) {
	localctx = NewReverseCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, PPLParserRULE_reverseCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(555)
		p.Match(PPLParserREVERSE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFlattenCommandContext is an interface to support dynamic dispatch.
type IFlattenCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FLATTEN() antlr.TerminalNode
	FieldReference() IFieldReferenceContext

	// IsFlattenCommandContext differentiates from other interfaces.
	IsFlattenCommandContext()
}

type FlattenCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFlattenCommandContext() *FlattenCommandContext {
	var p = new(FlattenCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_flattenCommand
	return p
}

func InitEmptyFlattenCommandContext(p *FlattenCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_flattenCommand
}

func (*FlattenCommandContext) IsFlattenCommandContext() {}

func NewFlattenCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FlattenCommandContext {
	var p = new(FlattenCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_flattenCommand

	return p
}

func (s *FlattenCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *FlattenCommandContext) FLATTEN() antlr.TerminalNode {
	return s.GetToken(PPLParserFLATTEN, 0)
}

func (s *FlattenCommandContext) FieldReference() IFieldReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldReferenceContext)
}

func (s *FlattenCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FlattenCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FlattenCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFlattenCommand(s)
	}
}

func (s *FlattenCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFlattenCommand(s)
	}
}

func (s *FlattenCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFlattenCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) FlattenCommand() (localctx IFlattenCommandContext) {
	localctx = NewFlattenCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, PPLParserRULE_flattenCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(557)
		p.Match(PPLParserFLATTEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(558)
		p.FieldReference()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBooleanValueContext is an interface to support dynamic dispatch.
type IBooleanValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode

	// IsBooleanValueContext differentiates from other interfaces.
	IsBooleanValueContext()
}

type BooleanValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanValueContext() *BooleanValueContext {
	var p = new(BooleanValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_booleanValue
	return p
}

func InitEmptyBooleanValueContext(p *BooleanValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_booleanValue
}

func (*BooleanValueContext) IsBooleanValueContext() {}

func NewBooleanValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanValueContext {
	var p = new(BooleanValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_booleanValue

	return p
}

func (s *BooleanValueContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanValueContext) TRUE() antlr.TerminalNode {
	return s.GetToken(PPLParserTRUE, 0)
}

func (s *BooleanValueContext) FALSE() antlr.TerminalNode {
	return s.GetToken(PPLParserFALSE, 0)
}

func (s *BooleanValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BooleanValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterBooleanValue(s)
	}
}

func (s *BooleanValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitBooleanValue(s)
	}
}

func (s *BooleanValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitBooleanValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) BooleanValue() (localctx IBooleanValueContext) {
	localctx = NewBooleanValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, PPLParserRULE_booleanValue)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(560)
		_la = p.GetTokenStream().LA(1)

		if !(_la == PPLParserTRUE || _la == PPLParserFALSE) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDescribeCommandContext is an interface to support dynamic dispatch.
type IDescribeCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DESCRIBE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode

	// IsDescribeCommandContext differentiates from other interfaces.
	IsDescribeCommandContext()
}

type DescribeCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDescribeCommandContext() *DescribeCommandContext {
	var p = new(DescribeCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_describeCommand
	return p
}

func InitEmptyDescribeCommandContext(p *DescribeCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_describeCommand
}

func (*DescribeCommandContext) IsDescribeCommandContext() {}

func NewDescribeCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DescribeCommandContext {
	var p = new(DescribeCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_describeCommand

	return p
}

func (s *DescribeCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *DescribeCommandContext) DESCRIBE() antlr.TerminalNode {
	return s.GetToken(PPLParserDESCRIBE, 0)
}

func (s *DescribeCommandContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *DescribeCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DescribeCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DescribeCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterDescribeCommand(s)
	}
}

func (s *DescribeCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitDescribeCommand(s)
	}
}

func (s *DescribeCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitDescribeCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) DescribeCommand() (localctx IDescribeCommandContext) {
	localctx = NewDescribeCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, PPLParserRULE_describeCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(562)
		p.Match(PPLParserDESCRIBE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(563)
		p.Match(PPLParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IShowDatasourcesCommandContext is an interface to support dynamic dispatch.
type IShowDatasourcesCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SHOWDATASOURCES() antlr.TerminalNode

	// IsShowDatasourcesCommandContext differentiates from other interfaces.
	IsShowDatasourcesCommandContext()
}

type ShowDatasourcesCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowDatasourcesCommandContext() *ShowDatasourcesCommandContext {
	var p = new(ShowDatasourcesCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_showDatasourcesCommand
	return p
}

func InitEmptyShowDatasourcesCommandContext(p *ShowDatasourcesCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_showDatasourcesCommand
}

func (*ShowDatasourcesCommandContext) IsShowDatasourcesCommandContext() {}

func NewShowDatasourcesCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowDatasourcesCommandContext {
	var p = new(ShowDatasourcesCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_showDatasourcesCommand

	return p
}

func (s *ShowDatasourcesCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowDatasourcesCommandContext) SHOWDATASOURCES() antlr.TerminalNode {
	return s.GetToken(PPLParserSHOWDATASOURCES, 0)
}

func (s *ShowDatasourcesCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowDatasourcesCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowDatasourcesCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterShowDatasourcesCommand(s)
	}
}

func (s *ShowDatasourcesCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitShowDatasourcesCommand(s)
	}
}

func (s *ShowDatasourcesCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitShowDatasourcesCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ShowDatasourcesCommand() (localctx IShowDatasourcesCommandContext) {
	localctx = NewShowDatasourcesCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, PPLParserRULE_showDatasourcesCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(565)
		p.Match(PPLParserSHOWDATASOURCES)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExplainCommandContext is an interface to support dynamic dispatch.
type IExplainCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXPLAIN() antlr.TerminalNode

	// IsExplainCommandContext differentiates from other interfaces.
	IsExplainCommandContext()
}

type ExplainCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExplainCommandContext() *ExplainCommandContext {
	var p = new(ExplainCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_explainCommand
	return p
}

func InitEmptyExplainCommandContext(p *ExplainCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_explainCommand
}

func (*ExplainCommandContext) IsExplainCommandContext() {}

func NewExplainCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExplainCommandContext {
	var p = new(ExplainCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_explainCommand

	return p
}

func (s *ExplainCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ExplainCommandContext) EXPLAIN() antlr.TerminalNode {
	return s.GetToken(PPLParserEXPLAIN, 0)
}

func (s *ExplainCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExplainCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExplainCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterExplainCommand(s)
	}
}

func (s *ExplainCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitExplainCommand(s)
	}
}

func (s *ExplainCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitExplainCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ExplainCommand() (localctx IExplainCommandContext) {
	localctx = NewExplainCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, PPLParserRULE_explainCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(567)
		p.Match(PPLParserEXPLAIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OrExpression() IOrExpressionContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) OrExpression() IOrExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrExpressionContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (s *ExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, PPLParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(569)
		p.OrExpression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOrExpressionContext is an interface to support dynamic dispatch.
type IOrExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAndExpression() []IAndExpressionContext
	AndExpression(i int) IAndExpressionContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsOrExpressionContext differentiates from other interfaces.
	IsOrExpressionContext()
}

type OrExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrExpressionContext() *OrExpressionContext {
	var p = new(OrExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_orExpression
	return p
}

func InitEmptyOrExpressionContext(p *OrExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_orExpression
}

func (*OrExpressionContext) IsOrExpressionContext() {}

func NewOrExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrExpressionContext {
	var p = new(OrExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_orExpression

	return p
}

func (s *OrExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *OrExpressionContext) AllAndExpression() []IAndExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAndExpressionContext); ok {
			len++
		}
	}

	tst := make([]IAndExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAndExpressionContext); ok {
			tst[i] = t.(IAndExpressionContext)
			i++
		}
	}

	return tst
}

func (s *OrExpressionContext) AndExpression(i int) IAndExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAndExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAndExpressionContext)
}

func (s *OrExpressionContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(PPLParserOR)
}

func (s *OrExpressionContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserOR, i)
}

func (s *OrExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterOrExpression(s)
	}
}

func (s *OrExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitOrExpression(s)
	}
}

func (s *OrExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitOrExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) OrExpression() (localctx IOrExpressionContext) {
	localctx = NewOrExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, PPLParserRULE_orExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(571)
		p.AndExpression()
	}
	p.SetState(576)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserOR {
		{
			p.SetState(572)
			p.Match(PPLParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(573)
			p.AndExpression()
		}

		p.SetState(578)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAndExpressionContext is an interface to support dynamic dispatch.
type IAndExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNotExpression() []INotExpressionContext
	NotExpression(i int) INotExpressionContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsAndExpressionContext differentiates from other interfaces.
	IsAndExpressionContext()
}

type AndExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAndExpressionContext() *AndExpressionContext {
	var p = new(AndExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_andExpression
	return p
}

func InitEmptyAndExpressionContext(p *AndExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_andExpression
}

func (*AndExpressionContext) IsAndExpressionContext() {}

func NewAndExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AndExpressionContext {
	var p = new(AndExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_andExpression

	return p
}

func (s *AndExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *AndExpressionContext) AllNotExpression() []INotExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INotExpressionContext); ok {
			len++
		}
	}

	tst := make([]INotExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INotExpressionContext); ok {
			tst[i] = t.(INotExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AndExpressionContext) NotExpression(i int) INotExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INotExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INotExpressionContext)
}

func (s *AndExpressionContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(PPLParserAND)
}

func (s *AndExpressionContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserAND, i)
}

func (s *AndExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AndExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAndExpression(s)
	}
}

func (s *AndExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAndExpression(s)
	}
}

func (s *AndExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAndExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) AndExpression() (localctx IAndExpressionContext) {
	localctx = NewAndExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, PPLParserRULE_andExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(579)
		p.NotExpression()
	}
	p.SetState(584)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserAND {
		{
			p.SetState(580)
			p.Match(PPLParserAND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(581)
			p.NotExpression()
		}

		p.SetState(586)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INotExpressionContext is an interface to support dynamic dispatch.
type INotExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NOT() antlr.TerminalNode
	NotExpression() INotExpressionContext
	ComparisonExpression() IComparisonExpressionContext

	// IsNotExpressionContext differentiates from other interfaces.
	IsNotExpressionContext()
}

type NotExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNotExpressionContext() *NotExpressionContext {
	var p = new(NotExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_notExpression
	return p
}

func InitEmptyNotExpressionContext(p *NotExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_notExpression
}

func (*NotExpressionContext) IsNotExpressionContext() {}

func NewNotExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NotExpressionContext {
	var p = new(NotExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_notExpression

	return p
}

func (s *NotExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *NotExpressionContext) NOT() antlr.TerminalNode {
	return s.GetToken(PPLParserNOT, 0)
}

func (s *NotExpressionContext) NotExpression() INotExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INotExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INotExpressionContext)
}

func (s *NotExpressionContext) ComparisonExpression() IComparisonExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonExpressionContext)
}

func (s *NotExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NotExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NotExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterNotExpression(s)
	}
}

func (s *NotExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitNotExpression(s)
	}
}

func (s *NotExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitNotExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) NotExpression() (localctx INotExpressionContext) {
	localctx = NewNotExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, PPLParserRULE_notExpression)
	p.SetState(590)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserNOT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(587)
			p.Match(PPLParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(588)
			p.NotExpression()
		}

	case PPLParserPLUS, PPLParserMINUS, PPLParserTRUE, PPLParserFALSE, PPLParserNULL, PPLParserCASE, PPLParserCOUNT, PPLParserSUM, PPLParserAVG, PPLParserMIN, PPLParserMAX, PPLParserDC, PPLParserDISTINCT_COUNT, PPLParserVAR, PPLParserVARP, PPLParserSTDEV, PPLParserSTDEVP, PPLParserPERCENTILE, PPLParserMEDIAN, PPLParserMODE, PPLParserEARLIEST, PPLParserLATEST, PPLParserVALUES, PPLParserRANGE, PPLParserLPAREN, PPLParserIDENTIFIER, PPLParserINTEGER, PPLParserDECIMAL, PPLParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(589)
			p.ComparisonExpression()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComparisonExpressionContext is an interface to support dynamic dispatch.
type IComparisonExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAdditiveExpression() []IAdditiveExpressionContext
	AdditiveExpression(i int) IAdditiveExpressionContext
	LIKE() antlr.TerminalNode
	IN() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ExpressionList() IExpressionListContext
	RPAREN() antlr.TerminalNode
	EQ() antlr.TerminalNode
	NEQ() antlr.TerminalNode
	LT() antlr.TerminalNode
	LTE() antlr.TerminalNode
	GT() antlr.TerminalNode
	GTE() antlr.TerminalNode

	// IsComparisonExpressionContext differentiates from other interfaces.
	IsComparisonExpressionContext()
}

type ComparisonExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonExpressionContext() *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_comparisonExpression
	return p
}

func InitEmptyComparisonExpressionContext(p *ComparisonExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_comparisonExpression
}

func (*ComparisonExpressionContext) IsComparisonExpressionContext() {}

func NewComparisonExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_comparisonExpression

	return p
}

func (s *ComparisonExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonExpressionContext) AllAdditiveExpression() []IAdditiveExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAdditiveExpressionContext); ok {
			len++
		}
	}

	tst := make([]IAdditiveExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAdditiveExpressionContext); ok {
			tst[i] = t.(IAdditiveExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ComparisonExpressionContext) AdditiveExpression(i int) IAdditiveExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAdditiveExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAdditiveExpressionContext)
}

func (s *ComparisonExpressionContext) LIKE() antlr.TerminalNode {
	return s.GetToken(PPLParserLIKE, 0)
}

func (s *ComparisonExpressionContext) IN() antlr.TerminalNode {
	return s.GetToken(PPLParserIN, 0)
}

func (s *ComparisonExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserLPAREN, 0)
}

func (s *ComparisonExpressionContext) ExpressionList() IExpressionListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionListContext)
}

func (s *ComparisonExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserRPAREN, 0)
}

func (s *ComparisonExpressionContext) EQ() antlr.TerminalNode {
	return s.GetToken(PPLParserEQ, 0)
}

func (s *ComparisonExpressionContext) NEQ() antlr.TerminalNode {
	return s.GetToken(PPLParserNEQ, 0)
}

func (s *ComparisonExpressionContext) LT() antlr.TerminalNode {
	return s.GetToken(PPLParserLT, 0)
}

func (s *ComparisonExpressionContext) LTE() antlr.TerminalNode {
	return s.GetToken(PPLParserLTE, 0)
}

func (s *ComparisonExpressionContext) GT() antlr.TerminalNode {
	return s.GetToken(PPLParserGT, 0)
}

func (s *ComparisonExpressionContext) GTE() antlr.TerminalNode {
	return s.GetToken(PPLParserGTE, 0)
}

func (s *ComparisonExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterComparisonExpression(s)
	}
}

func (s *ComparisonExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitComparisonExpression(s)
	}
}

func (s *ComparisonExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitComparisonExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ComparisonExpression() (localctx IComparisonExpressionContext) {
	localctx = NewComparisonExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, PPLParserRULE_comparisonExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(592)
		p.AdditiveExpression()
	}
	p.SetState(602)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(593)
			_la = p.GetTokenStream().LA(1)

			if !((int64((_la-59)) & ^0x3f) == 0 && ((int64(1)<<(_la-59))&63) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(594)
			p.AdditiveExpression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) == 2 {
		{
			p.SetState(595)
			p.Match(PPLParserLIKE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(596)
			p.AdditiveExpression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) == 3 {
		{
			p.SetState(597)
			p.Match(PPLParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(598)
			p.Match(PPLParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(599)
			p.ExpressionList()
		}
		{
			p.SetState(600)
			p.Match(PPLParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAdditiveExpressionContext is an interface to support dynamic dispatch.
type IAdditiveExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllMultiplicativeExpression() []IMultiplicativeExpressionContext
	MultiplicativeExpression(i int) IMultiplicativeExpressionContext
	AllPLUS() []antlr.TerminalNode
	PLUS(i int) antlr.TerminalNode
	AllMINUS() []antlr.TerminalNode
	MINUS(i int) antlr.TerminalNode

	// IsAdditiveExpressionContext differentiates from other interfaces.
	IsAdditiveExpressionContext()
}

type AdditiveExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAdditiveExpressionContext() *AdditiveExpressionContext {
	var p = new(AdditiveExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_additiveExpression
	return p
}

func InitEmptyAdditiveExpressionContext(p *AdditiveExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_additiveExpression
}

func (*AdditiveExpressionContext) IsAdditiveExpressionContext() {}

func NewAdditiveExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AdditiveExpressionContext {
	var p = new(AdditiveExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_additiveExpression

	return p
}

func (s *AdditiveExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *AdditiveExpressionContext) AllMultiplicativeExpression() []IMultiplicativeExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMultiplicativeExpressionContext); ok {
			len++
		}
	}

	tst := make([]IMultiplicativeExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMultiplicativeExpressionContext); ok {
			tst[i] = t.(IMultiplicativeExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AdditiveExpressionContext) MultiplicativeExpression(i int) IMultiplicativeExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMultiplicativeExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMultiplicativeExpressionContext)
}

func (s *AdditiveExpressionContext) AllPLUS() []antlr.TerminalNode {
	return s.GetTokens(PPLParserPLUS)
}

func (s *AdditiveExpressionContext) PLUS(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserPLUS, i)
}

func (s *AdditiveExpressionContext) AllMINUS() []antlr.TerminalNode {
	return s.GetTokens(PPLParserMINUS)
}

func (s *AdditiveExpressionContext) MINUS(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserMINUS, i)
}

func (s *AdditiveExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditiveExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AdditiveExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAdditiveExpression(s)
	}
}

func (s *AdditiveExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAdditiveExpression(s)
	}
}

func (s *AdditiveExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAdditiveExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) AdditiveExpression() (localctx IAdditiveExpressionContext) {
	localctx = NewAdditiveExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, PPLParserRULE_additiveExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(604)
		p.MultiplicativeExpression()
	}
	p.SetState(609)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserPLUS || _la == PPLParserMINUS {
		{
			p.SetState(605)
			_la = p.GetTokenStream().LA(1)

			if !(_la == PPLParserPLUS || _la == PPLParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(606)
			p.MultiplicativeExpression()
		}

		p.SetState(611)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMultiplicativeExpressionContext is an interface to support dynamic dispatch.
type IMultiplicativeExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllUnaryExpression() []IUnaryExpressionContext
	UnaryExpression(i int) IUnaryExpressionContext
	AllSTAR() []antlr.TerminalNode
	STAR(i int) antlr.TerminalNode
	AllSLASH() []antlr.TerminalNode
	SLASH(i int) antlr.TerminalNode
	AllPERCENT() []antlr.TerminalNode
	PERCENT(i int) antlr.TerminalNode

	// IsMultiplicativeExpressionContext differentiates from other interfaces.
	IsMultiplicativeExpressionContext()
}

type MultiplicativeExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMultiplicativeExpressionContext() *MultiplicativeExpressionContext {
	var p = new(MultiplicativeExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_multiplicativeExpression
	return p
}

func InitEmptyMultiplicativeExpressionContext(p *MultiplicativeExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_multiplicativeExpression
}

func (*MultiplicativeExpressionContext) IsMultiplicativeExpressionContext() {}

func NewMultiplicativeExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MultiplicativeExpressionContext {
	var p = new(MultiplicativeExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_multiplicativeExpression

	return p
}

func (s *MultiplicativeExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *MultiplicativeExpressionContext) AllUnaryExpression() []IUnaryExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUnaryExpressionContext); ok {
			len++
		}
	}

	tst := make([]IUnaryExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUnaryExpressionContext); ok {
			tst[i] = t.(IUnaryExpressionContext)
			i++
		}
	}

	return tst
}

func (s *MultiplicativeExpressionContext) UnaryExpression(i int) IUnaryExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnaryExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnaryExpressionContext)
}

func (s *MultiplicativeExpressionContext) AllSTAR() []antlr.TerminalNode {
	return s.GetTokens(PPLParserSTAR)
}

func (s *MultiplicativeExpressionContext) STAR(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserSTAR, i)
}

func (s *MultiplicativeExpressionContext) AllSLASH() []antlr.TerminalNode {
	return s.GetTokens(PPLParserSLASH)
}

func (s *MultiplicativeExpressionContext) SLASH(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserSLASH, i)
}

func (s *MultiplicativeExpressionContext) AllPERCENT() []antlr.TerminalNode {
	return s.GetTokens(PPLParserPERCENT)
}

func (s *MultiplicativeExpressionContext) PERCENT(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserPERCENT, i)
}

func (s *MultiplicativeExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicativeExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MultiplicativeExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterMultiplicativeExpression(s)
	}
}

func (s *MultiplicativeExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitMultiplicativeExpression(s)
	}
}

func (s *MultiplicativeExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitMultiplicativeExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) MultiplicativeExpression() (localctx IMultiplicativeExpressionContext) {
	localctx = NewMultiplicativeExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, PPLParserRULE_multiplicativeExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(612)
		p.UnaryExpression()
	}
	p.SetState(617)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64((_la-69)) & ^0x3f) == 0 && ((int64(1)<<(_la-69))&7) != 0 {
		{
			p.SetState(613)
			_la = p.GetTokenStream().LA(1)

			if !((int64((_la-69)) & ^0x3f) == 0 && ((int64(1)<<(_la-69))&7) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(614)
			p.UnaryExpression()
		}

		p.SetState(619)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnaryExpressionContext is an interface to support dynamic dispatch.
type IUnaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UnaryExpression() IUnaryExpressionContext
	PLUS() antlr.TerminalNode
	MINUS() antlr.TerminalNode
	PrimaryExpression() IPrimaryExpressionContext

	// IsUnaryExpressionContext differentiates from other interfaces.
	IsUnaryExpressionContext()
}

type UnaryExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnaryExpressionContext() *UnaryExpressionContext {
	var p = new(UnaryExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_unaryExpression
	return p
}

func InitEmptyUnaryExpressionContext(p *UnaryExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_unaryExpression
}

func (*UnaryExpressionContext) IsUnaryExpressionContext() {}

func NewUnaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryExpressionContext {
	var p = new(UnaryExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_unaryExpression

	return p
}

func (s *UnaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *UnaryExpressionContext) UnaryExpression() IUnaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnaryExpressionContext)
}

func (s *UnaryExpressionContext) PLUS() antlr.TerminalNode {
	return s.GetToken(PPLParserPLUS, 0)
}

func (s *UnaryExpressionContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PPLParserMINUS, 0)
}

func (s *UnaryExpressionContext) PrimaryExpression() IPrimaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionContext)
}

func (s *UnaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterUnaryExpression(s)
	}
}

func (s *UnaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitUnaryExpression(s)
	}
}

func (s *UnaryExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitUnaryExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) UnaryExpression() (localctx IUnaryExpressionContext) {
	localctx = NewUnaryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, PPLParserRULE_unaryExpression)
	var _la int

	p.SetState(623)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PPLParserPLUS, PPLParserMINUS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(620)
			_la = p.GetTokenStream().LA(1)

			if !(_la == PPLParserPLUS || _la == PPLParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(621)
			p.UnaryExpression()
		}

	case PPLParserTRUE, PPLParserFALSE, PPLParserNULL, PPLParserCASE, PPLParserCOUNT, PPLParserSUM, PPLParserAVG, PPLParserMIN, PPLParserMAX, PPLParserDC, PPLParserDISTINCT_COUNT, PPLParserVAR, PPLParserVARP, PPLParserSTDEV, PPLParserSTDEVP, PPLParserPERCENTILE, PPLParserMEDIAN, PPLParserMODE, PPLParserEARLIEST, PPLParserLATEST, PPLParserVALUES, PPLParserRANGE, PPLParserLPAREN, PPLParserIDENTIFIER, PPLParserINTEGER, PPLParserDECIMAL, PPLParserSTRING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(622)
			p.PrimaryExpression()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimaryExpressionContext is an interface to support dynamic dispatch.
type IPrimaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Literal() ILiteralContext
	FieldReference() IFieldReferenceContext
	FunctionCall() IFunctionCallContext
	CaseExpression() ICaseExpressionContext
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode

	// IsPrimaryExpressionContext differentiates from other interfaces.
	IsPrimaryExpressionContext()
}

type PrimaryExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryExpressionContext() *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_primaryExpression
	return p
}

func InitEmptyPrimaryExpressionContext(p *PrimaryExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_primaryExpression
}

func (*PrimaryExpressionContext) IsPrimaryExpressionContext() {}

func NewPrimaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_primaryExpression

	return p
}

func (s *PrimaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExpressionContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *PrimaryExpressionContext) FieldReference() IFieldReferenceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldReferenceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldReferenceContext)
}

func (s *PrimaryExpressionContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *PrimaryExpressionContext) CaseExpression() ICaseExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICaseExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICaseExpressionContext)
}

func (s *PrimaryExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserLPAREN, 0)
}

func (s *PrimaryExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *PrimaryExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserRPAREN, 0)
}

func (s *PrimaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterPrimaryExpression(s)
	}
}

func (s *PrimaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitPrimaryExpression(s)
	}
}

func (s *PrimaryExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitPrimaryExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) PrimaryExpression() (localctx IPrimaryExpressionContext) {
	localctx = NewPrimaryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, PPLParserRULE_primaryExpression)
	p.SetState(633)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 54, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(625)
			p.Literal()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(626)
			p.FieldReference()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(627)
			p.FunctionCall()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(628)
			p.CaseExpression()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(629)
			p.Match(PPLParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(630)
			p.Expression()
		}
		{
			p.SetState(631)
			p.Match(PPLParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER() antlr.TerminalNode
	DECIMAL() antlr.TerminalNode
	STRING() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	NULL() antlr.TerminalNode

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_literal
	return p
}

func InitEmptyLiteralContext(p *LiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_literal
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *LiteralContext) DECIMAL() antlr.TerminalNode {
	return s.GetToken(PPLParserDECIMAL, 0)
}

func (s *LiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(PPLParserSTRING, 0)
}

func (s *LiteralContext) TRUE() antlr.TerminalNode {
	return s.GetToken(PPLParserTRUE, 0)
}

func (s *LiteralContext) FALSE() antlr.TerminalNode {
	return s.GetToken(PPLParserFALSE, 0)
}

func (s *LiteralContext) NULL() antlr.TerminalNode {
	return s.GetToken(PPLParserNULL, 0)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterLiteral(s)
	}
}

func (s *LiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitLiteral(s)
	}
}

func (s *LiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, PPLParserRULE_literal)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(635)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-72)) & ^0x3f) == 0 && ((int64(1)<<(_la-72))&240518168583) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldReferenceContext is an interface to support dynamic dispatch.
type IFieldReferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	INTEGER() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode

	// IsFieldReferenceContext differentiates from other interfaces.
	IsFieldReferenceContext()
}

type FieldReferenceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldReferenceContext() *FieldReferenceContext {
	var p = new(FieldReferenceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fieldReference
	return p
}

func InitEmptyFieldReferenceContext(p *FieldReferenceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_fieldReference
}

func (*FieldReferenceContext) IsFieldReferenceContext() {}

func NewFieldReferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldReferenceContext {
	var p = new(FieldReferenceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_fieldReference

	return p
}

func (s *FieldReferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldReferenceContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(PPLParserIDENTIFIER)
}

func (s *FieldReferenceContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, i)
}

func (s *FieldReferenceContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(PPLParserDOT)
}

func (s *FieldReferenceContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserDOT, i)
}

func (s *FieldReferenceContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(PPLParserLBRACKET, 0)
}

func (s *FieldReferenceContext) INTEGER() antlr.TerminalNode {
	return s.GetToken(PPLParserINTEGER, 0)
}

func (s *FieldReferenceContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PPLParserRBRACKET, 0)
}

func (s *FieldReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldReferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFieldReference(s)
	}
}

func (s *FieldReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFieldReference(s)
	}
}

func (s *FieldReferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFieldReference(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) FieldReference() (localctx IFieldReferenceContext) {
	localctx = NewFieldReferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, PPLParserRULE_fieldReference)
	var _la int

	p.SetState(649)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 56, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(637)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(642)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == PPLParserDOT {
			{
				p.SetState(638)
				p.Match(PPLParserDOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(639)
				p.Match(PPLParserIDENTIFIER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

			p.SetState(644)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(645)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(646)
			p.Match(PPLParserLBRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(647)
			p.Match(PPLParserINTEGER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(648)
			p.Match(PPLParserRBRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) CopyAll(ctx *FunctionCallContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type FunctionCallWithArgsContext struct {
	FunctionCallContext
}

func NewFunctionCallWithArgsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionCallWithArgsContext {
	var p = new(FunctionCallWithArgsContext)

	InitEmptyFunctionCallContext(&p.FunctionCallContext)
	p.parser = parser
	p.CopyAll(ctx.(*FunctionCallContext))

	return p
}

func (s *FunctionCallWithArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallWithArgsContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *FunctionCallWithArgsContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserLPAREN, 0)
}

func (s *FunctionCallWithArgsContext) ExpressionList() IExpressionListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionListContext)
}

func (s *FunctionCallWithArgsContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserRPAREN, 0)
}

func (s *FunctionCallWithArgsContext) DISTINCT() antlr.TerminalNode {
	return s.GetToken(PPLParserDISTINCT, 0)
}

func (s *FunctionCallWithArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFunctionCallWithArgs(s)
	}
}

func (s *FunctionCallWithArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFunctionCallWithArgs(s)
	}
}

func (s *FunctionCallWithArgsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFunctionCallWithArgs(s)

	default:
		return t.VisitChildren(s)
	}
}

type FunctionCallNoArgsContext struct {
	FunctionCallContext
}

func NewFunctionCallNoArgsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionCallNoArgsContext {
	var p = new(FunctionCallNoArgsContext)

	InitEmptyFunctionCallContext(&p.FunctionCallContext)
	p.parser = parser
	p.CopyAll(ctx.(*FunctionCallContext))

	return p
}

func (s *FunctionCallNoArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallNoArgsContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(PPLParserIDENTIFIER, 0)
}

func (s *FunctionCallNoArgsContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserLPAREN, 0)
}

func (s *FunctionCallNoArgsContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserRPAREN, 0)
}

func (s *FunctionCallNoArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterFunctionCallNoArgs(s)
	}
}

func (s *FunctionCallNoArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitFunctionCallNoArgs(s)
	}
}

func (s *FunctionCallNoArgsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitFunctionCallNoArgs(s)

	default:
		return t.VisitChildren(s)
	}
}

type AggregationFunctionCallNoArgsContext struct {
	FunctionCallContext
}

func NewAggregationFunctionCallNoArgsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AggregationFunctionCallNoArgsContext {
	var p = new(AggregationFunctionCallNoArgsContext)

	InitEmptyFunctionCallContext(&p.FunctionCallContext)
	p.parser = parser
	p.CopyAll(ctx.(*FunctionCallContext))

	return p
}

func (s *AggregationFunctionCallNoArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationFunctionCallNoArgsContext) AggregationFunction() IAggregationFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationFunctionContext)
}

func (s *AggregationFunctionCallNoArgsContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserLPAREN, 0)
}

func (s *AggregationFunctionCallNoArgsContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserRPAREN, 0)
}

func (s *AggregationFunctionCallNoArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAggregationFunctionCallNoArgs(s)
	}
}

func (s *AggregationFunctionCallNoArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAggregationFunctionCallNoArgs(s)
	}
}

func (s *AggregationFunctionCallNoArgsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAggregationFunctionCallNoArgs(s)

	default:
		return t.VisitChildren(s)
	}
}

type AggregationFunctionCallContext struct {
	FunctionCallContext
}

func NewAggregationFunctionCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AggregationFunctionCallContext {
	var p = new(AggregationFunctionCallContext)

	InitEmptyFunctionCallContext(&p.FunctionCallContext)
	p.parser = parser
	p.CopyAll(ctx.(*FunctionCallContext))

	return p
}

func (s *AggregationFunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationFunctionCallContext) AggregationFunction() IAggregationFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationFunctionContext)
}

func (s *AggregationFunctionCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserLPAREN, 0)
}

func (s *AggregationFunctionCallContext) ExpressionList() IExpressionListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionListContext)
}

func (s *AggregationFunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PPLParserRPAREN, 0)
}

func (s *AggregationFunctionCallContext) DISTINCT() antlr.TerminalNode {
	return s.GetToken(PPLParserDISTINCT, 0)
}

func (s *AggregationFunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAggregationFunctionCall(s)
	}
}

func (s *AggregationFunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAggregationFunctionCall(s)
	}
}

func (s *AggregationFunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAggregationFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, PPLParserRULE_functionCall)
	var _la int

	p.SetState(674)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 59, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFunctionCallNoArgsContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(651)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(652)
			p.Match(PPLParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(653)
			p.Match(PPLParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewFunctionCallWithArgsContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(654)
			p.Match(PPLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(655)
			p.Match(PPLParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(657)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == PPLParserDISTINCT {
			{
				p.SetState(656)
				p.Match(PPLParserDISTINCT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(659)
			p.ExpressionList()
		}
		{
			p.SetState(660)
			p.Match(PPLParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewAggregationFunctionCallNoArgsContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(662)
			p.AggregationFunction()
		}
		{
			p.SetState(663)
			p.Match(PPLParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(664)
			p.Match(PPLParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewAggregationFunctionCallContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(666)
			p.AggregationFunction()
		}
		{
			p.SetState(667)
			p.Match(PPLParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(669)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == PPLParserDISTINCT {
			{
				p.SetState(668)
				p.Match(PPLParserDISTINCT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(671)
			p.ExpressionList()
		}
		{
			p.SetState(672)
			p.Match(PPLParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAggregationFunctionContext is an interface to support dynamic dispatch.
type IAggregationFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COUNT() antlr.TerminalNode
	SUM() antlr.TerminalNode
	AVG() antlr.TerminalNode
	MIN() antlr.TerminalNode
	MAX() antlr.TerminalNode
	DC() antlr.TerminalNode
	DISTINCT_COUNT() antlr.TerminalNode
	VAR() antlr.TerminalNode
	VARP() antlr.TerminalNode
	STDEV() antlr.TerminalNode
	STDEVP() antlr.TerminalNode
	PERCENTILE() antlr.TerminalNode
	MEDIAN() antlr.TerminalNode
	MODE() antlr.TerminalNode
	EARLIEST() antlr.TerminalNode
	LATEST() antlr.TerminalNode
	VALUES() antlr.TerminalNode
	RANGE() antlr.TerminalNode

	// IsAggregationFunctionContext differentiates from other interfaces.
	IsAggregationFunctionContext()
}

type AggregationFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationFunctionContext() *AggregationFunctionContext {
	var p = new(AggregationFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_aggregationFunction
	return p
}

func InitEmptyAggregationFunctionContext(p *AggregationFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_aggregationFunction
}

func (*AggregationFunctionContext) IsAggregationFunctionContext() {}

func NewAggregationFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationFunctionContext {
	var p = new(AggregationFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_aggregationFunction

	return p
}

func (s *AggregationFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationFunctionContext) COUNT() antlr.TerminalNode {
	return s.GetToken(PPLParserCOUNT, 0)
}

func (s *AggregationFunctionContext) SUM() antlr.TerminalNode {
	return s.GetToken(PPLParserSUM, 0)
}

func (s *AggregationFunctionContext) AVG() antlr.TerminalNode {
	return s.GetToken(PPLParserAVG, 0)
}

func (s *AggregationFunctionContext) MIN() antlr.TerminalNode {
	return s.GetToken(PPLParserMIN, 0)
}

func (s *AggregationFunctionContext) MAX() antlr.TerminalNode {
	return s.GetToken(PPLParserMAX, 0)
}

func (s *AggregationFunctionContext) DC() antlr.TerminalNode {
	return s.GetToken(PPLParserDC, 0)
}

func (s *AggregationFunctionContext) DISTINCT_COUNT() antlr.TerminalNode {
	return s.GetToken(PPLParserDISTINCT_COUNT, 0)
}

func (s *AggregationFunctionContext) VAR() antlr.TerminalNode {
	return s.GetToken(PPLParserVAR, 0)
}

func (s *AggregationFunctionContext) VARP() antlr.TerminalNode {
	return s.GetToken(PPLParserVARP, 0)
}

func (s *AggregationFunctionContext) STDEV() antlr.TerminalNode {
	return s.GetToken(PPLParserSTDEV, 0)
}

func (s *AggregationFunctionContext) STDEVP() antlr.TerminalNode {
	return s.GetToken(PPLParserSTDEVP, 0)
}

func (s *AggregationFunctionContext) PERCENTILE() antlr.TerminalNode {
	return s.GetToken(PPLParserPERCENTILE, 0)
}

func (s *AggregationFunctionContext) MEDIAN() antlr.TerminalNode {
	return s.GetToken(PPLParserMEDIAN, 0)
}

func (s *AggregationFunctionContext) MODE() antlr.TerminalNode {
	return s.GetToken(PPLParserMODE, 0)
}

func (s *AggregationFunctionContext) EARLIEST() antlr.TerminalNode {
	return s.GetToken(PPLParserEARLIEST, 0)
}

func (s *AggregationFunctionContext) LATEST() antlr.TerminalNode {
	return s.GetToken(PPLParserLATEST, 0)
}

func (s *AggregationFunctionContext) VALUES() antlr.TerminalNode {
	return s.GetToken(PPLParserVALUES, 0)
}

func (s *AggregationFunctionContext) RANGE() antlr.TerminalNode {
	return s.GetToken(PPLParserRANGE, 0)
}

func (s *AggregationFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AggregationFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterAggregationFunction(s)
	}
}

func (s *AggregationFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitAggregationFunction(s)
	}
}

func (s *AggregationFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitAggregationFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) AggregationFunction() (localctx IAggregationFunctionContext) {
	localctx = NewAggregationFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, PPLParserRULE_aggregationFunction)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(676)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-80)) & ^0x3f) == 0 && ((int64(1)<<(_la-80))&262143) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionListContext is an interface to support dynamic dispatch.
type IExpressionListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExpressionListContext differentiates from other interfaces.
	IsExpressionListContext()
}

type ExpressionListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionListContext() *ExpressionListContext {
	var p = new(ExpressionListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_expressionList
	return p
}

func InitEmptyExpressionListContext(p *ExpressionListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_expressionList
}

func (*ExpressionListContext) IsExpressionListContext() {}

func NewExpressionListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionListContext {
	var p = new(ExpressionListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_expressionList

	return p
}

func (s *ExpressionListContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionListContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionListContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PPLParserCOMMA)
}

func (s *ExpressionListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PPLParserCOMMA, i)
}

func (s *ExpressionListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterExpressionList(s)
	}
}

func (s *ExpressionListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitExpressionList(s)
	}
}

func (s *ExpressionListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitExpressionList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) ExpressionList() (localctx IExpressionListContext) {
	localctx = NewExpressionListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, PPLParserRULE_expressionList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(678)
		p.Expression()
	}
	p.SetState(683)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PPLParserCOMMA {
		{
			p.SetState(679)
			p.Match(PPLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(680)
			p.Expression()
		}

		p.SetState(685)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICaseExpressionContext is an interface to support dynamic dispatch.
type ICaseExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	END() antlr.TerminalNode
	AllWhenClause() []IWhenClauseContext
	WhenClause(i int) IWhenClauseContext
	ELSE() antlr.TerminalNode
	Expression() IExpressionContext

	// IsCaseExpressionContext differentiates from other interfaces.
	IsCaseExpressionContext()
}

type CaseExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCaseExpressionContext() *CaseExpressionContext {
	var p = new(CaseExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_caseExpression
	return p
}

func InitEmptyCaseExpressionContext(p *CaseExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_caseExpression
}

func (*CaseExpressionContext) IsCaseExpressionContext() {}

func NewCaseExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaseExpressionContext {
	var p = new(CaseExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_caseExpression

	return p
}

func (s *CaseExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *CaseExpressionContext) CASE() antlr.TerminalNode {
	return s.GetToken(PPLParserCASE, 0)
}

func (s *CaseExpressionContext) END() antlr.TerminalNode {
	return s.GetToken(PPLParserEND, 0)
}

func (s *CaseExpressionContext) AllWhenClause() []IWhenClauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IWhenClauseContext); ok {
			len++
		}
	}

	tst := make([]IWhenClauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IWhenClauseContext); ok {
			tst[i] = t.(IWhenClauseContext)
			i++
		}
	}

	return tst
}

func (s *CaseExpressionContext) WhenClause(i int) IWhenClauseContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhenClauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhenClauseContext)
}

func (s *CaseExpressionContext) ELSE() antlr.TerminalNode {
	return s.GetToken(PPLParserELSE, 0)
}

func (s *CaseExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CaseExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CaseExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CaseExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterCaseExpression(s)
	}
}

func (s *CaseExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitCaseExpression(s)
	}
}

func (s *CaseExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitCaseExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) CaseExpression() (localctx ICaseExpressionContext) {
	localctx = NewCaseExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, PPLParserRULE_caseExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(686)
		p.Match(PPLParserCASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(688)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == PPLParserWHEN {
		{
			p.SetState(687)
			p.WhenClause()
		}

		p.SetState(690)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(694)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PPLParserELSE {
		{
			p.SetState(692)
			p.Match(PPLParserELSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(693)
			p.Expression()
		}

	}
	{
		p.SetState(696)
		p.Match(PPLParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWhenClauseContext is an interface to support dynamic dispatch.
type IWhenClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHEN() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	THEN() antlr.TerminalNode

	// IsWhenClauseContext differentiates from other interfaces.
	IsWhenClauseContext()
}

type WhenClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhenClauseContext() *WhenClauseContext {
	var p = new(WhenClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_whenClause
	return p
}

func InitEmptyWhenClauseContext(p *WhenClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PPLParserRULE_whenClause
}

func (*WhenClauseContext) IsWhenClauseContext() {}

func NewWhenClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhenClauseContext {
	var p = new(WhenClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PPLParserRULE_whenClause

	return p
}

func (s *WhenClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WhenClauseContext) WHEN() antlr.TerminalNode {
	return s.GetToken(PPLParserWHEN, 0)
}

func (s *WhenClauseContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *WhenClauseContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *WhenClauseContext) THEN() antlr.TerminalNode {
	return s.GetToken(PPLParserTHEN, 0)
}

func (s *WhenClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhenClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhenClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.EnterWhenClause(s)
	}
}

func (s *WhenClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PPLParserListener); ok {
		listenerT.ExitWhenClause(s)
	}
}

func (s *WhenClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case PPLParserVisitor:
		return t.VisitWhenClause(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *PPLParser) WhenClause() (localctx IWhenClauseContext) {
	localctx = NewWhenClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, PPLParserRULE_whenClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(698)
		p.Match(PPLParserWHEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(699)
		p.Expression()
	}
	{
		p.SetState(700)
		p.Match(PPLParserTHEN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(701)
		p.Expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
