package dataconvert

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func jsonIndent(input string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(input), "", "    ")

	return out.String(), err
}

func jsonCompact(input string) (string, error) {
	var out bytes.Buffer
	err := json.Compact(&out, []byte(input))

	return out.String(), err
}

func base64URLEncode(input string) (string, error) {
	encoded := base64.URLEncoding.EncodeToString([]byte(input))

	return encoded, nil
}

func base64URLDecode(input string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(input)

	return string(decoded), err
}

func md5Checksum(input string) (string, error) {
	return fmt.Sprintf("%x", md5.Sum([]byte(input))), nil
}

func keyValueToJSON(input string) (string, error) {
	ctn := make(map[string]string, 10)

	parseKeyValueLines(input, func(k, v string) {
		ctn[k] = v
	})

	b, err := json.Marshal(ctn)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func keyValueToQueryString(input string) (string, error) {
	q := url.Values{}

	parseKeyValueLines(input, func(k, v string) {
		q.Add(k, v)
	})

	return q.Encode(), nil
}

func queryStringToKeyValue(input string) (string, error) {
	var result string

	v, err := url.ParseQuery(input)
	if err != nil {
		return "", err
	}

	for k, vList := range v {
		for _, v := range vList {
			result += fmt.Sprintf("%s=%s\n", k, v)
		}
	}

	return result, nil
}

func escapeNewline(input string) (string, error) {
	input = strings.Replace(input, "\r\n", "\n", -1)
	return strings.Replace(input, "\n", "\\n", -1), nil
}

func captureNewline(input string) (string, error) {
	return strings.Replace(input, "\\n", "\n", -1), nil
}
