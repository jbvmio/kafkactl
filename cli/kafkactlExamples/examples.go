package kafkactlExamples

// GetLag returns Example Usage.
func GetLag() string {
	return `  kafkactl get lag
  kafkactl get lag <groupName>
  
Passed as Flags:
  kafkactl get topic <topicName> --lag
  kafkactl get member <memberName> --lag
  kafkactl get group <groupName> --lag`
}
