package command

import (
	m "github.com/GitHub121380/golib/middleware/gin"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/cron"
	"github.com/gin-gonic/gin/cycle"
	"github.com/gin-gonic/gin/job"
)

func InitCycle(g *gin.Engine) (c *cycle.Cycle) {
	c = cycle.New(g).AddBeforeRun(cycleBeforeRun).AddAfterRun(cycleAfterRun)
	return c
}

func cycleBeforeRun(ctx *gin.Context) bool {
	m.UseMetadata(ctx)
	m.LoggerBeforeRun(ctx)
	return true
}

func cycleAfterRun(ctx *gin.Context) {
	m.LoggerAfterRun(ctx)
}

func InitCrontab(g *gin.Engine) (c *cron.Cron) {
	c = cron.New(g).AddBeforeRun(cronBeforeRun).AddAfterRun(cronAfterRun)
	c.Start()
	return c
}

// 生成trace，塞入ctx
func cronBeforeRun(ctx *gin.Context) bool {
	m.UseMetadata(ctx)
	m.LoggerBeforeRun(ctx)
	return true
}

func cronAfterRun(ctx *gin.Context) {
	m.LoggerAfterRun(ctx)
}

func NewJob(g *gin.Engine) (c *job.Job) {
	c = job.New(g).AddBeforeRun(jobBeforeRun).AddAfterRun(jobAfterRun)
	return c
}

// 生成trace，塞入ctx
func jobBeforeRun(newCtx *gin.Context, tempParentSpanContext interface{}) bool {
	m.UseMetadata(newCtx)
	m.LoggerBeforeRun(newCtx)
	return true
}

func jobAfterRun(ctx *gin.Context) {
	m.LoggerAfterRun(ctx)
}
