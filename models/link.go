package models

import "fmt"

type LinkType int

const (
	LinkTypePage LinkType = iota + 1 //友情链接页
	LinkTypeHome                     //主页友情链接
)

//友情链接表
type Link struct {
	Id    int64  `xorm:"not null pk autoincr BIGINT(20)"` //友情ID
	Title string `xorm:"VARCHAR(50)"`                     //友情名
	Url   string `xorm:"VARCHAR(255)"`                    //友情url
	Sort  int    `xorm:"default 0 INT(11)"`               //排序
	Type  int    `xorm:"default 1 INT(11)"`               //类型
}

// 获取指定友情链接
func GetLinkById(id int64) (*Link, error) {
	datas := &Link{}
	has, err := x.Id(id).Get(datas)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("Link does not exist by id = %v", id)
	}
	return datas, nil
}

//根据具体条件获取数据，并且进行分页
//只支持Link里面字段"="条件的过滤
//使用OrderType进行排序
func GetLinkByCnd(data *Link, pageSize int, offset int, orderParams map[string]OrderType) (datas []*Link, err error) {
	datas = []*Link{}
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
