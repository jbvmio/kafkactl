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

	"github.com/spf13/cobra"
)

const (
	//majorVer = `1.`
	//minorVer = `0.`
	//patchVer = `19`
	contact = `jbvm.io`
)

// buildTime - revision := ~Year
var (
	showLatest bool
	fullVer    string

	release     string
	nextRelease string
	latestMajor string
	latestMinor string
	latestPatch string
	revision    string
	buildTime   string
	commitHash  string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print kafkactl version and exit",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fullVer = latestMajor + `.` + latestMinor + `.` + latestPatch
		if release != "true" {
			fullVer = fullVer + `+` + revision
		} else {
			fullVer = nextRelease
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kafkactl   : %s\n", contact)
		fmt.Printf("Version    : %s\n", fullVer)
		fmt.Printf("Build      : %s\n", buildTime)
		fmt.Printf("Revision   : %s\n", revision)
		fmt.Printf("Commit     : %s\n", commitHash)
		if showLatest {
			switch {
			case latestMajor == "" || latestMinor == "" || latestPatch == "":
				fmt.Printf("Latest*    : N/A\n")
				fmt.Printf("NextRel*   : %s\n", nextRelease)
			default:
				fmt.Printf("Latest*    : %s\n", string(latestMajor+`.`+latestMinor+`.`+latestPatch))
				fmt.Printf("NextRel*   : %s\n", nextRelease)
			}
		}
	},
}

func init() {
	versionCmd.Flags().BoolVar(&showLatest, "latest", false, "Show Latest Release.")
	versionCmd.Flags().MarkHidden("latest")

	rootCmd.AddCommand(versionCmd)
}
