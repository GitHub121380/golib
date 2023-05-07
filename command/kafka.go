package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GitHub121380/golib/base"
	m "github.com/GitHub121380/golib/middleware/gin"
	"github.com/GitHub121380/golib/zlog"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"runtime"
	"time"
)

const KafkaBodyKey string = "KafkaMsg"

type KafkaConsumeConfig struct {
	Service string   `yaml:"service"`
	Version string   `yaml:"version"`
	Brokers []string `yaml:"brokers"`
	Newest  bool
}

type KafkaSubClient struct {
	Brokers []string
	Version sarama.KafkaVersion
	g       *gin.Engine
}

func InitKafkaSub(g *gin.Engine, subConf KafkaConsumeConfig) *KafkaSubClient {
	v, err := sarama.ParseKafkaVersion(subConf.Version)
	if err != nil {
		panic("Error parsing Kafka version: " + err.Error())
	}

	return &KafkaSubClient{
		Version: v,
		g:       g,
		Brokers: subConf.Brokers,
	}
}

type kafkaHandler func(*gin.Context) error

type KafkaConsumerOption struct {
	ConsumerFromNewest bool
}

func (c *KafkaSubClient) AddSubFunction(topics []string, groupID string, handler kafkaHandler, opts *KafkaConsumerOption) {
	config := sarama.NewConfig()
	config.Version = c.Version

	if opts == nil || opts.ConsumerFromNewest == true {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	//config.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(c.Brokers, groupID, config)
	if err != nil {
		panic("NewConsumerGroup error: " + err.Error())
	}

	consumerHandler := &KafkaConsumerGroup{
		handler: handler,
		Client:  c,
	}

	ctx, _ := context.WithCancel(context.Background())
	consumerHandler.Ready = make(chan bool, 0)

	go func() {
		defer func() {
			if err := consumerGroup.Close(); err != nil {
				zlog.Warn(nil, "Error closing ConsumerGroupClient: ", err.Error())
			}
		}()
		if err := consumerGroup.Consume(ctx, topics, consumerHandler); err != nil {
			zlog.Warn(nil, "Error from consumer: ", err.Error())
		}

		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
	}()
	// Await till the consumer has been set up
	<-consumerHandler.Ready
	return
}

type KafkaConsumerGroup struct {
	Ready   chan bool
	handler func(*gin.Context) error
	Client  *KafkaSubClient
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *KafkaConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *KafkaConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *KafkaConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		err := c.HandleMessage(message)
		if err != nil {
			continue
		}
		session.MarkMessage(message, "")
	}

	return nil
}

func (c *KafkaConsumerGroup) HandleMessage(message *sarama.ConsumerMessage) error {
	ctx := gin.CreateNewContext(c.Client.g)
	customCtx := gin.CustomContext{
		Handle:    c.handler,
		Desc:      message.Topic,
		Type:      "Kafka",
		StartTime: time.Now(),
	}
	ctx.CustomContext = customCtx

	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			info, _ := json.Marshal(map[string]interface{}{
				"time":      time.Now().Format("2006-01-02 15:04:05"),
				"level":     "error",
				"module":    "stack",
				"requestId": zlog.GetRequestID(ctx),
				"handle":    ctx.CustomContext.HandlerName(),
			})
			fmt.Printf("%s\n-------------------stack-start-------------------\n%+v\n-------------------stack-end-------------------\n", string(info), r)
		}
		gin.RecycleContext(c.Client.g, ctx)
	}()

	var body base.KafkaBody
	if err := json.Unmarshal(message.Value, &body); err != nil {
		return err
	}
	ctx.Set(KafkaBodyKey, body.Msg)
	//m.LoggerBeforeRun(ctx)

	err := c.handler(ctx)

	ctx.CustomContext.Error = err
	ctx.CustomContext.EndTime = time.Now()
	m.LoggerAfterRun(ctx)
	return err
}

func GetKafkaMsg(ctx *gin.Context) (msg interface{}, exist bool) {
	msg, exist = ctx.Get(KafkaBodyKey)
	return msg, exist
}
