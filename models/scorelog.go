package models

import (
	"time"
)

//积分类型
type ScoreChangeType int

const (
	ScoreChangeTypeRecharge ScoreChangeType = iota + 1
	ScoreChangeTypeStick
	ScoreChangeTypePublish
	ScoreChangeTypeUpdate
)

//积分修改日志
type ScoreLog struct {
	Id          int64     `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Uid         int64     `xorm:"BIGINT(20)"`                      //用户ID
	ChangeScore int       `xorm:"INT(11)"`                         //积分改变值
	Score       int       `xorm:"INT(11)"`                         //改变后的积分
	Type        int       `xorm:"default 1 INT(11)"`               //类型(充值1,消费2)
	Created     time.Time `xorm:"created DATETIME"`                //创建时间
}

//用户积分，关联用户名
type UserScoreLog struct {
	ScoreLog `xorm:"extends"`
	Name     string
	Account  string
}

func (UserScoreLog) TableName() string {
	return "score_log"
}

//根据具体条件获取数据，并且进行分页
//只支持ScoreLog里面字段"="条件的过滤
//使用OrderType进行排序
func GetScoreLogByUid(uid int64, pageSize int, page int, orderParams map[string]OrderType) (datas []*UserScoreLog, err error) {
	datas = []*UserScoreLog{}
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
	sess.Join("INNER", "user", "user.id = score_log.uid")
	if uid > 0 {
		sess.Where("score_log.uid = ?", uid)
	}
	err = sess.Find(&datas)
	return datas, err
}
