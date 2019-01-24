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
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	strParts string
	tParts   []int32
)

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail a Topic",
	Long:  `Example: kafkactl tail 5 -t myTopic (Default tail start from newest: 1)`,
	Run: func(cmd *cobra.Command, args []string) {
		if targetTopic == "" {
			closeFatal("specify a topic, eg. --topic")
		}
		tailTarget := validateTailArgs(args)
		if strParts == "" {
			tParts = []int32{}
		} else {
			parts := strings.Split(strParts, ",")
			for _, p := range parts {
				tParts = append(tParts, cast.ToInt32(p))
			}
			validateParts(tParts)
		}
		tailTopic(targetTopic, tailTarget, tParts...)
		return
	},
}

func init() {
	rootCmd.AddCommand(tailCmd)
	tailCmd.Flags().StringVarP(&strParts, "partitions", "p", "", `Comma Separated (eg: "0,1,7,9") Partitions to Tail (Default: All)`)
}
