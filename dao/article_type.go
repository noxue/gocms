package dao

import (
	"gocms/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// article type
type ArticleTypeDao struct {
}

func NewArticleTypeDao() *ArticleTypeDao {
	return &ArticleTypeDao{}
}

func (this *ArticleTypeDao) getCol(session *mgo.Session) *mgo.Collection {
	return getCollect(session, "article_type")
}

func (this *ArticleTypeDao) Insert(articleType *domain.ArticleType) error {
	sess := GetSession()
	defer sess.Close()
	articleType.CreatedAt = time.Now()
	coll := this.getCol(sess)
	return coll.Insert(articleType)
}

func (this *ArticleTypeDao) Update(name string,articleType *domain.ArticleType) error {
	sess := GetSession()
	defer sess.Close()
	coll := this.getCol(sess)
	return coll.Update(bson.M{"name": name},
		bson.M{
			"$set": bson.M{
				"name":  articleType.Name,
				"title": articleType.Title,
				"sort":  articleType.Sort,
			},
		})
}

func (this *ArticleTypeDao) Delete(name string) error {
	sess := GetSession()
	defer sess.Close()
	coll := this.getCol(sess)
	return coll.Remove(bson.M{"name": name})
}

func (this *ArticleTypeDao) Get(name string) (t *domain.ArticleType, err error) {
	sess := GetSession()
	defer sess.Close()
	coll := this.getCol(sess)
	err = coll.Find(bson.M{"name": name}).One(&t)
	if err != nil {
		return
	}
	return
}

func (this *ArticleTypeDao) Select(page, pagesize int, conditions []bson.M, sorts []string) (ts []domain.ArticleType, err error) {
	sess := GetSession()
	defer sess.Close()
	coll := this.getCol(sess)
	cond := bson.M{}
	if conditions != nil {
		cond = bson.M{"$and": conditions}
	}
	err = coll.Find(cond).Limit(pagesize).Skip((page - 1) * pagesize).Sort(sorts...).All(&ts)
	return
}
