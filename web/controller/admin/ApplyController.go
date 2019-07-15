package admin

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"time"
	"wxapp/config"
	"wxapp/models"
	"wxapp/pkg/seo"
	"wxapp/utils/conv"
	"wxapp/web/controller"
)

type ApplyController struct {
	controller.AuthController
}

func (c *ApplyController) Get() {
	c.Ctx.View("admin/apply/apply.html")
}
func (c *ApplyController) GetList() {
	params := c.Ctx.URLParams()
	pApply := &models.Apply{}
	pass, ok := params["pass"]
	if ok {
		pApply.Pass = conv.Int(pass)
	}
	conv.FillStructStr(params, pApply, false)
	applys, _ := models.GetApplyByCnd(pApply, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.GetApplyByCndCount(pApply)
	c.SendMsgData(iris.Map{"rows": applys, "total": total})
}

func (c *ApplyController) GetAdd() {
	city, _ := models.GetCityByPid(0)
	c.Ctx.ViewData("province", city)
	c.Ctx.View("admin/apply/add.html")
}

func (c *ApplyController) GetEditBy(id int64) {
	apply, err := models.GetApplyById(id)
	if err != nil {
		golog.Errorf("ApplyArticelService GetEditBy error : %v", err)
	}
	province, _ := models.GetCityByPid(0)
	if apply.CityId != 0 {
		city, _ := models.GetCityByPid(apply.ProvideId)
		c.Ctx.ViewData("city", city)
	}
	cate, err := models.GetCategoryById(apply.Cid)
	if err != nil {
		golog.Errorf("ApplyArticelService GetEditBy GetCategoryById error : %v", err)
	}
	c.Ctx.ViewData("Cname", cate.Title)
	c.Ctx.ViewData("province", province)
	c.Ctx.ViewData("apply", apply)
	c.Ctx.View("admin/apply/edit.html")
}

func (c *ApplyController) PostSave() {
	form := &ArticleForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	if form.Cid <= 0 {
		c.SendMsg(1, "请选择类目", nil)
		return
	}
	apply := c.form2Apply(form)
	apply.Uid = c.CurrentUserId()
	apply.CreateTime = time.Now()
	apply.Pass = conv.Int(models.PassTypeCommon)
	if len(apply.Screenshot) < 1 {
		c.SendMsg(1, "必须上传1张截图", nil)
		return
	}
	_, err := models.Insert(apply, nil)
	if err == nil {
		c.SendMsg(0, "保存成功", nil)
		return
	}
	if err != nil {
		golog.Errorf("ApplyArticelService PostSave error : %v", err)
	}
	c.SendMsg(1, "创建文章失败", nil)
}

func (c *ApplyController) PostUpdate() {
	form := &ArticleForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	apply := c.form2Apply(form)
	if len(apply.Screenshot) < 1 {
		c.SendMsg(1, "必须上传1张截图", nil)
		return
	}
	apply.Pass = conv.Int(models.PassTypeCommon)
	_, err := models.Update(apply, models.Apply{Id: apply.Id}, nil)
	if err != nil {
		c.SendMsg(1, "更新失败", nil)
		return
	}
	c.SendMsg(0, "更新成功", nil)
}

func (c *ApplyController) GetPassBy(id int64) {
	c.Ctx.ViewData("Id", id)
	apply := &models.Apply{Id: id}
	has, _ := models.Get(apply)
	if !has {
		c.SendMsg(1, "文章不存在", nil)
		return
	}
	c.Ctx.ViewData("Pass", apply.Pass)
	c.Ctx.View("admin/apply/pass.html")
}

func (c *ApplyController) PostPass() {
	id, err := c.Ctx.PostValueInt64("Id")
	if err != nil {
		c.SendMsg(1, "非法的文章ID", nil)
		return
	}
	pass, err := c.Ctx.PostValueInt("Pass")
	if err != nil {
		c.SendMsg(1, "非法审核类型必须为数字", nil)
		return
	}
	if pass != 1 && pass != 2 && pass != 3 {
		c.SendMsg(1, "非法审核类型", nil)
		return
	}
	apply := &models.Apply{Id: id}
	has, _ := models.Get(apply)
	if !has {
		c.SendMsg(1, "文章不存在", nil)
		return
	}
	apply.Pass = pass
	_, err = models.ApplyPass(apply)
	if err != nil {
		c.SendMsg(1, "操作失败", nil)
		golog.Errorf("ApplyController PostPass error : %v", err)
		return
	}
	//主动推送向百度
	golog.Errorf("ApplyController PostPass : %v", fmt.Sprintf("%s/detail/%d", config.AppUrl, apply.Id))
	seo.PushBaidu([]string{fmt.Sprintf("%s/app/detail/%d", config.AppUrl, apply.Id)})
	c.SendMsg(0, "操作通过", nil)
}

func (c *ApplyController) form2Apply(form *ArticleForm) *models.Apply {
	apply := &models.Apply{}
	apply.Id = form.Id
	apply.Cid = form.Cid
	apply.CityId = form.City
	apply.Content = form.Content
	apply.PicCover = form.PicCover
	apply.PicQrcode = form.PicQrcode
	apply.ProvideId = form.Province
	apply.Qq = form.Qq
	apply.Title = form.Title
	apply.Keywords = form.Keywords
	apply.Screenshot = []string{}
	if len(form.Screenshot0) > 0 {
		apply.Screenshot = append(apply.Screenshot, form.Screenshot0)
	}
	if len(form.Screenshot1) > 0 {
		apply.Screenshot = append(apply.Screenshot, form.Screenshot1)
	}
	if len(form.Screenshot2) > 0 {
		apply.Screenshot = append(apply.Screenshot, form.Screenshot2)
	}
	if len(form.Screenshot3) > 0 {
		apply.Screenshot = append(apply.Screenshot, form.Screenshot3)
	}
	if len(form.Screenshot4) > 0 {
		apply.Screenshot = append(apply.Screenshot, form.Screenshot4)
	}
	return apply
}
