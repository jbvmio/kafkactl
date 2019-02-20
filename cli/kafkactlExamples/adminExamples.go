package kafkactlExamples

// AdminGetOffsets returns Example Usage.
func AdminGetOffsets() string {
	return `  kafkactl admin get offsets <topicName>
  kafkactl admin get offsets <topicName> --groups`
}

// AdminsetOffsets returns Example Usage.
func AdminSetOffsets() string {
	return `  kafkactl admin set offsets <topicName>
  kafkactl admin set offsets <topicName> --group <groupName>
  kafkactl admin set offsets my.topic.here --group my.consumer.group --relative 10 --allparts
  kafkactl admin set offsets my.topic.here --group my.consumer.group --partition 3 --offset 250
  kafkavtl admin set offsets my.topic.here --group my.consumer.group --allparts --newest`
}

func AdminGetReplicas() string {
	return `  kafkactl admin get replicas <topicName>
  kafkactl admin get replicas <topicName> --brokers "1,3,5" --partitions "0,2"
  kafkactl admin get replicas <topicName> --rfactor 3 --allparts`
}

func AdminSetReplicas() string {
	return `  kafkactl admin set replicas <topicName>
  kafkactl admin set replicas <topicName> --brokers "1,3,5" --partitions "0,2"
  kafkactl admin set replicas <topicName> --rfactor 3 --allparts`
}

func AdminMoveFunc() string {
	return `  kafkactl describe topic <topicName> | kafkactl admin move --brokers "1,2,3"
  kafkactl get topic <topicName> --describe | kafkactl admin move --brokers 5`
}
