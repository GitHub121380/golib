/**
 * @Author: shaoying@zuoyebang.com
 * @Description: 资源调度器
 * @File:  ral.go
 * @Version: 1.0.0
 * @Date: 2019-12-26 22:21
 */

package ral

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/GitHub121380/golib/env"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const ( // 资源类型
	TYPE_RAL   = "ral"
	TYPE_ZNS   = "zns"
	TYPE_NMQ   = "nmq"
	TYPE_HTTP  = "http"
	TYPE_REDIS = "redis"
	TYPE_MYSQL = "mysql"
	TYPE_HBASE = "hbase"
)

const ( // 选取方法
	WITH_RANDOM = "random"
	WITH_HASH   = "hash"
	WITH_ORDER  = "order"
	WITH_FIRST  = "first"
	WITH_LAST   = "last"
)
const (
	DATA_CONTEXT = "_context"
)

var ( // 返回错误
	ERR_NOT_FOUND_MODULE   = errors.New("not found module")
	ERR_NOT_FOUND_RESOURCE = errors.New("not found resource")
	ERR_NOT_FOUND_INSTANCE = errors.New("not found instance")
	ERR_NOT_FOUND_CLIENT   = errors.New("not found client")
	ERR_NOT_FOUND_METHOD   = errors.New("not found method")
	ERR_NOT_FOUND_ENCODE   = errors.New("not found encode")
	ERR_NOT_FOUND_LOGID    = errors.New("not found logid")
	ERR_NOT_FOUND_PATH     = errors.New("not found path")
	ERR_NOT_FOUND_CMD      = errors.New("not found cmd")
	ERR_NOT_FOUND_ARG      = errors.New("not found arg")
)

type Instance struct {
	IP   string
	Port int

	Client interface{}
	subs   []*Instance
	res    *Resource
}

func (ins *Instance) Request(ctx *gin.Context, method string, data map[string]interface{}, head map[string]string) (interface{}, error) {
	if ins == nil {
		return nil, ERR_NOT_FOUND_INSTANCE
	}

	res := ins.res
	if res == nil {
		return nil, ERR_NOT_FOUND_RESOURCE
	}

	mod := res.mod
	if mod == nil {
		return nil, ERR_NOT_FOUND_MODULE
	}

	// acl inject to header
	if head == nil {
		head = map[string]string{
			"SERVICE": env.AppName,
		}
	} else {
		head["SERVICE"] = env.AppName
	}
	// new span and inject header
	injectSpanContextToHeader(ctx, head)

	if len(mod.Methods) > 0 {
		if method, ok := mod.Methods[method]; ok {
			return method(ctx, res, ins, data, head)
		}
	}

	if mod.Method != nil {
		return mod.Method(ctx, res, ins, method, data, head)
	}
	return nil, ERR_NOT_FOUND_METHOD
}

func injectSpanContextToHeader(ctx *gin.Context, head map[string]string) {
	h := make(http.Header)
	for k, v := range head {
		h.Set(k, v)
	}

	if v, ok := h["Trace-Id"]; ok && len(v) > 0 {
		head["Trace-Id"] = v[0]
	}
	head["x_bd_spanid"] = zlog.CreateSpan(ctx)
	head["x_bd_logid"] = zlog.GetLogID(ctx)
}

func (ins *Instance) Release() {
	if ins != nil && ins.res != nil {
		if ins.res.list != nil {
			return
		}
		ins.res.count++
	}
}

// 返回retry值
func (ins *Instance) Retry(cb func(res *Resource, ins *Instance) bool) int {
	res := ins.res
	if res.Retry < 0 {
		res.Retry = 0
	}
	if res.Retry > 3 {
		res.Retry = 3
	}
	curIns := ins
	for i := 0; i < res.Retry+1; i++ {
		if !cb(res, curIns) || i == res.Retry {
			return i
		}

		if next, err := GetInstance(nil, res.Type, res.Name); err == nil {
			next.Release()
			curIns = next
		}
	}
	return res.Retry
}

type Resource struct {
	// 基本信息
	Type string
	Name string

	// 连接超时
	ConnTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	Retry        int

	// 编码解码
	Encode string
	Decode string

	// 是否开启长连接
	LongConnect bool

	// Redis配置
	Redis struct {
		MaxIdle     int
		MaxActive   int
		IdleTimeout time.Duration
		Password    string
		Database    int
		Wait        bool
	}

	// Mysql配置
	Mysql struct {
		Username        string
		Password        string
		Database        string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifeTime time.Duration
	}

	// hbase配置
	HBase struct {
		MaxIdle        int
		MaxActive      int
		IdleTimeout    time.Duration
		MaxWaitTimeout time.Duration
	}

	// 选取策略
	Strategy string

	// 服务域名
	ZNS struct {
		Name string
		IDC  map[string]string
	}
	// 服务地址
	Manual map[string][]struct {
		IP   string
		Host string
		Port int
	}

	Depend []string

	// 选取参数
	order int
	rand  *rand.Rand

	// 实例计数
	limit int
	total int
	count int

	// 实例列表
	list []*Instance
	busy []*Instance

	mod *Module

	lock sync.RWMutex
}

type Module struct {
	Type string

	// 添加资源
	Config func(*Resource) *Resource

	// 资源依赖
	Append func(self *Resource, res *Resource, ins *Instance) (*Instance, error)
	Remove func(self *Instance, res *Resource, ins *Instance) bool

	// 方法列表
	Methods map[string]func(*gin.Context, *Resource, *Instance, map[string]interface{}, map[string]string) (interface{}, error)
	Method  func(*gin.Context, *Resource, *Instance, string, map[string]interface{}, map[string]string) (interface{}, error)
}

var modules = map[string]*Module{}
var resources = map[string]*Resource{}
var dependons = map[string][]*Resource{}

var lock sync.RWMutex

func AddModule(mod *Module) *Module {
	lock.Lock()
	defer lock.Unlock()
	modules[mod.Type] = mod
	return mod
}
func AddResource(res *Resource) *Resource {
	if res, ok := GetResource(res.Type, res.Name); ok {
		return res
	}

	lock.Lock()
	defer lock.Unlock()
	zlog.InfoLogger(nil, fmt.Sprintf("ral add resource %s:%s", res.Type, res.Name), zap.String("prot", "ral"))
	resources[res.Type+":"+res.Name] = res

	res.mod = modules[res.Type]
	if res.mod != nil && res.mod.Config != nil {
		res.mod.Config(res)
	}

	if len(res.Manual) == 0 && res.ZNS.Name != "" {
		res.Depend = []string{"zns:" + res.ZNS.Name}
	}

	for _, dep := range res.Depend {
		dependons[dep] = append(dependons[dep], res)
		zlog.InfoLogger(nil, fmt.Sprintf("ral add dependon %s->%s:%s", dep, res.Type, res.Name), zap.String("prot", "ral"))
	}
	return res
}
func GetResource(modType string, name string) (*Resource, bool) {
	lock.RLock()
	defer lock.RUnlock()
	res, ok := resources[modType+":"+name]
	return res, ok
}

func AddInstance(modType string, name string, ins *Instance) (*Instance, error) {
	if ins == nil {
		return nil, nil
	}
	res, ok := GetResource(modType, name)
	if !ok {
		res = AddResource(&Resource{Type: modType, Name: name})
	}

	res.lock.Lock()
	defer res.lock.Unlock()

	zlog.InfoLogger(nil, fmt.Sprintf("ral add instance %s:%s %s:%d", modType, name, ins.IP, ins.Port), zap.String("prot", "ral"))
	res.list, ins.res = append(res.list, ins), res
	res.count++
	res.total++

	for _, v := range dependons[modType+":"+name] {
		if mod := modules[v.Type]; mod != nil && mod.Append != nil {
			if sub, err := mod.Append(v, res, ins); err != nil {
				return nil, err
			} else if sub != nil {
				i, err := AddInstance(v.Type, v.Name, sub)
				if err != nil {
					return nil, err
				}
				ins.subs = append(ins.subs, i)
			}
		}
	}
	return ins, nil
}
func DelInstance(Type string, Name string, ins *Instance) error {
	res, ok := GetResource(Type, Name)
	if !ok {
		return ERR_NOT_FOUND_RESOURCE
	}

	res.lock.Lock()
	defer res.lock.Unlock()

	zlog.InfoLogger(nil, fmt.Sprintf("ral del instance %s:%s %s:%d", Type, Name, ins.IP, ins.Port), zap.String("prot", "ral"))
	for i := 0; i < len(res.list); i++ {
		if res.list[i] == ins {
			for j := i; j < len(res.list)-1; j++ {
				res.list[j] = res.list[j+1]
			}
			res.list = res.list[:len(res.list)-1]
		}
	}
	res.count--
	res.total--

	for _, v := range ins.subs {
		if v != nil && v.res != nil {
			if mod := modules[v.res.Type]; mod != nil && mod.Append != nil {
				if mod.Remove(v, res, ins) {
					DelInstance(v.res.Type, v.res.Name, v)
				}
			}
		}
	}
	return nil
}

func GetInstance(ctx *gin.Context, modType string, name string) (*Instance, error) {
	fields := []zap.Field{
		zap.String("prot", "ral"),
	}

	res, ok := GetResource(modType, name)
	if !ok {
		zlog.WarnLogger(ctx, "ral GetInstance nil, modType: "+modType+"name: "+name, fields...)
		return nil, ERR_NOT_FOUND_RESOURCE
	}
	if res.count <= 0 {
		zlog.WarnLogger(ctx, "ral GetInstance count<=0 "+modType+"name: "+name, fields...)
		return nil, ERR_NOT_FOUND_INSTANCE
	}

	res.lock.RLock()
	defer res.lock.RUnlock()

	which := 0
	switch res.Strategy {
	case WITH_RANDOM:
		if res.rand == nil {
			res.rand = rand.New(rand.NewSource(time.Now().Unix()))
		}
		which = res.rand.Intn(res.count)
	case WITH_HASH:
	case WITH_FIRST:
		which = 0
	case WITH_LAST:
		if res.count > 1 {
			which = res.count - 1
		}
	case WITH_ORDER:
		which = res.order % res.count
		res.order++
	}

	if res.list != nil && which < len(res.list) {
		zlog.DebugLogger(ctx, fmt.Sprintf("ral get instance %s:%s %s:%d %s:%d", modType, name, res.Strategy, which, res.list[which].IP, res.list[which].Port), fields...)
		return res.list[which], nil
	}
	res.count--
	return &Instance{res: res}, nil
}
func GetAllInstance(Type string, Name string) []*Instance {
	if res, ok := GetResource(Type, Name); ok {
		return res.list
	}
	return nil
}

var mod *Module

func init() {
	mod = AddModule(&Module{Type: TYPE_RAL})
}

// 应用基础信息
var AppInfo struct {
	Agent string
	Name  string
	IDC   string
	IP    string
}

const AppAgent = "ral/golib 1.0"

func Init(dir string) (err error) {
	AppInfo.Agent = AppAgent
	AppInfo.Name = env.AppName
	AppInfo.IDC = env.IDC
	AppInfo.IP = env.LocalIP

	list, err := ioutil.ReadDir(dir)
	if err != nil {
		zlog.ErrorLogger(nil, "load config: "+err.Error(), zap.String("prot", "ral"))
		return err
	}

	for _, file := range list {
		resource := &Resource{}
		if _, e := utils.Load(path.Join(dir, file.Name()), &resource); e != nil {
			zlog.ErrorLogger(nil, "load config: "+e.Error(), zap.String("prot", "ral"))
			continue
		}
		AddResource(resource)
	}

	for _, res := range resources {
		if mod := res.mod; len(res.Manual) > 0 {
			for _, v := range res.Manual[env.IDC] {
				if mod != nil && mod.Append != nil {
					ins, err := mod.Append(res, res, &Instance{IP: v.IP, Port: v.Port})
					if err != nil {
						return err
					}
					if _, err = AddInstance(res.Type, res.Name, ins); err != nil {
						return err
					}
				} else {
					i := &Instance{
						IP:   v.IP,
						Port: v.Port,
					}
					if _, err = AddInstance(res.Type, res.Name, i); err != nil {
						return err
					}
				}
			}
		} else {
			var znsName string
			if res.ZNS.Name != "" {
				znsName = res.ZNS.Name
			} else if z, exist := res.ZNS.IDC[env.IDC]; exist {
				znsName = z
			}

			if znsName != "" {
				// todo: 由于值拷贝存在锁问题，直接初始化。Resource 变更需要同步修改此处
				r := &Resource{
					Type:         TYPE_ZNS,
					Name:         znsName,
					ConnTimeOut:  res.ConnTimeOut,
					ReadTimeOut:  res.ReadTimeOut,
					WriteTimeOut: res.WriteTimeOut,
					Retry:        res.Retry,
					Encode:       res.Encode,
					Decode:       res.Decode,
					Redis:        res.Redis,
					Mysql:        res.Mysql,
					HBase:        res.HBase,
					Strategy:     res.Strategy,
				}
				AddResource(r)
			}
		}
	}
	return nil
}
