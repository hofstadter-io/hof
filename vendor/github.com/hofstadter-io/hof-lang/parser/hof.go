// Package 'parser' implements the basic parser for 'hof-lang'
// i.e. there is no ast or validation
package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hofstadter-io/hof-lang/ast"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case []interface{}:
		return v.([]interface{})
	default:
		return []interface{}{v}
	}
}

var g = &grammar{
	rules: []*rule{
		{
			name: "HOF",
			pos:  position{line: 25, col: 1, offset: 398},
			expr: &actionExpr{
				pos: position{line: 25, col: 8, offset: 405},
				run: (*parser).callonHOF1,
				expr: &seqExpr{
					pos: position{line: 25, col: 8, offset: 405},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 25, col: 8, offset: 405},
							label: "defs",
							expr: &ruleRefExpr{
								pos:  position{line: 25, col: 13, offset: 410},
								name: "Definitions",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 25, col: 25, offset: 422},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Definitions",
			pos:  position{line: 33, col: 1, offset: 514},
			expr: &actionExpr{
				pos: position{line: 33, col: 16, offset: 529},
				run: (*parser).callonDefinitions1,
				expr: &seqExpr{
					pos: position{line: 33, col: 16, offset: 529},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 33, col: 16, offset: 529},
							label: "defs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 33, col: 21, offset: 534},
								expr: &ruleRefExpr{
									pos:  position{line: 33, col: 23, offset: 536},
									name: "Definition",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 33, col: 37, offset: 550},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "Definition",
			pos:  position{line: 45, col: 1, offset: 716},
			expr: &actionExpr{
				pos: position{line: 45, col: 15, offset: 730},
				run: (*parser).callonDefinition1,
				expr: &seqExpr{
					pos: position{line: 45, col: 15, offset: 730},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 45, col: 15, offset: 730},
							name: "__",
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 18, offset: 733},
							name: "DEF",
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 22, offset: 737},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 45, col: 24, offset: 739},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 29, offset: 744},
								name: "ID",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 32, offset: 747},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 45, col: 34, offset: 749},
							label: "dsl",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 38, offset: 753},
								name: "DSL",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 42, offset: 757},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 45, col: 44, offset: 759},
							label: "body",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 49, offset: 764},
								name: "DefnBody",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 58, offset: 773},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "DefnBody",
			pos:  position{line: 55, col: 1, offset: 907},
			expr: &actionExpr{
				pos: position{line: 55, col: 13, offset: 919},
				run: (*parser).callonDefnBody1,
				expr: &seqExpr{
					pos: position{line: 55, col: 13, offset: 919},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 55, col: 13, offset: 919},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 17, offset: 923},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 55, col: 20, offset: 926},
							label: "defs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 55, col: 25, offset: 931},
								expr: &ruleRefExpr{
									pos:  position{line: 55, col: 27, offset: 933},
									name: "DefnField",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 40, offset: 946},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 55, col: 43, offset: 949},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 55, col: 47, offset: 953},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "DefnField",
			pos:  position{line: 68, col: 1, offset: 1114},
			expr: &actionExpr{
				pos: position{line: 68, col: 14, offset: 1127},
				run: (*parser).callonDefnField1,
				expr: &seqExpr{
					pos: position{line: 68, col: 14, offset: 1127},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 68, col: 14, offset: 1127},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 68, col: 17, offset: 1130},
							label: "val",
							expr: &choiceExpr{
								pos: position{line: 68, col: 23, offset: 1136},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 68, col: 23, offset: 1136},
										name: "TypeDecl",
									},
									&ruleRefExpr{
										pos:  position{line: 68, col: 34, offset: 1147},
										name: "Field",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 68, col: 42, offset: 1155},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "TypeDecl",
			pos:  position{line: 72, col: 1, offset: 1183},
			expr: &actionExpr{
				pos: position{line: 72, col: 13, offset: 1195},
				run: (*parser).callonTypeDecl1,
				expr: &seqExpr{
					pos: position{line: 72, col: 13, offset: 1195},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 72, col: 13, offset: 1195},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 72, col: 15, offset: 1197},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 18, offset: 1200},
								name: "ID",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 72, col: 21, offset: 1203},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 72, col: 23, offset: 1205},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 72, col: 27, offset: 1209},
								name: "TYPE",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 72, col: 32, offset: 1214},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 72, col: 34, offset: 1216},
							label: "obj",
							expr: &zeroOrOneExpr{
								pos: position{line: 72, col: 38, offset: 1220},
								expr: &ruleRefExpr{
									pos:  position{line: 72, col: 38, offset: 1220},
									name: "Object",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 72, col: 46, offset: 1228},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 85, col: 1, offset: 1400},
			expr: &actionExpr{
				pos: position{line: 85, col: 10, offset: 1409},
				run: (*parser).callonValue1,
				expr: &labeledExpr{
					pos:   position{line: 85, col: 10, offset: 1409},
					label: "val",
					expr: &choiceExpr{
						pos: position{line: 85, col: 16, offset: 1415},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 85, col: 16, offset: 1415},
								name: "Boolean",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 26, offset: 1425},
								name: "TypeDecl",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 37, offset: 1436},
								name: "Object",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 46, offset: 1445},
								name: "Array",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 54, offset: 1453},
								name: "Number",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 63, offset: 1462},
								name: "Integer",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 73, offset: 1472},
								name: "String",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 82, offset: 1481},
								name: "TypeRef",
							},
							&ruleRefExpr{
								pos:  position{line: 85, col: 92, offset: 1491},
								name: "Null",
							},
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 89, col: 1, offset: 1523},
			expr: &actionExpr{
				pos: position{line: 89, col: 10, offset: 1532},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 89, col: 10, offset: 1532},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 89, col: 10, offset: 1532},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 89, col: 13, offset: 1535},
							label: "tok",
							expr: &ruleRefExpr{
								pos:  position{line: 89, col: 17, offset: 1539},
								name: "Token",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 23, offset: 1545},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 89, col: 25, offset: 1547},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 29, offset: 1551},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 89, col: 31, offset: 1553},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 89, col: 35, offset: 1557},
								name: "Value",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 41, offset: 1563},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "Object",
			pos:  position{line: 98, col: 1, offset: 1677},
			expr: &actionExpr{
				pos: position{line: 98, col: 11, offset: 1687},
				run: (*parser).callonObject1,
				expr: &seqExpr{
					pos: position{line: 98, col: 11, offset: 1687},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 98, col: 11, offset: 1687},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 98, col: 15, offset: 1691},
							label: "fields",
							expr: &zeroOrMoreExpr{
								pos: position{line: 98, col: 22, offset: 1698},
								expr: &ruleRefExpr{
									pos:  position{line: 98, col: 24, offset: 1700},
									name: "Field",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 98, col: 33, offset: 1709},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 98, col: 36, offset: 1712},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Elem",
			pos:  position{line: 109, col: 1, offset: 1930},
			expr: &actionExpr{
				pos: position{line: 109, col: 9, offset: 1938},
				run: (*parser).callonElem1,
				expr: &seqExpr{
					pos: position{line: 109, col: 9, offset: 1938},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 109, col: 9, offset: 1938},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 109, col: 12, offset: 1941},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 109, col: 16, offset: 1945},
								name: "Value",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 109, col: 22, offset: 1951},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 109, col: 24, offset: 1953},
							val:        ",",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 109, col: 28, offset: 1957},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "Array",
			pos:  position{line: 113, col: 1, offset: 1982},
			expr: &actionExpr{
				pos: position{line: 113, col: 10, offset: 1991},
				run: (*parser).callonArray1,
				expr: &seqExpr{
					pos: position{line: 113, col: 10, offset: 1991},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 113, col: 10, offset: 1991},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 113, col: 14, offset: 1995},
							label: "elems",
							expr: &zeroOrMoreExpr{
								pos: position{line: 113, col: 20, offset: 2001},
								expr: &ruleRefExpr{
									pos:  position{line: 113, col: 22, offset: 2003},
									name: "Elem",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 113, col: 30, offset: 2011},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 113, col: 33, offset: 2014},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CodeBlock",
			pos:  position{line: 122, col: 1, offset: 2235},
			expr: &actionExpr{
				pos: position{line: 122, col: 14, offset: 2248},
				run: (*parser).callonCodeBlock1,
				expr: &seqExpr{
					pos: position{line: 122, col: 14, offset: 2248},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 122, col: 14, offset: 2248},
							name: "Code",
						},
						&ruleRefExpr{
							pos:  position{line: 122, col: 19, offset: 2253},
							name: "__",
						},
					},
				},
			},
		},
		{
			name: "Code",
			pos:  position{line: 127, col: 1, offset: 2303},
			expr: &zeroOrMoreExpr{
				pos: position{line: 127, col: 9, offset: 2311},
				expr: &choiceExpr{
					pos: position{line: 127, col: 11, offset: 2313},
					alternatives: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 127, col: 11, offset: 2313},
							expr: &seqExpr{
								pos: position{line: 127, col: 13, offset: 2315},
								exprs: []interface{}{
									&notExpr{
										pos: position{line: 127, col: 13, offset: 2315},
										expr: &charClassMatcher{
											pos:        position{line: 127, col: 14, offset: 2316},
											val:        "[{}]",
											chars:      []rune{'{', '}'},
											ignoreCase: false,
											inverted:   false,
										},
									},
									&anyMatcher{
										line: 127, col: 19, offset: 2321,
									},
								},
							},
						},
						&seqExpr{
							pos: position{line: 127, col: 26, offset: 2328},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 127, col: 26, offset: 2328},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 127, col: 30, offset: 2332},
									name: "Code",
								},
								&litMatcher{
									pos:        position{line: 127, col: 35, offset: 2337},
									val:        "}",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Number",
			pos:  position{line: 129, col: 1, offset: 2345},
			expr: &actionExpr{
				pos: position{line: 129, col: 11, offset: 2355},
				run: (*parser).callonNumber1,
				expr: &seqExpr{
					pos: position{line: 129, col: 11, offset: 2355},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 129, col: 11, offset: 2355},
							expr: &litMatcher{
								pos:        position{line: 129, col: 11, offset: 2355},
								val:        "-",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 129, col: 16, offset: 2360},
							name: "Integer",
						},
						&litMatcher{
							pos:        position{line: 129, col: 24, offset: 2368},
							val:        ".",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 129, col: 28, offset: 2372},
							expr: &ruleRefExpr{
								pos:  position{line: 129, col: 28, offset: 2372},
								name: "DecimalDigit",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 129, col: 42, offset: 2386},
							expr: &ruleRefExpr{
								pos:  position{line: 129, col: 42, offset: 2386},
								name: "Exponent",
							},
						},
					},
				},
			},
		},
		{
			name: "Index",
			pos:  position{line: 142, col: 1, offset: 2629},
			expr: &actionExpr{
				pos: position{line: 142, col: 10, offset: 2638},
				run: (*parser).callonIndex1,
				expr: &ruleRefExpr{
					pos:  position{line: 142, col: 10, offset: 2638},
					name: "Integer",
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 148, col: 1, offset: 2778},
			expr: &choiceExpr{
				pos: position{line: 148, col: 12, offset: 2789},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 148, col: 12, offset: 2789},
						val:        "0",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 148, col: 18, offset: 2795},
						run: (*parser).callonInteger3,
						expr: &seqExpr{
							pos: position{line: 148, col: 18, offset: 2795},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 148, col: 18, offset: 2795},
									name: "NonZeroDecimalDigit",
								},
								&zeroOrMoreExpr{
									pos: position{line: 148, col: 38, offset: 2815},
									expr: &ruleRefExpr{
										pos:  position{line: 148, col: 38, offset: 2815},
										name: "DecimalDigit",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Exponent",
			pos:  position{line: 161, col: 1, offset: 3057},
			expr: &seqExpr{
				pos: position{line: 161, col: 13, offset: 3069},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 161, col: 13, offset: 3069},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 161, col: 18, offset: 3074},
						expr: &charClassMatcher{
							pos:        position{line: 161, col: 18, offset: 3074},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 161, col: 24, offset: 3080},
						expr: &ruleRefExpr{
							pos:  position{line: 161, col: 24, offset: 3080},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 163, col: 1, offset: 3095},
			expr: &actionExpr{
				pos: position{line: 163, col: 11, offset: 3105},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 163, col: 11, offset: 3105},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 163, col: 11, offset: 3105},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 163, col: 15, offset: 3109},
							expr: &choiceExpr{
								pos: position{line: 163, col: 17, offset: 3111},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 163, col: 17, offset: 3111},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 163, col: 17, offset: 3111},
												expr: &ruleRefExpr{
													pos:  position{line: 163, col: 18, offset: 3112},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 163, col: 30, offset: 3124,
											},
										},
									},
									&seqExpr{
										pos: position{line: 163, col: 34, offset: 3128},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 163, col: 34, offset: 3128},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 163, col: 39, offset: 3133},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 163, col: 57, offset: 3151},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AlphaNumeric",
			pos:  position{line: 177, col: 1, offset: 3458},
			expr: &choiceExpr{
				pos: position{line: 177, col: 17, offset: 3474},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 177, col: 17, offset: 3474},
						name: "Alphabetic",
					},
					&ruleRefExpr{
						pos:  position{line: 177, col: 30, offset: 3487},
						name: "DecimalDigit",
					},
				},
			},
		},
		{
			name: "Alphabetic",
			pos:  position{line: 179, col: 1, offset: 3501},
			expr: &charClassMatcher{
				pos:        position{line: 179, col: 15, offset: 3515},
				val:        "[a-zA-Z]",
				ranges:     []rune{'a', 'z', 'A', 'Z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 181, col: 1, offset: 3525},
			expr: &charClassMatcher{
				pos:        position{line: 181, col: 16, offset: 3540},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 183, col: 1, offset: 3556},
			expr: &choiceExpr{
				pos: position{line: 183, col: 19, offset: 3574},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 183, col: 19, offset: 3574},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 183, col: 38, offset: 3593},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 185, col: 1, offset: 3608},
			expr: &charClassMatcher{
				pos:        position{line: 185, col: 21, offset: 3628},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 187, col: 1, offset: 3641},
			expr: &seqExpr{
				pos: position{line: 187, col: 18, offset: 3658},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 187, col: 18, offset: 3658},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 187, col: 22, offset: 3662},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 187, col: 31, offset: 3671},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 187, col: 40, offset: 3680},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 187, col: 49, offset: 3689},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 189, col: 1, offset: 3699},
			expr: &charClassMatcher{
				pos:        position{line: 189, col: 17, offset: 3715},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 191, col: 1, offset: 3722},
			expr: &charClassMatcher{
				pos:        position{line: 191, col: 24, offset: 3745},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 193, col: 1, offset: 3752},
			expr: &charClassMatcher{
				pos:        position{line: 193, col: 13, offset: 3764},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "Boolean",
			pos:  position{line: 195, col: 1, offset: 3775},
			expr: &choiceExpr{
				pos: position{line: 195, col: 12, offset: 3786},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 195, col: 12, offset: 3786},
						run: (*parser).callonBoolean2,
						expr: &litMatcher{
							pos:        position{line: 195, col: 12, offset: 3786},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 195, col: 62, offset: 3836},
						run: (*parser).callonBoolean4,
						expr: &litMatcher{
							pos:        position{line: 195, col: 62, offset: 3836},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 199, col: 1, offset: 3918},
			expr: &actionExpr{
				pos: position{line: 199, col: 9, offset: 3926},
				run: (*parser).callonNull1,
				expr: &litMatcher{
					pos:        position{line: 199, col: 9, offset: 3926},
					val:        "null",
					ignoreCase: false,
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 201, col: 1, offset: 3954},
			expr: &zeroOrMoreExpr{
				pos: position{line: 201, col: 6, offset: 3961},
				expr: &choiceExpr{
					pos: position{line: 201, col: 8, offset: 3963},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 201, col: 8, offset: 3963},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 201, col: 21, offset: 3976},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 201, col: 27, offset: 3982},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 202, col: 1, offset: 3993},
			expr: &zeroOrMoreExpr{
				pos: position{line: 202, col: 5, offset: 3999},
				expr: &choiceExpr{
					pos: position{line: 202, col: 7, offset: 4001},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 202, col: 7, offset: 4001},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 20, offset: 4014},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 204, col: 1, offset: 4051},
			expr: &charClassMatcher{
				pos:        position{line: 204, col: 14, offset: 4066},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 205, col: 1, offset: 4074},
			expr: &litMatcher{
				pos:        position{line: 205, col: 7, offset: 4082},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 206, col: 1, offset: 4087},
			expr: &choiceExpr{
				pos: position{line: 206, col: 7, offset: 4095},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 206, col: 7, offset: 4095},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 206, col: 7, offset: 4095},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 206, col: 10, offset: 4098},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 206, col: 16, offset: 4104},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 206, col: 16, offset: 4104},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 206, col: 18, offset: 4106},
								expr: &ruleRefExpr{
									pos:  position{line: 206, col: 18, offset: 4106},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 206, col: 37, offset: 4125},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 206, col: 43, offset: 4131},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 206, col: 43, offset: 4131},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 206, col: 46, offset: 4134},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 208, col: 1, offset: 4139},
			expr: &notExpr{
				pos: position{line: 208, col: 7, offset: 4147},
				expr: &anyMatcher{
					line: 208, col: 8, offset: 4148,
				},
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 211, col: 1, offset: 4152},
			expr: &anyMatcher{
				line: 211, col: 14, offset: 4167,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 212, col: 1, offset: 4169},
			expr: &choiceExpr{
				pos: position{line: 212, col: 11, offset: 4181},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 212, col: 11, offset: 4181},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 30, offset: 4200},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 213, col: 1, offset: 4218},
			expr: &seqExpr{
				pos: position{line: 213, col: 20, offset: 4239},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 213, col: 20, offset: 4239},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 213, col: 25, offset: 4244},
						expr: &seqExpr{
							pos: position{line: 213, col: 27, offset: 4246},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 213, col: 27, offset: 4246},
									expr: &litMatcher{
										pos:        position{line: 213, col: 28, offset: 4247},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 213, col: 33, offset: 4252},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 213, col: 47, offset: 4266},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 214, col: 1, offset: 4271},
			expr: &seqExpr{
				pos: position{line: 214, col: 36, offset: 4308},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 214, col: 36, offset: 4308},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 214, col: 41, offset: 4313},
						expr: &seqExpr{
							pos: position{line: 214, col: 43, offset: 4315},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 214, col: 43, offset: 4315},
									expr: &choiceExpr{
										pos: position{line: 214, col: 46, offset: 4318},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 214, col: 46, offset: 4318},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 214, col: 53, offset: 4325},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 214, col: 59, offset: 4331},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 214, col: 73, offset: 4345},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 215, col: 1, offset: 4350},
			expr: &seqExpr{
				pos: position{line: 215, col: 21, offset: 4372},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 215, col: 21, offset: 4372},
						val:        "//",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 215, col: 26, offset: 4377},
						expr: &seqExpr{
							pos: position{line: 215, col: 28, offset: 4379},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 215, col: 28, offset: 4379},
									expr: &ruleRefExpr{
										pos:  position{line: 215, col: 29, offset: 4380},
										name: "EOL",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 215, col: 33, offset: 4384},
									name: "SourceChar",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DEF",
			pos:  position{line: 217, col: 1, offset: 4399},
			expr: &litMatcher{
				pos:        position{line: 217, col: 8, offset: 4406},
				val:        "def",
				ignoreCase: false,
			},
		},
		{
			name: "TypeRef",
			pos:  position{line: 219, col: 1, offset: 4413},
			expr: &actionExpr{
				pos: position{line: 219, col: 12, offset: 4424},
				run: (*parser).callonTypeRef1,
				expr: &labeledExpr{
					pos:   position{line: 219, col: 12, offset: 4424},
					label: "val",
					expr: &seqExpr{
						pos: position{line: 219, col: 18, offset: 4430},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 219, col: 18, offset: 4430},
								name: "Alphabetic",
							},
							&zeroOrMoreExpr{
								pos: position{line: 219, col: 29, offset: 4441},
								expr: &choiceExpr{
									pos: position{line: 219, col: 31, offset: 4443},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 219, col: 31, offset: 4443},
											name: "TokenCharactors",
										},
										&charClassMatcher{
											pos:        position{line: 219, col: 49, offset: 4461},
											val:        "[.]",
											chars:      []rune{'.'},
											ignoreCase: false,
											inverted:   false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TYPE",
			pos:  position{line: 227, col: 1, offset: 4566},
			expr: &actionExpr{
				pos: position{line: 227, col: 9, offset: 4574},
				run: (*parser).callonTYPE1,
				expr: &labeledExpr{
					pos:   position{line: 227, col: 9, offset: 4574},
					label: "val",
					expr: &seqExpr{
						pos: position{line: 227, col: 15, offset: 4580},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 227, col: 15, offset: 4580},
								name: "Alphabetic",
							},
							&zeroOrMoreExpr{
								pos: position{line: 227, col: 26, offset: 4591},
								expr: &ruleRefExpr{
									pos:  position{line: 227, col: 26, offset: 4591},
									name: "AlphaNumeric",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ID",
			pos:  position{line: 235, col: 1, offset: 4703},
			expr: &actionExpr{
				pos: position{line: 235, col: 7, offset: 4709},
				run: (*parser).callonID1,
				expr: &labeledExpr{
					pos:   position{line: 235, col: 7, offset: 4709},
					label: "val",
					expr: &seqExpr{
						pos: position{line: 235, col: 13, offset: 4715},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 235, col: 13, offset: 4715},
								name: "Alphabetic",
							},
							&zeroOrMoreExpr{
								pos: position{line: 235, col: 24, offset: 4726},
								expr: &ruleRefExpr{
									pos:  position{line: 235, col: 24, offset: 4726},
									name: "AlphaNumeric",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DSL",
			pos:  position{line: 243, col: 1, offset: 4838},
			expr: &actionExpr{
				pos: position{line: 243, col: 8, offset: 4845},
				run: (*parser).callonDSL1,
				expr: &labeledExpr{
					pos:   position{line: 243, col: 8, offset: 4845},
					label: "val",
					expr: &seqExpr{
						pos: position{line: 243, col: 14, offset: 4851},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 243, col: 14, offset: 4851},
								name: "Alphabetic",
							},
							&zeroOrMoreExpr{
								pos: position{line: 243, col: 25, offset: 4862},
								expr: &choiceExpr{
									pos: position{line: 243, col: 27, offset: 4864},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 243, col: 27, offset: 4864},
											name: "AlphaNumeric",
										},
										&charClassMatcher{
											pos:        position{line: 243, col: 42, offset: 4879},
											val:        "[/]",
											chars:      []rune{'/'},
											ignoreCase: false,
											inverted:   false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Token",
			pos:  position{line: 251, col: 1, offset: 4984},
			expr: &actionExpr{
				pos: position{line: 251, col: 10, offset: 4993},
				run: (*parser).callonToken1,
				expr: &labeledExpr{
					pos:   position{line: 251, col: 10, offset: 4993},
					label: "val",
					expr: &seqExpr{
						pos: position{line: 251, col: 16, offset: 4999},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 251, col: 16, offset: 4999},
								name: "Alphabetic",
							},
							&zeroOrMoreExpr{
								pos: position{line: 251, col: 27, offset: 5010},
								expr: &ruleRefExpr{
									pos:  position{line: 251, col: 27, offset: 5010},
									name: "AlphaNumeric",
								},
							},
						},
					},
				},
			},
		},
	},
}

func (c *current) onHOF1(defs interface{}) (interface{}, error) {
	ret := ast.HofFile{
		Definitions: defs.([]ast.Definition),
	}

	return ret, nil
}

func (p *parser) callonHOF1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHOF1(stack["defs"])
}

func (c *current) onDefinitions1(defs interface{}) (interface{}, error) {

	ret := []ast.Definition{}
	vals := toIfaceSlice(defs)

	for _, def := range vals {
		ret = append(ret, def.(ast.Definition))
	}

	return ret, nil
}

func (p *parser) callonDefinitions1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDefinitions1(stack["defs"])
}

func (c *current) onDefinition1(name, dsl, body interface{}) (interface{}, error) {
	ret := ast.Definition{
		Name: name.(ast.Token),
		DSL:  dsl.(ast.Token),
		Body: body.([]ast.ASTNode),
	}

	return ret, nil
}

func (p *parser) callonDefinition1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDefinition1(stack["name"], stack["dsl"], stack["body"])
}

func (c *current) onDefnBody1(defs interface{}) (interface{}, error) {

	ret := []ast.ASTNode{}

	vals := toIfaceSlice(defs)

	for _, val := range vals {
		ret = append(ret, val.(ast.ASTNode))
	}

	return ret, nil
}

func (p *parser) callonDefnBody1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDefnBody1(stack["defs"])
}

func (c *current) onDefnField1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonDefnField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDefnField1(stack["val"])
}

func (c *current) onTypeDecl1(id, typ, obj interface{}) (interface{}, error) {
	ret := ast.TypeDecl{
		Name: id.(ast.Token),
		Type: typ.(ast.Token),
	}
	if obj != nil {
		objVal := obj.(ast.Object)
		ret.Extra = &objVal
	}

	return ret, nil
}

func (p *parser) callonTypeDecl1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeDecl1(stack["id"], stack["typ"], stack["obj"])
}

func (c *current) onValue1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onValue1(stack["val"])
}

func (c *current) onField1(tok, val interface{}) (interface{}, error) {
	ret := ast.Field{
		Key:   tok.(ast.Token),
		Value: val.(ast.ASTNode),
	}

	return ret, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["tok"], stack["val"])
}

func (c *current) onObject1(fields interface{}) (interface{}, error) {
	vals := toIfaceSlice(fields)
	ret := ast.Object{Fields: make([]ast.Field, 0, len(vals))}

	for _, val := range vals {
		ret.Fields = append(ret.Fields, val.(ast.Field))
	}

	return ret, nil
}

func (p *parser) callonObject1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObject1(stack["fields"])
}

func (c *current) onElem1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonElem1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onElem1(stack["val"])
}

func (c *current) onArray1(elems interface{}) (interface{}, error) {
	vals := toIfaceSlice(elems)
	ret := ast.Array{Elems: make([]ast.ASTNode, 0, len(vals))}
	for _, val := range vals {
		ret.Elems = append(ret.Elems, val.(ast.ASTNode))
	}
	return ret, nil
}

func (p *parser) callonArray1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArray1(stack["elems"])
}

func (c *current) onCodeBlock1() (interface{}, error) {
	text := string(c.text)
	return text, nil
}

func (p *parser) callonCodeBlock1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCodeBlock1()
}

func (c *current) onNumber1() (interface{}, error) {
	// JSON numbers have the same syntax as Go's, and are parseable using
	// strconv.
	val, err := strconv.ParseFloat(string(c.text), 64)
	if err != nil {
		return nil, err
	}

	ret := ast.Decimal{Value: val}

	return ret, nil
}

func (p *parser) callonNumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumber1()
}

func (c *current) onIndex1() (interface{}, error) {
	// JSON numbers have the same syntax as Go's, and are parseable using
	return strconv.ParseInt(string(c.text), 10, 64)
}

func (p *parser) callonIndex1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndex1()
}

func (c *current) onInteger3() (interface{}, error) {
	// JSON numbers have the same syntax as Go's, and are parseable using
	val, err := strconv.ParseInt(string(c.text), 10, 64)
	if err != nil {
		return nil, err
	}

	ret := ast.Integer{Value: int(val)}

	return ret, nil
}

func (p *parser) callonInteger3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInteger3()
}

func (c *current) onString1() (interface{}, error) {
	// TODO : the forward slash (solidus) is not a valid escape in Go, it will
	// fail if there's one in the string
	text, err := strconv.Unquote(string(c.text))
	if err != nil {
		return ast.Token{}, err
	}

	ret := ast.Token{
		Value: text,
	}
	return ret, nil
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onBoolean2() (interface{}, error) {
	return ast.Bool{Value: true}, nil
}

func (p *parser) callonBoolean2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolean2()
}

func (c *current) onBoolean4() (interface{}, error) {
	return ast.Bool{Value: false}, nil
}

func (p *parser) callonBoolean4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolean4()
}

func (c *current) onNull1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNull1()
}

func (c *current) onTypeRef1(val interface{}) (interface{}, error) {
	text := string(c.text)
	ret := ast.Token{
		Value: text,
	}
	return ret, nil
}

func (p *parser) callonTypeRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeRef1(stack["val"])
}

func (c *current) onTYPE1(val interface{}) (interface{}, error) {
	text := string(c.text)
	ret := ast.Token{
		Value: text,
	}
	return ret, nil
}

func (p *parser) callonTYPE1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTYPE1(stack["val"])
}

func (c *current) onID1(val interface{}) (interface{}, error) {
	text := string(c.text)
	ret := ast.Token{
		Value: text,
	}
	return ret, nil
}

func (p *parser) callonID1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onID1(stack["val"])
}

func (c *current) onDSL1(val interface{}) (interface{}, error) {
	text := string(c.text)
	ret := ast.Token{
		Value: text,
	}
	return ret, nil
}

func (p *parser) callonDSL1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDSL1(stack["val"])
}

func (c *current) onToken1(val interface{}) (interface{}, error) {
	text := string(c.text)
	ret := ast.Token{
		Value: text,
	}
	return ret, nil
}

func (p *parser) callonToken1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onToken1(stack["val"])
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// the globalStore allows the parser to store arbitrary values
	globalStore map[string]interface{}
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			globalStore: make(map[string]interface{}),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
