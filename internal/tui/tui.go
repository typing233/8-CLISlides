package tui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/clislides/internal/executor"
	"github.com/user/clislides/internal/model"
	"github.com/user/clislides/internal/renderer"
)

type mode int

const (
	modeNormal mode = iota
	modeSearch
	modeExec
)

type Model struct {
	pres       *model.Presentation
	current    int
	width      int
	height     int
	renderer   *renderer.Renderer
	mode       mode
	searchBuf  string
	searchHits []int
	searchMsg  string
	numBuf     string
	execOutput string
	lastG      bool
}

func NewModel(pres *model.Presentation) Model {
	return Model{
		pres:    pres,
		current: 0,
		width:   80,
		height:  24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.renderer = nil
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeSearch:
			return m.handleSearch(msg)
		case modeExec:
			if msg.String() == "q" || msg.String() == "esc" || msg.String() == "enter" {
				m.mode = modeNormal
				m.execOutput = ""
			}
			return m, nil
		default:
			return m.handleNormal(msg)
		}
	}
	return m, nil
}

func (m Model) handleNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "right", "l", "n", " ", "enter":
		n := m.getNum(1)
		m.numBuf = ""
		m.lastG = false
		m.searchMsg = ""
		m.navigate(n)
		return m, nil

	case "left", "h", "p", "backspace":
		n := m.getNum(1)
		m.numBuf = ""
		m.lastG = false
		m.searchMsg = ""
		m.navigate(-n)
		return m, nil

	case "G":
		if m.numBuf != "" {
			n, _ := strconv.Atoi(m.numBuf)
			m.numBuf = ""
			m.lastG = false
			if n > 0 && n <= len(m.pres.Slides) {
				m.current = n - 1
			}
		} else {
			m.current = len(m.pres.Slides) - 1
			m.lastG = false
		}
		return m, nil

	case "g":
		if m.lastG {
			m.current = 0
			m.lastG = false
			m.numBuf = ""
		} else {
			m.lastG = true
		}
		return m, nil

	case "/":
		m.mode = modeSearch
		m.searchBuf = ""
		m.searchMsg = ""
		m.lastG = false
		m.numBuf = ""
		return m, nil

	case "e", "x":
		output := m.executeCurrentCode()
		if output != "" {
			m.mode = modeExec
			m.execOutput = output
		}
		return m, nil

	default:
		if key >= "0" && key <= "9" {
			m.numBuf += key
			m.lastG = false
		} else {
			m.numBuf = ""
			m.lastG = false
		}
		return m, nil
	}
}

func (m Model) handleSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "enter":
		m.mode = modeNormal
		m.doSearch()
		return m, nil
	case "esc":
		m.mode = modeNormal
		m.searchBuf = ""
		return m, nil
	case "backspace":
		if len(m.searchBuf) > 0 {
			m.searchBuf = m.searchBuf[:len(m.searchBuf)-1]
		}
		return m, nil
	default:
		if len(key) == 1 {
			m.searchBuf += key
		}
		return m, nil
	}
}

func (m *Model) doSearch() {
	if m.searchBuf == "" {
		m.searchMsg = ""
		return
	}
	re, err := regexp.Compile("(?i)" + m.searchBuf)
	if err != nil {
		m.searchMsg = fmt.Sprintf("Invalid regex: %s", err.Error())
		return
	}
	for i := 1; i <= len(m.pres.Slides); i++ {
		idx := (m.current + i) % len(m.pres.Slides)
		if re.MatchString(m.pres.Slides[idx].Raw) {
			m.current = idx
			m.searchMsg = fmt.Sprintf("Found: /%s (slide %d)", m.searchBuf, idx+1)
			return
		}
	}
	m.searchMsg = fmt.Sprintf("No match: /%s", m.searchBuf)
}

func (m *Model) navigate(delta int) {
	m.current += delta
	if m.current < 0 {
		m.current = 0
	}
	if m.current >= len(m.pres.Slides) {
		m.current = len(m.pres.Slides) - 1
	}
}

func (m Model) getNum(def int) int {
	if m.numBuf == "" {
		return def
	}
	n, err := strconv.Atoi(m.numBuf)
	if err != nil || n < 1 {
		return def
	}
	return n
}

func (m Model) executeCurrentCode() string {
	slide := m.pres.Slides[m.current]
	codeRe := regexp.MustCompile("(?ms)^```(\\w*)\n(.*?)^```")
	matches := codeRe.FindStringSubmatch(slide.Raw)
	if len(matches) < 3 {
		return ""
	}
	lang := matches[1]
	code := matches[2]
	out, err := executor.ExecuteCodeBlock(code, lang)
	if err != nil {
		return fmt.Sprintf("Error: %v\n%s", err, out)
	}
	return out
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if m.mode == modeExec {
		return m.renderExecView()
	}

	rendered := m.renderSlide()
	statusBar := m.renderStatus()

	contentHeight := m.height - 1
	lines := strings.Split(rendered, "\n")
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n") + "\n" + statusBar
}

func (m Model) renderSlide() string {
	if m.current >= len(m.pres.Slides) {
		return ""
	}

	slide := m.pres.Slides[m.current]
	if slide.Rendered != "" {
		return slide.Rendered
	}

	r := m.renderer
	if r == nil {
		theme := m.pres.Meta.Theme
		if theme == "" {
			theme = "dark"
		}
		var err error
		r, err = renderer.New(theme, m.width-4)
		if err != nil {
			return slide.Raw
		}
	}

	rendered, err := r.Render(slide.Raw)
	if err != nil {
		return slide.Raw
	}
	return rendered
}

func (m Model) renderStatus() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	var left string
	switch {
	case m.mode == modeSearch:
		left = fmt.Sprintf("/%s█", m.searchBuf)
	case m.searchMsg != "":
		left = m.searchMsg
	default:
		parts := []string{}
		if m.pres.Meta.Title != "" {
			parts = append(parts, m.pres.Meta.Title)
		}
		if m.pres.Meta.Author != "" {
			parts = append(parts, m.pres.Meta.Author)
		}
		if m.pres.Meta.Date != "" {
			parts = append(parts, m.pres.Meta.Date)
		}
		left = strings.Join(parts, " │ ")
	}

	right := ""
	if m.pres.Meta.Pager {
		right = fmt.Sprintf(" %d/%d ", m.current+1, len(m.pres.Slides))
	}

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 0 {
		gap = 0
	}

	bar := style.Render(left + strings.Repeat(" ", gap) + right)
	return bar
}

func (m Model) renderExecView() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		Render("─── Code Output (press q/esc/enter to close) ───")

	content := m.execOutput
	lines := strings.Split(content, "\n")
	maxLines := m.height - 3
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	return header + "\n\n" + strings.Join(lines, "\n")
}

func Run(pres *model.Presentation) error {
	m := NewModel(pres)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
