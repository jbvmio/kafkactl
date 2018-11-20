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
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile   bool
	bootStrap string
	bsport    string
	exact     bool
	verbose   bool
	useFast   bool
	meta      bool

	targetTopic     string
	targetGroup     string
	targetPartition int32

	kafkaBrokers []string
	burrowEPs    []string
	zkServers    []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafkactl",
	Short: "kafkactl: Kafka Management Tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if configCmd.CalledAs() != "config" {
			validateBootStrap()
			launchClient()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if configCmd.CalledAs() != "config" && !zkCommandInvoked {
			closeClient()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(metaCmd, args)
	},
}

// Execute starts here.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&bootStrap, "broker", "b", "", "Bootstrap Kafka Broker")
	rootCmd.PersistentFlags().StringVarP(&targetTopic, "topic", "t", "", "Specify a Target Topic")
	rootCmd.PersistentFlags().StringVarP(&targetGroup, "group", "g", "", "Specify a Target Group")
	rootCmd.PersistentFlags().StringVar(&bsport, "port", "9092", "Port used for Bootstrap Kafka Broker")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display any additional info and error output.")
	rootCmd.PersistentFlags().BoolVarP(&useFast, "fast", "F", true, "Use go routine methods when available for faster results")
}

func initConfig() {
	if fileExists(configLocation) {
		kafkaBrokers, burrowEPs, zkServers = getEntries(configLocation)
		cfgFile = true
	}
}

/*
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kafkactl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".2kafkactl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
*/
