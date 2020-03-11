package ast

import (
	"fmt"
	// "github.com/kr/pretty"
)

func (N HofFile) Print(indent string) {
	fmt.Printf("%sFile %s/%s |%d| [\n", indent, N.Path, N.Name, len(N.Definitions) )
	for _, c := range N.Definitions {
		c.Print(indent + "  ")
	}
	fmt.Printf("%s]\n", indent)
}

func (N Definition) Print(indent string) {
	fmt.Printf("%sdef %s %s |%d| [\n", indent, N.Name, N.DSL, len(N.Body))
	for _, c := range N.Body {
		c.Print(indent + "  ")
	}
	fmt.Printf("%s]\n", indent)
}

func (N TypeDecl) Print(indent string) {
	fmt.Printf("%sdef %s %s", indent, N.Name, N.Type)
	if N.Extra != nil {
		N.Extra.Print(indent + "  ")
	} else {
		fmt.Println()
	}
}

func (N Object) Print(indent string) {
	fmt.Printf("%sObject |%d| {\n", indent, len(N.Fields))
	for _, c := range N.Fields {
		c.Print(indent + "  ")
	}
	fmt.Printf("%s}\n", indent)
}

func (N Field) Print(indent string) {
	fmt.Printf("%sField %q {\n", indent, N.Key)
	N.Value.Print(indent + "  ")
	fmt.Printf("%s}\n", indent)
}

func (N Array) Print(indent string) {
	fmt.Printf("%sArray |%d| [\n", indent, len(N.Elems))
	for _, c := range N.Elems {
		c.Print(indent + "  ")
	}
	fmt.Printf("%s]\n", indent)
}

func (N PathExpr) Print(indent string) {
	fmt.Printf("%sPathExpr |%d| {\n", indent, len(N.PathList))
	for _, c := range N.PathList {
		c.Print(indent + "  ")
	}
	fmt.Printf("%s}\n", indent)
}

func (N TokenPath) Print(indent string) {
	fmt.Printf("%sTokenPath %q\n", indent, N.Value)
}

func (N RangeExpr) Print(indent string) {
	if !N.Range {
		fmt.Printf("%sRangeExpr [%d]\n", indent, N.Low)
		return
	}

	low, high := "", ""
	if N.Low >= 0 {
		low = fmt.Sprintf("%d", N.Low)
	}
	if N.High >= 0 {
		high = fmt.Sprintf("%d", N.High)
	}

	fmt.Printf("%sRangeExpr [%s:%s]\n", indent, low, high)
}

func (N BracePath) Print(indent string) {
	fmt.Printf("%sBracePath |%d| {\n", indent, len(N.Exprs))
	for _, c := range N.Exprs {
		c.Print(indent + "  ")
	}
	fmt.Printf("%s}\n", indent)
}

func (N Token) Print(indent string) {
	fmt.Printf("%sToken %q\n", indent, N.Value)
}

func (N Integer) Print(indent string) {
	fmt.Printf("%sInteger %v\n", indent, N.Value)
}

func (N Decimal) Print(indent string) {
	fmt.Printf("%sDecimal %v\n", indent, N.Value)
}

func (N Bool) Print(indent string) {
	fmt.Printf("%sBool %v\n", indent, N.Value)
}
