package converter

import (
	"bytes"
	"fmt"

	htmltomd "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/yuin/goldmark"
)

// HTMLToMarkdown converts HTML bytes to Markdown bytes.
func HTMLToMarkdown(htmlContent []byte) ([]byte, error) {
	converter := htmltomd.NewConverter("", true, nil)
	md, err := converter.ConvertBytes(htmlContent)
	if err != nil {
		return nil, fmt.Errorf("HTMLToMarkdown: %w", err)
	}
	return md, nil
}

// MarkdownToHTML converts Markdown bytes to HTML bytes.
func MarkdownToHTML(mdContent []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(mdContent, &buf); err != nil {
		return nil, fmt.Errorf("MarkdownToHTML: %w", err)
	}
	return buf.Bytes(), nil
}
