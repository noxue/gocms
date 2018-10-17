package dao

import (
	"gocms/conf"
	"gopkg.in/mgo.v2"
)

func getCollect(session *mgo.Session, collectionName string) *mgo.Collection {
	return session.DB(conf.Conf.Db.DbName).C(collectionName)
}
