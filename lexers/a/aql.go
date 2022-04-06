package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// ArangoDB Query Language (AQL) lexer.
var Aql = internal.Register(MustNewLexer( // nolint: forbidigo
	&Config{
		Name:      "AQL",
		Aliases:   []string{"aql"},
		Filenames: []string{"*.aql"},
		MimeTypes: []string{"application/x-aql"},
		DotAll:    true, // ???
		EnsureNL:  true,
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
	"commentsandwhitespace": {
		{`\s+`, Text, nil},
		{`//.*?\n`, CommentSingle, nil}, // DotAll?
		{`/\*`, CommentMultiline, Push("multiline-comment")},
	},
	"multiline-comment": {
		{`\*/`, CommentMultiline, Pop(1)},
	}
	"slashstartsregex": {
		Include("commentsandwhitespace"),
		{`/(\\.|[^[/\\\n]|\[(\\.|[^\]\\\n])*])+/([gimuy]+\b|\B)`, LiteralStringRegex, Pop(1)},
		{`(?=/)`, Text, Push("#pop", "badregex")},
		Default(Pop(1)),
	},
	"badregex": {
		{`\n`, Text, Pop(1)},
	},
	"root": {
		Include("commentsandwhitespace"),
		{`0|[1-9][0-9]*`, LiteralNumberInteger, nil},
		{`0x[0-9a-f]+`, LiteralNumberHex, nil},
		{`0b[01]+`, LiteralNumberBin, nil},
		{`((0|[1-9][0-9]*)(\.[0-9]+)?|\.[0-9]+)(e[\-\+]?[0-9]+)?`, LiteralNumberFloat, nil},
		{`@@?(_+[a-z0-9]+[a-z0-9_]*|[a-z0-9][a-z0-9_]*)`, NameVariable, nil},
		{`[,(){}\[\]]`, Punctuation, nil},
		{`=~|!~|[=!<>]=?|[%?:/*+-]|::|\.\.|&&|\|\|`, Operator, nil},
		{`(AGGREGATE|ALL|AND|ANY|ASC|COLLECT|DESC|DISTINCT|FALSE|FILTER|FOR|GRAPH|IN|INBOUND|INSERT|INTO|K_PATHS|K_SHORTEST_PATHS|LET|LIKE|LIMIT|NONE|NOT|NULL|OR|OUTBOUND|REMOVE|REPLACE|RETURN|SHORTEST_PATH|SORT|TRUE|UPDATE|UPSERT|WITH|WINDOW)\b`, KeywordReserved, nil},
		{`(true|false|null)\b`, KeywordConstant, nil}, // remove from keyword list?
		// WITH or LET? {`()\b`, KeywordDeclaration, Push("slashstartsregex")},
		// Keyword-like? {`()\b`, Keyword, nil},
		// aql:: or functions or pseudo variables? {`()\b`, NameBuiltin, nil},
		{`(?:[$_\p{L}\p{N}]|\\u[a-fA-F0-9]{4})(?:(?:[$\p{L}\p{N}]|\\u[a-fA-F0-9]{4}))*`, NameOther, nil},
		{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
		{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
		{"`", LiteralStringBacktick, Push("interp")},
	},
	"interp": {
		{"`", LiteralStringBacktick, Pop(1)},
		{`\\\\`, LiteralStringBacktick, nil},
		{"\\\\`", LiteralStringBacktick, nil},
		{"\\\\[^`\\\\]", LiteralStringBacktick, nil},
		{`\$\{`, LiteralStringInterpol, Push("interp-inside")},
		{`\$`, LiteralStringBacktick, nil},
		{"[^`\\\\$]+", LiteralStringBacktick, nil},
	},
	"interp-inside": {
		{`\}`, LiteralStringInterpol, Pop(1)},
		Include("root"),
	},
}
