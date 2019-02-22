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

package burrow

import (
	"strings"

	"github.com/jbvmio/kafkactl/cli/cx"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/jbvmio/burrow"
)

type BurrowFlags struct {
	ShowID  bool
	Monitor bool
	Exact   bool
	ErrOnly bool
	Group   string
	Topic   string
	Context string
}

var (
	burClient *burrow.Client
	errd      error
	exact     bool
)

func LaunchBurrowClient(context *cx.Context, flags BurrowFlags) {
	var burrowEP []string
	switch {
	case len(context.Burrow) < 1:
		out.Failf("No Burrow Endpoints Defined.")
	default:
		burrowEP = context.Burrow
	}
	burClient, errd = burrow.NewBurrowClient(burrowEP)
	if errd != nil {
		out.Failf("Error initializing burrow client: %v", errd)
	}
	exact = flags.Exact
}

func SearchBurrowConsumers(flags BurrowFlags, consumers ...string) []burrow.Partition {
	var burrowConsumers []string
	var err error
	switch {
	case len(consumers) > 0:
		cl, err := burClient.GetConsumerList()
		if err != nil {
			out.Failf("Error obtaining consumer list: %v\n", err)
		}
		switch {
		case exact:
			for _, consumer := range consumers {
				for _, c := range cl {
					if c == consumer {
						burrowConsumers = append(burrowConsumers, c)
					}
				}
			}
		default:
			for _, consumer := range consumers {
				for _, c := range cl {
					if strings.Contains(c, consumer) {
						burrowConsumers = append(burrowConsumers, c)
					}
				}
			}
		}
	default:
		burrowConsumers, err = burClient.GetConsumerList()
		if err != nil {
			out.Failf("Error obtaining consumer list: %v\n", err)
		}
	}
	conParts, err := burClient.GetConsumerPartitions(burrowConsumers...)
	if err != nil {
		out.Failf("Error obtaining consumer data: %v\n", err)
	}
	if flags.Topic != "" {
		conParts = searchBurrowTopics(conParts, flags.Topic)
	}
	if flags.ErrOnly {
		conParts = filterERRs(conParts)
	}
	return conParts
}

func searchBurrowTopics(parts []burrow.Partition, topic string) []burrow.Partition {
	var topicCons []burrow.Partition
	switch {
	case exact:
		for _, p := range parts {
			if p.Topic == topic {
				topicCons = append(topicCons, p)
			}
		}
	default:
		for _, p := range parts {
			if strings.Contains(p.Topic, topic) {
				topicCons = append(topicCons, p)
			}
		}
	}
	return topicCons
}

func filterERRs(parts []burrow.Partition) []burrow.Partition {
	var topicCons []burrow.Partition
	for _, p := range parts {
		if p.PStatusCode > 1 {
			topicCons = append(topicCons, p)
		}
	}
	return topicCons
}
