package kafkactlExamples

// examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"

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

// AdminGetReplicas returns Example Usage.
func AdminGetReplicas() string {
	return `  kafkactl admin get replicas <topicName>

Perform Dry Runs for setting (set) replicas:
  kafkactl admin get replicas <topicName> --brokers "1,3,5" --partitions "0,2"
  kafkactl admin get replicas <topicName> --replicas 3 --allparts`
}

// AdminSetReplicas returns Example Usage.
func AdminSetReplicas() string {
	return `  kafkactl admin set replicas <topicName> --brokers "1,3,5" --partitions "0,2"
  kafkactl admin set replicas <topicName> --replicas 3 --allparts
  kafkactl admin set replicas <topicName> --replicas 5 --partitions 1 --dry-run`
}

// AdminMoveFunc returns Example Usage.
func AdminMoveFunc() string {
	return `  kafkactl describe topic <topicName> | kafkactl admin move --brokers "1,2,3"
  kafkactl get topic <topicName> --describe | kafkactl admin move --brokers 5`
}

// AdminGetTopics returns Example Usage.
func AdminGetTopics() string {
	return `  kafkactl admin get topic --config retention.ms <topicName> -x
  kafkactl admin get topic <topicName> --changed`
}

// AdminSetTopics returns Example Usage.
func AdminSetTopics() string {
	return `  kafkactl admin set topic --config retention.ms --value "86400000" <topicName>
  kafkactl admin set topic --reset <topicName>  //reset config back to defaults`
}
