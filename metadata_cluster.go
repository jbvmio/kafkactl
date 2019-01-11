package kafkactl

import (
	"fmt"
	"sort"

	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
)

type ClusterMeta struct {
	BrokerIDs      []int32
	Brokers        []string
	Topics         []string
	Groups         []string
	Controller     int32
	APIMinVersions map[int16]int16
	APIMaxVersions map[int16]int16
	Errors         []string
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

func (kc *KClient) BrokerList() ([]string, error) {
	var brokerlist []string
	res, err := kc.ReqMetadata()
	if err != nil {
		return brokerlist, err
	}
	for _, b := range res.Brokers {
		id := b.ID()
		addr := b.Addr()
		broker := string(addr + "/" + cast.ToString(id))
		brokerlist = append(brokerlist, broker)
	}
	return brokerlist, nil
}

func (kc *KClient) GetClusterMeta() (ClusterMeta, error) {
	cm := ClusterMeta{}
	res, err := kc.ReqMetadata()
	if err != nil {
		return cm, err
	}
	grps, errs := kc.ListGroups()
	if len(errs) > 0 {
		if len(grps) < 1 {
			return cm, fmt.Errorf("%v", errs[len(errs)-1])
		}
		if len(errs) > 0 {
			cm.Errors = append(cm.Errors, errs...)
		}
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
	cm.APIMaxVersions, cm.APIMinVersions, err = kc.GetAPIVersions()
	if err != nil {
		return cm, err
	}
	cm.Groups = grps
	sort.Strings(cm.Groups)
	sort.Strings(cm.Brokers)
	sort.Strings(cm.Topics)
	return cm, nil
}

func (kc *KClient) GetAPIVersions() (apiMaxVers, apiMinVers map[int16]int16, err error) {
	apiMaxVers = make(map[int16]int16)
	apiMinVers = make(map[int16]int16)
	apiReq := sarama.ApiVersionsRequest{}
	controller, err := kc.Controller()
	if err != nil {
		return
	}
	apiVers, err := controller.ApiVersions(&apiReq)
	if err != nil {
		return
	}
	for _, api := range apiVers.ApiVersions {
		apiMaxVers[api.ApiKey] = api.MaxVersion
		apiMinVers[api.ApiKey] = api.MinVersion
	}
	return
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

// APIKey Codes
const (
	APIKeyProduce                 int16 = 0
	APIKeyFetch                   int16 = 1
	APIKeyListOffsets             int16 = 2
	APIKeyMetadata                int16 = 3
	APIKeyLeaderAndIsr            int16 = 4
	APIKeyStopReplica             int16 = 5
	APIKeyUpdateMetadata          int16 = 6
	APIKeyControlledShutdown      int16 = 7
	APIKeyOffsetCommit            int16 = 8
	APIKeyOffsetFetch             int16 = 9
	APIKeyFindCoordinator         int16 = 10
	APIKeyJoinGroup               int16 = 11
	APIKeyHeartbeat               int16 = 12
	APIKeyLeaveGroup              int16 = 13
	APIKeySyncGroup               int16 = 14
	APIKeyDescribeGroups          int16 = 15
	APIKeyListGroups              int16 = 16
	APIKeySaslHandshake           int16 = 17
	APIKeyApiVersions             int16 = 18
	APIKeyCreateTopics            int16 = 19
	APIKeyDeleteTopics            int16 = 20
	APIKeyDeleteRecords           int16 = 21
	APIKeyInitProducerId          int16 = 22
	APIKeyOffsetForLeaderEpoch    int16 = 23
	APIKeyAddPartitionsToTxn      int16 = 24
	APIKeyAddOffsetsToTxn         int16 = 25
	APIKeyEndTxn                  int16 = 26
	APIKeyWriteTxnMarkers         int16 = 27
	APIKeyTxnOffsetCommit         int16 = 28
	APIKeyDescribeAcls            int16 = 29
	APIKeyCreateAcls              int16 = 30
	APIKeyDeleteAcls              int16 = 31
	APIKeyDescribeConfigs         int16 = 32
	APIKeyAlterConfigs            int16 = 33
	APIKeyAlterReplicaLogDirs     int16 = 34
	APIKeyDescribeLogDirs         int16 = 35
	APIKeySaslAuthenticate        int16 = 36
	APIKeyCreatePartitions        int16 = 37
	APIKeyCreateDelegationToken   int16 = 38
	APIKeyRenewDelegationToken    int16 = 39
	APIKeyExpireDelegationToken   int16 = 40
	APIKeyDescribeDelegationToken int16 = 41
	APIKeyDeleteGroups            int16 = 42
)

// APIKey Descriptions
var APIDescriptions = map[int16]string{
	APIKeyProduce:                 "Produce",
	APIKeyFetch:                   "Fetch",
	APIKeyListOffsets:             "ListOffsets",
	APIKeyMetadata:                "Metadata",
	APIKeyLeaderAndIsr:            "LeaderAndIsr",
	APIKeyStopReplica:             "StopReplica",
	APIKeyUpdateMetadata:          "UpdateMetadata",
	APIKeyControlledShutdown:      "ControlledShutdown",
	APIKeyOffsetCommit:            "OffsetCommit",
	APIKeyOffsetFetch:             "OffsetFetch",
	APIKeyFindCoordinator:         "FindCoordinator",
	APIKeyJoinGroup:               "JoinGroup",
	APIKeyHeartbeat:               "Heartbeat",
	APIKeyLeaveGroup:              "LeaveGroup",
	APIKeySyncGroup:               "SyncGroup",
	APIKeyDescribeGroups:          "DescribeGroups",
	APIKeyListGroups:              "ListGroups",
	APIKeySaslHandshake:           "SaslHandshake",
	APIKeyApiVersions:             "ApiVersions",
	APIKeyCreateTopics:            "CreateTopics",
	APIKeyDeleteTopics:            "DeleteTopics",
	APIKeyDeleteRecords:           "DeleteRecords",
	APIKeyInitProducerId:          "InitProducerId",
	APIKeyOffsetForLeaderEpoch:    "OffsetForLeaderEpoch",
	APIKeyAddPartitionsToTxn:      "AddPartitionsToTxn",
	APIKeyAddOffsetsToTxn:         "AddOffsetsToTxn",
	APIKeyEndTxn:                  "EndTxn",
	APIKeyWriteTxnMarkers:         "WriteTxnMarkers",
	APIKeyTxnOffsetCommit:         "TxnOffsetCommit",
	APIKeyDescribeAcls:            "DescribeAcls",
	APIKeyCreateAcls:              "CreateAcls",
	APIKeyDeleteAcls:              "DeleteAcls",
	APIKeyDescribeConfigs:         "DescribeConfigs",
	APIKeyAlterConfigs:            "AlterConfigs",
	APIKeyAlterReplicaLogDirs:     "AlterReplicaLogDirs",
	APIKeyDescribeLogDirs:         "DescribeLogDirs",
	APIKeySaslAuthenticate:        "SaslAuthenticate",
	APIKeyCreatePartitions:        "CreatePartitions",
	APIKeyCreateDelegationToken:   "CreateDelegationToken",
	APIKeyRenewDelegationToken:    "RenewDelegationToken",
	APIKeyExpireDelegationToken:   "ExpireDelegationToken",
	APIKeyDescribeDelegationToken: "DescribeDelegationToken",
	APIKeyDeleteGroups:            "DeleteGroups",
}
