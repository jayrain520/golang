package config

import (
	"encoding/json"
	"github.com/wonderivan/logger"
	"gopkg.in/olivere/elastic.v6"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testgin/api/response"
)

type Seting struct {
	//mysql configure
	DbConnectAddress string `json:"db_connect_address"`
	DbName           string `json:"db_name"`
	UserName         string `json:"user_name"`
	Port             string `json:"port"`
	PassWord         string `json:"password"`
	Charset          string `json:"charset"`
	SetMaxOpenConns	 int64  `json:"set_max_open_conns"`
	SetMaxIdleConns  int64  `json:"set_max_idle_conns"`
	//local storage path
	VideoSavePath string `json:"local_video_path"`

	//video and image src bind
	Bind string `json:"src_bind"`
	//server run port
	BindPort string `json:"server_run_port"`

	//aliyun oss storage
	BucketName      string `json:"bucket_name"`
	ResourceAddress string `json:"resource_address"`
	OssVideoPath    string `json:"oss_videos_path"`
	OpenOssServe    bool   `json:"open_oss_serve"`
	EndPoint        string `json:"end_point"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`

	//logger level set
	Level string `json:"level"`

	//elastic search configure
	SetSniff       bool   `json:"set_sniff"`
	OpenElastic    bool   `json:"open_elastic"`
	Index          string `json:"index"`
	Type           string `json:"type"`
	InitDataNumber int64  `json:"init_data_number"`
	ElsClient      *elastic.Client
	Address        string `json:"address"`


	//email config
	EmailTitle string `json:"email_title"`
	EmailUser string `json:"email_user"`
	EmailAccesskey string `json:"email_accesskey"`
	EmailHost string `json:"email_host"`
	EmailPort int64 `json:"email_port"`
}

var (
	Set *Seting
)

func init() {
	Set = newConfig()
	err := logger.SetLogger(`{"Console": {"level":"` + Set.Level + `","color": true}}`)
	if err != nil {
		log.Printf("[ERROR] :%v\n", err)
	}

	if Set.OpenElastic {
		Set.ElsClient = newElasticClient()
	} else {
		logger.Warn(response.NotOpenEls.ResponseMsg)
	}
}

func newConfig() *Seting {
	dir, _ := os.Getwd()
	file, err := os.Open(dir + "/seting.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//start decode json config
	var decode map[string]interface{}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &decode)
	if err != nil {
		panic(err)
	}
	var set Seting
	startDecode(decode, reflect.TypeOf(&set), reflect.ValueOf(&set))
	return &set
}

func startDecode(obj interface{}, rSetType reflect.Type, rSetValue reflect.Value) {
	rObj := reflect.ValueOf(obj)
	OneKey := rObj.MapKeys()

	for _, k := range OneKey {
		TowKey := rObj.MapIndex(k).Interface().(map[string]interface{})

		rvalue1 := reflect.ValueOf(TowKey)
		//拿到第二层的全部key
		OkTowKey := rvalue1.MapKeys()

		for _, k1 := range OkTowKey {
			for i := 0; i < rSetValue.Elem().NumField(); i++ {
				if rSetType.Elem().Field(i).Tag.Get("json") == k1.String() {

					switch rSetType.Elem().Field(i).Type.Kind() {
					case reflect.String:
						rSetValue.Elem().Field(i).SetString(rvalue1.MapIndex(k1).Interface().(string))
					case reflect.Bool:
						rSetValue.Elem().Field(i).SetBool(rvalue1.MapIndex(k1).Interface().(bool))
					case reflect.Int64:
						rSetValue.Elem().Field(i).SetInt(int64(rvalue1.MapIndex(k1).Interface().(float64)))
					}

				}
			}

		}
	}

}

func newElasticClient() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetSniff(Set.SetSniff), //在docker中运行,需要设置false
		elastic.SetURL(Set.Address),
	)
	if err != nil {
		logger.Emer(err)
		return nil
	}
	return client
}
