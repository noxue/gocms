package dao

import (
	"time"
	"gopkg.in/mgo.v2"
)

type Attachment struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Hash   string `json:"hash"`
	Size   int64  `json:"size"`
	Prefix string `json:"prefix"`
	CreateAt time.Time
}

func getAttachmentCollection(session *mgo.Session) *mgo.Collection{
	return getCollect(session,"attachment")
}

func AttachInsert(attachment *Attachment) error {
	sess:=GetSession()
	defer sess.Close()
	c := getAttachmentCollection(sess)
	attachment.CreateAt = time.Now()
	return c.Insert(attachment)
}
