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
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/jbvmio/kafkactl/cli/x"
	"github.com/spf13/cast"
)

type topicStdinData struct {
	topic     string
	partition int32
}

// ParseTopicStdin parses Stdin passed from kafkactl topic metadata
func ParseTopicStdin(r io.Reader) []topicStdinData {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		closeFatal("Failed to read from stdin: %v\n", err)
	}
	bits := bytes.TrimSpace(b)
	lines := string(bits)

	var topicData []topicStdinData
	a := strings.Split(lines, "\n")
	headers := strings.Fields(strings.TrimSpace(a[0]))
	if len(headers) < 3 {
		closeFatal("Invalid Stdin Passed")
	}
	if headers[0] != "TOPIC" || headers[1] != "PART" || headers[2] != "OFFSET" {
		closeFatal("Best to pass stdin through kafkactl itself.")
	}
	for _, b := range a[1:] {
		td := topicStdinData{}
		b := strings.TrimSpace(b)
		td.topic = x.CutField(b, 1)
		td.partition = cast.ToInt32(x.CutField(b, 2))
		topicData = append(topicData, td)
	}
	return topicData
}

func MovePartitionsStdin(moveData []topicStdinData, brokers []int32) RAPartList {
	var raparts []RAPartition
	for _, tm := range moveData {
		rap := RAPartition{
			Topic:     tm.topic,
			Partition: tm.partition,
			Replicas:  brokers,
		}
		raparts = append(raparts, rap)
	}
	rapList := RAPartList{
		Version:    1,
		Partitions: raparts,
	}
	return rapList
}
