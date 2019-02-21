package kafkactlExamples

// examples "github.com/jbvmio/kafkactl/cli/kafkactlExamples"

// ZKLS returns Example Usage.
func ZKLS() string {
	return `  kafkactl zk ls /
  kafkactl zk ls /brokers/ids/3
  kafkactl zk ls /brokers --recurse --depth 4`
}

// ZKCreate returns Example Usage.
func ZKCreate() string {
	return `  kafkactl zk create /my/new/path/here
  kafkactl zk create /my/new/value --value '{"value": "myCreatedValue"}'`
}
