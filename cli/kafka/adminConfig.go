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
	"encoding/json"
	"strings"

	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/jbvmio/kafkactl/cli/zookeeper"

	"github.com/jbvmio/kafkactl"
)

type TopicConfigFlags struct {
	Config         string
	Value          string
	Configs        []string
	GetNonDefaults bool
	SetDefault     bool
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

const (
	topicConfVersion      = 1
	topicEntityVersion    = 2
	topicConfigPath       = `/config/topics/`
	topicEntityChangePath = `/config/changes/config_change_`
	topicEntityPath       = `topics/`
)

type emptyConfig struct {
	Version int               `json:"version"`
	Config  map[string]string `json:"config"`
}

type configChange struct {
	Version    int    `json:"version"`
	EntityPath string `json:"entity_path"`
}

func GetTopicConfigs(configs []string, topics ...string) []TopicConfig {
	var topicConfig []TopicConfig
	switch true {
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

func GetNonDefaultConfigs(configs []TopicConfig) []TopicConfig {
	var topicConfig []TopicConfig
	for _, tc := range configs {
		if !tc.Default {
			topicConfig = append(topicConfig, tc)
		}
	}
	return topicConfig
}

func SetTopicConfig(config, value string, topics ...string) []TopicConfig {
	var topicConfigs []TopicConfig
	exact = true
	for _, topic := range topics {
		ts := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		if len(ts) != 1 {
			out.Warnf("Error validating topic: %v", topic)
		}
		// Validation Check:
		tc := GetTopicConfigs([]string{config}, topic)

		err := client.SetTopicConfig(topic, config, value)
		handleC("Error setting configuration: %v\n", err)

		tc = GetTopicConfigs([]string{config}, topic)
		topicConfigs = append(topicConfigs, tc...)
	}
	return topicConfigs
}

func SetDefaultConfig(config string, topics ...string) []TopicConfig {
	var topicConfigs []TopicConfig
	exact = true
	for _, topic := range topics {
		ts := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		if len(ts) != 1 {
			out.Warnf("Error validating topic: %v", topic)
		}
		// Validation Check:
		preTC := GetTopicConfigs([]string{config}, topic)

		configs := make(map[string]string, 1)
		defaultConfig := emptyConfig{Version: topicConfVersion, Config: configs}
		j, err := json.Marshal(defaultConfig)
		handleC("Error on config marshal: %v", err)
		sent := zkCreateDefaultConfig(topic, j)
		if sent {
			postTC := GetTopicConfigs([]string{config}, topic)
			for _, tc1 := range postTC {
				for _, tc2 := range preTC {
					if tc1.Config == tc2.Config {
						if tc1.Default != tc2.Default {
							topicConfigs = append(topicConfigs, tc1)
						}
					}
				}
			}
		}
	}
	if len(topicConfigs) < 1 {
		closeFatal("No configuration changes needed.")
	}
	return topicConfigs
}

func zkCreateDefaultConfig(topic string, data []byte) bool {
	topicPath := topicConfigPath + topic
	handleC("%v", zookeeper.KafkaZK(targetContext, verbose))
	check, err := zookeeper.ZKCheckExists(topicPath)
	handleC("Error: %v", err)
	if !check {
		closeFatal("Unable to locate topic in zookeeper.")
	}
	currentValue, err := zookeeper.ZKGetValue(topicPath)
	handleC("Error: %v", err)
	if string(currentValue) == string(data) {
		out.Warnf("WARN: topic already default value: %v", topic)
		return false
	}
	zookeeper.ZKCreate(topicPath, true, true, data...)
	zkTopicChangeNotify(topic)
	return true
}

func zkTopicChangeNotify(topic string) {
	entityPath := topicEntityPath + topic
	change := configChange{
		Version:    topicEntityVersion,
		EntityPath: entityPath,
	}
	j, err := json.Marshal(change)
	handleC("Error on config marshal: %v", err)

	zookeeper.ZKSetSeq()
	zookeeper.ZKCreate(topicEntityChangePath, true, false, j...)
	zookeeper.ZKRemoveFlags()

	/*
		handleC("%v", zookeeper.KafkaZK(targetContext, verbose))
		zookeeper.ZKSetSeq()
		zookeeper.ZKCreate(path, false, false, []byte(value)...)
	*/
}
