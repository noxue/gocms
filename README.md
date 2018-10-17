# gocms


## 说明

* 本项目只是后台接口，需要配合前端界面和前端管理界面

[前端界面项目地址](https://github.com/noxue/gocms-ui)
[后台界面项目地址](https://github.com/noxue/gocms-admin-ui)


### response error code

* `1000` post data error
* `1111` form request field is invalid


### restful api (v1)

#### Upload Api

* `GET /upload/token` Get Upload Token For QiNiu
    * `GET 127.0.0.1/api/v1/upload/token?filename=1.jpg&type=article`
    
* `POST /upload/callback` upload callback api for qiniu notify my server
    * file 
    * token xxxxxxxxxxxx
    * key article/2018/10/2a498b27-bd82-4c0b-8f41-85529dcd1a05.jpg

###### Article Type Api

* `POST /article-type` add article type
* `DELETE /article-type/{typeName}` delete article type
* `PUT /article-type/{typeName}` update article type
* `GET /article-type/{typeName}` get a article type
* `GET /article-types?word={keyword}&page={page}&size={pagesize}&sortby={fieldName}&order={asc|desc}` get article types

###### Article Api

* `POST /article` add article
* `DELETE /article/{articleName}` delete article
* `POST /article/{articleName}` add a chapter to the article
* `PUT /article/{articleName}` update article
* `GET /article/{typeName}/{articleName}` get a article
* `GET /articles?page={page}&size={pagesize}&sortby={fieldName}&order={asc|desc}` get articles 
* `GET /article/{articleName}?page={page}&size={pagesize}&sortby={fieldName}&order={asc|desc}` get articles chapter
* `GET /article/{typeName}/{articleName}/{chapterName}` get a chapter

* 上面接口有所改动，具体可以看`main.go`