package admin

import (
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/web/controller"
)

type CityController struct {
	controller.BaseController
}

func (c *CityController) Get() {
	c.Ctx.View("admin/city/city.html")
}
func (c *CityController) GetList() {
	params := c.Ctx.URLParams()
	pCity := &models.City{}
	conv.FillStructStr(params, pCity, false)
	city, _ := models.GetCityByCnd(pCity, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	c.SendMsgData(city)
}

// 获取下级城市
func (c *CityController) GetCityBy(id int64) {
	city, _ := models.GetCityByPid(id)
	c.SendMsgData(city)
}

func (c *CityController) GetAddBy(id int64) {
	if id == 0 {
		c.Ctx.ViewData("Pname", "中国")
	} else {
		city, err := models.GetCityById(id)
		if city == nil || err != nil {
			c.Ctx.ViewData("Pname", "中国")
		} else {
			c.Ctx.ViewData("Pname", city.Name)
		}
	}
	c.Ctx.ViewData("Pid", id)
	c.Ctx.View("admin/city/add.html")
}

func (c *CityController) GetEditBy(id int64) {
	city, err := models.GetCityById(id)
	if city.Pid == 0 || city == nil || err != nil {
		c.Ctx.ViewData("Pname", "中国")
	} else {
		c.Ctx.ViewData("Pname", city.Name)
	}
	c.Ctx.ViewData("city", city)
	c.Ctx.View("admin/city/edit.html")
}

func (c *CityController) PostSave() {
	form := &CityForm{}
	if c.SendValidateErr(form) == true {
		return
	}

	city := c.form2City(form)
	has, _ := models.Get(&models.City{Id: city.Id})
	if has == true {
		c.SendMsg(1, "城市ID已存在", nil)
		return
	}

	_, err := models.Insert(city, nil)
	if err != nil {
		c.SendMsg(1, "添加城市失败", nil)
		return
	}
	c.SendMsg(0, "保存成功", nil)
}

func (c *CityController) PostUpdate() {
	form := &CityForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	city := c.form2City(form)
	_, err := models.Update(city, &models.City{Id: city.Id}, nil)
	if err != nil {
		c.SendMsg(1, "更新失败", nil)
		return
	}
	c.SendMsg(0, "更新成功", nil)
}

func (c *CityController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的城市ID", nil)
		return
	}
	has, _ := models.Get(&models.City{Id: id})
	if !has {
		c.SendMsg(1, "城市不存在", nil)
		return
	}
	if count := models.Count(&models.City{Pid: id}); count > 0 {
		c.SendMsg(1, "包含下级城市，请先删除下级城市", nil)
		return
	}
	_, err = models.Delete(&models.City{Id: id}, nil)
	if err != nil {
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *CityController) form2City(form *CityForm) *models.City {
	city := &models.City{}
	city.Id = form.Id
	city.Name = form.Name
	city.Pid = form.Pid
	city.Group = form.Group
	return city
}
