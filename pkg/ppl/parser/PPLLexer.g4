// Copyright 2024 Quidditch Project
// Licensed under the Apache License, Version 2.0

lexer grammar PPLLexer;

// Keywords (Tier 0 Commands)
SEARCH: 'search' | 'source';
WHERE: 'where';
FIELDS: 'fields';
STATS: 'stats';
SORT: 'sort';
HEAD: 'head';
DESCRIBE: 'describe';
SHOWDATASOURCES: 'showdatasources';
EXPLAIN: 'explain';

// Keywords (Tier 1 Commands)
CHART: 'chart';
TIMECHART: 'timechart';
BIN: 'bin';
DEDUP: 'dedup';
TOP: 'top';
RARE: 'rare';
EVAL: 'eval';
RENAME: 'rename';
REPLACE: 'replace';
FILLNULL: 'fillnull';
PARSE: 'parse';
REX: 'rex';
LOOKUP: 'lookup';
APPEND: 'append';
JOIN: 'join';
TABLE: 'table';
EVENTSTATS: 'eventstats';
STREAMSTATS: 'streamstats';
REVERSE: 'reverse';
FLATTEN: 'flatten';
FILLNULL: 'fillnull';

// Stats keywords
BY: 'by';
AS: 'as';
WITH: 'with';
OUTPUT: 'output';
TYPE: 'type';
VALUE: 'value';

// Join keywords
INNER: 'inner';
LEFT: 'left';
RIGHT: 'right';
OUTER: 'outer';
FULL: 'full';

// Tier 1 command options
SPAN: 'span';
BINS: 'bins';
KEEPEVENTS: 'keepevents';
CONSECUTIVE: 'consecutive';
SORTBY: 'sortby';
LIMIT: 'limit';
COUNTFIELD: 'countfield';
PERCENTFIELD: 'percentfield';
SHOWCOUNT: 'showcount';
SHOWPERC: 'showperc';
USEOTHER: 'useother';
OTHERSTR: 'otherstr';
NULLSTR: 'nullstr';

// Sort order
ASC: 'asc';
DESC: 'desc';

// Logical operators
AND: 'and' | 'AND';
OR: 'or' | 'OR';
NOT: 'not' | 'NOT';

// Comparison operators
EQ: '=';
NEQ: '!=' | '<>';
LT: '<';
LTE: '<=';
GT: '>';
GTE: '>=';
LIKE: 'like' | 'LIKE';
IN: 'in' | 'IN';

// Arithmetic operators
PLUS: '+';
MINUS: '-';
STAR: '*';
SLASH: '/';
PERCENT: '%';

// Literals
TRUE: 'true' | 'TRUE';
FALSE: 'false' | 'FALSE';
NULL: 'null' | 'NULL';

// Case expression
CASE: 'case' | 'CASE';
WHEN: 'when' | 'WHEN';
THEN: 'then' | 'THEN';
ELSE: 'else' | 'ELSE';
END: 'end' | 'END';

// Aggregation functions (Tier 0)
COUNT: 'count' | 'COUNT';
SUM: 'sum' | 'SUM';
AVG: 'avg' | 'AVG';
MIN: 'min' | 'MIN';
MAX: 'max' | 'MAX';

// Aggregation functions (Tier 1)
DC: 'dc' | 'DC';
DISTINCT_COUNT: 'distinct_count' | 'DISTINCT_COUNT';
VAR: 'var' | 'VAR';
VARP: 'varp' | 'VARP';
STDEV: 'stdev' | 'STDEV';
STDEVP: 'stdevp' | 'STDEVP';
PERCENTILE: 'percentile' | 'PERCENTILE' | 'perc';
MEDIAN: 'median' | 'MEDIAN';
MODE: 'mode' | 'MODE';
EARLIEST: 'earliest' | 'EARLIEST' | 'first';
LATEST: 'latest' | 'LATEST' | 'last';
VALUES: 'values' | 'VALUES' | 'list';
RANGE: 'range' | 'RANGE';

// Other keywords
DISTINCT: 'distinct' | 'DISTINCT';

// Delimiters
PIPE: '|';
COMMA: ',';
DOT: '.';
LPAREN: '(';
RPAREN: ')';
LBRACKET: '[';
RBRACKET: ']';

// Identifiers and literals
IDENTIFIER: [a-zA-Z_] [a-zA-Z0-9_]*;
INTEGER: [0-9]+;
DECIMAL: [0-9]+ '.' [0-9]+;
STRING: '\'' (~'\'' | '\\\'')* '\''
      | '"' (~'"' | '\\"')* '"'
      | '`' (~'`' | '\\`')* '`';

// Whitespace
WS: [ \t\r\n]+ -> skip;

// Comments
LINE_COMMENT: '//' ~[\r\n]* -> skip;
BLOCK_COMMENT: '/*' .*? '*/' -> skip;
