package main

import (
	"encoding/json"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"gocms/conf"
	"log"
	"net/http"
	"strconv"
	"time"

	"gocms/api/v1"
	_ "gocms/dao"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"
//
//func helloHandler(c *gin.Context) {
//	claims := jwt.ExtractClaims(c)
//	user, _ := c.Get(identityKey)
//	c.JSON(200, gin.H{
//		"userID":   claims["id"],
//		"userName": user.(*User).UserName,
//		"text":     "Hello World.",
//	})
//}

// User demo
type User struct {
	UserName  string
}

func main() {

	port := conf.Conf.Server.Port
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(Cors())

	if port == 0 {
		port = 80
	}

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "web",
		Key:         []byte(conf.Conf.App.Key),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == conf.Conf.App.UserName && password == conf.Conf.App.Password) {
				return &User{
					UserName:  userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"msg": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	authMiddleware = authMiddleware
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}


	//
	//r.POST("/login", authMiddleware.LoginHandler)
	//
	r.NoRoute(authMiddleware.MiddlewareFunc(),func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	//
	//auth := r.Group("/auth")
	////Refresh time can be longer than token timeout
	//auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	routeV1 := r.Group("/api/v1")
	routeV1.Use(authMiddleware.MiddlewareFunc())

	// user api
	{
		uR := routeV1.Group("user")
		r.POST("/api/v1/user/login", authMiddleware.LoginHandler)
		//c.JSON(200, gin.H{"code": 0, "data": gin.H{"token": "admin"}})
		uR.GET("info", func(c *gin.Context) {
			var b interface{}
			json.Unmarshal([]byte(`{"code":0,"data":{"roles":["admin"],"name":"admin","avatar":"https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"}}`), &b)
			c.JSON(200, b)
		})
	}


	// translate title to English url
	{
		routeV1.GET("translate", (&v1.TranslateApi{}).TitleToUrl)
	}

	// Article Api
	{
		articleApi := v1.ArticleApi{}
		route := routeV1.Group("/article")
		r.GET("api/v1/articles", articleApi.Select)
		routeV1.POST("article", articleApi.Add)
		routeV1.POST("article/:article", articleApi.Add)
		route.DELETE(":type/:article", articleApi.Delete)
		route.DELETE(":type/:article/:chapter", articleApi.Delete)
		route.PUT(":type/:article", articleApi.Edit)
		route.PUT(":type/:article/:chapter", articleApi.Edit)
		r.GET("api/v1/article/:type/:article", articleApi.Get)
		r.GET("api/v1/article/:type/:article/:chapter", articleApi.Get)
	}

	// ArticleDao Type Api
	{
		articleType := &v1.ArticleTypeApi{}
		routeType := routeV1.Group("article-type")
		routeV1.POST("article-type", articleType.Add)
		routeType.DELETE(":name", articleType.Delete)
		routeType.PUT(":name", articleType.Edit)
		r.GET("api/v1/article-type/:name", articleType.Get)
		r.GET("api/v1/article-types", articleType.Select)
	}

	// file upload api
	{
		uploadApi := &v1.UploadApi{}
		route := routeV1.Group("upload")
		route.GET("token", uploadApi.Token)
		r.POST("/api/v1/upload/callback", uploadApi.Callback)
	}

	if err := http.ListenAndServe(":"+strconv.Itoa(port), r); err != nil {
		log.Fatal(err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		//放行所有OPTIONS方法
		method := c.Request.Method
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
			return
		}
		c.Next()
	}
}
