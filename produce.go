package kafkactl

import "github.com/Shopify/sarama"

// SendMSG sends a message to the targeted topic/partition defined in the message.
func (kc *KClient) SendMSG(message *Message) (partition int32, offset int64, err error) {
	producer, err := sarama.NewSyncProducerFromClient(kc.cl)
	if err != nil {
		return
	}
	partition, offset, err = producer.SendMessage(message.toSarama())
	producer.Close()
	return
}

// SendMessages sends groups of messages to the targeted topic/partition defined in each message.
func (kc *KClient) SendMessages(messages []*Message) (err error) {
	producer, err := sarama.NewSyncProducerFromClient(kc.cl)
	if err != nil {
		return
	}
	var msgs []*sarama.ProducerMessage
	for _, message := range messages {
		msgs = append(msgs, message.toSarama())
	}
	err = producer.SendMessages(msgs) //.SendMessage(message.toSarama())
	producer.Close()
	return
}
