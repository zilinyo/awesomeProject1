package main
import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/silenceper/wechat/v2"
	//"github.com/silenceper/wechat/v2/cache"
	//offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/sendgrid/rest"
	"github.com/gin-gonic/gin"



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
func main (){
	r := gin.Default()
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
func GetAccessToken()(str string,err error) {
	const host = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx870e0c515d19cde4&secret=ae4bf23de5e9fb9680d1fccfaf0fbbed"
	baseURL := host
	method := rest.Get
	request := rest.Request{
		Method:  method,
		BaseURL: baseURL,
	}
	response, err := rest.Send(request)
	if err != nil{
		fmt.Errorf("errosis%s",err.Error())
	}
	if bytes.Contains([]byte(response.Body), []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal([]byte(response.Body), &atr)
		if err != nil {
			return "",fmt.Errorf("发送get请求获取 atoken 返回数据json解析错误%s", err)
		}
		return atr.AccessToken,nil
	} else {
		return "",fmt.Errorf("发送get请求获取 atoken 返回数据json解析错误%s", err)

	}
}
