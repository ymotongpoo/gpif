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
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ParsePackage start parsing all .go files under pkgroot and returns AST tree in map.
func ParsePackage(pkgroot string) (map[string][]string, error) {
	buf := make(map[string][]string)
	err := filepath.Walk(pkgroot, DirParser(buf))
	return buf, err
}

// DirParser runs parser.ParseDir
func DirParser(buf map[string][]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		// Skip dot directories, especially ".git"
		if hasDotDir(path) {
			return nil
		}

		fset := token.NewFileSet()
		pkgs, first := parser.ParseDir(fset, path, nil, parser.ImportsOnly)
		if first != nil {
			return first
		}
		paths := traverseAST(pkgs)
		buf[path] = paths
		return nil
	}
}

func hasDotDir(path string) bool {
	list := strings.Split(path, string(filepath.Separator))
	for _, s := range list {
		if strings.HasPrefix(s, ".") {
			return true
		}
	}
	return false
}

const (
	bufSize      = 100
	totalBufSize = 10000
)

type packageVisitor struct {
	imports []string
}

func newPackageVisitor() *packageVisitor {
	imports := make([]string, 0, bufSize)
	return &packageVisitor{
		imports: imports,
	}
}

func (p *packageVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch n := node.(type) {
		case *ast.ImportSpec:
			is := (*ast.ImportSpec)(n)
			if is.Path != nil {
				p.imports = append(p.imports, (*is.Path).Value)
			}
			return p
		default:
			return p
		}
	}
	return p
}

func traverseAST(root map[string]*ast.Package) []string {
	buf := make(map[string]bool)
	for _, v := range root {
		if v != nil {
			pv := newPackageVisitor()
			ast.Walk(pv, v)

			for _, i := range pv.imports {
				if _, ok := buf[i]; !ok {
					buf[i] = true
				}
			}
		}
	}
	result := make([]string, 0, len(buf))
	for k, _ := range buf {
		result = append(result, k)
	}
	return result
}
