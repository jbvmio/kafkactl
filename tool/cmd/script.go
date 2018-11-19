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
	"strings"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/spf13/cobra"
)

var (
	listScripts bool
	box         packr.Box
)

var scriptCmd = &cobra.Command{
	Use:     "script",
	Aliases: []string{"scripts"},
	Short:   "script testing WiP*",
	Long:    `WiP*`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		box = packr.NewBox("./scripts")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalln("Running Kafka Scripts is Currently WiP*")
		timeFormat := "01/02/2006 15:04:05.000"
		targetTime := "11/02/2018 8:01:00.0"
		testTargetTime := strings.Split(targetTime, ".")
		fmt.Println(len(testTargetTime), testTargetTime)
		if len(testTargetTime) > 1 {
			addZeros := 3 - (len([]rune(testTargetTime[1])))
			if addZeros > 0 {
				for i := 0; i < addZeros; i++ {
					targetTime += "0"
				}
			}
		}
		//fmt.Println(t.Format("2006-01-02 15:04:05"))
		time, err := time.Parse(timeFormat, targetTime)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		timeSeconds := time.Unix()
		timeMilli := timeSeconds * 1000
		//fmt.Println(t.Format("01/02/2006 15:04:05.000"))
		fmt.Println(time)
		fmt.Println(timeSeconds, timeMilli)
		return
		if listScripts || len(args) < 1 {
			for _, s := range box.List() {
				fmt.Println(s)
			}
			return
		}
		targetScript := args[0]
		args = args[1:]
		b, err := box.Find(targetScript)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		scriptFile := fmt.Sprintf("%s", b)
		runScript(scriptFile, args...)
		//fmt.Printf("%s", b)
		return
	},
}

func init() {
	adminCmd.AddCommand(scriptCmd)
	scriptCmd.Flags().BoolVarP(&listScripts, "listscripts", "l", false, "List Included Scripts")
	//scriptCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")

}
