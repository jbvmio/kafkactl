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

	"github.com/jbvmio/kafkactl"
)

type GroupTopicOffsetMeta struct {
	TopicPartitionMeta
	GroupPartitionMeta
}

type GroupPartitionMeta struct {
	Group            string
	Partition        int32
	GroupOffset      int64
	Lag              int64
	GroupCoordinator string
}

type TopicPartitionMeta struct {
	Topic       string
	Partition   int32
	TopicOffset int64
	TopicLeader int32
}

func getGroupTopicOffsets(group, topic string) ([]GroupTopicOffsetMeta, error) {
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
	exact = true
	var GTMeta []GroupTopicOffsetMeta
	tSum := kafkactl.GetTopicSummary(searchTopicMeta(topic))
	gList := searchGroupListMeta(group)

	fmt.Printf("TOPIC:\n%+v\n", tSum)
	fmt.Printf("GROUP:\n%+v\n", gList)

	for _, ts := range tSum {
		for _, p := range ts.Partitions {
			for _, grp := range gList {
				tpm := TopicPartitionMeta{
					Topic:       ts.Topic,
					Partition:   p,
					TopicLeader: ts.Leader,
				}
				gpm := GroupPartitionMeta{
					Group:            grp.Group,
					GroupCoordinator: grp.Coordinator,
				}

				fmt.Printf("tpmNIL: %+v gpmNIL: %+v\n", (tpm == nil), (gpm == nil))

				fmt.Printf("tpm:\n%+v\n", tpm)
				fmt.Printf("gpm:\n%+v\n", gpm)

				offset, lag, err := client.OffSetAdmin().Group(grp.Group).Topic(ts.Topic).GetOffsetLag(p)
				if err != nil {
					log.Fatalf("Error: %v\n", err)
				}

				fmt.Printf("Offset: %+v Lag: %+v Error: %v\n", offset, lag, err)

				tpm.TopicOffset = (offset + lag)
				gpm.GroupOffset = offset
				gpm.Lag = lag
				gpm.Partition = p

				gtMeta := GroupTopicOffsetMeta{
					tpm,
					gpm,
				}

				GTMeta = append(GTMeta, gtMeta)
				fmt.Printf("GTMETA:\n%+v\n", gtMeta)
			}
		}
	}

	return GTMeta, err
}
