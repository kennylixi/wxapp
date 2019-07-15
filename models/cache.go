package models

import (
	"fmt"
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/muesli/cache2go"
	"math/rand"
	"time"
	"wxapp/utils/conv"
	rand2 "wxapp/utils/rand"
)

const (
	c_KEY_CATEGORY_NAV               string = "CATEGORY_NAV"               //分类导航
	c_KEY_CATEGORY_NAV_CHILD_IDS     string = "CATEGORY_NAV_CHILD_IDS"     //分类子类的ID
	c_KEY_CATEGORY_ARTICLE_RECOMMEND string = "CATEGORY_ARITCLE_RECOMMEND" //分类文字推荐
	c_KEY_CATEGORY_ARTICLE_STICK_TOP string = "CATEGORY_ARTICLE_STICK_TOP" //缓存指定分类置顶的文章
	c_KEY_ARTICLE_NEW                string = "ARTICLE_NEW"                //最新文章
	c_KEY_ARTICLE_HOT_RECOMMEND      string = "ARTICLE_HOT_RECOMMEND"      //文章的关联热门文章
	c_KEY_ARTICLE_NEW_RECOMMEND      string = "ARTICLE_NEW_RECOMMEND"      //文章的关联最新文章
	c_KEY_ARTICLE_SEE_RECOMMEND      string = "ARTICLE_SEE_RECOMMEND"      //文章的关联别人看过的文章
	c_KEY_ARTICLE_RANK               string = "ARTICLE_RANK"               //排行行
	c_KEY_ARTICLE_NEW_500            string = "ARTICLE_NEW_500"            //缓存最热的前500条数据的ID
	c_KEY_ARTICLE_HOT_500            string = "ARTICLE_HOT_500"            //缓存最热的前500条数据的ID
	c_KEY_RECOMMEND                  string = "RECOMMEND"                  //推荐
	c_KEY_CATEGORY_ARTICLE           string = "CATEGORY_ARTICLE"           //类目文章
	c_KEY_CATEGORY_ARTICLE_COUNT     string = "CATEGORY_ARTICLE_COUNT"     //类目文章总数
	c_KEY_SPECIAL_COUNT              string = "SPECIAL_COUNT"              //专题总数
	c_KEY_SPECIAL_ARTICLE            string = "SPECIAL_ARTICLE"            //专题和专题文章
	c_KEY_SPECIALS                   string = "SPECIALS"                   //专题列表
	c_KEY_RAND_SPECIAL               string = "RAND_SPECIAL"               //随机专题
	c_KEY_EVALUATE_COUNT             string = "EVALUATE_COUNT"             //评测文章总数
	c_KEY_EVALUATE                   string = "EVALUATE"                   //评测文章
	c_KEY_EVALUATE_HOT               string = "EVALUATE_HOT"               //评测热度榜
	c_KEY_MP_NEW                     string = "MP_NEW"                     //最新公众号
	c_KEY_MP_HOT_RECOMMEND           string = "MP_HOT_RECOMMEND"           //公众号的关联热门公众号
	c_KEY_MP_NEW_RECOMMEND           string = "MP_NEW_RECOMMEND"           //公众号的关联最新公众号
	c_KEY_MP_SEE_RECOMMEND           string = "MP_SEE_RECOMMEND"           //公众号的关联别人看过的公众号
	c_KEY_MP_NEW_500                 string = "MP_NEW_500"                 //缓存最热的前500条数据的ID
	c_KEY_MP_HOT_500                 string = "MP_HOT_500"                 //缓存最热的前500条数据的ID
	c_KEY_CATEGORY_MP_STICK_TOP      string = "CATEGORY_MP_STICK_TOP"      //缓存指定分类置顶的公众号
	c_KEY_CATEGORY_MP_COUNT          string = "CATEGORY_MP_COUNT"          //类目文章总数
	c_KEY_CATEGORY_MP                string = "CATEGORY_MP"                //类目文章
	c_KEY_MP_RANK                    string = "MP_RANK"                    //公众号排行榜
)

var (
	cache *cache2go.CacheTable
)

func InitCache() {
	cache = cache2go.Cache("myCache")
}

//清除缓存
func ClearCache() {
	if cache != nil {
		cache.Flush()
	}
}

//分类导航缓存
func cacheCategoryNav(pid int64, data []*Category) []*Category {
	liftTime, _ := time.ParseDuration("1h")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_NAV, pid), liftTime, data).Data().([]*Category)
}
func getCacheCategoryNav(pid int64) []*Category {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_NAV, pid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Category)
}

//分类导航子分类ID缓存
func cacheCategoryChildIds(id int64, data []int64) []int64 {
	liftTime, _ := time.ParseDuration("1h")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_NAV_CHILD_IDS, id), liftTime, data).Data().([]int64)
}
func getCacheCategoryChildIds(id int64) []int64 {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_NAV_CHILD_IDS, id))
	if data == nil {
		return nil
	}
	return data.Data().([]int64)
}

//分类文字推荐缓存
func cacheCategoryArticleRecommend(data []*Category) []*Category {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(c_KEY_CATEGORY_ARTICLE_RECOMMEND, liftTime, data).Data().([]*Category)
}
func getCategoryArticleRecommend() []*Category {
	data, _ := cache.Value(c_KEY_CATEGORY_ARTICLE_RECOMMEND)
	if data == nil {
		return nil
	}
	return data.Data().([]*Category)
}

//分类推荐
func cacheRecommendArticle(ptype RecommendType, cid int64, article []*Article) []*Article {
	liftTime, _ := time.ParseDuration("10m")
	randArticleArray(article)
	return cache.Add(fmt.Sprintf("%s_%d_%d", c_KEY_RECOMMEND, conv.Int(ptype), cid), liftTime, article).Data().([]*Article)
}
func getCacheRecommendArticle(ptype RecommendType, cid int64) []*Article {
	data, _ := cache.Value(fmt.Sprintf("%s_%d_%d", c_KEY_RECOMMEND, conv.Int(ptype), cid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//最新文章
func cacheArticleNew(article []*Article) []*Article {
	liftTime, _ := time.ParseDuration("10m")
	randArticleArray(article)
	return cache.Add(c_KEY_ARTICLE_NEW, liftTime, article).Data().([]*Article)
}
func getCacheArticleNew() []*Article {
	data, _ := cache.Value(c_KEY_ARTICLE_NEW)
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//缓存指定文章的关联热门文章
func cacheHotRecommedArticleList(aid int64, article []*Article) []*Article {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_ARTICLE_HOT_RECOMMEND, aid), liftTime, article).Data().([]*Article)
}
func getCacheHotRecommedArticleList(aid int64) []*Article {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_ARTICLE_HOT_RECOMMEND, aid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//缓存指定文章的关联最新文章
func cacheNewRecommedArticleList(aid int64, article []*Article) []*Article {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_ARTICLE_NEW_RECOMMEND, aid), liftTime, article).Data().([]*Article)
}
func getCacheNewRecommedArticleList(aid int64) []*Article {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_ARTICLE_NEW_RECOMMEND, aid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//缓存指定文章的关联别人看过的文章
func cacheSeeRecommedArticleList(aid int64, article []*Article) []*Article {
	liftTime, _ := time.ParseDuration("10m")
	//随机获取索引
	indexes := hashset.New()
	genRandArr(8, len(article), indexes)
	data := []*Article{}
	for _, index := range indexes.Values() {
		data = append(data, article[conv.Int(index)])
	}
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_ARTICLE_SEE_RECOMMEND, aid), liftTime, data).Data().([]*Article)
}
func getCacheSeeRecommedArticleList(aid int64) []*Article {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_ARTICLE_SEE_RECOMMEND, aid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//缓存最近更新前500条数据的ID
func cacheCategoryArticleNew(ids []int64) []int64 {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(c_KEY_ARTICLE_NEW_500, liftTime, ids).Data().([]int64)
}
func getCacheCategoryArticleNew() []int64 {
	data, _ := cache.Value(c_KEY_ARTICLE_NEW_500)
	if data == nil {
		return nil
	}
	return data.Data().([]int64)
}

//缓存最热的前500条数据的ID
func cacheCategoryArticleHot(ids []int64) []int64 {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(c_KEY_ARTICLE_HOT_500, liftTime, ids).Data().([]int64)
}
func getCacheCategoryArticleHot() []int64 {
	data, _ := cache.Value(c_KEY_ARTICLE_HOT_500)
	if data == nil {
		return nil
	}
	return data.Data().([]int64)
}

//缓存指定分类置顶的文章
func cacheCategoryArticleStickTop(cid int64, article []*Article) []*Article {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_ARTICLE_STICK_TOP, cid), liftTime, article).Data().([]*Article)
}
func getCacheCategoryArticleStickTop(cid int64) []*Article {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_ARTICLE_STICK_TOP, cid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//文章排行榜
func cacheArticleRank(article []*Article) []*Article {
	liftTime, _ := time.ParseDuration("30m")
	return cache.Add(c_KEY_ARTICLE_RANK, liftTime, article).Data().([]*Article)
}
func getCacheArticleRank() []*Article {
	data, _ := cache.Value(c_KEY_ARTICLE_RANK)
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//类目文章缓存
func cacheArticleListByCid(cid int64, page int, sortType int, data []*Article) []*Article {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d_%d_%d", c_KEY_CATEGORY_ARTICLE, cid, page, sortType), liftTime, data).Data().([]*Article)
}
func getCacheArticleListByCid(cid int64, page int, sortType int) []*Article {
	data, _ := cache.Value(fmt.Sprintf("%s_%d_%d_%d", c_KEY_CATEGORY_ARTICLE, cid, page, sortType))
	if data == nil {
		return nil
	}
	return data.Data().([]*Article)
}

//类目文章总数
func cacheArticleCountByCid(cid int64, count int64) int64 {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_ARTICLE_COUNT, cid), liftTime, count).Data().(int64)
}
func getCacheArticleCountByCid(cid int64) int64 {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_ARTICLE_COUNT, cid))
	if data == nil {
		return 0
	}
	return data.Data().(int64)
}

//专题总数
func cacheSpecialCount(count int64) int64 {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(c_KEY_SPECIAL_COUNT, liftTime, count).Data().(int64)
}
func getCacheSpecialCount() int64 {
	data, _ := cache.Value(c_KEY_SPECIAL_COUNT)
	if data == nil {
		return 0
	}
	return data.Data().(int64)
}

//缓存专题
func cacheSpecialAndArticle(sid int64, special *Special) *Special {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_SPECIAL_ARTICLE, sid), liftTime, special).Data().(*Special)
}
func getCacheSpecialAndArticle(sid int64) *Special {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_SPECIAL_ARTICLE, sid))
	if data == nil {
		return nil
	}
	return data.Data().(*Special)
}

//所有专题
func cacheSpecials(data []*Special) []*Special {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(c_KEY_SPECIALS, liftTime, data).Data().([]*Special)
}
func getCacheSpecials() []*Special {
	data, _ := cache.Value(c_KEY_SPECIALS)
	if data == nil {
		return nil
	}
	return data.Data().([]*Special)
}

//专题(包括缓存专题的前三个文章)
func cacheSpecialsByPage(page int, data []*Special) []*Special {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_SPECIALS, page), liftTime, data).Data().([]*Special)
}
func getCacheSpecialsByPage(page int) []*Special {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_SPECIALS, page))
	if data == nil {
		return nil
	}
	return data.Data().([]*Special)
}

//随机状态缓存
func cacheRandSpecial(cacheType string, data []*Special) []*Special {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%s", c_KEY_RAND_SPECIAL, cacheType), liftTime, data).Data().([]*Special)
}
func getCacheRandSpecial(cacheType string) []*Special {
	data, _ := cache.Value(fmt.Sprintf("%s_%s", c_KEY_RAND_SPECIAL, cacheType))
	if data == nil {
		return nil
	}
	return data.Data().([]*Special)
}

//评测总数
func cacheEvaluateListCount(count int64) int64 {
	liftTime, _ := time.ParseDuration("12h")
	return cache.Add(c_KEY_EVALUATE_COUNT, liftTime, count).Data().(int64)
}
func getCacheEvaluateListCount() int64 {
	data, _ := cache.Value(c_KEY_EVALUATE_COUNT)
	if data == nil {
		return 0
	}
	return data.Data().(int64)
}

//评测文章
func cacheEvaluateList(page int, data []*Evaluate) []*Evaluate {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_EVALUATE, page), liftTime, data).Data().([]*Evaluate)
}
func getCacheEvaluateList(page int) []*Evaluate {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_EVALUATE, page))
	if data == nil {
		return nil
	}
	return data.Data().([]*Evaluate)
}

//评测热度榜
func cacheEvaluateHotList(data []*Evaluate) []*Evaluate {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(c_KEY_EVALUATE_HOT, liftTime, data).Data().([]*Evaluate)
}
func getCacheEvaluateHotList() []*Evaluate {
	data, _ := cache.Value(c_KEY_EVALUATE_HOT)
	if data == nil {
		return nil
	}
	return data.Data().([]*Evaluate)
}

//=====================MP Start=====================================
//最新公众号
func cacheMpNew(mp []*Mp) []*Mp {
	liftTime, _ := time.ParseDuration("10m")
	randMpArray(mp)
	return cache.Add(c_KEY_MP_NEW, liftTime, mp).Data().([]*Mp)
}
func getCacheMpNew() []*Mp {
	data, _ := cache.Value(c_KEY_MP_NEW)
	if data == nil {
		return nil
	}
	return data.Data().([]*Mp)
}

//缓存指定公众号的关联热门公众号
func cacheHotRecommedMpList(mid int64, mp []*Mp) []*Mp {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_MP_HOT_RECOMMEND, mid), liftTime, mp).Data().([]*Mp)
}
func getCacheHotRecommedMpList(mid int64) []*Mp {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_MP_HOT_RECOMMEND, mid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Mp)
}

//缓存指定公众号的关联最新公众号
func cacheNewRecommedMpList(mid int64, mp []*Mp) []*Mp {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_MP_NEW_RECOMMEND, mid), liftTime, mp).Data().([]*Mp)
}
func getCacheNewRecommedMpList(mid int64) []*Mp {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_MP_NEW_RECOMMEND, mid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Mp)
}

//缓存指定公众号的关联别人看过的公众号
func cacheSeeRecommedMpList(mid int64, mp []*Mp) []*Mp {
	liftTime, _ := time.ParseDuration("10m")
	//随机获取索引
	indexes := hashset.New()
	genRandArr(8, len(mp), indexes)
	data := []*Mp{}
	for _, index := range indexes.Values() {
		data = append(data, mp[conv.Int(index)])
	}
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_MP_SEE_RECOMMEND, mid), liftTime, data).Data().([]*Mp)
}
func getCacheSeeRecommedMpList(mid int64) []*Mp {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_MP_SEE_RECOMMEND, mid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Mp)
}

//缓存最近更新前500条数据的ID
func cacheCategoryMpNew(ids []int64) []int64 {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(c_KEY_MP_NEW_500, liftTime, ids).Data().([]int64)
}
func getCacheCategoryMpNew() []int64 {
	data, _ := cache.Value(c_KEY_MP_NEW_500)
	if data == nil {
		return nil
	}
	return data.Data().([]int64)
}

//缓存最热的前500条数据的ID
func cacheCategoryMpHot(ids []int64) []int64 {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(c_KEY_MP_HOT_500, liftTime, ids).Data().([]int64)
}
func getCacheCategoryMpHot() []int64 {
	data, _ := cache.Value(c_KEY_MP_HOT_500)
	if data == nil {
		return nil
	}
	return data.Data().([]int64)
}

//缓存指定分类置顶的公众号
func cacheCategoryMpStickTop(cid int64, mp []*Mp) []*Mp {
	liftTime, _ := time.ParseDuration("24h")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_MP_STICK_TOP, cid), liftTime, mp).Data().([]*Mp)
}
func getCacheCategoryMpStickTop(cid int64) []*Mp {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_MP_STICK_TOP, cid))
	if data == nil {
		return nil
	}
	return data.Data().([]*Mp)
}

//类目公众号总数
func cacheMpCountByCid(cid int64, count int64) int64 {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_MP_COUNT, cid), liftTime, count).Data().(int64)
}
func getCacheMpCountByCid(cid int64) int64 {
	data, _ := cache.Value(fmt.Sprintf("%s_%d", c_KEY_CATEGORY_MP_COUNT, cid))
	if data == nil {
		return 0
	}
	return data.Data().(int64)
}

//类目公众号缓存
func cacheMpListByCid(cid int64, page int, sortType int, data []*Mp) []*Mp {
	liftTime, _ := time.ParseDuration("10m")
	return cache.Add(fmt.Sprintf("%s_%d_%d_%d", c_KEY_CATEGORY_MP, cid, page, sortType), liftTime, data).Data().([]*Mp)
}
func getCacheMpListByCid(cid int64, page int, sortType int) []*Mp {
	data, _ := cache.Value(fmt.Sprintf("%s_%d_%d_%d", c_KEY_CATEGORY_MP, cid, page, sortType))
	if data == nil {
		return nil
	}
	return data.Data().([]*Mp)
}

//公众号排行榜
func cacheMpRank(mp []*Mp) []*Mp {
	liftTime, _ := time.ParseDuration("30m")
	return cache.Add(c_KEY_MP_RANK, liftTime, mp).Data().([]*Mp)
}
func getCacheMpRank() []*Mp {
	data, _ := cache.Value(c_KEY_MP_RANK)
	if data == nil {
		return nil
	}
	return data.Data().([]*Mp)
}

//=====================MP End=====================================
func randArticleArray(in []*Article) []*Article {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := rr.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}

func randMpArray(in []*Mp) []*Mp {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(in)
	for i := l - 1; i > 0; i-- {
		r := rr.Intn(i)
		in[r], in[i] = in[i], in[r]
	}
	return in
}

func genRandArr(size int, max int, result sets.Set) {
	if result.Size() == size || result.Size() == max {
		return
	}
	n := rand2.Rand(0, max-1)
	if !result.Contains(n) {
		result.Add(n)
	}
	genRandArr(size, max, result)
}
