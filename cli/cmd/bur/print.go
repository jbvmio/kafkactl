package bur

import (
	"fmt"

	"github.com/fatih/color"
	jbbur "github.com/jbvmio/burrow"
	"github.com/rodaine/table"
)

func printBur(i interface{}, kind ...string) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []jbbur.Partition:
		switch {
		case len(kind) > 0:
			outKind := kind[0]
			switch {
			case outKind == "id":
				tbl = table.New("MEMBERID", "TOPIC", "PART", "LAG", "STATUS", "CODE")
				for _, v := range i {
					tbl.AddRow(v.ClientID, v.Topic, v.Partition, v.CurrentLag, v.PStatus, v.PStatusCode)
				}
			}
		default:
			tbl = table.New("GROUP", "TOPIC", "PART", "LAG", "STATUS", "CODE")
			for _, v := range i {
				tbl.AddRow(v.Group, v.Topic, v.Partition, v.CurrentLag, v.PStatus, v.PStatusCode)
			}
		}
	}
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.Print()
	fmt.Println()
}
