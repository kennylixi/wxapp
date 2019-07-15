package controller

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"io/ioutil"
	"qiniupkg.com/x/bytes.v7"
	"wxapp/config"
)

type UploadController struct {
	BaseController
}

func (c *UploadController) Post() {
	file, _, err := c.Ctx.FormFile("file")
	defer file.Close()
	if err != nil {
		c.SendMsg(1, "上传文件失败", nil)
		return
	}

	//fname := info.Filename
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.SendMsg(1, "上传文件失败", nil)
		return
	}
	dataLen := int64(len(data))

	putPolicy := storage.PutPolicy{
		Scope:         config.QiniuBucket,
		PersistentOps: config.QPersistentOps,
	}
	mac := qbox.NewMac(config.QAccessKey, config.QSecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, bytes.NewReader(data), dataLen, nil)
	if err != nil {
		c.SendMsg(1, "上传文件失败", nil)
		return
	}
	c.SendMsg(0, "上传文件成功", iris.Map{"url": fmt.Sprintf(config.QDomain, ret.Key)})
}

//simditor编辑器上传
func (c *UploadController) PostSimditor() {
	file, _, err := c.Ctx.FormFile("file")
	defer file.Close()
	if err != nil {
		c.SendMsgData(iris.Map{"success": false, "msg": "上传文件失败", "file_path": nil})
		return
	}

	//fname := info.Filename
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.SendMsgData(iris.Map{"success": false, "msg": "上传文件失败", "file_path": nil})
		return
	}
	dataLen := int64(len(data))

	putPolicy := storage.PutPolicy{
		Scope:         config.QiniuBucket,
		PersistentOps: config.QPersistentOps,
	}
	mac := qbox.NewMac(config.QAccessKey, config.QSecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, bytes.NewReader(data), dataLen, nil)
	if err != nil {
		c.SendMsgData(iris.Map{"success": false, "msg": "上传文件失败", "file_path": nil})
		return
	}
	c.SendMsgData(iris.Map{"success": true, "msg": "上传文件成功", "file_path": fmt.Sprintf(config.QDomain, ret.Key)})
}
