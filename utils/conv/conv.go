// 类型转换.
// 如果给定的interface{}参数不是指定转换的输出类型，那么会进行强制转换，效率会比较低，
// 建议已知类型的转换自行调用相关方法来单独处理。
package conv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Bytes(i interface{}) []byte {
	if i == nil {
		return nil
	}
	if r, ok := i.([]byte); ok {
		return r
	} else {
		return []byte(String(i))
	}
}

func String(i interface{}) string {
	if i == nil {
		return ""
	}
	if r, ok := i.(string); ok {
		return r
	} else {
		return fmt.Sprintf("%v", i)
	}
}

func Strings(i interface{}) []string {
	if i == nil {
		return nil
	}
	if r, ok := i.([]string); ok {
		return r
	} else if r, ok := i.([]interface{}); ok {
		strs := make([]string, len(r))
		for k, v := range r {
			strs[k] = String(v)
		}
		return strs
	}
	return []string{fmt.Sprintf("%v", i)}
}

//false: "", 0, false, off
func Bool(i interface{}) bool {
	if i == nil {
		return false
	}
	if v, ok := i.(bool); ok {
		return v
	}
	if s := String(i); s != "" && s != "0" && s != "false" && s != "off" {
		return true
	}
	return false
}

func Int(i interface{}) int {
	if i == nil {
		return 0
	}
	if v, ok := i.(int); ok {
		return v
	}
	v, _ := strconv.Atoi(String(i))
	return v
}

func Int64(i interface{}) int64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(int64); ok {
		return v
	}
	v, _ := strconv.ParseInt(String(i), 10, 64)
	return v
}

func Uint(i interface{}) uint {
	if i == nil {
		return 0
	}
	if v, ok := i.(uint); ok {
		return v
	}
	v, _ := strconv.ParseUint(String(i), 10, 8)
	return uint(v)
}

func Float32(i interface{}) float32 {
	if i == nil {
		return 0
	}
	if v, ok := i.(float32); ok {
		return v
	}
	v, _ := strconv.ParseFloat(String(i), 8)
	return float32(v)
}

func Float64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	if v, ok := i.(float64); ok {
		return v
	}
	v, _ := strconv.ParseFloat(String(i), 8)
	return v
}

//用map填充结构
func FillStructObj(data map[string]interface{}, obj interface{}, isStrict bool) error {
	for k, v := range data {
		con, err := setField(obj, k, v, isStrict)
		if con == false && err != nil {
			return err
		}
	}
	return nil
}

//用map填充结构
func FillStructStr(data map[string]string, obj interface{}, isStrict bool) error {
	for k, v := range data {
		con, err := setField(obj, k, v, isStrict)
		if con == false && err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
//isStrict严格模式true，false表示结构没有的字段直接跳过
func setField(obj interface{}, name string, value interface{}, isStrict bool) (bool, error) {
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return !isStrict, fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return false, fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = typeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return false, err
		}
	}

	structFieldValue.Set(val)
	return true, nil
}

//类型转换
func typeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		if value == "" {
			value = "0"
		}
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		if value == "" {
			value = "0"
		}
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		if value == "" {
			value = "0"
		}
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
