package redis

import "errors"

var (
	ParamsErr             = errors.New("init redis params err")
	HasServerNameErr      = errors.New("redis server name has registered")
	GetRedisConnErr       = errors.New("get redis conn fail")
	NotExistServerNameErr = errors.New("redis server name not exist")
	RemoveNodeErr         = errors.New("remove redis node err")
)
