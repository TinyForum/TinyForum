package utils

import (
	"io"

	"golang.org/x/net/html"
)

type Utils interface {
	TimeExpressionParser() TimeExpressionParser
	HTMLParser() HTMLParser
}

// HTMLParser 解析 HTML 节点，输出特定格式的字符串。
type HTMLParser interface {
	Parse(node *html.Node, w io.Writer) error
}
