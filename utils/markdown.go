package utils

import (
	"github.com/Depado/bfchroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/russross/blackfriday/v2"
)

const defaultBfExtensions = blackfriday.NoIntraEmphasis | blackfriday.Tables | blackfriday.FencedCode | blackfriday.Autolink | blackfriday.Strikethrough | blackfriday.SpaceHeadings | blackfriday.BackslashLineBreak | blackfriday.DefinitionLists | blackfriday.Footnotes

var chromaRender blackfriday.Renderer

func defaultChromaRender() blackfriday.Renderer {
	if chromaRender == nil {
		chromaRender = bfchroma.NewRenderer(
			bfchroma.WithoutAutodetect(),
			bfchroma.Style("tango"),
			bfchroma.ChromaOptions(html.WithLineNumbers(true)),
			bfchroma.Extend(
				blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
					Flags: blackfriday.UseXHTML | blackfriday.Smartypants | blackfriday.SmartypantsFractions |
						blackfriday.SmartypantsDashes | blackfriday.SmartypantsLatexDashes | blackfriday.TOC,
				}),
			),
		)
	}

	return chromaRender
}

func Markdown2Html(md []byte) []byte {
	return blackfriday.Run(md,
		blackfriday.WithExtensions(defaultBfExtensions),
		blackfriday.WithRenderer(defaultChromaRender()),
	)
}

func Markdown2HtmlWithOption(md []byte, opts ...blackfriday.Option) []byte {
	return blackfriday.Run(md, opts...)
}
