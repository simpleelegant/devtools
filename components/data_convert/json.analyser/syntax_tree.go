package analyser

import (
	"bytes"
	"encoding/json"
	"errors"
)

// ErrRootNodeShouldBeObjectOrArray ...
var ErrRootNodeShouldBeObjectOrArray = errors.New("Root node should be object or array")

// Node represent node on syntax tree
type Node interface {
	Append(i interface{})
	SetParent(p Node)
	GetParent() Node
}

// ObjectNode represent node of JSON object
type ObjectNode struct {
	Parent     Node   `json:"-"`
	NextKey    string `json:"-"`
	Properties map[string]interface{}
}

// SetParent set parent of node
func (n *ObjectNode) SetParent(p Node) {
	n.Parent = p
}

// GetParent get parent of node
func (n *ObjectNode) GetParent() Node {
	return n.Parent
}

// Append append sub node or other things
func (n *ObjectNode) Append(i interface{}) {
	if n.NextKey == "" {
		n.NextKey = i.(*Lexeme).Value
		return
	}

	if n.Properties == nil {
		n.Properties = map[string]interface{}{}
	}

	n.Properties[n.NextKey] = i
	n.NextKey = ""

	if o, ok := i.(Node); ok {
		o.SetParent(n)
	}
}

// ArrayNode represent node of JSON array
type ArrayNode struct {
	Parent Node `json:"-"`
	Items  []interface{}
}

// SetParent set parent of node
func (n *ArrayNode) SetParent(p Node) {
	n.Parent = p
}

// GetParent get parent of node
func (n *ArrayNode) GetParent() Node {
	return n.Parent
}

// Append append sub node or other things
func (n *ArrayNode) Append(i interface{}) {
	n.Items = append(n.Items, i)

	if o, ok := i.(Node); ok {
		o.SetParent(n)
	}
}

// SyntaxTree syntax tree
type SyntaxTree struct {
	Root    Node
	Current Node
	Remain  string
}

// Write ...
func (t *SyntaxTree) Write(l *LexemeList) error {
	lexeme, err := l.Read()
	// if EOF
	if err != nil {
		return nil
	}

	// init root node
	switch lexeme.Type {
	case ObjectOpen:
		t.Root = &ObjectNode{}
		t.Current = t.Root
	case ArrayOpen:
		t.Root = &ArrayNode{}
		t.Current = t.Root
	default:
		return ErrRootNodeShouldBeObjectOrArray
	}

	for true {
		lexeme, err := l.Read()
		if err != nil {
			return nil
		}

		switch lexeme.Type {
		case ObjectOpen:
			n := &ObjectNode{}
			t.Current.Append(n)
			t.Current = n
		case ArrayOpen:
			n := &ArrayNode{}
			t.Current.Append(n)
			t.Current = n
		case ObjectClose, ArrayClose:
			t.Current = t.Current.GetParent()
		case String, Number, Bool, Null:
			t.Current.Append(lexeme)
		default:
		}
	}

	return nil
}

// EncodeToStruct encode to Go struct
func (t *SyntaxTree) Convert(travel func(node interface{}, depth int, output *bytes.Buffer)) []byte {
	b := &bytes.Buffer{}

	travel(t.Root, 0, b)

	return b.Bytes()
}

// String ...
func (t *SyntaxTree) String() string {
	b, err := json.Marshal(t.Root)
	if err != nil {
		return err.Error()
	}

	return string(b)
}
