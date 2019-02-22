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

/*
This is WiP* and needs to be completely redone.
Porting over from an existing version of a burrow cli project for now, modify with kafkactl/burrow package later
*/

package burrow

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/gizak/termui"
	"github.com/tidwall/gjson"
	//termui "gopkg.in/gizak/termui.v2"
)

type burrowDash struct {
	Consumer string
	Topic    string
	Total    int64
	Metrics  burrowBundle
}

type burrowBundle struct {
	Labels []string
	Values []int64
}

func LaunchBurrowMonitor(flags BurrowFlags) {
	if flags.Topic == "" || flags.Group == "" {
		out.Failf("missing topic or group.")
	}
	cg := flags.Group
	top := flags.Topic
	exact = true
	u := getURLMatch([]string{cg})
	if len(u) < 1 {
		out.Failf("Group not found: %v\n  Confirm you are using the correct endpoint: %s\n", cg)
	}
	if err := termui.Init(); err != nil {
		out.Failf("Error launching monitor: %v\n", err)
	}
	defer termui.Close()
	par0 := termui.NewPar(string(" [Group]" + cg + " > [Topic]" + top + " > Partition Lag "))
	par0.Height = 1
	par0.Border = false
	par3 := termui.NewPar(string("                                           q to quit "))
	par3.Height = 1
	par3.Border = false
	bc := termui.NewBarChart()
	bc.Height = 15
	bc.BarWidth = 7
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorWhite

	par1 := termui.NewPar(string(" [ Group ] " + cg + " > [ Topic ] " + top + " > Total Lag "))
	par1.Height = 1
	par1.Border = false
	lc0 := termui.NewLineChart()
	lc0.Height = 20
	lc0.Mode = "dot"
	lc0.AxesColor = termui.ColorWhite
	lc0.LineColor = termui.ColorGreen | termui.AttrBold

	par2 := termui.NewList()
	par2.Height = 20
	par2.BorderLabel = string("| Seconds | Total Lag | (+/-) |")

	var lc0Data []float64
	var lc0Labels []string
	var lc0Count int
	var listItems []string
	var ongoing bool

	generateGraph := func(cg, top string, data []int, lc0Labels, labels, list []string, total int64, bc *termui.BarChart, lc0 *termui.LineChart, lc0Data []float64) {
		bc.Data = data
		bc.DataLabels = labels
		lc0.Data = lc0Data
		lc0.DataLabels = lc0Labels
		par1.Text = (string(" [Group]" + cg + " > [Topic]" + top + " > Total Lag: " + strconv.FormatInt(total, 10)))
		par2.Items = list
		termui.Render(termui.Body)
	}

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(9, 0, par0),
			termui.NewCol(9, 0, par3)),
		termui.NewRow(
			termui.NewCol(12, 0, bc)),
		termui.NewRow(
			termui.NewCol(9, 0, par1)),
		termui.NewRow(
			termui.NewCol(9, 0, lc0),
			termui.NewCol(3, 0, par2)))

	termui.Body.Align()

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd", func(termui.Event) {
		lc0Data = lc0Data[(len(lc0Data) - 1):]
		lc0Labels = lc0Labels[(len(lc0Labels) - 1):]
		lc0Labels[0] = "0"
		ongoing = false
	})
	termui.Handle("/timer/1s", func(termui.Event) {
		m := getTopicMatch(top, u)
		dash := createDash(m[0], cg, top)
		total := dash.Total
		var data []int
		for _, x := range dash.Metrics.Values {
			data = append(data, int(x))
		}
		lc0Count++
		lc0Data = append(lc0Data, float64(dash.Total))
		lc0Labels = append(lc0Labels, strconv.Itoa(lc0Count))
		labels := dash.Metrics.Labels
		if !ongoing {
			if lc0Count%100 == 99 {
				ongoing = true
			}
		}
		if ongoing {
			lc0Data = lc0Data[1:]
			lc0Labels = lc0Labels[1:]
		}
		var before float64
		var now float64
		if len(lc0Data) > 2 {
			before = lc0Data[len(lc0Data)-2]
			now = lc0Data[len(lc0Data)-1]
		}
		var progress string
		diff := now - before
		if diff == 0 {
			progress = ""
		}
		if diff < 0 {
			progress = string("Decreasing (" + strconv.Itoa(int(diff)) + ")")
		}
		if diff > 0 {
			progress = string("Increasing (+" + strconv.Itoa(int(diff)) + ")")
		}
		listItems = append(listItems, string(" ["+strconv.Itoa(lc0Count)+"]    "+strconv.Itoa(int(lc0Data[len(lc0Data)-1]))+"    "+progress))

		if len(listItems) > 15 {
			listItems = listItems[1:]
		}

		generateGraph(cg, top, data, lc0Labels, labels, listItems, total, bc, lc0, lc0Data)
	})
	termui.Loop()
}

func createDash(g gjson.Result, cg, top string) burrowDash {
	if !g.Get("0.topic").Exists() {
		out.Failf("Topic not found: %v\n  Confirm you are using the correct endpoint: %s\n", top)
	}
	bd := burrowDash{}
	bd.Consumer = cg
	bd.Topic = top
	bb := burrowBundle{}
	var tl int64
	elements := len(g.Array())
	for e := 0; e < elements; e++ {
		bb.Labels = append(bb.Labels, string("P"+strconv.FormatInt((g.Get(string(strconv.Itoa(e)+".partition")).Int()), 10)))
		v := g.Get(string(strconv.Itoa(e) + ".current_lag")).Int()
		tl += v
		bb.Values = append(bb.Values, v)
	}
	bd.Total = tl
	bd.Metrics = bb
	return bd
}

func getTopicMatch(topic string, paths []string) []gjson.Result {
	var matches []gjson.Result
	t := string("status.partitions.#[topic==" + topic + "]#")
	for _, p := range paths {
		r := gjson.GetBytes(curlIt(string(p+"/lag")), t)
		matches = append(matches, r)
	}
	return matches
}

func getURLMatch(terms []string) (consumerPaths []string) {
	var cluPaths []string
	allBaseUrls := burClient.GetClusterURLs()
	for _, base := range allBaseUrls {
		clu := gjson.GetBytes(curlIt(base), "clusters").Array()
		for _, i := range clu {
			path := string(base + "/" + i.String())
			cluPaths = append(cluPaths, path)
		}
		for _, i := range cluPaths {
			con := gjson.GetBytes(curlIt(string(i+"/consumer")), "consumers").Array()
			for _, x := range con {
				for _, y := range terms {
					if x.String() == y {
						path := string(i + "/consumer/" + x.String())
						consumerPaths = append(consumerPaths, path)
					}
				}
			}
		}
	}
	return
}

func curlIt(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		out.Failf("Error: %v\n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		out.Failf("Error: %v\n", err)
	}
	return body
}
