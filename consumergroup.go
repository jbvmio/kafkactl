package kafkactl

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type DEBUG struct {
	HasData bool
	IsNote  bool
	IsErr   bool
	Err     error
	Notification
}

type Notification struct {
	Type     string
	Claimed  map[string][]int32
	Released map[string][]int32
	Current  map[string][]int32
}

type ConsumerGroup struct {
	consumer      *cluster.Consumer
	debugChan     chan *DEBUG
	haveDebugChan chan bool
	debugEnabled  bool
}

// CG returns the underlying Consumer from sarama-cluster
func (cg *ConsumerGroup) CG() *cluster.Consumer {
	return cg.consumer
}

func (cg *ConsumerGroup) Close() error {
	if cg.debugEnabled {
		close(cg.debugChan)
	}
	return cg.consumer.Close()
}

// NewConsumerGroup returns a new ConsumerGroup
func (kc *KClient) NewConsumerGroup(groupID string, debug bool, topics ...string) (*ConsumerGroup, error) {
	var cg ConsumerGroup
	var dChan chan *DEBUG
	kc.config.Producer.RequiredAcks = sarama.WaitForAll
	kc.config.Producer.Return.Successes = true
	config := cluster.NewConfig()
	conf := kc.config
	config.Config = *conf
	config.Group.Mode = cluster.ConsumerModePartitions
	if debug {
		config.Consumer.Return.Errors = true
		config.Group.Return.Notifications = true
		dChan = make(chan *DEBUG, 256)
	}
	brokers, err := kc.BrokerList()
	if err != nil {
		return &cg, err
	}
	consumer, err := cluster.NewConsumer(brokers, groupID, topics, config)
	if err != nil {
		return &cg, err
	}
	cg.consumer = consumer
	if debug {
		cg.debugChan = dChan
		cg.haveDebugChan = make(chan bool, 1)
		cg.debugEnabled = true
		go startDEBUG(&cg)
	}
	return &cg, nil
}

// ConsumeMessages here
func (cg *ConsumerGroup) ConsumeMessages(pc cluster.PartitionConsumer, iterator func(msg *Message) bool) {
ProcessMessages:
	for m := range pc.Messages() {
		message := convertMsg(m)
		if !iterator(message) {
			break ProcessMessages
		}
		cg.consumer.MarkOffset(m, "") // mark message as processed
	}
	cg.consumer.CommitOffsets()
}

// DeBug here
func (cg *ConsumerGroup) DeBug() <-chan *DEBUG {
	return cg.debugChan
}

// DebugAvailale here
func (cg *ConsumerGroup) DebugAvailale() <-chan bool {
	return cg.haveDebugChan
}

// DeBug here
func startDEBUG(cg *ConsumerGroup) {
	for {
		select {
		case note, ok := <-cg.consumer.Notifications():
			if !ok {
				break
			}
			d := DEBUG{}
			d.Type = note.Type.String()
			d.Claimed = note.Claimed
			d.Released = note.Released
			d.Current = note.Current
			d.HasData = true
			d.IsNote = true
			//fmt.Printf("%+v\n", d)
			cg.debugChan <- &d
			cg.haveDebugChan <- true
		case err, ok := <-cg.consumer.Errors():
			if !ok {
				break
			}
			if err != nil {
				cg.debugChan <- &DEBUG{
					HasData: true,
					Err:     err,
				}
			}
		}
	}
}

// ProcessDEBUG here
func (cg *ConsumerGroup) ProcessDEBUG(iterator func(d *DEBUG) bool) {
ProcessDEBUG:
	for m := range cg.debugChan {
		if !iterator(m) {
			break ProcessDEBUG
		}
	}
}
