package models

//专题文章关联表
type SpecialArtice struct {
	Id  int64 `xorm:"not null pk autoincr BIGINT(20)"` //主键ID
	Aid int64 `xorm:"BIGINT(20)"`                      //文章ID
	Sid int64 `xorm:"BIGINT(20)"`                      //专题ID
}

//文章关联专题ID
type ArticleSpecial struct {
	Article `xorm:"extends"`
	Id      int64
}

func (ArticleSpecial) TableName() string {
	return "article"
}
