package kafkactl

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type GroupListMeta struct {
	Group       string
	Type        string
	Coordinator string
}

func (kc *KClient) ListGroups() ([]string, error) {
	var groups []string
	for _, broker := range kc.brokers {
		grps, err := broker.ListGroups(&sarama.ListGroupsRequest{})
		if err != nil {
			return groups, err
		}
		for k := range grps.Groups {
			groups = append(groups, k)
		}
	}
	return groups, nil
}

func (kc *KClient) BrokerGroups(brokerID int32) ([]string, error) {
	var groups []string
	for _, b := range kc.brokers {
		if b.ID() == brokerID {
			grps, err := b.ListGroups(&sarama.ListGroupsRequest{})
			if err != nil {
				return groups, err
			}
			for k := range grps.Groups {
				groups = append(groups, k)
			}
			return groups, nil
		}
	}
	return nil, fmt.Errorf("no broker with id %v found", brokerID)
}
