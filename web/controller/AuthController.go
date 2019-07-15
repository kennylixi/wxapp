package controller

type AuthController struct {
	BaseController
}

const SESSION_USER_ID_KEY = "SessionUserIdKey"

func (c *AuthController) CurrentUserId() int64 {
	userId := c.Session.GetInt64Default(SESSION_USER_ID_KEY, 0)
	return userId
}

func (c *AuthController) IsLoggedIn() bool {
	return c.CurrentUserId() > 0
}

func (c *AuthController) Logout() {
	c.Session.Destroy()
}
