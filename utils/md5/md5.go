// MD5
package md5

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	log "gopkg.in/clog.v1"
	"io"
	"os"
	"reflect"
)

// 将任意类型的变量进行md5摘要(注意map等非排序变量造成的不同结果)
func Encode(v interface{}) string {
	h := md5.New()
	if "string" == reflect.TypeOf(v).String() {
		h.Write([]byte(v.(string)))
	} else {
		b, err := json.Marshal(v)
		if err != nil {
			return ""
		} else {
			h.Write(b)
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 将字符串进行MD5哈希摘要计算
func EncodeString(v string) string {
	h := md5.New()
	h.Write([]byte(v))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 将文件内容进行MD5哈希摘要计算
func EncodeFile(path string) string {
	f, e := os.Open(path)
	if e != nil {
		log.Fatal(0, "file is not exist : v%", e)
	}
	h := md5.New()
	_, e = io.Copy(h, f)
	if e != nil {
		log.Fatal(0, "file is copy false : v%", e)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
