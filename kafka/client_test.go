package kafka

import (
	"sort"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
)

func TestNewClient(t *testing.T) {
	seedBroker, controllerBroker := getTestingBrokers(t)
	client, err := NewClient(seedBroker.Addr(), controllerBroker.Addr())
	if err != nil {
		t.Fatal(err)
	}
	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
	seedBroker.Close()
}

func TestCustomClient(t *testing.T) {
	seedBroker, controllerBroker := getTestingBrokers(t)
	conf := GetConf("testID")
	conf.Metadata.Retry.Max = 0
	client, err := NewCustomClient(conf, seedBroker.Addr(), controllerBroker.Addr())
	if err != nil {
		t.Fatal(err)
	}
	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
	seedBroker.Close()
}

func TestLogging(t *testing.T) {
	m := `testmessage`
	Log("Log Logging Validation:", m)
	Logf("Logf Logging Validation: %v\n", m)
	Warnf("Warnf Logging Validation: %v\n", m)
}

func TestClientLogging(t *testing.T) {
	seedBroker, controllerBroker := getTestingBrokers(t)
	client, err := NewClient(seedBroker.Addr(), controllerBroker.Addr())
	if err != nil {
		t.Fatal(err)
	}
	m := `testmessage`
	client.Log("Log Logging Validation:", m)
	client.Logf("Logf Logging Validation: %v\n", m)
	client.Warnf("Warnf Logging Validation: %v\n", m)
	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
	seedBroker.Close()
}

func TestClusterMetaRequest(t *testing.T) {
	clientTimeout := (time.Second * 5)
	clientRetries := 1
	seedBroker, controllerBroker := getTestingBrokers(t)
	defer seedBroker.Close()
	defer controllerBroker.Close()
	conf := GetConf()
	conf.Net.DialTimeout = clientTimeout
	conf.Net.ReadTimeout = clientTimeout
	conf.Net.WriteTimeout = clientTimeout
	conf.Metadata.Retry.Max = clientRetries
	conf.Version = MinKafkaVersion
	client, err := NewCustomClient(conf, seedBroker.Addr())
	if err != nil {
		t.Fatal(err)
	}
	cm, err := client.clusterMetaTest()
	if err != nil {
		t.Fatal(err)
	}
	if cm.BrokerCount() != 2 {
		t.Error("Client returned incorrect number of available brokers, expected 2, received:", cm.BrokerCount())
	}
	client.Logf("Found %v Brokers, %v Topics, %v Groups", cm.BrokerCount(), cm.TopicCount(), cm.GroupCount())
	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func (kc *KClient) clusterMetaTest() (ClusterMeta, error) {
	cm := ClusterMeta{}
	res, err := kc.ReqMetadata()
	if err != nil {
		return cm, err
	}
	grps, errs := kc.ListGroups()
	if len(errs) > 0 {
		cm.ErrorStack = append(cm.ErrorStack, errs...)
	}
	cm.Controller = res.ControllerID
	for _, b := range res.Brokers {
		id := b.ID()
		addr := b.Addr()
		broker := string(addr + "/" + cast.ToString(id))
		cm.Brokers = append(cm.Brokers, broker)
		cm.BrokerIDs = append(cm.BrokerIDs, id)
	}
	for _, t := range res.Topics {
		cm.Topics = append(cm.Topics, t.Name)
	}
	cm.Groups = grps
	sort.Strings(cm.Groups)
	sort.Strings(cm.Brokers)
	sort.Strings(cm.Topics)
	return cm, nil
}

func getTestingClient(seedBroker *sarama.MockBroker) (*KClient, error) {
	clientTimeout := (time.Second * 5)
	clientRetries := 1
	conf := GetConf()
	conf.Net.DialTimeout = clientTimeout
	conf.Net.ReadTimeout = clientTimeout
	conf.Net.WriteTimeout = clientTimeout
	conf.Metadata.Retry.Max = clientRetries
	conf.Version = MinKafkaVersion
	return NewCustomClient(conf, seedBroker.Addr())
}

func getTestingBrokers(t *testing.T) (seedBroker, controllerBroker *sarama.MockBroker) {
	seedBroker = sarama.NewMockBroker(t, 1)
	controllerBroker = sarama.NewMockBroker(t, 2)
	seedBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetController(controllerBroker.BrokerID()).
			SetBroker(seedBroker.Addr(), seedBroker.BrokerID()).
			SetBroker(controllerBroker.Addr(), controllerBroker.BrokerID()).
			SetLeader("testTopic", 0, seedBroker.BrokerID()).
			SetLeader("testTopic", 1, controllerBroker.BrokerID()),
	})
	return
}

/*
func getTestingBrokers(t *testing.T) (seedBroker, controllerBroker *sarama.MockBroker) {
	seedBroker = sarama.NewMockBroker(t, 1)
	controllerBroker = sarama.NewMockBroker(t, 2)
	seedBroker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetController(controllerBroker.BrokerID()).
			SetBroker(seedBroker.Addr(), seedBroker.BrokerID()).
			SetBroker(controllerBroker.Addr(), controllerBroker.BrokerID()).
			SetLeader("testTopic", 0, seedBroker.BrokerID()).
			SetLeader("testTopic", 1, controllerBroker.BrokerID()),
		"DescribeGroupsRequest": sarama.NewMockDescribeGroupsResponse(t).
			AddGroupDescription("testGroup", sarama.GroupDescription{
				Err:     nil,
				GroupID: "testGroup",
			}),
	})
	return
}
*/
