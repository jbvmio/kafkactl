package kafkactl

import (
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

type KClient struct {
	bootStrap string
	brokers   []*sarama.Broker

	cl     sarama.Client
	ca     clusterAdmin
	config *sarama.Config
	logger *log.Logger
}

func NewClient(bootStrap string) (*KClient, error) {
	conf, err := getConf()
	if err != nil {
		return nil, err
	}
	client, err := getClient(bootStrap, conf)
	if err != nil {
		return nil, err
	}
	_, err = client.Controller() // Connectivity check
	if err != nil {
		return nil, err
	}
	kc := KClient{
		bootStrap: bootStrap,
		cl:        client,
		config:    conf,
		brokers:   client.Brokers(),
	}
	kc.Connect()
	return &kc, nil
}

func NewCustomClient(bootStrap string, conf *sarama.Config) (*KClient, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	client, err := getClient(bootStrap, conf)
	if err != nil {
		return nil, err
	}
	_, err = client.Controller() // Connectivity check
	if err != nil {
		return nil, err
	}
	kc := KClient{
		bootStrap: bootStrap,
		cl:        client,
		config:    conf,
		brokers:   client.Brokers(),
	}
	return &kc, nil
}

func (kc *KClient) Controller() (*sarama.Broker, error) {
	return kc.cl.Controller()
}

func (kc *KClient) Connect() error {
	for _, broker := range kc.brokers {
		if ok, _ := broker.Connected(); !ok {
			if err := broker.Open(kc.config); err != nil {
				return fmt.Errorf("Error connecting to broker: %v", err)
			}
		}
	}
	return nil
}

func (kc *KClient) Close() error {
	errString := fmt.Sprintf("ERROR:\n")
	var errFound bool
	for _, b := range kc.brokers {
		if ok, _ := b.Connected(); ok {
			if err := b.Close(); err != nil {
				e := fmt.Sprintf("  Error closing broker, %v: %v\n", b.Addr(), err)
				errString = string(errString + e)
				errFound = true
			}
		}
	}
	if err := kc.cl.Close(); err != nil {
		e := fmt.Sprintf("  Error closing client: %v\n", err)
		errString = string(errString + e)
		errFound = true
	}
	if errFound {
		return fmt.Errorf("%v", errString)
	}
	return nil
}

func (kc *KClient) Logger(prefix string) {
	if prefix == "" {
		prefix = "[kafkactl] "
	}
	sarama.Logger = log.New(os.Stdout, prefix, log.LstdFlags)
}

func getConf() (*sarama.Config, error) {
	conf := sarama.NewConfig()
	conf.ClientID = "kafkactl"
	conf.Version = sarama.V1_1_0_0
	err := conf.Validate()
	return conf, err
}

func getClient(broker string, conf *sarama.Config) (sarama.Client, error) {
	return sarama.NewClient([]string{broker}, conf)
}
