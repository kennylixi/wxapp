package admin

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"time"
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/utils/json"
	"wxapp/web/controller"
)

type UserController struct {
	controller.AuthController
}

func (c *UserController) Get() {
	c.Ctx.View("admin/user/user.html")
}
func (c *UserController) GetList() {
	params := c.Ctx.URLParams()
	pUser := &models.User{}
	conv.FillStructStr(params, pUser, false)
	users, _ := models.GetUserByCnd(pUser, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.Count(pUser)
	c.SendMsgData(iris.Map{"rows": users, "total": total})
}

func (c *UserController) GetAdd() {
	c.Ctx.View("admin/user/add.html")
}

func (c *UserController) GetEditBy(id int64) {
	user, err := models.GetUserById(id)
	if err != nil {
		golog.Errorf("UserController GetEditBy error : %v", err)
	}
	c.Ctx.ViewData("user", user)
	c.Ctx.View("admin/user/edit.html")
}

func (c *UserController) PostSave() {
	userForm := &UserAddForm{}
	if c.SendValidateErr(userForm) == true {
		return
	}
	has, _ := models.Get(&models.User{Account: userForm.Account})
	if has {
		c.SendMsg(1, "该账号已存在", nil)
		return
	}
	user := c.addForm2User(userForm)
	user.Created = time.Time{}
	_, err := models.Insert(user, nil)
	if err != nil {
		golog.Errorf("UserController PostSave error : %v", err)
		c.SendMsg(1, "创建账号失败", nil)
		return
	}
	c.SendMsg(0, "保存成功", nil)
}

func (c *UserController) PostUpdate() {
	userForm := &UserEditForm{}
	if c.SendValidateErr(userForm) == true {
		return
	}
	user := c.editForm2User(userForm)
	_, err := models.Update(user, &models.User{Id: user.Id}, nil)
	if err != nil {
		golog.Errorf("UserController PostUpdate error : %v", err)
		c.SendMsg(1, "更新失败", nil)
		return
	}
	c.SendMsg(0, "更新成功", nil)
}

func (c *UserController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的用户ID", nil)
		return
	}
	has, _ := models.Get(&models.User{Id: id})
	if !has {
		c.SendMsg(1, "用户不存在", nil)
		return
	}
	_, err = models.Delete(&models.User{Id: id}, nil)
	if err != nil {
		golog.Errorf("UserController PostRemove error : %v", err)
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *UserController) PostRemoves() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要删除的用户", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		user := &models.User{Id: conv.Int64(id)}
		has, _ := models.Get(user)
		if !has {
			continue
		}
		_, err = models.Delete(user, nil)
		if err != nil {
			continue
		}
	}

	c.SendMsg(0, "删除成功", nil)
}

func (c *UserController) GetResetpwdBy(id int64) {
	c.Ctx.ViewData("Id", id)
	c.Ctx.View("admin/user/reset_pwd.html")
}

func (c *UserController) PostResetpwd() {
	form := &ResetPwdForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	user := &models.User{Id: form.Id}
	has, _ := models.Get(user)
	if !has {
		c.SendMsg(1, "用户不存在", nil)
		return
	}
	user.Pwd = form.Pwd
	_, err := models.Update(user, &models.User{Id: form.Id}, nil)
	if err != nil {
		golog.Errorf("UserController PostResetpwd error : %v", err)
		c.SendMsg(1, "密码重置失败", nil)
		return
	}
	c.SendMsg(0, "密码重置成功", nil)
}

func (c *UserController) GetRechargeBy(id int64) {
	c.Ctx.ViewData("Id", id)
	c.Ctx.View("admin/user/recharge.html")
}

func (c *UserController) PostRecharge() {
	var err error
	var id int64
	var score int
	id, err = c.Ctx.PostValueInt64("Id")
	if err != nil || id <= 0 {
		c.SendMsg(1, "非法的用户ID", nil)
		return
	}
	score, err = c.Ctx.PostValueInt("Score")
	if err != nil {
		c.SendMsg(1, "积分必须为数字", nil)
		return
	}
	err = models.Recharge(id, score)
	if err != nil {
		golog.Errorf("UserController PostRecharge error : %v", err)
		c.SendMsg(1, "充值失败", nil)
		return
	}

	c.SendMsg(0, "充值成功", nil)
}

func (c *UserController) addForm2User(form *UserAddForm) *models.User {
	user := &models.User{}
	user.Account = form.Account
	user.Name = form.Name
	user.Pwd = form.Pwd
	user.Qq = form.Qq
	user.Score = form.Score
	user.Type = form.Type
	user.Wx = form.Wx
	user.Status = form.Status
	return user
}

func (c *UserController) editForm2User(form *UserEditForm) *models.User {
	user := &models.User{}
	user.Id = form.Id
	user.Account = form.Account
	user.Name = form.Name
	user.Qq = form.Qq
	user.Score = form.Score
	user.Type = form.Type
	user.Wx = form.Wx
	user.Status = form.Status
	return user
}
