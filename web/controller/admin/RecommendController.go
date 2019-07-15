package admin

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"time"
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/utils/json"
	"wxapp/web/controller"
)

type RecommendController struct {
	controller.BaseController
}

func (c *RecommendController) Get() {
	c.Ctx.View("admin/recommend/recommend.html")
}
func (c *RecommendController) GetList() {
	params := c.Ctx.URLParams()
	pRecommend := &models.RecommendArticle{}
	conv.FillStructStr(params, pRecommend, false)
	recommend, _ := models.GetRecommendByCnd(pRecommend, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.Count(pRecommend.Recommend)
	c.SendMsgData(iris.Map{"rows": recommend, "total": total})
}

func (c *RecommendController) GetAdd() {
	c.Ctx.View("admin/recommend/add.html")
}

func (c *RecommendController) PostSave() {
	form := &RecommendForm{}
	if c.SendValidateErr(form) == true {
		return
	}
	article, err := models.GetArticleById(form.Aid)
	if article == nil {
		c.SendMsg(1, "小程序ID不存在", nil)
		return
	}
	if err != nil {
		golog.Errorf("RecommendController PostSave GetArticleById error : %v", err)
	}
	form.Cid = article.Cid
	recommend := c.form2Recommend(form)
	recommend.Status = conv.Int(models.RecommendStatusNo)
	_, err = models.Insert(recommend, nil)
	if err != nil {
		golog.Errorf("RecommendController PostSave Insert error : %v", err)
		c.SendMsg(1, "保存失败", nil)
		return
	}
	c.SendMsg(0, "保存成功", nil)
}

func (c *RecommendController) PostGrounding() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的参数", nil)
		return
	}
	recommend, err := models.GetRecommendById(id)
	if recommend == nil {
		c.SendMsg(1, "指定数据不存在", nil)
		return
	}
	if err != nil {
		golog.Errorf("RecommendController PostGrounding GetRecommendById error : %v", err)
		c.SendMsg(1, "指定数据不存在", nil)
		return
	}
	if recommend.Status == conv.Int(models.RecommendStatusNo) {
		recommend.Status = conv.Int(models.RecommendStatusOn)
		duration, err := time.ParseDuration(fmt.Sprintf("%dh", recommend.Rtime*30*24))
		if err != nil {
			golog.Errorf("RecommendController PostGrounding time.ParseDuration error : %v", err)
			c.SendMsg(0, "操作失败", nil)
			return
		}
		recommend.Etime = time.Now().Add(duration)
	} else if recommend.Status == conv.Int(models.RecommendStatusOn) {
		recommend.Status = conv.Int(models.RecommendStatusOff)
	} else if recommend.Status == conv.Int(models.RecommendStatusOff) {
		recommend.Status = conv.Int(models.RecommendStatusOn)
	}
	_, err = models.Update(recommend.Recommend, &models.Recommend{Id: id}, nil)
	if err != nil {
		golog.Errorf("RecommendController PostGrounding Update error : %v", err)
		c.SendMsg(0, "操作失败", nil)
		return
	}
	c.SendMsg(0, "操作成功", nil)
}
func (c *RecommendController) PostRemoves() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要删除的推荐", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		recommend := &models.Recommend{Id: conv.Int64(id)}
		has, err := models.Get(recommend)
		if err != nil {
			golog.Errorf("RecommendController PostRemoves Get error : %v", err)
		}
		if !has {
			continue
		}
		_, err = models.Delete(recommend, nil)
		if err != nil {
			golog.Errorf("RecommendController PostRemoves Delete error : %v", err)
			continue
		}
	}

	c.SendMsg(0, "删除成功", nil)
}

func (c *RecommendController) form2Recommend(form *RecommendForm) *models.Recommend {
	recommend := &models.Recommend{}
	recommend.Cid = form.Cid
	recommend.Aid = form.Aid
	recommend.IsCost = form.IsCost
	recommend.Rtime = form.Rtime
	recommend.Type = form.Type
	return recommend
}
