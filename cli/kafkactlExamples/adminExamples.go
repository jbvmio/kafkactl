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
