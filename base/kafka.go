package base

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/GitHub121380/golib/env"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/zlog"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

type KafkaProducerConfig struct {
	Service string `yaml:"service"`
	Addr    string `yaml:"addr"`
	Version string `yaml:"version"`

	SASL struct {
		Enable    bool   `yaml:"enable"`
		Handshake bool   `yaml:"handshake"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
	} `yaml:"sasl"`

	TLS struct {
		Enable                bool   `yaml:"enable"`
		CA                    string `yaml:"ca"`
		Cert                  string `yaml:"cert"`
		Key                   string `yaml:"key"`
		InsecureSkipTLSVerify bool   `yaml:"insecure_skip_tls_verify"`
	} `yaml:"tls"`
}
type KafkaPubClient struct {
	Conf     KafkaProducerConfig
	producer sarama.SyncProducer
}

type KafkaBody struct {
	Msg interface{}
}

func (conf *KafkaProducerConfig) GetKafkaConfig() (*sarama.Config, error) {
	defaultConfig := sarama.NewConfig()
	v, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		return nil, err
	}
	defaultConfig.Version = v
	if conf.SASL.Enable {
		defaultConfig.Net.SASL.Enable = true
		defaultConfig.Net.SASL.Handshake = conf.SASL.Handshake
		defaultConfig.Net.SASL.User = conf.SASL.User
		defaultConfig.Net.SASL.Password = conf.SASL.Password
	}
	if conf.TLS.Enable {
		defaultConfig.Net.TLS.Enable = true
		defaultConfig.Net.TLS.Config = &tls.Config{
			RootCAs:            x509.NewCertPool(),
			InsecureSkipVerify: conf.TLS.InsecureSkipTLSVerify,
		}
		if conf.TLS.CA != "" {
			ca, err := ioutil.ReadFile(conf.TLS.CA)
			if err != nil {
				panic("kafka pub CA error: %v" + err.Error())
			}
			defaultConfig.Net.TLS.Config.RootCAs.AppendCertsFromPEM(ca)
		}
	}
	defaultConfig.Producer.Return.Successes = true

	return defaultConfig, nil
}

func InitKafkaPub(conf KafkaProducerConfig) *KafkaPubClient {
	saramaConfig, err := conf.GetKafkaConfig()
	if err != nil {
		panic("kafka pub version error: %v" + err.Error())
	}

	producer, err := sarama.NewSyncProducer([]string{conf.Addr}, saramaConfig)
	if err != nil {
		panic("kafka pub new producer error: %v" + err.Error())
	}

	c := &KafkaPubClient{
		Conf:     conf,
		producer: producer,
	}
	return c
}

func (client *KafkaPubClient) CloseProducer() error {
	if client.producer != nil {
		return client.producer.Close()
	}
	return nil
}

func (client *KafkaPubClient) Pub(ctx *gin.Context, topic string, msg interface{}) error {
	if client.producer == nil {
		return errors.New("kafka producer not init")
	}
	// todo 可加业务相关
	kafkaBody := KafkaBody{
		Msg: msg,
	}
	body, err := json.Marshal(kafkaBody)
	if err != nil {
		return err
	}

	start := time.Now()
	kafkaMsg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(body)}
	partition, offset, err := client.producer.SendMessage(kafkaMsg)
	end := time.Now()

	ralCode := 0
	infoMsg := "kafka pub success"
	if err != nil {
		ralCode = -1
		infoMsg = err.Error()
	}

	fields := []zap.Field{
		zap.String("requestId", zlog.GetRequestID(ctx)),
		zap.String("localIp", env.LocalIP),
		zap.String("remoteAddr", client.Conf.Addr),
		zap.String("service", client.Conf.Service),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.Int("ralCode", ralCode),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
	}

	zlog.InfoLogger(nil, infoMsg, fields...)

	return nil
}
