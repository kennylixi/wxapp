package admin

import (
	"github.com/kataras/golog"
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/utils/json"
	"wxapp/web/controller"
	"github.com/emirpasic/gods/lists/arraylist"
)

type CategoryController struct {
	controller.BaseController
}

func (c *CategoryController) Get() {
	c.Ctx.View("admin/category/category.html")
}
func (c *CategoryController) GetList() {
	params := c.Ctx.URLParams()
	pCategory := &models.Category{}
	conv.FillStructStr(params, pCategory, false)
	category, _ := models.GetCategoryByCnd(pCategory, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	c.SendMsgData(category)
}

func (c *CategoryController) GetAddBy(id int64) {
	c.Ctx.ViewData("Pid", id)
	if id == 0 {
		c.Ctx.ViewData("Pname", "无")
	} else {
		cate, err := models.GetCategoryById(id)
		if err != nil {
			golog.Errorf("CategoryController GetAddBy GetCategoryById error : %v", err)
		}
		if cate == nil {
			c.Ctx.ViewData("Pname", "无")
		} else {
			c.Ctx.ViewData("Pname", cate.Title)
		}
	}
	c.Ctx.View("admin/category/add.html")
}

func (c *CategoryController) GetEditBy(id int64) {
	cate, err := models.GetCategoryById(id)
	if err != nil {
		golog.Errorf("CategoryController GetEditBy GetCategoryById error : %v", err)
	}
	if cate == nil || cate.Pid == 0 {
		c.Ctx.ViewData("Pname", "无")
	} else {
		c.Ctx.ViewData("Pname", cate.Title)
	}
	c.Ctx.ViewData("category", cate)
	c.Ctx.View("admin/category/edit.html")
}

func (c *CategoryController) PostSave() {
	form := &CategoryForm{}
	if c.SendValidateErr(form) == true {
		return
	}

	category := c.form2Category(form)
	_, err := models.Insert(category, nil)
	if err != nil {
		c.SendMsg(1, "添加类目失败", nil)
		return
	}
	c.SendMsg(0, "保存成功", nil)
}

func (c *CategoryController) PostUpdate() {
	form := &CategoryForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	category := c.form2Category(form)
	_, err := models.Update(category, &models.Category{Id: category.Id}, nil)
	if err != nil {
		c.SendMsg(1, "更新失败", nil)
		return
	}
	c.SendMsg(0, "更新成功", nil)
}

func (c *CategoryController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的类目ID", nil)
		return
	}
	has, _ := models.Get(&models.Category{Id: id})
	if !has {
		c.SendMsg(1, "类目不存在", nil)
		return
	}
	if count := models.Count(&models.Category{Pid: id}); count > 0 {
		c.SendMsg(1, "包含下级类目，请先删除下级类目", nil)
		return
	}
	_, err = models.Delete(&models.Category{Id: id}, nil)
	if err != nil {
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *CategoryController) GetTree() {
	params := c.Ctx.URLParams()
	ids := params["ids"]
	jids, err := json.DecodeToJson(conv.Bytes(ids))
	if err!=nil {
		c.SendMsg(1, "ids参数格式错误", nil)
	}
	idList := arraylist.New()
	for _, id := range jids.ToArray() {
		idList.Add(conv.String(id))
	}
	c.SendMsgData(models.GetCategoryTreeByIds(idList))
}

func (c *CategoryController) GetTreeview() {
	c.Ctx.View("admin/category/categoryTree.html")
}

func (c *CategoryController) form2Category(form *CategoryForm) *models.Category {
	category := &models.Category{}
	category.Id = form.Id
	category.Title = form.Title
	category.Status = form.Status
	category.Pid = form.Pid
	category.ShortTitle = form.ShortTitle
	category.Sort = form.Sort
	category.Short = form.Short
	category.SeoDesc = form.SeoDesc
	category.SeoKey = form.SeoKey
	return category
}
