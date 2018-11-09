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
	"log"
	"sort"
	"strings"

	"github.com/jbvmio/kafkactl"
)

func searchTopicMeta(topics ...string) []kafkactl.TopicMeta {
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
	tMeta, err := client.GetTopicMeta()
	if err != nil {
		log.Fatalf("Error getting topic metadata: %s\n", err)
	}
	var topicMeta []kafkactl.TopicMeta
	if len(topics) >= 1 {
		if topics[0] != "" {
			for _, t := range topics {
				for _, m := range tMeta {
					if exact {
						if m.Topic == t {
							topicMeta = append(topicMeta, m)
						}
					} else {
						if strings.Contains(m.Topic, t) {
							topicMeta = append(topicMeta, m)
						}
					}
				}
			}
		} else {
			topicMeta = tMeta
		}
	}
	sort.Slice(topicMeta, func(i, j int) bool {
		if topicMeta[i].Topic < topicMeta[j].Topic {
			return true
		}
		if topicMeta[i].Topic > topicMeta[j].Topic {
			return false
		}
		return topicMeta[i].Partition < topicMeta[j].Partition
	})
	return topicMeta
}

func getTopicOffsetMap(tm []kafkactl.TopicMeta) []kafkactl.TopicOffsetMap {
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
	return client.MakeTopicOffsetMap(tm)

}

func chanGetTopicOffsetMap(t []kafkactl.TopicMeta) []kafkactl.TopicOffsetMap {
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
	var TOM []kafkactl.TopicOffsetMap
	var count int
	tmMap := make(map[string][]kafkactl.TopicMeta)
	for _, tm := range t {
		tmMap[tm.Topic] = append(tmMap[tm.Topic], tm)
	}
	done := make(map[string]bool, len(tmMap))
	tomChan := make(chan []kafkactl.TopicOffsetMap, 100)
	for t, meta := range tmMap {
		if !done[t] {
			count++
			done[t] = true
		}
		go chanMakeTOM(client, meta, tomChan)
	}
	for i := 0; i < count; i++ {
		tom := <-tomChan
		TOM = append(TOM, tom...)
	}
	return TOM
}

func chanMakeTOM(client *kafkactl.KClient, tMeta []kafkactl.TopicMeta, tomChan chan []kafkactl.TopicOffsetMap) {
	tom := client.MakeTopicOffsetMap(tMeta)
	tomChan <- tom
	return
}

func refreshMetadata(topics ...string) {
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
	err = client.RefreshMetadata(topics...)
	if err != nil {
		log.Fatalf("Error refreshing topic metadata: %v\n", err)
	}
}
