package cos

import (
	"errors"
	"fmt"
	"github.com/GitHub121380/golib/utils"
	"github.com/GitHub121380/golib/zlog"
	"github.com/gin-gonic/gin"
	tcos "github.com/tencentyun/cos-go-sdk-v5"
	"image"
	"io"
	"net/url"
	"os"
	"strings"
)

// 1M
const sliceSize1M = 1048576

// 20M 大于20M的文件需要进行分片传输
const maxUnSliceFileSize = 20971520

/*
@Title  UploadLocalFile
@Description  上传本地文件到Cos
@param localFilePath  string 本地文件路径
@param fileName string 上传到cos的文件名
@param fileType string 上传到cos的文件名后缀
@param overwrite  string 同名文件是否强制覆盖，0-不强制覆盖，1-强制覆盖
@param needSize string 是否需要在pid中写宽高，0-不需要，1-需要
*/
func (b *Bucket) UploadLocalFile(ctx *gin.Context, localFilePath, fileName, fileType string, overwrite, needSize bool) (dwUrl string, err error) {
	f, err := os.Open(localFilePath)
	if err != nil {
		zlog.Warnf(ctx, "open local file error: %s", err.Error())
		return dwUrl, err
	}
	s, err := f.Stat()
	if err != nil {
		zlog.Warnf(ctx, "read file error: %s", err.Error())
		return dwUrl, err
	}

	fileSize := int(s.Size())
	if fileSize > b.FileSizeLimit {
		zlog.Warnf(ctx, "local file %s exceeds size limit %s ", localFilePath, b.FileSizeLimit)
		return dwUrl, errors.New("file exceeds max file size")
	}

	if fileType == "" {
		fileType = "jpg"
	}

	var objectKey string
	objectKey = fileName + "." + fileType
	if fileName == "" {
		// 默认规则生成存储到cos的object key
		name := utils.Md5(localFilePath)
		if needSize {
			if im, _, err := image.DecodeConfig(f); err == nil {
				name = fmt.Sprintf("%s_%d_%d", name, im.Width, im.Height)
			}
		}
		objectKey = fmt.Sprintf("%s%s.%s", b.FilePrefix, name, fileType)
	}

	if b.Directory != "" {
		objectKey = fmt.Sprintf("%s/%s", b.Directory, objectKey)
	}
	ur := b.getUploadUrl(objectKey)
	u, err := url.Parse(ur)
	if err != nil {
		zlog.Warn(ctx, "uploadUrl Parse error ")
		return dwUrl, err
	}

	basicUrl := &tcos.BaseURL{
		BucketURL: u,
	}
	client := tcos.NewClient(basicUrl, b.conn)
	if fileSize < maxUnSliceFileSize {
		// 直接上传
		opt := &tcos.ObjectPutOptions{
			ObjectPutHeaderOptions: &tcos.ObjectPutHeaderOptions{
				ContentLength: fileSize,
			},
		}
		_, err = client.Object.Put(ctx, objectKey, f, opt)
		if err != nil {
			return dwUrl, err
		}
	} else {
		// 分块上传
		up, _, err := client.Object.InitiateMultipartUpload(ctx, objectKey, nil)
		if err != nil {
			return dwUrl, err
		}
		parts := (fileSize / sliceSize1M) + 1
		opt := &tcos.CompleteMultipartUploadOptions{}
		for i := 1; i <= parts; i++ {
			buf := make([]byte, sliceSize1M)
			_, err := f.Read(buf)
			if err != nil {
				return dwUrl, err
			}

			etag, err := uploadPart(ctx, client, objectKey, up.UploadID, i, strings.NewReader(string(buf)))
			if err != nil {
				return etag, err
			}
			opt.Parts = append(opt.Parts, tcos.Object{
				PartNumber: i,
				ETag:       etag,
			},
			)
		}
		_, _, err = client.Object.CompleteMultipartUpload(ctx, objectKey, up.UploadID, opt)
		if err != nil {
			return dwUrl, err
		}
	}
	dst := b.getDownloadUrl(objectKey)
	return dst["sourceUrl"], nil
}

func (b *Bucket) UploadFileContent(ctx *gin.Context, content, fileName, fileType string, overwrite bool) (dwUrl string, err error) {
	if content == "" {
		return dwUrl, errors.New("content is empty")
	}

	fileSize := len(content)
	if fileSize > b.FileSizeLimit {
		return dwUrl, errors.New("file content length exceeds size limit")
	}

	if fileType == "" {
		fileType = "jpg"
	}

	var objectKey string
	objectKey = fileName + "." + fileType
	if objectKey == "" {
		// 默认规则生成存储到cos的object key
		objectKey = b.FilePrefix + utils.Md5(content) + "." + fileType
	}

	if b.Directory != "" {
		objectKey = fmt.Sprintf("%s/%s", b.Directory, objectKey)
	}
	ur := b.getUploadUrl(objectKey)
	u, err := url.Parse(ur)
	if err != nil {
		zlog.Warn(ctx, "uploadUrl Parse error ")
		return dwUrl, err
	}

	basicUrl := &tcos.BaseURL{
		BucketURL: u,
	}
	client := tcos.NewClient(basicUrl, b.conn)
	if fileSize < maxUnSliceFileSize {
		// 直接上传
		opt := &tcos.ObjectPutOptions{
			ObjectPutHeaderOptions: &tcos.ObjectPutHeaderOptions{
				ContentLength: fileSize,
			},
		}
		_, err = client.Object.Put(ctx, objectKey, strings.NewReader(content), opt)
		if err != nil {
			return dwUrl, err
		}
	} else {
		// 分块上传
		up, _, err := client.Object.InitiateMultipartUpload(ctx, objectKey, nil)
		if err != nil {
			return dwUrl, err
		}

		blockSize := sliceSize1M
		parts := (fileSize / blockSize) + 1
		opt := &tcos.CompleteMultipartUploadOptions{}
		for i := 1; i <= parts; i++ {
			pos := sliceSize1M
			if len(content) < sliceSize1M {
				pos = len(content)
			}
			buf := content[:pos]
			content = content[pos:]
			etag, err := uploadPart(ctx, client, objectKey, up.UploadID, i, strings.NewReader(buf))
			if err != nil {
				return etag, err
			}
			e := tcos.Object{
				PartNumber: i,
				ETag:       etag,
			}
			opt.Parts = append(opt.Parts, e)
		}
		_, _, err = client.Object.CompleteMultipartUpload(ctx, objectKey, up.UploadID, opt)
		if err != nil {
			return dwUrl, err
		}
	}
	dst := b.getDownloadUrl(objectKey)
	return dst["sourceUrl"], nil
}

func uploadPart(ctx *gin.Context, client *tcos.Client, objectKey, UploadID string, partNo int, buf io.Reader) (etag string, error error) {
	resp, err := client.Object.UploadPart(ctx, objectKey, UploadID, partNo, buf, nil)
	if err != nil {
		return etag, err
	}
	etag = resp.Header.Get("Etag")
	return etag, nil
}

func (b *Bucket) getUploadUrl(dstPath string) string {
	return fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", b.Name, b.AppID, b.Region)
}

func (b *Bucket) getDownloadUrl(objectKey string) (urls map[string]string) {
	srcUrl := fmt.Sprintf("https://%s-%s.%s.myqcloud.com", b.Name, b.AppID, b.PictureRegion)
	accessUrl := fmt.Sprintf("https://%s-%s.file.myqcloud.com", b.Name, b.AppID)
	sourceUrl := fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", b.Name, b.AppID, b.Region)

	if objectKey != "" {
		srcUrl = fmt.Sprintf("%s/%s", srcUrl, objectKey)
		accessUrl = fmt.Sprintf("%s/%s", accessUrl, objectKey)
		sourceUrl = fmt.Sprintf("%s/%s", sourceUrl, objectKey)
	}

	urls = map[string]string{
		"srcUrl":    srcUrl,
		"accessUrl": accessUrl,
		"sourceUrl": sourceUrl,
	}
	return urls
}

func (b *Bucket) Download2Local(ctx *gin.Context, srcFileName, dstFileName string) (err error) {
	if srcFileName == "" {
		return errors.New("fileName is empty")
	}

	objectKey := srcFileName
	if b.Directory != "" {
		objectKey = fmt.Sprintf("%s/%s", b.Directory, srcFileName)
	}

	ur := b.getDownloadUrl("")
	//zlog.Debugf(ctx, "download2Local get %s download uri is %s", srcFileName, ur["srcUrl"])

	u, _ := url.Parse(ur["sourceUrl"])
	baseUrl := &tcos.BaseURL{
		BucketURL: u,
	}
	cObj := tcos.NewClient(baseUrl, b.conn)
	_, err = cObj.Object.GetToFile(ctx, objectKey, dstFileName, nil)
	if err != nil {
		zlog.Warnf(ctx, "Tcos GetToFile error: %s", err.Error())
		return err
	}
	return nil
}

func (b *Bucket) GetObjectList(ctx *gin.Context, bucket string, opt *tcos.BucketGetOptions) (objectList []tcos.Object, err error) {
	u, err := getBucketUrl(bucket, b.AppID, b.Region)
	if err != nil {
		return objectList, err
	}

	basicUrl := &tcos.BaseURL{
		BucketURL: u,
	}
	// 查询对象列表
	client := tcos.NewClient(basicUrl, b.conn)
	v, _, err := client.Bucket.Get(ctx, opt)
	if err != nil {
		return objectList, err
	}
	return v.Contents, nil
}
