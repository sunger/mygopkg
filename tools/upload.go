package tools

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/sunger/mygopkg/config"
)

// 接收两个参数 一个文件流 一个 bucket 你的七牛云标准空间的名字
func Upload(file *multipart.FileHeader) (err error, path string, key string, bucket2 string) {
	c := config.GetConfig()
	defaultoss := c.GetString("oss.default")

	bkt := c.GetString(defaultoss + ".bucket")
	cb := c.GetString(defaultoss + ".notify")

	putPolicy := storage.PutPolicy{
		Scope: bkt,
	}

	if cb != "" {
		putPolicy = storage.PutPolicy{
			Scope:        bkt,
			CallbackURL:  cb,
			CallbackBody: "key=$(key)&hash=$(etag)&bucket=$(bucket)&fsize=$(fsize)&name=$(x:name)",
		}
	}

	ak := c.GetString(defaultoss + ".access-key")
	sk := c.GetString(defaultoss + ".secret-key")
	fmt.Println("ak", ak)
	fmt.Println("sk", sk)
	mac := qbox.NewMac(ak, sk)
	fmt.Println(mac)
	upToken := putPolicy.UploadToken(mac)
	fmt.Println("upToken", upToken)
	cfg := storage.Config{}
	// 空间对应的机房

	zone := c.GetString(defaultoss + ".zone")
	if zone == "huadong" {
		cfg.Zone = &storage.ZoneHuadong
	} else if zone == "huanan" {
		cfg.Zone = &storage.ZoneHuanan
	} else if zone == "huabei" {
		cfg.Zone = &storage.ZoneHuabei
	}

	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	//putExtra := storage.PutExtra{
	//	Params: map[string]string{
	//		"x:name": "github logo",
	//	},
	//}
	f, e := file.Open()
	if e != nil {
		fmt.Println(e)
		return e, "", "", ""
	}
	dataLen := file.Size
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	fmt.Println("fileKey", fileKey)
	err = formUploader.Put(context.Background(), &ret, upToken, fileKey, f, dataLen, nil)
	fmt.Println("err-----", err)
	if err != nil {
		fmt.Println("upload file fail:", err)
		return err, "", "", ""
	}
	return err, c.GetString(defaultoss+".img-path") + "/" + ret.Key, ret.Key, bkt
}

func DeleteFile(key string) error {
	c := config.GetConfig()
	defaultoss := c.GetString("oss.default")
	ak := c.GetString(defaultoss + ".access-key")
	sk := c.GetString(defaultoss + ".secret-key")
	mac := qbox.NewMac(ak, sk)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(c.GetString(defaultoss+".bucket"), key)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
