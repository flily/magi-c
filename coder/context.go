package coder

type VariableInfo struct {
	Name     string
	Assigned string
}

type VariableMap struct {
	Variables map[string]*VariableInfo
}

func NewVariableMap() *VariableMap {
	m := &VariableMap{
		Variables: make(map[string]*VariableInfo),
	}

	return m
}

func (m *VariableMap) Add(nameInSource string, nameInCode string) {
	m.Variables[nameInSource] = &VariableInfo{
		Name: nameInCode,
	}
}

func (m *VariableMap) Get(name string) (*VariableInfo, bool) {
	info, found := m.Variables[name]
	return info, found
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
