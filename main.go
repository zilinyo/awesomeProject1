package main

import (
	"awesomeProject/log"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/broadcast"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/rest"
	"github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
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
	log.Logs.Log("日志开启")
	r := gin.Default()
	//http.HandleFunc("/", checkout)
	r.GET("/", checkout)
	r.GET("/ping", Ping)
	r.POST("/userInfo", UserInfo)

	errChan := make(chan error)

	go func() {
		fmt.Println("Http Server start at port:8073")
		errChan <- r.Run(":80")
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //监听强制退出
		errChan <- fmt.Errorf("%s", <-c)
	}()
	_ = log.GetLogger().Log("异常退出", <-errChan)
	//r.Run(":80")
}

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

func UserInfo(c *gin.Context) {
	signature := c.Request.URL.RawQuery
	log.GetLogger().Log("userinfo", signature)
	c.JSON(200, gin.H{"msg": "success"})
	return
}

//群发消息
func Ping(c *gin.Context) {
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        "127.0.0.1:6379",
		Database:    1,
		MaxActive:   10,
		MaxIdle:     10,
		IdleTimeout: 60, //second
	}
	redisCache := cache.NewRedis(redisOpts)
	cfg := &offConfig.Config{
		AppID:     "wx870e0c515d19cde4",
		AppSecret: "ae4bf23de5e9fb9680d1fccfaf0fbbed",
		Token:     "token",
		//EncodingAESKey: "xxxx",
		Cache: redisCache,
	}
	oa := wc.GetOfficialAccount(cfg)
	bd := oa.GetBroadcast()
	users := &User{
		TagID:  1,
		OpenID: []string{"ojMDM6lPYLERook4WX9qVWPPY944", "ojMDM6iEvA73V1e9BO_boCLTqVts"},
	}

	text, err := bd.SendText((*broadcast.User)(users), "sssss")
	log.Logs.Log("日志开启",
		map[string]interface{}{
			"text": text,
			"err":  err,
		})

	if err != nil {
		c.JSON(400, gin.H{"test": text, "err": err})
		return
	}
	c.JSON(200, gin.H{"msg": text})
	return
}
func checkout(c *gin.Context) {
	//解析URL参数
	// token

	//accessToken, _ := GetAccessToken()
	token := "token"
	// 获取参数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
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
	log.Logs.Log("参数为", map[string]interface{}{
		"signature":  signature,
		"timestamp":  timestamp,
		"nonce":      nonce,
		"echostr":    echostr,
		"sha1String": sha1String,
	})
	//获得加密后的字符串可与signature对比
	if sha1String == signature {
		log.Logs.Log("对比成功。。。")
		c.Writer.Write([]byte(echostr))
		//c.JSON(http.StatusOK, []byte(echostr))
		return
	} else {
		log.Logs.Log("验证失败。。。")
	}
	log.Logs.Log("验证成功。。。")
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
	log.GetLogger().Log("data", response.Body)
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
