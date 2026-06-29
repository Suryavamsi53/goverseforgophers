package markdown

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	htmlrenderer "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Heading struct {
	Level int
	Text  string
	ID    string
}

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

// ExtractTOC parses the markdown and returns a list of headings
func (r *Renderer) ExtractTOC(source []byte) ([]Heading, error) {
	ctx := parser.NewContext()
	doc := r.engine.Parser().Parse(text.NewReader(source), parser.WithContext(ctx))
	
	var headings []Heading

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading {
			heading := n.(*ast.Heading)
			
			var textBuf []byte
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				textBuf = append(textBuf, c.Text(source)...)
			}
			
			idAttr, ok := n.AttributeString("id")
			var id string
			if ok {
				if idBytes, ok := idAttr.([]byte); ok {
					id = string(idBytes)
				}
			}
			
			headings = append(headings, Heading{
				Level: heading.Level,
				Text:  string(textBuf),
				ID:    id,
			})
		}
		return ast.WalkContinue, nil
	})
	
	return headings, err
}
