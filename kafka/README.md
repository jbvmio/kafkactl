# kafkactl - Kafka Management Tool - **[wiki](https://github.com/jbvmio/kafkactl/wiki)**
kafkactl - Package and CLI Tool (written in Go) for mgmt of Apache Kafka and related components.

### The kafkactl tool currently features the following:

* Search and Manage Groups and Topics
* Show Topic Partition Details - Offsets, Replicas, Leaders, etc.
* Show Group Details - Lag, Offsets, Members, etc.
* Display / Modify Topic Configs
* Create / Delete Topics
* Delete Consumer Groups
* Reset Partition Offsets for a Group
* Tail a Topic in realtime
* Launch a Consumer Group for a Topic
* Produce to a Specific Topic Partition or all Partitions at once
* Perform a Preferred Replica Election for a target topic or for all topics
* Increase / Decrease Replicas
* Perform Topic Partition Migrations between Brokers
* Query burrow details
* Monitor burrow lag via Terminal Graph
* Explore Zookeeper Paths
* Create and Delete Zookeeper Values
* Pass stdin to create Kafka messages or Zookeeper Values

### ToDo:
- Add Security Features
- Add Metrics Testing

kafkactl is actively developed with new features being added and tested. Thus, ongoing optimization and re-factoring will occur so ensure you are aware of the [latest releases](https://github.com/jbvmio/kafkactl/releases).

# Package
The package is mostly a wrapper around the excellent [sarama library](https://github.com/Shopify/sarama) and for the kafkactl tool itself. The kafkactl tool utilizes a context style config, similar to the kubernetes tool - kubectl, grouping a set of Kafka clusters to corresponding [Zookeeper](https://zookeeper.apache.org/) and [Burrow](https://github.com/linkedin/Burrow) instances.

### package example usage:
Get package:
```
go get -u github.com/jbvmio/kafkactl
```
List Topics:
```
package main

import (
	"fmt"
	"log"

	"github.com/jbvmio/kafkactl"
)

func main() {

	client, err := kafkactl.NewClient("localhost:9092")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	topics, err := client.GetTopicMeta()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	for _, t := range topics {
		fmt.Printf("TOPIC> %v  PARTITION> %v  LEADER> %v\n", t.Topic, t.Partition, t.Leader)
	}

}

```

# kafkactl tool - Get Started

### **Manual Download**
- Download the [latest](https://github.com/jbvmio/kafkactl/releases) kafkactl tool and extract to a $PATH directory.
- Run "kafkactl config --sample" to generate a sample config at $HOME/.kafkactl.yaml
- Edit the config file and save.

* Alternatively, pass all arguments to the command when running.
```
# kafkactl --broker brokerhost01:9092 admin create --topic newtopic01 --partitions 15 --rfactor 3

```

### **Using Docker**
**Run with Docker using cmdline args:**
```
docker run --rm jbvmio/kafkactl -b brokerAddress:9092
docker run --rm jbvmio/kafkactl -b brokerAddress:9092 topics
docker run --rm jbvmio/kafkactl -b brokerAddress:9092 topics mytopic -x -m
docker run --rm jbvmio/kafkactl -b brokerAddress:9092 groups
docker run --rm jbvmio/kafkactl -b brokerAddress:9092 groups mygroup --lag
docker run --rm jbvmio/kafkactl -b brokerAddress:9092 admin create --topic mytopic2 --partitions 5 --rfactor 3

```
**Using a config file with Docker:**
- Generate a sample config:
```
docker run --rm -v ${PWD}:/home/kafkactl jbvmio/kafkactl config --sample

```
- This will generate a sample config named .kafkactl.yaml, edit this file with your values.
- At minimum, specify at least one entry containing a Kafka broker.
```
# Minimum Config:
current: Cluster1
entries:
- name: Cluster1
  kafka:
  - broker01:9092
```
- Now mount your .kafkactl.yaml file whenever running the Docker image:
```
docker run --rm -v ${PWD}/.kafkactl.yaml:/home/kafkactl/.kafkactl.yaml jbvmio/kafkactl
docker run --rm -v ${PWD}/.kafkactl.yaml:/home/kafkactl/.kafkactl.yaml jbvmio/kafkactl topics
docker run --rm -v ${PWD}/.kafkactl.yaml:/home/kafkactl/.kafkactl.yaml jbvmio/kafkactl groups

```

# kafkactl tool - Example Commands
### example config commands
```
kafkactl config --sample
kafkactl config --show
kafkactl config --use testcluster1
```

### example config file (~/.kafkactl.yaml)
```
current: testCluster1
entries:
- name: testCluster1
  kafka:
  - brokerHost1:9092
  - brokerHost2:9092
  burrow:
  - http://burrow1:3000
  - http://burrow2:3000
  zookeeper:
  - http://zk1:2181
  - http://zk2:2181
- name: testCluster2
  kafka:
  - brokerHost1:9092
  - brokerHost2:9092
  burrow:
  - http://burrow1:3000
  - http://burrow2:3000
  zookeeper:
  - http://zk1:2181
  - http://zk2:2181

```

# Usage

### kafkactl -h
```
Usage of kafkactl:
Usage:
  kafkactl [flags]
  kafkactl [command]

Available Commands:
  admin       Perform Various Kafka Administration Tasks
  burrow      Query Burrow
  config      Manage kafkactl Configuration
  consume     Consume from Topics using Consumer Groups.
  describe    Return Topic or Group details
  group       Search and Retrieve Group Info
  help        Help about any command
  message     Retreive Targeted Messages from a Kafka Topic
  meta        Return Metadata
  produce     Produce messages to a Kafka topic
  tail        Tail a Topic
  topic       Search and Retrieve Available Topics
  version     Print kafkactl version and exit
  zk          Perform Various Zookeeper Administration Tasks

Flags:
  -b, --broker string   Bootstrap Kafka Broker
  -g, --group string    Specify a Target Group
  -h, --help            help for kafkactl
      --port string     Port used for Bootstrap Kafka Broker (default "9092")
  -t, --topic string    Specify a Target Topic
  -v, --verbose         Display any additional info and error output.

```

### kafkactl admin -h
```
  conf        Get and Set Available Kafka Related Configs (Topics Only for Now)
  create      Create Topics
  delete      Delete Topics or Groups
  offset      Manage Offsets
  pre         Perform Preferred Replica Election Tasks
```

### kafkactl zk -h
```
Available Commands:
  ls          List or Get Values of a Given Path in Zookeeper

Flags:
  -c, --create string   Create a Zookeeper Path (Use with --value for setting a value)
  -d, --delete string   Delete a Zookeeper Path/Value
  -f, --force           Force Operation
  -h, --help            help for zk
  -s, --server string   Specify a targeted Zookeeper Server and Port (eg. localhost:2181
      --value string    Create a Zookeeper Value (Use with --create to specify the path for the value) Wins over StdIn
```


# Example output:

### kafka
```
# kafkactl

Brokers:  5
 Topics:  208
 Groups:  126

Cluster: (Kafka: v1.1)
  brokerhost01:9092/1
* brokerhost02:9092/2
  brokerhost03:9092/3
  brokerhost04:9092/4
  brokerhost05:9092/5

(*)Controller

```

```
# kafkactl group mygroup05

GROUPTYPE  GROUP      COORDINATOR
consumer   mygroup05  broker03:9092/3

```

```
kafkactl group mygroup05 --lag

GROUP      TOPIC      PART  MEMBER                    OFFSET  LAG
mygroup05  mytopic05  14    mygroup05-77885078-77abc  97      0
mygroup05  mytopic05  1     mygroup05-77885078-77abc  76      0
mygroup05  mytopic05  2     mygroup05-77885078-77abc  84      0
mygroup05  mytopic05  9     mygroup05-77885078-77abc  84      0
mygroup05  mytopic05  12    mygroup05-77885078-77abc  84      0
mygroup05  mytopic05  3     mygroup05-77885078-77abc  135     0
mygroup05  mytopic05  8     mygroup05-77885078-77abc  90      0
mygroup05  mytopic05  7     mygroup05-77885078-77abc  83      0
mygroup05  mytopic05  4     mygroup05-77885078-77abc  111     0
mygroup05  mytopic05  6     mygroup05-77885078-77abc  114     0
mygroup05  mytopic05  0     mygroup05-77885078-77abc  2329    0
mygroup05  mytopic05  5     mygroup05-77885078-77abc  85      0
mygroup05  mytopic05  11    mygroup05-77885078-77abc  139     0
mygroup05  mytopic05  10    mygroup05-77885078-77abc  125     0
mygroup05  mytopic05  13    mygroup05-77885078-77abc  97      0

```

### zookeeper
```
# kafkactl zk ls /cluster/id

VALUE: /cluster/id
 {"version":"1","id":"wwI7vYwOShKnEqVzDhhUYqi"}

```

```
# kafkactl zk ls /

PATH: /
 cluster
 controller
 brokers
 zookeeper
 admin
 isr_change_notification
 log_dir_event_notification
 controller_epoch
 consumers
 burrow
 latest_producer_id_block
 config

# kafkactl zk --create /testpath
Successfully Created: /testpath

# echo 'my testvalue from stdin' | kafkactl zk --create /testpath/testvalue
Successfully Created: /testpath/testvalue

# kafkactl zk ls /testpath/testvalue

VALUE: /testpath/testvalue
 my testvalue from stdin

```

### burrow
```
kafkactl burrow mygroup07 -x

BURROW      GROUP      TOPIC      PARTITION  LAG  TOPICLAG  STATUS
dc01-kafka  mygroup07  mytopic07  0          0    0         OK
dc01-kafka  mygroup07  mytopic07  1          0    0         OK
dc01-kafka  mygroup07  mytopic07  2          0    0         OK
dc01-kafka  mygroup07  mytopic07  3          0    0         OK
dc01-kafka  mygroup07  mytopic07  4          0    0         OK
dc01-kafka  mygroup07  mytopic07  5          0    0         OK
dc01-kafka  mygroup07  mytopic07  6          0    0         OK
dc01-kafka  mygroup07  mytopic07  7          0    0         OK
dc01-kafka  mygroup07  mytopic07  8          0    0         OK
dc01-kafka  mygroup07  mytopic07  9          0    0         OK
dc01-kafka  mygroup07  mytopic07  10         0    0         OK
dc01-kafka  mygroup07  mytopic07  11         0    0         OK
dc01-kafka  mygroup07  mytopic07  12         0    0         OK
dc01-kafka  mygroup07  mytopic07  13         0    0         OK
dc01-kafka  mygroup07  mytopic07  14         0    0         OK

```
