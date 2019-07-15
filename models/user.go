package models

import (
	"fmt"
	"time"
	"wxapp/utils/conv"
)

//用户类型
const (
	USER_TYPE_COMMON int = iota + 1
	USER_TYPE_ADMIN
)

//用户状态
const (
	USER_STATUS_NORMAL int = iota + 1
	USER_STATUS_UMNORMAL
)

//用户表
type User struct {
	Id      int64     `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Name    string    `xorm:"VARCHAR(100)"`                    //昵称
	Account string    `xorm:"VARCHAR(100) unique"`             //账号
	Pwd     string    `xorm:"VARCHAR(20)"`                     //密码
	Score   int       `xorm:"INT(11)"`                         //积分
	Type    int       `xorm:"default 1 INT(1)"`                //类型
	Wx      string    `xorm:"VARCHAR(50)"`                     //微信
	Qq      string    `xorm:"VARCHAR(50)"`                     //QQ
	Status  int       `xorm:"default 1 INT(1)"`                //状态
	Created time.Time `xorm:"created DATETIME"`                //创建时间
}

// 获取指定用户
func GetUserById(id int64) (*User, error) {
	datas := &User{}
	has, err := x.Id(id).Get(datas)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, fmt.Errorf("User does not exist by id = %v", id)
	}
	return datas, nil
}

//根据具体条件获取数据，并且进行分页
//只支持User里面字段"="条件的过滤
//使用OrderType进行排序
func GetUserByCnd(data *User, pageSize int, page int, orderParams map[string]OrderType) (datas []*User, err error) {
	datas = []*User{}
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
	err = sess.Find(&datas, data)
	return datas, err
}

//根据账号密码获取用户信息
func GetByAccountAndPwd(account string, pwd string) (*User, error) {
	users, err := GetUserByCnd(&User{Account: account, Pwd: pwd}, 0, 0, nil)
	if len(users) > 0 {
		return users[0], err
	}
	return nil, err
}

//充值
func Recharge(id int64, value int) error {
	user, err := GetUserById(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("not exits user id : %v", id)
	}
	e := NewSession()
	defer e.Close()
	if err := e.Begin(); err != nil {
		return err
	}
	user = &User{Score: user.Score + value}
	_, err = Update(user, User{Id: id}, e)
	if err != nil {
		e.Rollback()
		return err
	}
	scoreLog := &ScoreLog{Uid: id, ChangeScore: value, Score: user.Score, Type: conv.Int(ScoreChangeTypeRecharge), Created: time.Time{}}
	_, err = Insert(scoreLog, e)
	if err != nil {
		e.Rollback()
		return err
	}
	return e.Commit()
}
