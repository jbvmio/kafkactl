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
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
)

var (
	targetKey   string
	targetMsg   string
	partitioner string
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "Produce messages to a Kafka topic",
	Long: `Produces messages to a Kafka topic.
  Either by:
    Command arguments,
    Stdin,
    A Console Producer (launched if no arguments or Stdin is found.`,
	Aliases: []string{"pro", "producer"},
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			sarama.Logger = log.New(os.Stdout, "[kafkactl] ", log.LstdFlags)
		}
		conf := gimmeConf()
		conf.Producer.RequiredAcks = sarama.WaitForAll
		conf.Producer.Return.Successes = true

		switch partitioner {
		case "":
			if targetPartition >= 0 {
				conf.Producer.Partitioner = sarama.NewManualPartitioner
			} else {
				conf.Producer.Partitioner = sarama.NewHashPartitioner
			}
		case "hash":
			conf.Producer.Partitioner = sarama.NewHashPartitioner
		case "random":
			conf.Producer.Partitioner = sarama.NewRandomPartitioner
		case "manual":
			conf.Producer.Partitioner = sarama.NewManualPartitioner
			if targetPartition == -1 {
				log.Fatalf("--partition is required when partitioning manually")
			}
		default:
			log.Fatalf("Partitioner %s not supported.", partitioner)
		}

		message := &sarama.ProducerMessage{
			Topic:     targetTopic,
			Partition: int32(targetPartition),
		}
		if targetKey != "" {
			message.Key = sarama.StringEncoder(targetKey)
		}
		producer, err := sarama.NewSyncProducer([]string{bootStrap}, conf)
		if err != nil {
			log.Fatalf("Error initializing Kafka Producer: %v\n", err)
		}
		defer func() {
			if err := producer.Close(); err != nil {
				log.Println("Failed to close Kafka producer cleanly:", err)
			}
		}()
		if targetMsg != "" {
			message.Value = sarama.StringEncoder(targetMsg)
		} else if stdinAvailable() {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("Failed to read data from the standard input: %v\n", err)
			}
			bits := bytes.TrimSpace(b)
			message.Value = sarama.ByteEncoder(bits)
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := scanner.Text()
				//l := strings.TrimSpace(line)
				if strings.TrimSpace(line) != "" {
					message.Value = sarama.StringEncoder(line)
					partition, offset, err := producer.SendMessage(message)
					if err != nil {
						log.Fatalf("Failed to produce message: %v\n", err)
					}
					if verbose {
						log.Printf("topic=%s\tpartition=%d\toffset=%d\n", targetTopic, partition, offset)
					}
				}
			}
		}
		partition, offset, err := producer.SendMessage(message)
		if err != nil {
			log.Fatalf("Failed to produce message: %v", err)
		}
		if verbose {
			log.Printf("topic=%s\tpartition=%d\toffset=%d\n", targetTopic, partition, offset)
		}
	},
}

func init() {
	rootCmd.AddCommand(produceCmd)
	produceCmd.Flags().StringVarP(&targetKey, "key", "k", "", "Optional key to produce/send")
	produceCmd.Flags().StringVarP(&targetMsg, "message", "m", "", "Message to produce/send, Wins over Stdin.")
	produceCmd.Flags().StringVar(&partitioner, "partitioner", "", `The partitioning scheme to use. Can be 'hash', 'manual', or 'random'`)
	produceCmd.Flags().IntVarP(&targetPartition, "partition", "p", -1, "Specific partition to produce/send")
	produceCmd.MarkPersistentFlagRequired("topic")

}
