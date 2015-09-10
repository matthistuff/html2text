package html2text

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strconv"
)

// Formatter converts an html.Node to a string.
type Formatter interface {
	Format(*html.Node, int) (string, error)
}

// defaultFormatter does nothing
type defaultFormatter struct{}

func (defaultFormatter) Format(node *html.Node, childIndex int) (string, error) {
	return "%s", nil
}

// aFormatter formats <a> tags
type aFormatter struct{}

func (aFormatter) Format(node *html.Node, childIndex int) (string, error) {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			quoted := strconv.Quote(attr.Val)
			return quoted[1 : len(quoted)-1], nil
		}
	}
	return "%s", nil
}

// liFormatter formats a <li> tag
type liFormatter struct{}

func (liFormatter) Format(node *html.Node, childIndex int) (string, error) {
	if node.Parent.DataAtom == atom.Ol {
		return fmt.Sprintf("%d. %%s\n", childIndex+1), nil
	}
	return "* %s\n", nil
}

// brFormatter formats a <br> tag
type brFormatter struct{}

func (brFormatter) Format(node *html.Node, childIndex int) (string, error) {
	return "\n", nil
}

// pFormatter formats a <p> tag
type pFormatter struct{}

func (pFormatter) Format(node *html.Node, childIndex int) (string, error) {
	return "\n\n%s\n\n", nil
}

// textFormatter formats text
type textFormatter struct{}

func (textFormatter) Format(node *html.Node, childIndex int) (string, error) {
	if IsPreformatted(node) {
		return node.Data, nil
	}

	data := doubleSpaceRe.ReplaceAllString(
		spacingRe.ReplaceAllString(node.Data, " "),
		" ")

	return data, nil
}
