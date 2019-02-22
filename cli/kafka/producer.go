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

package kafka

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/Shopify/sarama"
	kafkactl "github.com/jbvmio/kafka"
)

const delim string = "!7dd###755^^^557D$!"

type SendFlags struct {
	Key           string
	Value         string
	Delimiter     string
	DelimiterUsed bool
	Partitions    []string
	Partition     int32
	AllPartitions bool
	FromStdin     bool
	NoSplit       bool
	LineSplit     string
}

type sendData struct {
	key   string
	value string
}

func ProduceFromFile(flags SendFlags, data io.Reader, topics ...string) {
	if data == nil {
		data = bytes.NewBufferString(flags.Key + delim + flags.Value)
	}
	d, err := ioutil.ReadAll(data)
	handleC("Failed to read from stdin: %v\n", err)
	b := bytes.TrimSpace(d)
	var stringLines []string
	var allData []sendData
	switch true {
	case flags.FromStdin:
		if flags.NoSplit {
			stringLines = append(stringLines, string(b))
		} else {
			stringLines = strings.Split(string(b), flags.LineSplit)
		}
		subMatch := true
		switch subMatch {
		case flags.Value != "":
			for _, line := range stringLines {
				sd := sendData{
					key:   line,
					value: flags.Value,
				}
				allData = append(allData, sd)
			}
		case flags.Key != "":
			for _, line := range stringLines {
				sd := sendData{
					key:   flags.Key,
					value: line,
				}
				allData = append(allData, sd)
			}
		case flags.Delimiter != "":
			for _, line := range stringLines {
				keyVal := strings.Split(line, flags.Delimiter)
				if len(keyVal) == 2 {
					sd := sendData{
						key:   keyVal[0],
						value: keyVal[1],
					}
					allData = append(allData, sd)
				} else {
					out.Warnf("[WARN] Unable to parse data:\n%v\n", line)
				}
			}
		default:
			for _, line := range stringLines {
				sd := sendData{
					value: line,
				}
				allData = append(allData, sd)
			}
		}
	default:
		keyVal := strings.Split(string(b), delim)
		if len(keyVal) == 2 {
			sd := sendData{
				key:   keyVal[0],
				value: keyVal[1],
			}
			allData = append(allData, sd)
		} else {
			out.Warnf("[WARN] Unable to parse data:\n%v\n", keyVal)
		}
	}
	if len(allData) < 1 {
		closeFatal("Error: Missing Key/Values\n")
	}
	exact = true
	for _, topic := range topics {
		var msgs []*kafkactl.Message
		var parts []int32
		var hash bool
		switch true {
		case flags.AllPartitions:
			topicSummary := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
			if len(topicSummary) != 1 {
				closeFatal("Error isolating topic: %v\n", topic)
			}
			parts = topicSummary[0].Partitions
		case len(flags.Partitions) > 0:
			parts = validateParts(flags.Partitions)
		case flags.Partition == -1:
			parts = append(parts, flags.Partition)
			hash = true
		default:
			parts = append(parts, flags.Partition)
		}
		for _, sd := range allData {
			msgs = append(msgs, makeMessages(topic, sd.key, sd.value, parts...)...)
		}
		sendMessages(msgs, hash)
	}
}

func sendMessages(msgs []*kafkactl.Message, hash bool) {
	client.SaramaConfig().Producer.RequiredAcks = sarama.WaitForAll
	client.SaramaConfig().Producer.Return.Successes = true
	client.SaramaConfig().Producer.Return.Errors = true
	if !hash {
		client.SaramaConfig().Producer.Partitioner = sarama.NewManualPartitioner
	}
	errd = client.SendMessages(msgs)
	if errd != nil {
		closeFatal("Error sending messages: %v\n", errd)
	}
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
					out.Warnf("Error sending messages: %v", errd)
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
