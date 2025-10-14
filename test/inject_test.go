package test

import (
    "Thor/utils/injectv2"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGetBeanByName(t *testing.T) {
    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{}})
    beans.Provide(&injectv2.Object{Name: "dog", Value: &Dog{}})
    _ = beans.Populate()

    assert.Equal(t, "cat:", beans.GetByName("cat").(*Cat).GetName())
    assert.Equal(t, "dog:", beans.GetByName("dog").(*Dog).GetName())
}

func TestInjectByNamed(t *testing.T) {
    type Zoo struct {
        Cat *Cat `inject:"cat"`
    }

    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{Name: "kitty"}})
    beans.Provide(&injectv2.Object{Name: "zoo", Value: &Zoo{}})
    _ = beans.Populate()

    zoo := beans.GetByName("zoo").(*Zoo)
    assert.Equal(t, "cat:kitty", zoo.Cat.GetName())
}

func TestInjectByTyped(t *testing.T) {
    type Zoo struct {
        Cat *Cat `inject:""`
    }

    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{Name: "kitty"}})
    beans.Provide(&injectv2.Object{Name: "zoo", Value: &Zoo{}})
    _ = beans.Populate()

    zoo := beans.GetByName("zoo").(*Zoo)
    assert.Equal(t, "cat:kitty", zoo.Cat.GetName())
}

func TestInjectByUnnamedInterface(t *testing.T) {
    type Zoo struct {
        Animal IAnimal `inject:""`
    }

    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{Name: "kitty"}})
    beans.Provide(&injectv2.Object{Name: "zoo", Value: &Zoo{}})
    _ = beans.Populate()

    zoo := beans.GetByName("zoo").(*Zoo)
    assert.Equal(t, "cat:kitty", zoo.Animal.GetName())
}

func TestInjectByNamedInterface(t *testing.T) {
    type Zoo struct {
        Animal IAnimal `inject:"cat"`
    }

    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{Name: "kitty"}})
    beans.Provide(&injectv2.Object{Name: "zoo", Value: &Zoo{}})
    _ = beans.Populate()

    zoo := beans.GetByName("zoo").(*Zoo)
    assert.Equal(t, "cat:kitty", zoo.Animal.GetName())
}

func TestInjectByList(t *testing.T) {
    type Zoo struct {
        Animals []IAnimal `inject:""`
    }

    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{Name: "kitty"}})
    beans.Provide(&injectv2.Object{Name: "dog", Value: &Dog{Name: "fido"}})
    beans.Provide(&injectv2.Object{Name: "zoo", Value: &Zoo{}})
    _ = beans.Populate()

    zoo := beans.GetByName("zoo").(*Zoo)
    assert.Equal(t, "cat:kitty", zoo.Animals[0].GetName())
    assert.Equal(t, "dog:fido", zoo.Animals[1].GetName())
}

func TestInjectByMap(t *testing.T) {
    type Zoo struct {
        Animals map[string]IAnimal `inject:""`
    }

    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{Name: "kitty"}})
    beans.Provide(&injectv2.Object{Name: "dog", Value: &Dog{Name: "fido"}})
    beans.Provide(&injectv2.Object{Name: "zoo", Value: &Zoo{}})
    _ = beans.Populate()

    zoo := beans.GetByName("zoo").(*Zoo)
    assert.Equal(t, "cat:kitty", zoo.Animals["cat"].GetName())
    assert.Equal(t, "dog:fido", zoo.Animals["dog"].GetName())
}

func TestCycleInject(t *testing.T) {
    type Zoo struct {
        Animals map[string]IAnimal `inject:""`
    }

    type City struct {
        Zoo *Zoo `inject:""`
    }

    beans := injectv2.NewGraph()
    beans.Provide(&injectv2.Object{Name: "city", Value: &City{}})
    beans.Provide(&injectv2.Object{Name: "cat", Value: &Cat{Name: "kitty"}})
    beans.Provide(&injectv2.Object{Name: "dog", Value: &Dog{Name: "fido"}})
    beans.Provide(&injectv2.Object{Name: "zoo", Value: &Zoo{}})
    _ = beans.Populate()

    city := beans.GetByName("city").(*City)
    assert.Equal(t, "cat:kitty", city.Zoo.Animals["cat"].GetName())
    assert.Equal(t, "dog:fido", city.Zoo.Animals["dog"].GetName())
}
