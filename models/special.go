package models

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/golog"
	"wxapp/utils/conv"
)

const (
	SPECIAL_PAGE_SIZE int = 8 //专题页长度
)

//专题表
type Special struct {
	Id      int64  `xorm:"not null pk autoincr BIGINT(20)"` //专题ID
	Title   string `xorm:"VARCHAR(50)"`                     //专题名
	Desc    string `xorm:"VARCHAR(255)"`                    //专题描述
	Pic     string `xorm:"VARCHAR(255)"`                    //专题图片
	Deleted int    `xorm:"deleted INT(11)"`                 //软删除标识

	Articles []*ArticleSpecial `xorm:"-"`
}

// 获取指定文章
func GetSpecialById(id int64) (*Special, error) {
	datas := &Special{}
	has, err := x.Id(id).Get(datas)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Special does not exist by id = %v", id)
	}
	return datas, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Special里面字段"="条件的过滤
//使用OrderType进行排序
func GetSpecialByCnd(data *Special, pageSize int, offset int, orderParams map[string]OrderType) (datas []*Special, err error) {
	datas = []*Special{}
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

//绑定文章到专题
func BindArticleForSpecial(aid int64, sid int64) error {
	sa := &SpecialArtice{Aid: aid, Sid: sid}
	has, err := Get(sa)
	if err != nil {
		return err
	}
	if has {
		return nil
	}
	_, err = Insert(sa, nil)
	return err
}

//解除绑定
func UnBindArticleForSpecial(id int64) error {
	sa := &SpecialArtice{Id: id}
	has, err := Get(sa)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}
	_, err = Delete(sa, nil)
	return err
}

//获取指定专题的文章列表
func GetArticleBySpecialId(sid int64, limit int) (datas []*ArticleSpecial, err error) {
	datas = []*ArticleSpecial{}
	sess := x.Alias("article")
	if limit > 0 {
		sess.Limit(limit, 0)
	}
	sess.Join("LEFT", "special_artice", "article.id = special_artice.aid")
	sess.Where("special_artice.sid = ?", sid)
	err = sess.Find(&datas)

	return datas, err
}

//获取专题和专题相关的文章
func GetSpecialArticleById(sid int64) *Special {
	special := getCacheSpecialAndArticle(sid)
	if special != nil {
		return special
	}
	special, err := GetSpecialById(sid)
	if err != nil {
		golog.Errorf("GetSpecialArticleById GetSpecialById error : %v", err)
	}
	articles, err := GetArticleBySpecialId(sid, 0)
	if err != nil {
		golog.Errorf("GetSpecialArticleById GetArticleBySpecialId error : %v", err)
	}
	special.Articles = articles
	return cacheSpecialAndArticle(sid, special)
}

//获取所有专题以及文章
func GetSpecialList(page int) ([]*Special, map[string]interface{}) {
	var err error
	//统计总数
	count := getCacheSpecialCount()
	if count == 0 {
		count, err = x.Count(&Special{})
		if err != nil {
			count = 0
		}
		count = cacheSpecialCount(count)
	}
	specials := getCacheSpecialsByPage(page)
	if specials != nil {
		return specials, Paginator(page, SPECIAL_PAGE_SIZE, count)
	}

	specials, err = GetSpecialByCnd(&Special{}, SPECIAL_PAGE_SIZE, (page-1)*SPECIAL_PAGE_SIZE, nil)
	if err != nil {
		golog.Errorf("GetSpecialList error : %v", err)
	}
	//获取专题下的文章
	for _, sp := range specials {
		aList, err := GetArticleBySpecialId(sp.Id, 3)
		if err != nil {
			golog.Errorf("GetSpecialList GetArticleBySpecialId error : %v", err)
		}
		sp.Articles = aList
	}
	return cacheSpecialsByPage(page, specials), Paginator(page, SPECIAL_PAGE_SIZE, count)
}

//获取随机专题
func GetRandSpecialList(cacheType string) []*Special {
	randSpecials := getCacheRandSpecial(cacheType)
	if randSpecials != nil {
		return randSpecials
	}
	specilAll := getAllSpecialList()
	indexes := hashset.New()
	genRandArr(4, len(specilAll), indexes)
	for _, index := range indexes.Values() {
		randSpecials = append(randSpecials, specilAll[conv.Int(index)])
	}
	return cacheRandSpecial(cacheType, randSpecials)
}

func getAllSpecialList() []*Special {
	specials := getCacheSpecials()
	if specials != nil {
		return specials
	}

	specials, err := GetSpecialByCnd(&Special{}, 0, 0, nil)
	if err != nil {
		golog.Errorf("getAllSpecialList error : %v", err)
	}

	return cacheSpecials(specials)
}
