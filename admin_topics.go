package kafkactl

import (
	"github.com/Shopify/sarama"
)

func (kc *KClient) AddTopic(name string, partitions int32, replication int16) error {
	details := sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replication,
	}
	return kc.Admin().CreateTopic(name, &details, false)
}

func (kc *KClient) RemoveTopic(name string) error {
	return kc.Admin().DeleteTopic(name)
}

// AddPartitions increases the partition count for the given topic.
func (kc *KClient) AddPartitions(name string, count int32, validateOnly bool) error {
	return kc.Admin().CreatePartitions(name, count, nil, validateOnly)
}

func (kc *KClient) GetTopicConfig(topic string, configName ...string) ([]sarama.ConfigEntry, error) {
	if len(configName) > 0 {
		if configName[0] == "" {
			configName = []string{}
		}
	}
	resource := sarama.ConfigResource{
		Name:        topic,
		Type:        sarama.TopicResource,
		ConfigNames: configName,
	}
	return kc.Admin().DescribeConfig(resource)
}

func (kc *KClient) SetTopicConfig(topic, configName, value string) error {
	entry := make(map[string]*string, 1)
	entry[configName] = &value
	return kc.Admin().AlterConfig(sarama.TopicResource, topic, entry, false)
}

func (kc *KClient) DeleteToOffset(topic string, partition int32, offset int64) error {
	partitionOffset := make(map[int32]int64, 1)
	partitionOffset[partition] = offset
	return kc.Admin().DeleteRecords(topic, partitionOffset)
}

func (kc *KClient) RefreshMetadata(topics ...string) error {
	return kc.cl.RefreshMetadata(topics...)
}
