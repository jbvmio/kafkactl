package kafka

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	metrics "github.com/rcrowley/go-metrics"
)

// Metric Types:
const (
	meterMetricType = `meter`
	histoMetricType = `histogram`
)

// Metric Measurement Names:
const (
	// Global Measurements
	MetricCount = `count`
	// Metere Measurements
	MetricOneRate     = `1m.rate`
	MetricFiveRate    = `5m.rate`
	MetricFifteenRate = `15m.rate`
	MetricMeanRate    = `mean.rate`
	// Histogram Measurements
	MetricMin         = `min`
	MetricMax         = `max`
	MetricSeventyFive = `75%`
	MetricNinetyNine  = `99%`
)

type metricType map[string]bool

var meterMetric = metricType{meterMetricType: true}
var histoMetric = metricType{histoMetricType: true}

// RawMetric contains the raw metric measurements and values.
type RawMetric struct {
	Measurement string
	Values      map[string]interface{}
	Type        metricType
}

// Update updates the RawMetric with current values from a metrics registry.
func (m *RawMetric) Update(r *metrics.Registry) {
	all := *r
	allMetrics := all.GetAll()
	m.Values = allMetrics[m.Measurement]
	m.getMetricType()
}

func (m *RawMetric) getMetricType() {
	if m.Type == nil {
		switch {
		case m.Measurement == "" || m.Values == nil:
			Warnf("metric not initialized!")
			break
		default:
			m.Type = make(map[string]bool, 1)
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
}

// ConvertToMetric converts the RawMetric and returns a Metric.
func (m *RawMetric) ConvertToMetric() *Metric {
	var KM Metric
	switch {
	case m.Type[meterMetricType]:
		KM = &MeterMetric{
			Measurement: m.Measurement,
			Type:        meterMetricType,
			Count:       m.Values[MetricCount].(int64),
			OneRate:     m.Values[MetricOneRate].(float64),
			FiveRate:    m.Values[MetricFiveRate].(float64),
			FifteenRate: m.Values[MetricFifteenRate].(float64),
			MeanRate:    m.Values[MetricMeanRate].(float64),
		}
	case m.Type[histoMetricType]:
		KM = &HistoMetric{
			Measurement: m.Measurement,
			Type:        histoMetricType,
			Count:       m.Values[MetricCount].(int64),
			Min:         m.Values[MetricMin].(int64),
			Max:         m.Values[MetricMax].(int64),
			SeventyFive: m.Values[MetricSeventyFive].(float64),
			NinetyNine:  m.Values[MetricNinetyNine].(float64),
		}
	}
	return &KM
}

// MetricCollection contains a collection of Metrics.
type MetricCollection struct {
	Meters     []MeterMetric
	Histograms []HistoMetric
}

// MeterCount returns the number of meters currently in the collection.
func (mc *MetricCollection) MeterCount() int {
	return len(mc.Meters)
}

// HistoCount returns the number of meters currently in the collection.
func (mc *MetricCollection) HistoCount() int {
	return len(mc.Histograms)
}

// AddFromRaw recieves a single or collection of RawMetric types and appends to its appropriate collection type.
func (mc *MetricCollection) AddFromRaw(rawMetrics ...*RawMetric) {
	for _, raw := range rawMetrics {
		mc.Add(raw.ConvertToMetric())
	}
}

// Add recieves a Metric and appends it to its appropriate collection type.
func (mc *MetricCollection) Add(metrics ...*Metric) {
	for _, metric := range metrics {
		m := *metric
		switch {
		case m.IsMeter():
			mc.Meters = append(mc.Meters, *m.(*MeterMetric))
		case m.IsHisto():
			mc.Histograms = append(mc.Histograms, *m.(*HistoMetric))
		}
	}
}

// Metric represents a metric measurement from Kafka.
type Metric interface {
	GetType() string
	IsMeter() bool
	IsHisto() bool
}

// MeterMetric contains Meter values.
type MeterMetric struct {
	Measurement string
	Type        string
	Count       int64
	OneRate     float64
	FiveRate    float64
	FifteenRate float64
	MeanRate    float64
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
	Count       int64
	Min         int64
	Max         int64
	SeventyFive float64
	NinetyNine  float64
}

// GetType returns the metric type.
func (m *HistoMetric) GetType() string {
	return m.Type
}

// IsMeter returns true if the metric type is a meter.
func (m *HistoMetric) IsMeter() bool {
	return m.Type == meterMetricType
}

// IsHisto returns true if the metric type is a histogram.
func (m *HistoMetric) IsHisto() bool {
	return m.Type == histoMetricType
}

// From Sarama GoDocs:
const (
	metricsReservoirSize = 1028
	metricsAlphaFactor   = 0.015
)

func getOrRegisterHistogram(name string, r metrics.Registry) metrics.Histogram {
	return r.GetOrRegister(name, func() metrics.Histogram {
		return metrics.NewHistogram(metrics.NewExpDecaySample(metricsReservoirSize, metricsAlphaFactor))
	}).(metrics.Histogram)
}

func getMetricNameForBroker(name string, broker *sarama.Broker) string {
	// Use broker id like the Java client as it does not contain '.' or ':' characters that
	// can be interpreted as special character by monitoring tool (e.g. Graphite)
	return fmt.Sprintf(name+"-for-broker-%d", broker.ID())
}

func getOrRegisterBrokerMeter(name string, broker *sarama.Broker, r metrics.Registry) metrics.Meter {
	return metrics.GetOrRegisterMeter(getMetricNameForBroker(name, broker), r)
}

func getOrRegisterBrokerHistogram(name string, broker *sarama.Broker, r metrics.Registry) metrics.Histogram {
	return getOrRegisterHistogram(getMetricNameForBroker(name, broker), r)
}

func getMetricNameForTopic(name string, topic string) string {
	// Convert dot to _ since reporters like Graphite typically use dot to represent hierarchy
	// cf. KAFKA-1902 and KAFKA-2337
	return fmt.Sprintf(name+"-for-topic-%s", strings.Replace(topic, ".", "_", -1))
}

func getOrRegisterTopicMeter(name string, topic string, r metrics.Registry) metrics.Meter {
	return metrics.GetOrRegisterMeter(getMetricNameForTopic(name, topic), r)
}

func getOrRegisterTopicHistogram(name string, topic string, r metrics.Registry) metrics.Histogram {
	return getOrRegisterHistogram(getMetricNameForTopic(name, topic), r)
}
