package test

type IHello interface {
    SayHello() string
    SetWord(word string)
}

type Hello struct {
    SayHello func() string
    SetWord  func(word string)
}
type HelloWorld struct {
    Word string
}

func (h *HelloWorld) SayHello() string {
    return h.Word
}

func (h *HelloWorld) SetWord(word string) {
    h.Word = word
}

func (h *HelloWorld) GetWord() string {
    return h.Word
}

type HelloWorldV2 struct {
    Word string
}
