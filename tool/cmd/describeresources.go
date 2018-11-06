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
	"strings"

	"github.com/jbvmio/kafkactl"
)

func searchGroupMeta(group ...string) []kafkactl.GroupMeta {
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
	grpMeta, err := client.GetGroupMeta()
	if err != nil {
		log.Fatalf("Error getting group metadata: %s\n", err)
	}
	if len(group) >= 1 {
		if group[0] != "" {
			var groupMeta []kafkactl.GroupMeta
			for _, g := range group {
				for _, gm := range grpMeta {
					if exact {
						if gm.Group == g {
							groupMeta = append(groupMeta, gm)
						}
					} else {
						if strings.Contains(gm.Group, g) {
							groupMeta = append(groupMeta, gm)
						}
					}
				}
			}
			return groupMeta
		}
	}
	return grpMeta
}

func groupMetaByTopic(topic string, grpMeta []kafkactl.GroupMeta) []kafkactl.GroupMeta {
	var groupMeta []kafkactl.GroupMeta
	for _, gm := range grpMeta {
		gMeta := gm
		gMeta.MemberAssignments = nil
		for _, ma := range gm.MemberAssignments {
			mm := kafkactl.MemberMeta{
				ClientHost:      ma.ClientHost,
				ClientID:        ma.ClientID,
				Active:          ma.Active,
				TopicPartitions: make(map[string][]int32),
			}
			for t, p := range ma.TopicPartitions {
				if exact {
					if t == topic {
						mm.TopicPartitions[t] = append(mm.TopicPartitions[t], p...)
						gMeta.MemberAssignments = append(gMeta.MemberAssignments, mm)
					}
				} else {
					if strings.Contains(t, topic) {
						mm.TopicPartitions[t] = append(mm.TopicPartitions[t], p...)
						gMeta.MemberAssignments = append(gMeta.MemberAssignments, mm)
					}
				}

			}
		}
		groupMeta = append(groupMeta, gMeta)
	}
	return groupMeta
}

func getGroupLag(grpMeta []kafkactl.GroupMeta) []PartitionLag {
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

	var partitionLag []PartitionLag
	for _, gm := range grpMeta {
		for _, m := range gm.MemberAssignments {
			for topic, partitions := range m.TopicPartitions {
				for _, p := range partitions {
					offset, lag, err := client.OffSetAdmin().Group(gm.Group).Topic(topic).GetOffsetLag(p)
					if err != nil {
						lag = -7777
					}
					pl := PartitionLag{
						Group:     gm.Group,
						Topic:     topic,
						Partition: p,
						Member:    m.ClientID,
						Offset:    offset,
						Lag:       lag,
					}
					partitionLag = append(partitionLag, pl)
				}
			}
		}
	}
	return partitionLag
}

/*
func getGroupLag(grpMeta []kafkactl.GroupMeta) []GroupLag {
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
	var groupLag []GroupLag
	glChan := make(chan GroupLag, 100)
	doneChan := make(chan bool, 100)
	for _, gm := range grpMeta {
		go goGetGroupLag(client, gm, glChan, doneChan)
	}
	i := 0
	for i < len(grpMeta) {
		select {
		case gl := <-glChan:
			groupLag = append(groupLag, gl)
		case <-doneChan:
			i++
		}
	}
	return groupLag
}

func goGetGroupLag(client *kafkactl.KClient, gm kafkactl.GroupMeta, glChan chan GroupLag, doneChan chan bool) {
	memDone := make(chan bool, 100)
	grpName := gm.Group
	for _, m := range gm.MemberAssignments {
		cID := m.ClientID
		for t, p := range m.TopicPartitions {
			go goGetMemberLag(client, grpName, t, cID, p, glChan, memDone)
		}
		for i := 0; i < len(m.TopicPartitions); i++ {
			select {
			case <-memDone:
				fmt.Println("memDONE")
			}
		}

	}
	doneChan <- true
}

func goGetMemberLag(client *kafkactl.KClient, grpName, topName, memID string, p []int32, glChan chan GroupLag, memDone chan bool) {
	totalLag, err := client.OffSetAdmin().Group(grpName).Topic(topName).GetTotalLag(p)
	if err != nil {
		totalLag = -77777
	}
	gl := GroupLag{
		Group:      grpName,
		Topic:      topName,
		Partitions: p,
		Member:     memID,
		Lag:        totalLag,
	}
	glChan <- gl
	memDone <- true
}
*/

func testOSA() {
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

	off, lag, err := client.OffSetAdmin().Group("v2-pps-storage-consumer-auditor-cpx-auditor").Topic("action.local.pif.invoices").GetOffsetLag(5)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Println(off, lag)
}
