package inject

/**
 * 根据facebook的inject库，重新实现的依赖注入工具
 */
import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
)

// Logger allows for simple logging as inject traverses and populates the
// object graph.
type Logger interface {
	Debugf(format string, v ...interface{})
}

// Populate is a short-hand for populating a graph with the given incomplete
// object values.
func Populate(values ...interface{}) error {
	var g Graph
	for _, v := range values {
		if err := g.Provide(&Object{Value: v}); err != nil {
			return err
		}
	}
	return g.Populate()
}

// An Object in the Graph.
type Object struct {
	Value        interface{}
	Name         string             // Optional
	Complete     bool               // If true, the Value will be considered complete
	Fields       map[string]*Object // Populated with the field names that were injected and their corresponding *Object.
	reflectType  reflect.Type
	reflectValue reflect.Value
	private      bool // If true, the Value will not be used and will only be populated
	created      bool // If true, the Object was created by us
	embedded     bool // If true, the Object is an embedded struct provided internally
}

// String representation suitable for human consumption.
func (o *Object) String() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, o.reflectType)
	if o.Name != "" {
		fmt.Fprintf(&buf, " named %s", o.Name)
	}
	return buf.String()
}

func (o *Object) addDep(field string, dep *Object) {
	if o.Fields == nil {
		o.Fields = make(map[string]*Object)
	}
	o.Fields[field] = dep
}

// The Graph of Objects.
type Graph struct {
	Logger Logger // Optional, will trigger debug logging.
	named  map[string]*Object
}

// Provide objects to the Graph. The Object documentation describes
// the impact of various fields.
func (g *Graph) Provide(objects ...*Object) error {
	for _, o := range objects {
		o.reflectType = reflect.TypeOf(o.Value)
		o.reflectValue = reflect.ValueOf(o.Value)

		if o.Fields != nil {
			return fmt.Errorf(
				"fields were specified on object %s when it was provided",
				o,
			)
		}

		if o.Name == "" {
			if !isStructPtr(o.reflectType) {
				return fmt.Errorf(
					"expected unnamed object value to be a pointer to a struct but got type %s "+
						"with value %v",
					o.reflectType,
					o.Value,
				)
			}
			// 当没有指定Name时，默认Name为结构体类型名
			o.Name = o.reflectType.Elem().String()
			//fmt.Println("default name:", o.reflectType.Elem().String())
		}

		if g.named == nil {
			g.named = make(map[string]*Object)
		}

		if g.named[o.Name] != nil {
			return fmt.Errorf("provided two instances named %s", o.Name)
		}
		g.named[o.Name] = o

		if g.Logger != nil {
			if o.created {
				g.Logger.Debugf("created %s", o)
			} else if o.embedded {
				g.Logger.Debugf("provided embedded %s", o)
			} else {
				g.Logger.Debugf("provided %s", o)
			}
		}
	}
	return nil
}

// Populate the incomplete Objects.
func (g *Graph) Populate() error {
	// step1. 注入所有struct类型的Object
	for _, o := range g.named {
		if o.Complete {
			continue
		}

		// 注入指针类型的Object
		if err := g.populateExplicit(o); err != nil {
			return err
		}
	}

	// step2. 注入所有接口类型的Object
	for _, o := range g.named {
		if o.Complete {
			continue
		}

		if err := g.populateInterface(o); err != nil {
			return err
		}
	}

	// step3. 注入所有map类型的Object
	for _, o := range g.named {
		if o.Complete {
			continue
		}

		if err := g.populateMap(o); err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph) populateExplicit(o *Object) error {
	for i := 0; i < o.reflectValue.Elem().NumField(); i++ {
		field := o.reflectValue.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := o.reflectType.Elem().Field(i).Tag
		fieldName := o.reflectType.Elem().Field(i).Name
		// Skip fields without a ttag.
		ttag := parseTagV2(fieldTag)
		if ttag == nil {
			continue
		}

		// Cannot be used with unexported fields.
		if !field.CanSet() {
			return fmt.Errorf(
				"inject requested on unexported field `%s` in type `%s`",
				o.reflectType.Elem().Field(i).Name,
				o.reflectType,
			)
		}

		// 第一轮注入，只注入struct类型的字段
		if fieldType.Kind() != reflect.Struct && fieldType.Kind() != reflect.Ptr {
			continue
		}

		// Inline ttag on anything besides a struct is considered invalid.
		if ttag.Inline && fieldType.Kind() != reflect.Struct {
			return fmt.Errorf(
				"inline requested on non inlined field %s in type %s",
				o.reflectType.Elem().Field(i).Name,
				o.reflectType,
			)
		}

		// Don't overwrite existing values.
		if !isNilOrZero(field, fieldType) {
			continue
		}

		// 没有指定Name时，默认Name为字段类型名
		if ttag.Name == "" {
			ttag.Name = fieldType.Elem().String()
			//fmt.Printf("default inject name: %s, field name: %s\n", ttag.Name, fieldName)
		}

		existing := g.named[ttag.Name]
		if existing == nil {
			return fmt.Errorf(
				"did not find object named %s required by field %s in type %s",
				ttag.Name, fieldName, o.reflectType,
			)
		}

		if !existing.reflectType.AssignableTo(fieldType) {
			return fmt.Errorf(
				"object ttag named `%s` of type `%s` is not assignable to field `%s` in type `%s`",
				ttag.Name,
				fieldType,
				fieldName,
				o.reflectType,
			)
		}

		field.Set(reflect.ValueOf(existing.Value))
		if g.Logger != nil {
			g.Logger.Debugf("assigned %s to field %s in %s", existing, o.reflectType.Elem().Field(i).Name, o)
		}
		o.addDep(fieldName, existing)
	}
	return nil
}

func (g *Graph) populateInterface(o *Object) error {
loop:
	for i := 0; i < o.reflectValue.Elem().NumField(); i++ {
		field := o.reflectValue.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := o.reflectType.Elem().Field(i).Tag
		fieldName := o.reflectType.Elem().Field(i).Name

		tag := parseTagV2(fieldTag)
		if tag == nil {
			continue
		}
		// We only handle interface injection here. Other cases including errors
		// are handled in the first pass when we inject pointers.
		if fieldType.Kind() != reflect.Interface {
			continue
		}

		// Interface injection can't be private because we can't instantiate new
		// instances of an interface.
		if tag.Private {
			return fmt.Errorf(
				"found private inject tag on interface field %s in type %s",
				o.reflectType.Elem().Field(i).Name,
				o.reflectType,
			)
		}

		// Don't overwrite existing values.
		if !isNilOrZero(field, fieldType) {
			continue
		}

		// Find one, and only one assignable value for the field.
		for _, existing := range g.named {
			if existing.private || (0 < len(tag.Name) && existing.Name != tag.Name) {
				continue
			}
			if !existing.reflectType.AssignableTo(fieldType) {
				continue
			}
			field.Set(reflect.ValueOf(existing.Value))
			if g.Logger != nil {
				g.Logger.Debugf(
					"assigned existing %s to interface field %s in %s",
					existing,
					o.reflectType.Elem().Field(i).Name,
					o,
				)
			}
			o.addDep(fieldName, existing)
			break loop
		}

		// If we didn't find an assignable value, we're missing something.
		return fmt.Errorf(
			"found no assignable value for field %s in type %s",
			o.reflectType.Elem().Field(i).Name,
			o.reflectType,
		)
	}
	return nil
}

func (g *Graph) populateMap(o *Object) error {
	for i := 0; i < o.reflectValue.Elem().NumField(); i++ {
		field := o.reflectValue.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := o.reflectType.Elem().Field(i).Tag
		fieldName := o.reflectType.Elem().Field(i).Name

		tag := parseTagV2(fieldTag)
		if tag == nil {
			continue
		}

		// Interface injection can't be private because we can't instantiate new
		// instances of an interface.
		if tag.Private {
			return fmt.Errorf(
				"found private inject tag on interface field %s in type %s",
				o.reflectType.Elem().Field(i).Name,
				o.reflectType,
			)
		}

		// Don't overwrite existing values.
		if !isNilOrZero(field, fieldType) {
			continue
		}

		// 第三轮注入，只注入Map/List类型的字段
		if fieldType.Kind() == reflect.Map {
			mp := reflect.MakeMap(fieldType)
			for _, existing := range g.named {
				if existing.private || !existing.reflectType.AssignableTo(fieldType.Elem()) {
					continue
				}
				mp.SetMapIndex(reflect.ValueOf(existing.Name), reflect.ValueOf(existing.Value))
			}
			field.Set(mp)
			o.addDep(fieldName, &Object{Value: mp, Complete: true, created: true})
		}

		if fieldType.Kind() == reflect.Slice {
			founds := reflect.MakeSlice(fieldType, 0, 0)
			for _, existing := range g.named {
				if existing.private || !existing.reflectType.AssignableTo(fieldType.Elem()) {
					continue
				}
				founds = reflect.Append(founds, reflect.ValueOf(existing.Value))
			}

			field.Set(founds)
			o.addDep(fieldName, &Object{Value: founds, Complete: true, created: true})
		}
	}
	return nil
}

// Objects returns all known objects, named as well as unnamed. The returned
// elements are not in a stable order.
func (g *Graph) Objects() []*Object {
	objects := make([]*Object, 0, len(g.named))
	for _, o := range g.named {
		if !o.embedded {
			objects = append(objects, o)
		}
	}
	// randomize to prevent callers from relying on ordering
	for i := 0; i < len(objects); i++ {
		j := rand.Intn(i + 1)
		objects[i], objects[j] = objects[j], objects[i]
	}
	return objects
}

func (g *Graph) Get(name string) any {
	if obj, ok := g.named[name]; ok {
		return obj.Value
	}
	panic("bean not found: " + name)
}

func (g *Graph) Query(t reflect.Type) []*Object {
	results := make([]*Object, 0)
	for _, existing := range g.named {
		if existing.private || !existing.reflectType.AssignableTo(t) {
			continue
		}
		results = append(results, existing)
	}
	return results
}

type tag struct {
	Name    string
	Inline  bool
	Private bool
}

func parseTagV2(t reflect.StructTag) *tag {
	value, ok := t.Lookup("inject")
	if !ok {
		return nil
	}
	if value == "" {
		return &tag{}
	}
	if value == "inline" {
		return &tag{Inline: true}
	}
	if value == "private" {
		return &tag{Private: true}
	}
	return &tag{Name: value}
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func isNilOrZero(v reflect.Value, t reflect.Type) bool {
	switch v.Kind() {
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(t).Interface())
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
}
