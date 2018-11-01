package kafkactl

type TopicMeta struct {
	Topic           string
	Partition       int32
	Leader          int32
	Replicas        []int32
	ISRs            []int32
	OfflineReplicas []int32
}

func (kc *KClient) GetTopicMeta() ([]TopicMeta, error) {
	var topicMeta []TopicMeta
	res, err := kc.ReqMetadata()
	if err != nil {
		return topicMeta, err
	}
	for _, t := range res.Topics {
		topicName := t.Name
		for _, x := range t.Partitions {
			tm := TopicMeta{
				Topic:           topicName,
				Partition:       x.ID,
				Leader:          x.Leader,
				Replicas:        x.Replicas,
				ISRs:            x.Isr,
				OfflineReplicas: x.OfflineReplicas,
			}
			topicMeta = append(topicMeta, tm)
		}
	}
	return topicMeta, err
}

func (kc *KClient) ListTopics() ([]string, error) {
	res, err := kc.ReqMetadata()
	if err != nil {
		return nil, err
	}
	var topics = make([]string, 0, len(res.Topics))
	for _, t := range res.Topics {
		topics = append(topics, t.Name)
	}
	return topics, nil
}
