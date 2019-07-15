package admin

import (
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type LoginForm struct {
	UserName string //`form:"username"`
	Password string //`form:"password"`
}

// 预处理
func (f *LoginForm) Preprocess() {

}

// 规则
func (f *LoginForm) Rules() map[string]string {
	return map[string]string{
		"UserName": "required|length:1,100",
		"Password": "required|length:1,20",
	}
}

// 规则对应的效验错误描述
func (f *LoginForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"UserName": map[string]string{
			"required": "用户名不能为空",
			"length":   "用户名长度为:min到:max个字符",
		},
		"Password": map[string]string{
			"required": "密码不能为空",
			"length":   "密码长度为:min到:max个字符",
		},
	}
}

// 用户表单
type UserAddForm struct {
	Name    string //昵称
	Account string //账号
	Pwd     string //密码
	Score   int    //积分
	Type    int    //类型
	Wx      string //微信
	Qq      string //QQ
	Status  int    //状态
}

// 预处理
func (f *UserAddForm) Preprocess() {
	f.Name = strings.Trim(f.Name, " ")
	f.Account = strings.Trim(f.Account, " ")
	f.Pwd = strings.Trim(f.Pwd, " ")
	f.Wx = strings.Trim(f.Wx, " ")
	f.Qq = strings.Trim(f.Qq, " ")
}

// 规则
func (f *UserAddForm) Rules() map[string]string {
	return map[string]string{
		"Id":      "integer",
		"Name":    "required|length:1,100",
		"Account": "required|length:1,100",
		"Pwd":     "required|length:6,20",
		"Score":   "integer",
		"Type":    "required|in:1,2",
		"Status":  "required|in:1,2",
		"Qq":      "qq|length:5,50",
		"Wx":      "length:1,50",
	}
}

// 规则对应的效验错误描述
func (f *UserAddForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Name": map[string]string{
			"required": "用户名不能为空",
			"length":   "用户名长度为:min到:max个字符",
		},
		"Account": map[string]string{
			"required": "账号不能为空",
			"length":   "账号长度为:min到:max个字符",
		},
		"Pwd": map[string]string{
			"required": "密码不能为空",
			"length":   "密码长度为:min到:max个字符",
		},
		"Score": map[string]string{
			"integer": "积分必须为数字",
		},
		"Type": map[string]string{
			"required": "用户类型不能为空",
			"in":       "用户类型非法",
		},
		"Status": map[string]string{
			"required": "用户状态不能为空",
			"in":       "用户状态非法",
		},
		"Qq": map[string]string{
			"qq":     "QQ格式不正确",
			"length": "QQ长度为:min到:max个字符",
		},
		"Wx": map[string]string{
			"length": "微信长度为:min到:max个字符",
		},
	}
}

// 用户表单
type UserEditForm struct {
	Id      int64  //Id
	Name    string //昵称
	Account string //账号
	Score   int    //积分
	Type    int    //类型
	Wx      string //微信
	Qq      string //QQ
	Status  int    //状态
}

// 预处理
func (f *UserEditForm) Preprocess() {
	f.Name = strings.Trim(f.Name, " ")
	f.Account = strings.Trim(f.Account, " ")
	f.Wx = strings.Trim(f.Wx, " ")
	f.Qq = strings.Trim(f.Qq, " ")
}

// 规则
func (f *UserEditForm) Rules() map[string]string {
	return map[string]string{
		"Id":      "integer",
		"Name":    "required|length:1,100",
		"Account": "required|length:1,100",
		"Score":   "integer",
		"Type":    "required|in:1,2",
		"Status":  "required|in:1,2",
		"Qq":      "qq|length:5,50",
		"Wx":      "length:1,50",
	}
}

// 规则对应的效验错误描述
func (f *UserEditForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Name": map[string]string{
			"required": "用户名不能为空",
			"length":   "用户名长度为:min到:max个字符",
		},
		"Account": map[string]string{
			"required": "账号不能为空",
			"length":   "账号长度为:min到:max个字符",
		},
		"Score": map[string]string{
			"integer": "积分必须为数字",
		},
		"Type": map[string]string{
			"required": "用户类型不能为空",
			"in":       "用户类型非法",
		},
		"Status": map[string]string{
			"required": "用户状态不能为空",
			"in":       "用户状态非法",
		},
		"Qq": map[string]string{
			"qq":     "QQ格式不正确",
			"length": "QQ长度为:min到:max个字符",
		},
		"Wx": map[string]string{
			"length": "微信长度为:min到:max个字符",
		},
	}
}

// 用户表单
type ResetPwdForm struct {
	Id  int64  //Id
	Pwd string //密码
}

// 预处理
func (f *ResetPwdForm) Preprocess() {
	f.Pwd = strings.Trim(f.Pwd, " ")
}

// 规则
func (f *ResetPwdForm) Rules() map[string]string {
	return map[string]string{
		"Id":  "required|integer",
		"Pwd": "required|length:6,20",
	}
}

// 规则对应的效验错误描述
func (f *ResetPwdForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Id": map[string]string{
			"required": "Id不能为空",
			"integer":  "非法的Id",
		},
		"Pwd": map[string]string{
			"required": "密码不能为空",
			"length":   "密码长度为:min到:max个字符",
		},
	}
}

// 目录表单
type CategoryForm struct {
	Title      string //类目名
	Id         int64  //类目ID
	Pid        int64  //上级ID
	Sort       int    //排序号
	Short      string //拼音
	Status     int    //是否隐藏
	ShortTitle string //简称
	SeoKey     string //seo关键字
	SeoDesc    string //seo描述
}

// 预处理
func (f *CategoryForm) Preprocess() {
	f.Title = strings.Trim(f.Title, " ")
	f.ShortTitle = strings.Trim(f.ShortTitle, " ")
	f.Short = strings.Trim(f.Short, " ")
	f.SeoKey = strings.Trim(f.SeoKey, " ")
	f.SeoDesc = strings.Trim(f.SeoDesc, " ")
}

// 规则
func (f *CategoryForm) Rules() map[string]string {
	return map[string]string{
		"Pid":        "integer",
		"Id":         "integer",
		"Title":      "required|length:1,50",
		"Sort":       "integer",
		"Status":     "required|in:1,2",
		"ShortTitle": "required|length:1,50",
		"Short":      "required|length:1,50",
		"SeoKey":     "length:0,150",
		"SeoDesc":    "length:0,150",
	}
}

// 规则对应的效验错误描述
func (f *CategoryForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Id": map[string]string{
			"integer": "ID必须为数字",
		},
		"Pid": map[string]string{
			"integer": "上级ID必须为数字",
		},
		"Title": map[string]string{
			"required": "类目名称不能为空",
			"length":   "类目名称长度为:min到:max个字符",
		},
		"Sort": map[string]string{
			"integer": "排序必须为数字",
		},
		"Status": map[string]string{
			"required": "隐藏不能为空",
			"in":       "隐藏类型非法",
		},
		"ShortTitle": map[string]string{
			"required": "类目简称不能为空",
			"length":   "类目简称长度为:min到:max个字符",
		},
		"SeoKey": map[string]string{
			"length": "SEO搜索关键字长度为:min到:max个字符",
		},
		"SeoDesc": map[string]string{
			"length": "SEO搜索描述长度为:min到:max个字符",
		},
	}
}

// 城市表单
type CityForm struct {
	Name  string //类目名
	Id    int64  //类目ID
	Pid   int64  //上级ID
	Group string //分组
}

// 预处理
func (f *CityForm) Preprocess() {
	f.Name = strings.Trim(f.Name, " ")
	f.Group = strings.Trim(f.Group, " ")
}

// 规则s
func (f *CityForm) Rules() map[string]string {
	return map[string]string{
		"Pid":  "integer",
		"Id":   "integer",
		"Name": "required|length:1,50",
	}
}

// 规则对应的效验错误描述
func (f *CityForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Id": map[string]string{
			"integer": "ID必须为数字",
		},
		"Pid": map[string]string{
			"integer": "上级ID必须为数字",
		},
		"Name": map[string]string{
			"required": "类目名称不能为空",
			"length":   "类目名称长度为:min到:max个字符",
		},
	}
}

// 专题表单
type SpecialForm struct {
	Title string //专题名
	Id    int64  //专题ID
	Desc  string //专题描述
	Pic   string //专题图片
	Cname string //类目名
}

// 预处理
func (f *SpecialForm) Preprocess() {
	f.Title = strings.Trim(f.Title, " ")
	f.Desc = strings.Trim(f.Desc, " ")
	f.Cname = strings.Trim(f.Cname, " ")
}

// 规则s
func (f *SpecialForm) Rules() map[string]string {
	return map[string]string{
		"Id":    "integer",
		"Title": "required|length:1,50",
		"Pic":   "required|length:1,255",
		"Desc":  "required|length:1,255",
	}
}

// 规则对应的效验错误描述
func (f *SpecialForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Id": map[string]string{
			"integer": "ID必须为数字",
		},
		"Title": map[string]string{
			"required": "专题名称不能为空",
			"length":   "专题名称长度为:min到:max个字符",
		},
		"Pic": map[string]string{
			"required": "专题图片不能为空",
			"length":   "专题图片长度为:min到:max个字符",
		},
		"Desc": map[string]string{
			"required": "专题描述不能为空",
			"length":   "专题描述长度为:min到:max个字符",
		},
	}
}

// 文章表单
type ArticleForm struct {
	Id          int64  //主键ID
	Cid         int64  //分类ID
	Cname       string //分类名
	Uid         int64  //用户ID
	Province    int64  //省
	City        int64  //市
	Title       string //标题
	Author      string //主体信息
	Qq          string //QQ号
	Keywords    string //tag
	Content     string //内容
	PicCover    string //封面图片
	PicQrcode   string //二维码图片
	Screenshot0 string
	Screenshot1 string
	Screenshot2 string
	Screenshot3 string
	Screenshot4 string
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
		"Title":     "required|length:1,50",
		"Author":    "required|length:1,50",
		"Qq":        "required|qq|length:5,20",
		"Keywords":  "required|length:1,255",
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
			"required": "小程序标签不能为空",
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

// 友情链接表单
type LinkForm struct {
	Title string //友情链接名
	Id    int64  //友情链接ID
	Url   string //专题Url
	Sort  int    //排序号
	Type  int    //类型
}

// 预处理
func (f *LinkForm) Preprocess() {
	f.Title = strings.Trim(f.Title, " ")
	f.Url = strings.Trim(f.Url, " ")
}

// 规则
func (f *LinkForm) Rules() map[string]string {
	return map[string]string{
		"Id":    "integer",
		"Title": "required|length:1,50",
		"Url":   "required|length:1,255",
		"Sort":  "integer",
		"Type":  "in:1,2",
	}
}

// 规则对应的效验错误描述
func (f *LinkForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Id": map[string]string{
			"integer": "ID必须为数字",
		},
		"Title": map[string]string{
			"required": "友情链接名称不能为空",
			"length":   "友情链接名称长度为:min到:max个字符",
		},
		"Url": map[string]string{
			"required": "友情链接Url不能为空",
			"length":   "友情链接Url长度为:min到:max个字符",
		},
		"Sort": map[string]string{
			"integer": "排序必须为数字",
		},
		"Type": map[string]string{
			"in": "友情链接类型非法",
		},
	}
}

// 评测
type EvaluateForm struct {
	Id       int64  //主键ID
	Title    string //类目名
	Aid      int64  //文章id
	Source   string //文章来源
	Score    int    //平分
	Content  string //内容
	Desc     string //简单描述
	Keywords string //标签
	PicCover string //图片
}

// 预处理
func (f *EvaluateForm) Preprocess() {
	f.Title = strings.Trim(f.Title, " ")
	f.Source = strings.Trim(f.Source, " ")
	policy := bluemonday.UGCPolicy()
	policy.AllowElements("h1","h2","h3","h4","h5","h6","b","i","u","strike","span","hr","p","blockquote","ol","li","ul","br","table","colgroup","col","thead","tr","th","tbody","td")
	policy.AllowAttrs("style","margin-left","text-align").OnElements("span","p")
	policy.AllowAttrs("height","width").OnElements("col")
	f.Content = policy.Sanitize(f.Content)
	f.Desc = bluemonday.UGCPolicy().Sanitize(f.Desc)
	f.Keywords = strings.Trim(f.Keywords, " ")
	f.PicCover = strings.Trim(f.PicCover, " ")
}

// 规则
func (f *EvaluateForm) Rules() map[string]string {
	return map[string]string{
		"Id":       "integer",
		"Title":    "required|length:1,250",
		"Aid":      "required|integer",
		"Source":   "required|length:1,100",
		"Score":    "required|integer",
		"Content":  "required|min-length:1",
		"Desc":     "required|length:1,500",
		"Keywords": "required|length:1,100",
		"PicCover": "length:1,255",
	}
}

// 规则对应的效验错误描述
func (f *EvaluateForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Id": map[string]string{
			"integer": "ID必须为数字",
		},
		"Title": map[string]string{
			"required": "评测名称不能为空",
			"length":   "评测名称长度为:min到:max个字符",
		},
		"Aid": map[string]string{
			"required": "应用ID不能为空",
			"integer":  "应用ID必须为数字",
		},
		"Source": map[string]string{
			"required": "来源不能为空",
			"length":   "来源长度为:min到:max个字符",
		},
		"Score": map[string]string{
			"required": "评分不能为空",
			"integer":  "评分必须为数字",
		},
		"Content": map[string]string{
			"required":   "评测名称不能为空",
			"min-length": "评测名称长至少:min个字符",
		},
		"Desc": map[string]string{
			"required": "评测简介不能为空",
			"length":   "评测简介长度为:min到:max个字符",
		},
		"Keywords": map[string]string{
			"required": "关键字不能为空",
			"length":   "关键字长度为:min到:max个字符",
		},
		"PicCover": map[string]string{
			"length": "图片长度为:min到:max个字符",
		},
	}
}

type RecommendForm struct {
	Aid    int64  //文章ID
	Cid    int64  //分类ID
	Cname  string //分类名
	Type   int    //类型（1首页推荐，2首页分类推荐，3分类推荐，4详情页推荐）
	Rtime  int    //推广时间（以一个月为单位）
	IsCost int    //是否付费推广（1没付费，2付费）
}

// 预处理
func (f *RecommendForm) Preprocess() {
	f.Cname = strings.Trim(f.Cname, " ")
}

// 规则
func (f *RecommendForm) Rules() map[string]string {
	return map[string]string{
		"Aid":    "required|integer|min:1",
		"Cid":    "integer|min:0",
		"Type":   "required|in:1,2,3,4",
		"IsCost": "required|in:1,2",
		"Rtime":  "required|integer|min:1",
	}
}

// 规则对应的效验错误描述
func (f *RecommendForm) Msgs() map[string]interface{} {
	return map[string]interface{}{
		"Aid": map[string]string{
			"required": "小程序ID不能为空",
			"integer":  "小程序ID必须为数字",
			"min":      "小程序ID必须大于:min",
		},
		"Cid": map[string]string{
			"integer": "小程序分类必须为数字",
			"min":     "小程序分类必须大于:min",
		},
		"Type": map[string]string{
			"required": "推荐类型不能为空",
			"in":       "推荐类型非法",
		},
		"IsCost": map[string]string{
			"required": "是否付费不能为空",
			"in":       "是否付费非法",
		},
		"Rtime": map[string]string{
			"required": "推广时间不能为空",
			"integer":  "推广时间必须为数字",
			"min":      "推广时间必须大于:min",
		},
	}
}
