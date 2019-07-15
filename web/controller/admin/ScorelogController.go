package admin

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/web/controller"
)

type ScorelogController struct {
	controller.AuthController
}

func (c *ScorelogController) Get() {
	c.Ctx.View("admin/scorelog/scorelog.html")
}
func (c *ScorelogController) GetList() {
	params := c.Ctx.URLParams()
	account := params["Account"]
	var uid int64
	if len(account) > 0 {
		user := &models.User{Account: account}
		has, err := models.Get(user)
		if err != nil {
			golog.Errorf("ScorelogController GetList error : %v", err)
			c.SendMsgData(iris.Map{"rows": []*models.UserScoreLog{}, "total": 0})
			return
		} else if !has {
			golog.Errorf("ScorelogController GetList error not has user by account : %v", account)
			c.SendMsgData(iris.Map{"rows": []*models.UserScoreLog{}, "total": 0})
			return
		}
		uid = user.Id
	}
	scorelogs, _ := models.GetScoreLogByUid(uid, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.Count(&models.ScoreLog{Uid: uid})
	c.SendMsgData(iris.Map{"rows": scorelogs, "total": total})
}
