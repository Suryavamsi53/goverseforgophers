package markdown

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type Renderer struct {
	engine goldmark.Markdown
}

func NewRenderer() *Renderer {
	engine := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			extension.Table,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"), // Dark premium style
				highlighting.WithFormatOptions(
					html.WithClasses(false), // Use inline styles for easy setup
					html.TabWidth(4),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			htmlrenderer.WithHardWraps(),
			htmlrenderer.WithUnsafe(),
		),
	)

	return &Renderer{
		engine: engine,
	}
}

// Render converts markdown string to HTML string
func (r *Renderer) Render(source []byte) (string, error) {
	var buf bytes.Buffer
	if err := r.engine.Convert(source, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
