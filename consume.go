package kafkactl

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

// Message represents a Kafka message received from a topic partition
type Message struct {
	Key, Value []byte
	Topic      string
	Partition  int32
	Offset     int64
	Timestamp  time.Time
}

type Header struct {
	Key   []byte
	Value []byte
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

// ConsumeOffset retreives a single message from the given topic, partition and literal offset.
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
	return
}

// ChanPartitionConsume retreives messages from the given topic, partition and literal offset.
// Meant to be used via a goroutine, a Message and bool channel (for stopping the process) should be passed.
// Will return an error through the msgChan if any are encountered before initializing the Consume Loop.
func (kc *KClient) ChanPartitionConsume(topic string, partition int32, offset int64, msgChan chan *Message, stopChan chan bool) { //stopChan chan os.Signal) {
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
		log.Fatalln(err)
	}
}
