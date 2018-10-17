package srv

import (
	"gocms/dao"
	"gocms/domain"
	"qiniupkg.com/x/errors.v7"
	"time"
)

type ArticleSrv struct {
}

func NewArticleSrv() *ArticleSrv {
	return &ArticleSrv{}
}

func (this *ArticleSrv) Add(article *domain.Article) (err error) {

	articleTypeDao := dao.NewArticleTypeDao()

	t, err := articleTypeDao.Get(article.Type)
	if t == nil || err != nil {
		err = errors.New("this type is not exists:" + article.Type)
		return
	}

	// if the article is exists, return
	articleDao := dao.NewArticleDao()
	_, err = articleDao.Get(article.Type, article.Name, "")
	if err == nil {
		err = errors.New("this article is already exists:" + article.Name)
		return
	}

	article.CreatedAt = time.Now()
	article.Hits = 0

	err = articleDao.Insert("", article)

	return
}

func (this *ArticleSrv) AddChapter(article string, chapter *domain.Article) (err error) {

	articleDao := dao.NewArticleDao()
	// check the article
	_, err = articleDao.Get(chapter.Type, article, "")
	if err != nil {
		err = errors.New("the article is not exists:" + article)
		return
	}

	// check the chapter
	_, err = articleDao.Get(chapter.Type, article, chapter.Name)
	if err == nil {
		err = errors.New("the article is exists:" + chapter.Name)
		return
	}

	chapter.Type = ""
	chapter.CreatedAt = time.Now()
	chapter.Hits = 0
	err = articleDao.Insert(article, chapter)

	return
}

func (this *ArticleSrv) Get(ty, article, chapter string) (a *domain.Article, err error) {
	articleDao := dao.NewArticleDao()

	// check the article is exists
	a, err = articleDao.Get(ty, article, chapter)
	if err != nil {
		err = errors.New("the article is not exists")
		return
	}
	return
}

func (this *ArticleSrv) Edit(oldType, oldName string, article *domain.Article) (err error) {

	articleTypeDao := dao.NewArticleTypeDao()

	t, err := articleTypeDao.Get(article.Type)
	if t == nil || err != nil {
		err = errors.New("this type is not exists:" + article.Type)
		return
	}

	// if the article is exists, return
	articleDao := dao.NewArticleDao()
	if oldName != article.Name {
		_, err = articleDao.Get(article.Type, article.Name, "")
		if err == nil {
			err = errors.New("this article is already exists:" + article.Name)
			return
		}
	}

	err = articleDao.UpdateArticle(oldType, oldName, article)

	return
}

func (this *ArticleSrv) EditChapter(oldType, oldName, oldChapter string, chapter *domain.Article) (err error) {

	articleDao := dao.NewArticleDao()
	// check the article
	_, err = articleDao.Get(oldType, oldName, oldChapter)
	if err != nil {
		err = errors.New("the article is not exists:" + oldChapter)
		return
	}

	// check the chapter
	if oldChapter != chapter.Name {
		_, err = articleDao.Get(oldType, oldName, chapter.Name)
		if err == nil {
			err = errors.New("the chapter is exists:" + chapter.Name)
			return
		}
	}

	chapter.Type = ""
	err = articleDao.UpdateChapter(oldType, oldName, oldChapter, chapter)

	return
}

func (this *ArticleSrv)Delete(typeName, articleName string) (err error) {
	dao.NewArticleDao().Delete(typeName, articleName)
	return
}

func (this *ArticleSrv)DeleteChapter(typeName, articleName, chapterName string) (err error) {
	err = dao.NewArticleDao().DeleteChapter(typeName, articleName, chapterName)
	return
}