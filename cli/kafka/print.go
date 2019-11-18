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
	"sort"

	"github.com/jbvmio/kafkactl/cli/x"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/fatih/color"
	kafkactl "github.com/jbvmio/kafka"
	"github.com/rodaine/table"
)

func PrintOut(i interface{}) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []*Broker:
		tbl = table.New("BROKER", "ID", "GRPs", "LDR.REPLICAS", "PEER.REPLICAS", "TOTAL.REPLICAS", "MIGRATING.REPLICAS", "OVERLOAD")
		for _, v := range i {
			tbl.AddRow(v.Address, v.ID, v.GroupsCoordinating, v.LeaderReplicas, v.PeerReplicas, v.TotalReplicas, v.MigratingReplicas, v.Overload)
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
	case []PartitionLag:
		tbl = table.New("GROUP", "TOPIC", "PART", "MEMBER", "OFFSET", "LAG", "HOST")
		for _, v := range i {
			tbl.AddRow(v.Group, v.Topic, v.Partition, v.Member, v.Offset, v.Lag, v.Host)
		}
	case []TotalLag:
		tbl = table.New("GROUP", "TOPIC", "TOTALLAG")
		for _, v := range i {
			tbl.AddRow(v.Group, v.Topic, v.TotalLag)
		}
	case []APIVersion:
		tbl = table.New("NAME", "KEY", "MAXVERSION")
		for _, v := range i {
			tbl.AddRow(v.Name, v.Key, v.MaxVersion)
		}
	case OffsetRangeMap:
		tbl = table.New("TOPIC", "PARTITION", "OFFSET.RANGE", "COUNT")
		for topic, parts := range i.Ranges {
			var P []int32
			for p := range parts {
				P = append(P, p)
			}
			sort.SliceStable(P, func(i, j int) bool {
				return P[i] < P[j]
			})
			for _, x := range P {
				offsets := fmt.Sprintf("%v-%v", parts[x][0], parts[x][1])
				tbl.AddRow(topic, x, offsets, (parts[x][1] - parts[x][0]))
			}
		}
	}
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.Print()
	fmt.Println()
}

func PrintMSGs(msgs []*kafkactl.Message, outFlags out.OutFlags) {
	switch true {
	case outFlags.Header:
		for _, msg := range msgs {
			out.Infof("%s", msg.Value)
		}
	default:
		headerFmt := color.New(color.FgGreen).SprintfFunc()
		for _, msg := range msgs {
			h := headerFmt("TOPIC: %v, PARTITION: %v, OFFSET: %v, TIMESTAMP: %v\n", msg.Topic, msg.Partition, msg.Offset, msg.Timestamp)
			out.Infof("%v%s\n", h, msg.Value)
		}
	}
}

// PrintMSG returns messages displayed by the desired format while following a topic.
func PrintMSG(msg *kafkactl.Message, outFlags out.OutFlags) {
	switch true {
	case outFlags.Format != "":
		out.Marshal(msg, outFlags.Format)
	case outFlags.Header:
		out.Infof("%s", msg.Value)
	default:
		headerFmt := color.New(color.FgGreen).SprintfFunc()
		h := headerFmt("TOPIC: %v, PARTITION: %v, OFFSET: %v, TIMESTAMP: %v\n", msg.Topic, msg.Partition, msg.Offset, msg.Timestamp)
		out.Infof("%v%s\n", h, msg.Value)
	}
}
