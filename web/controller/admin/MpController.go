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

type MpController struct {
	controller.AuthController
}

func (c *MpController) Get() {
	c.Ctx.View("admin/mp/mp.html")
}
func (c *MpController) GetList() {
	params := c.Ctx.URLParams()
	pMp := &models.Mp{}
	status, ok := params["status"]
	if ok {
		pMp.Status = conv.Int(status)
	}
	conv.FillStructStr(params, pMp, false)
	mps, _ := models.GetMpByCnd(pMp, conv.Int(params["limit"]), conv.Int(params["offset"]), nil)
	total := models.Count(pMp)
	c.SendMsgData(iris.Map{"rows": mps, "total": total})
}
//上下架
func (c *MpController) PostGrounding() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的公众号ID", nil)
		return
	}
	mp := &models.Mp{Id: id}
	has, _ := models.Get(mp)
	if !has {
		c.SendMsg(1, "公众号不存在", nil)
		return
	}
	oldStatus := mp.Status
	unGrounding := conv.Int(models.MpStatusUnGrounding)
	grounding := conv.Int(models.MpStatusGrounding)
	var operStr string
	if oldStatus == unGrounding {
		mp.Status = grounding
		operStr = "上架"
	} else {
		mp.Status = unGrounding
		operStr = "下架"
	}

	_, err = models.Update(mp, &models.Mp{Id: id}, nil)
	if err != nil {
		c.SendMsg(1, fmt.Sprintf("%s失败", operStr), nil)
		golog.Errorf("MpController PostGrounding error : %v", err)
		return
	}
	c.SendMsg(0, fmt.Sprintf("%s成功", operStr), nil)
}

func (c *MpController) PostRemove() {
	id, err := c.Ctx.PostValueInt64("id")
	if err != nil {
		c.SendMsg(1, "非法的公众号ID", nil)
		return
	}
	has, _ := models.Get(&models.Mp{Id: id})
	if !has {
		c.SendMsg(1, "公众号不存在", nil)
		return
	}
	_, err = models.Delete(&models.Mp{Id: id},nil)
	if err != nil {
		c.SendMsg(1, "删除失败", nil)
		return
	}
	c.SendMsg(0, "删除成功", nil)
}

func (c *MpController) PostRemoves() {
	ids := c.Ctx.PostValueTrim("ids")
	if len(ids) <= 0 {
		c.SendMsg(1, "请选择需要删除的公众号", nil)
		return
	}
	json, err := json.DecodeToJson(conv.Bytes(ids))
	if err != nil {
		c.SendMsg(1, "参数格式不正确", nil)
		return
	}
	for _, id := range json.ToArray() {
		mp := &models.Mp{Id: conv.Int64(id)}
		has, err := models.Get(mp)
		if err != nil {
			golog.Errorf("MpController PostRemoves Get error : %v", err)
		}
		if !has {
			continue
		}
		_, err = models.Delete(&models.Mp{Id: conv.Int64(id)},nil)
		if err != nil {
			golog.Errorf("MpController PostRemoves Delete error : %v", err)
			continue
		}
	}

	c.SendMsg(0, "删除成功", nil)
}