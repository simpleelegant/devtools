package tree

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	ErrRootNodeMustBeObjectOrArray = errors.New("Root node must be object or array")
	ErrJSONInvalid                 = errors.New("JSON invalid")
)

// Node represent node on syntax tree
type Node interface {
	Append(i interface{}) error
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
func (n *ObjectNode) Append(i interface{}) error {
	if n.NextKey == "" {
		v, ok := i.(string)
		if !ok {
			return ErrJSONInvalid
		}

		// FIXME need full check
		if strings.ContainsAny(v, "\"") {
			return ErrJSONInvalid
		}

		n.NextKey = v
		return nil
	}

	if n.Properties == nil {
		n.Properties = map[string]interface{}{}
	}
	n.Properties[n.NextKey] = i
	n.NextKey = ""

	if o, ok := i.(Node); ok {
		o.SetParent(n)
	}

	return nil
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
func (n *ArrayNode) Append(i interface{}) error {
	n.Items = append(n.Items, i)

	if o, ok := i.(Node); ok {
		o.SetParent(n)
	}

	return nil
}

// SyntaxTree syntax tree
type SyntaxTree struct {
	Root    Node
	Current Node
	Remain  string
}

// Write ...
func (t *SyntaxTree) Write(p []byte) (n int, err error) {
	for _, b := range p {
		n += 1

		if isSpace(b) {
			continue
		}

		if strings.HasPrefix(t.Remain, "\"") && !strings.HasSuffix(t.Remain, "\"") {
			t.appendRemain(b)

			continue
		}

		switch b {
		case '{':
			n := &ObjectNode{}
			if t.Root == nil {
				t.Root = n
			} else {
				t.Current.Append(n)
			}
			t.Current = n
		case '[':
			n := &ArrayNode{}
			if t.Root == nil {
				t.Root = n
			} else {
				t.Current.Append(n)
			}
			t.Current = n
		case '}', ']':
			if t.Current == nil {
				return n, ErrRootNodeMustBeObjectOrArray
			}
			if err := t.consumeRemain(); err != nil {
				return n, err
			}

			t.Current = t.Current.GetParent()
		case ':', ',':
			if t.Current == nil {
				return n, ErrRootNodeMustBeObjectOrArray
			}
			if err := t.consumeRemain(); err != nil {
				return n, err
			}
		default:
			if t.Current == nil {
				return n, ErrRootNodeMustBeObjectOrArray
			}
			t.appendRemain(b)
		}
	}

	return
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

func (t *SyntaxTree) parseRemain() interface{} {
	switch t.Remain {
	case "null":
		return nil
	case "true":
		return true
	case "false":
		return false
	default:
		// string
		if strings.HasPrefix(t.Remain, "\"") {
			return strings.Trim(t.Remain, "\"")
		}
		if f, err := strconv.ParseFloat(t.Remain, 64); err == nil {
			return f
		}

		return t.Remain
	}
}

func (t *SyntaxTree) appendRemain(c byte) {
	t.Remain += string(c)
}

func (t *SyntaxTree) consumeRemain() error {
	if t.Remain == "" {
		return nil
	}

	if err := t.Current.Append(t.parseRemain()); err != nil {
		return err
	}

	t.Remain = ""

	return nil
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
	case nil:
		w.Write([]byte("interface{}"))
	case bool:
		w.Write([]byte("bool   "))
	case string:
		w.Write([]byte("string "))
	case float64:
		w.Write([]byte("float64"))
	default:
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
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
