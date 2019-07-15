package admin

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"wxapp/config"
	"wxapp/models"
	"wxapp/pkg/seo"
	"wxapp/utils/conv"
	"wxapp/utils/json"
	"wxapp/web/controller"
)

type EvaluateController struct {
	controller.AuthController
}

func (c *EvaluateController) Get() {
	c.Ctx.View("admin/evaluate/evaluate.html")
}
func (c *EvaluateController) GetList() {
	params := c.Ctx.URLParams()
	pEvaluate := &models.Evaluate{}
	conv.FillStructStr(params, pEvaluate, false)
	evaluates, _ := models.GetEvaluateByCnd(pEvaluate, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.Count(pEvaluate)
	c.SendMsgData(iris.Map{"rows": evaluates, "total": total})
}

func (c *EvaluateController) GetAdd() {
	c.Ctx.View("admin/evaluate/add.html")
}

func (c *EvaluateController) GetEditBy(id int64) {
	evaluate, err := models.GetEvaluateById(id)
	if err != nil {
		golog.Errorf("EvaluateController GetEditBy error : %v", err)
	}
	c.Ctx.ViewData("evaluate", evaluate)
	c.Ctx.View("admin/evaluate/edit.html")
}

func (c *EvaluateController) PostSave() {
	form := &EvaluateForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	evaluate := c.form2Evaluate(form)
	_, err := models.Insert(evaluate, nil)
	if err != nil {
		golog.Errorf("EvaluateController PostSave error : %v", err)
		c.SendMsg(1, fmt.Sprintf("保存评测失败:v%", err), nil)
		return
	}
	//主动推送向百度
	seo.PushBaidu([]string{fmt.Sprintf("%s/evaluate/%d", config.AppUrl, evaluate.Id)})
	c.SendMsg(0, "保存成功", nil)
}

func (c *EvaluateController) PostUpdate() {
	form := &EvaluateForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	evaluate := c.form2Evaluate(form)
	_, err := models.Update(evaluate, &models.Evaluate{Id: evaluate.Id}, nil)
	if err != nil {
		golog.Errorf("EvaluateController PostUpdate error : %v", err)
		c.SendMsg(1, fmt.Sprintf("更新失败:v%", err), nil)
		return
	}
	c.SendMsg(0, "更新成功", nil)
}

func (c *EvaluateController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的评测ID", nil)
		return
	}
	has, err := models.Get(&models.Evaluate{Id: id})
	if !has {
		c.SendMsg(1, "评测不存在", nil)
		return
	}
	if err != nil {
		golog.Errorf("EvaluateController PostRemove Get error : %v", err)
		c.SendMsg(1, "删除失败", nil)
		return
	}
	_, err = models.Delete(&models.Evaluate{Id: id}, nil)
	if err != nil {
		golog.Errorf("EvaluateController PostRemove Delete error : %v", err)
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *EvaluateController) PostRemoves() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要删除的评测", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		evaluate := &models.Evaluate{Id: conv.Int64(id)}
		has, _ := models.Get(evaluate)
		if !has {
			continue
		}
		_, err = models.Delete(evaluate, nil)
		if err != nil {
			continue
		}
	}

	c.SendMsg(0, "删除成功", nil)
}

func (c *EvaluateController) form2Evaluate(form *EvaluateForm) *models.Evaluate {
	evaluate := &models.Evaluate{}
	evaluate.Id = form.Id
	evaluate.Title = form.Title
	evaluate.Content = form.Content
	evaluate.Desc = form.Desc
	evaluate.Keywords = form.Keywords
	evaluate.Score = form.Score
	evaluate.Source = form.Source
	evaluate.Aid = form.Aid
	evaluate.PicCover = form.PicCover
	evaluate.SearchWords = fmt.Sprintf("%s %s", form.Title, form.Keywords)
	return evaluate
}
