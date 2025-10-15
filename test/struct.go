package test

/**
 * 测试用的结构体
 */
type IAnimal interface {
    GetName() string
    SetName(name string)
}

type DefaultAnimal struct {
    GetName func() string
    SetName func(name string)
}

type Cat struct {
    Name string
}
type Dog struct {
    Name string
}

func (c *Cat) GetName() string {
    return "cat:" + c.Name
}
func (c *Cat) SetName(name string) {
    c.Name = name
}
func (d *Dog) GetName() string {
    return "dog:" + d.Name
}
func (d *Dog) SetName(name string) {
    d.Name = name
}
