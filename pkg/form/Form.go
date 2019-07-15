package form

import (
	"github.com/fatih/structs"
	"wxapp/utils/valid"
)

type Form interface {
	Rules() map[string]string
	Msgs() map[string]interface{}
	Preprocess()
}

type BaseForm struct {
	Csrf csrf `form:"csrf"`
}

type csrf struct {
	Token string
}

func Validate(f Form, formStruct interface{}) map[string]string {
	return valid.CheckMap(structs.Map(formStruct), f.Rules(), f.Msgs())
}
