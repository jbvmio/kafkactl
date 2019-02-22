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
	"fmt"
	"os"
	"os/signal"
	"strings"

	kafkactl "github.com/jbvmio/abstraction/kafka"
)

func launchCG(groupID string, debug bool, topics ...string) {
	consumer, err := client.NewConsumerGroup(groupID, debug, topics...)
	if err != nil {
		closeFatal("Error creating consumer group: %v\n", err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	var commitErr error
ConsumeLoop:
	for {
		select {
		case part, ok := <-consumer.CG().Partitions():
			if !ok {
				fmt.Println("ERROR ON MSG CHAN")
				return
			}
			go consumer.ConsumeMessages(part, func(msg *kafkactl.Message) bool {
				fmt.Printf("[%v] %v > %s\n", msg.Partition, msg.Offset, msg.Value)
				return true
			})
		case debugAvailable, ok := <-consumer.DebugAvailale():
			if !ok {
				fmt.Println("ERROR ON DEBUG CHAN")
				return
			}
			if debugAvailable {
				go consumer.ProcessDEBUG(func(debug *kafkactl.DEBUG) bool {
					if debug.HasData {
						if debug.IsErr {
							fmt.Printf("ERROR: %v\n", debug.Err)
							return true
						}
						if debug.IsNote {
							if strings.Contains(debug.Type, "start") {
								client.Logf("[NOTICE*][%v]\n", debug.Type)
							}
							if strings.Contains(debug.Type, "OK") {
								for k := range debug.Released {
									client.Logf("[NOTICE*][%v] Partitions Removed [%v] %v\n", debug.Type, k, debug.Released[k])
								}
								for k := range debug.Claimed {
									client.Logf("[NOTICE*][%v] Partitions Added [%v] %v\n", debug.Type, k, debug.Claimed[k])
								}
								for k := range debug.Current {
									client.Logf("[NOTICE*][%v] Partitions Assigned [%v] %v\n", debug.Type, k, debug.Current[k])
								}
							}
							if strings.Contains(debug.Type, "error") {
								client.Logf("[ERROR*][%v]:\n%+v\n%+v\n%+v\n", debug.Type, debug.Released, debug.Claimed, debug.Current)
							}
							return true
						}
					}
					return true
				})
			}
		case <-signals:
			commitErr = consumer.CG().CommitOffsets()
			break ConsumeLoop
		}
	}
	if commitErr != nil {
		fmt.Println("Warning > Offset Commit Error:", commitErr)
	}
	if debug {
		if err := consumer.Close(); err != nil {
			closeFatal("Error closing Consumer: %v\n", err)
		}
	}
}
