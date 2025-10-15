package inject

import (
    "fmt"
    "github.com/samber/lo"
    "reflect"
)

// Object 表示一个受管理的bean对象
type Object struct {
    Name      string        // Bean名称
    Value     any           // Bean实例，必须是指针类型
    RfType    reflect.Type  // Bean实例的反射类型
    RfValue   reflect.Value // Bean实例的反射值
    Completed bool          // 完整的Bean实例，无需依赖注入
}

// Graph 表示一组bean对象的管理容器
type Graph struct {
    objects map[string]*Object // 所有受管理的bean对象，key为bean名称
}

// NewGraph 创建一个新的Graph实例
func NewGraph() *Graph {
    return &Graph{objects: make(map[string]*Object)}
}

// Provide 向容器中添加一个或多个bean对象
func (g *Graph) Provide(objs ...*Object) {
    for _, obj := range objs {
        obj.RfType = reflect.TypeOf(obj.Value)
        obj.RfValue = reflect.ValueOf(obj.Value)
        // 检查Bean实例是否是指针类型
        if obj.RfType.Kind() != reflect.Ptr {
            panic("Bean实例必须是指针类型")
        }
        // 检查Bean的名称是否为空，如果为空，则根据bean的类型构造名称
        obj.Name = lo.Ternary(len(obj.Name) == 0, obj.RfType.Elem().Name(), obj.Name)
        // 检查是否有相同名称的bean对象
        if _, ok := g.objects[obj.Name]; ok {
            panic("provided two instances named `" + obj.Name + "`")
        }
        // 根据Bean的名称，加入容器管理
        g.objects[obj.Name] = obj
    }
}

func (g *Graph) Populate() error {
    // 遍历所有bean对象，进行依赖注入
    for _, obj := range g.objects {
        // 检查Bean是否已经完成依赖注入，节省注入成本
        if obj.Completed {
            continue
        }
        // 拿到Bean实例的所有属性值，准备注入
        for i := 0; i < obj.RfType.Elem().NumField(); i++ {
            fieldValue := obj.RfValue.Elem().Field(i)
            // 跳过不能设置的属性
            if !fieldValue.CanSet() {
                continue
            }
            // 分别获取属性的 类型 和 标签
            field := obj.RfType.Elem().Field(i)
            tag, ok := field.Tag.Lookup("inject")
            // 跳过没有inject标签的属性
            if !ok {
                continue
            }
            // 根据属性的 类型 和 标签，进行依赖注入
            factory.Inject(g, fieldValue, field.Type, tag)
        }
    }
    return nil
}

func (g *Graph) GetByName(name string) any {
    return g.objects[name].Value
}

/**
 * 依赖注入器
 * 负责根据bean对象的属性标签，进行依赖注入
 */
var factory = &InjectFactory{
    injector: []Injector{
        &NamedStructInjector{},
        &UnnamedStructInjector{},
        &NamedInterfaceInjector{},
        &UnnamedInterfaceInjector{},
        &ListInjector{},
        &MapInjector{},
    },
}

type InjectFactory struct {
    injector []Injector
}

func (f *InjectFactory) Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, fieldTag string) {
    for _, inject := range f.injector {
        if inject.Validate(fieldType, fieldTag) {
            inject.Inject(g, fieldVal, fieldType, fieldTag)
            return
        }
    }
    panic(fmt.Sprintf("did not find suitable injector for field [%s]", fieldType.Name()))
}

type Injector interface {
    Validate(fieldType reflect.Type, fieldTag string) bool
    Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, fieldTag string)
}

type NamedStructInjector struct{}
type UnnamedStructInjector struct{}
type NamedInterfaceInjector struct{}
type UnnamedInterfaceInjector struct{}
type ListInjector struct{}
type MapInjector struct{}

func (n *NamedInterfaceInjector) Validate(fieldType reflect.Type, fieldTag string) bool {
    // 注入属性必须是接口类型
    if fieldType.Kind() != reflect.Interface {
        return false
    }
    // 注入属性标签不能为空
    if len(fieldTag) == 0 {
        return false
    }
    // 可执行依赖注入
    return true
}

func (n *NamedInterfaceInjector) Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, fieldTag string) {
    // 从容器中根据属性标签获取管理的bean对象
    obj, ok := g.objects[fieldTag]
    // 检查容器中是否存在目标bean对象
    if !ok {
        panic(fmt.Sprintf("bean `%s` not found", fieldTag))
    }
    // 检查依赖的bean对象是否实现了目标接口
    if !obj.RfType.Implements(fieldType) {
        panic(fmt.Sprintf("bean `%s` does not implement interface `%s`", obj.Name, fieldType.Name()))
    }
    // 进行依赖注入
    fieldVal.Set(reflect.ValueOf(obj.Value))
}
func (i *UnnamedInterfaceInjector) Validate(fieldType reflect.Type, fieldTag string) bool {
    // 属性标签必须为空
    if len(fieldTag) != 0 {
        return false
    }
    // 注入属性必须是接口类型
    if fieldType.Kind() != reflect.Interface {
        return false
    }
    // 可执行依赖注入
    return true
}

func (i *UnnamedInterfaceInjector) Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, fieldTag string) {
    // 从容器中循环获取管理的bean对象
    for _, obj := range g.objects {
        // 检查依赖的bean对象必须实现目标接口
        if obj.RfType.Implements(fieldType) {
            // 进行依赖注入
            fieldVal.Set(reflect.ValueOf(obj.Value))
        }
    }
}

func (m *MapInjector) Validate(fieldType reflect.Type, fieldTag string) bool {
    // 注入属性必须是map类型
    if fieldType.Kind() != reflect.Map {
        return false
    }
    // 注入类型是map时，属性标签必须为空
    if len(fieldTag) != 0 {
        panic("map type inject must have no tag.")
    }

    // 可执行依赖注入
    return true
}

func (m *MapInjector) Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, _ string) {
    mValues := reflect.MakeMap(fieldType)
    // 从容器中循环获取管理的bean对象
    for _, obj := range g.objects {
        // 检查依赖的bean对象是否是目标属性的类型或子类型
        if obj.RfType.AssignableTo(fieldType.Elem()) {
            // 进行依赖注入
            mValues.SetMapIndex(reflect.ValueOf(obj.Name), reflect.ValueOf(obj.Value))
        }
    }
    // 设置属性值
    fieldVal.Set(mValues)
}

func (l *ListInjector) Validate(fieldType reflect.Type, fieldTag string) bool {
    // 注入属性必须是切片类型
    if fieldType.Kind() != reflect.Slice {
        return false
    }
    // 属性标签必须为空
    if len(fieldTag) != 0 {
        panic("slice type inject must have no tag.")
    }

    // 可执行依赖注入
    return true
}

func (l *ListInjector) Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, _ string) {
    lValues := reflect.MakeSlice(fieldType, 0, 0)
    // 从容器中循环获取管理的bean对象
    for _, obj := range g.objects {
        // 检查依赖的bean对象是否是目标属性的类型或子类型
        if obj.RfType.AssignableTo(fieldType.Elem()) {
            // 进行依赖注入
            lValues = reflect.Append(lValues, reflect.ValueOf(obj.Value))
        }
    }
    // 设置属性值
    fieldVal.Set(lValues)
}

func (t *UnnamedStructInjector) Validate(fieldType reflect.Type, fieldTag string) bool {
    // 属性标签必须为空
    if len(fieldTag) != 0 {
        return false
    }
    // 注入属性必须是指针类型
    if fieldType.Kind() != reflect.Ptr {
        return false
    }
    // 可执行依赖注入
    return true
}

func (t *UnnamedStructInjector) Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, _ string) {
    // 从容器中循环获取管理的bean对象
    for _, obj := range g.objects {
        // 检查依赖的bean对象是否是目标属性的类型或子类型
        if obj.RfType.Elem().AssignableTo(fieldType.Elem()) {
            // 进行依赖注入
            fieldVal.Set(reflect.ValueOf(obj.Value))
            return
        }
    }
    // 未找到合适的bean对象，抛出异常
    panic(fmt.Sprintf("did not find object of type [%s] required by type [%s]", fieldType.Name(), fieldType.Name()))
}

func (n *NamedStructInjector) Inject(g *Graph, fieldVal reflect.Value, fieldType reflect.Type, fieldTag string) {
    // 从容器中获取依赖的bean对象
    existing, ok := g.objects[fieldTag]
    if !ok {
        panic(fmt.Sprintf("did not find object named [%s] required by field [%s] in type [%s]", fieldTag, fieldVal.Type().Name(), fieldType.Name()))
    }

    // 检查依赖的bean对象是否是目标属性的类型或子类型
    if !existing.RfType.Elem().AssignableTo(fieldType.Elem()) {
        panic(fmt.Sprintf("object named [%s] of type [%s] is not assignable to field [%s] in type [%s]", existing.Name, existing.RfType.Name(), fieldVal.Type().Name(), fieldType.Name()))
    }

    // 进行依赖注入
    fieldVal.Set(reflect.ValueOf(existing.Value))
}

func (n *NamedStructInjector) Validate(fieldType reflect.Type, fieldTag string) bool {
    // 属性标签不能为空
    if len(fieldTag) == 0 {
        return false
    }
    // 注入属性必须是指针类型
    if fieldType.Kind() != reflect.Ptr {
        return false
    }
    // 可执行依赖注入
    return true
}
