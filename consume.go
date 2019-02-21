package kafkactl

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// Message represents a Kafka message sent/received to/from a topic partition
type Message struct {
	Key, Value []byte
	Topic      string
	Partition  int32
	Offset     int64
	Timestamp  time.Time
}

func convertMsg(m *sarama.ConsumerMessage) (msg *Message) {
	return &Message{
		Key:       m.Key,
		Value:     m.Value,
		Topic:     m.Topic,
		Partition: m.Partition,
		Offset:    m.Offset,
		Timestamp: m.Timestamp,
	}
}

func (m *Message) toSarama() *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic:     m.Topic,
		Key:       sarama.ByteEncoder(m.Key),
		Value:     sarama.ByteEncoder(m.Value),
		Partition: m.Partition,
	}
}

// ConsumeOffsetMsg retreives a single message from the given topic, partition and literal offset.
func (kc *KClient) ConsumeOffsetMsg(topic string, partition int32, offset int64) (message *Message, err error) {
	consumer, err := sarama.NewConsumerFromClient(kc.cl)
	if err != nil {
		return
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return
	}
	msg := <-partitionConsumer.Messages()
	message = convertMsg(msg)
	err = partitionConsumer.Close()
	if err != nil {
		return
	}
	err = consumer.Close()
	if err != nil {
		return
	}
	return
}

// ChanPartitionConsume retreives messages from the given topic, partition and literal offset.
// Meant to be used via a goroutine, a Message channel and bool channel (for stopping the process) should be passed.
// Will return an error through the msgChan if any are encountered before initializing the Consume Loop.
func (kc *KClient) ChanPartitionConsume(topic string, partition int32, offset int64, msgChan chan *Message, stopChan chan bool) {
	consumer, err := sarama.NewConsumerFromClient(kc.cl)
	if err != nil {
		errMsg := fmt.Sprintf("ERROR: %v", err)
		msgChan <- &Message{
			Value: []byte(errMsg),
		}
		return
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		errMsg := fmt.Sprintf("ERROR: %v", err)
		msgChan <- &Message{
			Value: []byte(errMsg),
		}
		return
	}
ConsumeLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			msgChan <- convertMsg(msg)
		case <-stopChan:
			break ConsumeLoop
		}
	}
	if err := partitionConsumer.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := consumer.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// PartitionOffsetByTime retreives the most recent available offset at the given time (ms) for the specified topic and partition.
func (kc *KClient) PartitionOffsetByTime(topic string, partition int32, time int64) (int64, error) {
	return kc.cl.GetOffset(topic, partition, time)
}

// OffsetMsgByTime retreives a single message from the given topic, partition and reference time in UTC.
// (datetime string Ex: "11/12/2018 03:59:38.508").
func (kc *KClient) OffsetMsgByTime(topic string, partition int32, datetime string) (message *Message, err error) {
	timeFormat := "01/02/2006 15:04:05.000"
	targetTime := roundTime(datetime)
	time, err := time.Parse(timeFormat, targetTime)
	if err != nil {
		return
	}
	timeMilli := (time.Unix() * 1000)
	offset, err := kc.PartitionOffsetByTime(topic, partition, timeMilli)
	if err != nil {
		return
	}
	return kc.ConsumeOffsetMsg(topic, partition, offset)
}

func roundTime(targetTime string) string {
	testTargetTime := strings.Split(targetTime, ".")
	if len(testTargetTime) > 1 {
		addZeros := 3 - (len([]rune(testTargetTime[1])))
		if addZeros > 0 {
			for i := 0; i < addZeros; i++ {
				targetTime += "0"
			}
		}
		return targetTime
	}
	if len(testTargetTime) == 1 {
		targetTime += ".000"
	}
	return targetTime
}
