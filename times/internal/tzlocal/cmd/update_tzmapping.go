// Copyright 2023 The Blocky Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blockysource/go-pkg/times/internal/tzdata"
)

func main() {
	path, _ := filepath.Abs("./tzmapping.go")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString("// Code generated by tzlocal/update_tzmapping.go DO NOT EDIT.\n")
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString("package tzlocal\n\n")
	if err != nil {
		panic(err)
	}
	if err = tzdata.UpdateWindowsTZMapping(f); err != nil {
		fmt.Println(err)
	}
}
