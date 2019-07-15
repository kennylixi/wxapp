package seo

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"io/ioutil"
	"net/http"
	"strings"
	"wxapp/config"
	"wxapp/utils/conv"
)

//主动向百度推送Url
func PushBaidu(urls []string) bool {
	//线上版本才推送
	if !config.ProdMode {
		return false
	}
	bdurl := fmt.Sprintf("http://data.zz.baidu.com/urls?site=%s&token=%s", config.Domain, config.SeoBaiduPushToken)
	resp, err := http.Post(bdurl, "application/x-www-form-urlencoded", strings.NewReader(strings.Join(urls, "\n")))
	if err != nil {
		golog.Errorf("PushBaidu http.Post error : %v", err)
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		golog.Errorf("PushBaidu ioutil.ReadAll error : %v", err)
		return false
	}
	var bodyData map[string]interface{}
	if err = json.Unmarshal(body, &bodyData); err == nil {
		return conv.Int(bodyData["success"]) > 0
	} else {
		golog.Errorf("PushBaidu json.Unmarshal error : %v", err)
		return false
	}
	return true
}
