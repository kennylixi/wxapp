package models

import "fmt"

//城市表
type City struct {
	Id      int64  `xorm:"not null pk BIGINT(20)"` //主键ID
	Name    string `xorm:"VARCHAR(50)"`            //名字
	Group   string `xorm:"VARCHAR(10)"`            //简写
	Pid     int64  `xorm:"BIGINT(20)"`             //上级ID
	Deleted int    `xorm:"deleted INT(11)"`        //软删除标识
}

// 获取指定文章
func GetCityById(id int64) (*City, error) {
	citys := &City{}
	has, err := x.Id(id).Get(citys)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("City does not exist by id = %v", id)
	}
	return citys, nil
}

//根据具体条件获取数据，并且进行分页
//只支持City里面字段"="条件的过滤
//使用OrderType进行排序
func GetCityByCnd(city *City, pageSize int, page int, orderParams map[string]OrderType) (citys []*City, err error) {
	citys = []*City{}
	sess := x.NewSession()
	defer sess.Close()
	//是否进行分页
	if pageSize > 0 && page > 0 {
		sess.Limit(pageSize, (page-1)*pageSize)
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
	err = sess.Find(&citys, city)
	return citys, err
}

//根据pid获取城市列表
func GetCityByPid(pid int64) (citys []*City, err error) {
	citys = []*City{}
	err = x.MustCols("pid").Find(&citys, &City{Pid: pid})
	return citys, err
}
