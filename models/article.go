package models

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-xorm/builder"
	"github.com/kataras/golog"
	"time"
	"wxapp/utils/conv"
)

type ArticleStatus int

const (
	ArticleStatusGrounding ArticleStatus = iota + 1
	ArticleStatusUnGrounding
)

const (
	ArticleHomeNewSize  int = 8  //首页显示最新文章的个数
	CategoryArticleSize int = 48 //分类页显示文章个数
	RankTopSize         int = 50 //排行榜
)

//文章表
type Article struct {
	Id           int64     `xorm:"not null pk BIGINT(20)"` //主键ID
	Cid          int64     `xorm:"BIGINT(20)"`             //分类ID
	Uid          int64     `xorm:"BIGINT(20)"`             //用户ID
	ProvideId    int64     `xorm:"BIGINT(20)"`             //省
	CityId       int64     `xorm:"BIGINT(20)"`             //市
	Title        string    `xorm:"VARCHAR(100) unique"`    //标题
	Author       string    `xorm:"VARCHAR(100)"`            //主体信息
	Qq           string    `xorm:"VARCHAR(20)"`            //QQ
	Content      string    `xorm:"VARCHAR(500)"`           //内容
	PicCover     string    `xorm:"VARCHAR(255)"`           //封面图片
	PicQrcode    string    `xorm:"VARCHAR(255)"`           //二维码图片
	Keywords     string    `xorm:"VARCHAR(100)"`           //关键字
	Browse       int       `xorm:"INT(11)"`                //浏览次数
	Status       int       `xorm:"INT(1) default 1"`       //状态 1上架，2下架
	CreateTime   time.Time `xorm:"DATETIME"`               //发布时间
	UpdateTime   time.Time `xorm:"DATETIME"`               //更新时间
	StickEndTime time.Time `xorm:"DATETIME"`               //置顶时间
	Screenshot   []string  `xorm:"json TEXT"`              //截图
	SearchWords  string    `xorm:"VARCHAR(255)"`           //搜索关键字
}

// 获取指定文章
func GetArticleById(id int64) (*Article, error) {
	data := &Article{}
	has, err := x.Id(id).Get(data)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Article does not exist by id = %v", id)
	}
	return data, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Article里面字段"="条件的过滤
//使用OrderType进行排序
func GetArticleByCnd(data *Article, pageSize int, offset int, orderParams map[string]OrderType) (datas []*Article, err error) {
	datas = []*Article{}
	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if pageSize > 0 && offset >= 0 {
		sess.Limit(pageSize, offset)
	}
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
	return datas, err
}

//获取最新文章列表
func GetNewArticleList() []*Article {
	data := getCacheArticleNew()
	if data != nil {
		return data
	}
	data, err := GetArticleByCnd(&Article{}, ArticleHomeNewSize, 0, map[string]OrderType{"id": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetNewArticleList error : %v", err)
	}
	return cacheArticleNew(data)
}

//获取指定文章的关联热门文章
func GetHotRecommedArticleList(aid int64) []*Article {
	data := getCacheHotRecommedArticleList(aid)
	if data != nil {
		return data
	}
	//随机获取缓存中的6条数据
	top_ids := getArticleHot()
	indexes := hashset.New()
	genRandArr(6, len(top_ids), indexes)
	var ids []int64
	for _, index := range indexes.Values() {
		ids = append(ids, top_ids[conv.Int(index)])
	}
	err := x.In("id", ids).Table("article").Find(&data)
	if err != nil {
		golog.Errorf("GetHotRecommedArticleList error : %v", err)
	}
	return cacheHotRecommedArticleList(aid, data)
}

//获取指定文章的关联最新文章
func GetNewRecommedArticleList(aid int64) []*Article {
	data := getCacheNewRecommedArticleList(aid)
	if data != nil {
		return data
	}
	//随机获取缓存中的6条数据
	top_ids := getArticleNew()
	indexes := hashset.New()
	genRandArr(6, len(top_ids), indexes)
	var ids []int64
	for _, index := range indexes.Values() {
		ids = append(ids, top_ids[conv.Int(index)])
	}
	err := x.In("id", ids).Table("article").Find(&data)
	if err != nil {
		golog.Errorf("GetNewRecommedArticleList error : %v", err)
	}
	return cacheNewRecommedArticleList(aid, data)
}

//获取指定文章的关联别人看过的文章
func GetSeeRecommedArticleList(aid int64, cid int64) []*Article {
	data := getCacheSeeRecommedArticleList(aid)
	if data != nil {
		return data
	}
	data = getCategoryArticleStickTop(cid)
	return cacheSeeRecommedArticleList(aid, data)
}

//获取最近更新的前500条数据的ID
func getArticleNew() []int64 {
	data := getCacheCategoryArticleNew()
	if data != nil {
		return data
	}

	var ids []int64
	sess := x.NewSession()
	defer sess.Close()
	sess.Desc("update_time")
	sess.Limit(500, 0)
	err := sess.Table("article").Cols("id").Find(&ids)
	if err != nil {
		golog.Errorf("getCategoryArticleNew error : %v", err)
	}

	return cacheCategoryArticleNew(ids)
}

//获取最热的前500条数据的ID
func getArticleHot() []int64 {
	data := getCacheCategoryArticleHot()
	if data != nil {
		return data
	}

	var ids []int64
	sess := x.NewSession()
	defer sess.Close()
	sess.Desc("browse")
	sess.Limit(500, 0)
	err := sess.Table("article").Cols("id").Find(&ids)
	if err != nil {
		golog.Errorf("getCategoryArticleHot error : %v", err)
	}

	return cacheCategoryArticleHot(ids)
}

//获取指定分类置顶的文章
func getCategoryArticleStickTop(cid int64) []*Article {
	data := getCacheCategoryArticleStickTop(cid)
	if data != nil {
		return data
	}
	data, err := GetArticleByCnd(&Article{Cid: cid}, CategoryArticleSize, 0, map[string]OrderType{"stick_end_time": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetCategoryArticleStickTop error : %v", err)
	}
	return cacheCategoryArticleStickTop(cid, data)
}

//获取分类下的所有文章
func GetArticleList(page int, sortType int) ([]*Article, map[string]interface{}) {
	var err error
	var data []*Article
	var cid int64 = 0
	//统计总数
	count := getCacheArticleCountByCid(cid)
	if count == 0 {
		count, err = x.Count(&Article{})
		if err != nil {
			count = 0
		}
		count = cacheArticleCountByCid(cid, count)
	}
	data = getCacheArticleListByCid(cid, page, sortType)
	if data != nil {
		return data, Paginator(page, CategoryArticleSize, count)
	}
	order := make(map[string]OrderType)
	if sortType == 1 {
		order["id"] = ORDER_TYPE_DESC
	} else if sortType == 2 {
		order["browse"] = ORDER_TYPE_DESC
	} else {
		order["stick_end_time"] = ORDER_TYPE_DESC
	}
	data, err = GetArticleByCnd(&Article{}, CategoryArticleSize, (page-1)*CategoryArticleSize, order)
	if err != nil {
		golog.Errorf("GetArticleList error : %v", err)
	}

	return cacheArticleListByCid(cid, page, sortType, data), Paginator(page, CategoryArticleSize, count)
}

//获取分类下的所有文章
func GetArticleListByCategory(cid int64, page int, sortType int) ([]*Article, map[string]interface{}) {
	var err error
	var data []*Article
	cids := getCategoryChildIds(cid)
	//统计总数
	count := getCacheArticleCountByCid(cid)
	if count == 0 {
		count, err = x.Where(builder.In("cid", cids)).Count(&Article{})
		if err != nil {
			count = 0
		}
		count = cacheArticleCountByCid(cid, count)
	}
	data = getCacheArticleListByCid(cid, page, sortType)
	if data != nil {
		return data, Paginator(page, CategoryArticleSize, count)
	}
	order := make(map[string]OrderType)
	if sortType == 1 {
		order["id"] = ORDER_TYPE_DESC
	} else if sortType == 2 {
		order["browse"] = ORDER_TYPE_DESC
	} else {
		order["stick_end_time"] = ORDER_TYPE_DESC
	}
	data, err = GetArticleByCnd(&Article{Cid: cid}, CategoryArticleSize, (page-1)*CategoryArticleSize, order)
	if err != nil {
		golog.Errorf("GetArticleListByCategory error : %v", err)
	}

	return cacheArticleListByCid(cid, page, sortType, data), Paginator(page, CategoryArticleSize, count)
}

func SearchArticle(q string, page int) ([]*Article, map[string]interface{}) {
	result := []*Article{}
	count, err := x.Where(builder.Like{"search_words", q}).Count(&Article{})
	if err != nil {
		count = 0
	}
	//如果没有搜索到直接返回
	if count == 0 {
		return result, Paginator(page, CategoryArticleSize, count)
	}

	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if CategoryArticleSize > 0 && (page-1)*CategoryArticleSize >= 0 {
		sess.Limit(CategoryArticleSize, (page-1)*CategoryArticleSize)
	}
	//进行排序
	sess.Desc("browse")
	err = sess.Where(builder.Like{"search_words", q}).Find(&result)
	if err != nil {
		golog.Errorf("SearchArticle error : %v", err)
	}

	return result, Paginator(page, CategoryArticleSize, count)
}

//排行榜
func GetArticleRank() []*Article {
	data := getCacheArticleRank()
	if data != nil {
		return data
	}
	data, err := GetArticleByCnd(&Article{}, RankTopSize, 0, map[string]OrderType{"browse": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetArticleRank error : %v", err)
	}
	return cacheArticleRank(data)
}

func DeleteArticel(articelId int64) (bool, error) {
	sess := NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return false, err
	}
	_, err := sess.Delete(&Article{Id: articelId})
	if err != nil {
		golog.Errorf("DeleteArticel Delete error : %v", err)
		return false, sess.Rollback()
	}
	_, err = sess.Delete(&SpecialArtice{Aid: articelId})
	if err != nil {
		golog.Errorf("DeleteArticel Delete SpecialArtice error : %v", err)
		return false, sess.Rollback()
	}
	return true, sess.Commit()
}
