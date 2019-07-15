package models

import (
	"fmt"
	"github.com/go-xorm/builder"
	"time"
	"wxapp/utils/conv"
)

type PassType int

const (
	PassTypeCommon PassType = iota + 1
	PassTypeUnPass
	PassTypePass
)

//文章表
type Apply struct {
	Id         int64     `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Cid        int64     `xorm:"BIGINT(20)"`                      //分类ID
	Uid        int64     `xorm:"BIGINT(20)"`                      //用户ID
	ProvideId  int64     `xorm:"BIGINT(20)"`                      //省
	CityId     int64     `xorm:"BIGINT(20)"`                      //市
	Title      string    `xorm:"VARCHAR(100) unique"`             //标题
	Author     string    `xorm:"VARCHAR(100)"`                     //主体信息
	Qq         string    `xorm:"VARCHAR(20)"`                     //QQ
	Content    string    `xorm:"VARCHAR(500)"`                    //内容
	PicCover   string    `xorm:"VARCHAR(255)"`                    //封面图片
	PicQrcode  string    `xorm:"VARCHAR(255)"`                    //二维码图片
	Browse     int       `xorm:"INT(11)"`                         //浏览次数
	Pass       int       `xorm:"INT(1) default 1"`                //审核是否通过(1待审核，2审核不通过，3审核通过)
	CreateTime time.Time `xorm:"DATETIME"`                        //发布时间
	Screenshot []string  `xorm:"json TEXT"`                       //截图
	Keywords   string    `xorm:"VARCHAR(100)"`                    //标签
}

// 获取指定文章
func GetApplyById(id int64) (*Apply, error) {
	data := &Apply{}
	has, err := x.Id(id).Get(data)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Apply does not exist by id = %v", id)
	}
	return data, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Article里面字段"="条件的过滤
//使用OrderType进行排序
func GetApplyByCnd(data *Apply, pageSize int, offset int, orderParams map[string]OrderType) (datas []*Apply, err error) {
	datas = []*Apply{}
	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if pageSize > 0 && offset >= 0 {
		sess.Limit(pageSize, offset)
	}
	cid := data.Cid
	//如果分类id指定就查询指定分类下面的所有文章包括子分类
	if data.Cid > 0 {
		cids := getCategoryChildIds(data.Cid)
		sess.Where(builder.In("cid", cids))
		data.Cid = 0
	}
	//进行排序
	if orderParams != nil {
		for feild, orderType := range orderParams {
			if orderType == ORDER_TYPE_ASE {
				sess.Asc(feild)
			} else if orderType == ORDER_TYPE_DESC {
				sess.Desc(feild)
			}
		}
	}
	err = sess.Find(&datas, data)
	//执行查询后还原数据
	data.Cid = cid
	return datas, err
}

func GetApplyByCndCount(data *Apply) int64 {
	//如果分类id指定就查询指定分类下面的所有文章包括子分类
	sess := NewSession()
	if data.Cid > 0 {
		cids := getCategoryChildIds(data.Cid)
		sess.Where(builder.In("cid", cids))
		data.Cid = 0
	}
	count, _ := sess.Count(data)
	return count
}

func ApplyPass(apply *Apply) (bool, error) {
	if apply.Pass == conv.Int(PassTypePass) {
		articel, _ := GetArticleById(apply.Id)
		if articel == nil {
			return insertApplyArticel(apply)
		} else {
			return updateApplyArticel(apply)
		}
	} else {
		_, err := Update(apply, Apply{Id: apply.Id}, nil)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}
func insertApplyArticel(apply *Apply) (bool, error) {
	sess := NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return false, err
	}
	_, err := Update(apply, Apply{Id: apply.Id}, sess)
	if err != nil {
		return false, sess.Rollback()
	}
	articel := apply2Artice(apply)
	articel.Status = conv.Int(ArticleStatusGrounding)
	_, err = Insert(articel, sess)
	if err != nil {
		return false, sess.Rollback()
	}
	return true, sess.Commit()
}

func updateApplyArticel(apply *Apply) (bool, error) {
	sess := NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return false, err
	}
	_, err := Update(apply, Apply{Id: apply.Id}, sess)
	if err != nil {
		return false, sess.Rollback()
	}
	articel := apply2Artice(apply)
	articel.Status = conv.Int(ArticleStatusGrounding)
	_, err = Update(articel, &Article{Id: articel.Id}, sess)
	if err != nil {
		return false, sess.Rollback()
	}
	return true, sess.Commit()
}

func apply2Artice(apply *Apply) *Article {
	article := &Article{}
	article.Id = apply.Id
	article.Author = apply.Author
	article.Screenshot = apply.Screenshot
	article.ProvideId = apply.ProvideId
	article.CityId = apply.CityId
	article.Cid = apply.Cid
	article.Content = apply.Content
	article.PicCover = apply.PicCover
	article.PicQrcode = apply.PicQrcode
	article.Qq = apply.Qq
	article.Title = apply.Title
	article.Uid = apply.Uid
	article.Keywords = apply.Keywords
	article.CreateTime = apply.CreateTime
	// 合成搜索关键字
	article.SearchWords = fmt.Sprintf("%s %s", apply.Title, apply.Keywords)
	return article
}
