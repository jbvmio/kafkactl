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

	"github.com/fatih/color"
	"github.com/jbvmio/zk"
)

var zkClient *zk.ZooKeeper

func launchZKClient() {
	if len(zkServers) < 1 {
		log.Fatalf("No Zookeeper Servers Defined.\n")
	}
	zkClient = zk.NewZooKeeper()
	zkClient.EnableLogger(verbose)
	zkClient.SetServers(zkServers)
}

func zkLS(path ...string) {
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
