//    Copyright 2016 Yoshi Yamaguchi
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package gpif

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// ParsePackage start parsing all .go files under pkgroot and returns AST tree in map.
func ParsePackage(pkgroot string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, first := parser.ParseDir(fset, pkgroot, nil, parser.ImportsOnly)
	if first != nil {
		return nil, first
	}
	return pkgs, nil
}

type packageVisitor func(node ast.Node) ast.Visitor

func (p packageVisitor) Visit(node ast.Node) ast.Visitor {
	return p(node)
}

func ShowImport(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BasicLit:
		bl := (*ast.BasicLit)(n)
		fmt.Printf("%v\n", bl.Value)
		return packageVisitor(ShowImport)
	default:
		fmt.Printf("%#v\n", n)
		return packageVisitor(ShowImport)
	}
	return nil
}

func traverseAST(root map[string]*ast.Package) {
	for k, v := range root {
		fmt.Println(k)
		if v != nil {
			ast.Walk(packageVisitor(ShowImport), v)
		}
	}
}
