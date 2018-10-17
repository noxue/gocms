package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var Conf Config

type Config struct {
	Db              Db                         `json:"db"`
	Server          Server                     `json:"server"`
	App             App                        `json:"app"`
	Upload          UploadConfig               `json:"upload"`
	TranslateConfig map[string]TranslateConfig `json:"translate"`
}

type Server struct {
	Port int `json:"port"`
}

type App struct {
	Key      string `json:"key"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type TypeConfig struct {
	Prefix string   `json:"prefix"`
	Allow  []string `json:"allow"`
	Size   int64    `json:"size"`
}

type UploadConfig struct {
	QiNiu QiNiuConfig           `json:"qiniu"`
	Types map[string]TypeConfig `json:"types"`
}

type QiNiuConfig struct {
	Id       string `json:"id"`
	Key      string `json:"key"`
	Bucket   string `json:"bucket"`
	Callback string `json:"callback"`
}

type TranslateConfig struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type Db struct {
	Host     string `json:"host"`
	DbName   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("read config file error:", err)
	}

	err = json.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatal("config file parse error:", err)
	}
}
