package utils

import (
	"Thor/ctx"
	"errors"
	"github.com/samber/lo"
	"reflect"
)
import "github.com/facebookgo/inject"

func ScanInject(name string, obj interface{}) {
	var v = reflect.ValueOf(obj)
	if v.Type().Kind() != reflect.Ptr {
		panic("对象必须是指针！")
	}
	var t = v.Type()
	if t.Elem().Kind() == reflect.Struct {
		for i := 0; i < t.Elem().NumField(); i++ {
			var typeItem = t.Elem().Field(i)
			var bean, ok = typeItem.Tag.Lookup("bean")
			if ok {
				if bean == "" {
					panic(t.String() + "." + typeItem.Name + "提供的bean的tag必须提供名称")
				}
				err := ctx.Beans.Provide(&inject.Object{
					Name:  bean,
					Value: v.Elem().Field(i).Addr().Interface(),
				})
				if err != nil {
					panic("Bean Inject Error. " + err.Error())
				}
			}
		}
	}

	err := ctx.Beans.Provide(&inject.Object{Name: name, Value: obj})
	if err != nil {
		panic("Bean Inject Error. " + err.Error())
	}
}

func GetBean[T any](name string) (*T, error) {
	objs := ctx.Beans.Objects()
	obj, flag := lo.Find(objs, func(item *inject.Object) bool { return item.Name == name })
	if !flag {
		return nil, errors.New("found bean from inject graph error, bean name: " + name)
	}
	t := obj.Value.(*T)
	return t, nil
}
