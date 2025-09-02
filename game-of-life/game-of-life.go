package lifegame

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sdball/funterm/game-of-life/life"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* -- TICK -- */
type tickMsg time.Time

func tickCmd(d time.Duration) tea.Cmd {
	if d <= 0 {
		d = 50 * time.Millisecond
	}
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type styles struct {
	alive         lipgloss.Style
	selectedAlive lipgloss.Style
	selectedDead  lipgloss.Style
	status        lipgloss.Style
}

/* -- MODEL -- */

type model struct {
	width, height int
	cells         life.Set
	running       bool
	speed         time.Duration
	minSpeed      time.Duration
	maxSpeed      time.Duration
	cellWidth     int
	selected      life.Cell
	styles        styles
}

func InitialModel(width, height int) model {
	m := model{
		cells:     life.NewSet(),
		running:   false,
		speed:     100 * time.Millisecond,
		minSpeed:  1 * time.Millisecond,
		maxSpeed:  2 * time.Second,
		cellWidth: 2,
		styles: styles{
			alive:         lipgloss.NewStyle().Background(lipgloss.Color("#00AA00")),
			selectedAlive: lipgloss.NewStyle().Background(lipgloss.Color("#8cff00ff")).Bold(true),
			selectedDead:  lipgloss.NewStyle().Background(lipgloss.Color("#f6ff00ff")),
			status:        lipgloss.NewStyle().Background(lipgloss.Color("#AAAAAA")),
		},
		width:  width,
		height: height,
	}

	seed0 := life.NewSet(life.Cell{X: 21, Y: 14})
	seed1 := life.NewSet(life.Cell{X: 22, Y: 16})
	seed2 := life.NewSet(life.Cell{X: 20, Y: 15}, life.Cell{X: 21, Y: 15}, life.Cell{X: 22, Y: 15})
	for c := range seed0 {
		m.cells.Add(c)
	}
	for c := range seed1 {
		m.cells.Add(c)
	}
	for c := range seed2 {
		m.cells.Add(c)
	}
	return m
}

// tea.Model Init
func (m model) Init() tea.Cmd {
	return nil
}

// tea.Model Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, tea.EnableMouseCellMotion
	case tickMsg:
		if m.running {
			m.cells = life.Step(m.cells)
			return m, tickCmd(m.speed)
		}
		return m, nil
	case tea.MouseMsg:
		if msg.Action != tea.MouseActionPress {
			return m, nil
		}
		gridWidth, gridHeight := m.gridSize()
		x, y := msg.X, msg.Y
		if x < 0 || y < 0 || y > gridHeight || gridWidth <= 0 {
			return m, nil
		}
		cx := x / m.cellWidth
		cy := y
		if cx >= 0 && cx < gridWidth && cy >= 0 && cy < gridHeight {
			cell := life.Cell{X: cx, Y: cy}
			if m.cells.Contains(cell) {
				m.cells.Remove(cell)
			} else {
				m.cells.Add(cell)
			}
			m.selected = cell
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c", "esc":
			return m, tea.Quit
		case "c":
			if !m.running {
				m.cells = life.NewSet()
			}
			return m, nil
		case "r":
			if !m.running {
				for x := range m.width {
					for y := range m.height {
						if rand.Float32() > 0.5 {
							m.cells.Add(life.Cell{X: x, Y: y})
						}
					}
				}
				return m, nil
			}
		case " ", "enter":
			m.running = !m.running
			if m.running {
				return m, tickCmd(m.speed)
			}
			return m, nil
		case "+", "=":
			m.speed -= 50 * time.Millisecond
			if m.speed < m.minSpeed {
				m.speed = m.minSpeed
			}
			return m, nil
		case "-", "_":
			if m.speed == m.minSpeed {
				m.speed = 50 * time.Millisecond
			} else {
				m.speed += 50 * time.Millisecond
			}
			if m.speed > m.maxSpeed {
				m.speed = m.maxSpeed
			}
			return m, nil
		case "tab":
			if !m.running {
				if m.cells.Contains(m.selected) {
					m.cells.Remove(m.selected)
				} else {
					m.cells.Add(m.selected)
				}
			}
			return m, nil
		case "left", "h":
			if !m.running {
				if m.selected.X > 0 {
					m.selected.X--
				}
			}
			return m, nil
		case "right", "l":
			if !m.running {
				m.selected.X++
			}
			return m, nil
		case "up", "k":
			if !m.running {
				m.selected.Y--
			}
			return m, nil
		case "down", "j":
			if !m.running {
				m.selected.Y++
			}
			return m, nil
		case "n", "s":
			if !m.running {
				m.cells = life.Step(m.cells)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "loading...\n"
	}

	gridW, gridH := m.gridSize()
	space := strings.Repeat(" ", m.cellWidth)
	var b strings.Builder
	b.Grow((m.cellWidth*gridW + 1) * gridH)

	for y := range gridH {
		for x := range gridW {
			c := life.Cell{X: x, Y: y}
			isAlive := m.cells.Contains(c)
			isSelected := (!m.running && c == m.selected)

			switch {
			case isSelected && isAlive:
				b.WriteString(m.styles.selectedAlive.Render(space))
			case isSelected && !isAlive:
				b.WriteString(m.styles.selectedDead.Render(space))
			case isAlive:
				b.WriteString(m.styles.alive.Render(space))
			default:
				b.WriteString(space)
			}
		}
		b.WriteByte('\n')
	}

	help := m.helpLine()
	status := m.statusLine()
	b.WriteString(help)
	b.WriteByte('\n')
	b.WriteString(status)

	return b.String()
}

func (m model) gridSize() (w, h int) {
	// reserve one line for status and one line for help
	h = max(m.height-2, 0)
	m.cellWidth = max(m.cellWidth, 1)
	w = m.width / m.cellWidth
	return
}

func (m model) helpLine() string {
	var help string
	if m.running {
		help = "[q] quit  [space|enter] pause   [+/-] speed  [click] add cell"
	} else {
		help = "[q] quit  [space|enter] run  [click] toggle  [arrows] select  [tab] toggle   [c] clear  [r] random  [s/n] step"
	}
	neededPadding := m.width - len(help)
	for range neededPadding {
		help = help + " "
	}
	return m.styles.status.Render(help)
}

func (m model) statusLine() string {
	var status string
	if m.running {
		status = fmt.Sprintf("RUNNING | refresh: %+v  alive: %d", m.speed, len(m.cells))
	} else {
		status = fmt.Sprintf("PAUSED |  refresh: %+v  alive: %d", m.speed, len(m.cells))
	}

	neededPadding := m.width - len(status)
	for range neededPadding {
		status = status + " "
	}

	return m.styles.status.Render(status)
}
