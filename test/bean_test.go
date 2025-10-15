package test

import (
    "Thor/utils/inject"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestBeanAutoProvide(t *testing.T) {
    type Zoo struct {
        Cat *Cat `bean:"cat"`
    }

    g := inject.NewGraph()
    g.CycleProvide(&inject.Object{Value: &Zoo{}})
    assert.Equal(t, "cat:", g.GetByName("cat").(*Cat).GetName())
}

func TestBeanAutoPopulate(t *testing.T) {
    type Zoo struct {
        Cat *Cat `bean:"cat" inject:"cat"`
    }

    g := inject.NewGraph()
    g.CycleProvide(&inject.Object{Name: "zoo", Value: &Zoo{}})
    _ = g.Populate()

    zoo := g.GetByName("zoo").(*Zoo)
    assert.Equal(t, "cat:", zoo.Cat.GetName())
}
