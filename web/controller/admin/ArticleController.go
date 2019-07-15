package admin

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"wxapp/models"
	"wxapp/utils/conv"
	"wxapp/utils/json"
	"wxapp/web/controller"
)

type ArticleController struct {
	controller.AuthController
}

func (c *ArticleController) Get() {
	c.Ctx.View("admin/article/article.html")
}
func (c *ArticleController) GetList() {
	params := c.Ctx.URLParams()
	pArticle := &models.Article{}
	status, ok := params["status"]
	if ok {
		pArticle.Status = conv.Int(status)
	}
	conv.FillStructStr(params, pArticle, false)
	articles, _ := models.GetArticleByCnd(pArticle, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.Count(pArticle)
	c.SendMsgData(iris.Map{"rows": articles, "total": total})
}

func (c *ArticleController) PostGrounding() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的文章ID", nil)
		return
	}
	article := &models.Article{Id: id}
	has, _ := models.Get(article)
	if !has {
		c.SendMsg(1, "文章不存在", nil)
		return
	}
	oldStatus := article.Status
	unGrounding := conv.Int(models.ArticleStatusUnGrounding)
	grounding := conv.Int(models.ArticleStatusGrounding)
	var operStr string
	if oldStatus == unGrounding {
		article.Status = grounding
		operStr = "上架"
	} else {
		article.Status = unGrounding
		operStr = "下架"
	}

	_, err = models.Update(article, &models.Article{Id: id}, nil)
	if err != nil {
		c.SendMsg(1, fmt.Sprintf("%s失败", operStr), nil)
		golog.Errorf("ArticleController PostGrounding error : %v", err)
		return
	}
	c.SendMsg(0, fmt.Sprintf("%s成功", operStr), nil)
}

func (c *ArticleController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的文章ID", nil)
		return
	}
	has, _ := models.Get(&models.Article{Id: id})
	if !has {
		c.SendMsg(1, "文章不存在", nil)
		return
	}
	_, err = models.DeleteArticel(id)
	if err != nil {
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *ArticleController) PostRemoves() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要删除的文章", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		article := &models.Article{Id: conv.Int64(id)}
		has, err := models.Get(article)
		if err != nil {
			golog.Errorf("ArticleController PostRemoves Get error : %v", err)
		}
		if !has {
			continue
		}
		_, err = models.DeleteArticel(article.Id)
		if err != nil {
			golog.Errorf("ArticleController PostRemoves DeleteArticel error : %v", err)
			continue
		}
	}

	c.SendMsg(0, "删除成功", nil)
}

func (c *ArticleController) GetSpecialBy(id int64) {
	article, err := models.GetArticleById(id)
	if err != nil {
		golog.Errorf("ArticleController GetSpecialBy GetArticleById error : %v", err)
	}
	special, err := models.GetSpecialByCnd(&models.Special{}, 0, 0, nil)
	if err != nil {
		golog.Errorf("ArticleController GetSpecialBy GetByCnd error : %v", err)
	}
	c.Ctx.ViewData("special", special)
	c.Ctx.ViewData("article", article)
	c.Ctx.View("admin/article/special.html")
}

//绑定专题
func (c *ArticleController) PostSpecial() {
	id, err := c.Ctx.PostValueInt64("Id")
	if err != nil || id == 0 {
		c.SendMsg(1, "非法的文章ID", nil)
		return
	}
	has, _ := models.Get(&models.Article{Id: id})
	if !has {
		c.SendMsg(1, "文章不存在", nil)
		return
	}
	specialId, err := c.Ctx.PostValueInt64("SpecialId")
	if err != nil {
		c.SendMsg(1, "非法的专题ID", nil)
		return
	}
	if specialId == 0 {
		c.SendMsg(1, "请选择专题", nil)
		return
	}
	has, _ = models.Get(&models.Special{Id: specialId})
	if !has {
		c.SendMsg(1, "专题不存在", nil)
		return
	}

	err = models.BindArticleForSpecial(id, specialId)
	if err != nil {
		c.SendMsg(1, "专题添加文章失败", nil)
		return
	}
	c.SendMsg(0, "专题添加文章成功", nil)
}

// 取消绑定
func (c *ArticleController) PostUnbind() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要取消绑定的文章", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		err = models.UnBindArticleForSpecial(conv.Int64(id))
		if err != nil {
			golog.Errorf("ArticleController PostUnbind UnBindArticle error : %v", err)
			continue
		}
	}

	c.SendMsg(0, "取消绑定成功", nil)
}
