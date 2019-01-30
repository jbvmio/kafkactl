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
	"strings"

	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/jbvmio/kafkactl"
)

type TopicConfigFlags struct {
	Config  string
	Value   string
	Configs []string
}

// TopicConfig struct def:
type TopicConfig struct {
	Topic     string
	Config    string
	Value     string
	ReadOnly  bool
	Default   bool
	Sensitive bool
}

func GetTopicConfigs(configs []string, topics ...string) []TopicConfig {
	var topicConfig []TopicConfig
	match := true
	switch match {
	case exact || len(configs) < 1:
		for _, t := range topics {
			c, err := client.GetTopicConfig(t, configs...)
			if err != nil {
				closeFatal("Error getting config for topic %v: %v\n", t, err)
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
	default:
		for _, t := range topics {
			c, err := client.GetTopicConfig(t)
			if err != nil {
				closeFatal("Error getting config for topic %v: %v\n", t, err)
			}
			for _, config := range configs {
				for _, v := range c {

					if strings.Contains(v.Name, config) {
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
			}
		}
	}
	if len(topicConfig) < 1 {
		closeFatal("Config not found: %v\n", configs)
	}
	return topicConfig
}

func SearchTopicConfigs(configs []string, topics ...string) []TopicConfig {
	var tops []string
	for _, topic := range topics {
		ts := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		if len(ts) < 1 {
			closeFatal("unable to isolate topic: %v\n", topic)
		}
		for _, t := range ts {
			tops = append(tops, t.Topic)
		}
	}
	return GetTopicConfigs(configs, tops...)
}

func SetTopicConfig(config, value string, topics ...string) []TopicConfig {
	var topicConfigs []TopicConfig
	exact = true
	for _, topic := range topics {
		ts := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		if len(ts) != 1 {
			out.Warnf("Error validating topic: %v", topic)
		}
		tc := GetTopicConfigs([]string{config}, topic)

		err := client.SetTopicConfig(topic, config, value)
		handleC("Error setting configuration: %v\n", err)

		tc = GetTopicConfigs([]string{config}, topic)
		topicConfigs = append(topicConfigs, tc...)
	}
	return topicConfigs
}
