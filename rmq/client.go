package rmq

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/GitHub121380/golib/zlog"
	"github.com/GitHub121380/golib/zns"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

// auth 提供链接到Broker所需要的验证信息（按需配置）
type auth struct {
	AccessKey string `json:"ak,omitempty"`
	SecretKey string `json:"sk,omitempty"`
}

// ClientConfig 包含链接到RocketMQ服务所需要的各配置项
type ClientConfig struct {
	Service string `yaml:"service"`
	// 提供名字服务器的地址列表，例如: [ "127.0.0.1:9876" ]
	NameServers []string `json:"nameservers" yaml:"nameservers"`
	// 或者可以选择使用ZNS
	NameServerZNS string `json:"nameserverzns" yaml:"nameserverzns"`
	// 生产/消费者组名称，各业务线间需要保持唯一
	Group string `json:"group" yaml:"group"`
	// 要消费/订阅的主题
	Topic string `json:"topic" yaml:"topic"`
	// 如果配置了ACL，需提供验证信息
	Auth auth `json:"auth" yaml:"auth"`
	// 是否是广播消费模式
	Broadcast bool `json:"broadcast" yaml:"broadcast"`
	// 是否是顺序消费模式
	Orderly bool `json:"orderly" yaml:"orderly"`
	// 生产失败时的重试次数
	Retry int `json:"retry" yaml:"retry"`
	// 生产超时时间
	Timeout int `json:"timeout" yaml:"timeout"`
}

// Client 为客户端主体结构
type client struct {
	*ClientConfig

	mu sync.Mutex

	producer       *rmqProducer
	pushConsumer   *rmqPushConsumer
	namingListener net.Listener
}

func (c *client) startNamingHandler() error {
	var err error
	c.namingListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		zlog.ErrorLogger(nil, "failed to create naming listener")
		return err
	}
	go func() {
		err = http.Serve(c.namingListener, c.createNamingHandler())
		zlog.ErrorLogger(nil, "naming handler stopped", zap.String("error", err.Error()))
	}()
	return nil
}
func (c *client) createNamingHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if c.ClientConfig.NameServerZNS != "" {
			zlog.DebugLogger(nil, "try serve through zns", zap.String("zns", c.ClientConfig.NameServerZNS))
			res, err := zns.GetZnsInstance(c.ClientConfig.NameServerZNS)
			if err == nil && len(res) > 0 {
				for _, z := range res {

					zlog.DebugLogger(nil, "found naming",
						zap.String("zns", c.ClientConfig.NameServerZNS),
						zap.String("entry", fmt.Sprintf("%s:%d", z.IP, z.Port)))

					_, err = resp.Write([]byte(fmt.Sprintf("%v:%d;", z.IP, z.Port)))
					if err != nil {
						zlog.ErrorLogger(nil, "write response failed", zap.String("error", err.Error()))
					}
				}
				return
			}
			zlog.WarnLogger(nil, "get ns through zns failed",
				zap.String("zns", c.ClientConfig.NameServerZNS),
			)

		}
		if len(c.ClientConfig.NameServers) > 0 {
			zlog.WarnLogger(nil, "try serve through static config",
				zap.Strings("ns", c.ClientConfig.NameServers),
			)

			// no zns configured, fall back to ns in config file
			for _, ns := range c.ClientConfig.NameServers {
				_, err := resp.Write([]byte(fmt.Sprintf("%s;", ns)))
				if err != nil {
					zlog.WarnLogger(nil, "write response failed: "+err.Error())
				}
			}
			return
		}

		// no ns available
		resp.WriteHeader(http.StatusNotFound)
	}
}
func (c *client) getNameserverDomain() (string, error) {
	if c.namingListener != nil {
		return "http://" + c.namingListener.Addr().String(), nil
	}
	return "", ErrRmqSvcInvalidOperation
}

// DelayLevel 定义消息延迟发送的级别
type DelayLevel int

const (
	Second = DelayLevel(iota)
	Seconds5
	Seconds10
	Seconds30
	Minute1
	Minutes2
	Minutes3
	Minutes4
	Minutes5
	Minutes6
	Minutes7
	Minutes8
	Minutes9
	Minutes10
	Minutes20
	Minutes30
	Hour1
	Hours2
)

type rlogger struct {
	actualLogger *zap.Logger
	initOnce     sync.Once
	verbose      bool
}

func (r *rlogger) logger() *zap.Logger {
	r.initOnce.Do(func() {
		if os.Getenv("RMQ_SDK_VERBOSE") != "" {
			r.verbose = true
		} else {
			r.verbose = false
		}
		r.actualLogger = zlog.GetZapLogger().With(
			zap.String("prot", "rmq-sdk"))
	})
	return r.actualLogger
}

func (r *rlogger) withField(fields map[string]interface{}) *zap.Logger {
	l := r.logger()
	if r.verbose {
		var f []zap.Field
		for k, v := range fields {
			f = append(f, zap.Reflect(k, v))
		}
		l = l.With(f...)
	}
	return l
}
func (r *rlogger) Debug(msg string, fields map[string]interface{}) {
	if r.verbose {
		r.withField(fields).Debug(msg)
	}
}
func (r *rlogger) Info(msg string, fields map[string]interface{}) {
	if r.verbose {
		r.withField(fields).Info(msg)
	}
}
func (r *rlogger) Warning(msg string, fields map[string]interface{}) {
	if r.verbose {
		r.withField(fields).Warn(msg)
	}
}
func (r *rlogger) Error(msg string, fields map[string]interface{}) {
	r.withField(fields).Error(msg)
}
func (r *rlogger) Fatal(msg string, fields map[string]interface{}) {
	r.withField(fields).Fatal(msg)
}

func init() {
	rlog.SetLogger(&rlogger{})
}

// Message 消息提供的接口定义
type Message interface {
	WithTag(string) Message
	WithShard(string) Message
	WithDelay(DelayLevel) Message
	Send() (msgID string, err error)
	GetContent() []byte
	GetTag() string
	GetShard() string
	GetID() string
}

type messageWrapper struct {
	msg      *primitive.Message
	client   *client
	offsetID string
}

// WithTag 设置消息的标签Tag
func (m *messageWrapper) WithTag(tag string) Message {
	m.msg = m.msg.WithTag(tag)
	return m
}

// WithShard 设置消息的分片键
func (m *messageWrapper) WithShard(shard string) Message {
	m.msg = m.msg.WithShardingKey(shard)
	return m
}

// WithDelay 设置消息的延迟等级
func (m *messageWrapper) WithDelay(lvl DelayLevel) Message {
	m.msg = m.msg.WithDelayTimeLevel(int(lvl))
	return m
}

// Send 发送消息
func (m *messageWrapper) Send() (msgID string, err error) {
	if m.client == nil {
		zlog.WarnLogger(nil, "client is not specified")
		return "", ErrRmqSvcInvalidOperation
	}
	m.client.mu.Lock()
	prod := m.client.producer
	m.client.mu.Unlock()
	if prod == nil {
		zlog.WarnLogger(nil, "producer not started")
		return "", ErrRmqSvcInvalidOperation
	}
	queue, id, offset, err := m.client.producer.SendMessage(m.msg)
	if err != nil {
		zlog.ErrorLogger(nil, "failed to send message",
			zap.String("error", err.Error()),
			zap.String("message", m.msg.String()),
			zap.String("prot", "rmq"),
		)
		return "", err
	}

	fields := []zap.Field{
		zap.String("prot", "rmq"),
		zap.String("message", m.msg.String()),
		zap.String("queue", queue),
		zap.String("msgid", id),
		zap.String("offsetid", offset),
	}
	zlog.DebugLogger(nil, "sent message", fields...)

	return offset, nil
}

// GetContent 获取消息体内容
func (m *messageWrapper) GetContent() []byte {
	return m.msg.Body
}

// GetTag 获取消息标签
func (m *messageWrapper) GetTag() string {
	return m.msg.GetTags()
}

// GetShard 获取消息分片键
func (m *messageWrapper) GetShard() string {
	return m.msg.GetShardingKey()
}

// GetID 获取消息ID
func (m *messageWrapper) GetID() string {
	return m.offsetID
}
