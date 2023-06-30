package kafka

import (
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

// SendMSG sends a message to the targeted topic/partition defined in the message.
func (kc *KClient) SendMSG(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	producer, err := sarama.NewSyncProducerFromClient(kc.cl)
	if err != nil {
		return
	}
	partition, offset, err = producer.SendMessage(message)
	producer.Close()
	return
}

// SendMessages sends groups of messages to the targeted topic/partition defined in each message.
func (kc *KClient) SendMessages(messages []*sarama.ProducerMessage) (err error) {
	producer, err := sarama.NewSyncProducerFromClient(kc.cl)
	if err != nil {
		return
	}
	err = producer.SendMessages(messages) //.SendMessage(message.toSarama())
	producer.Close()
	return
}

// Producer is the implementation of an AsyncProducer.
type Producer struct {
	producer         sarama.AsyncProducer
	cl               *KClient
	killChan         chan struct{}
	errors           chan *sarama.ProducerError
	input, successes chan *sarama.ProducerMessage
	wg               sync.WaitGroup
	noClose          bool
}

// NewProducer returns a new Producer with Successes and Error channels that must be read from.
func NewProducer(addrs []string, config *sarama.Config) (*Producer, error) {
	var producer Producer
	if config == nil {
		config = GetConf("")
		config.Version = RecKafkaVersion
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
		config.Producer.RequiredAcks = 1
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	}
	kc, err := NewCustomClient(config, addrs...)
	if err != nil {
		return &producer, err
	}
	p, err := sarama.NewAsyncProducerFromClient(kc.cl)
	if err != nil {
		return &producer, err
	}
	producer.producer = p
	producer.cl = kc
	producer.killChan = make(chan struct{})
	producer.errors = make(chan *sarama.ProducerError)
	producer.input = make(chan *sarama.ProducerMessage)
	producer.successes = make(chan *sarama.ProducerMessage)
	producer.wg = sync.WaitGroup{}
	go producer.start()
	return &producer, nil
}

// NewProducer returns a new Producer with Successes and Error channels that must be read from.
func (kc *KClient) NewProducer() (*Producer, error) {
	/*	Producer Options to be aware of:
		config.Version = RecKafkaVersion
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
		config.Producer.RequiredAcks = 1
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	*/
	p, err := sarama.NewAsyncProducerFromClient(kc.cl)
	if err != nil {
		return &Producer{}, err
	}
	producer := Producer{
		producer:  p,
		cl:        kc,
		killChan:  make(chan struct{}),
		errors:    make(chan *sarama.ProducerError),
		input:     make(chan *sarama.ProducerMessage),
		successes: make(chan *sarama.ProducerMessage),
		wg:        sync.WaitGroup{},
		noClose:   true,
	}
	go producer.start()
	return &producer, nil
}

// RequiredAcks sets the desired number of replica acknowledgements it must see from the broker
// when producing messages. Default is 1.
// 0 = NoResponse doesn't send any response, the TCP ACK is all you get.
// 1 = WaitForLocal waits for only the local commit to succeed before responding. (default)
// 2 = WaitForAll waits for all in-sync replicas to commit before responding.
func (p *Producer) RequiredAcks(n int) {
	switch n {
	case 0:
		p.cl.SaramaConfig().Producer.RequiredAcks = 0
	case 2:
		p.cl.SaramaConfig().Producer.RequiredAcks = -1
	default:
		p.cl.SaramaConfig().Producer.RequiredAcks = 1
	}
}

// SetPartitioner sets the desired Partitioner which decides which partition messages are sent.
// 0 - ManualPartitioner
// 1 - RandomPartitioner
// 2 - HashPartitioner
// 3 - RoundRobinPartitioner (default)
func (p *Producer) SetPartitioner(n int) {
	switch n {
	case 0:
		p.cl.SaramaConfig().Producer.Partitioner = sarama.NewManualPartitioner
	case 1:
		p.cl.SaramaConfig().Producer.Partitioner = sarama.NewRandomPartitioner
	case 2:
		p.cl.SaramaConfig().Producer.Partitioner = sarama.NewHashPartitioner
	default:
		p.cl.SaramaConfig().Producer.Partitioner = sarama.NewRoundRobinPartitioner
	}
}

// AsyncClose triggers a shutdown of the producer. The shutdown has completed
// when both the Errors and Successes channels have been closed. When calling
// AsyncClose, you *must* continue to read from those channels in order to
// drain the results of any messages in flight.
func (p *Producer) AsyncClose() {
	p.producer.AsyncClose()
}

// Close shuts down the producer and waits for any buffered messages to be
// flushed. You must call this function before a producer object passes out of
// scope, as it may otherwise leak memory. You must call this before calling
// Close on the underlying client.
func (p *Producer) Close() error {
	return p.producer.Close()
}

// Shutdown starts the shutdown of the producer, draining both Errors and Successes channels, then closes the underlying client.
// A given timeout value (in seconds) can be used to specify an allotted time for draining the channels. (default: 10)
// If the timer expires and error will be returned.
func (p *Producer) Shutdown(timeout ...int) error {
	var to int
	switch {
	case len(timeout) == 1:
		to = timeout[0]
	default:
		to = 10
	}
	var errd error
	//doneMap := make(map[uint8]bool, 2)
	var done uint8
	doneChan := make(chan uint8)
	close(p.killChan)
	p.producer.AsyncClose()
	timer := time.NewTicker(time.Duration(to) * time.Second)

	go func() {
		for range p.producer.Successes() {
		}
		doneChan <- 1
	}()

	go func() {
		for range p.producer.Errors() {
		}
		doneChan <- 1
	}()

drainLoop:
	for {
		select {
		case <-timer.C:
			errd = fmt.Errorf("Timed out Draining Producer")
			break drainLoop
		case <-doneChan:
			done++
			if done == 2 {
				break drainLoop
			}
		}
	}
	if !p.noClose {
		p.cl.Close()
	}

	/*
		drainLoop:
			for {
				select {
				case <-timer.C:
					errd = fmt.Errorf("Timed out Draining Producer")
					break drainLoop
				case _, ok := <-p.producer.Successes():
					if !ok {
						doneMap[1] = true
					}
				case _, ok := <-p.producer.Errors():
					if !ok {
						doneMap[2] = true
					}
				default:
					if doneMap[1] && doneMap[2] {
						break drainLoop
					}
				}
			}
			if !p.noClose {
				p.cl.Close()
			}
	*/
	return errd
}

// Input is the input channel for the user to write messages to that they
// wish to send.
func (p *Producer) Input() chan<- *sarama.ProducerMessage {
	return p.input
}

// Successes is the success output channel back to the user when Return.Successes is
// enabled. If Return.Successes is true, you MUST read from this channel or the
// Producer will deadlock. It is suggested that you send and read messages
// together in a single select statement.
func (p *Producer) Successes() <-chan *sarama.ProducerMessage {
	return p.successes
}

// Errors is the error output channel back to the user. You MUST read from this
// channel or the Producer will deadlock when the channel is full. Alternatively,
// you can set Producer.Return.Errors in your config to false, which prevents
// errors to be returned.
func (p *Producer) Errors() <-chan *sarama.ProducerError {
	return p.errors
}

func (p *Producer) start() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for m := range p.producer.Errors() {
			p.errors <- m
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for m := range p.producer.Successes() {
			p.successes <- m
		}
	}()
inputLoop:
	for {
		select {
		case <-p.killChan:
			break inputLoop
		case m := <-p.input:
			p.producer.Input() <- m
		}
	}
	wg.Wait()
}
