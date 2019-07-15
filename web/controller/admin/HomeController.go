package admin

import (
	"github.com/kataras/iris/view"
	"wxapp/models"
	"wxapp/web/controller"
)

type HomeController struct {
	controller.AuthController
}

func (c HomeController) Get() {
	c.Ctx.View("admin/index_v1.html")
}

// 前往登录界面
func (c *HomeController) GetLogin() {
	if c.IsLoggedIn() {
		c.Logout()
	}

	c.Ctx.ViewLayout(view.NoLayout)
	c.Ctx.View("admin/login.html")
}

// 执行登录
func (c *HomeController) PostLogin() {
	loginForm := &LoginForm{}
	if c.SendValidateErr(loginForm) == true {
		return
	}

	user, _ := models.GetByAccountAndPwd(loginForm.UserName, loginForm.Password)

	if user == nil {
		c.SendMsg(1, "账号或密码错误", nil)
		return
	}

	if user.Type != models.USER_TYPE_ADMIN {
		c.SendMsg(1, "只允许管理员账号登录", nil)
		return
	}

	if user.Status == models.USER_STATUS_UMNORMAL {
		c.SendMsg(1, "你的账号已被锁定，请联系管理员", nil)
		return
	}

	c.Session.Set(controller.SESSION_USER_ID_KEY, user.Id)

	c.SendMsg(0, "登录成功", nil)
}

// 登出
func (c *HomeController) AnyLogout() {
	if c.IsLoggedIn() {
		c.Logout()
	}
	c.Ctx.Redirect("/login")
}
