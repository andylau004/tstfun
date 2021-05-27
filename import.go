package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strconv"
)

func int2string(num int) string {
	return strconv.Itoa(num)
}

func ExtractPkg() {
	m, _ := ParseFile("./somefun.go")
	fmt.Printf("%#v", m)
}

func ParseFile(fileName string) (map[string]string, error) {
	fset := token.NewFileSet()
	// 解析文件，主要解析token
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("parse file faield! err=", err)
		return nil, err
	}

	// 类型检查, 得到常量的值
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("cmd/hello", fset, []*ast.File{f}, nil)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, v := range f.Scope.Objects {
		if v.Kind == ast.Con {
			d := v.Decl.(*ast.ValueSpec)
			m[pkg.Scope().Lookup(v.Name).(*types.Const).Val().String()] = d.Comment.Text()
		}
	}
	return m, nil
}
