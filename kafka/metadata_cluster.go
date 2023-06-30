package kafka

import (
	"fmt"
	"sort"

	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
)

// ClusterMeta contains metadata for a Kafka Cluster.
type ClusterMeta struct {
	BrokerIDs      []int32
	Brokers        []string
	Topics         []string
	Groups         []string
	Controller     int32
	APIMaxVersions map[int16]int16
	ErrorStack     []string
}

// BrokerCount returns the number of brokers.
func (cm ClusterMeta) BrokerCount() int {
	return len(cm.Brokers)
}

// TopicCount returns the number of Topics.
func (cm ClusterMeta) TopicCount() int {
	return len(cm.Topics)
}

// GroupCount returns the number of Groups.
func (cm ClusterMeta) GroupCount() int {
	return len(cm.Groups)
}

// BrokerList returns the list of brokers.
func (kc *KClient) BrokerList() ([]string, error) {
	var brokerlist []string
	res, err := kc.ReqMetadata()
	if err != nil {
		return brokerlist, err
	}
	for _, b := range res.Brokers {
		brokerlist = append(brokerlist, b.Addr())
	}
	return brokerlist, nil
}

// BrokerIDMap returns broker addresses by their corresponding IDs.
func (kc *KClient) BrokerIDMap() (map[int32]string, error) {
	brokerMap := make(map[int32]string, len(kc.brokers))
	res, err := kc.ReqMetadata()
	if err != nil {
		return brokerMap, err
	}
	for _, b := range res.Brokers {
		brokerMap[b.ID()] = b.Addr()
	}
	return brokerMap, nil
}

// GetClusterMeta returns Cluster Metadata.
func (kc *KClient) GetClusterMeta() (ClusterMeta, error) {
	cm := ClusterMeta{}
	res, err := kc.ReqMetadata()
	if err != nil {
		return cm, err
	}
	grps, errs := kc.ListGroups()
	if len(errs) > 0 {
		cm.ErrorStack = append(cm.ErrorStack, errs...)
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
	cm.APIMaxVersions, err = kc.GetAPIVersions()
	if err != nil {
		if len(grps) > 0 {
			cm.ErrorStack = append(cm.ErrorStack, err.Error())
		} else {
			return cm, err
		}
	}
	cm.Groups = grps
	sort.Strings(cm.Groups)
	sort.Strings(cm.Brokers)
	sort.Strings(cm.Topics)
	return cm, nil
}

// ReqMetadata returns a metadata response from the first responsive broker.
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
	return nil, fmt.Errorf("No Metadata Response Received")
}

// APIKey Codes
const (
	APIKeyProduce int16 = iota
	APIKeyFetch
	APIKeyListOffsets
	APIKeyMetadata
	APIKeyLeaderAndIsr
	APIKeyStopReplica
	APIKeyUpdateMetadata
	APIKeyControlledShutdown
	APIKeyOffsetCommit
	APIKeyOffsetFetch
	APIKeyFindCoordinator
	APIKeyJoinGroup
	APIKeyHeartbeat
	APIKeyLeaveGroup
	APIKeySyncGroup
	APIKeyDescribeGroups
	APIKeyListGroups
	APIKeySaslHandshake
	APIKeyAPIVersions
	APIKeyCreateTopics
	APIKeyDeleteTopics
	APIKeyDeleteRecords
	APIKeyInitProducerID
	APIKeyOffsetForLeaderEpoch
	APIKeyAddPartitionsToTxn
	APIKeyAddOffsetsToTxn
	APIKeyEndTxn
	APIKeyWriteTxnMarkers
	APIKeyTxnOffsetCommit
	APIKeyDescribeAcls
	APIKeyCreateAcls
	APIKeyDeleteAcls
	APIKeyDescribeConfigs
	APIKeyAlterConfigs
	APIKeyAlterReplicaLogDirs
	APIKeyDescribeLogDirs
	APIKeySaslAuthenticate
	APIKeyCreatePartitions
	APIKeyCreateDelegationToken
	APIKeyRenewDelegationToken
	APIKeyExpireDelegationToken
	APIKeyDescribeDelegationToken
	APIKeyDeleteGroups
	APIKeyElectPreferredLeaders
	APIKeyIncrementalAlterConfigs
	APIKeyAlterPartitionReassignments
	APIKeyListPartitionReassignments
	APIKeyOffsetDelete
	APIKeyDescribeClientQuotas
	APIKeyAlterClientQuotas
)

// APIDescriptions for APIs.
// https://kafka.apache.org/protocol
var APIDescriptions = map[int16]string{
	APIKeyProduce:                     "Produce",
	APIKeyFetch:                       "Fetch",
	APIKeyListOffsets:                 "ListOffsets",
	APIKeyMetadata:                    "Metadata",
	APIKeyLeaderAndIsr:                "LeaderAndIsr",
	APIKeyStopReplica:                 "StopReplica",
	APIKeyUpdateMetadata:              "UpdateMetadata",
	APIKeyControlledShutdown:          "ControlledShutdown",
	APIKeyOffsetCommit:                "OffsetCommit",
	APIKeyOffsetFetch:                 "OffsetFetch",
	APIKeyFindCoordinator:             "FindCoordinator",
	APIKeyJoinGroup:                   "JoinGroup",
	APIKeyHeartbeat:                   "Heartbeat",
	APIKeyLeaveGroup:                  "LeaveGroup",
	APIKeySyncGroup:                   "SyncGroup",
	APIKeyDescribeGroups:              "DescribeGroups",
	APIKeyListGroups:                  "ListGroups",
	APIKeySaslHandshake:               "SaslHandshake",
	APIKeyAPIVersions:                 "ApiVersions",
	APIKeyCreateTopics:                "CreateTopics",
	APIKeyDeleteTopics:                "DeleteTopics",
	APIKeyDeleteRecords:               "DeleteRecords",
	APIKeyInitProducerID:              "InitProducerId",
	APIKeyOffsetForLeaderEpoch:        "OffsetForLeaderEpoch",
	APIKeyAddPartitionsToTxn:          "AddPartitionsToTxn",
	APIKeyAddOffsetsToTxn:             "AddOffsetsToTxn",
	APIKeyEndTxn:                      "EndTxn",
	APIKeyWriteTxnMarkers:             "WriteTxnMarkers",
	APIKeyTxnOffsetCommit:             "TxnOffsetCommit",
	APIKeyDescribeAcls:                "DescribeAcls",
	APIKeyCreateAcls:                  "CreateAcls",
	APIKeyDeleteAcls:                  "DeleteAcls",
	APIKeyDescribeConfigs:             "DescribeConfigs",
	APIKeyAlterConfigs:                "AlterConfigs",
	APIKeyAlterReplicaLogDirs:         "AlterReplicaLogDirs",
	APIKeyDescribeLogDirs:             "DescribeLogDirs",
	APIKeySaslAuthenticate:            "SaslAuthenticate",
	APIKeyCreatePartitions:            "CreatePartitions",
	APIKeyCreateDelegationToken:       "CreateDelegationToken",
	APIKeyRenewDelegationToken:        "RenewDelegationToken",
	APIKeyExpireDelegationToken:       "ExpireDelegationToken",
	APIKeyDescribeDelegationToken:     "DescribeDelegationToken",
	APIKeyDeleteGroups:                "DeleteGroups",
	APIKeyElectPreferredLeaders:       "ElectPreferredLeaders",
	APIKeyIncrementalAlterConfigs:     "IncrementalAlterConfigs",
	APIKeyAlterPartitionReassignments: "AlterPartitionReassignments",
	APIKeyListPartitionReassignments:  "ListPartitionReassignments",
	APIKeyOffsetDelete:                "OffsetDelete",
	APIKeyDescribeClientQuotas:        "DescribeClientQuotas",
	APIKeyAlterClientQuotas:           "AlterClientQuotas",
}

// Kafka API Versions:
var (
	MinKafkaVersion     = sarama.MinVersion
	MaxKafkaVersion     = sarama.MaxVersion
	VER210KafkaVersion  = sarama.V2_1_0_0
	RecKafkaVersion     = sarama.V1_1_0_0
	MinCreatePartsVer   = sarama.V1_0_0_0
	MinDeleteRecordsVer = sarama.V0_11_0_0
	MinTopicOpsVer      = sarama.V0_10_1_0
)

// GetAPIVersions returns a API version Mapping.
func (kc *KClient) GetAPIVersions() (apiMaxVers map[int16]int16, err error) {
	apiMaxVers = make(map[int16]int16)
	apiVers, err := kc.apiVersions()
	if err != nil {
		return
	}
	for _, api := range apiVers.ApiKeys {
		apiMaxVers[api.ApiKey] = api.MaxVersion
	}
	return
}

func (kc *KClient) apiVersions() (*sarama.ApiVersionsResponse, error) {
	var apiRes *sarama.ApiVersionsResponse
	controller, err := kc.cl.Controller()
	if err != nil {
		return apiRes, err
	}
	apiReq := sarama.ApiVersionsRequest{}
	apiRes, err = controller.ApiVersions(&apiReq)
	if err != nil {
		return apiRes, err
	}
	return apiRes, nil
}

// BrokerAPIVersions returns the available API Versions for the given broker.
func BrokerAPIVersions(conf *sarama.Config, broker string) (apiMaxVers map[int16]int16, err error) {
	b := sarama.NewBroker(broker)
	if conf == nil {
		conf = GetConf()
		conf.ClientID = makeHex(3)
		conf.Version = RecKafkaVersion
	}
	b.Open(conf)
	apiReq := sarama.ApiVersionsRequest{}
	apiVers, err := b.ApiVersions(&apiReq)
	if err != nil {
		return
	}
	apiMaxVers = make(map[int16]int16)
	for _, api := range apiVers.ApiKeys {
		apiMaxVers[api.ApiKey] = api.MaxVersion
	}
	return
}

// MatchKafkaVersion parses the given versiona and returns the corresponding KafkaVersion from sarama.
func MatchKafkaVersion(version string) (sarama.KafkaVersion, error) {
	return sarama.ParseKafkaVersion(version)
}
