package zns

import (
	"fmt"
	"github.com/GitHub121380/golib/ral"
	"github.com/GitHub121380/golib/zlog"
	"github.com/GitHub121380/golib/zns/bns"
	"github.com/GitHub121380/golib/zns/util"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"time"
)

type Instance struct {
	IP   string
	Port int
}

// 根据 znsName 查询有效ip:port
func GetZnsInstance(name string) (list []Instance, err error) {
	insInfoList, err := Resolve(name)
	if err != nil {
		return nil, err
	}
	for _, value := range insInfoList {
		if *value.InstanceStatus.Status != 0 {
			continue
		}
		in := Instance{
			IP:   util.UInt32IpToString(*value.HostIp),
			Port: int(*value.InstanceStatus.Port),
		}
		list = append(list, in)
	}
	return list, nil
}

func Resolve(name string) ([]*bns.InstanceInfo, error) {
	bnsClient := bns.New(config.LocalAddr, config.Timeout)
	if err := bnsClient.Connect(); err != nil {
		zlog.ErrorLogger(nil, "zns connect error:"+err.Error(), zap.String("prot", "zns"))
		return nil, err
	}
	defer bnsClient.Close()

	req := &bns.LocalNamingRequest{
		ServiceName: proto.String(name),
		All:         proto.Bool(true),
		Type:        proto.Int32(0),
	}
	var rsp bns.LocalNamingResponse
	if err := bnsClient.Call(req, &rsp); err != nil {
		zlog.ErrorLogger(nil, "call zns error: "+err.Error()+" name:"+name, zap.String("prot", "zns"))
		return nil, err
	}

	return rsp.InstanceInfo, nil
}

func Download(list []*ral.Resource) {
	// 连接服务器
	bnsClient := bns.New(config.LocalAddr, config.Timeout)
	if err := bnsClient.Connect(); err != nil {
		zlog.ErrorLogger(nil, "zns connect error:"+err.Error(), zap.String("prot", "zns"))
		return
	}
	defer bnsClient.Close()

	for _, s := range list {
		fields := []zap.Field{
			zap.String("prot", "zns"),
			zap.String("service", s.Name),
		}
		req := &bns.LocalNamingRequest{
			ServiceName: proto.String(s.Name),
			All:         proto.Bool(true),
			Type:        proto.Int32(0),
		}
		rsp := &bns.LocalNamingResponse{}

		// 下载主机列表
		if err := bnsClient.Call(req, rsp); err != nil {
			zlog.ErrorLogger(nil, "call zns error: "+err.Error(), fields...)
			continue
		}

		if len(rsp.InstanceInfo) == 0 {
			// 如果本次没有拿到结果，不做处理，使用上次获得的数据
			zlog.WarnLogger(nil, "call zns get num: 0 , pls manual check. ", fields...)
			continue
		}

		// 查询主机列表
		list := map[string]*ral.Instance{}
		for _, v := range ral.GetAllInstance(ral.TYPE_ZNS, s.Name) {
			list[fmt.Sprintf("%s:%d", v.IP, v.Port)] = v
		}

		for _, value := range rsp.InstanceInfo {
			if *value.InstanceStatus.Status == 0 {
				ip := util.UInt32IpToString(*value.HostIp)
				port := *value.InstanceStatus.Port

				if node := fmt.Sprintf("%s:%d", ip, port); list[node] == nil {
					// 添加主机地址
					zlog.InfoLogger(nil, "znsService add "+s.Name+" "+node, fields...)
					if _, err := ral.AddInstance(ral.TYPE_ZNS, s.Name, &ral.Instance{IP: ip, Port: int(port)}); err != nil {
						zlog.ErrorLogger(nil, "znsService add "+s.Name+""+node+" error: "+err.Error(), fields...)
					}
				} else {
					delete(list, node)
				}
			}
		}

		// 删除多余主机
		for k, v := range list {
			zlog.InfoLogger(nil, "znsService del "+s.Name+" "+k, fields...)
			if err := ral.DelInstance(ral.TYPE_ZNS, s.Name, v); err != nil {
				zlog.ErrorLogger(nil, "ral.DelInstance error:"+err.Error(), fields...)
			}
		}
	}
	return
}

func Refresh() {
	defer func() {
		if err := recover(); err != nil {
			zlog.ErrorLogger(nil, fmt.Sprintf("zns download panic, error: %+v", err))
			Refresh()
		}
	}()

	for {
		Download(resList)
		time.Sleep(config.Interval)
	}
}

var resList []*ral.Resource

var mod *ral.Module

func init() {
	mod = ral.AddModule(&ral.Module{Type: ral.TYPE_ZNS, Config: func(res *ral.Resource) *ral.Resource {
		resList = append(resList, res)
		return res
	}})
}

type Config struct {
	Timeout    time.Duration `yaml:"timeout"`
	Interval   time.Duration `yaml:"interval"`
	LocalAddr  string        `yaml:"localAddr"`
	RemoteAddr string        `yaml:"remoteAddr"`
}

var config Config

func (conf *Config) checkConf() {
	if conf.Timeout == 0 {
		conf.Timeout = 5 * time.Second
	}

	if conf.Interval == 0 {
		conf.Interval = 10 * time.Second
	}

	if conf.LocalAddr == "" {
		conf.LocalAddr = "localhost:793"
		conf.RemoteAddr = "localhost:793"
	} else {
		conf.RemoteAddr = conf.LocalAddr
	}
}

func Init(conf Config) {
	conf.checkConf()
	config = conf

	go Refresh()
}
