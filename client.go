package kafkactl

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

const defaultTimestampFormat = time.RFC3339

type KClient struct {
	apiVers map[int16]int16
	brokers []*sarama.Broker

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
	kc := KClient{
		cl:      client,
		config:  conf,
		brokers: client.Brokers(),
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
	_, err = client.Controller() // Connectivity check
	if err != nil {
		return nil, err
	}
	kc := KClient{
		cl:      client,
		config:  conf,
		brokers: client.Brokers(),
	}
	kc.Connect()
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
	logger := log.New()
	logger.Out = os.Stdout
	logger.Formatter = &log.TextFormatter{
		TimestampFormat: defaultTimestampFormat,
		FullTimestamp:   true,
	}
	sarama.Logger = logger
	//sarama.Logger = log.New(os.Stdout, prefix, log.LstdFlags)
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
	random := makeHex(3)
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
	conf, _ := GetConf()
	for _, b := range brokers {
		broker := sarama.NewBroker(b)
		broker.Open(conf)
		ok, _ := broker.Connected()
		if ok {
			return b, nil
		}
	}
	return "", fmt.Errorf("Error: No Connectivity to Provided Brokers.")
}

func randomBytes(n int) []byte {
	return makeByte(n)
}

func makeHex(s int) string {
	b := randomBytes(s)
	hexstring := hex.EncodeToString(b)
	return hexstring
}

func makeByte(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
