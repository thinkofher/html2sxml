package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "html2sxml: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading from stdin: %w", err)
	}

	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return fmt.Errorf("html parse: %w", err)
	}

	var root *html.Node
	root = findFirstElement(doc)
	if root == nil {
		return fmt.Errorf("no html element")
	}

	sxml := nodeToSxml(root)

	fmt.Println(sxml)

	return nil
}

func findFirstElement(n *html.Node) *html.Node {
	if n.Type == html.ElementNode {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findFirstElement(c); found != nil {
			return found
		}
	}

	return nil
}

func nodeToSxml(n *html.Node) string {
	switch n.Type {
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text == "" {
			return ""
		}

		return fmt.Sprintf(`"%s"`, escapeQuotes(text))
	case html.ElementNode:
		tagName := ":" + n.Data

		var attrParts []string
		for _, attr := range n.Attr {
			attrParts = append(attrParts, fmt.Sprintf("(%s \"%s\")", attr.Key, escapeQuotes(attr.Val)))
		}

		var attrSection string
		if len(attrParts) > 0 {
			attrSection = fmt.Sprintf("(@ %s)", strings.Join(attrParts, " "))
		}

		var childParts []string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			childSxml := nodeToSxml(c)
			if childSxml != "" {
				childParts = append(childParts, childSxml)
			}
		}

		parts := []string{tagName}
		if attrSection != "" {
			parts = append(parts, attrSection)
		}
		parts = append(parts, childParts...)

		return fmt.Sprintf("(%s)", strings.Join(parts, " "))
	default:
		return ""
	}
}

func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
