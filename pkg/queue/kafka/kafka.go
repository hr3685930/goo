package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang-module/carbon"
	"goo/pkg/queue"
	"reflect"
	"strings"
)

type Kafka struct {
	cli            sarama.Client
	Brokers        []string
	ConsumerTopics []string
	ProducerTopic  string
	Prefix         string
}

func NewKafka(urls, prefix string) queue.Queue {
	brokers := strings.Split(urls, ",")
	return &Kafka{Prefix: prefix, Brokers: brokers}
}

type Msg struct {
	GroupId string
	Message queue.JobBase
}

type consumerGroupHandler struct {
	k       *Kafka
	jobName string
	GroupId string
	Message queue.JobBase
}

// connect
func (k *Kafka) Connect() error {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	client, err := sarama.NewClient(k.Brokers, config)
    if err != nil {
        return err
    }
	k.cli = client
	return nil
}

// producer connect
func (k *Kafka) ProducerConnect() queue.Queue {
	return &Kafka{cli: k.cli, Prefix: k.Prefix, ProducerTopic: k.ProducerTopic, ConsumerTopics: k.ConsumerTopics, Brokers: k.Brokers}
}

// consumer connect
func (k *Kafka) ConsumerConnect() queue.Queue {
	return &Kafka{cli: k.cli, Prefix: k.Prefix, ProducerTopic: k.ProducerTopic, ConsumerTopics: k.ConsumerTopics, Brokers: k.Brokers}
}

// topic
func (k *Kafka) Topic(topic string) {
	k.ProducerTopic = topic
	k.ConsumerTopics = []string{topic}
}

//error handler
func (k *Kafka) failOnErr(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func (k *Kafka) Producer(job queue.JobBase) {
	p, err := sarama.NewSyncProducerFromClient(k.cli)
	k.failOnErr(err, "producer create producer error")

	msg := &sarama.ProducerMessage{
		Topic: k.ProducerTopic,
		Key:   nil,
	}

	jName := reflect.TypeOf(job).String()
	groupID := strings.ToLower(strings.Replace(k.Prefix+"."+jName[1:], ".", "_", -1))
	jsonMsg := &Msg{groupID, job}
	message, err := json.Marshal(jsonMsg)
	k.failOnErr(err, "Umarshal failed")
	msg.Value = sarama.ByteEncoder(message)
	_, _, err = p.SendMessage(msg)
	k.failOnErr(err, "send error")
	_ = p.Close()
}

// no sleep retry
func (k *Kafka) Consumer(job queue.JobBase, sleep, retry int32) {
	jName := reflect.TypeOf(job).String()
	jobName := jName[1:]
	groupID := strings.ToLower(strings.Replace(k.Prefix+"."+jobName, ".", "_", -1))
	group, err := sarama.NewConsumerGroupFromClient(groupID, k.cli)
	k.failOnErr(err, "Consumer group err")
	ctx := context.Background()
	for { // 避免被挤掉
		topics := k.ConsumerTopics
		handler := &consumerGroupHandler{k: k, jobName: groupID, Message: job}
		err := group.Consume(ctx, topics, handler)
		k.failOnErr(err, "Consumer err")
	}
}

// report
func (k *Kafka) Err(failed queue.FailedJobs) {
	queue.ErrJob <- failed
}

func (k *Kafka) Close() {
	_ = k.cli.Close()
}

func (k *Kafka) ExportErr(err error, msg, groupID string) {
	e := err.(*queue.Error)
	k.Err(queue.FailedJobs{
		Connection: "kafka",
		Queue:      groupID,
		Message:    msg,
		Exception:  err.Error(),
		Stack:      e.Stack(),
		FiledAt:    carbon.Now(),
	})
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := json.Unmarshal(msg.Value, c)
		if err != nil {
            sess.MarkMessage(msg, "")
		}

		if c.GroupId == c.jobName {
			err = c.Message.Handler()
			queueErr := err.(*queue.Error)
			if queueErr != nil {
				c.k.ExportErr(queue.Err(err), string(msg.Value), c.GroupId)
			}
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
