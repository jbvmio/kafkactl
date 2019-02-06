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
	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/kafkactl/cli/x/out"
)

func findKafkaVersion(context *cx.Context) string {
	bootStrap, err := kafkactl.ReturnFirstValid(context.Brokers...)
	if err != nil {
		out.Failf("Error connecting to Kafka: %v", err)
	}
	apiVer, err := kafkactl.BrokerAPIVersions(bootStrap)
	if err != nil {
		out.Infof("ERR: %v", err)
	}
	return getKafkaVersion(apiVer)
}

func getKafkaVersion(apiKeys map[int16]int16) string {
	switch true {
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
