package inshape

import (
	"fmt"

	"github.com/google/uuid"
)

type Client struct {
	Tocken    string
	Manager   *ShapeMan
	Instances []*Shape
}
type ShapeResponce struct {
	Methods []string
	Name    string
}

func (c *Client) RequestShape(Name string) ShapeResponce {
	fmt.Println(Name)
	fmt.Println(c.Manager.Shapes)
	m := c.Manager.Shapes[Name].Methods
	Methods := make([]string, len(m))
	i := 0
	for k := range m {
		Methods[i] = k
		i++
	}
	return ShapeResponce{Methods, Name}
}

type Shape struct {
	Methods    map[string]Rat
	InstanceID string
	listeners  map[string]*Client
	Init       func(args any)
}

type Rat func(Return func(res, err any), call Call)

type ShapeMan struct {
	Clients map[string]*Client
	Shapes  map[string]*Shape
}

func (sm *ShapeMan) NewClient() *Client {
	nt := uuid.NewString()
	fmt.Println(nt)
	newClient := &Client{Tocken: nt, Manager: sm}
	sm.Clients[nt] = newClient
	return newClient
}

func (sm *ShapeMan) NewShape(Name string, Methods map[string]Rat) {
	sm.Shapes[Name] = &Shape{Methods: Methods, InstanceID: uuid.NewString(), Init: nil}
}
func (shape *Shape) New() {

}

type Call struct {
	Type uint8
	Args any
}

func NewShapeMan() *ShapeMan {
	return &ShapeMan{
		Clients: make(map[string]*Client),
		Shapes:  make(map[string]*Shape),
	}
}
