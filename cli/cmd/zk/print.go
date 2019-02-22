// Copyright © 2018 NAME HERE <jbonds@jbvm.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zk

import (
	"fmt"

	"github.com/jbvmio/kafkactl/cli/zookeeper"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func printZK(i interface{}) {
	switch i := i.(type) {
	case []zookeeper.ZKPathValue:
		fmt.Println()
		for _, zk := range i {
			fmt.Println(color.YellowString("TYPE [%v]: %v", zk.Type, zk.Key))
			for _, v := range zk.Value {
				fmt.Printf("%v\n", v)
			}
		}
		fmt.Println()
	case []zookeeper.ZKPath:
		printOut(i)
	}
}

func printOut(i interface{}) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []zookeeper.ZKPath:
		tbl = table.New("TYPE", "EMPTY", "PATH")
		for _, v := range i {
			tbl.AddRow(v.Type, v.EmptyValue, v.Key)
		}
	}
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.Print()
	fmt.Println()
}
