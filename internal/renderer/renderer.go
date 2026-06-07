package renderer

import (
	"github.com/charmbracelet/glamour"
)

type Renderer struct {
	gr *glamour.TermRenderer
}

func New(theme string, width int) (*Renderer, error) {
	styleName := "dark"
	switch theme {
	case "light":
		styleName = "light"
	case "ascii":
		styleName = "ascii"
	case "notty":
		styleName = "notty"
	}

	gr, err := glamour.NewTermRenderer(
		glamour.WithStylePath(styleName),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}
	return &Renderer{gr: gr}, nil
}

func (r *Renderer) Render(markdown string) (string, error) {
	return r.gr.Render(markdown)
}
