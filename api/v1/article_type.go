package v1

import (
	"github.com/gin-gonic/gin"
	"gocms/dao"
	"gocms/domain"
	"gocms/srv"
	"gocms/utils"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
)

type ArticleTypeApi struct {
}

func (this *ArticleTypeApi)Add(c *gin.Context) {
	var articleType domain.ArticleType

	err := c.BindJSON(&articleType)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"code": 1000, "msg": "Post Data Format Must JSON"})
		return
	}

	errs := utils.ValidateStruct(articleType)
	if len(errs) > 0 {
		c.JSON(200, gin.H{"code": 1111, "msg": "field valid failed", "errors": errs})
		return
	}

	at :=dao.NewArticleTypeDao()
	_, err = at.Get(articleType.Name)
	if err==nil {
		c.JSON(200, gin.H{"code": 1001, "msg": "the type name is exists:"+articleType.Name})
		return
	}

	if err := dao.NewArticleTypeDao().Insert(&articleType); err != nil {
		c.JSON(200, gin.H{"code": 1001, "msg": "insert type failed"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

func (this *ArticleTypeApi)Edit(c *gin.Context) {

	_, err := dao.NewArticleTypeDao().Get(c.Param("name"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(404)
		return
	}

	var t domain.ArticleType

	err = c.BindJSON(&t)
	if err != nil {
		c.JSON(200, gin.H{"code": 1000, "msg": "Post Data Err"})
		return
	}

	errs := utils.ValidateStruct(t)
	if len(errs) > 0 {
		c.JSON(200, gin.H{"code": 1111, "msg": "field valid failed", "errors": errs})
		return
	}

	if c.Param("name") != t.Name {
		_, err = dao.NewArticleTypeDao().Get(t.Name)
		if err == nil {
			c.JSON(200, gin.H{"code": 1001, "msg": "the type name is exists"})
			return
		}
	}

	if err := srv.NewArticleTypeSrv().Edit(c.Param("name"),&t); err != nil {
		c.JSON(200, gin.H{"code": 1001, "msg": "update type failed"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

func (this *ArticleTypeApi)Delete(c *gin.Context) {
	err:=srv.NewArticleTypeSrv().Delete(c.Param("name"))
	if  err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "article type delete failed:" + err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

func (this *ArticleTypeApi)Get(c *gin.Context) {
	name := c.Param("name")
	t, err := dao.NewArticleTypeDao().Get(name)
	if err!=nil {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": t})
}

func (this *ArticleTypeApi)Select(c *gin.Context) {

	word := c.Query("word")
	page := c.DefaultQuery("page","1")
	size := c.DefaultQuery("size","10")
	sort := c.DefaultQuery("sortby","createdat")
	order := c.DefaultQuery("order","asc")

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "page param must be a number"})
		return
	}

	pageSize, err := strconv.Atoi(size)
	if err != nil {
		c.JSON(200, gin.H{"code": -2, "msg": "size param must be a number"})
		return
	}

	conditions := []bson.M{
		bson.M{"title": bson.M{"$regex": bson.RegEx{word, "i"}}},
	}

	if order == "desc" {
		sort = "-" + sort
	}
	orders := []string{sort}

	types, err := dao.NewArticleTypeDao().Select(pageNum, pageSize, conditions, orders)
	if err != nil {
		c.JSON(200, gin.H{"code": -2, "msg": "select type list failed"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "types": types})
}

