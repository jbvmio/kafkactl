package kafkactl

import (
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type GroupListMeta struct {
	Group           string
	Type            string
	CoordinatorAddr string
	CoordinatorID   int32
	State           string
}

type GroupMeta struct {
	Group             string
	Topics            []string
	Members           []string
	TopicMemberMap    map[string][]string
	MemberAssignments []MemberMeta
}

type QuickGroupMeta struct {
	Group        string
	State        string
	ProtocolType string
	Protocol     string
}

type MemberMeta struct {
	ClientHost      string
	ClientID        string
	Active          bool
	TopicPartitions map[string][]int32
}

func (kc *KClient) ListGroups() (groups []string, errors []string) {
	var wg sync.WaitGroup
	for _, broker := range kc.brokers {
		wg.Add(1)
		go func(broker *sarama.Broker) {
			grps, err := broker.ListGroups(&sarama.ListGroupsRequest{})
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				for k := range grps.Groups {
					groups = append(groups, k)
				}
			}
			wg.Done()
		}(broker)
	}
	wg.Wait()
	return
}

func (kc *KClient) ListGroups2() (groups []string, errors []string) {
	grpChan := make(chan []string, 100)
	timerChan := make(chan string, len(kc.brokers))
	var toCount uint8
	for _, broker := range kc.brokers {
		go func(broker *sarama.Broker) {
			var tmpGrps []string
			grps, err := broker.ListGroups(&sarama.ListGroupsRequest{})
			if err != nil {
				errors = append(errors, err.Error())
			} else {
				for k := range grps.Groups {
					tmpGrps = append(tmpGrps, k)
				}
			}
			grpChan <- tmpGrps
		}(broker)
		go func(br *sarama.Broker) {
			time.Sleep(time.Second * 10)
			timerChan <- br.Addr()
			br.Close()
		}(broker)
	}
	for i := 0; i < len(kc.brokers); i++ {
		select {
		case grps := <-grpChan:
			if len(grps) > 0 {
				groups = append(groups, grps...)
			}
		case br := <-timerChan:
			errors = append(errors, fmt.Sprintf("%v timed out retrieving metadata.", br))
		}
	}
	if toCount > 0 {
		errors = append(errors, fmt.Sprintf("%v brokers timed out.", toCount))
	}
	return
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
	qgMeta, err := kc.QuickGroupMeta()
	if err != nil {
		return groups, err
	}
	for _, broker := range kc.brokers {
		grps, err := broker.ListGroups(&sarama.ListGroupsRequest{})
		if err != nil {
			return groups, err
		}
		id := broker.ID()
		coord := broker.Addr()
		for k, v := range grps.Groups {
			grp := GroupListMeta{
				Group:           k,
				Type:            v,
				CoordinatorAddr: coord,
				CoordinatorID:   id,
			}
			for _, qg := range qgMeta {
				if qg.Group == k {
					grp.State = qg.State
				}
			}
			groups = append(groups, grp)
		}
	}
	return groups, nil
}

func (kc *KClient) GetGroupMeta() ([]GroupMeta, error) {
	var groupMeta []GroupMeta
	gl, errs := kc.ListGroups()
	if len(errs) > 0 {
		if len(gl) < 1 {
			return groupMeta, fmt.Errorf("%v", errs[len(errs)-1])
		}
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

// QuickGroupMeta returns any additional information not present from a GetGroupMeta request.
func (kc *KClient) QuickGroupMeta() ([]QuickGroupMeta, error) {
	var qGroupMeta []QuickGroupMeta
	gl, errs := kc.ListGroups()
	if len(errs) > 0 {
		if len(gl) < 1 {
			return qGroupMeta, fmt.Errorf("%v", errs[len(errs)-1])
		}
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
		if desc != nil {
			for _, g := range desc.Groups {
				code := int16(g.Err)
				if code == 0 {
					groups = append(groups, g)
				}
			}
		}
	}
	for _, grp := range groups {
		qgm := QuickGroupMeta{
			Group:        grp.GroupId,
			State:        grp.State,
			ProtocolType: grp.ProtocolType,
			Protocol:     grp.Protocol,
		}
		qGroupMeta = append(qGroupMeta, qgm)
	}
	return qGroupMeta, nil
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
