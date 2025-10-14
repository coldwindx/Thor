package test

import (
    "fmt"
    "reflect"
    "testing"
)

func TestGetMethodFromStruct(t *testing.T) {
    obj := &HelloWorld{Word: "World"}

    objType := reflect.TypeOf(obj)
    objValue := reflect.ValueOf(obj)

    for i := 0; i < objType.Elem().NumField(); i++ {
        field := objType.Elem().Field(i)
        value := objValue.Elem().Field(i)
        fmt.Println("[field] >>>", field.Name, ", CanSet:", value.CanSet())
    }

    for i := 0; i < objType.NumMethod(); i++ {
        method := objType.Method(i)
        fmt.Println("[method] >>>", method.Name, ", CanSet:", method.Func.CanSet())
    }
}
