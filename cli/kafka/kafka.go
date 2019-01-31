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
	"time"

	"github.com/Shopify/sarama"
	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/kafkactl/cli/x/out"
)

var (
	// exact specifies exact match querying, otherwise querying returns wildcard matches.
	exact bool
	// verbose enables additional details to print.
	verbose bool
)

type ClientFlags struct {
	Exact   bool
	Verbose bool
	Context string
	Version string
}

// client variables
var (
	Alive         bool
	client        *kafkactl.KClient
	errd          error
	clientVer     sarama.KafkaVersion
	clientTimeout = (time.Second * 5)
	clientRetries = 1
)

func LaunchClient(context *cx.Context, flags ClientFlags) {
	exact = flags.Exact
	verbose = flags.Verbose
	if verbose {
		kafkactl.Logger("")
	}
	conf, err := kafkactl.GetConf()
	if err != nil {
		out.Failf("Error: %v", err)
	}
	match := true
	switch match {
	case flags.Version != "":
		context.ClientVersion = flags.Version
	case context.ClientVersion == "":
		context.ClientVersion = findKafkaVersion(context)
	}
	conf.Version, err = kafkactl.MatchKafkaVersion(context.ClientVersion)
	if err != nil {
		out.Warnf("WARN: %v", err)
		conf.Version = kafkactl.MinKafkaVersion
	}
	clientVer = conf.Version
	conf.Net.DialTimeout = clientTimeout
	conf.Net.ReadTimeout = clientTimeout
	conf.Net.WriteTimeout = clientTimeout
	conf.Metadata.Retry.Max = clientRetries
	client, errd = kafkactl.NewCustomClient(conf, context.Brokers...)
	if errd != nil {
		out.Failf("Error: %v", errd)
	}
	Alive = true
}

// Client returns the kafkactl client.
func Client() *kafkactl.KClient {
	return client
}

func CloseClient() {
	if Alive {
		if connected := client.IsConnected(); connected {
			if errd = client.Close(); errd != nil {
				out.Failf("Error closing client: %v\n", errd)
			}
		}
	}
}

func closeFatal(format string, args ...interface{}) {
	CloseClient()
	out.Failf(format, args...)
}

func handleF(format string, err error) {
	if err != nil {
		out.Failf(format, err)
	}
}

func handleW(format string, err error) {
	if err != nil {
		out.Warnf(format, err)
	}
}

func handleC(format string, err error) {
	if err != nil {
		closeFatal(format, err)
	}
}
