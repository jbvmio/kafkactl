package kafkactl

import (
	"fmt"

	"github.com/Shopify/sarama"
)

const (
	// OffsetNewest stands for the log head offset, i.e. the offset that will be
	// assigned to the next message that will be produced to the partition. You
	// can send this to a client's GetOffset method to get this offset, or when
	// calling ConsumePartition to start consuming new messages.
	OffsetNewest int64 = -1
	// OffsetOldest stands for the oldest offset available on the broker for a
	// partition. You can send this to a client's GetOffset method to get this
	// offset, or when calling ConsumePartition to start consuming from the
	// oldest offset that is still available on the broker.
	OffsetOldest int64 = -2
)

type OffsetAdmin interface {
	Group(group string) OffsetAdmin
	Topic(topic string) OffsetAdmin
	Valid() bool
	GetOffsetLag(partition int32) (int64, int64, error)
	GetTotalLag(partitions []int32) (int64, error)
	ResetOffset(partition int32, targetOffset int64) error
}

type offsetAdmin struct {
	grp    string
	top    string
	client sarama.Client
	om     sarama.OffsetManager
	pom    sarama.PartitionOffsetManager
}

func (kc *KClient) OffSetAdmin() OffsetAdmin {
	return &offsetAdmin{
		client: kc.cl,
	}
}

func (oa *offsetAdmin) Group(group string) OffsetAdmin {
	oa.grp = group
	return oa
}

func (oa *offsetAdmin) Topic(topic string) OffsetAdmin {
	oa.top = topic
	return oa
}

func (oa *offsetAdmin) Valid() bool {
	if oa.grp == "" || oa.top == "" {
		return false
	}
	return true
}

// GetOffsetLag returns the current group offset and lag for the given partition.
func (oa *offsetAdmin) GetOffsetLag(partition int32) (groupOffset int64, partitionLag int64, err error) {
	if !oa.Valid() {
		err = fmt.Errorf("No specified Group and/or Topic")
		return
	}
	oa.om, err = sarama.NewOffsetManagerFromClient(oa.grp, oa.client)
	if err != nil {
		return
	}
	oa.pom, err = oa.om.ManagePartition(oa.top, partition)
	if err != nil {
		return
	}
	groupOffset, _ = oa.pom.NextOffset()
	partOffset, err := oa.client.GetOffset(oa.top, partition, sarama.OffsetNewest)
	if err != nil {
		return
	}
	if groupOffset == -1 {
		groupOffset = partOffset
	}
	partitionLag = (partOffset - groupOffset)
	oa.om.Close()
	oa.pom.Close()
	return
}

// Improve or Remove this:
// GetTotalLag returns the sum of total lag given for a group of partitions.
func (oa *offsetAdmin) GetTotalLag(partitions []int32) (totalLag int64, err error) {
	if !oa.Valid() {
		err = fmt.Errorf("No specified Group and/or Topic")
		return
	}
	oa.om, err = sarama.NewOffsetManagerFromClient(oa.grp, oa.client)
	if err != nil {
		return
	}
	for _, partition := range partitions {
		var groupOffset int64
		var partOffset int64
		oa.pom, err = oa.om.ManagePartition(oa.top, partition)
		if err != nil {
			return
		}
		groupOffset, _ = oa.pom.NextOffset()
		partOffset, err = oa.client.GetOffset(oa.top, partition, sarama.OffsetNewest)
		if err != nil {
			return
		}
		if groupOffset == -1 {
			groupOffset = partOffset
		}
		totalLag = (totalLag + (partOffset - groupOffset))
		oa.pom.Close()
	}
	oa.om.Close()
	return
}

func (oa *offsetAdmin) ResetOffset(partition int32, targetOffset int64) (err error) {
	if !oa.Valid() {
		err = fmt.Errorf("No specified Group and/or Topic")
		return
	}
	oa.om, err = sarama.NewOffsetManagerFromClient(oa.grp, oa.client)
	if err != nil {
		return
	}
	oa.pom, err = oa.om.ManagePartition(oa.top, partition)
	if err != nil {
		return
	}
	oa.pom.ResetOffset(targetOffset, "reset_offset")
	oa.om.Close()
	oa.pom.Close()
	return
}
