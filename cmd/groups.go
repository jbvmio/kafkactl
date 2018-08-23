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
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
)

type groupMetadata struct {
	Group      string
	Topic      string
	Partitions string
	Host       string
	ClientID   string
	Member     string
	State      string
	Lag        int64
}

type groupListMeta struct {
	Group       string
	Type        string
	Coordinator string
}

func (gm groupMetadata) isTopic(t string) (bool, bool) {
	var (
		full    bool
		partial bool
	)
	if gm.Topic == t {
		full = true
	}
	if strings.Contains(gm.Topic, t) {
		partial = true
	}
	return full, partial
}

func getGroupList() []string {
	if verbose {
		sarama.Logger = log.New(os.Stdout, "[kafkactl] ", log.LstdFlags)
	}
	config := gimmeConf()
	client := gimmeClient(bootStrap, config)
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("Error closing client: %v\n", err)
		}
	}()
	allBrokers := client.Brokers()
	defer func() {
		for _, b := range allBrokers {
			if ok, _ := b.Connected(); ok {
				if err := b.Close(); err != nil {
					log.Fatalf("Error closing broker, %v: %v\n", b.Addr(), err)
				}
			}
		}
	}()
	req := sarama.ListGroupsRequest{}
	var groups []string
	for _, broker := range allBrokers {
		if ok, _ := broker.Connected(); !ok {
			if err := broker.Open(config); err != nil {
				log.Fatalf("Error connecting to broker: %v\n", err)
			}
		}
		grps, err := broker.ListGroups(&req)
		if err != nil {
			log.Fatalf("Error listing groups: %v\n", err)
		}
		for k := range grps.Groups {
			groups = append(groups, k)
		}
		if err = broker.Close(); err != nil {
			log.Fatalf("Error closing broker: %v\n", err)
		}
	}
	return groups
}

func getAllGroups(c sarama.Client) []groupListMeta {
	allBrokers := c.Brokers()
	defer func() {
		for _, b := range allBrokers {
			if ok, _ := b.Connected(); ok {
				if err := b.Close(); err != nil {
					log.Fatalf("Error closing broker, %v: %v\n", b.Addr(), err)
				}
			}
		}
	}()
	req := sarama.ListGroupsRequest{}
	var groups []groupListMeta

	for _, broker := range allBrokers {
		if ok, _ := broker.Connected(); !ok {
			if err := broker.Open(c.Config()); err != nil {
				log.Fatalf("Error connecting to broker: %v\n", err)
			}
		}
		coord := string(broker.Addr() + "/" + cast.ToString(broker.ID()))
		grps, err := broker.ListGroups(&req)
		if err != nil {
			log.Fatalf("Error listing groups: %v\n", err)
		}

		for k, v := range grps.Groups {
			glm := groupListMeta{
				Coordinator: coord,
			}
			glm.Group = k
			glm.Type = v
			groups = append(groups, glm)
		}

		if err = broker.Close(); err != nil {
			log.Fatalf("Error closing broker: %v\n", err)
		}
	}
	return groups
}

func describeGroups(c sarama.Client, g []string) []groupMetadata {
	div := 1
	if len(g) > numCPU {
		div = numCPU
	}

	dlChunks := stringsToMaps(g, div)

	var descriptions []*sarama.DescribeGroupsResponse
	descChan := make(chan *sarama.DescribeGroupsResponse, 100)
	doneChan := make(chan int, 100)
	var finished int

	for _, x := range dlChunks {
		go goDescribeGroups(c, x, descChan, doneChan)
	}

	for finished < len(dlChunks) {
		select {
		case desc := <-descChan:
			descriptions = append(descriptions, desc)
		case _ = <-doneChan:
			finished++
		}
	}

	gChunks := groupsToMaps(descriptions, div)
	var groups []groupMetadata
	metaChan := make(chan groupMetadata, 100)
	finished = 0

	for _, f := range gChunks {
		go filterGroupsDesc(f, metaChan, doneChan)
	}

	for finished < len(gChunks) {
		select {
		case mc, _ := <-metaChan:
			groups = append(groups, mc)
		case _ = <-doneChan:
			finished++
		}
	}
	return groups
}

func gimmeConf() *sarama.Config {
	conf := sarama.NewConfig()
	conf.ClientID = "kafkactl"
	conf.Version = sarama.V0_10_0_0
	err := conf.Validate()
	if err != nil {
		log.Fatalf("Invalid config: %v\n", err)
	}
	return conf
}

func gimmeClient(broker string, conf *sarama.Config) sarama.Client {
	client, err := sarama.NewClient([]string{broker}, conf)
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}
	return client
}

func dnsResolve(ip string) string {
	hn, err := net.LookupAddr(ip)
	if err == nil {
		return strings.Trim(hn[0], ".")
	}
	return ip
}

func groupsToMaps(someSlice []*sarama.DescribeGroupsResponse, chunks int) map[string][]*sarama.DescribeGroupsResponse {
	q, r := divMod(len(someSlice), chunks)
	myMap := make(map[string][]*sarama.DescribeGroupsResponse)
	label := "chunk"
	var perChunk = q
	var lastChunk = q + r
	var progress int
	var count int

	for progress < len(someSlice)-lastChunk {
		count++
		k := string(label + strconv.Itoa(count))
		v := progress + perChunk
		myMap[k] = someSlice[progress:v]
		progress = v
	}
	myMap[string(label+strconv.Itoa(chunks))] = someSlice[progress:]
	return myMap
}

func getController(b string) string {
	var controller string
	conf := gimmeConf()
	client := gimmeClient(b, conf)
	if c, err := client.Controller(); err == nil {
		controller = c.Addr()
	} else {
		log.Fatalf("Error retrieving active controller: %v\n", err)
	}
	return controller
}

func goDescribeGroups(c sarama.Client, g []string, descChan chan *sarama.DescribeGroupsResponse, isDone chan int) {
	for _, r := range g {
		request := sarama.DescribeGroupsRequest{
			Groups: []string{r},
		}
		broker, err := c.Coordinator(r)
		if err != nil {
			log.Fatalf("Error connecting to group coordinator: %v\n", err)
		}
		if ok, _ := broker.Connected(); !ok {
			if err := broker.Open(c.Config()); err != nil {
				if verbose {
					log.Printf("Error connecting broker, %v: %v\n", broker.Addr(), err)
				}
			}
		}
		desc, err := broker.DescribeGroups(&request)
		if err != nil {
			fmt.Printf("Error getting %v metadata: %v\n", r, err)
		} else {
			descChan <- desc
		}
	}
	isDone <- 1
}

func filterGroupsDesc(descriptions []*sarama.DescribeGroupsResponse, gmChan chan groupMetadata, doneChan chan int) {
	for _, desc := range descriptions {
		for _, g := range desc.Groups {
			for k, v := range g.Members {
				assignments, err := v.GetMemberAssignment()
				if err != nil {
					if err.Error() != `kafka: insufficient data to decode packet, more bytes expected` {
						log.Fatalf("Error getting assignments: %v\n", err)
					} else {
						isIdle = true
					}
				}
				for x, y := range assignments.Topics {
					var z []string
					for _, n := range y {
						z = append(z, cast.ToString(n))
					}
					parts := strings.Join(z[:], ",")
					grp := groupMetadata{}
					grp.Topic = x
					grp.Partitions = parts
					grp.Group = g.GroupId
					grp.ClientID = v.ClientId
					grp.Host = dnsResolve(strings.Trim(v.ClientHost, "/"))
					grp.Member = k
					grp.State = g.State
					if searchByGroup && searchByTopic {
						grp.Lag = getTotalLag(client, y, grp.Group, grp.Topic)
					}
					if searchByTopic {
						full, partial := grp.isTopic(targetTopic)
						if exact {
							if full {
								gmChan <- grp
							}
						} else {
							if partial {
								gmChan <- grp
							}
						}
					} else {
						gmChan <- grp
					}
				}
			}
		}

	}
	doneChan <- 1
}

func getTotalLag(client sarama.Client, parts []int32, group, topic string) int64 {
	omc, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		log.Fatalf("Error creating OMC: %v\n", err)
	}
	defer func() {
		go closeOMC(omc)
	}()
	var totalLag int64
	lagChan := make(chan int64, 100)
	for _, part := range parts {
		go getLag(omc, part, group, topic, lagChan)
	}
	for i := 0; i < len(parts); i++ {
		lag := <-lagChan
		totalLag = totalLag + lag
	}
	return totalLag
}

func getLag(omc sarama.OffsetManager, part int32, group, topic string, iChan chan int64) {
	pom, err := omc.ManagePartition(topic, part)
	if err != nil {
		log.Fatalf("Error creating POM: %v\n", err)
	}
	defer func() {
		go closePOM(pom)
	}()
	groupOffSet, _ := pom.NextOffset()
	partOffSet, err := client.GetOffset(topic, part, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to get offset for topic: %v, partition %v, Error: %v\n", topic, part, err)
	}
	if groupOffSet == -1 {
		groupOffSet = partOffSet
	}
	lag := partOffSet - groupOffSet
	iChan <- lag

}

func closePOM(pom sarama.PartitionOffsetManager) {
	err := pom.Close()
	if err != nil {
		fmt.Println("Error Closing POM: %V\n", err)
	}
}

func closeOMC(omc sarama.OffsetManager) {
	err := omc.Close()
	if err != nil {
		fmt.Println("Error Closing OMC: %V\n", err)
	}
}
