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
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var (
	groupList     []string
	config        *sarama.Config
	client        sarama.Client
	clientID      string
	searchByTopic bool
	searchByGroup bool
	isIdle        bool
)

var groupCmd = &cobra.Command{
	Use:     "group",
	Short:   "Search and Retrieve Group details",
	Long:    `Example kafkactl group group1 group2`,
	Aliases: []string{"groups"},
	Run: func(cmd *cobra.Command, args []string) {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(43)
		s := spinner.New(spinner.CharSets[r], 100*time.Millisecond)
		s.Prefix = fmt.Sprintf("Searching Groups by Topic - Please Wait ... ")
		if targetTopic != "" {
			searchByTopic = true
		}
		if verbose {
			sarama.Logger = log.New(os.Stdout, "[kafkactl] ", log.LstdFlags)
		}
		config = gimmeConf()
		client = gimmeClient(bootStrap, config)
		defer func() {
			if err := client.Close(); err != nil {
				log.Fatalf("Error closing client: %v\n", err)
			}
		}()
		allBrokers := client.Brokers()
		defer func() {
			for _, b := range allBrokers {
				if ok, _ := b.Connected(); ok {
					if err := b.Close(); err != nil {
						log.Fatalf("Error closing broker, %v: %v\n", b.Addr(), err)
					}
				}
			}
		}()
		if len(args) == 0 {
			if !searchByTopic && clientID == "" {
				gl := getAllGroups(client)
				sort.Slice(gl, func(i, j int) bool {
					return gl[i].Group < gl[j].Group
				})
				formatOutput(gl)
				return
			}
			groupList = getGroupList()
		} else {
			for _, i := range getGroupList() {
				for _, s := range args {
					if exact {
						if i == s {
							groupList = append(groupList, i)
						}
					} else {
						if strings.Contains(i, s) {
							groupList = append(groupList, i)
						}
					}
				}
			}
			searchByGroup = true
		}
		if len(groupList) == 0 {
			if searchByTopic {
				fmt.Printf("No results found for groups: %v and topic: %v\n", args, targetTopic)
				return
			}
			fmt.Printf("No results found for groups: %v\n", args)
		} else {
			if (searchByTopic || clientID != "") && len(groupList) > 150 {
				if clientID != "" {
					s.Prefix = fmt.Sprintf("Filtering Groups by ClientID - Please Wait ... ")
				}
				if searchByTopic && clientID != "" {
					s.Prefix = fmt.Sprintf("Filtering Groups by Topic and ClientID - Please Wait ... ")
				}
				s.Start()
			}
			grp := describeGroups(client, groupList)
			if len(grp) == 0 {
				fmt.Printf("\nNo Results Found: %v %v\n\n", targetTopic, groupList)
				if isIdle {
					fmt.Println("  It is possible these groups are just idle or inactive.\n")
				}
				return
			}
			s.Stop()
			if clientID != "" {
				var tmp []groupMetadata
				for _, g := range grp {
					if exact {
						if g.ClientID == clientID {
							tmp = append(tmp, g)
						}
					} else {
						if strings.Contains(g.ClientID, clientID) {
							tmp = append(tmp, g)
						}
					}
				}
				grp = tmp
				if len(grp) == 0 {
					if searchByTopic {
						fmt.Printf("\nNo results found for clientID: %v and topic: %v\n\n", clientID, targetTopic)
						return
					}
					fmt.Printf("\nNo results found for clientID: %v\n\n", clientID)
					return
				}
			}
			formatOutput(grp)
		}
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	groupCmd.Flags().StringVarP(&clientID, "clientid", "i", "", "Find groups by ClientID")
}
