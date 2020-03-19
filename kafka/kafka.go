package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gladmo/kafka2mongo/conf"
	"github.com/gladmo/kafka2mongo/mongodb"
)

type Msg struct {
	Timestamp      time.Time   `json:"timestamp"`
	BlockTimestamp time.Time   `json:"block_timestamp"`
	Key            interface{} `json:"key"`
	Value          interface{} `json:"value"`
	Topic          string      `json:"topic"`
	Partition      int32       `json:"partition"`
	Offset         int64       `json:"offset"`
	CreateAt       time.Time   `json:"create_at"`
}

func byte2interface(b []byte) (i interface{}, err error) {
	err = bson.UnmarshalExtJSON(b, true, &i)
	if err != nil {
		return
	}

	return
}

type allTopic struct{}

func (allTopic) Setup(_ ConsumerGroupSession) error   { return nil }
func (allTopic) Cleanup(_ ConsumerGroupSession) error { return nil }
func (h allTopic) ConsumeClaim(sess ConsumerGroupSession, claim ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		key, err := byte2interface(msg.Key)
		if err != nil {
			key = string(msg.Key)
		}

		value, err := byte2interface(msg.Value)
		if err != nil {
			value = string(msg.Value)
		}

		doc := Msg{
			Timestamp:      msg.Timestamp,
			BlockTimestamp: msg.BlockTimestamp,
			Key:            key,
			Value:          value,
			Topic:          msg.Topic,
			Partition:      msg.Partition,
			Offset:         msg.Offset,
			CreateAt:       time.Now(),
		}

		err = mongodb.Store(&mongodb.Document{
			Database:   "kafka2mongo",
			Collection: "all-topics",
			Doc:        doc,
		})

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func Run() {
	// Init config, specify appropriate version
	config := NewConfig()
	config.Producer.Compression = CompressionSnappy
	config.Producer.MaxMessageBytes = 4000000
	config.Producer.Return.Successes = true
	config.Version = V2_2_0_0

	address := conf.Getenv("KAFKA_ADDR", "localhost:9092")
	addr := strings.Split(strings.Trim(address, ","), ",")

	// Start with a client
	client, err := NewClient(addr, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	// Start a new consumer group
	group, err := NewConsumerGroupFromClient("kafka2mongo-group", client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	allTopics, err := client.Topics()
	if err != nil {
		panic(err)
	}
	excludeTopics := conf.Getenv("EXCLUDE_TOPICS", "")

	topics := ExcludeTopic(allTopics, excludeTopics)
	fmt.Println("allTopics:", allTopics)
	fmt.Println("excludeTopics:", excludeTopics)
	fmt.Println("useTopics:", topics)
	fmt.Println("Brokers addr:", client.Brokers()[0].Addr())
	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		err = group.Consume(ctx, topics, allTopic{})
		if err != nil {
			panic(err)
		}
	}
}

func ExcludeTopic(topics []string, excludeTopics string) (out []string) {
	// Exclude topic
	excludeTopicList := strings.Split(strings.Trim(excludeTopics, ","), ",")

	var exclude = make(map[string]struct{})
	for _, excludeTopic := range excludeTopicList {
		exclude[excludeTopic] = struct{}{}
	}

	for _, topic := range topics {
		if _, ok := exclude[topic]; ok {
			continue
		}

		out = append(out, topic)
	}

	return
}
