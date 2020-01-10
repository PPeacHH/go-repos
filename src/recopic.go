//基于百度云的图片识别接口，将车牌图片中相关的车牌号识别输出
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func main() {
	start := time.Now()
	handler := PlateHandler{}
	ch := make(chan string)

	appKey := "------"
	secret := "------"
	accessToken, err := handler.GetAccessToken(appKey, secret)

	if err != nil {
		log.Println("error:", err)
		return
	}
	//log.Println("获取到的accessToken:",accessToken)
	dir, _ := os.Getwd()
	if len(os.Args[1:]) < 1 {
		//加载固定路径图片
		log.Print("loading...")
		return
	} else {
		for _, picPath := range os.Args[1:] {
			picPath := filepath.Join(dir, picPath)
			go handler.GetPlate(picPath, accessToken, ch)
		}
	}
	for range os.Args[1:] {
		fmt.Println("获取到的车牌:" + <-ch)
	}
	fmt.Printf("%.2fs 耗费时间", time.Since(start).Seconds())
}

type accessTokenInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in""`
}

type WordResult struct {
	Number string `json:"number"`
}
type Data struct {
	WordsResult WordResult `json:"words_result"`
}
type PlateHandler struct {
}

func (handler *PlateHandler) GetAccessToken(appKey string, appSecret string) (accessToken string, err error) {
	//todo 添加accessToken缓存
	//accessToken是否存在，如果存在则进行expire_in判断
	//如果不存在，则请求新的accessToken
	url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id="+appKey+"&client_secret="+appSecret
	response, err := http.Get(url)
	if err != nil {
		return "-1", err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "-2", err
	}
	info := accessTokenInfo{}
	json.Unmarshal(data, &info)
	log.Print("请求accessToken返回的数据:", string(data))
	return info.AccessToken, nil
}

func (handler *PlateHandler) GetPlate(pictureUrl string, accessToken string, ch chan<- string) {
	data, error := ioutil.ReadFile(pictureUrl)
	if error != nil {
		ch <- fmt.Sprint(error)
	}
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(data)

	toUrl := "https://aip.baidubce.com/rest/2.0/ocr/v1/license_plate?access_token="+accessToken

	values := url.Values{}
	values.Add("image", sourcestring)
	values.Add("multi_detect", "false")
	rsp2, err := http.PostForm(toUrl, values)
	defer rsp2.Body.Close()
	if err != nil {
		log.Fatal(err)
		ch <- fmt.Sprint(error)
	}
	result, error := ioutil.ReadAll(rsp2.Body)
	if error != nil {
		log.Fatal(error)
		ch <- fmt.Sprint(error)
	}
	//log.Println("请求车牌返回的数据:",string(result))
	m := Data{}
	err = json.Unmarshal(result, &m)
	if err != nil {
		log.Fatal(err)
		ch <- fmt.Sprint(err)
	}
	ch <- m.WordsResult.Number
}
