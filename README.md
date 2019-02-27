# kafkactl - Kafka Management Tool - **[wiki](https://github.com/jbvmio/kafkactl/wiki)**
kafkactl - CLI for Apache Kafka, Zookeeper and Burrow Management.

[![GitHub release](https://img.shields.io/github/release/jbvmio/kafkactl.svg)](
  https://github.com/jbvmio/kafkactl/releases/latest)
[![Travis-CI Build Status](https://travis-ci.com/jbvmio/kafkactl.svg?branch=master)](
  https://travis-ci.com/jbvmio/kafkactl)
  [![Go Report Card](https://goreportcard.com/badge/github.com/jbvmio/kafkactl)](
  https://goreportcard.com/report/github.com/jbvmio/kafkactl)

#### Current Features:
- [x] Search / Filter for Groups and Topics
- [x] Show Topic Partition Details - Offsets, Replicas, Leaders, etc.
- [x] Show Group Details - Lag, Offsets, Members, etc.
- [x] Display / Modify Topic Configs
- [x] Create / Delete Topics
- [x] Delete Consumer Groups
- [x] Reset Partition Offsets for a Group
- [x] Tail Topics in realtime
- [x] Produce to a Specific Topic Partition or all Partitions at once
- [x] Perform a Preferred Replica Election for a target topic or for all topics
- [x] Increase / Decrease Replicas
- [x] Migrate Topics between Brokers
- [x] Query burrow details
- [x] Monitor burrow lag via Terminal Graph
- [x] Explore Zookeeper Paths
- [x] Create and Delete Zookeeper Values
- [x] Pass stdin to create Kafka messages or Zookeeper Values

#### ToDo:
- [ ] Add Security Features
- [ ] Add Metrics Testing
- [ ] Test Avro Serialization

kafkactl is actively developed with new features being added and tested. Thus, ongoing optimization and re-factoring will occur so ensure you are aware of the [latest releases](https://github.com/jbvmio/kafkactl/releases/latest).


# kafkactl tool - Get Started

### **From Source**
* Requires Go Version 1.11+
```
git clone https://github.com/jbvmio/kafkactl
cd kafkactl/
go build -o $GOPATH/bin/kafkactl

# If inside your go path:
GO111MODULE=on go build -o $GOPATH/bin/kafkactl
```
### **Homebrew**
```
# First time install:
brew tap jbvmio/tap
brew install jbvmio/tap/kafkactl

# Upgrading
brew upgrade jbvmio/tap/kafkactl

```
### **Docker**
```
docker pull docker.io/jbvmio/kafkactl:latest
docker run --rm jbvmio/kafkactl -B broker:9092 get topics --lag
docker run --rm jbvmio/kafkactl -Z zookeeper:2181 ls /brokers/ids/2
docker run --rm jbvmio/kafkactl config --sample
```
### **Manual Download and Install (Linux, Windows, Mac)**
- Download the [latest kafkactl tool](https://github.com/jbvmio/kafkactl/releases/latest) appropriate for your platform, extract and install to a valid directory, preferably located somewhere in your configured $PATH.
- Run "kafkactl config --sample" to generate and display a sample config. Create your config and place as $HOME/.kafkactl.yaml

#### **Alternatively, pass all arguments using flags instead of using a config file when running kafkactl:**
```
# kafkactl --broker brokerhost01:9092 admin create topic newtopic01 --partitions 15 --rfactor 3

```

# kafkactl - Example Commands
### example config commands
```
kafkactl config --sample
kafkactl config view
kafkactl config get-contexts
kafkactl config current-context
kafkactl config use cluster1
```

### example config file (~/.kafkactl.yaml)
```
contexts:
  prod-atl01:
    name: prod-atl01
    brokers:
    - broker01:9092
    - broker02:9092
    - broker03:9092
    burrow:
    - burrow01:8080
    - burrow02:8080
    - burrow03:8080
    zookeeper:
    - zkhost01:2181
    - zkhost02:2181
    - zkhost03:2181
    clientVersion: ""
  prod-atl02:
    name: prod-atl02
    brokers:
    - broker01:9092
    - broker02:9092
    - broker03:9092
    burrow:
    - burrow01:8080
    - burrow02:8080
    - burrow03:8080
    zookeeper:
    - zkhost01:2181
    - zkhost02:2181
    - zkhost03:2181
    clientVersion: ""
  prod-atl03:
    name: prod-atl03
    brokers:
    - broker01:9092
    - broker02:9092
    - broker03:9092
    burrow:
    - burrow01:8080
    - burrow02:8080
    - burrow03:8080
    zookeeper:
    - zkhost01:2181
    - zkhost02:2181
    - zkhost03:2181
    clientVersion: ""
current-context: prod-atl01
config-version: 1

```

# Usage

### kafkactl -h
```
kafkactl: Kafka Management Tool

Usage:
  kafkactl [flags]
  kafkactl [command]

Examples:
  kafkactl --context <contextName> get brokers

Available Commands:
  admin       Kafka Admin Actions
  burrow      Show Burrow Lag Evaluations <wip>
  config      Show and Edit kafkactl config
  describe    Get Kafka Details
  get         Get Kafka Information
  help        Help about any command
  logs        Get Messages from a Kafka Topic
  send        Send/Produce Messages to a Kafka Topic
  version     Print kafkactl version and exit
  zk          Zookeeper Actions

Flags:
  -B, --broker string      Specify a single broker target host:port - Overrides config.
      --burrow string      Specify a single burrow endpoint http://host:port - Overrides config.
      --cfg string         config file (default is $HOME/.kafkactl.yaml)
  -C, --context string     Specify a context.
  -x, --exact              Find exact matches.
  -h, --help               help for kafkactl
  -o, --out string         Change Output Format - yaml|json.
  -v, --verbose            Display additional info or errors.
      --version string     Specify a client version.
  -Z, --zookeeper string   Specify a single zookeeper target host:port - Overrides config.

```

### kafkactl admin -h
```
Available Commands:
  create      Create Kafka Resources
  delete      Delete Kafka Resources
  get         Get Kafka Configurations
  move        Move Partitions using Stdin
  set         Set Kafka Configurations
```

### kafkactl zk -h
```
Zookeeper Actions

Usage:
  kafkactl zk [flags]
  kafkactl zk [command]

Available Commands:
  create      Create Zookeeper Paths and Values
  delete      Delete Zookeeper Paths and Values
  ls          Print Zookeeper Paths and Values

Flags:
  -h, --help         help for zk
  -o, --out string   Change Output Format - yaml|json.

Global Flags:
  -B, --broker string      Specify a single broker target host:port - Overrides config.
      --burrow string      Specify a single burrow endpoint http://host:port - Overrides config.
      --cfg string         config file (default is $HOME/.kafkactl.yaml)
  -C, --context string     Specify a context.
  -x, --exact              Find exact matches.
  -v, --verbose            Display additional info or errors.
      --version string     Specify a client version.
  -Z, --zookeeper string   Specify a single zookeeper target host:port - Overrides config.
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
# kafkactl get group mygroup05

GROUPTYPE  GROUP      COORDINATOR
consumer   mygroup05  broker03:9092/3

```

```
kafkactl get group mygroup05 --lag

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
