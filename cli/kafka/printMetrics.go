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

package kafka

import (
	"fmt"

	"github.com/fatih/color"
	kafkactl "github.com/jbvmio/kafkactl/kafka"
	"github.com/rodaine/table"
)

func PrintMetricCollection(MC kafkactl.MetricCollection) {
	PrintMetrics(MC.Meters)
	PrintMetrics(MC.Histograms)
}

func PrintMetrics(i interface{}) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	var tbl table.Table
	switch i := i.(type) {
	case []*kafkactl.RawMetric:
		tbl = table.New("MEASUREMENT", "VALUES", "TYPE")
		for _, v := range i {
			tbl.AddRow(v.Measurement, v.Values, v.Type)
		}
	case []kafkactl.MeterMetric:
		tbl = table.New("MEASUREMENT", "COUNT", "MEAN", "1MINUTE", "5MINUTE", "15MINUTE", "TYPE")
		for _, v := range i {
			tbl.AddRow(v.Measurement, v.Count, v.MeanRate, v.OneRate, v.FiveRate, v.FifteenRate, v.Type)
		}
	case []kafkactl.HistoMetric:
		tbl = table.New("MEASUREMENT", "COUNT", "MIN", "MAX", "75%", "99%", "TYPE")
		for _, v := range i {
			tbl.AddRow(v.Measurement, v.Count, v.Min, v.Max, v.SeventyFive, v.NinetyNine, v.Type)
		}
	}
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.Print()
	fmt.Println()
}
