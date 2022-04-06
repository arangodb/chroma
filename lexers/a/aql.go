package a

import (
	. "github.com/alecthomas/chroma/v2" // nolint
)

// ArangoDB Query Language (AQL) lexer.
var Aql = Register(MustNewLexer( // nolint: forbidigo
	&Config{
		Name:      "AQL",
		Aliases:   []string{"aql"},
		Filenames: []string{"*.aql"},
		MimeTypes: []string{"application/x-aql"},
		DotAll: true,
		EnsureNL: true,
		CaseInsensitive: true,
		// Generally case-insensitive (keywords) but not in case of pseudo variables
	},
	AqlRules,
))

/*
Keywords:
- (AGGREGATE|ALL|AND|ANY|ASC|COLLECT|DESC|DISTINCT|FALSE|FILTER|FOR|GRAPH|IN|INBOUND|INSERT|INTO|K_PATHS|K_SHORTEST_PATHS|LET|LIKE|LIMIT|NONE|NOT|NULL|OR|OUTBOUND|REMOVE|REPLACE|RETURN|SHORTEST_PATH|SORT|TRUE|UPDATE|UPSERT|WITH|WINDOW)\b

Predefined type literals (remove from keyword list?)
- (NULL|TRUE|FALSE)\b

Keyword-like (not reserved):
- KEEP (COLLECT)
- COUNT (WITH COUNT INTO)
- OPTIONS (FOR, SEARCH, COLLECT, INSERT, UPDATE, REPLACE, UPSERT, REMOVE)
- PRUNE (FOR)
- SEARCH
- TO (Shortest Path, k Shortest Paths, k Paths)

Special/pseudo variables (case-sensitive):
- (?-i)(CURRENT|NEW|OLD)\b

Operators: =~ !~ == != >= > <= < = ! && || + - * / % ? :: : ..

Punctuation: , ( ) { } [ ]

Identifiers: (($?|_+)[a-zA-Z]+[_a-zA-Z0-9]*)

String enclosed in backtick or forward tick, double or single quotes (all supporting backslash escapes)

Number literals:
- int decimal: (0|[1-9][0-9]*)
- int binary: (0[bB][01]+)
- int hex: (0[xX][0-9a-fA-F]+)
- double decimal: ((0|[1-9][0-9]*)(\.[0-9]+)?|\.[0-9]+)([eE][\-\+]?[0-9]+)?

Bind parameters:
- @(_+[a-zA-Z0-9]+[a-zA-Z0-9_]*|[a-zA-Z0-9][a-zA-Z0-9_]*)
- @@(_+[a-zA-Z0-9]+[a-zA-Z0-9_]*|[a-zA-Z0-9][a-zA-Z0-9_]*)

Whitespace: [ \t\r\n]+

Comments: single, multi
*/

var AqlRules = Rules{
	"multiline-comment": {
		{`\*/`, CommentMultiline, Pop(1)},
		{`.+`, CommentMultiline, nil},
	},
	"double-quote": {
		{`\\"`, LiteralStringDouble, nil},
		{`"`, LiteralStringDouble, Pop(1)},
		{`.`, LiteralStringDouble, nil},
	},
	"single-quote": {
		{`\\'`, LiteralStringSingle, nil},
		{`'`, LiteralStringSingle, Pop(1)},
		{`.+`, LiteralStringSingle, nil},
	},
	"backtick": {
		{"\\\\`", Name, nil},
		{"`", Name, Pop(1)},
		{`.+`, Name, nil},
	},
	"forwardtick": {
		{"\\\\´", Name, nil},
		{"´", Name, Pop(1)},
		{`.+`, Name, nil},
	},
	"root": {
		{`\s+`, Text, nil},
		{`//.*?\n`, CommentSingle, nil},
		{`/\*`, CommentMultiline, Push("multiline-comment")},
		{`0|[1-9][0-9]*`, LiteralNumberInteger, nil},
		{`0x[0-9a-f]+`, LiteralNumberHex, nil},
		{`0b[01]+`, LiteralNumberBin, nil},
		{`((0|[1-9][0-9]*)(\.[0-9]+)?|\.[0-9]+)(e[\-\+]?[0-9]+)?`, LiteralNumberFloat, nil},
		{`@@?(_+[a-z0-9]+[a-z0-9_]*|[a-z0-9][a-z0-9_]*)`, NameVariable, nil},
		{`[,(){}\[\]]`, Punctuation, nil},
		{`=~|!~|[=!<>]=?|[%?:/*+-]|::|\.\.|&&|\|\|`, Operator, nil},
		{`(AGGREGATE|ALL|AND|ANY|ASC|COLLECT|DESC|DISTINCT|FILTER|FOR|GRAPH|IN|INBOUND|INSERT|INTO|K_PATHS|K_SHORTEST_PATHS|LIKE|LIMIT|NONE|NOT|OR|OUTBOUND|REMOVE|REPLACE|RETURN|SHORTEST_PATH|SORT|UPDATE|UPSERT|WITH|WINDOW)\b`, KeywordReserved, nil},
		{`LET\b`, KeywordDeclaration, nil}, // also WITH but only at beginning of query
		{`(true|false|null)\b`, KeywordConstant, nil},
		// Keyword-like? {`()\b`, Keyword, nil},
		// aql:: or functions or pseudo variables? {`()\b`, NameBuiltin, nil},
		{`"`, LiteralStringDouble, Push("double-quote")},
		{`'`, LiteralStringSingle, Push("single-quote")},
		{"`", Name, Push("backtick")},
		{"´", Name, Push("forwardtick")},
	},
}
