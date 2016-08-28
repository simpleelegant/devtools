package dataconvert

import (
	"bytes"

	"github.com/simpleelegant/devtools/components/data_convert/json.analyser"
)

func jsonToGoStruct(input string) (string, error) {
	a := &analyser.LexemeList{}
	if _, err := a.Write([]byte(input)); err != nil {
		return "", err
	}

	t := analyser.SyntaxTree{}
	if err := t.Write(a); err != nil {
		return "", err
	}

	return string(t.Convert(jsonToGoStructTravel)), nil
}

func jsonToGoStructTravel(node interface{}, depth int, output *bytes.Buffer) {
	switch i := node.(type) {
	case *analyser.ArrayNode:
		if len(i.Items) == 0 {
			output.WriteString("[]interface{}")
		} else {
			output.WriteString("[]")
			jsonToGoStructTravel(i.Items[0], depth, output)
		}
	case *analyser.ObjectNode:
		output.WriteString("struct {\n")

		spin := multipleString("    ", depth+1)
		for k, v := range i.Properties {
			output.WriteString(spin + insureLen(toGoFieldName(k), 12) + " ")
			jsonToGoStructTravel(v, depth+1, output)
			output.WriteString(" `json:\"" + k + "\"`\n")
		}

		output.WriteString(multipleString("    ", depth) + "}")
	case *analyser.Lexeme:
		switch i.Type {
		case analyser.String:
			output.WriteString("string ")
		case analyser.Number:
			output.WriteString("float64")
		case analyser.Bool:
			output.WriteString("bool   ")
		case analyser.Null:
			output.WriteString("interface{}")
		default:
		}
	default:
	}
}
