package front

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/mojocn/base64Captcha"
	"html/template"
	"strings"
	"time"
	"wxapp/models"
	"wxapp/pkg/mpcrawler"
	"wxapp/utils/conv"
	"wxapp/web/controller"
	"wxapp/pkg/seo"
	"wxapp/config"
)

type HomeController struct {
	controller.BaseController
}

// 首页
func (c *HomeController) Get() {
	//首页推荐
	ar := models.GetRecommendArticle(models.RecommendTypeHomePage, 0)
	c.Ctx.ViewData("ar", ar)
	//最新发布
	an := models.GetNewArticleList()
	c.Ctx.ViewData("an", an)
	//首页分类推荐
	ca := models.GetRecommendHomeCategoryArticle()
	c.Ctx.ViewData("ca", ca)
	links, err := models.GetLinkByCnd(&models.Link{Type: conv.Int(models.LinkTypeHome)}, 0, 0, nil)
	if err != nil {
		golog.Errorf("HomeController.Get error : %v", err)
	}
	c.Ctx.ViewData("links", links)
	randSpecial := models.GetRandSpecialList("HOME")
	c.Ctx.ViewData("randSpecial", randSpecial)
	c.Ctx.ViewData("cid", 0)
	c.Ctx.ViewData("menu", 0)
	c.Ctx.View("front/index.html")
}

func (c *HomeController) GetApp() {
	c.getApp(models.CATE_ROOT,models.CATE_APP_NAME)
}
func (c *HomeController) GetAppBy(short string) {
	c.getApp(models.CATE_APP_ID,short)
}
func (c *HomeController) GetGame() {
	c.getApp(models.CATE_ROOT,models.CATE_GAME_NAME)
}
func (c *HomeController) GetGameBy(short string) {
	c.getApp(models.CATE_GAME_ID,short)
}
func (c *HomeController) getApp(pid int64,short string) {
	sort, err := c.Ctx.URLParamInt("sort")
	if err != nil {
		sort = 0
	}
	//设置排序方式
	c.Ctx.ViewData("sort", sort)
	page, err := c.Ctx.URLParamInt("page")
	if err != nil || page <= 0 {
		page = 1
	}
	//分类信心
	cate, err := models.GetCategoryByShort(pid,short)
	if err != nil {
		golog.Errorf("HomeController.getApp GetCategoryByShort error : %v", err)
	}
	if cate == nil {
		c.Ctx.NotFound()
		return
	}
	//分类导航
	if short == models.CATE_APP_NAME || short == models.CATE_GAME_NAME {
		pid = cate.Id
	} else {
		pid = cate.Pid
	}
	cnav := models.GetCategoryNavList(pid)
	c.Ctx.ViewData("cnav", cnav)

	pcate := cate
	//Pid大于零的时候获取第一级菜单
	if cate.Pid > 0 {
		pcate, err = models.GetCategoryById(cate.Pid)
		if err != nil {
			golog.Errorf("HomeController.getApp GetCategoryById error : %v", err)
			c.Ctx.NotFound()
			return
		}
	}
	if pcate.Short == models.CATE_GAME_NAME {
		c.Ctx.ViewData("menu", 1001)
	} else {
		c.Ctx.ViewData("menu", 1000)
	}
	c.Ctx.ViewData("pcate", pcate)
	c.Ctx.ViewData("cate", cate)
	c.Ctx.ViewData("t", "app")
	//文章列表
	articles, paginator := models.GetArticleListByCategory(cate.Id, page, sort)
	c.Ctx.ViewData("articles", articles)
	c.Ctx.ViewData("paginator", paginator)
	//随机专题
	randSpecial := models.GetRandSpecialList(fmt.Sprintf("CATEGORY_%d", cate.Id))
	c.Ctx.ViewData("randSpecial", randSpecial)

	c.Ctx.View("front/category.html")
}

func (c *HomeController) GetMp() {
	c.getMp(models.CATE_ROOT,models.CATE_MP_NAME)
}
func (c *HomeController) GetMpBy(short string) {
	c.getMp(models.CATE_MP_ID,short)
}
func (c *HomeController) getMp(pid int64,short string){
	sort, err := c.Ctx.URLParamInt("sort")
	if err != nil {
		sort = 0
	}
	//设置排序方式
	c.Ctx.ViewData("sort", sort)
	page, err := c.Ctx.URLParamInt("page")
	if err != nil || page <= 0 {
		page = 1
	}
	//分类信心
	cate, err := models.GetCategoryByShort(pid,short)
	if err != nil {
		golog.Errorf("HomeController.getMp GetCategoryByShort error : %v", err)
	}
	if cate == nil {
		c.Ctx.NotFound()
		return
	}
	//分类导航
	if short == models.CATE_MP_NAME {
		pid = cate.Id
	} else {
		pid = cate.Pid
	}
	cnav := models.GetCategoryNavList(pid)
	c.Ctx.ViewData("cnav", cnav)

	pcate := cate
	//Pid大于零的时候获取第一级菜单
	if cate.Pid > 0 {
		pcate, err = models.GetCategoryById(cate.Pid)
		if err != nil {
			golog.Errorf("HomeController.getMp GetCategoryById error : %v", err)
			c.Ctx.NotFound()
			return
		}
	}
	c.Ctx.ViewData("menu", 1002)
	c.Ctx.ViewData("pcate", pcate)
	c.Ctx.ViewData("cate", cate)
	c.Ctx.ViewData("t", "mp")
	//公众号列表
	mps, paginator := models.GetMpListByCategory(cate.Id, page, sort)
	c.Ctx.ViewData("mps", mps)
	c.Ctx.ViewData("paginator", paginator)

	c.Ctx.View("front/categorymp.html")
}

// 详情页
func (c *HomeController) GetDetailBy(id int64) {
	c.GetAppDetailBy(id)
}

// 详情页
func (c *HomeController) GetAppDetailBy(id int64) {
	article, err := models.GetArticleById(id)
	if err != nil {
		golog.Errorf("HomeController.GetDetailBy GetArticleById id: %d error : %v", id, err)
	}
	if article == nil {
		c.Ctx.NotFound()
		return
	}
	//热门推荐
	ar := models.GetHotRecommedArticleList(id)
	c.Ctx.ViewData("ar", ar)
	//最新推荐
	an := models.GetNewRecommedArticleList(id)
	c.Ctx.ViewData("an", an)
	//别人看过的
	as := models.GetSeeRecommedArticleList(id, article.Cid)
	c.Ctx.ViewData("as", as)

	eva, err := models.GetEvaluateByAid(id)
	if err != nil {
		golog.Errorf("HomeController.GetDetailBy GetEvaluateByAid id: %d error : %v", id, err)
		c.Ctx.ViewData("hasEva", false)
	} else {
		c.Ctx.ViewData("hasEva", true)
	}
	c.Ctx.ViewData("eva", eva)
	//增加浏览次数
	article.Browse = article.Browse + 1
	_, err = models.Update(article, models.Article{Id: article.Id}, nil)
	if err != nil {
		golog.Errorf("HomeController.GetDetailBy Update id: %d error : %v", id, err)
	}

	cate, err := models.GetCategoryById(article.Cid)
	if err != nil {
		golog.Errorf("HomeController.GetDetailBy GetCategoryById id: %d error : %v", id, err)
	}
	c.Ctx.ViewData("cate", cate)
	c.Ctx.ViewData("cid", article.Cid)

	pcate, err := models.GetCategoryById(cate.Pid)
	if err != nil {
		golog.Errorf("HomeController.GetBy GetCategoryById error : %v", err)
		c.Ctx.NotFound()
		return
	}
	c.Ctx.ViewData("pcate", pcate)

	if pcate.Short == "game" {
		c.Ctx.ViewData("menu", 1001)
	} else {
		c.Ctx.ViewData("menu", 1000)
	}
	c.Ctx.ViewData("t", "app")
	c.Ctx.ViewData("article", article)
	c.Ctx.View("front/detail.html")
}

// 公账号详情页
func (c *HomeController) GetMpDetailBy(id int64) {
	mp, err := models.GetMpById(id)
	if err != nil {
		golog.Errorf("HomeController.GetMpById GetMpById id: %d error : %v", id, err)
	}
	if mp == nil {
		c.Ctx.NotFound()
		return
	}
	//热门推荐
	ar := models.GetHotRecommedMpList(id)
	c.Ctx.ViewData("ar", ar)
	//最新推荐
	an := models.GetNewRecommedMpList(id)
	c.Ctx.ViewData("an", an)
	//别人看过的
	as := models.GetSeeRecommedMpList(id, mp.Cid)
	c.Ctx.ViewData("as", as)

	//增加浏览次数
	mp.Browse = mp.Browse + 1
	_, err = models.Update(mp, models.Mp{Id: mp.Id}, nil)
	if err != nil {
		golog.Errorf("HomeController.GetMpById Update id: %d error : %v", id, err)
	}

	cate, err := models.GetCategoryById(mp.Cid)
	if err != nil {
		golog.Errorf("HomeController.GetMpById GetCategoryById id: %d error : %v", id, err)
	}
	c.Ctx.ViewData("t", "mp")
	c.Ctx.ViewData("cate", cate)

	pcate, err := models.GetCategoryById(cate.Pid)
	if err != nil {
		golog.Errorf("HomeController.GetBy GetCategoryById error : %v", err)
		c.Ctx.NotFound()
		return
	}
	c.Ctx.ViewData("pcate", pcate)

	c.Ctx.ViewData("cid", mp.Cid)
	c.Ctx.ViewData("menu", 1002)
	c.Ctx.ViewData("mp", mp)
	c.Ctx.View("front/mp.html")
}

//搜索
func (c *HomeController) GetSearch() {
	q := c.Ctx.URLParam("q")
	t := c.Ctx.URLParam("t")
	page, err := c.Ctx.URLParamInt("page")
	if err != nil || page <= 0 {
		page = 1
	}

	c.search(q,t,page)
}

func (c *HomeController) PostSearch() {
	q := c.Ctx.PostValueDefault("q", "")
	t := c.Ctx.PostValueDefault("t", "")
	c.search(q,t,0)
}

func (c *HomeController) search(q,t string, page int)  {
	//默认app
	if len(t) <= 0 {
		t = "app"
	}
	c.Ctx.ViewData("q", q)
	c.Ctx.ViewData("t", t)
	c.Ctx.ViewData("menu", -1)
	q = strings.Replace(q, " ", "", -1)
	if len(q) <= 0 {
		paginator := models.Paginator(0, models.CategoryArticleSize, 0)
		c.Ctx.ViewData("paginator", paginator)
		c.Ctx.View("front/search.html")
		return
	}
	var articles, paginator interface{}
	switch t {
	case "mp":
		articles, paginator = models.SearchMp(q, page)
	case "eva":
		articles, paginator = models.SearchEvaluate(q, page)
	default:
		articles, paginator = models.SearchArticle(q, page)
	}
	c.Ctx.ViewData("articles", articles)
	c.Ctx.ViewData("paginator", paginator)
	c.Ctx.View("front/search.html")
}

func (c *HomeController) GetLinks() {
	c.Ctx.ViewData("menu", -1)

	links, err := models.GetLinkByCnd(&models.Link{Type: conv.Int(models.LinkTypePage)}, 0, 0, nil)
	if err != nil {
		golog.Errorf("HomeController.GetLinks error : %v", err)
	}
	c.Ctx.ViewData("links", links)
	c.Ctx.View("front/links.html")
}

func (c *HomeController) GetBusiness() {
	c.Ctx.ViewData("menu", -1)

	c.Ctx.View("front/business.html")
}

func (c *HomeController) GetApply() {
	//分类导航
	cnav := models.GetCategoryNavList(0)
	c.Ctx.ViewData("cnav", cnav)

	c.Ctx.ViewData("menu", -1)
	province, err := models.GetCityByPid(0)
	if err != nil {
		golog.Errorf("HomeController.GetApply error : %v", err)
	}
	c.Ctx.ViewData("province", province)

	//获取验证码
	idKey, base64string := c.generateCaptcha()
	c.Ctx.ViewData("CaptchaId", idKey)
	c.Ctx.ViewData("Captcha", base64string)

	c.Ctx.View("front/apply.html")
}

func (c *HomeController) GetApplyBy(id int64) {
	//分类导航
	cnav := models.GetCategoryNavList(0)
	c.Ctx.ViewData("cnav", cnav)

	c.Ctx.ViewData("menu", -1)
	apply, err := models.GetApplyById(id)
	if err != nil {
		golog.Errorf("HomeController.GetApplyBy GetApplyById error : %v", err)
	}
	if apply == nil {
		c.Ctx.NotFound()
		return
	}
	province, err := models.GetCityByPid(0)
	if err != nil {
		golog.Errorf("HomeController.GetApplyBy GetCityByCnd error : %v", err)
	}
	c.Ctx.ViewData("province", province)
	if apply.CityId != 0 {
		city, _ := models.GetCityByPid(apply.ProvideId)
		c.Ctx.ViewData("city", city)
	}
	c.Ctx.ViewData("province", province)
	c.Ctx.ViewData("apply", apply)
	c.Ctx.ViewData("t", "app")
	//获取验证码
	idKey, base64string := c.generateCaptcha()
	c.Ctx.ViewData("CaptchaId", idKey)
	c.Ctx.ViewData("Captcha", base64string)

	c.Ctx.View("front/applyedit.html")
}

func (c *HomeController) PostApply() {
	form := &ArticleForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	if form.Cid <= 0 {
		c.SendMsg(1, "请选择类目", nil)
		return
	}
	verifyResult := base64Captcha.VerifyCaptcha(form.CaptchaId, form.Captcha)
	if !verifyResult {
		c.SendMsg(1, "验证码错误", nil)
		return
	}

	tapply := &models.Apply{Title: form.Title}
	has, err := models.Get(tapply)
	if err != nil {
		golog.Errorf("PostApply PostSave Get error : %v", err)
		c.SendMsg(1, "保存失败", nil)
		return
	}
	apply := c.form2Apply(form)
	apply.CreateTime = time.Now()
	apply.Pass = conv.Int(models.PassTypeCommon)
	if len(apply.Screenshot) < 1 {
		c.SendMsg(1, "必须上传1张截图", nil)
		return
	}
	if form.Id > 0 {
		if has == false || tapply.Id == apply.Id {
			_, err = models.Update(apply, models.Apply{Id: form.Id}, nil)
		} else {
			c.SendMsg(1, "已有相同名称的小程序", nil)
			return
		}
	} else {
		if has {
			c.SendMsg(1, "已有相同名称的小程序", nil)
			return
		}
		_, err = models.Insert(apply, nil)
	}
	if err == nil {
		c.SendMsg(0, "保存成功", nil)
		return
	}
	if err != nil {
		golog.Errorf("PostApply PostSave error : %v", err)
	}
	c.SendMsg(1, "申请失败", nil)
}

//申请收录公众号
func (c *HomeController) GetMpApply() {
	//分类导航
	cnav := models.GetCategoryNavList(0)
	c.Ctx.ViewData("cnav", cnav)

	c.Ctx.ViewData("menu", -1)
	province, err := models.GetCityByPid(0)
	if err != nil {
		golog.Errorf("HomeController.GetMpApply error : %v", err)
	}
	c.Ctx.ViewData("province", province)
	c.Ctx.ViewData("t", "mp")
	//获取验证码
	idKey, base64string := c.generateCaptcha()
	c.Ctx.ViewData("CaptchaId", idKey)
	c.Ctx.ViewData("Captcha", base64string)

	c.Ctx.View("front/applymp.html")
}

//保存公众号
func (c *HomeController) PostMpApply() {
	form := &MpForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	if form.Cid <= 0 {
		c.SendMsg(1, "请选择类目", nil)
		return
	}
	verifyResult := base64Captcha.VerifyCaptcha(form.CaptchaId, form.Captcha)
	if !verifyResult {
		c.SendMsg(1, "验证码错误", nil)
		return
	}
	mpInfo := mpcrawler.CrawlerInfo(form.Url)
	if mpInfo.IsNull() {
		c.SendMsg(1, "公众号文章地址非法", nil)
		return
	}
	tapply := &models.Mp{UserName: mpInfo.UserName}
	has, err := models.Get(tapply)
	if err != nil {
		golog.Errorf("PostMpApply Get Get error : %v", err)
		c.SendMsg(1, "保存失败", nil)
		return
	}
	apply := c.form2MpApply(tapply, mpInfo, form)
	if has {
		_, err = models.Update(apply, models.Mp{Id: apply.Id}, nil)
	} else {
		apply.Status = conv.Int(models.MpStatusGrounding)
		_, err = models.Insert(apply, nil)
	}
	if err == nil {
		seo.PushBaidu([]string{fmt.Sprintf("%s/mp/detail/%d", config.AppUrl, apply.Id)})
		c.SendMsg(0, "收录成功", nil)
		return
	}
	if err != nil {
		golog.Errorf("PostMpApply error : %v", err)
	}
	c.SendMsg(1, "申请失败", nil)
}

func (c *HomeController) GetEvaluates() {
	page, err := c.Ctx.URLParamInt("page")
	if err != nil || page <= 0 {
		page = 1
	}
	//最新推荐
	na := models.GetNewRecommedArticleList(0)
	c.Ctx.ViewData("na", na)
	//最热的评测
	eh := models.GetEvaluateHotList()
	c.Ctx.ViewData("eh", eh)
	//文章列表
	evaluates, paginator := models.GetEvaluateList(page)
	c.Ctx.ViewData("evaluates", evaluates)
	c.Ctx.ViewData("paginator", paginator)
	c.Ctx.ViewData("t", "eva")
	c.Ctx.ViewData("menu", 4)
	c.Ctx.View("front/evaluates.html")
}

func (c *HomeController) GetEvaluateBy(id int64) {
	evaluate, err := models.GetEvaluateById(id)
	if err != nil {
		golog.Errorf("HomeController.GetEvaluateBy GetEvaluateById error : %v", err)
	}
	if evaluate == nil {
		c.Ctx.NotFound()
		return
	}
	//最新推荐
	na := models.GetNewRecommedArticleList(0)
	c.Ctx.ViewData("na", na)
	//最热的评测
	eh := models.GetEvaluateHotList()
	c.Ctx.ViewData("eh", eh)
	//增加浏览次数
	evaluate.Browse = evaluate.Browse + 1
	_, err = models.Update(evaluate, models.Evaluate{Id: evaluate.Id}, nil)
	if err != nil {
		golog.Errorf("HomeController.GetEvaluateBy Update id: %d error : %v", id, err)
	}
	//获取关联文章
	article, err := models.GetArticleById(evaluate.Aid)
	if err != nil {
		golog.Errorf("HomeController.GetEvaluateBy GetArticleById error : %v", err)
	}
	if article != nil {
		evaluate.Article = article
	}
	c.Ctx.ViewData("evaluate", evaluate)
	c.Ctx.ViewData("t", "eva")
	c.Ctx.ViewData("menu", 4)
	c.Ctx.View("front/evaluate.html")
}

// 获取下级城市
func (c *HomeController) GetCityBy(id int64) {
	city, err := models.GetCityByPid(id)
	if err != nil {
		golog.Errorf("HomeController.GetCityBy error : %v", err)
	}
	c.SendMsgData(city)
}

//获取所有专题
func (c *HomeController) GetSpecials() {
	page, err := c.Ctx.URLParamInt("page")
	if err != nil || page <= 0 {
		page = 1
	}

	specials, paginator := models.GetSpecialList(page)
	c.Ctx.ViewData("specials", specials)
	c.Ctx.ViewData("paginator", paginator)
	c.Ctx.ViewData("menu", 2)
	c.Ctx.View("front/specials.html")
}

//专题
func (c *HomeController) GetSpecialBy(id int64) {
	c.Ctx.ViewData("menu", 2)

	special := models.GetSpecialArticleById(id)
	c.Ctx.ViewData("special", special)
	//随机推荐专题
	randSpecial := models.GetRandSpecialList(fmt.Sprintf("SPECIAL_%d", id))
	c.Ctx.ViewData("randSpecial", randSpecial)

	c.Ctx.View("front/special.html")
}

func (c *HomeController) GetRank() {
	articles := models.GetArticleRank()
	c.Ctx.ViewData("top", articles[:3])
	c.Ctx.ViewData("top1", articles[3:])
	c.Ctx.ViewData("menu", 3)
	c.Ctx.View("front/rank.html")
}

//清除缓存
func (c *HomeController) GetClearcache() {
	models.ClearCache()
}

// 获取验证码
func (c *HomeController) GetCaptcha() {
	idKey, base64string := c.generateCaptcha()
	c.SendMsgData(iris.Map{"CaptchaId": idKey, "Captcha": base64string})
}

//生成验证码
func (c *HomeController) generateCaptcha() (string, template.URL) {
	//字符,公式,验证码配置
	var config = base64Captcha.ConfigCharacter{
		Height: 32,
		Width:  100,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    true,
		IsShowSlimeLine:    true,
		IsUseSimpleFont:    true,
		IsShowSineLine:     false,
		CaptchaLen:         4,
	}
	//创建字符公式验证码.
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uuid.
	idKey, cap := base64Captcha.GenerateCaptcha("", config)
	//以base64编码
	base64string := base64Captcha.CaptchaWriteToBase64Encoding(cap)

	return idKey, template.URL(base64string)
}

func (c *HomeController) form2Apply(form *ArticleForm) *models.Apply {
	apply := &models.Apply{}
	apply.Id = form.Id
	apply.Cid = form.Cid
	apply.Author = form.Author
	apply.CityId = form.City
	apply.Content = form.Content
	apply.PicCover = form.PicCover
	apply.PicQrcode = form.PicQrcode
	apply.ProvideId = form.Province
	apply.Qq = form.Qq
	apply.Title = form.Title
	apply.Keywords = form.Keywords
	apply.Screenshot = form.Screenshot
	return apply
}

func (c *HomeController) form2MpApply(mp *models.Mp, mpInfo *mpcrawler.MpCrawlInfo, form *MpForm) *models.Mp {
	mp.Cid = form.Cid
	mp.Qq = form.Qq
	mp.Keywords = form.Keywords
	mp.MsgLink = mpInfo.MsgLink
	mp.UserName = mpInfo.UserName
	mp.PicCover = mpInfo.OriHeadImgUrl
	mp.PicQrcode = mpInfo.QrImgUrl
	mp.Title = mpInfo.Nickname
	mp.Content = mpInfo.Content
	mp.ShowId = mpInfo.Id
	mp.SearchWords = fmt.Sprintf("%s %s",mp.Title,mp.Keywords)
	return mp
}
