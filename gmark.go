package gmark

type Markdown struct {
	Inline   Lexer
	Block    Lexer
	Renderer Renderer

	Tokens []*Token
}

func (m *Markdown) Parse(text string) string {
	m.Tokens = m.Block.Process(text, nil)
	return m.Output()
}

func (m *Markdown) Output() (out string) {
	for _, tok := range m.Tokens {
		out += m.Renderer.Render(tok)
	}
	return
}

func Convert(text string) string {
	m := Markdown{
		Block:    DefaultBlockLexer,
		Renderer: DefaultRenderer,
	}
	return m.Parse(text)
}
