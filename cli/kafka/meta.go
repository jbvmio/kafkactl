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
	"strings"

	kafkactl "github.com/jbvmio/kafka"
)

// ClusterDetails prints details for the current context.
func ClusterDetails() {
	meta, err := client.GetClusterMeta()
	if err != nil {
		closeFatal("Error getting cluster metadata: %v\n", err)
	}
	kafkaVer, _ := kafkactl.MatchKafkaVersion(getKafkaVersion(meta.APIMaxVersions))
	c, err := client.Controller()
	if err != nil {
		closeFatal("Error obtaining controller: %v\n", err)
	}
	if len(meta.ErrorStack) > 0 {
		fmt.Println("ERRORs:")
		for _, e := range meta.ErrorStack {
			fmt.Printf(" %v\n", e)
		}
	}
	fmt.Println("\nBrokers: ", meta.BrokerCount())
	fmt.Println(" Topics: ", meta.TopicCount())
	fmt.Println(" Groups: ", meta.GroupCount())
	fmt.Printf("\nClient:  (Using: %v)\n", clientVer)
	fmt.Printf("Cluster: (Kafka: %v)\n", kafkaVer)
	for _, b := range meta.Brokers {
		if strings.Contains(b, c.Addr()) {
			fmt.Println("*", b)
		} else {
			fmt.Println(" ", b)
		}
	}
	fmt.Printf("\n(*)Controller\n\n")
}

func MetaData() kafkactl.ClusterMeta {
	meta, err := client.GetClusterMeta()
	handleC("Error getting cluster metadata: %v\n", err)
	return meta
}
