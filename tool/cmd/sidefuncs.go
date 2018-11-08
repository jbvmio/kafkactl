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

package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/jbvmio/kafkactl"
	"github.com/rodaine/table"
)

func printOutput(i interface{}) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []kafkactl.TopicSummary:
		tbl = table.New("TOPIC", "PART", "RFactor", "ISRs", "OFFLINE", "LEADER")
		for _, v := range i {
			tbl.AddRow(v.Topic, v.Parts, v.RFactor, v.ISRs, v.OfflineReplicas, v.Leader)
		}
	case []kafkactl.TopicMeta:
		tbl = table.New("TOPIC", "PART", "REPLICAS", "ISRs", "OFFLINE", "LEADER")
		for _, v := range i {
			tbl.AddRow(v.Topic, v.Partition, v.Replicas, v.ISRs, v.OfflineReplicas, v.Leader)
		}
	case []kafkactl.GroupListMeta:
		tbl = table.New("GROUPTYPE", "GROUP", "COORDINATOR")
		for _, v := range i {
			tbl.AddRow(v.Type, v.Group, v.Coordinator)
		}
	case []GroupTopicOffsetMeta:
		tbl = table.New("GROUP", "PARTITION", "PART", "GrpOFFSET", "TopicOFFSET", "LAG", "GrpCoordinator")
		for _, v := range i {
			tbl.AddRow(v.Group, v.Topic, v.Partition, v.GroupOffset, v.TopicOffset, v.Lag, v.GroupCoordinator)
		}
	case []kafkactl.GroupMeta:
		tbl = table.New("GROUP", "TOPIC", "PART", "MEMBER")
		for _, v := range i {
			grpName := truncateString(v.Group, 64)
			for _, m := range v.MemberAssignments {
				cID := m.ClientID
				for t, p := range m.TopicPartitions {
					tbl.AddRow(grpName, t, kafkactl.MakeSeqStr(p), cID)
				}
			}
		}
	case []PartitionLag:
		tbl = table.New("GROUP", "TOPIC", "PART", "MEMBER", "OFFSET", "LAG")
		for _, v := range i {
			tbl.AddRow(v.Group, v.Topic, v.Partition, v.Member, v.Offset, v.Lag)
		}
	case []TopicConfig:
		tbl = table.New("TOPIC", "CONFIG", "VALUE", "READONLY", "DEFAULT", "SENSITIVE")
		for _, v := range i {
			tbl.AddRow(v.Topic, v.Config, v.Value, v.ReadOnly, v.Default, v.Sensitive)
		}
	}
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.Print()
	fmt.Println()
}

func truncateString(str string, num int) string {
	s := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		s = str[0:num] + "..."
	}
	return s
}

func filterUnique(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func testFunc() {
	client, err := kafkactl.NewClient(bootStrap)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("Error closing client: %v\n", err)
		}
	}()
	if verbose {
		client.Logger("")
	}
	offset, lag, err := client.OffSetAdmin().Group("jblap").Topic("testtopic").GetOffsetLag(0)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Println(offset, lag)
}
