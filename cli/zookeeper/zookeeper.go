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

package zookeeper

import (
	"fmt"
	"log"
	"sort"

	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/zk"
)

var (
	// verbose enables additional details to print.
	verbose bool
)

type ZKFlags struct {
	Context string
	Depth   uint8
	Recurse bool
	Values  bool
	Verbose bool
}

type ZKPathValue struct {
	Type       string
	Key        string
	Value      []string
	EmptyValue bool
}

type ZKPath struct {
	Type       string
	Key        string
	EmptyValue bool
}

// client variables
var (
	zkClient *zk.ZooKeeper
)

func LaunchZKClient(context *cx.Context, flags ZKFlags) {
	if len(context.Zookeeper) < 1 {
		log.Fatalf("No Zookeeper Servers Defined.\n")
	}
	zkClient = zk.NewZooKeeper()
	zkClient.EnableLogger(flags.Verbose)
	zkClient.SetServers(context.Zookeeper)
}

func ZKls(path ...string) []ZKPathValue {
	var ZKP []ZKPathValue
	for _, p := range path {
		zkp := ZKPathValue{
			Key: p,
		}
		match := true
		sp, err := zkClient.Children(p)

		switch match {
		case err != nil:
			log.Fatalf("Error retrieving path: %v\n", p)
		case len(sp) == 0:
			val, err := zkClient.Get(p)
			if err != nil {
				log.Fatalf("Error retrieving value: %v\n", p)
			}
			value := fmt.Sprintf("%s", val)
			zkp.Type = "value"
			zkp.EmptyValue = (value == "")
			if !zkp.EmptyValue {
				zkp.Value = []string{value}
			}
			//zkp.Value = []string{fmt.Sprintf("%s", val)}
			//zkp.EmptyValue = (zkp.Value == "")
		default:
			zkp.Type = "path"
			zkp.Value = sp
		}
		ZKP = append(ZKP, zkp)
	}
	return ZKP
}

func ZKRecurseLS(depth uint8, path ...string) []ZKPath {
	var zkPaths []ZKPath
	var limit uint8
	pathMap := make(map[string]bool)
	ZKP := ZKls(path...)
	zkPaths = append(zkPaths, zkConvertValPath(ZKP...)...)
	count := len(ZKP)
RecurseLoop:
	for count > 0 {
		var newPaths []string
		for _, parent := range ZKP {
			if parent.Type == "path" {
				if !pathMap[parent.Key] {
					pathMap[parent.Key] = true
				}
				for _, child := range parent.Value {
					if !pathMap[child] {
						var new string
						if parent.Key == "/" {
							new = string(parent.Key + child)
						} else {
							new = string(parent.Key + "/" + child)
						}
						newPaths = append(newPaths, new)
						pathMap[new] = true
					}
				}
			}
		}
		ZKP = ZKls(newPaths...)
		zkPaths = append(zkPaths, zkConvertValPath(ZKP...)...)
		count = len(ZKP)
		limit++
		if limit == depth {
			break RecurseLoop
		}
	}
	sort.Slice(zkPaths, func(i, j int) bool {
		return zkPaths[i].Key < zkPaths[j].Key
	})
	return zkPaths
}

func zkConvertValPath(ZKP ...ZKPathValue) []ZKPath {
	var zkPaths []ZKPath
	for _, zkpv := range ZKP {
		zkPaths = append(zkPaths, ZKPath{
			Type:       zkpv.Type,
			Key:        zkpv.Key,
			EmptyValue: zkpv.EmptyValue,
		})
	}
	return zkPaths
}

func ZKFilterAllVals(zkp []ZKPath) []ZKPath {
	var zkPaths []ZKPath
	for _, z := range zkp {
		if !z.EmptyValue && (z.Type == "value") {
			zkPaths = append(zkPaths, z)
		}
	}
	return zkPaths
}

/*
func ZKls(path ...string) {
	if len(path) == 0 {
		path = []string{"/"}
	}
	subPath := make(map[string][]string, len(path))
	for _, p := range path {
		var keyVal = string("PATH: " + p)
		sp, err := zkClient.Children(p)
		if err != nil {
			log.Fatalf("Error retrieving path: %v\n", p)
		}
		if len(sp) == 0 {
			val, err := zkClient.Get(p)
			if err == nil {
				keyVal = string("VALUE: " + p)
				sp = []string{
					fmt.Sprintf("%s", val),
				}
			}
		}
		subPath[keyVal] = sp
	}
	fmt.Println()
	for k, v := range subPath {
		fmt.Println(color.YellowString(k))
		for _, sp := range v {
			fmt.Printf(" %v\n", sp)
		}
	}
	fmt.Println()
}

func zkCreateValue(path string, value []byte) {
	if path == "" {
		log.Fatalf("Error: No Path Specified.\n")
	}
	check, err := zkClient.Exists(path)
	if err != nil {
		log.Fatalf("Error Validating Path: %v\n", err)
	}
	if check {
		if value == nil {
			if !zkForceUpdate {
				log.Fatalf("Empty Value Submitted. Use --force to override.\n")
			}
		}
		_, err := zkClient.Set(path, value)
		if err != nil {
			log.Fatalf("Error Setting Value: %v\n", err)
		}
		fmt.Println("Successfully Updated:", path)
		return
	}
	str, err := zkClient.Create(path, value, "", false)
	if err != nil {
		log.Fatalf("Error Setting Value: %v\n", err)
	}
	fmt.Println("Successfully Created:", str)
	return
}

func zkDeleteValue(path string) {
	if path == "" {
		log.Fatalf("Error: No Path Specified.\n")
	}
	err := zkClient.Delete(path)
	if err != nil {
		log.Fatalf("Error Deleting Path: %v\n", err)
	}
	fmt.Println("Successfully Deleted:", path)
}

func zkCreatePRE(path string, value []byte) {
	if path == "" {
		log.Fatalf("Error: No Path Specified.\n")
	}
	check, err := zkClient.Exists(path)
	if err != nil {
		log.Fatalf("Error Validating Path: %v\n", err)
	}
	if check {
		log.Fatalf("Error: Preferred Replica Election Already in Progress.\n")
	}
	_, err = zkClient.Create(path, value, "", false)
	if err != nil {
		log.Fatalf("Error Creating Preferred Replica Election: %v\n", err)
	}
	fmt.Println("\n\n", "Successfully Created Preferred Replica Election.\n")
	return
}

func zkCreateReassignPartitions(path string, value []byte) {
	if path == "" {
		log.Fatalf("Error: No Path Specified.\n")
	}
	check, err := zkClient.Exists(path)
	if err != nil {
		log.Fatalf("Error Validating Path: %v\n", err)
	}
	if check {
		log.Fatalf("Error: Reassign Partitions Already in Progress.\n")
	}
	_, err = zkClient.Create(path, value, "", false)
	if err != nil {
		log.Fatalf("Error Creating Reassign Partitions: %v\n", err)
	}
	fmt.Println("\n\n", "Successfully Started Reassign Partition Process.\n")
	return
}
*/
