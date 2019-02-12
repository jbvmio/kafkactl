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
	"sort"

	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/jbvmio/zk"
)

var (
	// verbose enables additional details to print.
	verbose bool
)

type ZKFlags struct {
	Context string
	Value   string
	Depth   uint8
	Force   bool
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
		//kafka.CloseClient()
		out.Failf("No Zookeeper Servers Defined.\n")
	}
	zkClient = zk.NewZooKeeper()
	zkClient.EnableLogger(flags.Verbose)
	zkClient.SetServers(context.Zookeeper)
	ok, err := zkClient.Exists("/")
	if !ok || err != nil {
		out.Failf("Error Validating Zookeeper Configuration.", context.Zookeeper)
	}
}

func KafkaZK(ctx *cx.Context, verbose bool) error {
	//ctx := cfg.GetContext(context)
	if len(ctx.Zookeeper) < 1 {
		out.Failf("No Zookeeper Servers Defined.\n")
	}
	zkClient = zk.NewZooKeeper()
	zkClient.EnableLogger(verbose)
	zkClient.SetServers(ctx.Zookeeper)
	ok, err := zkClient.Exists("/admin")
	if !ok || err != nil {
		return fmt.Errorf("Error Validating Zookeeper Configuration.", ctx.Zookeeper)
	}
	return nil
}

func ZKCheckExists(path string) (bool, error) {
	return zkClient.Exists(path)
}

func ZKls(path ...string) []ZKPathValue {
	var ZKP []ZKPathValue
	for _, p := range path {
		zkp := ZKPathValue{
			Key: p,
		}
		sp, err := zkClient.Children(p)
		switch true {
		case err != nil:
			//kafka.CloseClient()
			out.Failf("Error retrieving path: %v", p)
		case len(sp) == 0:
			val, err := zkClient.Get(p)
			if err != nil {
				//kafka.CloseClient()
				out.Failf("Error retrieving value: %v", p)
			}
			value := fmt.Sprintf("%s", val)
			zkp.Type = "value"
			zkp.EmptyValue = (value == "")
			if !zkp.EmptyValue {
				zkp.Value = []string{value}
			}
		default:
			zkp.Type = "path"
			zkp.Value = sp
		}
		ZKP = append(ZKP, zkp)
	}
	return ZKP
}

func ZKGetValue(path string) ([]byte, error) {
	return zkClient.Get(path)
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

func ZKCreate(path string, silent, force bool, value ...byte) {
	switch true {
	case path != "" && value == nil:
		zkCreatePath(path, silent, force)
	case path != "" && value != nil:
		zkCreateValue(path, silent, force, value)
	default:
		//kafka.CloseClient()
		out.Failf("Invalid Path or Value")
	}
}

func ZKSetSeq() {
	zkClient.SetSequencial()
}

func ZKRemoveFlags() {
	zkClient.SetFlags(0)
}

func zkCreatePath(path string, silent, force bool) {
	var pathString string
	var errd error
	check, err := zkClient.Exists(path)
	switch true {
	case err != nil:
		//kafka.CloseClient()
		out.Failf("Error Validating Path: %v", err)
	case check:
		//kafka.CloseClient()
		out.Failf("Error: Path Exists.")
	case force:
		pathString, errd = zkClient.Create(path, nil, "", true)
	default:
		pathString, errd = zkClient.Create(path, nil, "", false)
	}
	if errd != nil {
		//kafka.CloseClient()
		out.Failf("Error Creating Path: %v. Try --force", errd)
	}
	if !silent {
		out.Infof("ZK Successfully Created: %v", pathString)
	}
}

func zkCreateValue(path string, silent, force bool, value []byte) {
	var pathString string
	var errd error
	var children bool
	check, err := zkClient.Exists(path)
	children, errd = zkClient.HasChildren(path)
	switch true {
	case err != nil || (errd != nil && children):
		//kafka.CloseClient()
		fmt.Println(check, err)
		fmt.Println(children, errd)
		out.Failf("Error Validating Path:\n %v\n %v", err, errd)
	case children:
		//kafka.CloseClient()
		out.Failf("Error: Child SubPaths Detected.")
	case check && force:
		_, errd = zkClient.Set(path, value)
		pathString = path
	case check:
		//kafka.CloseClient()
		out.Failf("Error: Path Exists. Use --force to override.")
	case force:
		pathString, errd = zkClient.Create(path, value, "", true)
	default:
		pathString, errd = zkClient.Create(path, value, "", false)
	}
	if errd != nil {
		//kafka.CloseClient()
		out.Failf("Error Creating Path: %v", errd)
	}
	if !silent {
		out.Infof("ZK Successfully Created: %v", pathString)
	}
}

func ZKDelete(path string, RMR bool) {
	var errd error
	check, err := zkClient.Exists(path)
	switch true {
	case path == "":
		//kafka.CloseClient()
		out.Failf("Empty Path Entered.")
	case !check:
		//kafka.CloseClient()
		out.Failf("Path Does Not Exist.")
	case err != nil:
		//kafka.CloseClient()
		out.Failf("Error Validating Path: %v", err)
	case RMR:
		errd = zkClient.DeleteRecursive(path)
	default:
		errd = zkClient.Delete(path)
	}
	if errd != nil {
		//kafka.CloseClient()
		out.Failf("Error Deleting Path: %v\n", errd)
	}
	out.Infof("Successfully Deleted: %v", path)
}

/*
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
