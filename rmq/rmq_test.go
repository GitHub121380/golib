package rmq_test

import (
	"fmt"
	"math/rand"
	"sync/atomic"

	"github.com/GitHub121380/golib/rmq"
)

var service = "rmqtest"
var serviceConf = "service/rmqtest.json"
var count = uint64(100)

func Example_producer() {
	// 加载服务topic、集群地址等信息
	err := rmq.Register(service, serviceConf)
	if err != nil {
		fmt.Println("register rmqtest service failed", err)
		return
	}
	fmt.Println("registered rmqtest service")
	// 启动生产者
	err = rmq.StartProducer(service)
	if err != nil {
		fmt.Println("start rmqtest producer failed", err)
		return
	}
	fmt.Println("started rmqtest producer")

	// 创建新消息
	msg, err := rmq.NewMessage(service, []byte("this is a message for test"))
	if err != nil {
		fmt.Println("create rmqtest message failed", err)
		return
	}

	// 可选设置Tag、分区、延迟，并进行消息发送
	id, err := msg.WithShard(fmt.Sprintf("%d", rand.Int63()%10)).
		WithTag(fmt.Sprintf("%d", rand.Int63()%10)).WithDelay(rmq.Seconds30).Send()
	if err != nil {
		fmt.Println("send rmqtest message failed", err)
		return
	}
	fmt.Println("sent message", "id", id, "shard", shard, "tag", tag)

	// 停止生产者
	err = rmq.StopProducer(service)
	if err != nil {
		fmt.Println("stop rmqtest producer failed", err)
		return
	}
	fmt.Println("stopped rmqtest producer")
}

func Example_consumer() {
	// 加载服务topic、集群地址等信息
	err := rmq.Register(service, serviceConf)
	if err != nil {
		fmt.Println("register rmqtest service failed", err)
		return
	}
	fmt.Println("registered rmqtest service")

	// 启动消费者，使用回调方式
	err = rmq.StartConsumer(service, nil, func(msg rmq.Message) error {
		fmt.Println("got message", "id", msg.GetID(), "shard", msg.GetShard(), "tag", msg.GetTag())
		atomic.AddUint64(&count, 1)
		return nil
	})
	if err != nil {
		fmt.Println("start rmqtest consumer failed", err)
		return
	}
	fmt.Println("started rmqtest consumer")

	// ... 其他业务逻辑

	// 停止消费者
	err = rmq.StopConsumer(service)
	if err != nil {
		fmt.Println("stop rmqtest consumer failed", err)
		return
	}
	fmt.Println("stopped rmqtest consumer", count)
}
