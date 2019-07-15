package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"wxapp/config"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/kataras/golog"
	"math"
)

// Engine represents a XORM engine or session.
type Engine interface {
	Delete(interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	Id(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Iterate(interface{}, xorm.IterFunc) error
	Sql(string, ...interface{}) *xorm.Session
	Table(interface{}) *xorm.Session
	Where(interface{}, ...interface{}) *xorm.Session
	Update(interface{}, ...interface{}) (int64, error)
}

//排序类型
type OrderType int

const (
	ORDER_TYPE_ASE  OrderType = iota //升序
	ORDER_TYPE_DESC                  //倒序
)

var (
	x                             *xorm.Engine
	tables                        []interface{}
	isInitDBConfig                bool
	dbHost, dbName, dbUser, dbPwd string
)

func init() {
	tables = append(tables, new(User), new(Special),
		new(SpecialArtice), new(ScoreLog), new(City), new(Category), new(Article),
		new(Link), new(Apply), new(Recommend), new(Evaluate), new(Mp))

	gonicNames := []string{"SSL"}
	for _, name := range gonicNames {
		core.LintGonicMapper[name] = true
	}
}

func loadConfigs() {
	if isInitDBConfig {
		return
	}
	sec := config.Cfg.Section("database")
	dbHost = sec.Key("HOST").String()
	dbName = sec.Key("NAME").String()
	dbUser = sec.Key("USER").String()
	dbPwd = sec.Key("PASSWD").String()
	isInitDBConfig = true
}

func getEngine() (*xorm.Engine, error) {
	loadConfigs()

	connStr := ""
	var Param string = "?"
	if strings.Contains(dbName, Param) {
		Param = "&"
	}
	connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s%scharset=utf8mb4&parseTime=true",
		dbUser, dbPwd, dbHost, dbName, Param)

	return xorm.NewEngine("mysql", connStr)
}

func NewTestEngine(x *xorm.Engine) (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})
	return x.StoreEngine("InnoDB").Sync2(tables...)
}

func SetEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})

	if err != nil {
		return fmt.Errorf("Fail to create 'xorm.log': %v", err)
	}

	if config.ProdMode {
		x.SetLogger(xorm.NewSimpleLogger3(golog.Default.Printer.Output, xorm.DEFAULT_LOG_PREFIX, xorm.DEFAULT_LOG_FLAG, core.LOG_WARNING))
	} else {
		x.SetLogger(xorm.NewSimpleLogger(golog.Default.Printer.Output))
	}
	x.SetMaxOpenConns(5)
	x.ShowSQL(true)
	return nil
}

func NewEngine() (err error) {
	if err = SetEngine(); err != nil {
		return err
	}

	if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}
	return nil
}

func Close() {
	x.Close()
}

func Ping() error {
	return x.Ping()
}

func Insert(bean interface{}, e Engine) (int64, error) {
	if e == nil {
		e = x
	}
	return e.InsertOne(bean)
}

func Delete(bean interface{}, e Engine) (int64, error) {
	if e == nil {
		e = x
	}
	return e.Delete(bean)
}

func Update(bean interface{}, cnd interface{}, e Engine) (int64, error) {
	if e == nil {
		e = x
	}
	return e.Update(bean, cnd)
}

func Get(bean interface{}) (bool, error) {
	return x.Get(bean)
}

func Count(bean interface{}) int64 {
	count, _ := x.Count(bean)
	return count
}

//创建Session
func NewSession() *xorm.Session {
	return x.NewSession()
}

// The version table. Should have only one row with id==1
type Version struct {
	ID      int64
	Version int64
}

// DumpDatabase dumps all data from database to file system in JSON format.
func DumpDatabase(dirPath string) (err error) {
	os.MkdirAll(dirPath, os.ModePerm)
	// Purposely create a local variable to not modify global variable
	tables := append(tables, new(Version))
	for _, table := range tables {
		tableName := strings.TrimPrefix(fmt.Sprintf("%T", table), "*models.")
		tableFile := path.Join(dirPath, tableName+".json")
		f, err := os.Create(tableFile)
		if err != nil {
			return fmt.Errorf("fail to create JSON file: %v", err)
		}

		if err = x.Asc("id").Iterate(table, func(idx int, bean interface{}) (err error) {
			enc := json.NewEncoder(f)
			return enc.Encode(bean)
		}); err != nil {
			f.Close()
			return fmt.Errorf("fail to dump table '%s': %v", tableName, err)
		}
		f.Close()
	}
	return nil
}

//分页方法，根据传递过来的页数，每页数，总数，返回分页的内容 7个页数 前 1，2，3，4，5 后 的格式返回,小于5页返回具体页数
func Paginator(page, prepage int, nums int64) map[string]interface{} {

	var firstpage int //前一页地址
	var lastpage int  //后一页地址
	//根据nums总数，和prepage每页数量 生成分页总数
	totalpages := int(math.Ceil(float64(nums) / float64(prepage))) //page总数
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= 3 && totalpages > 5:
		if page > totalpages-3 && totalpages > 5 { //最后5页
			start := totalpages - 5 + 1
			firstpage = page - 1
			lastpage = int(math.Min(float64(totalpages), float64(page+1)))
			pages = make([]int, 5)
			for i, _ := range pages {
				pages[i] = start + i
			}
		} else {
			start := page - 3 + 1
			pages = make([]int, 5)
			firstpage = page - 1
			for i, _ := range pages {
				pages[i] = start + i
			}
			firstpage = page - 1
			lastpage = page + 1
		}

	default:
		pages = make([]int, int(math.Min(float64(5), float64(totalpages))))
		for i, _ := range pages {
			pages[i] = i + 1
		}
		firstpage = int(math.Max(float64(1), float64(page-1)))
		lastpage = page + 1
	}
	if lastpage > totalpages {
		lastpage = totalpages
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalpages
	paginatorMap["firstpage"] = firstpage
	paginatorMap["lastpage"] = lastpage
	paginatorMap["currpage"] = page
	paginatorMap["totals"] = nums
	return paginatorMap
}
