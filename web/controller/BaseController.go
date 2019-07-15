package controller

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"wxapp/pkg/form"
)

type BaseController struct {
	Ctx iris.Context

	Session *sessions.Session
}

func (c *BaseController) SendMsg(code int, msg string, data interface{}) {
	_, err := c.Ctx.JSON(iris.Map{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	if err != nil {
		golog.Errorf("BaseController SendMsg error : %v", err)
	}
}

func (c *BaseController) SendMsgData(data interface{}) {
	_, err := c.Ctx.JSON(data)
	if err != nil {
		golog.Errorf("BaseController SendMsgData error : %v", err)
	}
}

func (c *BaseController) SendValidateErr(f form.Form) bool {
	err := c.Ctx.ReadForm(f)
	if err != nil {
		golog.Errorf("BaseController SendValidateErr error : %v", err)
	}
	//预处理
	f.Preprocess()
	msgs := form.Validate(f, f)
	if len(msgs) > 0 {
		_, err = c.Ctx.JSON(iris.Map{
			"code": -1,
			"msg":  msgs,
		})
		if err != nil {
			golog.Errorf("Form conversion error : %v", err)
		}
		return true
	}
	return false
}
