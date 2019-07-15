package models

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-xorm/builder"
	"github.com/kataras/golog"
	"time"
	"wxapp/utils/conv"
)

type MpStatus int

const (
	MpStatusGrounding MpStatus = iota + 1
	MpStatusUnGrounding
)
const (
	MpHomeNewSize  int = 8  //首页显示最新公众号的个数
	MpCategorySize int = 48 //分类页显示公众号个数
	MpRankTopSize  int = 50 //排行榜
)

//公众号
type Mp struct {
	Id           int64     `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Cid          int64     `xorm:"BIGINT(20)"`                      //分类ID
	Uid          int64     `xorm:"BIGINT(20)"`                      //用户ID
	Title        string    `xorm:"VARCHAR(100)"`                    //公众号的名字
	MsgLink      string    `xorm:"VARCHAR(255)"`                    //抓取公众号基本信息的链接
	UserName     string    `xorm:"VARCHAR(100) unique"`             //公众号的永久ID
	ShowId       string    `xorm:"VARCHAR(100)"`                    //公众号的显示ID
	Qq           string    `xorm:"VARCHAR(20)"`                     //QQ
	Content      string    `xorm:"VARCHAR(500)"`                    //公众号内容
	PicCover     string    `xorm:"VARCHAR(255)"`                    //封面图片
	PicQrcode    string    `xorm:"VARCHAR(255)"`                    //二维码图片
	Keywords     string    `xorm:"VARCHAR(100)"`                    //关键字
	Browse       int       `xorm:"INT(11)"`                         //浏览次数
	Status       int       `xorm:"INT(1) default 1"`                //状态 1上架，2下架
	CreatedAt    time.Time `xorm:"created"`                         //发布时间
	UpdatedAt    time.Time `xorm:"updated"`                         //更新时间
	StickEndTime time.Time `xorm:"DATETIME"`                        //置顶时间
	SearchWords  string    `xorm:"VARCHAR(255)"`                    //搜索关键字
}

// 获取指定公众号
func GetMpById(id int64) (*Mp, error) {
	data := &Mp{}
	has, err := x.Id(id).Get(data)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Mp does not exist by id = %v", id)
	}
	return data, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Mp里面字段"="条件的过滤
//使用OrderType进行排序
func GetMpByCnd(data *Mp, pageSize int, offset int, orderParams map[string]OrderType) (datas []*Mp, err error) {
	datas = []*Mp{}
	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if pageSize > 0 && offset >= 0 {
		sess.Limit(pageSize, offset)
	}
	//如果分类id指定就查询指定分类下面的所有公众号包括子分类
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

//获取最新公众号列表
func GetNewMpList() []*Mp {
	data := getCacheMpNew()
	if data != nil {
		return data
	}
	data, err := GetMpByCnd(&Mp{}, MpHomeNewSize, 0, map[string]OrderType{"id": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetNewMpList error : %v", err)
	}
	return cacheMpNew(data)
}

//获取指定公众号的关联热门公众号
func GetHotRecommedMpList(mid int64) []*Mp {
	data := getCacheHotRecommedMpList(mid)
	if data != nil {
		return data
	}
	//随机获取缓存中的6条数据
	top_ids := getMpHot()
	indexes := hashset.New()
	genRandArr(6, len(top_ids), indexes)
	var ids []int64
	for _, index := range indexes.Values() {
		ids = append(ids, top_ids[conv.Int(index)])
	}
	err := x.In("id", ids).Table("mp").Find(&data)
	if err != nil {
		golog.Errorf("GetHotRecommedMpList error : %v", err)
	}
	return cacheHotRecommedMpList(mid, data)
}

//获取指定公众号的关联最新公众号
func GetNewRecommedMpList(mid int64) []*Mp {
	data := getCacheNewRecommedMpList(mid)
	if data != nil {
		return data
	}
	//随机获取缓存中的6条数据
	top_ids := getMpNew()
	indexes := hashset.New()
	genRandArr(6, len(top_ids), indexes)
	var ids []int64
	for _, index := range indexes.Values() {
		ids = append(ids, top_ids[conv.Int(index)])
	}
	err := x.In("id", ids).Table("mp").Find(&data)
	if err != nil {
		golog.Errorf("GetNewRecommedMpList error : %v", err)
	}
	return cacheNewRecommedMpList(mid, data)
}

//获取指定公众号的关联别人看过的公众号
func GetSeeRecommedMpList(mid int64, cid int64) []*Mp {
	data := getCacheSeeRecommedMpList(mid)
	if data != nil {
		return data
	}
	data = getCategoryMpStickTop(cid)
	return cacheSeeRecommedMpList(mid, data)
}

//获取最近更新的前500条数据的ID
func getMpNew() []int64 {
	data := getCacheCategoryMpNew()
	if data != nil {
		return data
	}

	var ids []int64
	sess := x.NewSession()
	defer sess.Close()
	sess.Desc("updated_at")
	sess.Limit(500, 0)
	err := sess.Table("mp").Cols("id").Find(&ids)
	if err != nil {
		golog.Errorf("getMpNew error : %v", err)
	}

	return cacheCategoryMpNew(ids)
}

//获取最热的前500条数据的ID
func getMpHot() []int64 {
	data := getCacheCategoryMpHot()
	if data != nil {
		return data
	}

	var ids []int64
	sess := x.NewSession()
	defer sess.Close()
	sess.Desc("browse")
	sess.Limit(500, 0)
	err := sess.Table("mp").Cols("id").Find(&ids)
	if err != nil {
		golog.Errorf("getCategoryMpHot error : %v", err)
	}

	return cacheCategoryMpHot(ids)
}

//获取指定分类置顶的公众号
func getCategoryMpStickTop(cid int64) []*Mp {
	data := getCacheCategoryMpStickTop(cid)
	if data != nil {
		return data
	}
	data, err := GetMpByCnd(&Mp{Cid: cid}, MpCategorySize, 0, map[string]OrderType{"stick_end_time": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetCategoryMpStickTop error : %v", err)
	}
	return cacheCategoryMpStickTop(cid, data)
}

//获取分类下的所有公众号
func GetMpList(page int, sortType int) ([]*Mp, map[string]interface{}) {
	var err error
	var data []*Mp
	var cid int64 = 0
	//统计总数
	count := getCacheMpCountByCid(cid)
	if count == 0 {
		count, err = x.Count(&Mp{})
		if err != nil {
			count = 0
		}
		count = cacheMpCountByCid(cid, count)
	}
	data = getCacheMpListByCid(cid, page, sortType)
	if data != nil {
		return data, Paginator(page, MpCategorySize, count)
	}
	order := make(map[string]OrderType)
	if sortType == 1 {
		order["id"] = ORDER_TYPE_DESC
	} else if sortType == 2 {
		order["browse"] = ORDER_TYPE_DESC
	} else {
		order["stick_end_time"] = ORDER_TYPE_DESC
	}
	data, err = GetMpByCnd(&Mp{}, MpCategorySize, (page-1)*MpCategorySize, order)
	if err != nil {
		golog.Errorf("GetMpList error : %v", err)
	}

	return cacheMpListByCid(cid, page, sortType, data), Paginator(page, MpCategorySize, count)
}

//获取分类下的所有公众号
func GetMpListByCategory(cid int64, page int, sortType int) ([]*Mp, map[string]interface{}) {
	var err error
	var data []*Mp
	cids := getCategoryChildIds(cid)
	//统计总数
	count := getCacheMpCountByCid(cid)
	if count == 0 {
		count, err = x.Where(builder.In("cid", cids)).Count(&Mp{})
		if err != nil {
			count = 0
		}
		count = cacheMpCountByCid(cid, count)
	}
	data = getCacheMpListByCid(cid, page, sortType)
	if data != nil {
		return data, Paginator(page, MpCategorySize, count)
	}
	order := make(map[string]OrderType)
	if sortType == 1 {
		order["id"] = ORDER_TYPE_DESC
	} else if sortType == 2 {
		order["browse"] = ORDER_TYPE_DESC
	} else {
		order["stick_end_time"] = ORDER_TYPE_DESC
	}
	data, err = GetMpByCnd(&Mp{Cid: cid}, MpCategorySize, (page-1)*MpCategorySize, order)
	if err != nil {
		golog.Errorf("GetMpListByCategory error : %v", err)
	}

	return cacheMpListByCid(cid, page, sortType, data), Paginator(page, MpCategorySize, count)
}

func SearchMp(q string, page int) ([]*Mp, map[string]interface{}) {
	result := []*Mp{}
	count, err := x.Where(builder.Like{"search_words", q}).Count(&Mp{})
	if err != nil {
		count = 0
	}
	//如果没有搜索到直接返回
	if count == 0 {
		return result, Paginator(page, MpCategorySize, count)
	}

	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if MpCategorySize > 0 && (page-1)*MpCategorySize >= 0 {
		sess.Limit(MpCategorySize, (page-1)*MpCategorySize)
	}
	//进行排序
	sess.Desc("browse")
	err = sess.Where(builder.Like{"search_words", q}).Find(&result)
	if err != nil {
		golog.Errorf("SearchMp error : %v", err)
	}

	return result, Paginator(page, MpCategorySize, count)
}

//排行榜
func GetMpRank() []*Mp {
	data := getCacheMpRank()
	if data != nil {
		return data
	}
	data, err := GetMpByCnd(&Mp{}, MpRankTopSize, 0, map[string]OrderType{"browse": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetMpRank error : %v", err)
	}
	return cacheMpRank(data)
}
