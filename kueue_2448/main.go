package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const STRUCT_NAME = "resource.FlavorResource"
const KEY_1 = "Flavor"
const KEY_2 = "Resource"
const SOURCE_TYPE = "FlavorResourceQuantities"
const DEST_TYPE = "resource.FlavorResourceQuantities"

func flattenMapEntry(key1 ast.Expr, key2 ast.Expr, value ast.Expr) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key: &ast.CompositeLit{Type: &ast.Ident{Name: STRUCT_NAME}, Elts: []ast.Expr{
			&ast.KeyValueExpr{Key: &ast.Ident{Name: KEY_1}, Value: key1},
			&ast.KeyValueExpr{Key: &ast.Ident{Name: KEY_2}, Value: key2},
		}},
		Value: value}
}

func flattenMap(a *ast.CompositeLit) {
	elements := make([]ast.Expr, 0)
	for _, elem := range a.Elts {
		outerKV := elem.(*ast.KeyValueExpr)
		for _, innerElem := range outerKV.Value.(*ast.CompositeLit).Elts {
			innerKV := innerElem.(*ast.KeyValueExpr)
			elements = append(elements, flattenMapEntry(outerKV.Key, innerKV.Key, innerKV.Value))
		}
	}
	a.Elts = elements
}

type mapFlattener struct {
}

func (v mapFlattener) Visit(node ast.Node) ast.Visitor {
	if t1, isType := node.(*ast.CompositeLit); isType {
		if strings.Contains(fmt.Sprint(t1.Type), SOURCE_TYPE) {
			flattenMap(t1)
			t1.Type = &ast.Ident{Name: DEST_TYPE}
		}
	}
	return &v
}

type nestedMap map[string]map[string]int64

func main() {
	// if you want to run this program on its source ;)
	// Note: will have to change SOURCE_TYPE to nestedMap
	_ = nestedMap{"a1": {"a2": 5}, "b1": {"b2": 6, "b3": 3}, "a2":{}}

	fset := token.FileSet{}
	f, err := parser.ParseFile(&fset, os.Args[2], nil, parser.ParseComments)
	if err != nil {
		panic("failure parsing")
	}
	for _, decl := range f.Decls {
		ast.Walk(mapFlattener{}, decl)
	}

	err = format.Node(os.Stdout, &fset, f)
	if err != nil {
		panic("failure fomratting")
	}
}
