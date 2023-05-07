package pool

import "github.com/gin-gonic/gin"

type Pool interface {
	Get(ctx *gin.Context) (interface{}, error)

	Put(interface{}) error

	Close(interface{}) error

	Release()

	Len() int
}
