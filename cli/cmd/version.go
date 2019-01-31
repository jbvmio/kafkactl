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
	yTime   = 1546300800
	maint   = "19"
	version = "1.0."
)

var (
	revision   string
	buildTime  string
	commitHash string
	fullVer    string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print kafkactl version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version   : %s\n", fullVer)
		fmt.Printf("Build     : %s\n", buildTime)
		fmt.Printf("Commit    : %s\n", commitHash)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
