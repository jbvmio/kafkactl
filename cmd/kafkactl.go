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
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func formatOutput(k interface{}) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColumnSeparator("")

	switch k := k.(type) {
	case []groupListMeta:
		table.SetHeader([]string{"GroupType", "GroupName", "Coordinator"})
		for _, v := range k {
			x := []string{v.Type, v.Group, v.Coordinator}
			table.Append(x)
		}
	case []groupMetadata:
		if searchByGroup && searchByTopic {
			table.SetHeader([]string{"Group", "Topic", "Parts", "Lag", "ClientID", "Host"})
			for _, v := range k {
				k := []string{v.Group, v.Topic, v.Partitions, strconv.FormatInt(v.Lag, 10), v.ClientID, v.Host}
				table.Append(k)
			}
		} else {
			table.SetHeader([]string{"Group", "Topic", "Parts", "ClientID", "Host"})
			for _, v := range k {
				k := []string{v.Group, v.Topic, v.Partitions, v.ClientID, v.Host}
				table.Append(k)
			}
		}
	case []datameta:
		table.SetHeader([]string{"Topic", "Partition", "Replicas", "ISRs", "Partition Leader"})
		for _, v := range k {
			k := []string{v.Topic, strconv.Itoa(v.Partition), v.Replicas, v.ISRs, v.Leader}
			table.Append(k)
		}
	case []topicSummary:
		table.SetHeader([]string{"Topic", "Partitions", "RFactor", "ISRs"})
		for _, v := range k {
			k := []string{v.Topic, strconv.Itoa(v.Partitions), strconv.Itoa(v.RFactor), strconv.Itoa(v.ISR)}
			table.Append(k)
		}
	}

	fmt.Println()
	table.Render()
	fmt.Println()

}

func parseHost(s string) string {
	return strings.Split(s, ".")[0]
}
