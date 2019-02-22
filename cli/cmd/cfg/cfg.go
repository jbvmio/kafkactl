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

package cfg

import (
	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const configVersion = 1

// Config holds all values for a given context.
type Config struct {
	Contexts       map[string]cx.Context `yaml:"contexts"`
	CurrentContext string                `yaml:"current-context"`
	ConfigVersion  int                   `yaml:"config-version"`
}

// CXFlags is used for ad-hoc Contexts.
type CXFlags struct {
	Broker    string
	Zookeeper string
	Burrow    string
}

func GetConfig() *Config {
	var config Config
	viper.Unmarshal(&config)
	config.CurrentContext = viper.GetString("current-context")
	config.ConfigVersion = viper.GetInt("config-version")
	return &config
}

func GetContextList() map[string][]string {
	config := GetConfig()
	contexts := make(map[string][]string, len(config.Contexts))
	for k := range config.Contexts {
		if k == config.CurrentContext {
			contexts["contexts"] = append(contexts["contexts"], string(k+" [current-context]"))
		} else {
			contexts["contexts"] = append(contexts["contexts"], k)
		}
	}
	return contexts
}

// GetContext returns the configuration for the given context, or the current context if none is specified.
func GetContext(context ...string) *cx.Context {
	switch true {
	case len(context) > 1:
		out.Failf("Error: too many contexts specified, only 1 allowed")
	case len(context) < 1:
		return getCurrentCtx()
	case context[0] == "":
		return getCurrentCtx()
	}
	config := GetConfig()
	ctx := config.Contexts[context[0]]
	if ctx.Name == "" {
		out.Failf("Error: no context named %v", context[0])
	}
	return &ctx
}

func getCurrentCtx() *cx.Context {
	current := viper.GetString("current-context")
	config := GetConfig()
	ctx := config.Contexts[current]
	if ctx.Name == "" {
		converted := testReplaceOldConfig()
		if !converted {
			out.Failf("Error: invalid config or context")
		}
		viper.ReadInConfig()
		current = viper.GetString("current-context")
		config = GetConfig()
		ctx = config.Contexts[current]
		if ctx.Name == "" {
			out.Failf("Error: still invalid config or context")
		}
	}
	return &ctx
}

// AdhocContext return an adhoc Context built from Flag parameters.
func AdhocContext(cxFlags CXFlags) *cx.Context {
	return &cx.Context{
		Brokers:   []string{cxFlags.Broker},
		Burrow:    []string{cxFlags.Burrow},
		Zookeeper: []string{cxFlags.Zookeeper},
	}
}

// GenConfig generates and outputs sample config.
func GenSample() {
	var config Config
	config.Contexts = make(map[string]cx.Context)
	for i := 1; i < 4; i++ {
		thatOne := cast.ToString(i)
		ctx := cx.Context{
			Name: string("prod-atl0" + thatOne),
		}
		for x := 1; x < 4; x++ {
			thisOne := cast.ToString(x)
			ctx.Brokers = append(ctx.Brokers, string("broker0"+thisOne+":9092"))
			ctx.Burrow = append(ctx.Burrow, string("burrow0"+thisOne+":8080"))
			ctx.Zookeeper = append(ctx.Zookeeper, string("zkhost0"+thisOne+":2181"))
		}
		config.Contexts[ctx.Name] = ctx
	}
	config.CurrentContext = "prod-atl01"
	config.ConfigVersion = configVersion
	out.Marshal(config, "yaml")
}
