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
	"log"

	"github.com/spf13/cobra"
)

var (
	genSample           bool
	showConfig          bool
	showConfigFull      bool
	changeCurrentTarget string
	configLocation      = string(homeDir() + "/.2kafkactl.yaml")
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage kafkactl Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if genSample {
			generateSampleConfig(configLocation)
			return
		}
		if showConfig {
			printConfigSummary(configLocation)
			return
		}
		if showConfigFull {
			printConfig(configLocation)
		}
		if cmd.Flags().Changed("use") {
			if changeCurrentTarget == "" {
				log.Fatalf("Enter a config entry name to switch.\n")
			}
			changeCurrent(changeCurrentTarget, configLocation)
			return
		}
		if cmd.Flags().Changed("delete") {
			if changeCurrentTarget == "" {
				log.Fatalf("Enter a config entry name to delete.\n")
			}
			writeConfig(configLocation, removeEntry(changeCurrentTarget, returnConfig(readConfig(configLocation))))
			return
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVar(&genSample, "sample", false, "Generate a Sample Configuration at ~/.2kafkactl (will not overwrite existing if found)")
	configCmd.Flags().BoolVar(&showConfig, "show", false, "Show available config targets")
	configCmd.Flags().BoolVar(&showConfigFull, "show-full", false, "Print current config and exit")
	configCmd.Flags().StringVar(&changeCurrentTarget, "use", "", "Switch Current Target in Config")
	configCmd.Flags().StringVar(&changeCurrentTarget, "delete", "", "Delete an Entry in the Current Config")
}
