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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/ymotongpoo/gpif"
)

const (
	Interval = 30 * time.Second
)

// GopathSrc is tentative GOPATH for `go get`
var GopathSrc string

func init() {
	os.Setenv("GOPATH", "/tmp")
	GopathSrc = filepath.Join(os.Getenv("GOPATH"), "src")
}

func main() {
	data, err := ioutil.ReadFile("./repos.json")
	if err != nil {
		panic(err)
	}
	list := make(map[string][]string)
	err = json.Unmarshal(data, &list)
	if err != nil {
		panic(err)
	}
	repos := list["repos"]
	for _, r := range repos {
		fmt.Printf("[Processing]: %s\n", r)
		err := GoGet(r)
		if err != nil {
			fmt.Printf("[Error] %s: %v\n", r, err)
			time.Sleep(Interval)
			continue
		}
		fmt.Printf("[Done]: %s\n", r)

		pkgpath := filepath.Join(GopathSrc, r)
		pkgs, err := gpif.ParsePackage(pkgpath)
		if err != nil {
			fmt.Printf("[Error] %s: %v\n", r, err)
		}

		for k, v := range pkgs {
			if len(v) == 0 {
				continue
			}
			fmt.Println(k, ":")
			sort.Strings(v)
			for _, pkg := range v {
				fmt.Println("\t", pkg)
			}
		}
		time.Sleep(Interval)
	}
}

// GoGet fetch repo using `go get`.
func GoGet(pkg string) error {
	cmd := exec.Command("go", "get", "-u", pkg)
	err := cmd.Run()
	return err
}
