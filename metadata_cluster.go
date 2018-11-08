package kafkactl

import (
	"sort"

	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
)

type ClusterMeta struct {
	Brokers    []string
	Topics     []string
	Groups     []string
	Controller int32
	Version    int16
}

func (cm ClusterMeta) BrokerCount() int {
	return len(cm.Brokers)
}

func (cm ClusterMeta) TopicCount() int {
	return len(cm.Topics)
}

func (cm ClusterMeta) GroupCount() int {
	return len(cm.Groups)
}

func (kc *KClient) GetClusterMeta() (ClusterMeta, error) {
	cm := ClusterMeta{}
	res, err := kc.ReqMetadata()
	if err != nil {
		return cm, err
	}
	grps, err := kc.ListGroups()
	if err != nil {
		return cm, err
	}
	cm.Controller = res.ControllerID
	for _, b := range res.Brokers {
		id := b.ID()
		addr := b.Addr()
		broker := string(addr + "/" + cast.ToString(id))
		cm.Brokers = append(cm.Brokers, broker)
	}
	for _, t := range res.Topics {
		cm.Topics = append(cm.Topics, t.Name)
	}
	cm.Version = res.Version
	cm.Groups = grps
	sort.Strings(cm.Groups)
	sort.Strings(cm.Brokers)
	sort.Strings(cm.Topics)
	return cm, nil
}

func (kc *KClient) ReqMetadata() (*sarama.MetadataResponse, error) {
	var res *sarama.MetadataResponse
	var err error
	var req = sarama.MetadataRequest{
		AllowAutoTopicCreation: false,
	}
	for _, b := range kc.brokers {
		res, err = b.GetMetadata(&req)
		if err == nil {
			return res, nil
		}
	}
	return res, err
}

/*
func (kc *KClient) GetGroupMetaLite() ([]GroupMetaLite, error) {
	var groupMetaLite []GroupMetaLite
	gl, err := kc.ListGroups()
	if err != nil {
		return groupMetaLite, err
	}
	glm, err := kc.GetGroupListMeta()
	if err != nil {
		return groupMetaLite, err
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
		var topics []string
		for _, v := range grp.Members {
			assign, err := v.GetMemberAssignment()
			for k := range assign.Topics {
				topics = append(topics, k)
			}
		}
		topics = filterUnique(topics)
		groupMetaLite = append(groupMetaLite, group)
	}
	return groupMetaLite, nil
}
*/
