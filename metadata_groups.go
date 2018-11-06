package kafkactl

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
)

type GroupListMeta struct {
	Group       string
	Type        string
	Coordinator string
}

type GroupMeta struct {
	Group             string
	Topics            []string
	Members           []string
	TopicMemberMap    map[string][]string
	MemberAssignments []MemberMeta
}

type MemberMeta struct {
	ClientHost      string
	ClientID        string
	Active          bool
	TopicPartitions map[string][]int32
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

func (kc *KClient) GetGroupListMeta() ([]GroupListMeta, error) {
	var groups []GroupListMeta
	for _, broker := range kc.brokers {
		grps, err := broker.ListGroups(&sarama.ListGroupsRequest{})
		if err != nil {
			return groups, err
		}
		coord := string(broker.Addr() + "/" + cast.ToString(broker.ID()))
		for k, v := range grps.Groups {
			grp := GroupListMeta{
				Group:       k,
				Type:        v,
				Coordinator: coord,
			}
			groups = append(groups, grp)
		}
	}
	return groups, nil
}

func (kc *KClient) GetGroupMeta() ([]GroupMeta, error) {
	var groupMeta []GroupMeta
	gl, err := kc.ListGroups()
	if err != nil {
		return groupMeta, err
	}
	request := sarama.DescribeGroupsRequest{
		Groups: gl,
	}
	var groups []*sarama.GroupDescription
	for _, broker := range kc.brokers {
		desc, err := broker.DescribeGroups(&request)
		if err != nil {
			fmt.Println("ERROR on desc:", err)
		}
		for _, g := range desc.Groups {
			code := int16(g.Err)
			if code == 0 {
				groups = append(groups, g)
			}
		}
	}
	for _, grp := range groups {
		group := GroupMeta{
			Group:          grp.GroupId,
			TopicMemberMap: make(map[string][]string),
		}
		var topics []string
		var members []string
		for _, v := range grp.Members {
			memberMeta := MemberMeta{
				ClientID:   v.ClientId,
				ClientHost: v.ClientHost,
				Active:     true,
			}
			members = append(members, v.ClientId)
			assign, err := v.GetMemberAssignment()
			if err != nil {
				if err.Error() != `kafka: insufficient data to decode packet, more bytes expected` {
					return groupMeta, err
				} else {
					memberMeta.Active = false
					group.MemberAssignments = append(group.MemberAssignments, memberMeta)
				}
			}
			if memberMeta.Active {
				memberMeta.TopicPartitions = assign.Topics
				group.MemberAssignments = append(group.MemberAssignments, memberMeta)
				for k := range assign.Topics {
					topics = append(topics, k)
					group.TopicMemberMap[k] = append(group.TopicMemberMap[k], memberMeta.ClientID)
				}
			}
		}
		group.Topics = filterUnique(topics)
		group.Members = filterUnique(members)
		groupMeta = append(groupMeta, group)
	}
	return groupMeta, nil
}

// FilterUnique here
func filterUnique(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
