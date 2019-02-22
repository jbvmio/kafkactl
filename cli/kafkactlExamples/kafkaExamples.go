package kafkactlExamples

// examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"

// GetTopics returns Example Usage.
func GetTopics() string {
	return `  kafkactl get topic
  kafkactl get topic <topicName>
  kafkactl get topic <topicName> -x
  kafkactl get topic <topicName> --describe
  kafkactl get topic <topicName> --groups
  kafkactl get topic <topicName> --lag`
}

// GetGroups returns Example Usage.
func GetGroups() string {
	return `  kafkactl get group
  kafkactl get group <groupName>
  kafkactl get group <groupName> -x
  kafkactl get group <groupName> --describe
  kafkactl get group <groupName> --lag`
}

// GetLag returns Example Usage.
func GetLag() string {
	return `  kafkactl get lag
  kafkactl get lag <groupName>
  
Passed as Flags:
  kafkactl get topic <topicName> --lag
  kafkactl get member <memberName> --lag
  kafkactl get group <groupName> --lag`
}

// Describe returns Example Usage.
func Describe() string {
	return `  kafkactl describe topic <topicName>
  kafkactl describe group <groupName>
  kafkactl describe broker <brokerName>`
}

// SEND returns Example Usage.
func SEND() string {
	return `  kafkactl send --allparts --key "myKey" --value "MyTestValue" <topicName>
  kafkactl send --partition 5 --value "MyTestValue" <topicName>
  kafkactl send --partitions "0,1,2,3" --value "MyTestValue" <topicName>

Passed via Stdin:
  printf "myKey,myValue" | kafkactl send --p 5 <topicName>
  printf "key1,value1\nkey2,value2\nkey3,value3" | kafkactl send --allparts <topicName>  //sends 3 messages to all partitions`
}

// LOGS returns Example Usage.
func LOGS() string {
	return `  kafkactl logs <topicName>
  kafkactl logs --tail 5 <topicName>
  kafkactl logs --follow <topicName>
  kafkactl logs --partitions "1,2" --tail 3 --follow <topicName>`
}

// Config returns Example Usage.
func Config() string {
	return `  kafkactl config view
  kafkactl config current-context
  kafkactl config get-contexts
  kafkactl config get-context <contextName>
  kafkactl config use <contextName>`
}
