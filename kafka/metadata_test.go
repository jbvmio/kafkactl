package kafka

import (
	"testing"
)

func TestGetTopicMetadata(t *testing.T) {
	seedBroker, controllerBroker := getTestingBrokers(t)
	defer seedBroker.Close()
	defer controllerBroker.Close()
	client, err := getTestingClient(seedBroker)
	if err != nil {
		t.Fatal(err)
	}
	tm, err := client.GetTopicMeta()
	if err != nil {
		t.Fatal(err)
	}
	tom := client.MakeTopicOffsetMap(tm)
	if len(tom) != 1 {
		t.Error("Client returned incorrect number of topics, expected 1, received:", len(tom))
	}
	if len(tom) == 1 {
		toMap := tom[0]
		if len(toMap.PartitionOffsets) != 2 {
			t.Error("Client returned incorrect number of partitions, expected 2, received:", len(toMap.PartitionOffsets))
		}
	}
	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
}
