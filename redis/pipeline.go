package redis

import (
	"errors"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type Pipeliner interface {
	Exec(ctx *gin.Context) ([]interface{}, error)
	Put(ctx *gin.Context, cmd string, args ...interface{}) error
}

type commands struct {
	cmd   string
	args  []interface{}
	reply interface{}
	err   error
}

type Pipeline struct {
	cmds  []commands
	err   error
	redis *Redis
}

func (objRedis *Redis) Pipeline() Pipeliner {
	return &Pipeline{
		redis: objRedis,
	}
}

func (p *Pipeline) Put(ctx *gin.Context, cmd string, args ...interface{}) error {
	if len(args) < 1 {
		return errors.New("no key found in args")
	}
	c := commands{
		cmd:  cmd,
		args: args,
	}
	p.cmds = append(p.cmds, c)

	field := []zap.Field{
		zap.String("prot", "redis"),
		zap.String("command", cmd),
		zap.String("commandVal", utils.JoinArgs(logForRedisValue, args)),
		zap.String("service", p.redis.r.Service),
	}

	zlog.InfoLogger(ctx, "pipeline put", field...)
	return nil
}

func (p *Pipeline) Exec(ctx *gin.Context) (res []interface{}, err error) {
	start := time.Now()
	zlog.CreateSpan(ctx)

	conn := p.redis.r.pool.Get()
	defer conn.Close()

	for i := range p.cmds {
		err = conn.Send(p.cmds[i].cmd, p.cmds[i].args...)
	}

	err = conn.Flush()

	var msg string
	if err == nil {
		for i := range p.cmds {
			var reply interface{}
			reply, err = conn.Receive()
			res = append(res, reply)
			p.cmds[i].reply, p.cmds[i].err = reply, err
		}

		msg = "pipeline exec succ"
	} else {
		p.err = err
		msg = "pipeline exec error: " + err.Error()
	}

	end := time.Now()

	field := []zap.Field{
		zap.String("prot", "redis"),
		zap.String("service", p.redis.r.Service),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
	}

	zlog.InfoLogger(ctx, msg, field...)

	return res, err
}
