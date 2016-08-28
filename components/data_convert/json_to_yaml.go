package dataconvert

import (
	"bytes"
	"strings"

	"github.com/simpleelegant/devtools/components/data_convert/json.analyser"
)

func jsonToYAML(input string) (string, error) {
	a := &analyser.LexemeList{}
	if _, err := a.Write([]byte(input)); err != nil {
		return "", err
	}

	t := analyser.SyntaxTree{}
	if err := t.Write(a); err != nil {
		return "", err
	}

	return string(t.Convert(jsonToYAMLTravel)), nil
}

func jsonToYAMLTravel(n interface{}, depth int, output *bytes.Buffer) {
	switch i := n.(type) {
	case *analyser.ArrayNode:
		if len(i.Items) == 0 {
			output.WriteString(" []\n")
		} else {
			output.WriteString("\n")
			for _, v := range i.Items {
				output.WriteString(multipleString("  ", depth) + "-")
				jsonToYAMLTravel(v, depth+1, output)
			}
		}
	case *analyser.ObjectNode:
		if len(i.Properties) == 0 {
			output.WriteString(" {}\n")
		} else {
			output.WriteString("\n")
			for k, v := range i.Properties {
				output.WriteString(multipleString("  ", depth) + yamlStringWrap(k) + ":")
				jsonToYAMLTravel(v, depth+1, output)
			}
		}
	case *analyser.Lexeme:
		output.WriteString(" ")
		defer output.WriteString("\n")

		switch i.Type {
		case analyser.String:
			output.WriteString(yamlStringWrap(i.Value))
		default:
			output.WriteString(i.Value)
		}
	default:
	}
}

func yamlStringWrap(s string) string {
	if s == "" {
		return "\"\""
	}
	if len(strings.TrimSpace(s)) != len(s) {
		return "\"" + s + "\""
	}

	return s
}
