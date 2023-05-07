// Package rmq 提供了访问RocketMQ服务的能力
package rmq

import (
	"context"
	"fmt"
	"github.com/GitHub121380/golib/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/GitHub121380/golib/zlog"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var (
	// ErrRmqSvcConfigInvalid 服务配置无效
	ErrRmqSvcConfigInvalid = fmt.Errorf("requested rmq service is not correctly configured")
	// ErrRmqSvcNotRegiestered 服务尚未被注册
	ErrRmqSvcNotRegiestered = fmt.Errorf("requested rmq service is not registered")
	// ErrRmqSvcInvalidOperation 当前操作无效
	ErrRmqSvcInvalidOperation = fmt.Errorf("requested rmq service is not suitable for current operation")
)

var (
	rmqServices   = make(map[string]*client)
	rmqServicesMu sync.Mutex
)

// MessageCallback 定义业务方接收消息的回调接口
type MessageCallback func(ctx *gin.Context, msg Message) error

func (config ClientConfig) checkConfig() error {
	if config.Group == "" {
		return ErrRmqSvcConfigInvalid
	}
	if config.Topic == "" {
		return ErrRmqSvcConfigInvalid
	}
	if len(config.NameServers) == 0 && config.NameServerZNS == "" {
		return ErrRmqSvcConfigInvalid
	}
	return nil
}

func InitRmq(service string, config ClientConfig) (err error) {
	if err = config.checkConfig(); err != nil {
		return err
	}

	clnt := &client{
		ClientConfig: &config,
	}
	rmqServicesMu.Lock()
	defer rmqServicesMu.Unlock()

	err = clnt.startNamingHandler()
	if err != nil {
		return err
	}

	rmqServices[service] = clnt
	return nil
}

// RegisterWithConfig 根据提供的配置结构体注册RocketMQ服务并命名，后续可使用该名称调用服务
func Register(service, configFile string) (err error) {
	rmqServicesMu.Lock()
	defer rmqServicesMu.Unlock()

	clnt := &client{
		ClientConfig: &ClientConfig{},
	}

	clnt.Conf, err = utils.Load(configFile, clnt.ClientConfig)

	f := []zap.Field{
		zap.String("prot", "rmq"),
		zap.String("service", service),
		zap.String("config", configFile),
	}
	if err != nil {
		zlog.ErrorLogger(nil, "load failed", f...)
		return err
	}
	if clnt.ClientConfig.Group == "" {
		zlog.ErrorLogger(nil, "group not specified", f...)
		return ErrRmqSvcConfigInvalid
	} else if clnt.ClientConfig.Topic == "" {
		zlog.ErrorLogger(nil, "topic not specified", f...)
		return ErrRmqSvcConfigInvalid
	} else if (len(clnt.ClientConfig.NameServers) == 0) &&
		(clnt.ClientConfig.NameServerZNS == "") {
		zlog.ErrorLogger(nil, "either name servers or name server zns must be provided", f...)
		return ErrRmqSvcConfigInvalid
	}

	err = clnt.startNamingHandler()
	if err != nil {
		return err
	}

	rmqServices[service] = clnt
	return nil
}

// Deprecated
func RegisterWithConfig(service string, config *ClientConfig) (err error) {
	return nil
}

// StartProducer 启动指定已注册的RocketMQ生产服务
func StartProducer(service string) error {
	if client, ok := rmqServices[service]; ok {
		client.mu.Lock()
		defer client.mu.Unlock()
		if client.producer != nil {
			return ErrRmqSvcInvalidOperation
		}
		var err error
		var nsDomain string
		nsDomain, err = client.getNameserverDomain()
		if err != nil {
			return err
		}
		client.producer, err = newProducer(
			client.ClientConfig.Auth.AccessKey, client.ClientConfig.Auth.SecretKey,
			service, client.ClientConfig.Group, nsDomain,
			client.ClientConfig.Retry, time.Duration(client.ClientConfig.Timeout)*time.Millisecond)
		if err != nil {
			return err
		}
		return client.producer.start()
	}
	return ErrRmqSvcNotRegiestered
}

// StopProducer 停止指定已注册的RocketMQ生产服务
func StopProducer(service string) error {
	if client, ok := rmqServices[service]; ok {
		client.mu.Lock()
		defer client.mu.Unlock()
		if client.producer == nil {
			return ErrRmqSvcInvalidOperation
		}
		err := client.producer.stop()
		client.producer = nil
		return err
	}
	return ErrRmqSvcNotRegiestered
}

// StartConsumer 启动指定已注册的RocketMQ消费服务， 同时指定要消费的消息标签，以及消费回调
func StartConsumer(g *gin.Engine, service string, tags []string, callback MessageCallback) error {
	if client, ok := rmqServices[service]; ok {
		client.mu.Lock()
		defer client.mu.Unlock()
		if client.pushConsumer != nil || callback == nil {
			return ErrRmqSvcInvalidOperation
		}
		var err error
		var nsDomain string
		nsDomain, err = client.getNameserverDomain()
		if err != nil {
			return err
		}

		cb := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, m := range msgs {
				if ctx.Err() != nil {
					zlog.WarnLogger(nil, "stop consume cause ctx cancelled", zap.String("prot", "rmq"))
					return consumer.SuspendCurrentQueueAMoment, ctx.Err()
				}
				ctx := gin.CreateNewContext(g)
				err := callback(ctx, &messageWrapper{
					msg:      &m.Message,
					offsetID: m.OffsetMsgId,
				})
				if err != nil {
					zlog.WarnLogger(ctx, "failed to consume message",
						zap.String("message", m.String()),
						zap.String("error", err.Error()))
					gin.RecycleContext(g, ctx)
					return consumer.SuspendCurrentQueueAMoment, nil
				}
				gin.RecycleContext(g, ctx)
			}
			return consumer.ConsumeSuccess, nil
		}

		client.pushConsumer, err = newPushConsumer(
			client.ClientConfig.Auth.AccessKey,
			client.ClientConfig.Auth.SecretKey,
			service,
			client.ClientConfig.Group,
			client.ClientConfig.Topic,
			client.ClientConfig.Broadcast,
			client.ClientConfig.Orderly,
			client.ClientConfig.Retry,
			tags,
			nsDomain,
			cb)
		if err != nil {
			return err
		}
		return client.pushConsumer.start()
	}
	return ErrRmqSvcNotRegiestered
}

// StopConsumer 停止指定已注册的RocketMQ消费服务
func StopConsumer(service string) error {
	if client, ok := rmqServices[service]; ok {
		client.mu.Lock()
		defer client.mu.Unlock()
		if client.pushConsumer == nil {
			return ErrRmqSvcInvalidOperation
		}
		err := client.pushConsumer.stop()
		client.pushConsumer = nil
		return err
	}
	return ErrRmqSvcNotRegiestered

}

// NewMessage 创建一条新的消息
func NewMessage(service string, content []byte) (Message, error) {
	if client, ok := rmqServices[service]; ok {
		return &messageWrapper{
			client: client,
			msg:    primitive.NewMessage(client.ClientConfig.Topic, content),
		}, nil
	}
	return nil, ErrRmqSvcNotRegiestered
}

var consumers []string

func Use(g *gin.Engine, service string, tags []string, handler MessageCallback) {
	if err := StartConsumer(g, service, tags, handler); err != nil {
		panic("Start consumer  error: " + err.Error())
	}
	consumers = append(consumers, service)
}

func StopRocketMqConsume() {
	for _, svc := range consumers {
		_ = StopConsumer(svc)
	}
}
