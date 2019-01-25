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

package kafka

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/x"
	"github.com/rodaine/table"
)

func PrintOut(i interface{}) {
	var highlightColumn = true
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []*Broker:
		tbl = table.New("BROKER", "ID", "GRPs", "P.LEADERS", "P.REPLICAS", "P.TOTAL", "P.NOTLEADER")
		for _, v := range i {
			tbl.AddRow(v.Address, v.ID, v.GroupCoordinating, v.LeaderPartitions, v.ReplicaPartitions, v.TotalPartitions, v.NotLeader)
		}
	case []kafkactl.TopicSummary:
		tbl = table.New("TOPIC", "PART", "RFactor", "ISRs", "OFFLINE")
		for _, v := range i {
			tbl.AddRow(v.Topic, v.Parts, v.RFactor, v.ISRs, v.OfflineReplicas)
		}
	case []kafkactl.TopicOffsetMap:
		tbl = table.New("TOPIC", "PART", "OFFSET", "LEADER", "REPLICAS", "ISRs", "OFFLINE")
		for _, v := range i {
			for _, p := range v.TopicMeta {
				tbl.AddRow(p.Topic, p.Partition, v.PartitionOffsets[p.Partition], p.Leader, p.Replicas, p.ISRs, p.OfflineReplicas)
			}
		}
	case []kafkactl.GroupListMeta:
		tbl = table.New("GROUPTYPE", "GROUP", "COORDINATOR")
		for _, v := range i {
			tbl.AddRow(v.Type, v.Group, v.CoordinatorAddr)
		}
	case []kafkactl.GroupMeta:
		tbl = table.New("GROUP", "TOPIC", "PART", "MEMBER")
		for _, v := range i {
			grpName := x.TruncateString(v.Group, 64)
			for _, m := range v.MemberAssignments {
				cID := m.ClientID
				for t, p := range m.TopicPartitions {
					tbl.AddRow(grpName, t, x.MakeSeqStr(p), cID)
				}
			}
		}
	}
	if highlightColumn {
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	} else {
		tbl.WithHeaderFormatter(headerFmt)
	}

	tbl.Print()
	fmt.Println()
}
