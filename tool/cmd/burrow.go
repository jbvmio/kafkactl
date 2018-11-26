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
	"strings"

	"github.com/jbvmio/burrow"
)

var burClient *burrow.Client

func launchBurrowClient() {
	if len(burrowEPs) < 1 {
		log.Fatalf("No Burrow Endpoints Defined.\n")
	}
	burClient, errd = burrow.NewBurrowClient(burrowEPs)
	if errd != nil {
		log.Fatalf("Error initializing burrow client: %v\n", errd)
	}
}

func searchBurrowConsumers(consumers ...string) []burrow.Partition {
	var burrowConsumers []string
	cl, err := burClient.GetConsumerList()
	if err != nil {
		log.Fatalf("Error obtaining consumer list: %v\n", err)
	}
	for _, consumer := range consumers {
		for _, c := range cl {
			if exact {
				if c == consumer {
					burrowConsumers = append(burrowConsumers, c)
				}
			} else {
				if strings.Contains(c, consumer) {
					burrowConsumers = append(burrowConsumers, c)
				}
			}
		}
	}
	conParts, err := burClient.GetConsumerPartitions(burrowConsumers...)
	if err != nil {
		log.Fatalf("Error obtaining consumer data: %v\n", err)
	}
	return conParts
}
