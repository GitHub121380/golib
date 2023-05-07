package cos

import (
	"encoding/json"
	"fmt"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"regexp"
)

type ImageMeta struct {
	Format   string `json:"format"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Size     string `json:"size"`
	Md5      string `json:"md5"`
	PhotoRgb string `json:"photoRgb"`

	// 兼容php，回填数据部分
	Bucket    string `json:"bucket"`
	Pid       string `json:"pid"`
	SourceUrl string `json:"sourceUrl"`
}

// 获取图片基本信息
func (b *Bucket) GetImageMeta(ctx *gin.Context, pid, fileType string) (m ImageMeta, err error) {
	ur := b.GetImageUrlByPid(ctx, pid, fileType)
	ur += "?imageInfo"
	req, err := http.NewRequestWithContext(ctx, "GET", ur, nil)
	if err != nil {
		return m, err
	}
	req.Header.Set("User-Agent", "cos-php-sdk-v4.3.7")

	resp, err := b.conn.Do(req)
	if err != nil {
		return m, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return m, err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &m); err != nil {
		zlog.Warnf(ctx, "GetImageMeta response not valid: %s", string(body))
		return m, err
	}

	m.Bucket = b.Name
	m.Pid = pid
	m.SourceUrl = ur

	return m, nil
}

func (b *Bucket) GetImageUrlByPid(ctx *gin.Context, pid, fileType string) (url string) {
	if pid == "" {
		return url
	}
	if match, err := regexp.MatchString(b.FilePrefix, pid); err != nil || !match {
		zlog.Warnf(ctx, "Invalid pid!")
		return url
	}

	name := pid + "." + fileType
	imgUrl := b.GetUrlByFileName(name)
	return imgUrl
}

/*
@Title  GetThumbnailUrlByPid
@Description  根据快速缩略模板提供缩略图
@param pid  string 图片pid
@param thumbnail string 图片pid 缩略方案，可以指定限制最大宽度或最大高度或同时限制，格式为"w/number","h/number","w/num/h/num"
@param fileType  string 在cos的文件名后缀，默认为jpg
@param outType string 指定缩略图的格式，默认和原图一致
@return url	string	缩略图url
*/
func (b *Bucket) GetThumbnailUrlByPid(ctx *gin.Context, pid, thumbnail, fileType, outType string) (url string) {
	if pid == "" {
		return url
	}
	if b.Thumbnail != 1 {
		zlog.Warnf(ctx, "CosService not support thumbnail!")
		return url
	}

	if match, err := regexp.MatchString(b.FilePrefix, pid); err != nil || !match {
		zlog.Warnf(ctx, "Invalid pid!")
		return url
	}

	name := pid + "." + fileType
	imgUrl := b.GetUrlByFileName(name)
	url = fmt.Sprintf("%s?imageView2/2/%s", imgUrl, thumbnail)
	if outType != "" {
		url = fmt.Sprintf("%s/format/%s", url, outType)
	}
	return url
}

func (b *Bucket) GetUrlByFileName(name string) string {
	if name == "" {
		return ""
	}
	objectKey := name
	if b.Directory != "" {
		objectKey = b.Directory + "/" + name
	}
	return fmt.Sprintf("http://%s-%s.cos.%s.myqcloud.com/%s", b.Name, b.AppID, b.Region, objectKey)
}

// pid是否符合命名规范
func CheckPidValid(pid string) bool {
	pattern := `/^(zyb|qa)([\d]*)_[0-9a-zA-Z]+(\.[0-9a-zA-Z]+)?$/`
	if match, err := regexp.MatchString(pattern, pid); err == nil && match {
		return true
	}
	return false
}
