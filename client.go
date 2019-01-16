package kafkactl

import (
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/jbvmio/randstr"
)

type KClient struct {
	bootStrap string
	brokers   []*sarama.Broker

	cl     sarama.Client
	ca     clusterAdmin
	config *sarama.Config
	logger *log.Logger
}

func NewClient(brokerList ...string) (*KClient, error) {
	conf, err := GetConf()
	if err != nil {
		return nil, err
	}
	client, err := getClient(conf, brokerList...)
	if err != nil {
		return nil, err
	}
	_, err = client.Controller() // Connectivity check
	if err != nil {
		return nil, err
	}
	bootStrap, err := ReturnFirstValid(brokerList...)
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

func NewCustomClient(conf *sarama.Config, brokerList ...string) (*KClient, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	client, err := getClient(conf, brokerList...)
	if err != nil {
		return nil, err
	}
	C, err := client.Controller() // Connectivity check
	if err != nil {
		return nil, err
	}
	kc := KClient{
		bootStrap: C.Addr(),
		cl:        client,
		config:    conf,
		brokers:   client.Brokers(),
	}
	kc.Connect()
	return &kc, nil
}

func (kc *KClient) Controller() (*sarama.Broker, error) {
	return kc.cl.Controller()
}

func (kc *KClient) APIVersions() (*sarama.ApiVersionsResponse, error) {
	var apiRes *sarama.ApiVersionsResponse
	controller, err := kc.cl.Controller()
	if err != nil {
		return apiRes, err
	}
	apiReq := sarama.ApiVersionsRequest{}
	apiRes, err = controller.ApiVersions(&apiReq)
	if err != nil {
		return apiRes, err
	}
	return apiRes, nil
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

func (kc *KClient) IsConnected() bool {
	for _, broker := range kc.brokers {
		if ok, _ := broker.Connected(); ok {
			return true
		}
	}
	return false
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

func Logger(prefix string) {
	if prefix == "" {
		prefix = "[kafkactl] "
	}
	sarama.Logger = log.New(os.Stdout, prefix, log.LstdFlags)
}

func (kc *KClient) Logf(format string, v ...interface{}) {
	sarama.Logger.Printf(format, v...)
}

func (kc *KClient) Log(v ...interface{}) {
	sarama.Logger.Println(v...)
}

// SaramaConfig returns and allows modification of the underlying sarama client config
func (kc *KClient) SaramaConfig() *sarama.Config {
	return kc.config
}

func GetConf() (*sarama.Config, error) {
	random := randstr.Hex(3)
	conf := sarama.NewConfig()
	conf.ClientID = string("kafkactl" + "-" + random)
	conf.Version = MinKafkaVersion
	err := conf.Validate()
	return conf, err
}

func getClient(conf *sarama.Config, brokers ...string) (sarama.Client, error) {
	return sarama.NewClient(brokers, conf)
}

// ReturnFirstValid returns the first available, connectable broker provided from a broker list
func ReturnFirstValid(brokers ...string) (string, error) {
	for _, b := range brokers {
		broker := sarama.NewBroker(b)
		broker.Open(nil)
		ok, _ := broker.Connected()
		if ok {
			return b, nil
		}
	}
	return "", fmt.Errorf("Error: No Connectivity to Provided Brokers.")
}

func ReturnAPIVersions(broker string) (apiMaxVers, apiMinVers map[int16]int16, err error) {
	b := sarama.NewBroker(broker)
	conf, err := GetConf()
	if err != nil {
		return
	}
	b.Open(conf)
	apiReq := sarama.ApiVersionsRequest{}
	apiVers, err := b.ApiVersions(&apiReq)
	if err != nil {
		return
	}
	apiMaxVers = make(map[int16]int16)
	apiMinVers = make(map[int16]int16)
	for _, api := range apiVers.ApiVersions {
		apiMaxVers[api.ApiKey] = api.MaxVersion
		apiMinVers[api.ApiKey] = api.MinVersion
	}
	return
}

func MatchKafkaVersion(version string) (sarama.KafkaVersion, error) {
	return sarama.ParseKafkaVersion(version)
}
