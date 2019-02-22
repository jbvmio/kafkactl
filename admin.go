package kafkactl

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type clusterAdmin struct {
	client sarama.Client
	conf   *sarama.Config
}

func (kc *KClient) Admin() sarama.ClusterAdmin {
	return &clusterAdmin{
		client: kc.cl,
		conf:   kc.config,
	}
}

func (ca *clusterAdmin) Close() error {
	return ca.client.Close()
}

func (ca *clusterAdmin) Controller() (*sarama.Broker, error) {
	return ca.client.Controller()
}

func (ca *clusterAdmin) CreateTopic(topic string, detail *sarama.TopicDetail, validateOnly bool) error {

	if topic == "" {
		return sarama.ErrInvalidTopic
	}

	if detail == nil {
		return fmt.Errorf("You must specify topic details")
	}

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[topic] = detail

	request := &sarama.CreateTopicsRequest{
		TopicDetails: topicDetails,
		ValidateOnly: validateOnly,
		Timeout:      ca.conf.Admin.Timeout,
	}

	if ca.conf.Version.IsAtLeast(sarama.V0_11_0_0) {
		request.Version = 1
	}
	if ca.conf.Version.IsAtLeast(sarama.V1_0_0_0) {
		request.Version = 2
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.CreateTopics(request)
	if err != nil {
		return err
	}

	topicErr, ok := rsp.TopicErrors[topic]
	if !ok {
		return sarama.ErrIncompleteResponse
	}

	if topicErr.Err != sarama.ErrNoError {
		return topicErr.Err
	}

	return nil
}

func (ca *clusterAdmin) DeleteTopic(topic string) error {

	if topic == "" {
		return sarama.ErrInvalidTopic
	}

	request := &sarama.DeleteTopicsRequest{
		Topics:  []string{topic},
		Timeout: ca.conf.Admin.Timeout,
	}

	if ca.conf.Version.IsAtLeast(sarama.V0_11_0_0) {
		request.Version = 1
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.DeleteTopics(request)
	if err != nil {
		return err
	}

	topicErr, ok := rsp.TopicErrorCodes[topic]
	if !ok {
		return sarama.ErrIncompleteResponse
	}

	if topicErr != sarama.ErrNoError {
		return topicErr
	}
	return nil
}

func (ca *clusterAdmin) CreatePartitions(topic string, count int32, assignment [][]int32, validateOnly bool) error {
	if topic == "" {
		return sarama.ErrInvalidTopic
	}

	topicPartitions := make(map[string]*sarama.TopicPartition)
	topicPartitions[topic] = &sarama.TopicPartition{Count: count, Assignment: assignment}

	request := &sarama.CreatePartitionsRequest{
		TopicPartitions: topicPartitions,
		Timeout:         ca.conf.Admin.Timeout,
		ValidateOnly:    validateOnly,
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.CreatePartitions(request)
	if err != nil {
		return err
	}

	topicErr, ok := rsp.TopicPartitionErrors[topic]
	if !ok {
		return sarama.ErrIncompleteResponse
	}

	if topicErr.Err != sarama.ErrNoError {
		return topicErr.Err
	}

	return nil
}

func (ca *clusterAdmin) DeleteRecords(topic string, partitionOffsets map[int32]int64) error {

	if topic == "" {
		return sarama.ErrInvalidTopic
	}

	topics := make(map[string]*sarama.DeleteRecordsRequestTopic)
	topics[topic] = &sarama.DeleteRecordsRequestTopic{PartitionOffsets: partitionOffsets}
	request := &sarama.DeleteRecordsRequest{
		Topics:  topics,
		Timeout: ca.conf.Admin.Timeout,
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.DeleteRecords(request)
	if err != nil {
		return err
	}

	_, ok := rsp.Topics[topic]
	if !ok {
		return sarama.ErrIncompleteResponse
	}

	//todo since we are dealing with couple of partitions it would be good if we return slice of errors
	//for each partition instead of one error
	return nil
}

func (ca *clusterAdmin) DescribeConfig(resource sarama.ConfigResource) ([]sarama.ConfigEntry, error) {

	var entries []sarama.ConfigEntry
	var resources []*sarama.ConfigResource
	resources = append(resources, &resource)

	request := &sarama.DescribeConfigsRequest{
		Resources: resources,
	}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DescribeConfigs(request)
	if err != nil {
		return nil, err
	}

	for _, rspResource := range rsp.Resources {
		if rspResource.Name == resource.Name {
			if rspResource.ErrorMsg != "" {
				//return nil, errors.New(rspResource.ErrorMsg)
				return nil, fmt.Errorf("%v", rspResource.ErrorMsg)
			}
			for _, cfgEntry := range rspResource.Configs {
				entries = append(entries, *cfgEntry)
			}
		}
	}
	return entries, nil
}

func (ca *clusterAdmin) AlterConfig(resourceType sarama.ConfigResourceType, name string, entries map[string]*string, validateOnly bool) error {

	var resources []*sarama.AlterConfigsResource
	resources = append(resources, &sarama.AlterConfigsResource{
		Type:          resourceType,
		Name:          name,
		ConfigEntries: entries,
	})

	request := &sarama.AlterConfigsRequest{
		Resources:    resources,
		ValidateOnly: validateOnly,
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	rsp, err := b.AlterConfigs(request)
	if err != nil {
		return err
	}

	for _, rspResource := range rsp.Resources {
		if rspResource.Name == name {
			if rspResource.ErrorMsg != "" {
				//return errors.New(rspResource.ErrorMsg)
				return fmt.Errorf(rspResource.ErrorMsg)
			}
		}
	}
	return nil
}

func (ca *clusterAdmin) CreateACL(resource sarama.Resource, acl sarama.Acl) error {
	var acls []*sarama.AclCreation
	acls = append(acls, &sarama.AclCreation{resource, acl})
	request := &sarama.CreateAclsRequest{AclCreations: acls}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	_, err = b.CreateAcls(request)
	return err
}

func (ca *clusterAdmin) ListAcls(filter sarama.AclFilter) ([]sarama.ResourceAcls, error) {

	request := &sarama.DescribeAclsRequest{AclFilter: filter}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DescribeAcls(request)
	if err != nil {
		return nil, err
	}

	var lAcls []sarama.ResourceAcls
	for _, rAcl := range rsp.ResourceAcls {
		lAcls = append(lAcls, *rAcl)
	}
	return lAcls, nil
}

func (ca *clusterAdmin) DeleteACL(filter sarama.AclFilter, validateOnly bool) ([]sarama.MatchingAcl, error) {
	var filters []*sarama.AclFilter
	filters = append(filters, &filter)
	request := &sarama.DeleteAclsRequest{Filters: filters}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DeleteAcls(request)
	if err != nil {
		return nil, err
	}

	var mAcls []sarama.MatchingAcl
	for _, fr := range rsp.FilterResponses {
		for _, mACL := range fr.MatchingAcls {
			mAcls = append(mAcls, *mACL)
		}

	}
	return mAcls, nil
}
