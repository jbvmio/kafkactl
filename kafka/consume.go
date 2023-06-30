package kafka

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

// BatchFetchMessages retrieves a batch of messages identified by a topic and partition offset map.
// WIP* Do not use.
/*
func (kc *KClient) BatchFetchMessages(topic string, poMap map[int32][]int64) (messages []*Message, err error) {
	var response *sarama.FetchResponse
	fetchRequest := sarama.FetchRequest{}

	for part, offsets := range poMap {
		for _, o := range offsets {
			fetchRequest.AddBlock(topic, part, o, kc.config.Consumer.Fetch.Default)
		}
	}

	var errd error
	for _, b := range kc.brokers {
		response, errd = b.Fetch(&fetchRequest)
		if errd == nil && response != nil {
			break
		}
		if errd != nil {
			err = errd
		}
	}
	if err != nil {
		return
	}
	if response == nil {
		err = fmt.Errorf("%v", "received nil response")
		return
	}

	blocks, ok := response.Blocks[topic]

	if ok {
		for part, block := range blocks {
			for _, rs := range block.RecordsSet {
				for _, msg := range rs.MsgSet.Messages {
					messages = append(messages, &Message{
						Topic:     topic,
						Key:       msg.Msg.Key,
						Value:     msg.Msg.Value,
						Partition: part,
						Offset:    msg.Offset,
						Timestamp: msg.Msg.Timestamp,
					})
				}
			}
		}
	}
	return
}
*/

// ConsumeOffsetMsg retreives a single message from the given topic, partition and literal offset.
func (kc *KClient) ConsumeOffsetMsg(topic string, partition int32, offset int64) (message *sarama.ConsumerMessage, err error) {
	return kc.GetOffsetMsg(topic, partition, offset)
}

// GetOffsetMsg retreives a single message from the given topic, partition and literal offset without converstion.
func (kc *KClient) GetOffsetMsg(topic string, partition int32, offset int64) (message *sarama.ConsumerMessage, err error) {
	consumer, err := sarama.NewConsumerFromClient(kc.cl)
	if err != nil {
		return
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return
	}
	message = <-partitionConsumer.Messages()
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
// Meant to be used via a goroutine, a Message channel needs to be created and be passed which is used to receive any messages.
// Calling StopPartitionConsumers will stop all ChanPartitionConsume processes.
// Any errors will be passed through the msgChan if received initializing the Consume Loop.
func (kc *KClient) ChanPartitionConsume(topic string, partition int32, offset int64, msgChan chan *sarama.ConsumerMessage) {
	var stopNow bool
	if kc.stopChan == nil {
		kc.stopChan = make(chan none, 1)
	}
	consumer, err := sarama.NewConsumerFromClient(kc.cl)
	if err != nil {
		errMsg := fmt.Sprintf("ERROR: %v", err)
		msgChan <- &sarama.ConsumerMessage{
			Value: []byte(errMsg),
		}
		return
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		errMsg := fmt.Sprintf("ERROR: %v", err)
		msgChan <- &sarama.ConsumerMessage{
			Value: []byte(errMsg),
		}
		return
	}
ConsumeLoop:
	for {
		select {
		case <-kc.stopChan:
			stopNow = true
			break ConsumeLoop
		case msg := <-partitionConsumer.Messages():
			msgChan <- msg
			if stopNow {
				break ConsumeLoop
			}
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

// StopPartitionConsumers signals to stop all spawned ChanPartitionConsume processes.
func (kc *KClient) StopPartitionConsumers() {
	close(kc.stopChan)
}

// PartitionOffsetByTime retreives the most recent available offset at the given time (ms) for the specified topic and partition.
func (kc *KClient) PartitionOffsetByTime(topic string, partition int32, time int64) (int64, error) {
	return kc.cl.GetOffset(topic, partition, time)
}

// OffsetMsgByTime retreives a single message from the given topic, partition and reference time in UTC.
// (datetime string Ex: "11/12/2018 03:59:38.508").
func (kc *KClient) OffsetMsgByTime(topic string, partition int32, datetime string) (message *sarama.ConsumerMessage, err error) {
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
