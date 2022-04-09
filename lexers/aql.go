package lexers

import (
	. "github.com/alecthomas/chroma/v2" // nolint
)

// ArangoDB Query Language (AQL) lexer.
var Aql = Register(MustNewLexer(
	&Config{
		Name:      "AQL",
		Aliases:   []string{"aql"},
		Filenames: []string{"*.aql"},
		MimeTypes: []string{"application/x-aql"},
		DotAll: true,
		EnsureNL: true,
		CaseInsensitive: true, // except pseudo variables (CURRENT, NEW, OLD)
	},
	AqlRules,
))

const (
	aqlIdentifierPattern = "(?:$?|_+)[a-z]+[_a-z0-9]*"
	aqlBindVariablePattern = "@(?:_+[a-z0-9]+[a-z0-9_]*|[a-z0-9][a-z0-9_]*)"
	aqlUserFunctionsPattern = `[a-zA-Z0-9][a-zA-Z0-9_]*(?:::[a-zA-Z0-9_]+)+(?=\s*\()`
	aqlBuiltinFunctionsPattern = "(?:" +
		"to_bool|to_number|to_string|to_array|to_list|is_null|is_bool|is_number|is_string|is_array|is_list|is_object|is_document|is_datestring|" +
		"typename|json_stringify|json_parse|concat|concat_separator|char_length|lower|upper|substring|left|right|trim|reverse|contains|" +
		"log|log2|log10|exp|exp2|sin|cos|tan|asin|acos|atan|atan2|radians|degrees|pi|regex_test|regex_replace|" +
		"like|floor|ceil|round|abs|rand|sqrt|pow|length|count|min|max|average|avg|sum|product|median|variance_population|variance_sample|variance|" +
		"bit_and|bit_or|bit_xor|bit_negate|bit_test|bit_popcount|bit_shift_left|bit_shift_right|bit_construct|bit_deconstruct|bit_to_string|bit_from_string|" +
		"first|last|unique|outersection|interleave|in_range|jaccard|matches|merge|merge_recursive|has|attributes|values|unset|unset_recursive|keep|keep_recursive|" +
		"near|within|within_rectangle|is_in_polygon|distance|fulltext|stddev_sample|stddev_population|stddev|" +
		"slice|nth|position|contains_array|translate|zip|call|apply|push|append|pop|shift|unshift|remove_value|remove_values|" +
		"remove_nth|replace_nth|date_now|date_timestamp|date_iso8601|date_dayofweek|date_year|date_month|date_day|date_hour|" +
		"date_minute|date_second|date_millisecond|date_dayofyear|date_isoweek|date_leapyear|date_quarter|date_days_in_month|date_trunc|date_round|" +
		"date_add|date_subtract|date_diff|date_compare|date_format|date_utctolocal|date_localtoutc|date_timezone|date_timezones|" +
		"fail|passthru|v8|sleep|schema_get|schema_validate|call_greenspun|version|noopt|noeval|not_null|" +
		"first_list|first_document|parse_identifier|current_user|current_database|collection_count|pregel_result|" +
		"collections|document|decode_rev|range|union|union_distinct|minus|intersection|flatten|is_same_collection|check_document|" +
		"ltrim|rtrim|find_first|find_last|split|substitute|ipv4_to_number|ipv4_from_number|is_ipv4|md5|sha1|sha512|crc32|fnv64|hash|random_token|to_base64|" +
		"to_hex|encode_uri_component|soundex|assert|warn|is_key|sorted|sorted_unique|count_distinct|count_unique|" +
		"levenshtein_distance|levenshtein_match|regex_matches|regex_split|ngram_match|ngram_similarity|ngram_positional_similarity|uuid|" +
		"tokens|exists|starts_with|phrase|min_match|boost|analyzer|" +
		"geo_point|geo_multipoint|geo_polygon|geo_multipolygon|geo_linestring|geo_multilinestring|geo_contains|geo_intersects|" +
		"geo_equals|geo_distance|geo_area|geo_in_range" +
		")(?=\\s*\\()" // Will not recognize function if comment between name and opening parenthesis
)
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

Special punctuation (or operator?): .
Followed by identifier, bind parameter, or bind data source

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

func AqlRules() Rules {
	return Rules{
		"commentsandwhitespace": {
			{`\s+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("multiline-comment")},
		},
		"multiline-comment": {
			{`[^*]+`, CommentMultiline, nil},
			{`\*/`, CommentMultiline, Pop(1)},
			{`\*`, CommentMultiline, nil},
		},
		"double-quote": {
			{`\\.`, LiteralStringDouble, nil},
			{`[^"\\]+`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
		"single-quote": {
			{`\\.`, LiteralStringSingle, nil},
			{`[^'\\]+`, LiteralStringSingle, nil},
			{`'`, LiteralStringSingle, Pop(1)},
		},
		"backtick": {
			{"\\\\.", Name, nil},
			{"[^`\\\\]+", Name, nil},
			{"`", Name, Pop(1)},
		},
		"forwardtick": {
			{"\\\\.", Name, nil},
			{"[^´\\\\]+", Name, nil},
			{"´", Name, Pop(1)},
		},
		"identifier": {
			{aqlIdentifierPattern, Name, nil},
			{"`", Name, Push("backtick")},
			{"´", Name, Push("forwardtick")},
		},
		"bind-variable": {
			{aqlBindVariablePattern, NameVariable, nil},
		},
		"into": {
			Include("commentsandwhitespace"),
			{`KEEP\b`, KeywordPseudo, Pop(1)}, // false positives: INTO keep kEEP, INTO coll LET keep
			Include("identifier"),
			Include("bind-variable"),
			Default(Pop(1)),
		},
		"root": {
			Include("commentsandwhitespace"),
			{`0b[01]+`, LiteralNumberBin, nil},
			{`0x[0-9a-f]+`, LiteralNumberHex, nil},
			{`(?:(?:0|[1-9][0-9]*)(?:\.[0-9]+)?|\.[0-9]+)(?:e[\-\+]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`0|[1-9][0-9]*`, LiteralNumberInteger, nil},
			{`@` + aqlBindVariablePattern, NameVariableGlobal, nil}, // bind data source
			Include("bind-variable"),
			{`[.,(){}\[\]]`, Punctuation, nil},
			{aqlUserFunctionsPattern, NameFunction, nil},
			{`=~|!~|[=!<>]=?|[%?:/*+-]|\.\.|&&|\|\|`, Operator, nil},
			{`(WITH)(\s+)(COUNT)(\s+)(INTO)\b`, ByGroups(KeywordReserved, Text, KeywordPseudo, Text, KeywordReserved), nil},
			//{"(INTO)(\\s+)([`´]?" + aqlIdentifierPattern + "[`´]?)(\\s+)(KEEP)\b", ByGroups(KeywordReserved, Text, Name, Text, KeywordPseudo), nil}, // TODO: bind var? Escaped identifier?
			{`INTO\b`, KeywordReserved, Push("into")},
			//{`IN <name> SEARCH`}, // bind var!
			//{`??? PRUNE`},
			//{`??? TO`},
			{`OPTIONS\s*\{`, KeywordPseudo, nil},
			{`(?:AGGREGATE|ALL|AND|ANY|ASC|COLLECT|DESC|DISTINCT|FILTER|FOR|GRAPH|IN|INBOUND|INSERT|INTO|K_PATHS|K_SHORTEST_PATHS|LIKE|LIMIT|NONE|NOT|OR|OUTBOUND|REMOVE|REPLACE|RETURN|SHORTEST_PATH|SORT|UPDATE|UPSERT|WITH|WINDOW)\b`, KeywordReserved, nil},
			{`LET\b`, KeywordDeclaration, nil}, // also WITH but only at the beginning of a query
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(?-i)(CURRENT|NEW|OLD)\b`, NameBuiltinPseudo, nil},
			{aqlBuiltinFunctionsPattern, NameFunction, nil},
			{`"`, LiteralStringDouble, Push("double-quote")},
			{`'`, LiteralStringSingle, Push("single-quote")},
			Include("identifier"),
		},
	}
}
