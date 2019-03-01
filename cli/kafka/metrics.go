package kafka

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	metrics "github.com/rcrowley/go-metrics"
)

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

func GetOrRegisterBrokerMeter(name string, broker *sarama.Broker, r metrics.Registry) metrics.Meter {
	return metrics.GetOrRegisterMeter(getMetricNameForBroker(name, broker), r)
}

func GetOrRegisterBrokerHistogram(name string, broker *sarama.Broker, r metrics.Registry) metrics.Histogram {
	return getOrRegisterHistogram(getMetricNameForBroker(name, broker), r)
}

func getMetricNameForTopic(name string, topic string) string {
	// Convert dot to _ since reporters like Graphite typically use dot to represent hierarchy
	// cf. KAFKA-1902 and KAFKA-2337
	return fmt.Sprintf(name+"-for-topic-%s", strings.Replace(topic, ".", "_", -1))
}

func GetOrRegisterTopicMeter(name string, topic string, r metrics.Registry) metrics.Meter {
	return metrics.GetOrRegisterMeter(getMetricNameForTopic(name, topic), r)
}

func GetOrRegisterTopicHistogram(name string, topic string, r metrics.Registry) metrics.Histogram {
	return getOrRegisterHistogram(getMetricNameForTopic(name, topic), r)
}

func TestMetrics() {
	//metricsRegistry = metrics.NewRegistry()
	r := *MetricR()
	//c, _ := client.Controller()
	//meter := GetOrRegisterBrokerMeter("request-rate", c, r)

	/*
		for i := 0; i < 10; i++ {
			d := time.Second * 1
			//metrics.Write(r, d, os.Stdout)
			//metrics.Get
			//fmt.Printf("%+v", stats)
			metrics.WriteOnce(r, os.Stdout)
			time.Sleep(time.Second * 2)
		}
	*/
	meter := r.Get("request-rate").(metrics.Meter)
	fmt.Println("Rate1 :", meter.Rate1())
	fmt.Println("Rate5 :", meter.Rate5())
	fmt.Println("Rate15:", meter.Rate15())
	fmt.Println("Mean  :", meter.RateMean())

	//metrics.WriteOnce(r, os.Stdout)
	time.Sleep(time.Second * 10)
	fmt.Println("Rate1 :", meter.Rate1())
	fmt.Println("Rate5 :", meter.Rate5())
	fmt.Println("Rate15:", meter.Rate15())
	fmt.Println("Mean  :", meter.RateMean())

	//metrics.WriteOnce(r, os.Stdout)
	//fmt.Println("Final:", meter.Count())
	//fmt.Println(meter.Rate1())
	//meter.Stop()
	r.UnregisterAll()
}
