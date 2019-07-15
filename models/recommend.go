package models

import (
	"fmt"
	"github.com/kataras/golog"
	"time"
	"wxapp/utils/conv"
)

//推荐类型
type RecommendType int

const (
	RecommendTypeHomePage         RecommendType = iota + 1 //首页推荐
	RecommendTypeHomePageCategory                          //首页分类推荐
)

//状态
type RecommendStatus int

const (
	RecommendStatusNo  RecommendStatus = iota + 1 //未上架
	RecommendStatusOn                             //上架
	RecommendStatusOff                            //下架
)

//是否扣费
type RecommendConst int

const (
	RecommendConstOff RecommendConst = iota + 1 //没扣费
	RecommendConstOn                            //已扣费
)

//推荐表
type Recommend struct {
	Id     int64     `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Aid    int64     `xorm:"not null BIGINT(20)"`             //文章ID
	Cid    int64     `xorm:"BIGINT(20)"`                      //分类ID
	Type   int       `xorm:"default 1 INT(11)"`               //类型（1首页推荐，2首页分类推荐，3分类推荐，4详情页推荐）
	Status int       `xorm:"default 1 INT(11)"`               //状态（1未上架，2上架，3下架）
	Rtime  int       `xorm:"default 1 INT(11)"`               //推广时间（以一个月为单位）
	Etime  time.Time `xorm:"DATETIME"`                        //推广结束时间（上架的时候设置时间）
	IsCost int       `xorm:"default 1 INT(11)"`               //是否付费推广（1没付费，2付费）
}

//推荐文章关联表
type RecommendArticle struct {
	Recommend `xorm:"extends"`
	Title     string
}

func (RecommendArticle) TableName() string {
	return "recommend"
}

// 获取指定文章
func GetRecommendById(id int64) (*RecommendArticle, error) {
	datas := &RecommendArticle{}
	sess := x.NewSession()
	defer sess.Close()
	sess.Join("INNER", "article", "article.id = recommend.aid")
	has, err := sess.Id(id).Get(datas)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Recommend does not exist by id = %v", id)
	}
	return datas, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Category里面字段"="条件的过滤
//使用OrderType进行排序
func GetRecommendByCnd(data *RecommendArticle, pageSize int, offset int, orderParams map[string]OrderType) (datas []*RecommendArticle, err error) {
	datas = []*RecommendArticle{}
	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if pageSize > 0 && offset >= 0 {
		sess.Limit(pageSize, offset)
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
	sess.Join("INNER", "article", "article.id = recommend.aid")
	err = sess.Find(&datas, data)
	return datas, err
}

//获取推荐
//@param ptype 	推荐的类型
//@param cid	分类
func GetRecommendArticle(ptype RecommendType, cid int64) []*Article {
	articleList := getCacheRecommendArticle(ptype, cid)
	if articleList != nil {
		return articleList
	}

	sess := x.NewSession()
	defer sess.Close()
	sess.Table("recommend").Cols("Aid")
	sess.Where("type=? and status=?", conv.Int(ptype), conv.Int(RecommendStatusOn))
	if cid != 0 {
		sess.Where("cid=?", cid)
	}

	ids := []int64{}
	err := sess.Find(&ids)
	if err != nil {
		golog.Errorf("GetRecommendArticle param Type = %d, Cid = %d error : %v", conv.Int(ptype), cid, err)
	}

	articleList = []*Article{}
	x.In("id", ids).Find(&articleList)

	return cacheRecommendArticle(ptype, cid, articleList)
}

//获取首页分类，以及分类推荐
func GetRecommendHomeCategoryArticle() []*Category {
	cList := getCategoryArticleRecommend()
	if cList == nil {
		var err error
		cList, err = GetCategoryByCnd(&Category{Status:1}, 0, 0, map[string]OrderType{"sort": ORDER_TYPE_ASE})
		if err != nil {
			golog.Errorf("GetRecommendHomeCategory GetCategoryByCnd error : %v", err)
		}
		cList = cacheCategoryArticleRecommend(cList)
	}

	for _, c := range cList {
		c.Article = GetRecommendArticle(RecommendTypeHomePageCategory, c.Id)
	}
	return cList
}
