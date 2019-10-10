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
	"sort"

	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/kafkactl/cli/x/out"

	kafkactl "github.com/jbvmio/kafka"
)

// APIVersion describes an API Version Key and its Max Version.
type APIVersion struct {
	Name       string
	Key        int16
	MaxVersion int16
}

func GetAPIVersions() []APIVersion {
	var apiVersions []APIVersion
	apiVers, err := client.GetAPIVersions()
	handleC("Error: %v", err)
	for k := range apiVers {
		api := APIVersion{
			Name:       kafkactl.APIDescriptions[k],
			Key:        k,
			MaxVersion: apiVers[k],
		}
		apiVersions = append(apiVersions, api)
	}
	sort.Slice(apiVersions, func(i, j int) bool {
		return apiVersions[i].Key < apiVersions[j].Key
	})
	return apiVersions
}

func findKafkaVersion(context *cx.Context) string {
	bootStrap, err := kafkactl.ReturnFirstValid(context.Brokers...)
	if err != nil {
		out.Failf("Error connecting to Kafka: %v", err)
	}
	apiVer, err := kafkactl.BrokerAPIVersions(bootStrap)
	if err != nil {
		if verbose {
			kafkactl.Warnf("%v", err)
		}
	}
	return getKafkaVersion(apiVer)
}

func getKafkaVersion(apiKeys map[int16]int16) string {
	switch {
	case apiKeys[kafkactl.APIKeyIncrementalAlterConfigs] == 0:
		return "2.3.0"
	case apiKeys[kafkactl.APIKeyOffsetForLeaderEpoch] >= 2:
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
