package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	//"github.com/silenceper/wechat/v2"
	//"github.com/silenceper/wechat/v2/cache"
	//offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/rest"
)

/*
@Time : 2021/9/18 9:19 AM
@Author : apple
@File : main
@Software: GoLand
*/
type User struct {
	TagID  int64
	OpenID []string
}

func main() {
	r := gin.Default()
	http.HandleFunc("/", checkout)
	r.GET("/ping", func(c *gin.Context) {
		token, _ := GetAccessToken()
		c.JSON(200, gin.H{
			"token": token,
		})
	})
	r.Run(":80")
	//wc := wechat.NewWechat()
	//redisOpts := &cache.RedisOpts{
	//	Host:        "127.0.0.1:6379",
	//	Database:    1,
	//	MaxActive:   10,
	//	MaxIdle:     10,
	//	IdleTimeout: 60, //second
	//}
	//redisCache := cache.NewRedis(redisOpts)
	//cfg := &offConfig.Config{
	//	AppID:     "wx870e0c515d19cde4",
	//	AppSecret: "ae4bf23de5e9fb9680d1fccfaf0fbbed",
	//	Token:     token,
	//	//EncodingAESKey: "xxxx",
	//	Cache: redisCache,
	//}
	//oa := wc.GetOfficialAccount(cfg)
	//bd:=oa.GetBroadcast()
	//
	//text, err := bd.SendText( &User{
	//	TagID: 1,
	//	OpenID: "",
	//}, "sssss")
	//if err != nil{
	//		fmt.Errorf("error is %s",err)
	//}
	//fmt.Printf("msg is %v",text)
	//officialAccount := wc.GetOfficialAccount(cfg)
}

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

func checkout(response http.ResponseWriter, request *http.Request) {
	//解析URL参数
	err := request.ParseForm()
	if err != nil {
		fmt.Println("URL解析失败！")
		return
	}
	// token
	var token string = "iwuqing"
	// 获取参数
	signature := request.FormValue("signature")
	timestamp := request.FormValue("timestamp")
	nonce := request.FormValue("nonce")
	echostr := request.FormValue("echostr")
	//将token、timestamp、nonce三个参数进行字典序排序
	var tempArray = []string{token, timestamp, nonce}
	sort.Strings(tempArray)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	var sha1String string = ""
	for _, v := range tempArray {
		sha1String += v
	}
	h := sha1.New()
	h.Write([]byte(sha1String))
	sha1String = hex.EncodeToString(h.Sum([]byte("")))
	//获得加密后的字符串可与signature对比
	if sha1String == signature {
		_, err := response.Write([]byte(echostr))
		if err != nil {
			fmt.Println("响应失败。。。")
		}
	} else {
		fmt.Println("验证失败")
	}
}
func GetAccessToken() (str string, err error) {
	const host = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx870e0c515d19cde4&secret=ae4bf23de5e9fb9680d1fccfaf0fbbed"
	baseURL := host
	method := rest.Get
	request := rest.Request{
		Method:  method,
		BaseURL: baseURL,
	}
	response, err := rest.Send(request)
	if err != nil {
		fmt.Errorf("errosis%s", err.Error())
	}
	if bytes.Contains([]byte(response.Body), []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal([]byte(response.Body), &atr)
		if err != nil {
			return "", fmt.Errorf("发送get请求获取 atoken 返回数据json解析错误%s", err)
		}
		return atr.AccessToken, nil
	} else {
		return "", fmt.Errorf("发送get请求获取 atoken 返回数据json解析错误%s", err)

	}
}
