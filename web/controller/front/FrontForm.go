package front

import (
	"strings"
	"wxapp/pkg/form"
)

// 文章表单
type ArticleForm struct {
	form.BaseForm
	Id         int64  //主键ID
	Cid        int64  //分类ID
	Province   int64  //省
	City       int64  //市
	Title      string //标题
	Author     string //主体信息
	Qq         string //QQ号
	Keywords   string //tag
	Content    string //内容
	PicCover   string //封面图片
	PicQrcode  string //二维码图片
	Screenshot []string
	CaptchaId  string //验证码ID
	Captcha    string //验证码
}

// 预处理
func (f *ArticleForm) Preprocess() {
	f.Title = strings.Trim(f.Title, " ")
	f.Author = strings.Trim(f.Author, " ")
	f.Qq = strings.Trim(f.Qq, " ")
	f.Keywords = strings.Trim(f.Keywords, " ")
	f.Content = strings.Trim(f.Content, " ")
}

// 规则
func (f *ArticleForm) Rules() map[string]string {
	return map[string]string{
		"Cid":       "required|integer",
		"Province":  "required|integer",
		"City":      "required|integer",
		"Title":     "required|length:1,100",
		"Author":    "required|length:1,100",
		"Qq":        "required|qq|length:5,20",
		"Keywords":  "length:1,100",
		"Content":   "required|length:50,500",
		"PicCover":  "required|length:1,255",
		"PicQrcode": "required|length:1,255",
	}
}

// 规则对应的效验错误描述
func (f *ArticleForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Cid": map[string]string{
			"required": "请选择分类",
			"integer":  "请选择分类",
		},
		"Province": map[string]string{
			"required": "请选择省份",
			"integer":  "请选择省份",
		},
		"City": map[string]string{
			"required": "请选择城市",
			"integer":  "请选择城市",
		},
		"Title": map[string]string{
			"required": "小程序名不能为空",
			"length":   "小程序名长度为:min到:max个字符",
		},
		"Author": map[string]string{
			"required": "主体信息不能为空",
			"length":   "主体信息长度为:min到:max个字符",
		},
		"Qq": map[string]string{
			"required": "联系QQ不能为空",
			"length":   "联系QQ长度为:min到:max个字符",
			"qq":       "请输入正确QQ",
		},
		"Keywords": map[string]string{
			"length":   "小程序标签长度为:min到:max个字符",
		},
		"Content": map[string]string{
			"required": "小程序介绍不能为空",
			"length":   "小程序介绍长度为:min到:max个字符",
		},
		"PicCover": map[string]string{
			"required": "请上传小程序图标",
			"length":   "小程序图标长度为:min到:max个字符",
		},
		"PicQrcode": map[string]string{
			"required": "请上传小程序二维码",
			"length":   "小程序二维码长度为:min到:max个字符",
		},
	}
}

// 公众号表单
type MpForm struct {
	form.BaseForm
	Url       string //公众号文章链接
	Cid       int64  //分类ID
	Qq        string //QQ号
	Keywords  string //tag
	CaptchaId string //验证码ID
	Captcha   string //验证码
}

// 预处理
func (f *MpForm) Preprocess() {
	f.Url = strings.Trim(f.Url, " ")
	f.Qq = strings.Trim(f.Qq, " ")
	f.Keywords = strings.Trim(f.Keywords, " ")
}

// 规则
func (f *MpForm) Rules() map[string]string {
	return map[string]string{
		"Cid":      "required|integer",
		"Qq":       "required|qq|length:5,20",
		"Keywords": "length:1,100",
		"Url":      "required|length:1,255",
	}
}

// 规则对应的效验错误描述
func (f *MpForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Cid": map[string]string{
			"required": "请选择分类",
			"integer":  "请选择分类",
		},
		"Qq": map[string]string{
			"required": "联系QQ不能为空",
			"length":   "联系QQ长度为:min到:max个字符",
			"qq":       "请输入正确QQ",
		},
		"Keywords": map[string]string{
			"length":   "公众号标签长度为:min到:max个字符",
		},
		"Url": map[string]string{
			"required": "公众号文章链接不能为空",
			"length":   "公众号文章链接长度为:min到:max个字符",
		},

	}
}
