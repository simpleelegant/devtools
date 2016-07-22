package analyser

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
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
func (t *SyntaxTree) EncodeToStruct() []byte {
	w := &bytes.Buffer{}

	travel(t.Root, 0, w)

	return w.Bytes()
}

// String ...
func (t *SyntaxTree) String() string {
	b, err := json.Marshal(t.Root)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

func travel(n interface{}, depth int, w io.Writer) {
	switch i := n.(type) {
	case *ArrayNode:
		if len(i.Items) == 0 {
			w.Write([]byte("[]interface{}"))
		} else {
			w.Write([]byte("[]"))
			travel(i.Items[0], depth, w)
		}
	case *ObjectNode:
		w.Write([]byte("struct {\n"))

		spin := multipleString("    ", depth+1)
		for k, v := range i.Properties {
			// TODO check k valid

			w.Write([]byte(spin + insureLen(toCamelCase(k), 12) + " "))
			travel(v, depth+1, w)
			w.Write([]byte(" `json:\"" + k + "\"`\n"))
		}

		w.Write([]byte(multipleString("    ", depth) + "}"))
	case *Lexeme:
		switch i.Type {
		case String:
			w.Write([]byte("string "))
		case Number:
			w.Write([]byte("float64"))
		case Bool:
			w.Write([]byte("bool   "))
		case Null:
			w.Write([]byte("interface{}"))
		default:
		}
	default:
	}
}

func toCamelCase(s string) string {
	return strings.Replace(strings.Title(strings.Replace(s, "_", " ", -1)), " ", "", -1)
}

func multipleString(s string, multiple int) string {
	if multiple <= 0 {
		return ""
	}

	b := ""
	for i := 0; i < multiple; i++ {
		b += s
	}

	return b
}

func insureLen(s string, min int) string {
	return s + multipleString(" ", min-len(s))
}
