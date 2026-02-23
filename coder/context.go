package coder

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/coder/csyntax"
)

type VariableInfo struct {
	SourceName string
	SourceType ast.Type
	CodeName   string
	CodeType   csyntax.Type
	Assigned   string
}

type VariableMap struct {
	Variables []*VariableInfo
}

func NewVariableMap() *VariableMap {
	m := &VariableMap{
		Variables: make([]*VariableInfo, 0, 16),
	}

	return m
}

func (m *VariableMap) Add(nameInSource string, nameInCode string) {
	info := &VariableInfo{
		SourceName: nameInSource,
		CodeName:   nameInCode,
	}
	m.Variables = append(m.Variables, info)
}

func (m *VariableMap) Get(name string) (*VariableInfo, bool) {
	for _, info := range m.Variables {
		if info.SourceName == name {
			return info, true
		}
	}

	return nil, false
}

type Frame struct {
	Variables *VariableMap
	Next      *Frame
}

func NewFrameOn(next *Frame) *Frame {
	frame := &Frame{
		Variables: NewVariableMap(),
		Next:      next,
	}

	return frame
}

func NewFrame() *Frame {
	return NewFrameOn(nil)
}

func (f *Frame) IsRoot() bool {
	return f.Next == nil
}

func (f *Frame) GetName(name string) (*VariableInfo, bool) {
	return f.Variables.Get(name)
}

func (f *Frame) AddName(nameInSource string, nameInCode string) bool {
	if _, found := f.Variables.Get(nameInSource); found {
		return false
	}

	f.Variables.Add(nameInSource, nameInCode)
	return true
}

type Context struct {
	Global        *Frame
	FunctionIn    *VariableMap
	FunctionOut   *VariableMap
	FunctionFrame *Frame
}

func NewContext() *Context {
	ctx := &Context{
		Global:      NewFrame(),
		FunctionIn:  NewVariableMap(),
		FunctionOut: NewVariableMap(),
	}

	return ctx
}

func (c *Context) IsGlobalContext() bool {
	return c.FunctionFrame == nil
}

func (c *Context) Find(name string) (*VariableInfo, bool) {
	frame := c.FunctionFrame
	for frame != nil {
		if info, found := frame.GetName(name); found {
			return info, true
		}

		frame = frame.Next
	}

	if info, found := c.Global.GetName(name); found {
		return info, true
	}

	return nil, false
}

func (c *Context) RegisterVariable(nameInSource string, nameInCode string) bool {
	if c.IsGlobalContext() {
		if _, found := c.Global.GetName(nameInSource); found {
			return false
		}

		return c.Global.AddName(nameInSource, nameInCode)
	}

	top := c.FunctionFrame
	return top.AddName(nameInSource, nameInCode)
}

func (c *Context) PushFrame() *Frame {
	top := c.FunctionFrame
	frame := NewFrameOn(top)
	c.FunctionFrame = frame
	return frame
}

func (c *Context) PopFrame() {
	if c.FunctionFrame != nil {
		c.FunctionFrame = c.FunctionFrame.Next
	}
}
