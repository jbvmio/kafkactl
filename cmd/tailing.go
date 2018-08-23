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
	"strconv"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func readLast(p kafka.Partition) {
	ctx := context.Background()
	conn, err := kafka.DialLeader(ctx, "tcp", string(p.Leader.Host+":"+strconv.Itoa(p.Leader.Port)), p.Topic, p.ID)
	if err != nil {
		log.Fatalf("Cannot establish a connection to %v\n", bootStrap)
	}
	defer conn.Close()
	last, _ := conn.ReadLastOffset()
	if last > 0 {
		_, err = conn.Seek(offset, 2)
		if err != nil {
			if !strings.Contains(err.Error(), "Offset Out Of Range") {
				log.Fatalf("Could not tail %v from last offset %v: %v\n", offset, last, err)
			}
		}
	}
	var count int64
	for count <= offset {
		msg, _ := conn.ReadMessage(20480)
		fmt.Println(string(msg.Value))
		count++
	}

	if verbose == true {
		o, _ := conn.ReadLastOffset()
		if o == last {
			fmt.Printf("[kafkactl] Reached end of topic %v [%v] > offset %v\n", p.Topic, p.ID, o)
		}
	}
}

func tailPartition(p kafka.Partition, d *kafka.Dialer) {
	config := kafka.ReaderConfig{
		Brokers:         []string{string(p.Leader.Host + ":" + strconv.Itoa(p.Leader.Port))},
		Topic:           p.Topic,
		Partition:       p.ID,
		MinBytes:        1,
		MaxBytes:        10e6, // 10MB
		Dialer:          d,
		QueueCapacity:   1,
		MaxWait:         15 * time.Minute,
		ReadLagInterval: -1,
	}
	if verbose {
		config.Logger = log.New(os.Stdout, "[kafkactl][INFO] ", log.LstdFlags)
		config.ErrorLogger = log.New(os.Stdout, "[kafkactl][ERROR] ", log.LstdFlags)
	}
	ctx := context.Background()
	r := kafka.NewReader(config)
	r.SetOffset(-2)
	defer r.Close()
	readLast(p)
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		fmt.Println(string(m.Value))
		if verbose == true {
			conn, _ := kafka.DialLeader(ctx, "tcp", string(p.Leader.Host+":"+strconv.Itoa(p.Leader.Port)), m.Topic, m.Partition)
			o, _ := conn.ReadLastOffset()
			conn.Close()
			if o == m.Offset+1 {
				fmt.Printf("[kafkactl] Reached end of topic %v [%v] > offset %v\n", m.Topic, m.Partition, o)
			}
		}
	}
	wg.Done()
}

func divMod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func stringsToMaps(someSlice []string, chunks int) map[string][]string {
	q, r := divMod(len(someSlice), chunks)
	myMap := make(map[string][]string)
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
