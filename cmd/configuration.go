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
	"os"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// KafkactlConfig holds the kafkactl configuration
type KafkactlConfig struct {
	Targets map[string]string `yaml:"targets"`
	Current map[string]string `yaml:"current"`
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

func printConfig() {
	conf := KafkactlConfig{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("error reading config:", err)
	}
	yaml, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println("error reading config:", err)
	}
	fmt.Printf("\n%v\n\n", string(yaml))
}

func returnConfig() KafkactlConfig {
	conf := KafkactlConfig{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("error reading config:", err)
	}
	return conf
}

func useTarget(name string) {
	get := viper.Get(string("targets." + name))
	if get == "" || get == nil {
		fmt.Printf("\nEntry for %v not found\n\n", name)
		os.Exit(1)
	}
	conf := KafkactlConfig{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("error reading config:", err)
	}
	new := make(map[string]string, 1)
	new[name] = get.(string)
	conf.Current = new
	yaml, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println("error reading config:", err)
	}
	writeConfig(yaml)
	fmt.Printf("\nTarget set >\n  %v: %v\n\n", name, new[name])
}

func addTarget(name, addr string) {
	conf := KafkactlConfig{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("error reading config:", err)
	}
	conf.Targets[name] = addr
	yaml, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println("error writing config:", err)
	}
	writeConfig(yaml)
	fmt.Printf("\nTarget added >\n  %v: %v\n\n", name, conf.Targets[name])
}

func removeTarget(name string) {
	get := viper.Get(string("targets." + name))
	if get == "" || get == nil {
		fmt.Printf("\nEntry for %v not found\n\n", name)
		os.Exit(1)
	}
	conf := KafkactlConfig{}
	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("error reading config:", err)
	}
	delete(conf.Targets, name)
	yaml, err := yaml.Marshal(conf)
	if err != nil {
		fmt.Println("error writing config:", err)
	}
	writeConfig(yaml)
	fmt.Printf("\nTarget removed >\n  %v: %v\n\n", name, get.(string))
}

func writeConfig(data []byte) error {
	dest := configLocation
	err := ioutil.WriteFile(dest, data, 0644)
	return err
}

func getCurrentTarget() string {
	var broker string
	get := viper.Get("current")
	if get == "" || get == nil {
		fmt.Printf("\nEntries not found\n\n")
		os.Exit(1)
	}
	for _, v := range get.(map[string]interface{}) {
		broker = v.(string)
		return broker
	}
	return broker
}
