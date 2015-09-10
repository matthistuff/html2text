package html2text

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"strings"
)

// Formatters is the global formatter list.
var Formatters = map[string]Formatter{
	"_default": defaultFormatter{},
	"_text":    textFormatter{},
	"a":        aFormatter{},
	"br":       brFormatter{},
	"p":        pFormatter{},
	"li":       liFormatter{},
}

func textify(node *html.Node, buf *bytes.Buffer, childIndex int) error {
	var err error

	switch node.Type {
	case html.TextNode:
		if str, err := Formatters["_text"].Format(node, childIndex); err == nil {
			buf.WriteString(str)
		}
	case html.ElementNode:
		var (
			formatter Formatter
			exists    bool
		)

		if formatter, exists = Formatters[node.Data]; !exists {
			formatter = Formatters["_default"]
		}

		if format, err := formatter.Format(node, childIndex); err == nil {
			// Only drill down on child nodes when there is some fmt verb
			if strings.Contains(format, "%") {
				contentBuf := &bytes.Buffer{}
				if err = recurse(node, contentBuf); err == nil {
					nodeStr := fmt.Sprintf(format, contentBuf.String())
					nodeStr = doubleSpaceRe.ReplaceAllString(nodeStr, " ")
					nodeStr = spaceNewLineRe.ReplaceAllString(nodeStr, "\n")
					nodeStr = multiNewLineRe.ReplaceAllString(nodeStr, "\n\n")

					buf.WriteString(nodeStr)
				}
			} else {
				buf.WriteString(format)
			}
		}
	default:
		recurse(node, buf)
	}

	if err != nil {
		return err
	}
	return nil
}

func recurse(node *html.Node, buf *bytes.Buffer) error {
	var err error

	childIndex := 0
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err = textify(c, buf, childIndex); err != nil {
			return err
		}
		if c.Type == html.ElementNode {
			childIndex++
		}
	}

	return nil
}

// IsPreformatted returns true when a node belongs to a preformatted element (<pre>)
func IsPreformatted(node *html.Node) bool {
	for n := node; n != nil; n = n.Parent {
		if n.DataAtom == atom.Pre {
			return true
		}
	}
	return false
}

func FromReader(reader io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	doc, err := html.Parse(reader)
	if err != nil {
		return "", err
	}
	if err = textify(doc, buf, 0); err != nil {
		return "", err
	}
	text := strings.TrimSpace(buf.String())
	return text, nil
}

func FromString(input string) (string, error) {
	text, err := FromReader(strings.NewReader(input))
	if err != nil {
		return "", err
	}
	return text, nil
}
