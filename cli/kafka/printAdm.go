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
	"github.com/rodaine/table"
)

func PrintAdm(i interface{}) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []TopicConfig:
		tbl = table.New("TOPIC", "CONFIG", "VALUE", "READONLY", "DEFAULT", "SENSITIVE")
		for _, v := range i {
			tbl.AddRow(v.Topic, v.Config, v.Value, v.ReadOnly, v.Default, v.Sensitive)
		}
	case PRETopicMeta:
		tbl = table.New("TOPIC", "PART", "LEADER", "PREFERRED")
		for _, v := range i.Partitions {
			tbl.AddRow(v.Topic, v.Partition, v.Replicas[0], v.Leader)
		}
	case RAPartList:
		tbl = table.New("TOPIC", "PARTITION", "BROKERS")
		for _, v := range i.Partitions {
			tbl.AddRow(v.Topic, v.Partition, v.Replicas)
		}
	case ReplicaDetails:
		tbl = table.New("TOPIC", "PART", "LEADER", "BROKERS", "INSYNC", "PRE")
		for _, v := range i.TopicMetadata {
			var pre string
			if len(v.Replicas) > 0 {
				if v.Leader != v.Replicas[0] {
					pre = "x"
				}
			}
			tbl.AddRow(v.Topic, v.Partition, v.Leader, v.Replicas, v.ISRs, pre)
		}
	case PRESummary:
		tbl = table.New("TOPICS.NEED.PRE", "PARTITIONS.NEED.PRE")
		var topicCount uint16
		var partCount uint32
		for _, v := range i.Topics {
			topicCount++
			partCount = partCount + i.PRECount[v]
		}
		tbl.AddRow(topicCount, partCount)
	case OffsetDetails:
		switch {
		case i.IncludesGroups:
			tbl = table.New("TOPIC", "PARTITION", "OLDEST", "NEWEST", "GROUP", "OFFSET", "LAG")
			for _, v := range i.Details {
				tbl.AddRow(v.Topic, v.Partition, v.TopicOffsetOldest, v.TopicOffsetNewest, v.Group, v.GroupOffset, v.Lag)
			}
		default:
			tbl = table.New("TOPIC", "PARTITION", "OLDEST", "NEWEST")
			for _, v := range i.Details {
				tbl.AddRow(v.Topic, v.Partition, v.TopicOffsetOldest, v.TopicOffsetNewest)
			}
		}
	}
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.Print()
	fmt.Println()
}
