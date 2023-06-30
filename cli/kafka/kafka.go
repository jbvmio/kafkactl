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

	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/kafkactl/cli/x/out"
	metrics "github.com/rcrowley/go-metrics"

	"github.com/Shopify/sarama"
	kafkactl "github.com/jbvmio/kafkactl/kafka"
)

var (
	// exact specifies exact match querying, otherwise querying returns wildcard matches.
	exact bool
	// verbose enables additional details to print.
	verbose bool
	// targetContext stores the targeted context.
	targetContext *cx.Context
	// FORCE bypasses any configured checks.
	FORCE bool
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
	conf          *sarama.Config
	clientVer     sarama.KafkaVersion
	clientTimeout = (time.Second * 5)
	clientRetries = 1
	clientID      = `kafkactl`
)

func LaunchClient(context *cx.Context, flags ClientFlags) {
	exact = flags.Exact
	verbose = flags.Verbose
	targetContext = context
	if verbose {
		kafkactl.Logger()
	}
	conf = kafkactl.GetConf()
	conf.ClientID = clientID
	conf.Net.DialTimeout = clientTimeout
	conf.Net.ReadTimeout = clientTimeout
	conf.Net.WriteTimeout = clientTimeout
	conf.Metadata.Retry.Max = clientRetries
	conf.MetricRegistry = metrics.NewRegistry()
	ssl := context.Ssl
	if ssl != nil {
		tlsConfig, err := cx.SetupCerts(ssl.TLSCert, ssl.TLSCA, ssl.TLSKey, ssl.Insecure)
		if err != nil {
			out.Failf("Error setting up certs: %v", err)
		}
		conf.Net.TLS.Enable = true
		conf.Net.TLS.Config = tlsConfig
	}
	sasl := context.Sasl
	if sasl != nil {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = context.Sasl.Username
		conf.Net.SASL.Password = context.Sasl.Password
	}
	switch {
	case flags.Version != "":
		context.ClientVersion = flags.Version
	case context.ClientVersion == "":
		context.ClientVersion = findKafkaVersion(conf, context)
	}
	conf.Version, errd = kafkactl.MatchKafkaVersion(context.ClientVersion)
	if errd != nil {
		if verbose {
			kafkactl.Warnf("%v Defaulting to %v", errd, kafkactl.RecKafkaVersion)
		}
		conf.Version = kafkactl.RecKafkaVersion
	}
	clientVer = conf.Version
	client, errd = kafkactl.NewCustomClient(conf, context.Brokers...)
	if errd != nil {
		out.Failf("Error: %v", errd)
	}
	Alive = true
}

// MetricR returns the kafkactl metrics registry.
func MetricR() *metrics.Registry {
	return &conf.MetricRegistry
}

// Client returns the kafkactl client.
func Client() *kafkactl.KClient {
	return client
}

// ClientVersion returns the kafkactl client.
func ClientVersion() sarama.KafkaVersion {
	return clientVer
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
