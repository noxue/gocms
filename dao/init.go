package dao

import (
	"gocms/conf"
	"log"
	"gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session

func init() {
	var err error
	mgoSession, err = mgo.Dial(conf.Conf.Db.Host)  //connect database
	if err != nil {
		log.Fatal(err)
	}
}

func GetSession() *mgo.Session{
	return mgoSession.Copy()
}