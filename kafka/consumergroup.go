package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// ConsumerGroup implements a sarama ConsumerGroup.
type ConsumerGroup struct {
	GroupID  string
	errors   chan error
	done     chan struct{}
	handlers CGHandler
	consumer sarama.ConsumerGroup
	ctx      context.Context
}

// GET configures the ConsumerGroup to process the given topic with the corresponding ProcessMessageFunc.
func (cg *ConsumerGroup) GET(topic string, handler ProcessMessageFunc) {
	cg.handlers.Handles[topic] = handler
}

// GETALL configures the ConsumerGroup to process all configured topics with the corresponding ProcessMessageFunc.
func (cg *ConsumerGroup) GETALL(handler ProcessMessageFunc) {
	for topic := range cg.handlers.Handles {
		cg.handlers.Handles[topic] = handler
	}
}

// Consume joins a cluster of consumers for a given list of topics and
// starts a blocking ConsumerGroupSession through the ConsumerGroupHandler.
func (cg *ConsumerGroup) Consume() error {
	var err error
	if len(cg.handlers.Handles) < 1 {
		panic("No topics configured for Consumer Group ")
	}
	var topics []string
	for t := range cg.handlers.Handles {
		topics = append(topics, t)
	}
ConsumeLoop:
	for {
		err = cg.consumer.Consume(cg.ctx, topics, cg.handlers)
		if err != nil {
			break ConsumeLoop
		}
		if cg.ctx.Err() != nil {
			break ConsumeLoop
		}
	}
	return err
}

// Errors returns a read channel of errors that occurred during the consumer life-cycle.
func (cg *ConsumerGroup) Errors() <-chan error {
	return cg.errors
}

// Close stops the ConsumerGroup and detaches any running sessions.
func (cg *ConsumerGroup) Close() error {
	close(cg.done)
	return cg.consumer.Close()
}

// ResumeAll resumes all partitions that have been paused.
func (cg *ConsumerGroup) ResumeAll() {
	cg.consumer.ResumeAll()
}

// PauseAll suspends fetching from all partitions.
func (cg *ConsumerGroup) PauseAll() {
	cg.consumer.PauseAll()
}

// CGHandler implements sarama ConsumerGroupHandlers.
type CGHandler struct {
	ctx               context.Context
	Handles           TopicHandlerMap
	consumer          sarama.ConsumerGroup
	errors            chan error
	done              chan struct{}
	autoCommitEnabled bool
}

func newCGHandler(ctx context.Context, cg sarama.ConsumerGroup, errChan chan error, done chan struct{}, autoCommitEnabled bool) CGHandler {
	return CGHandler{
		ctx:               ctx,
		Handles:           newTopicHandlerMap(),
		consumer:          cg,
		errors:            errChan,
		done:              done,
		autoCommitEnabled: autoCommitEnabled,
	}
}

// ProcessMessageFunc is used to iterate over Messages. Returning true to continue processing or false to stop and exit.
// Error will be sent to the Errors Channel for the Consumer, retrievable via Errors().
// The sarama.ConsumerGroupSession is available for direct access, which can be used for manual tasks.
// If sarama.Config.Consumer.Offsets.AutoCommit.Enable is false, be sure manually Mark Offsets and Commit.
// If sarama.Config.Consumer.Offsets.AutoCommit.Enable is true, offsets will be Marked and Committed automatically (obviously).
type ProcessMessageFunc func(sarama.ConsumerGroupSession, *sarama.ConsumerMessage) (bool, error)

// DefaultMessageFunc is default ProcessMessageFunc that can be used with a TopicHandlerMap.
func DefaultMessageFunc(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) (bool, error) {
	if msg != nil {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n%s\n\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		return true, nil
	}
	return false, fmt.Errorf("recieved nil message")
}

// TopicHandlerMap maps topics to ProcessMessageFuncs.
type TopicHandlerMap map[string]ProcessMessageFunc

func newTopicHandlerMap() TopicHandlerMap {
	return make(map[string]ProcessMessageFunc)
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h CGHandler) Setup(sess sarama.ConsumerGroupSession) error {
	go func() {
	errLoop:
		for {
			select {
			case <-h.done:
				break errLoop
			case <-h.ctx.Done():
				break errLoop
			case err := <-h.consumer.Errors():
				h.errors <- err
			}
		}
	}()
	// Possibly auto-commit anyways?
	/*
		if !h.autoCommitEnabled {
			go func() {
				ticker := time.NewTicker(time.Second * 2)
			commitLoop:
				for {
					select {
					case <-h.ctx.Done():
						sess.Commit()
						break commitLoop
					case <-ticker.C:
						sess.Commit()
					}
				}
			}()
		}
	*/
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (h CGHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	/* For Future TroubleShooting or Debugging options:
	for t, parts := range sess.Claims() {
		for _, p := range parts {
			fmt.Println(h.cg.Metadata.Members[t][p], "> Topic:", t, ">", p)
		}
	}
	*/
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h CGHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var (
		good bool
		err  error
	)
	handleFunc, ok := h.Handles[claim.Topic()]
	if h.autoCommitEnabled {
		handleFunc = h.makeAutoProcessMessageFunc(handleFunc)
	}
	if ok {
	ProcessMessages:
		for {
			select {
			case <-h.ctx.Done():
				break ProcessMessages
			case msg := <-claim.Messages():
				good, err = handleFunc(sess, msg)
				if err != nil {
					h.errors <- err
				}
				if !good {
					break ProcessMessages
				}
			}
		}
	} else {
		panic("No handler for given topic " + claim.Topic())
	}
	return err
}

func (h CGHandler) makeAutoProcessMessageFunc(fn ProcessMessageFunc) ProcessMessageFunc {
	return func(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) (good bool, err error) {
		good, err = fn(sess, msg)
		if good {
			sess.MarkMessage(msg, "")
		}
		return
	}
}

// NewConsumerGroup returns a new ConsumerGroup.
func NewConsumerGroup(ctx context.Context, addrs []string, groupID string, config *sarama.Config, topics ...string) (*ConsumerGroup, error) {
	if config == nil {
		config = GetConf("")
		config.Version = RecKafkaVersion
		config.Consumer.Return.Errors = true
		config.Consumer.Group.Session.Timeout = time.Second * 10
		config.Consumer.Group.Heartbeat.Interval = time.Second * 3
		config.Consumer.Group.Rebalance.Retry.Backoff = time.Second * 2
		config.Consumer.Group.Rebalance.Retry.Max = 3
		config.Consumer.MaxProcessingTime = time.Millisecond * 500
		config.Consumer.Offsets.AutoCommit.Enable = true
	}
	group, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		return nil, err
	}
	cg := &ConsumerGroup{
		GroupID:  groupID,
		consumer: group,
		errors:   make(chan error, config.ChannelBufferSize),
		done:     make(chan struct{}),
		ctx:      ctx,
	}
	cg.handlers = newCGHandler(ctx, group, cg.errors, cg.done, config.Consumer.Offsets.AutoCommit.Enable)
	if len(topics) > 0 {
		for _, topic := range topics {
			cg.handlers.Handles[topic] = DefaultMessageFunc
		}
	}
	return cg, nil
}

// NewConsumerGroup returns a new ConsumerGroup using an existing KClient.
func (kc *KClient) NewConsumerGroup(ctx context.Context, groupID string, topics ...string) (*ConsumerGroup, error) {
	/*	Consumer Group Options to be aware of:
		kc.config.Consumer.Return.Errors = true
		kc.config.Consumer.Group.Session.Timeout = time.Second * 10
		kc.config.Consumer.Group.Heartbeat.Interval = time.Second * 3
		kc.config.Consumer.Group.Rebalance.Retry.Backoff = time.Second * 2
		kc.config.Consumer.Group.Rebalance.Retry.Max = 3
		kc.config.Consumer.MaxProcessingTime = time.Millisecond * 500
	*/
	group, err := sarama.NewConsumerGroupFromClient(groupID, kc.cl)
	if err != nil {
		return nil, err
	}
	cg := &ConsumerGroup{
		GroupID:  groupID,
		consumer: group,
		errors:   make(chan error, kc.config.ChannelBufferSize),
		done:     make(chan struct{}),
		ctx:      ctx,
	}
	cg.handlers = newCGHandler(ctx, group, cg.errors, cg.done, kc.config.Consumer.Offsets.AutoCommit.Enable)
	if len(topics) > 0 {
		for _, topic := range topics {
			cg.handlers.Handles[topic] = DefaultMessageFunc
		}
	}
	return cg, nil
}
