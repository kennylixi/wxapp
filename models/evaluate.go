package models

import (
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/kataras/golog"
	"time"
)

const (
	EvaluatePageSize int = 8 //页长
	EvaluateHotSize  int = 6 //热门长度
)

//评测
type Evaluate struct {
	Id          int64     `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Title       string    `xorm:"VARCHAR(250)"`                    //类目名
	Aid         int64     `xorm:"not null BIGINT(20) unique"`      //文章id
	Source      string    `xorm:"VARCHAR(100)"`                    //文章来源
	Score       int       `xorm:"INT(11)"`                         //平分
	Content     string    `xorm:"TEXT"`                            //内容
	Desc        string    `xorm:"VARCHAR(500)"`                    //简单描述
	Browse      int       `xorm:"INT(11) default 0"`               //浏览次数
	Keywords    string    `xorm:"VARCHAR(100)"`                    //标签
	PicCover    string    `xorm:"VARCHAR(255)`                     //图片
	CreatedAt   time.Time `xorm:"created"`                         //发布时间
	SearchWords string    `xorm:"VARCHAR(400)"`                    //搜索关键字

	Article *Article `xorm:"-"`
}

// 获取指定评测
func GetEvaluateById(id int64) (*Evaluate, error) {
	datas := &Evaluate{}
	has, err := x.Id(id).Get(datas)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Evaluate does not exist by id = %v", id)
	}
	return datas, nil
}

// 根据文章id获取评测
func GetEvaluateByAid(aid int64) (*Evaluate, error) {
	eva := &Evaluate{Aid: aid}
	has, err := x.Get(eva)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Evaluate does not exist by aid = %v", aid)
	}
	return eva, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Link里面字段"="条件的过滤
//使用OrderType进行排序
func GetEvaluateByCnd(data *Evaluate, pageSize int, offset int, orderParams map[string]OrderType) (datas []*Evaluate, err error) {
	datas = []*Evaluate{}
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
	err = sess.Find(&datas, data)
	return datas, err
}

//获取分类下的所有文章
func GetEvaluateList(page int) ([]*Evaluate, map[string]interface{}) {
	var err error
	var data []*Evaluate
	//统计总数
	count := getCacheEvaluateListCount()
	if count == 0 {
		count, err = x.Count(&Evaluate{})
		if err != nil {
			count = 0
		}
		count = cacheEvaluateListCount(count)
	}
	data = getCacheEvaluateList(page)
	if data != nil {
		return data, Paginator(page, EvaluatePageSize, count)
	}

	data, err = GetEvaluateByCnd(&Evaluate{}, EvaluatePageSize, (page-1)*EvaluatePageSize, map[string]OrderType{"id": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetEvaluateList error : %v", err)
	}

	return cacheEvaluateList(page, data), Paginator(page, EvaluatePageSize, count)
}

//获取热度前4个
func GetEvaluateHotList() []*Evaluate {
	data := getCacheEvaluateHotList()
	data, err := GetEvaluateByCnd(&Evaluate{}, EvaluateHotSize, 0, map[string]OrderType{"browse": ORDER_TYPE_DESC})
	if err != nil {
		golog.Errorf("GetEvaluateHotList error : %v", err)
	}
	return cacheEvaluateHotList(data)
}

func SearchEvaluate(q string, page int) ([]*Evaluate, map[string]interface{}) {
	result := []*Evaluate{}
	count, err := x.Where(builder.Like{"search_words", q}).Count(&Evaluate{})
	if err != nil {
		count = 0
	}
	//如果没有搜索到直接返回
	if count == 0 {
		return result, Paginator(page, EvaluatePageSize, count)
	}

	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if EvaluatePageSize > 0 && (page-1)*EvaluatePageSize >= 0 {
		sess.Limit(EvaluatePageSize, (page-1)*EvaluatePageSize)
	}
	//进行排序
	sess.Desc("browse")
	err = sess.Where(builder.Like{"search_words", q}).Find(&result)
	if err != nil {
		golog.Errorf("SearchEvaluate error : %v", err)
	}

	return result, Paginator(page, EvaluatePageSize, count)
}
