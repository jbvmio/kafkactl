package cmd

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	"github.com/jbvmio/kafkactl"
	"github.com/rodaine/table"
	"github.com/spf13/cast"
)

func printOutput(i interface{}) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []kafkactl.TopicMeta:
		tbl = table.New("TOPIC", "PART", "REPLs", "ISRs", "OFF", "LDR")
		for _, v := range i {
			tbl.AddRow(v.Topic, v.Partition, v.Replicas, v.ISRs, v.OfflineReplicas, v.Leader)
		}
	case []kafkactl.GroupListMeta:
		tbl = table.New("GROUPTYPE", "GROUP", "COORDINATOR")
		for _, v := range i {
			tbl.AddRow(v.Type, v.Group, v.Coordinator)
		}
	case []kafkactl.GroupMeta:
		tbl = table.New("GROUP", "TOPIC", "PART", "MEMBER")
		for _, v := range i {
			grpName := truncateString(v.Group, 64)
			for _, m := range v.MemberAssignments {
				cID := m.ClientID
				for t, p := range m.TopicPartitions {
					tbl.AddRow(grpName, t, makeSeqStr(p), cID)
				}
			}
		}
	case []PartitionLag:
		tbl = table.New("GROUP", "TOPIC", "PART", "MEMBER", "OFFSET", "LAG")
		for _, v := range i {
			tbl.AddRow(v.Group, v.Topic, v.Partition, v.Member, v.Offset, v.Lag)
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

func makeSeqStr(nums []int32) string {
	seqMap := make(map[int][]int32)
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	var mapCount int
	var done int
	var switchInt int
	seqMap[mapCount] = append(seqMap[mapCount], nums[done])
	done++
	switchInt = done
	for done < len(nums) {
		if nums[done] == ((seqMap[mapCount][(switchInt - 1)]) + 1) {
			seqMap[mapCount] = append(seqMap[mapCount], nums[done])
			switchInt++
		} else {
			mapCount++
			seqMap[mapCount] = append(seqMap[mapCount], nums[done])
			switchInt = 1
		}
		done++
	}
	var seqStr string
	for k, v := range seqMap {
		if k > 0 {
			seqStr += ","
		}
		if len(v) > 1 {
			seqStr += cast.ToString(v[0])
			seqStr += "-"
			seqStr += cast.ToString(v[len(v)-1])
		} else {
			seqStr += cast.ToString(v[0])
		}
	}
	return seqStr
}
