package cos

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	tcos "github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"time"
)

// XCosACL 相关属性
const BucketAclPublicReadWrite = "public-read-write"
const BucketAclPublicRead = "public-read"
const BucketAclPrivate = "private"

type Option struct {
	AppID     string `yaml:"app_id"`
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

func (conf *Option) checkConf() {
}

// 腾讯云对象存储服务
type Client struct {
	AppID     string
	SecretID  string
	SecretKey string
	conn      *http.Client
}

func NewTcos(o *Option) *Client {
	o.checkConf()
	c := &http.Client{
		Transport: &tcos.AuthorizationTransport{
			SecretID:  o.SecretID,
			SecretKey: o.SecretKey,
			Expire:    time.Hour,
			Transport: &cosRequestTransport{},
		},
	}
	return &Client{
		AppID:     o.AppID,
		SecretID:  o.SecretID,
		SecretKey: o.SecretKey,
		conn:      c,
	}
}

type Bucket struct {
	// bucket 基础信息
	Name   string
	Region string

	// tcos client 基础信息
	conn      *http.Client
	AppID     string
	SecretID  string
	SecretKey string

	// bucket 个性化信息
	PictureRegion string
	FileSizeLimit int
	Thumbnail     int
	Directory     string
	FilePrefix    string
	IsNeedProxy   bool
}

type BucketOption struct {
	PictureRegion string
	FileSizeLimit int
	Thumbnail     int
	Directory     string
	FilePrefix    string
	IsNeedProxy   bool
}

func (c *Client) Bucket(name, region string, option *BucketOption) *Bucket {
	b := &Bucket{
		Name:      name,
		Region:    region,
		conn:      c.conn,
		AppID:     c.AppID,
		SecretID:  c.SecretID,
		SecretKey: c.SecretKey,
	}
	return b
}

func (c *Client) GetBucketList(ctx *gin.Context) (b []tcos.Bucket, err error) {
	u, _ := c.getServiceUrl()
	basicUrl := &tcos.BaseURL{
		ServiceURL: u,
	}
	client := tcos.NewClient(basicUrl, c.conn)
	s, _, err := client.Service.Get(context.Background())
	if err != nil {
		return b, err
	}

	return s.Buckets, nil
}

func (c *Client) getServiceUrl() (*url.URL, error) {
	return url.Parse("http://service.cos.myqcloud.com")
}

// 创建bucket
func (c *Client) CreateBucket(ctx *gin.Context, bucket, region string, opt *tcos.BucketPutOptions) (err error) {
	u, err := getBucketUrl(bucket, c.AppID, region)
	if err != nil {
		return err
	}
	basicUrl := &tcos.BaseURL{
		BucketURL: u,
	}

	cli := tcos.NewClient(basicUrl, c.conn)
	// 创建存储桶
	_, err = cli.Bucket.Put(context.Background(), opt)
	if err != nil {
		return err
	}
	return nil
}

func getBucketUrl(bucket, appID, region string) (*url.URL, error) {
	bucketUrl := fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", bucket, appID, region)
	return url.Parse(bucketUrl)
}

// bucket 相关配置
type BucketConfig struct {
	Bucket        string `yaml:"bucket"`
	AppID         string `yaml:"app_id"`
	SecretID      string `yaml:"secret_id"`
	SecretKey     string `yaml:"secret_key"`
	Region        string `yaml:"region"`
	PictureRegion string `yaml:"picture_region"`
	FileSizeLimit int    `yaml:"filesize_limit"`
	Thumbnail     int    `yaml:"thumbnail"`
	Directory     string `yaml:"directory"`
	FilePrefix    string `yaml:"file_prefix"`

	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	IdleConnTimeout     time.Duration
	ConnectTimeout      time.Duration
}

func (conf *BucketConfig) checkConf() {
	// TODO: 后续支持用户配置连接池相关信息，暂用默认的
	conf.MaxIdleConns = 300
	conf.MaxIdleConnsPerHost = 100
	conf.IdleConnTimeout = 30 * time.Second
}

// 根据php版本cos直接初始化bucket操作对象
func NewBucket(cfg BucketConfig) (b *Bucket) {
	cfg.checkConf()
	c := &http.Client{
		Transport: &tcos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
			Expire:    time.Hour,
			Transport: &cosRequestTransport{
				Transport: &http.Transport{
					MaxIdleConns:        cfg.MaxIdleConns,
					MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
					IdleConnTimeout:     cfg.IdleConnTimeout,
				},
			},
		},
	}

	return &Bucket{
		Name:          cfg.Bucket,
		Region:        cfg.Region,
		conn:          c,
		PictureRegion: cfg.PictureRegion,
		FileSizeLimit: cfg.FileSizeLimit,
		Thumbnail:     cfg.Thumbnail,
		Directory:     cfg.Directory,
		FilePrefix:    cfg.FilePrefix,
		AppID:         cfg.AppID,
		SecretID:      cfg.SecretID,
		SecretKey:     cfg.SecretKey,
	}
}
