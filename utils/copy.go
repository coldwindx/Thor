package utils

import (
	"fmt"
	"reflect"
)

// DeepCopier deep copies a struct to/from a struct.
type DeepCopier struct {
	dst   any
	src   any
	err   error
	funcs []func(src any, dst any) error
}

// Copy sets source or destination.
func Copy(src any) *DeepCopier {
	return &DeepCopier{src: src, funcs: make([]func(src any, dst any) error, 0)}
}

func (dc *DeepCopier) With(funcs ...func(src any, dst any) error) *DeepCopier {
	dc.funcs = append(dc.funcs, funcs...)
	return dc
}

// To sets the destination.
func (dc *DeepCopier) To(dst any) error {
	dc.dst = dst
	return process(dc)
}

func (dc *DeepCopier) Create(t reflect.Type) any {
	dc.dst = reflect.New(t).Interface()
	if err := process(dc); err != nil {
		fmt.Printf("DeepCopier.Create failed, err: %v\n\n", err)
	}
	return dc.dst
}

// From sets the given the source as destination and destination as source.
func (dc *DeepCopier) From(src any) error {
	dc.src = src
	return process(dc)
}

// process deep copies a struct to/from a struct.
func process(dc *DeepCopier) error {
	srcValue := reflect.Indirect(reflect.ValueOf(dc.src))
	dstValue := reflect.Indirect(reflect.ValueOf(dc.dst))

	// dst必须是可寻址类型
	if !dstValue.CanAddr() {
		return fmt.Errorf("destination %+v is unaddressable", dstValue.Interface())
	}

	// 遍历dst的所有字段
	for i := 0; i < dstValue.NumField(); i++ {
		dstField := dstValue.Type().Field(i)
		// 寻找src中对应的字段
		srcField, ok := srcValue.Type().FieldByName(dstField.Name)
		// 找不到，或者src字段不可导出，或者无法将src字段的值转换为dst字段的类型，跳过
		if !ok || !srcField.IsExported() || !srcField.Type.AssignableTo(dstField.Type) {
			continue
		}
		// 复制字段值
		dstValue.Field(i).Set(srcValue.FieldByName(dstField.Name))
	}

	// 调用后处理函数
	for _, f := range dc.funcs {
		if err := f(dc.src, dc.dst); err != nil {
			return err
		}
	}
	return nil
}
