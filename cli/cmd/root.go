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

	"github.com/jbvmio/kafkactl/cli/cmd/cfg"
	"github.com/jbvmio/kafkactl/cli/cmd/describe"
	"github.com/jbvmio/kafkactl/cli/cmd/get"
	"github.com/jbvmio/kafkactl/cli/cmd/msg"
	"github.com/jbvmio/kafkactl/cli/kafka"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	targetContext string
	//verbose       bool
)

var kafkaFlags kafka.ClientFlags

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafkactl",
	Short: "kafkactl: Kafka Management Tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		kafka.LaunchClient(cfg.GetContext(kafkaFlags.Context), kafkaFlags)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		kafka.CloseClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		kafka.ClusterDetails(kafka.Client())
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "cfg", "", "config file (default is $HOME/.2kafkactl.yaml)")
	rootCmd.PersistentFlags().StringVar(&kafkaFlags.Context, "context", "", "Specify a context.")
	rootCmd.PersistentFlags().StringVar(&kafkaFlags.Version, "version", "", "Specify a client version.")
	rootCmd.PersistentFlags().BoolVarP(&kafkaFlags.Verbose, "verbose", "v", false, "Display additional info or errors.")
	rootCmd.PersistentFlags().BoolVarP(&kafkaFlags.Exact, "exact", "x", false, "Find exact matches.")

	rootCmd.AddCommand(cfg.CmdConfig)
	rootCmd.AddCommand(get.CmdGet)
	rootCmd.AddCommand(describe.CmdDescribe)
	rootCmd.AddCommand(msg.CmdLogs)

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".2kafkactl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
