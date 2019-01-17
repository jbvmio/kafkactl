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
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jbvmio/kafkactl"
	yaml "gopkg.in/yaml.v2"
)

var (
	client        *kafkactl.KClient
	errd          error
	clientTimeout = (time.Second * 5)
)

type apiVersion struct {
	key         int16
	maxVer      int16
	description string
}

func launchClient() {
	if verbose {
		kafkactl.Logger("")
	}
	conf, err := kafkactl.GetConf()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	conf.Version, err = sarama.ParseKafkaVersion(clientVer)
	if err != nil {
		conf.Version = kafkactl.MinKafkaVersion
	}
	conf.Net.DialTimeout = clientTimeout
	conf.Net.ReadTimeout = clientTimeout
	conf.Net.WriteTimeout = clientTimeout
	conf.Metadata.Retry.Max = 1
	client, errd = kafkactl.NewCustomClient(conf, kafkaBrokers...)
	if errd != nil {
		log.Fatalf("Error: %v\n", errd)
	}
}

func launchClient2() {
	if verbose {
		kafkactl.Logger("")
	}
	client, errd = kafkactl.NewClient(bootStrap)
	if errd != nil {
		log.Fatalf("Error: %v\n", errd)
	}
}

func closeFatal(format string, msg ...interface{}) {
	if connected := client.IsConnected(); connected {
		closeClient()
	}
	log.Fatalf(format, msg...)
}

func closeClient() {
	if errd = client.Close(); errd != nil {
		log.Fatalf("Error closing client: %v\n", errd)
	}
}

func validateBootStrap() {
	if bootStrap == "" {
		if cfgFile {
			bootStrap, errd = kafkactl.ReturnFirstValid(kafkaBrokers...)
			if errd != nil {
				log.Fatalf("Error reading config: %v\n", errd)
			}
		}
	}
	if !strings.Contains(bootStrap, ":") {
		bootStrap = net.JoinHostPort(bootStrap, bsport)
	}
	if bootStrap != "" {
		kafkaBrokers = []string{bootStrap}
	}
	if clientVer == "query" {
		apiVer, err := kafkactl.BrokerAPIVersions(bootStrap)
		if err != nil {
			fmt.Println("ERR", err)
		}
		clientVer = getKafkaVersion(apiVer)
		ver, _ := kafkactl.MatchKafkaVersion(clientVer)
		kafkaVer = ver.String()
	}
	return
}

func getKafkaVersion(apiKeys map[int16]int16) string {
	match := true
	switch match {
	case apiKeys[kafkactl.APIKeyOffsetForLeaderEpoch] == 2:
		return "2.1.0"
	case apiKeys[kafkactl.APIKeyOffsetForLeaderEpoch] == 1:
		return "2.0.0"
	case apiKeys[kafkactl.APIKeyFetch] == 7:
		return "1.1.0"
	case apiKeys[kafkactl.APIKeyFetch] == 6:
		return "1.0.0"
	case apiKeys[kafkactl.APIKeyFetch] == 5:
		return "0.11.0.0"
	case apiKeys[kafkactl.APIKeyUpdateMetadata] == 3:
		return "0.10.2.0"
	case apiKeys[kafkactl.APIKeyFetch] == 3:
		return "0.10.1.0"
	}
	return "UnknownVersion"
}

func getAPIVersions() []apiVersion {
	var apis []apiVersion
	max, err := client.GetAPIVersions()
	if err != nil {
		closeFatal("Unable to Retrieve API Versions.")
	}
	for k := range max {
		api := apiVersion{
			key:         k,
			maxVer:      max[k],
			description: kafkactl.APIDescriptions[k],
		}
		apis = append(apis, api)
	}
	sort.Slice(apis, func(i, j int) bool {
		return apis[i].key < apis[j].key
	})
	return apis
}

// Config contains a collection of cluster entries
type Config struct {
	Current string  `json:"current" yaml:"current"`
	Entries []Entry `json:"entries" yaml:"entries"`
}

// Entry contains kafka and burrow node details for a cluster
type Entry struct {
	Name      string   `json:"name" yaml:"name"`
	Kafka     []string `json:"kafka yaml:"kafka"`
	Burrow    []string `json:"burrow" yaml:"burrow"`
	Zookeeper []string `json:"zookeeper" yaml:"zookeeper"`
}

func getEntries(path string) (kafka, burrow, zookeeper []string) {
	entry := getCurrentEntry(path)
	return entry.Kafka, entry.Burrow, entry.Zookeeper
}

func tryByBroker(broker, path string) (kafka, burrow, zookeeper []string) {
	entry := getEntriesByBroker(broker, path)
	return entry.Kafka, entry.Burrow, entry.Zookeeper
}

func printCurrentEntry(path string) {
	current := getCurrentEntry(path)
	y, err := yaml.Marshal(current)
	if err != nil {
		log.Fatalf("Error marshaling current entry: %v\n", err)
	}
	fmt.Printf("\nCURRENT:\n%s\n", y)
}

func getCurrentEntry(path string) Entry {
	return getCurrentFromConfig(returnConfig(readConfig(path)))
}

func getEntriesByBroker(broker, path string) Entry {
	return findEntriesByBroker(broker, returnConfig(readConfig(path)))
}

func getCurrentFromConfig(config Config) Entry {
	current := config.Current
	for _, e := range config.Entries {
		if e.Name == current {
			return e
		}
	}
	log.Fatalf("Error reading current entry: Not Found")
	return Entry{}
}

// Ensure broker validation happens before running this
func findEntriesByBroker(broker string, config Config) Entry {
	if strings.Contains(broker, ":") {
		broker = strings.Split(broker, ":")[0]
	}
	for _, e := range config.Entries {
		for _, k := range e.Kafka {
			if strings.Contains(k, broker) {
				return e
			}
		}
	}
	return Entry{}
}

func returnConfig(config []byte) Config {
	conf := Config{}
	err := yaml.Unmarshal(config, &conf)
	if err != nil {
		log.Fatalf("Error returning config: %v\n", err)
	}
	return conf
}

func printConfig(path string) {
	fmt.Printf("%s", readConfig(path))
}

func printConfigSummary(path string) {
	config := returnConfig(readConfig(path))
	fmt.Printf("\nCURRENT: %v\nAvailable Cluster Entries:\n", config.Current)
	for _, e := range config.Entries {
		fmt.Printf("  Name: %v\n", e.Name)
	}
	fmt.Println()
}

func changeCurrent(name, configPath string) {
	config := returnConfig(readConfig(configPath))
	for _, e := range config.Entries {
		if e.Name == name {
			config.Current = name
			y, err := yaml.Marshal(config)
			if err != nil {
				log.Fatalf("Error changing config: %v\n", err)
			}
			writeConfig(configPath, y)
			return
		}
	}
	log.Fatalf("Error: no entry for %v found.\n", name)
}

func readConfig(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}
	return file
}

func writeConfig(path string, config []byte) {
	err := ioutil.WriteFile(path, config, 0644)
	if err != nil {
		log.Fatalf("Error writing config: %v\n", err)
	}
}

// alternate from below
func removeFromConfig(name string, config *Config) {
	for i := len(config.Entries) - 1; i >= 0; i-- {
		if config.Entries[i].Name != name {
			config.Entries = append(config.Entries[:i], config.Entries[i+1:]...)
		}
	}
}

func removeEntry(name string, config Config) []byte {
	tmp := config.Entries[:0]
	for _, c := range config.Entries {
		if c.Name != name {
			tmp = append(tmp, c)
		}
	}
	config.Entries = tmp
	y, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Error removing entry from config: %v\n", err)
	}
	return y
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func generateSampleConfig(path string) {
	if fileExists(path) {
		log.Fatalf("Error: Existing Config Found: %v\n", path)
	} else {
		writeConfig(path, sampleConfigBytes())
	}
}

func sampleConfigBytes() []byte {
	return []byte(sampleConfig())
}

func sampleConfig() string {
	return `current: testCluster1
entries:
- name: testCluster1
  kafka:
  - brokerHost1:9092
  - brokerHost2:9092
  burrow:
  - http://burrow1:3000
  - http://burrow2:3000
  zookeeper:
  - http://zk1:2181
  - http://zk2:2181
- name: testCluster2
  kafka:
  - brokerHost1:9092
  - brokerHost2:9092
  burrow:
  - http://burrow1:3000
  - http://burrow2:3000
  zookeeper:
  - http://zk1:2181
  - http://zk2:2181
`
}
