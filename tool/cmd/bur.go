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
	"github.com/spf13/cobra"
)

var (
	showClientID   bool
	targetBurrowEP string
)

var burrowCmd = &cobra.Command{
	Use:     "burrow",
	Aliases: []string{"b", "bur"},
	Short:   "Query Burrow",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		nonMainCMD = true
		if cmd.Flags().Changed("ep") {
			burrowEPs = []string{targetBurrowEP}
		}
		launchBurrowClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			args = []string{""}
		}
		printOutput(searchBurrowConsumers(args...))
		return
	},
}

func init() {
	rootCmd.AddCommand(burrowCmd)
	burrowCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	burrowCmd.Flags().BoolVarP(&showClientID, "clientid", "i", false, "Toggle/Show ClientIDs instead of Group Names")
	burrowCmd.Flags().StringVarP(&targetBurrowEP, "ep", "e", "", "Specify a targeted Burrow EndPoint (eg. http://localhost:8000)")
}
