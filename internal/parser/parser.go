package parser

import (
	"strings"

	"github.com/user/clislides/internal/model"
	"gopkg.in/yaml.v3"
)

func Parse(content string) (*model.Presentation, error) {
	pres := &model.Presentation{
		Meta: model.Metadata{Pager: true},
	}

	body := content
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		parts := strings.SplitN(strings.TrimSpace(content), "\n", 2)
		if len(parts) > 1 {
			rest := parts[1]
			end := strings.Index(rest, "\n---")
			if end != -1 {
				frontmatter := rest[:end]
				body = strings.TrimSpace(rest[end+4:])
				_ = yaml.Unmarshal([]byte(frontmatter), &pres.Meta)
			}
		}
	}

	pages := splitSlides(body)
	for i, page := range pages {
		pres.Slides = append(pres.Slides, model.Slide{
			Raw:   strings.TrimSpace(page),
			Index: i,
		})
	}

	if len(pres.Slides) == 0 {
		pres.Slides = append(pres.Slides, model.Slide{Raw: "", Index: 0})
	}

	return pres, nil
}

func splitSlides(body string) []string {
	lines := strings.Split(body, "\n")
	var slides []string
	var current []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			slides = append(slides, strings.Join(current, "\n"))
			current = nil
		} else {
			current = append(current, line)
		}
	}
	if len(current) > 0 || len(slides) > 0 {
		slides = append(slides, strings.Join(current, "\n"))
	}
	return slides
}
