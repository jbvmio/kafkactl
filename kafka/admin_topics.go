package kafka

import (
	"github.com/Shopify/sarama"
)

// AddTopic creates a new Topic.
func (kc *KClient) AddTopic(name string, partitions int32, replication int16) error {
	details := sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replication,
	}
	return kc.Admin().CreateTopic(name, &details, false)
}

// RemoveTopic deletes a Topic.
func (kc *KClient) RemoveTopic(name string) error {
	return kc.Admin().DeleteTopic(name)
}

// AddPartitions increases the partition count for the given topic.
func (kc *KClient) AddPartitions(name string, count int32, validateOnly bool) error {
	return kc.Admin().CreatePartitions(name, count, nil, validateOnly)
}

// GetTopicConfig returns the configuration for the given topic.
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

// SetTopicConfig sets a Topic Configuration parameter.
func (kc *KClient) SetTopicConfig(topic, configName, value string) error {
	entry := make(map[string]*string)
	existing, err := kc.GetTopicConfig(topic)
	if err != nil {
		return err
	}
	for _, e := range existing {
		val := e.Value
		entry[e.Name] = &val
	}
	entry[configName] = &value
	return kc.Admin().AlterConfig(sarama.TopicResource, topic, entry, false)
}

// ResetTopicConfig resets a Topics configuration.
func (kc *KClient) ResetTopicConfig(topic string) error {
	entry := make(map[string]*string)
	return kc.Admin().AlterConfig(sarama.TopicResource, topic, entry, false)
}

// DeleteToOffset deletes a partition to the desired offset.
func (kc *KClient) DeleteToOffset(topic string, partition int32, offset int64) error {
	partitionOffset := make(map[int32]int64, 1)
	partitionOffset[partition] = offset
	return kc.Admin().DeleteRecords(topic, partitionOffset)
}

// RefreshMetadata refreshes metadata for the given topic.
func (kc *KClient) RefreshMetadata(topics ...string) error {
	return kc.cl.RefreshMetadata(topics...)
}
