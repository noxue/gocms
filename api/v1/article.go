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

type ArticleApi struct {
}

func (this *ArticleApi) Add(c *gin.Context) {
	var article domain.Article
	err := c.BindJSON(&article)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"code": 1000, "msg": "Post Data Format Must JSON"})
		return
	}

	errs := utils.ValidateStruct(article)
	if len(errs) > 0 {
		c.JSON(200, gin.H{"code": 1111, "msg": "field valid failed", "errors": errs})
		return
	}

	articleSrv := srv.NewArticleSrv()

	pa := c.Param("article")
	if pa == "" {
		err = articleSrv.Add(&article)
	} else {
		err = articleSrv.AddChapter(pa, &article)
	}

	if err != nil {
		c.JSON(200, gin.H{"code": 222, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

func (this *ArticleApi) Edit(c *gin.Context) {
	var article domain.Article
	err := c.BindJSON(&article)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{"code": 1000, "msg": "Post Data Format Must JSON"})
		return
	}

	oldType := c.Param("type")
	oldArticle := c.Param("article")
	oldChapter := c.Param("chapter")

	errs := utils.ValidateStruct(article)
	if len(errs) > 0 {
		c.JSON(200, gin.H{"code": 1111, "msg": "field valid failed", "errors": errs})
		return
	}

	articleSrv := srv.NewArticleSrv()

	if oldChapter == "" {
		err = articleSrv.Edit(oldType, oldArticle, &article)
	} else {
		err = articleSrv.EditChapter(oldType, oldArticle, oldChapter, &article)
	}
	if err != nil {
		c.JSON(200, gin.H{"code": 222, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

func (this *ArticleApi) Delete(c *gin.Context) {
	ty := c.Param("type")
	article := c.Param("article")
	chapter := c.Param("chapter")

	articleSrv := srv.NewArticleSrv()
	if chapter != "" {
		err := articleSrv.DeleteChapter(ty, article, chapter)
		if err != nil {
			c.JSON(200, gin.H{"code": 222, "msg": "delete chapter error:" + err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 0, "msg": "success"})
		return
	}

	err := articleSrv.Delete(ty, article)
	if err != nil {
		c.JSON(200, gin.H{"code": 222, "msg": "delete article error:" + err.Error()})
	}
	c.JSON(200, gin.H{"code": 0, "msg": "success"})
}

func (this *ArticleApi) Get(c *gin.Context) {
	ty := c.Param("type")
	article := c.Param("article")
	chapter := c.Param("chapter")
	articleSrv := srv.NewArticleSrv()
	a, err := articleSrv.Get(ty, article, chapter)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	if chapter != "" && len(a.Chapters) == 1 {
		a.Chapters[0].Type = a.Type
		a = a.Chapters[0]
		c.JSON(200, gin.H{"code": 0, "article": a})
		return
	}
	c.JSON(200, gin.H{"code": 0, "article": a})
}

func (this *ArticleApi) Select(c *gin.Context) {
	word := c.Query("word")
	page := c.DefaultQuery("page", "1")
	size := c.DefaultQuery("size", "10")
	sort := c.DefaultQuery("sortby", "createdat")
	order := c.DefaultQuery("order", "asc")
	typeName := c.Query("type")
	article := c.Query("article")

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

	conditions := bson.M{}

	if word != "" {
		if article == "" {
			conditions["title"] = bson.M{"$regex": bson.RegEx{word, "i"}}
		} else {
			conditions["chapters.title"] = bson.M{"$regex": bson.RegEx{word, "i"}}
		}
	}

	if typeName != "" {
		conditions["type"] = typeName
	}

	if order == "desc" {
		sort = "-" + sort
	}
	orders := []string{sort}

	var articles []domain.Article
	var a domain.Article
	if article == "" {
		articles, err = dao.NewArticleDao().Select(pageNum, pageSize, conditions, orders)
	} else {
		conditions["name"] = article
		a, err = dao.NewArticleDao().SelectChapter(article, pageNum, pageSize, conditions, orders)
	}
	if err != nil {
		log.Println(err)
	}
	if article == "" {
		c.JSON(200, gin.H{"code": 0, "articles": articles})
	} else {
		c.JSON(200, gin.H{"code": 0, "article": a})
	}
}
