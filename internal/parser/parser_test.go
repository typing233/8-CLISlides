package parser

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	input := `---
title: "Test Pres"
author: "Tester"
theme: "light"
---

# Slide 1

Hello world

---

# Slide 2

Goodbye
`
	pres, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if pres.Meta.Title != "Test Pres" {
		t.Errorf("expected title 'Test Pres', got %q", pres.Meta.Title)
	}
	if pres.Meta.Author != "Tester" {
		t.Errorf("expected author 'Tester', got %q", pres.Meta.Author)
	}
	if pres.Meta.Theme != "light" {
		t.Errorf("expected theme 'light', got %q", pres.Meta.Theme)
	}
	if len(pres.Slides) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(pres.Slides))
	}
	if pres.Slides[0].Index != 0 {
		t.Errorf("expected slide 0 index 0, got %d", pres.Slides[0].Index)
	}
}

func TestParseNoFrontmatter(t *testing.T) {
	input := `# First

content

---

# Second

more content
`
	pres, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(pres.Slides) != 2 {
		t.Fatalf("expected 2 slides, got %d", len(pres.Slides))
	}
}

func TestParseSingleSlide(t *testing.T) {
	input := `# Only slide

Some content here
`
	pres, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(pres.Slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(pres.Slides))
	}
}
