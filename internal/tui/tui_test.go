package tui

import (
	"strings"
	"testing"

	"github.com/user/clislides/internal/model"
)

func TestPagerFalseHidesPageNumber(t *testing.T) {
	pres := &model.Presentation{
		Meta: model.Metadata{
			Title:  "Test",
			Author: "Alice",
			Date:   "2026-01-01",
			Pager:  false,
		},
		Slides: []model.Slide{
			{Raw: "# Slide 1", Index: 0},
			{Raw: "# Slide 2", Index: 1},
		},
	}
	m := NewModel(pres)
	m.width = 80
	m.height = 24

	view := m.View()
	if strings.Contains(view, "1/2") {
		t.Error("pager:false should not show page numbers, but found '1/2' in output")
	}
	if !strings.Contains(view, "Alice") {
		t.Error("author should still be visible when pager:false")
	}
	if !strings.Contains(view, "2026-01-01") {
		t.Error("date should still be visible when pager:false")
	}
}

func TestPagerFalseSearchNoPageNumber(t *testing.T) {
	pres := &model.Presentation{
		Meta: model.Metadata{
			Pager: false,
		},
		Slides: []model.Slide{
			{Raw: "# Hello", Index: 0},
			{Raw: "# World", Index: 1},
		},
	}
	m := NewModel(pres)
	m.width = 80
	m.height = 24
	m.mode = modeSearch
	m.searchBuf = "world"

	view := m.View()
	if strings.Contains(view, "1/2") {
		t.Error("pager:false in search mode should not show page numbers")
	}
	if !strings.Contains(view, "/world") {
		t.Error("search input should be visible in status bar")
	}
}

func TestPagerFalseSearchResultNoPageNumber(t *testing.T) {
	pres := &model.Presentation{
		Meta: model.Metadata{
			Pager: false,
		},
		Slides: []model.Slide{
			{Raw: "# Hello", Index: 0},
			{Raw: "# World", Index: 1},
		},
	}
	m := NewModel(pres)
	m.width = 80
	m.height = 24
	m.searchMsg = "No match: /xyz"

	view := m.View()
	if strings.Contains(view, "1/2") {
		t.Error("pager:false with search result should not show page numbers")
	}
	if !strings.Contains(view, "No match") {
		t.Error("search result message should be visible")
	}
}

func TestPagerTrueShowsPageNumber(t *testing.T) {
	pres := &model.Presentation{
		Meta: model.Metadata{
			Title:  "Demo",
			Author: "Bob",
			Date:   "2026-06-07",
			Pager:  true,
		},
		Slides: []model.Slide{
			{Raw: "# First", Index: 0},
			{Raw: "# Second", Index: 1},
			{Raw: "# Third", Index: 2},
		},
	}
	m := NewModel(pres)
	m.width = 80
	m.height = 24

	view := m.View()
	if !strings.Contains(view, "1/3") {
		t.Error("pager:true should show page numbers '1/3'")
	}
	if !strings.Contains(view, "Bob") {
		t.Error("author should be visible when pager:true")
	}
	if !strings.Contains(view, "2026-06-07") {
		t.Error("date should be visible when pager:true")
	}
}
