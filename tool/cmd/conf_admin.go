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
	"log"

	"github.com/jbvmio/kafkactl"
)

// TopicConfig struct def:
type TopicConfig struct {
	Topic     string
	Config    string
	Value     string
	ReadOnly  bool
	Default   bool
	Sensitive bool
}

func getTopicConfig(topics, configNames []string) []TopicConfig {

	/*
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
	*/

	var topicConfig []TopicConfig
	for _, t := range topics {
		c, err := client.GetTopicConfig(t, configNames...)
		if err != nil {
			log.Fatalf("Error getting config for topic %v: %v\n", t, err)
		}
		for _, v := range c {
			tc := TopicConfig{
				Topic:     t,
				Config:    v.Name,
				Value:     v.Value,
				ReadOnly:  v.ReadOnly,
				Default:   v.Default,
				Sensitive: v.Sensitive,
			}
			topicConfig = append(topicConfig, tc)
		}
	}
	return topicConfig
}

func searchTopicConfig(topic string, configNames ...string) []TopicConfig {
	var topics []string
	ts := kafkactl.GetTopicSummaries(searchTopicMeta(topic))
	if len(ts) < 1 {
		log.Fatalf("unable to locate specified topic: %v\n", topic)
	}
	for _, t := range ts {
		topics = append(topics, t.Topic)
	}
	return getTopicConfig(topics, configNames)
}

func setTopicConfig(topic, configName, value string) error {
	if configName == "" || value == "" {
		log.Fatalf("Error: Missing Key and/or Value\n")
	}
	exact = true
	ts := kafkactl.GetTopicSummaries(searchTopicMeta(topic))
	if len(ts) != 1 {
		log.Fatalf("Error validating topic: %v\n", topic)
	}
	return client.SetTopicConfig(topic, configName, value)
}
