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
	"sort"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/fatih/color"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/cast"
)

type datameta struct {
	Topic     string
	Partition int
	Replicas  string
	ISRs      string
	Leader    string
}

func getMetaConn(ctx context.Context) {
	conn, err := kafka.DialContext(ctx, "tcp", bootStrap)
	if err != nil {
		log.Fatalf("Cannot establish a connection to %v\n", bootStrap)
	}
	defer conn.Close()
	var p []kafka.Partition
	if len(topicList) != 0 {
		if targetTopic != "" {
			topicList = append(topicList, targetTopic)
		}
		p, err = conn.ReadPartitions(topicList[0:]...)
	} else if targetTopic != "" {
		p, err = conn.ReadPartitions(targetTopic)
	} else {
		p, err = conn.ReadPartitions()
	}
	if err != nil {
		log.Fatalf("Could not request Metadata: %v\n", err)
	}
	if len(p) == 0 {
		log.Fatalf("Topic %v does not exist or has no partitions.\n", targetTopic)
	}
	var dm []datameta
	for _, i := range p {
		md := datameta{}
		var replicas []string
		var isr []string
		for _, r := range i.Replicas {
			x := strconv.Itoa(r.ID)
			replicas = append(replicas, x)
		}
		for _, b := range i.Isr {
			x := strconv.Itoa(b.ID)
			isr = append(isr, x)
		}
		reps := strings.Join(replicas[:], ",")
		isrs := strings.Join(isr[:], ",")

		md.Topic = i.Topic
		md.Partition = i.ID
		md.Replicas = reps
		md.ISRs = isrs
		md.Leader = string(i.Leader.Host + "/" + strconv.Itoa(i.Leader.ID))
		dm = append(dm, md)
	}

	formatOutput(dm)

}

func getMetaDial(ctx context.Context) {
	var typo bool
	dialer := &kafka.Dialer{
		ClientID:  "kafkactl",
		DualStack: true,
	}
	var p []kafka.Partition
	if len(topicList) != 0 {
		if targetTopic != "" {
			topicList = append(topicList, targetTopic)
		}
		for _, topic := range topicList {
			i, err := dialer.LookupPartitions(ctx, "tcp", bootStrap, topic)
			if err != nil {
				log.Printf("Could not find topic %v on %v: %v\n", topic, bootStrap, err)
				typo = true
			} else if len(i) == 0 {
				log.Printf("Could not find topic %v on %v\n", topic, bootStrap)
				typo = true
			} else {
				for _, x := range i {
					p = append(p, x)
				}
			}
		}
	} else if targetTopic != "" {
		var err error
		p, err = dialer.LookupPartitions(ctx, "tcp", bootStrap, targetTopic)
		if err != nil {
			log.Fatalf("Could not find topic %v on %v: %v\n", targetTopic, bootStrap, err)
		}
	} else {
		conn, err := kafka.DialContext(ctx, "tcp", bootStrap)
		if err != nil {
			log.Fatalf("Cannot establish a connection to %v\n", bootStrap)
		}
		defer conn.Close()
		p, err = conn.ReadPartitions()
		if err != nil {
			log.Fatalf("Could not request Metadata: %v\n", err)
		}
	}
	if len(p) == 0 {
		if len(topicList) != 0 {
			log.Printf("%v do not exist or have no partitions.\n\n", topicList[:])
			fmt.Println("If", color.GreenString("auto.create.topics.enable"), "is set, they will have been created.")
			os.Exit(1)
		} else {
			log.Printf("Topic %v does not exist or has no partitions.\n\n", targetTopic)
			fmt.Println("If", color.GreenString("auto.create.topics.enable"), "is set, it will have been created.")
			os.Exit(1)
		}
	}
	var dm []datameta
	for _, i := range p {
		md := datameta{}
		var replicas []string
		var isr []string
		for _, r := range i.Replicas {
			x := strconv.Itoa(r.ID)
			replicas = append(replicas, x)
		}
		for _, b := range i.Isr {
			x := strconv.Itoa(b.ID)
			isr = append(isr, x)
		}
		reps := strings.Join(replicas[:], ",")
		isrs := strings.Join(isr[:], ",")

		md.Topic = i.Topic
		md.Partition = i.ID
		md.Replicas = reps
		md.ISRs = isrs
		md.Leader = string(i.Leader.Host + "/" + strconv.Itoa(i.Leader.ID))
		dm = append(dm, md)
	}

	formatOutput(dm)
	if typo {
		fmt.Println("If", color.GreenString("auto.create.topics.enable"), "is set, any topics not found will have been created.")
	}
}

func getTopicMetadata(ctx context.Context, topics []string) {
	dialer := &kafka.Dialer{
		ClientID:  "kafkactl",
		DualStack: true,
	}
	var p []kafka.Partition
	for _, topic := range topics {
		i, err := dialer.LookupPartitions(ctx, "tcp", bootStrap, topic)
		if err != nil {
			log.Printf("Could not find topic %v on %v: %v\n", topic, bootStrap, err)
		} else if len(i) == 0 {
			log.Printf("Could not find topic %v on %v\n", topic, bootStrap)
		} else {
			for _, x := range i {
				p = append(p, x)
			}
		}
	}
	var dm []datameta
	for _, i := range p {
		md := datameta{}
		var replicas []string
		var isr []string
		for _, r := range i.Replicas {
			x := strconv.Itoa(r.ID)
			replicas = append(replicas, x)
		}
		for _, b := range i.Isr {
			x := strconv.Itoa(b.ID)
			isr = append(isr, x)
		}
		reps := strings.Join(replicas[:], ",")
		isrs := strings.Join(isr[:], ",")

		md.Topic = i.Topic
		md.Partition = i.ID
		md.Replicas = reps
		md.ISRs = isrs
		md.Leader = string(i.Leader.Host + "/" + strconv.Itoa(i.Leader.ID))
		dm = append(dm, md)
	}
	formatOutput(dm)
}

type clusterMeta struct {
	Brokers    []string
	Topics     []string
	Controller int32
}

func (cm clusterMeta) BrokerCount() int {
	return len(cm.Brokers)
}

func (cm clusterMeta) TopicCount() int {
	return len(cm.Topics)
}

func getClusterMeta() clusterMeta {
	res := reqMetadata()
	cm := clusterMeta{}
	cm.Controller = res.ControllerID
	for _, b := range res.Brokers {
		id := b.ID()
		addr := b.Addr()
		broker := string(addr + "/" + cast.ToString(id))
		cm.Brokers = append(cm.Brokers, broker)
	}
	for _, t := range res.Topics {
		cm.Topics = append(cm.Topics, t.Name)
	}
	sort.Strings(cm.Brokers)
	sort.Strings(cm.Topics)
	return cm
}

func reqMetadata() *sarama.MetadataResponse {
	var req = sarama.MetadataRequest{
		AllowAutoTopicCreation: false,
	}
	broker := sarama.NewBroker(bootStrap)
	conf := gimmeConf()
	if err := broker.Open(conf); err != nil {
		log.Fatalf("Error Opening Broker Connection: %v\n", err)
	}
	defer broker.Close()
	res, err := broker.GetMetadata(&req)
	if err != nil {
		log.Fatalf("Error obtaining metadata: %v\n", err)
	}
	return res
}
