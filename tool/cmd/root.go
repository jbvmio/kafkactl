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
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	bootStrap string
	bsport    string
	exact     bool
	verbose   bool
	useFast   bool
	meta      bool

	targetTopic     string
	targetGroup     string
	targetPartition int32
)

var wg sync.WaitGroup
var ctx = context.Background()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafkactl",
	Short: "kafkactl: Kafka Management Tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		/*
			if bootStrap == "" {
				if fileExists(configLocation) {
					bootStrap = getCurrentTarget()
				}
			}
		*/
		if !strings.Contains(bootStrap, ":") {
			bootStrap = net.JoinHostPort(bootStrap, bsport)
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
