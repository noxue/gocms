package utils

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"gocms/conf"
)

var QiNiuMac *qbox.Mac
var BucketManager *storage.BucketManager

func initQiniu(){
	QiNiuMac = qbox.NewMac(conf.Conf.Upload.QiNiu.Id, conf.Conf.Upload.QiNiu.Key)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	BucketManager = storage.NewBucketManager(QiNiuMac, &cfg)
}
