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
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/jbvmio/kafkactl"
)

func createMsgSend(topic, key, value string, partition int32) {
	msg := &kafkactl.Message{
		Topic:     topic,
		Value:     []byte(value),
		Partition: partition,
	}
	if key != "" {
		msg.Key = []byte(key)
	}
	sendMsg(msg)
}

func createMsgSendParts(topic, key, value string, partitions ...int32) {
	var msgs []*kafkactl.Message
	for _, part := range partitions {
		msg := &kafkactl.Message{
			Topic:     topic,
			Value:     []byte(value),
			Partition: part,
		}
		if key != "" {
			msg.Key = []byte(key)
		}
		msgs = append(msgs, msg)
	}
	sendMsgToPartitions(msgs)
}

func makeMessages(topic, key, value string, partitions ...int32) []*kafkactl.Message {
	var msgs []*kafkactl.Message
	for _, part := range partitions {
		msg := &kafkactl.Message{
			Topic:     topic,
			Value:     []byte(value),
			Partition: part,
		}
		if key != "" {
			msg.Key = []byte(key)
		}
		msgs = append(msgs, msg)
	}
	return msgs
}

func sendMsg(msg *kafkactl.Message) (part int32, off int64) {
	if msg.Partition < -1 {
		log.Fatalf("Error: Invalid Partition Specified - %v\n", msg.Partition)
	}
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
	client.SaramaConfig().Producer.RequiredAcks = sarama.WaitForAll
	client.SaramaConfig().Producer.Return.Successes = true
	client.SaramaConfig().Producer.Return.Errors = true
	if msg.Partition == -1 {
		client.SaramaConfig().Producer.Partitioner = sarama.NewHashPartitioner
	} else {
		client.SaramaConfig().Producer.Partitioner = sarama.NewManualPartitioner
	}
	part, off, err = client.SendMSG(msg)
	if err != nil {
		log.Fatalf("Error sending message: %v\n", err)
	}
	return
}

func sendMsgToPartitions(msgs []*kafkactl.Message) {
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
	client.SaramaConfig().Producer.RequiredAcks = sarama.WaitForAll
	client.SaramaConfig().Producer.Return.Successes = true
	client.SaramaConfig().Producer.Return.Errors = true
	client.SaramaConfig().Producer.Partitioner = sarama.NewManualPartitioner
	err = client.SendMessages(msgs)
	if err != nil {
		log.Fatalf("Error sending messages: %v\n", err)
	}
}

func launchConsoleProducer(topic, key string, partitions ...int32) {
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
	client.SaramaConfig().Producer.RequiredAcks = sarama.WaitForAll
	client.SaramaConfig().Producer.Return.Successes = true
	client.SaramaConfig().Producer.Return.Errors = true
	if len(partitions) == 0 {
		client.SaramaConfig().Producer.Partitioner = sarama.NewHashPartitioner
		partitions = []int32{-1}
		fmt.Printf("Started Console Producer ... Partitions: Hashing\n\n")
	} else {
		client.SaramaConfig().Producer.Partitioner = sarama.NewManualPartitioner
		fmt.Printf("Started Console Producer ... Partitions: %v\n\n", partitions)
	}
	fmt.Printf("[kafkactl] ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			msgs := makeMessages(topic, key, line, partitions...)
			err = client.SendMessages(msgs)
			if err != nil {
				log.Printf("Error sending messages: %v\n", err)
			}
			fmt.Printf("[kafkactl] ")
		} else {
			fmt.Printf("[kafkactl] ")
		}
	}
	fmt.Println("Stopping Console Producer")
}
