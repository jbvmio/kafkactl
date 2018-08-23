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
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	kafka "github.com/segmentio/kafka-go"
)

func createTopic(ctx context.Context, topics []string, p, r int) {
	conf := kafka.TopicConfig{}
	conf.NumPartitions = p
	conf.ReplicationFactor = r
	controller := getController(bootStrap)
	conn, err := kafka.DialContext(ctx, "tcp", controller)
	if err != nil {
		log.Fatalf("Cannot establish a connection to %v\n", controller)
	}
	defer conn.Close()
	for _, t := range topics {
		conf.Topic = t
		err = conn.CreateTopics(conf)
		if err != nil {
			log.Printf("Could not create \"%v\" %v\n", t, err)
		} else {
			fmt.Println("Created", t)
		}
	}
}

func deleteTopic(ctx context.Context, topics []string) {
	controller := getController(bootStrap)
	conn, err := kafka.DialContext(ctx, "tcp", controller)
	if err != nil {
		log.Fatalf("Cannot establish a connection to %v\n", controller)
	}
	defer conn.Close()
	if confirmation == true {
		for _, t := range topics {
			err := conn.DeleteTopics(t)
			if err != nil {
				log.Printf("Could not delete Topic \"%v\" %v\n", t, err)
			} else {
				fmt.Println("Deleted", t)
			}
		}
	} else {
		fmt.Println("You are attempting to delete", topics[:])
		fmt.Println("If this is what you really want then re-run the command with \"--confirmed\"")
		fmt.Println()
	}
}

func getAllTopics() []string {
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[kafkactl] ", log.LstdFlags)
	}
	config := gimmeConf()
	client := gimmeClient(bootStrap, config)
	allTopics, err := client.Topics()
	if err != nil {
		log.Fatalf("Error collecting topics: %v\n", err)
	}
	return allTopics
}

func searchForTopics(topics, searches []string) []string {
	var results []string
	for _, s := range searches {
		for _, t := range topics {
			if exact {
				if t == s {
					results = append(results, t)
				}
			} else if strings.Contains(t, s) {
				results = append(results, t)
			}
		}
	}
	return results
}

type topicSummary struct {
	Topic      string
	Partitions int
	RFactor    int
	ISR        int
}

func getTopicSummaries(ctx context.Context, topics []string) []topicSummary {
	div := 1
	if len(topics) > numCPU {
		div = numCPU
	}
	tChunks := stringsToMaps(topics, div)

	var summaries []topicSummary
	sumChan := make(chan topicSummary, 100)
	doneChan := make(chan int, 100)
	var finished int

	for _, t := range tChunks {
		go goGetTopicSummaries(ctx, t, sumChan, doneChan)
	}

	for finished < len(tChunks) {
		select {
		case summary := <-sumChan:
			summaries = append(summaries, summary)
		case _ = <-doneChan:
			finished++

		}
	}
	return summaries
}

func goGetTopicSummaries(ctx context.Context, topics []string, sumChan chan topicSummary, doneChan chan int) {
	dialer := &kafka.Dialer{
		ClientID:  "kafkactl",
		DualStack: true,
	}
	for _, topic := range topics {
		var p []kafka.Partition
		i, err := dialer.LookupPartitions(ctx, "tcp", bootStrap, topic)
		if err != nil {
			log.Fatalf("Error looking up topics: %v\n", err)
		}
		for _, x := range i {
			p = append(p, x)
		}
		var ts = topicSummary{
			Topic:      topic,
			Partitions: len(p),
		}
		var replicas int
		var isr int
		for _, i := range p {
			for range i.Replicas {
				replicas++
			}
			for range i.Isr {
				isr++
			}
		}
		ts.RFactor = replicas / ts.Partitions
		ts.ISR = isr / ts.Partitions
		sumChan <- ts
	}
	doneChan <- 1
}
