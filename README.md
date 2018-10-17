# gocms


## 说明

* 本项目只是不学网的一部分，全部项目组合起来构成不学网

[后台API](https://github.com/noxue/gocms)

[前端界面](https://github.com/noxue/gocms-ui)

[后台界面](https://github.com/noxue/gocms-admin-ui)

[SEO处理工具](https://github.com/noxue/gocms-seo)

## 网站整体架构

* 关于网站设计思路和开发中遇到的问题，可以看这里

http://noxue.com/article/unclassified/gocms

## 部署

* 把配置文件 config.default.json 复制一份改名为 config.json，根据文件中的配置修改配置内容
* 编译成可执行程序，和config.json放到一个目录下就可以运行
* **注意** 如果要上传文件，必须要配置七牛

## 演示地址

[不学网](http://noxue.com)

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