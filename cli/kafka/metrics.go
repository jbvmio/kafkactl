package kafka

import (
	"fmt"
	"sort"
	"time"

	kafkactl "github.com/jbvmio/kafka"
)

// MetricFlags defines flags for reporting metrics.
type MetricFlags struct {
	Broker    bool
	Intervals int
	Seconds   int
}

func GetKakfaMetrics(flags MetricFlags) kafkactl.MetricCollection {
	var MC kafkactl.MetricCollection
	switch {
	case flags.Broker:
		MC = GetBrokerMetrics(flags.Intervals, flags.Seconds)
	}
	return MC
}

// GetBrokerMetrics return broker related metrics.
func GetBrokerMetrics(intervals, seconds int) kafkactl.MetricCollection {
	var MC kafkactl.MetricCollection
	var allM []*kafkactl.RawMetric
	client.GetClusterMeta()
	r := *MetricR()
	allMetrics := r.GetAll()
	for m := range allMetrics {
		metric := kafkactl.RawMetric{
			Measurement: m,
			Values:      allMetrics[m],
		}
		allM = append(allM, &metric)
	}
	// Periodically log metrics to stdout:
	//go metrics.Log(*MetricR(), 5*time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))
	for i := 1; i <= intervals; i++ {
		if len(r.GetAll()) > len(allMetrics) {
			var tmpM []*kafkactl.RawMetric
			allMetrics = r.GetAll()
			for m := range allMetrics {
				metric := kafkactl.RawMetric{
					Measurement: m,
					Values:      allMetrics[m],
				}
				tmpM = append(tmpM, &metric)
			}
			allM = tmpM
		}
		cm, err := client.GetClusterMeta()
		if err != nil {
			client.Warnf("encountered error %v\n", err)
		} else {
			if verbose {
				verb := fmt.Sprintf("%v Brokers, %v Groups, %v Topics", cm.BrokerCount(), cm.GroupCount(), cm.TopicCount())
				client.Logf("Completed Interval %v - %v\n", i, verb)
			}
		}
		for _, m := range allM {
			m.Update(MetricR())
		}
		time.Sleep(time.Duration(seconds) * time.Second)
	}
	MC.AddFromRaw(allM...)
	/*
		for _, m := range allM {
			MC.Add(m.ConvertToMetric())
		}
	*/
	sort.Slice(MC.Meters, func(i, j int) bool {
		return MC.Meters[i].Measurement < MC.Meters[j].Measurement
	})
	sort.Slice(MC.Histograms, func(i, j int) bool {
		return MC.Histograms[i].Measurement < MC.Histograms[j].Measurement
	})
	r.UnregisterAll()
	return MC
}

/*
func TestMetrics2() {
	//metricsRegistry = metrics.NewRegistry()
	var allM []RawMetric
	r := *MetricR()

	allMetrics := r.GetAll()
	for m := range allMetrics {

		metric := RawMetric{
			Measurement: m,
			Values:      allMetrics[m],
		}
		allM = append(allM, metric)
		//fmt.Printf("%v :", m)
		//fmt.Println(allMetrics[m])

	}
	out.Marshal(allM, "json")
	//c, _ := client.Controller()
	//meter := GetOrRegisterBrokerMeter("request-rate", c, r)

	meter := r.Get("request-rate").(metrics.Meter)
	fmt.Println("Rate1 :", meter.Rate1())
	fmt.Println("Rate5 :", meter.Rate5())
	fmt.Println("Rate15:", meter.Rate15())
	fmt.Println("Mean  :", meter.RateMean())

	fmt.Println()
	//metrics.WriteOnce(r, os.Stdout)

	time.Sleep(time.Second * 10)
	fmt.Println("Rate1 :", meter.Rate1())
	fmt.Println("Rate5 :", meter.Rate5())
	fmt.Println("Rate15:", meter.Rate15())
	fmt.Println("Mean  :", meter.RateMean())

	out.Marshal(meter, "json")
	//metrics.WriteOnce(r, os.Stdout)
	//fmt.Println("Final:", meter.Count())
	//fmt.Println(meter.Rate1())
	//meter.Stop()

	allMetrics = r.GetAll()
	for m := range allMetrics {
		fmt.Printf("%v :", m)
		fmt.Println(allMetrics[m])

	}

	meter.Stop()
	r.UnregisterAll()
}
*/

/*
// RawMetric contains the raw metric measurements and values.
type RawMetric struct {
	Measurement string
	Values      map[string]interface{}
	Type        metricType
}

func (m *RawMetric) getMetricType() {
	if m.Type == nil {
		m.Type = make(map[string]bool, 1)
	}
	switch {
	case m.Measurement == "" || m.Values == nil:
		kafkactl.Warnf("metric not initialized!")
		break
	case m.Type[meterMetricType] || m.Type[histoMetricType]:
		break
	default:
		for k := range m.Values {
			switch {
			case strings.Contains(k, "rate"):
				m.Type = meterMetric
				break
			case !strings.Contains(k, "rate") && !strings.Contains(k, "count"):
				m.Type = histoMetric
				break
			}
		}
	}
}

// KafkaMetric represents a metric measurement from Kafka.
type KafkaMetric interface {
	GetType() string
	IsMeter() bool
	IsHisto() bool
}

// MeterMetric contains Meter values.
type MeterMetric struct {
	Measurement string
	Type        string
	Count       uint32
	OneRate     float32
	FiveRate    float32
	FifteenRate float32
	MeanRate    float32
}

// GetType returns the metric type.
func (m *MeterMetric) GetType() string {
	return m.Type
}

// IsMeter returns true if the metric type is a meter.
func (m *MeterMetric) IsMeter() bool {
	return m.Type == meterMetricType
}

// IsHisto returns true if the metric type is a histogram.
func (m *MeterMetric) IsHisto() bool {
	return m.Type == histoMetricType
}

// HistoMetric contains Histogram values.
type HistoMetric struct {
	Measurement string
	Type        string
	Count       uint32
	Min         uint32
	Max         uint32
	SeventyFive float32
	NinetyNine  float32
}

func (m *RawMetric) update(r *metrics.Registry) {
	all := *r
	allMetrics := all.GetAll()
	m.Values = allMetrics[m.Measurement]
}

// Metric Types:
const (
	meterMetricType = `meter`
	histoMetricType = `histogram`
)

type metricType map[string]bool

var meterMetric = metricType{meterMetricType: true}
var histoMetric = metricType{histoMetricType: true}
*/
