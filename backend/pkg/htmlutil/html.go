package htmlutil

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type HTMLParser interface {
	Parse(node *html.Node, w io.Writer) error
}

// ============================================================================
// 具体解析器实现（SRP：每个解析器只负责一种格式）
// ============================================================================

// MarkdownParser 将 HTML 转换为 Markdown。
type MarkdownParser struct{}

// Parse 实现 HTMLParser 接口，将 HTML 节点转换为 Markdown 并写入 writer。
func (MarkdownParser) Parse(node *html.Node, w io.Writer) error {
	return parseNodeToMarkdown(node, w)
}

// TextParser 将 HTML 转换为纯文本。
type TextParser struct{}

// Parse 实现 HTMLParser 接口，将 HTML 节点转换为纯文本并写入 writer。
func (TextParser) Parse(node *html.Node, w io.Writer) error {
	return parseNodeToText(node, w)
}

// ============================================================================
// Markdown 转换核心逻辑
// ============================================================================

func parseNodeToMarkdown(node *html.Node, w io.Writer) error {
	switch node.Type {
	case html.TextNode:
		text := strings.TrimSpace(node.Data)
		if text != "" {
			_, err := w.Write([]byte(text))
			if err != nil {
				return err
			}
		}
	case html.ElementNode:
		switch node.Data {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			level := int(node.Data[1] - '0')
			prefix := strings.Repeat("#", level) + " "
			_, err := w.Write([]byte(prefix))
			if err != nil {
				return err
			}
			err = parseChildrenToMarkdown(node, w)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("\n\n"))
			return err

		case "p":
			err := parseChildrenToMarkdown(node, w)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("\n\n"))
			return err

		case "strong", "b":
			_, err := w.Write([]byte("**"))
			if err != nil {
				return err
			}
			err = parseChildrenToMarkdown(node, w)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("**"))
			return err

		case "em", "i":
			_, err := w.Write([]byte("*"))
			if err != nil {
				return err
			}
			err = parseChildrenToMarkdown(node, w)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("*"))
			return err

		case "a":
			href := getAttr(node, "href")
			_, err := w.Write([]byte("["))
			if err != nil {
				return err
			}
			err = parseChildrenToMarkdown(node, w)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("](" + href + ")"))
			return err

		case "img":
			src := getAttr(node, "src")
			alt := getAttr(node, "alt")
			_, err := w.Write([]byte(fmt.Sprintf("![%s](%s)", alt, src)))
			return err

		case "ul", "ol":
			err := parseListToMarkdown(node, w, node.Data == "ol")
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("\n"))
			return err

		case "li":
			// 由父级 ul/ol 处理，这里不单独输出
			return parseChildrenToMarkdown(node, w)

		case "blockquote":
			_, err := w.Write([]byte("> "))
			if err != nil {
				return err
			}
			err = parseChildrenToMarkdown(node, w)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("\n\n"))
			return err

		case "br":
			_, err := w.Write([]byte("\n"))
			return err

		default:
			// 其他元素只处理子节点
			return parseChildrenToMarkdown(node, w)
		}
	}
	return nil
}

func parseChildrenToMarkdown(node *html.Node, w io.Writer) error {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if err := parseNodeToMarkdown(child, w); err != nil {
			return err
		}
	}
	return nil
}

func parseListToMarkdown(node *html.Node, w io.Writer, ordered bool) error {
	idx := 1
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "li" {
			if ordered {
				_, err := fmt.Fprintf(w, "%d. ", idx)
				if err != nil {
					return err
				}
				idx++
			} else {
				_, err := w.Write([]byte("- "))
				if err != nil {
					return err
				}
			}
			if err := parseChildrenToMarkdown(child, w); err != nil {
				return err
			}
			_, err := w.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ============================================================================
// 纯文本转换核心逻辑
// ============================================================================

func parseNodeToText(node *html.Node, w io.Writer) error {
	switch node.Type {
	case html.TextNode:
		text := strings.TrimSpace(node.Data)
		if text != "" {
			_, err := w.Write([]byte(text))
			if err != nil {
				return err
			}
		}
	case html.ElementNode:
		switch node.Data {
		case "br", "p", "div", "h1", "h2", "h3", "h4", "h5", "h6", "li":
			err := parseChildrenToText(node, w)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte("\n"))
			return err
		default:
			return parseChildrenToText(node, w)
		}
	}
	return nil
}

func parseChildrenToText(node *html.Node, w io.Writer) error {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if err := parseNodeToText(child, w); err != nil {
			return err
		}
	}
	return nil
}

// ============================================================================
// 辅助函数
// ============================================================================

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// ============================================================================
// 组合解析器（职责链，OCP：可无限扩展）
// ============================================================================

// ParserHtmlChain 实现 HTMLParser 接口，按顺序尝试多个解析器。
type ParserHtmlChain struct {
	parsers []HTMLParser
}

// NewParserHtmlChain 创建一个解析器链。
func NewHtmlParserHtmlChain(parsers ...HTMLParser) *ParserHtmlChain {
	return &ParserHtmlChain{parsers: parsers}
}

// Add 在链末尾添加解析器（返回自身，支持链式调用）。
func (c *ParserHtmlChain) Add(p HTMLParser) *ParserHtmlChain {
	c.parsers = append(c.parsers, p)
	return c
}

// Parse 实现 HTMLParser 接口，依次调用链中解析器，第一个成功即返回。
func (c *ParserHtmlChain) Parse(node *html.Node, w io.Writer) error {
	for _, p := range c.parsers {
		if err := p.Parse(node, w); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no parser could handle the HTML node")
}

// ============================================================================
// 带错误返回的封装器（方便业务层使用）
// ============================================================================

// HTMLConverter 封装一个 HTMLParser，提供返回 error 的 Parse 方法。
type HTMLConverter struct {
	parser HTMLParser
}

// NewHTMLConverter 创建带错误返回的转换器。
func NewHTMLConverter(parser HTMLParser) *HTMLConverter {
	return &HTMLConverter{parser: parser}
}

// Convert 解析 HTML 字符串，输出指定格式的字符串。
func (c *HTMLConverter) Convert(htmlStr string) (string, error) {
	node, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}
	buf := &bytes.Buffer{}
	if err := c.parser.Parse(node, buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ============================================================================
// 包级默认实例（向后兼容 + 开闭原则：可替换）
// ============================================================================

var (
	// defaultChain 默认解析链（目前只有一个，但可以扩展）
	defaultChain = NewHtmlParserHtmlChain(MarkdownParser{})

	// DefaultMDConverter 默认的 Markdown 转换器
	DefaultMDConverter = NewHTMLConverter(MarkdownParser{})

	// DefaultTextConverter 默认的纯文本转换器
	DefaultTextConverter = NewHTMLConverter(TextParser{})
)

// HTMLToMarkdown 将 HTML 字符串转换为 Markdown（便捷函数）。
func HTMLToMarkdown(htmlStr string) (string, error) {
	return DefaultMDConverter.Convert(htmlStr)
}

// HTMLToText 将 HTML 字符串转换为纯文本（便捷函数）。
func HTMLToText(htmlStr string) (string, error) {
	return DefaultTextConverter.Convert(htmlStr)
}
