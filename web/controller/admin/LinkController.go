package admin

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/utils/json"
	"wxapp/web/controller"
)

type LinkController struct {
	controller.AuthController
}

func (c *LinkController) Get() {
	c.Ctx.View("admin/link/link.html")
}
func (c *LinkController) GetList() {
	params := c.Ctx.URLParams()
	pLink := &models.Link{}
	conv.FillStructStr(params, pLink, false)
	links, _ := models.GetLinkByCnd(pLink, conv.Int(params["limit"]), conv.Int(params["offset"]), map[string]models.OrderType{"sort": models.ORDER_TYPE_ASE})
	total := models.Count(pLink)
	c.SendMsgData(iris.Map{"rows": links, "total": total})
}

func (c *LinkController) GetAdd() {
	c.Ctx.View("admin/link/add.html")
}

func (c *LinkController) GetEditBy(id int64) {
	link, err := models.GetLinkById(id)
	if err != nil {
		golog.Errorf("LinkController GetEditBy error : %v", err)
	}
	c.Ctx.ViewData("link", link)
	c.Ctx.View("admin/link/edit.html")
}

func (c *LinkController) PostSave() {
	form := &LinkForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	link := c.form2Link(form)
	_, err := models.Insert(link, nil)
	if err != nil {
		golog.Errorf("LinkController PostSave error : %v", err)
		c.SendMsg(1, "创建友情链接失败", nil)
		return
	}
	c.SendMsg(0, "保存成功", nil)
}

func (c *LinkController) PostUpdate() {
	form := &LinkForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	link := c.form2Link(form)
	_, err := models.Update(link, &models.Link{Id: link.Id}, nil)
	if err != nil {
		golog.Errorf("LinkController PostUpdate error : %v", err)
		c.SendMsg(1, "更新失败", nil)
		return
	}
	c.SendMsg(0, "更新成功", nil)
}

func (c *LinkController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的友情链接ID", nil)
		return
	}
	has, err := models.Get(&models.Link{Id: id})
	if !has {
		c.SendMsg(1, "友情链接不存在", nil)
		return
	}
	if err != nil {
		golog.Errorf("LinkController PostRemove Get error : %v", err)
		c.SendMsg(1, "删除失败", nil)
		return
	}
	_, err = models.Delete(&models.Link{Id: id}, nil)
	if err != nil {
		golog.Errorf("LinkController PostRemove Delete error : %v", err)
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *LinkController) PostRemoves() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要删除的友情链接", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		link := &models.Link{Id: conv.Int64(id)}
		has, _ := models.Get(link)
		if !has {
			continue
		}
		_, err = models.Delete(link, nil)
		if err != nil {
			continue
		}
	}

	c.SendMsg(0, "删除成功", nil)
}

func (c *LinkController) form2Link(form *LinkForm) *models.Link {
	link := &models.Link{}
	link.Id = form.Id
	link.Url = form.Url
	link.Title = form.Title
	link.Sort = form.Sort
	link.Type = form.Type
	return link
}
