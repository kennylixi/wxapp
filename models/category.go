package models

import (
	"fmt"
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/golog"
	"wxapp/pkg/tree"
	"wxapp/utils/conv"
	"github.com/emirpasic/gods/lists"
)

const (
	CATE_ROOT      int64  = 0
	CATE_APP_ID    int64  = 1
	CATE_GAME_ID   int64  = 2
	CATE_MP_ID     int64  = 3
	CATE_APP_NAME  string = "app"
	CATE_GAME_NAME string = "game"
	CATE_MP_NAME   string = "mp"
)

//分类表
type Category struct {
	Id         int64  `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Title      string `xorm:"VARCHAR(50)"`                     //类目名
	ShortTitle string `xorm:"VARCHAR(50)"`                     //简称类目名
	Short      string `xorm:"VARCHAR(50)"`                     //拼音
	Pid        int64  `xorm:"default 0 BIGINT(20)"`            //上级ID
	Sort       int    `xorm:"default 0 INT(11)"`               //排序号
	Status     int    `xorm:"default 1 INT(11)"`               //是否隐藏
	Deleted    int    `xorm:"deleted INT(11)"`                 //软删除标识
	SeoKey     string `xorm:"VARCHAR(150)"`                    //SEO搜索关键字
	SeoDesc    string `xorm:"VARCHAR(150)"`                    //SEO描述

	Article []*Article `xorm:"-"`
}

// 获取指定文章
func GetCategoryById(id int64) (*Category, error) {
	datas := &Category{}
	has, err := x.Id(id).Get(datas)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Category does not exist by id = %v", id)
	}
	return datas, nil
}

// 获取指定文章
func GetCategoryByShort(pid int64,short string) (*Category, error) {
	datas := &Category{Short: short, Pid:pid}
	has, err := x.Get(datas)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Category does not exist by short = %v pid = %v", short, pid)
	}
	return datas, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Category里面字段"="条件的过滤
//使用OrderType进行排序
func GetCategoryByCnd(data *Category, pageSize int, offset int, orderParams map[string]OrderType) (datas []*Category, err error) {
	datas = []*Category{}
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

func GetCategoryTreeByIds(ids lists.List) *tree.Tree {
	categorys, err := GetCategoryByCnd(&Category{}, 0, 0, nil)
	if err != nil {
		golog.Errorf("GetCategoryTree : %v", err)
		return nil
	}
	trees := make([]*tree.Tree, 0, len(categorys))
	for _, category := range categorys {
		state := map[string]interface{}{"opened": true}
		t := tree.NewTree(conv.String(category.Id), conv.String(category.Pid), category.Title, state)
		trees = append(trees, t)
	}
	root := &tree.Tree{}
	root.Id = "-1"
	root.ParentId = ""
	root.HasParent = false
	root.HasChildren = true
	root.Checked = true
	root.Children = make([]*tree.Tree,0)
	root.Text = "顶级节点"
	root.State = map[string]interface{}{"opened": true}

	ts:=tree.BuildList(trees,"0")
	for _, t := range ts {
		if ids.Contains(t.Id) {
			root.Children = append(root.Children, t)
		}
	}
	return root
}
func GetCategoryTree() *tree.Tree {
	categorys, err := GetCategoryByCnd(&Category{}, 0, 0, nil)
	if err != nil {
		golog.Errorf("GetCategoryTree : %v", err)
		return nil
	}
	trees := make([]*tree.Tree, 0, len(categorys))
	for _, category := range categorys {
		state := map[string]interface{}{"opened": true}
		t := tree.NewTree(conv.String(category.Id), conv.String(category.Pid), category.Title, state)
		trees = append(trees, t)
	}
	return tree.Build(trees)
}

//获取指定父节点下面的所有数据的ID
func getCategoryChildIds(id int64) []int64 {
	ids := getCacheCategoryChildIds(id)
	if ids != nil {
		return ids
	}
	tree := GetCategoryTree()

	idset := hashset.New()
	idset.Add(id)

	getCategoryChildIdsRecursion(tree, conv.String(id), idset)

	ids = make([]int64, idset.Size())
	for index, cid := range idset.Values() {
		ids[index] = conv.Int64(cid)
	}
	return cacheCategoryChildIds(id, ids)
}

//获取指定父节点下面的所有数据的ID 递归函数
func getCategoryChildIdsRecursion(tree *tree.Tree, pid string, ids sets.Set) {
	for _, child := range tree.Children {
		if child.ParentId == pid {
			ids.Add(conv.Int64(child.Id))
			getCategoryChildIdsRecursion(child, child.Id, ids)
		} else {
			getCategoryChildIdsRecursion(child, pid, ids)
		}
	}
}

//获取导航分类列表
func GetCategoryNavList(pid int64) []*Category {
	categoryList := getCacheCategoryNav(pid)
	if categoryList != nil {
		return categoryList
	}
	categoryList, err := GetCategoryByCnd(&Category{Pid: pid, Status: 1}, 0, 0, map[string]OrderType{"sort": ORDER_TYPE_ASE})
	if err != nil {
		golog.Errorf("GetCategoryNavList error : %v", err)
	}
	return cacheCategoryNav(pid, categoryList)
}
