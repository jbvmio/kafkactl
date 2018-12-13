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
	"strings"

	"github.com/jbvmio/kafkactl"
)

func getMSG(topic string, partition int32, offset int64) *kafkactl.Message {
	msg, err := client.ConsumeOffsetMsg(topic, partition, offset)
	if err != nil {
		closeFatal("Error retrieving message: %v\n", err)
	}
	return msg
}

func getMSGByTime(topic string, partition int32, datetime string) *kafkactl.Message {
	msg, err := client.OffsetMsgByTime(topic, partition, datetime)
	if err != nil {
		if strings.Contains(err.Error(), "parsing time") && strings.Contains(err.Error(), "cannot parse") {
			errMsg := fmt.Sprintf(`datetime parse error: format should be in the form "mm/dd/YYYY HH:MM:SS.000".`)
			closeFatal("Error retrieving message:\n  %v", errMsg)
		}
		closeFatal("Error retrieving message: %v\n", err)
	}
	return msg
}
