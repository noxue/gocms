package v1

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"gocms/conf"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type TranslateApi struct {

}

func (this *TranslateApi)TitleToUrl(c *gin.Context) {
	title := c.Query("title")

	re1, _ := regexp.Compile(`[\s~!@<>#$%\+\^\&\*\(\)\./\\]+`)
	title = re1.ReplaceAllString(strings.ToLower(title), ",")

	appid := conf.Conf.TranslateConfig["baidu"].Id
	key := conf.Conf.TranslateConfig["baidu"].Key
	salt := time.Now().Unix()
	sign := fmt.Sprintf("%s%s%d%s", appid, title, salt, key);

	// md5(sign)
	data := []byte(sign)
	has := md5.Sum(data)
	sign = fmt.Sprintf("%x", has) //将[]byte转成16进制

	urlStr := fmt.Sprintf("http://api.fanyi.baidu.com/api/trans/vip/translate?q=%s&appid=%s&salt=%d&from=cn&to=en&sign=%s&_=%d",
		title,
		appid,
		salt,
		sign,
		time.Now().Unix(),
	)
	resp, err := http.Get(urlStr)
	if err != nil {
		c.JSON(200,gin.H{"code":-1,"msg":"request baidu translate api error"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(200,gin.H{"code":-2,"msg":"request baidu translate api error"})
	}

	reg := regexp.MustCompile(`.*?"dst":"(.*?)"}]}`)
	rets:=reg.FindStringSubmatch(string(body))

	if len(rets)!=2 {
		c.JSON(200, gin.H{"code": 3, "msg": "tanslate failed"})
		return
	}

	ret := rets[1]
	re, _ := regexp.Compile(`[^0-9a-z]+`)
	ret = re.ReplaceAllString(strings.ToLower(ret), "-")
	ret=strings.Trim(ret,"-")
	c.JSON(200, gin.H{"code": 0, "data": ret})
}