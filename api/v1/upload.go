package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/storage"
	"gocms/conf"
	"path/filepath"
	"strings"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
	"gocms/dao"
	"encoding/json"
	"gocms/utils"
	"log"
)

type UploadApi struct {

}

func (this *UploadApi)Token(c *gin.Context) {

	filename := c.Query("filename")
	uploadType := c.Query("type")

	upConfig := conf.Conf.Upload

	upType, ok := upConfig.Types[uploadType]
	if !ok {
		c.JSON(200, gin.H{"code": -1, "msg": "don't support the type:" + uploadType})
		return
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if len(ext) == 0 {
		c.JSON(200, gin.H{"code": -2, "msg": "don't support this file extension:" + ext, "allow": upType.Allow})
		return
	}

	allow := false
	ext = ext[1:]
	for _, t := range upType.Allow {
		if ext == t {
			allow = true
			break
		}
	}

	if !allow {
		c.JSON(200, gin.H{"code": -2, "msg": "don't support this file extension:" + ext, "allow": upType.Allow})
		return
	}
	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(200, gin.H{"code": -3, "msg": "new uuid error"})
		return
	}

	savePath := fmt.Sprintf("%s/%d/%d/%s."+ext, upType.Prefix, time.Now().Year(), time.Now().Month(), id)

	putPolicy := storage.PutPolicy{
		Scope:            fmt.Sprintf("%s:%s", upConfig.QiNiu.Bucket, savePath),
		CallbackURL:      upConfig.QiNiu.Callback,
		CallbackBody:     `{"path":"$(key)", "name":"$(fname)","hash":"$(etag)","size":$(fsize),"prefix":"` + upType.Prefix + `"}`,
		CallbackBodyType: "application/json",
	}

	token := putPolicy.UploadToken(utils.QiNiuMac)
	c.JSON(200, gin.H{"code": 0, "token": token, "size": upType.Size, "path": savePath})
}

type UploadRet struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
	Prefix string `json:"prefix"`
}

func (this *UploadApi)Callback(c *gin.Context) {

	var uploadRet UploadRet
	err := c.BindJSON(&uploadRet)
	if err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "request json format invalid"})
		return
	}

	upConfig := conf.Conf.Upload
	upType, ok := upConfig.Types[uploadRet.Prefix]
	if !ok {
		c.JSON(200, gin.H{"code": -2, "msg": "don't support the type:" + uploadRet.Prefix})
		return
	}

	// if the file is too large, delete it and notify client
	if uploadRet.Size > (upType.Size << 20) {
		c.JSON(200, gin.H{"code": -3, "msg": "file is too large, max size is " + fmt.Sprintf("%d",upType.Size) + "M"})
		err := utils.BucketManager.Delete(upConfig.QiNiu.Bucket,uploadRet.Path)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	var attachment dao.Attachment
	bs,_:=json.Marshal(uploadRet)
	json.Unmarshal(bs,&attachment)

	dao.AttachInsert(&attachment)


	c.JSON(200, gin.H{"code": 0, "msg": "upload file success", "data": uploadRet})
}
