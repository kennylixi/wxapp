// 公众号文章搬运工
package mpcrawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/robertkrimen/otto"
	"strings"
	"wxapp/utils/conv"
)

const (
	c_MP_QR string = "http://open.weixin.qq.com/qr/code?username=%s"
)

type MpCrawlInfo struct {
	MsgLink       string //原来的url
	UserName      string //微信公众唯一id
	OriHeadImgUrl string //微信公众号头像
	QrImgUrl      string //微信公众二维码
	Nickname      string //微信公众号名字
	Id            string //公众号id，可变的
	Content       string //微信公众号内容
}

func (m *MpCrawlInfo) IsNull() bool {
	if len(m.MsgLink) <= 0 || len(m.UserName) <= 0 || len(m.Nickname) <= 0 || len(m.Id) <= 0 {
		return true
	}
	return false
}

func CrawlerInfo(url string) *MpCrawlInfo {
	c := colly.NewCollector(
		colly.AllowedDomains("mp.weixin.qq.com"),
		colly.CacheDir("./mpcrawlercache"),
	)

	mpInfo := &MpCrawlInfo{}

	c.OnHTML("#js_profile_qrcode > div", func(e *colly.HTMLElement) {
		mpInfo.Nickname = e.ChildText("strong")
		mpInfo.Id = e.ChildText("p:nth-child(3) > span")
		mpInfo.Content = e.ChildText("p:nth-child(4) > span")
	})

	c.OnHTML("script", func(e *colly.HTMLElement) {
		vm := otto.New()
		vm.Run(e.Text)
		msgLink, _ := vm.Get("msg_link")
		if msgLink.IsUndefined() {
			return
		}
		mpInfo.MsgLink = strings.Replace(conv.String(msgLink), "&amp;", "&", -1)
		userName, _ := vm.Get("user_name")
		if userName.IsUndefined() {
			return
		}
		mpInfo.UserName = conv.String(userName)
		mpInfo.QrImgUrl = fmt.Sprintf(c_MP_QR, mpInfo.UserName)
		oriHeadImgUrl, _ := vm.Get("ori_head_img_url")
		if oriHeadImgUrl.IsUndefined() {
			return
		}
		mpInfo.OriHeadImgUrl = conv.String(oriHeadImgUrl)
	})

	c.OnHTML("#activity-detail div.page_msg div.text_area div.global_error_msg",
		func(e *colly.HTMLElement) {
			fmt.Println("error")
		})

	c.Visit(url)
	return mpInfo
}
