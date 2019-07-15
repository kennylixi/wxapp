package template

import (
	"github.com/kataras/iris/view"
	"html/template"
	"regexp"
	"strings"
	"time"
	"wxapp/config"
)

// 绑定模板方法
func BindTmpFunc(tmpl *view.HTMLEngine) {
	tmpl.AddFunc("getIndexValue", getIndexValue)
	tmpl.AddFunc("AppUrl", appUrl)
	tmpl.AddFunc("Html", html)
	tmpl.AddFunc("ProMode", proMode)
	tmpl.AddFunc("StaticUrl", staticUrl)
	tmpl.AddFunc("DateFormat", dateFormat)
	tmpl.AddFunc("Replace", replace)
	tmpl.AddFunc("string2Arr", string2Arr)
	tmpl.AddFunc("add", add)
	tmpl.AddFunc("isOdd", isOdd)
}

// 获取指定数组指定下标的值
func getIndexValue(arr []string, index int) string {
	if index < 0 || index >= len(arr) {
		return ""
	}
	return arr[index]
}

func appUrl() template.URL {
	return template.URL(config.AppUrl)
}

func html(str string) template.HTML {
	return template.HTML(str)
}

func proMode() bool {
	return config.ProdMode
}

func staticUrl() template.URL {
	return template.URL(config.StaticUrl)
}

//时间格式化
func dateFormat(date time.Time) string {
	return date.Format("2006-01-02")
}

//替换
func replace(str string, repl string) string {
	reg := regexp.MustCompile(" ")
	return reg.ReplaceAllString(str, repl)
}

//一个字符串以空格切割成数组
func string2Arr(str string) []string {
	tags := strings.Fields(str)
	return tags
}

//两个数相加
func add(a, b int) int {
	return a + b
}

//判断是否是奇数
func isOdd(a int) bool {
	return a%2 == 1
}
