package analyser

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// ErrEOF represent no more lexeme
var ErrEOF = errors.New("No more lexeme")

type lexemeType int

// all possible JSON lexeme
const (
	String lexemeType = iota
	Number
	Bool
	Null

	ObjectOpen
	ObjectClose
	ArrayOpen
	ArrayClose

	Comma
	Colon
)

// Lexeme ...
type Lexeme struct {
	Type  lexemeType
	Value string
}

// LexemeList ...
type LexemeList struct {
	list   []*Lexeme
	seek   int
	remain []byte
}

// Read ...
func (l *LexemeList) Read() (*Lexeme, error) {
	if l.seek >= len(l.list) {
		return nil, ErrEOF
	}

	l.seek++

	return l.list[l.seek-1], nil
}

func (l *LexemeList) parseRemain() {
	if len(l.remain) != 0 {
		e := &Lexeme{}
		e.Type, e.Value = parseValue(string(l.remain))
		l.list = append(l.list, e)
		l.remain = []byte{}
	}
}

// Write accept JSON string
func (l *LexemeList) Write(p []byte) (n int, err error) {
	// check json string is valid
	var scrap bytes.Buffer
	if err := json.Compact(&scrap, p); err != nil {
		return 0, err
	}

	var (
		doubleQuote = []byte{'"'}
		lexeme      *Lexeme
	)

	for _, b := range p {
		n++

		// the byte may be part of a string
		if bytes.HasPrefix(l.remain, doubleQuote) &&
			(len(l.remain) == 1 || !bytes.HasSuffix(l.remain, doubleQuote)) {
			l.remain = append(l.remain, b)
			continue
		}

		if isSpace(b) {
			continue
		}

		switch b {
		case '{':
			lexeme = &Lexeme{Type: ObjectOpen}
		case '}':
			lexeme = &Lexeme{Type: ObjectClose}
		case '[':
			lexeme = &Lexeme{Type: ArrayOpen}
		case ']':
			lexeme = &Lexeme{Type: ArrayClose}
		case ':':
			lexeme = &Lexeme{Type: Colon}
		case ',':
			lexeme = &Lexeme{Type: Comma}
		default:
			l.remain = append(l.remain, b)
			continue
		}

		l.parseRemain()

		l.list = append(l.list, lexeme)
	}

	l.parseRemain()

	return
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func parseValue(p string) (t lexemeType, v string) {
	switch p {
	case "null":
		return Null, p
	case "true", "false":
		return Bool, p
	default:
		// string
		if strings.HasPrefix(p, "\"") {
			return String, strings.Trim(p, "\"")
		}

		// number
		if _, err := strconv.ParseFloat(p, 64); err == nil {
			return Number, p
		}

		return String, p
	}
}
