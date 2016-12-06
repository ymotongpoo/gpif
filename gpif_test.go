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
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestParsePackage(t *testing.T) {
	in := []string{
		"./test",
	}

	for _, p := range in {
		pkgs, err := ParsePackage(p)
		if err != nil {
			t.Errorf("%s", err)
		}
		if len(pkgs) == 0 {
			t.Errorf("no direcotry")
		}

		for k, v := range pkgs {
			t.Logf("%v: %v", k, v) // Verbose testing
		}
	}
}

func TestTraverseAST(t *testing.T) {
	in := "./test"
	fset := token.NewFileSet()
	pkgs, first := parser.ParseDir(fset, in, nil, parser.ImportsOnly)
	if first != nil {
		t.Errorf("%s", first)
	}
	traverseAST(pkgs)
}

func TestDirParser(t *testing.T) {
	in := "./test"
	buf := make(map[string][]string)
	dirParser := DirParser(buf)
	fp, err := os.Open(in)
	if err != nil {
		t.Errorf("%s", err)
	}
	dir, err := fp.Stat()
	if err != nil {
		t.Errorf("%s", err)
	}
	dirParser(in, dir, nil)
}
