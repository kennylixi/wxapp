package admin

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/utils/json"
	"wxapp/web/controller"
)

type SpecialController struct {
	controller.AuthController
}

func (c *SpecialController) Get() {
	c.Ctx.View("admin/special/special.html")
}
func (c *SpecialController) GetList() {
	params := c.Ctx.URLParams()
	pSpecial := &models.Special{}
	conv.FillStructStr(params, pSpecial, false)
	specials, _ := models.GetSpecialByCnd(pSpecial, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.Count(pSpecial)
	c.SendMsgData(iris.Map{"rows": specials, "total": total})
}

func (c *SpecialController) GetAdd() {
	c.Ctx.View("admin/special/add.html")
}

func (c *SpecialController) GetEditBy(id int64) {
	special, err := models.GetSpecialById(id)
	if err != nil {
		golog.Errorf("SpecialController GetEditBy GetSpecialById error : %v", err)
	}
	c.Ctx.ViewData("special", special)
	c.Ctx.View("admin/special/edit.html")
}

func (c *SpecialController) PostSave() {
	form := &SpecialForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	special := c.form2Special(form)
	_, err := models.Insert(special, nil)
	if err != nil {
		golog.Errorf("SpecialController PostSave Insert error : %v", err)
		c.SendMsg(1, "创建专题失败", nil)
		return
	}
	c.SendMsg(0, "保存成功", nil)
}

func (c *SpecialController) PostUpdate() {
	form := &SpecialForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	special := c.form2Special(form)
	_, err := models.Update(special, &models.Special{Id: special.Id}, nil)
	if err != nil {
		golog.Errorf("SpecialController PostUpdate Update error : %v", err)
		c.SendMsg(1, "更新失败", nil)
		return
	}
	c.SendMsg(0, "更新成功", nil)
}

func (c *SpecialController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的专题ID", nil)
		return
	}
	has, _ := models.Get(&models.Special{Id: id})
	if !has {
		c.SendMsg(1, "专题不存在", nil)
		return
	}
	_, err = models.Delete(&models.Special{Id: id}, nil)
	if err != nil {
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *SpecialController) PostRemoves() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要删除的专题", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		special := &models.Special{Id: conv.Int64(id)}
		has, _ := models.Get(special)
		if !has {
			continue
		}
		_, err = models.Delete(special, nil)
		if err != nil {
			continue
		}
	}

	c.SendMsg(0, "删除成功", nil)
}

func (c *SpecialController) GetArticle() {
	special, err := models.GetSpecialByCnd(&models.Special{}, 0, 0, nil)
	if err != nil {
		golog.Errorf("SpecialController GetArticle error : %v", err)
	}
	c.Ctx.ViewData("special", special)
	c.Ctx.View("admin/special/article.html")
}

func (c *SpecialController) GetArticleList() {
	params := c.Ctx.URLParams()
	id := conv.Int64(params["id"])
	specials, _ := models.GetArticleBySpecialId(id, 0)
	c.SendMsgData(specials)
}

func (c *SpecialController) form2Special(form *SpecialForm) *models.Special {
	special := &models.Special{}
	special.Desc = form.Desc
	special.Id = form.Id
	special.Pic = form.Pic
	special.Title = form.Title
	return special
}
