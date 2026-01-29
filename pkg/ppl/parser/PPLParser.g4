// Copyright 2024 Quidditch Project
// Licensed under the Apache License, Version 2.0

parser grammar PPLParser;

options {
    tokenVocab=PPLLexer;
}

// Top-level rule: A query is either a search query, metadata command, or explain query
query
    : explainCommand searchQuery EOF
    | searchQuery EOF
    | metadataCommand EOF
    ;

// Search query: starts with searchCommand, followed by optional processing commands
searchQuery
    : searchCommand (PIPE processingCommand)*
    ;

// Processing commands (cannot be first)
processingCommand
    : whereCommand
    | fieldsCommand
    | statsCommand
    | sortCommand
    | headCommand
    // Tier 1 commands
    | chartCommand
    | timechartCommand
    | binCommand
    | dedupCommand
    | topCommand
    | rareCommand
    | evalCommand
    | renameCommand
    | replaceCommand
    | fillnullCommand
    | parseCommand
    | rexCommand
    | lookupCommand
    | appendCommand
    | joinCommand
    | tableCommand
    | eventstatsCommand
    | streamstatsCommand
    | reverseCommand
    | flattenCommand
    | fillnullCommand
    ;

// Metadata commands (standalone, no search required)
metadataCommand
    : describeCommand
    | showDatasourcesCommand
    ;

// search source=<index> or source=<index>
searchCommand
    : SEARCH SEARCH EQ IDENTIFIER  # SearchWithKeyword
    | SEARCH EQ IDENTIFIER          # SearchWithSource
    ;

// where <expression>
whereCommand
    : WHERE expression
    ;

// fields <field1>, <field2>, ... or fields - <field1>, <field2>
fieldsCommand
    : FIELDS MINUS fieldList       # FieldsExclude
    | FIELDS fieldList             # FieldsInclude
    ;

fieldList
    : expression (COMMA expression)*
    ;

// stats <agg1> [as <alias>], <agg2> [as <alias>] [by <field1>, <field2>]
statsCommand
    : STATS aggregationList (BY fieldList)?
    ;

aggregationList
    : aggregation (COMMA aggregation)*
    ;

aggregation
    : expression (AS IDENTIFIER)?
    ;

// sort <field1> [asc|desc], <field2> [asc|desc]
sortCommand
    : SORT sortFieldList
    ;

sortFieldList
    : sortField (COMMA sortField)*
    ;

sortField
    : expression (ASC | DESC)?
    ;

// head <n>
headCommand
    : HEAD INTEGER
    ;

// ============================================================================
// Tier 1 Commands
// ============================================================================

// chart: Visualization with aggregations
// chart avg(latency) by host
// chart count() by status span=1h
chartCommand
    : CHART aggregationList (BY fieldList)? chartOptions*
    ;

chartOptions
    : SPAN EQ timeSpan
    | LIMIT EQ INTEGER
    | USEOTHER EQ booleanValue
    | OTHERSTR EQ STRING
    | NULLSTR EQ STRING
    ;

// timechart: Time-based aggregation with automatic time bucketing
// timechart span=1h count() by status
// timechart avg(latency), max(latency)
timechartCommand
    : TIMECHART timechartOptions* aggregationList (BY fieldList)?
    ;

timechartOptions
    : SPAN EQ timeSpan
    | BINS EQ INTEGER
    | LIMIT EQ INTEGER
    | USEOTHER EQ booleanValue
    ;

// Time span specification: <number><unit>
// Units: s (seconds), m (minutes), h (hours), d (days), w (weeks), mon (months)
timeSpan
    : INTEGER IDENTIFIER  // e.g., 1h, 5m, 1d
    ;

// bin: Create time or numeric buckets
// bin timestamp span=1h
// bin latency bins=10
// bin timestamp span=auto
binCommand
    : BIN fieldReference binOptions*
    ;

binOptions
    : SPAN EQ (timeSpan | IDENTIFIER)  // IDENTIFIER allows 'auto'
    | BINS EQ INTEGER
    ;

// dedup: Remove duplicate rows
// dedup host
// dedup 5 host, status  // Keep first 5 duplicates
// dedup host keepevents=true
dedupCommand
    : DEDUP INTEGER? fieldList dedupOptions*
    ;

dedupOptions
    : KEEPEVENTS EQ booleanValue
    | CONSECUTIVE EQ booleanValue
    | SORTBY sortFieldList
    ;

// top: Get most common values
// top host
// top 10 host, status
// top host by datacenter
topCommand
    : TOP INTEGER? fieldList (BY fieldList)? topOptions*
    ;

topOptions
    : COUNTFIELD EQ STRING
    | PERCENTFIELD EQ STRING
    | SHOWCOUNT EQ booleanValue
    | SHOWPERC EQ booleanValue
    | LIMIT EQ INTEGER
    | USEOTHER EQ booleanValue
    | OTHERSTR EQ STRING
    ;

// rare: Get least common values (opposite of top)
// rare host
// rare 10 host, status
// rare host by datacenter
rareCommand
    : RARE INTEGER? fieldList (BY fieldList)? topOptions*  // Reuses topOptions
    ;

// eval: Compute new fields from expressions
// eval duration_sec = duration / 1000
// eval is_error = status >= 400
evalCommand
    : EVAL evalAssignment (COMMA evalAssignment)*
    ;

evalAssignment
    : IDENTIFIER EQ expression
    ;

// rename: Rename fields
// rename old_field as new_field
// rename field1 as alias1, field2 as alias2
renameCommand
    : RENAME renameAssignment (COMMA renameAssignment)*
    ;

renameAssignment
    : IDENTIFIER AS IDENTIFIER
    ;

// replace: Replace values in a field
// replace oldval1 with newval1, oldval2 with newval2 in field
replaceCommand
    : REPLACE replaceMapping (COMMA replaceMapping)* IN IDENTIFIER
    ;

replaceMapping
    : expression WITH expression
    ;

// fillnull: Fill null/missing values in fields
// fillnull field1=value1, field2=value2 or fillnull value=default fields field1, field2
fillnullCommand
    : FILLNULL fillnullAssignment (COMMA fillnullAssignment)*
    ;

fillnullAssignment
    : IDENTIFIER EQ expression
    ;

// parse: Extract fields from text using regex patterns with named captures
// parse message "(?<user>\w+) logged in from (?<ip>\d+\.\d+\.\d+\.\d+)"
// parse field=message "user=(?<user>\w+) action=(?<action>\w+)"
parseCommand
    : PARSE (IDENTIFIER EQ)? IDENTIFIER STRING
    ;

// rex: Extract fields using regular expressions
// rex "(?<error_code>\d{3}): (?<error_msg>.*)"
// rex field=message "user=(?<user>\w+)"
rexCommand
    : REX (IDENTIFIER EQ IDENTIFIER)? STRING
    ;

// lookup: Enrich data with external lookup tables
// lookup products product_id OUTPUT name, price
// lookup users user_id AS uid OUTPUT username AS user
lookupCommand
    : LOOKUP IDENTIFIER IDENTIFIER (AS IDENTIFIER)? OUTPUT lookupOutputList
    ;

lookupOutputList
    : lookupOutputField (COMMA lookupOutputField)*
    ;

lookupOutputField
    : IDENTIFIER (AS IDENTIFIER)?
    ;

// append: Concatenate results from a subsearch
// source=logs_2024 | append [source=logs_2023]
appendCommand
    : APPEND LBRACKET searchQuery RBRACKET
    ;

// join: Combine datasets with SQL-like joins
// source=orders | join user_id [source=users]
// source=orders | join type=left user_id [source=users]
joinCommand
    : JOIN (TYPE EQ joinType)? IDENTIFIER LBRACKET searchQuery RBRACKET
    ;

joinType
    : INNER | LEFT | RIGHT | OUTER | FULL
    ;

// table: Select and order specific columns
// table host, status, latency
tableCommand
    : TABLE fieldList
    ;

// eventstats: Compute running statistics across all events
// eventstats avg(latency) as avg_lat by host
eventstatsCommand
    : EVENTSTATS aggregationList (BY fieldList)?
    ;

// streamstats: Compute running statistics in streaming fashion
// streamstats count() as running_count
// streamstats window=10 avg(latency) as rolling_avg
streamstatsCommand
    : STREAMSTATS streamstatsOptions* aggregationList (BY fieldList)?
    ;

streamstatsOptions
    : IDENTIFIER EQ (INTEGER | STRING | booleanValue)  // Generic key=value options
    ;

// reverse (no parameters)
reverseCommand
    : REVERSE
    ;

// flatten: Flatten nested arrays/objects into separate rows
// flatten <field>
flattenCommand
    : FLATTEN fieldReference
    ;

// fillnull: Fill NULL/missing values with a default value
// fillnull value=<value> [fields <field_list>]
fillnullCommand
    : FILLNULL VALUE EQ literalValue (FIELDS fieldList)?
    ;

// Boolean value helper
booleanValue
    : TRUE
    | FALSE
    ;

// describe <source>
describeCommand
    : DESCRIBE IDENTIFIER
    ;

// showdatasources
showDatasourcesCommand
    : SHOWDATASOURCES
    ;

// explain
explainCommand
    : EXPLAIN
    ;

// Expressions (ordered by precedence, lowest to highest)

expression
    : orExpression
    ;

// OR has lowest precedence
orExpression
    : andExpression (OR andExpression)*
    ;

// AND has higher precedence than OR
andExpression
    : notExpression (AND notExpression)*
    ;

// NOT has higher precedence than AND
notExpression
    : NOT notExpression
    | comparisonExpression
    ;

// Comparison operators
comparisonExpression
    : additiveExpression
      ( (EQ | NEQ | LT | LTE | GT | GTE) additiveExpression
      | LIKE additiveExpression
      | IN LPAREN expressionList RPAREN
      )?
    ;

// Additive operators (+ -)
additiveExpression
    : multiplicativeExpression ((PLUS | MINUS) multiplicativeExpression)*
    ;

// Multiplicative operators (* / %)
multiplicativeExpression
    : unaryExpression ((STAR | SLASH | PERCENT) unaryExpression)*
    ;

// Unary operators (- +)
unaryExpression
    : (PLUS | MINUS) unaryExpression
    | primaryExpression
    ;

// Primary expressions
primaryExpression
    : literal
    | fieldReference
    | functionCall
    | caseExpression
    | LPAREN expression RPAREN
    ;

// Literals
literal
    : INTEGER
    | DECIMAL
    | STRING
    | TRUE
    | FALSE
    | NULL
    ;

// Field references (can include dots for nested fields)
fieldReference
    : IDENTIFIER (DOT IDENTIFIER)*
    | IDENTIFIER LBRACKET INTEGER RBRACKET  // Array indexing
    ;

// Function calls
functionCall
    : IDENTIFIER LPAREN RPAREN                                    # FunctionCallNoArgs
    | IDENTIFIER LPAREN DISTINCT? expressionList RPAREN           # FunctionCallWithArgs
    | aggregationFunction LPAREN RPAREN                           # AggregationFunctionCallNoArgs
    | aggregationFunction LPAREN DISTINCT? expressionList RPAREN  # AggregationFunctionCall
    ;

// Aggregation functions (Tier 0 + Tier 1)
aggregationFunction
    : COUNT
    | SUM
    | AVG
    | MIN
    | MAX
    // Tier 1 aggregation functions
    | DC
    | DISTINCT_COUNT
    | VAR
    | VARP
    | STDEV
    | STDEVP
    | PERCENTILE
    | MEDIAN
    | MODE
    | EARLIEST
    | LATEST
    | VALUES
    | RANGE
    ;

expressionList
    : expression (COMMA expression)*
    ;

// Case expression
caseExpression
    : CASE whenClause+ (ELSE expression)? END
    ;

whenClause
    : WHEN expression THEN expression
    ;
