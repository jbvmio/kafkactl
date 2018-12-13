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
	"os/signal"
	"strings"
	"time"

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
		closeFatal("Error: Invalid Partition Specified - %v\n", msg.Partition)
	}
	client.SaramaConfig().Producer.RequiredAcks = sarama.WaitForAll
	client.SaramaConfig().Producer.Return.Successes = true
	client.SaramaConfig().Producer.Return.Errors = true
	if msg.Partition == -1 {
		client.SaramaConfig().Producer.Partitioner = sarama.NewHashPartitioner
	} else {
		client.SaramaConfig().Producer.Partitioner = sarama.NewManualPartitioner
	}
	part, off, errd = client.SendMSG(msg)
	if errd != nil {
		closeFatal("Error sending message: %v\n", errd)
	}
	return
}

func sendMsgToPartitions(msgs []*kafkactl.Message) {
	client.SaramaConfig().Producer.RequiredAcks = sarama.WaitForAll
	client.SaramaConfig().Producer.Return.Successes = true
	client.SaramaConfig().Producer.Return.Errors = true
	client.SaramaConfig().Producer.Partitioner = sarama.NewManualPartitioner
	errd = client.SendMessages(msgs)
	if errd != nil {
		closeFatal("Error sending messages: %v\n", errd)
	}
}

func launchConsoleProducer(topic, key string, partitions ...int32) {
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
	sigChan := make(chan os.Signal, 1)
	stringChan := make(chan string)
	signal.Notify(sigChan, os.Interrupt)
	fmt.Printf("[kafkactl] # ")
	go consoleScanner(stringChan, sigChan)

ProducerLoop:
	for {
		select {
		case <-sigChan:
			fmt.Printf("signal: interrupt\n  Stopping Console Producer ...\n")
			sigChan <- os.Interrupt
			break ProducerLoop
		case line := <-stringChan:
			if strings.TrimSpace(line) != "" {
				msgs := makeMessages(topic, key, line, partitions...)
				errd = client.SendMessages(msgs)
				if errd != nil {
					log.Printf("Error sending messages: %v\n", errd)
				}
				time.Sleep(time.Millisecond * 200)
				fmt.Printf("[kafkactl] # ")
			} else {
				fmt.Printf("[kafkactl] # ")
			}
		}
	}
	fmt.Println("Stopped Console Producer")
}

func consoleScanner(stringChan chan string, sigChan chan os.Signal) {
	scanner := bufio.NewScanner(os.Stdin)
ConsoleLoop:
	for scanner.Scan() {
		select {
		case <-sigChan:
			fmt.Printf("signal: interrupt\n")
			break ConsoleLoop
		default:
			line := scanner.Text()
			stringChan <- line
		}
	}
	fmt.Println("STOPPED")
}
