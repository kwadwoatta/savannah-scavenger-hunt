package internal

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Direction int

var pEmojis = []rune{
	[]rune("🦁")[0],
	[]rune("🐊")[0],
	[]rune("🐯")[0],
	[]rune("🐻")[0],
	[]rune("🐺")[0],
}

const (
	Rest Direction = iota
	Up
	Down
	Left
	Right
)

// type sessionState uint

// const (
// 	gameView sessionState = iota
// 	triviaView
// )

type Model struct {
	// state sessionState
	tea.Model
	scavenger          Location
	scavengerDirection Direction
	Food               Food
	Viewport           Location
	ViewportPainted    bool
	Pressed            bool
	GameOver           bool
	Predators          Predators
	Keys               KeyMap
	Help               help.Model
	Score              int
}

func InitialModel() Model {
	return Model{
		scavenger:       Location{},
		Food:            Food{},
		Viewport:        Location{},
		ViewportPainted: false,
		Pressed:         false,
		GameOver:        false,
		Predators:       NewPredators(),
		Score:           0,
		Keys: KeyMap{
			Up:    key.NewBinding(key.WithKeys("up", "w"), key.WithHelp("↑/w", "up")),
			Down:  key.NewBinding(key.WithKeys("down", "s"), key.WithHelp("↓/s", "down")),
			Left:  key.NewBinding(key.WithKeys("left", "a"), key.WithHelp("←/a", "left")),
			Right: key.NewBinding(key.WithKeys("right", "d"), key.WithHelp("→/d", "right")),
			Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
		},
		Help: help.New(),
	}
}

type TickMsg time.Time

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second/8, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			fmt.Println()
			return m, tea.Quit

		case key.Matches(msg, m.Keys.Up):
			if !m.Pressed && m.scavenger.Y > 0 {
				m.scavenger.Y--
				m.scavengerDirection = Up
			}
		case key.Matches(msg, m.Keys.Down):
			if !m.Pressed && m.scavenger.Y < m.Viewport.Y {
				m.scavenger.Y++
				m.scavengerDirection = Down
			}
		case key.Matches(msg, m.Keys.Left):
			if !m.Pressed && m.scavenger.X > 0 {
				m.scavenger.X--
				m.scavengerDirection = Left
			}
		case key.Matches(msg, m.Keys.Right):
			if !m.Pressed && m.scavenger.X < m.Viewport.X {
				m.scavenger.X++
				m.scavengerDirection = Right
			}

			m.Pressed = true
		}
	case tea.WindowSizeMsg:

		if !m.ViewportPainted {
			// m.state = gameView
			m.ViewportPainted = true
			m.Help.Width = msg.Width
			m.Viewport.X = msg.Width - 5
			m.Viewport.Y = msg.Height / 2
			m.scavenger = Location{5, 5}
			m.Food = Food{&Location{m.Viewport.X / 2, m.Viewport.Y / 2}, &([]rune("🥩")[0])}
		}
	case TickMsg:
		return m.Frame()
	}
	return m, nil
}

func (m Model) Frame() (tea.Model, tea.Cmd) {
	if !m.Pressed {
		if m.scavenger.X <= 0 {
			m.scavenger.X = m.Viewport.X
		} else if m.scavenger.X >= m.Viewport.X {
			m.scavenger.X = 0
		} else if m.scavenger.Y <= 0 {
			m.scavenger.Y = m.Viewport.Y
		} else if m.scavenger.Y >= m.Viewport.Y {
			m.scavenger.Y = 0
		}

		switch m.scavengerDirection {
		case Up:
			m.scavenger.Y--
		case Down:
			m.scavenger.Y++
		case Left:
			m.scavenger.X--
		case Right:
			m.scavenger.X++
		case Rest:
		}
	}
	m.Pressed = false

	var occupiedLocations []Location

	for _, preds := range m.Predators {

		if preds.Collides(m.scavenger) {
			m.GameOver = true
			return m, tea.Quit
		}

	}

	if m.Food.Collides(m.scavenger) {
		// if m.state == gameView {
		// 	m.state = triviaView
		// } else {
		// 	m.state = gameView
		// }

		m.Score++
		m.scavengerDirection = Rest
		occupiedLocations = make([]Location, 12)
		maxX, maxY := m.Viewport.X, m.Viewport.Y
		occupiedLocations = append(occupiedLocations, Location{m.scavenger.X, m.scavenger.Y})
		newLoc := generateUniqueLocations(occupiedLocations, maxX, maxY, 1)[0]
		m.Food.Location = &newLoc
		occupiedLocations = append(occupiedLocations, newLoc)

		for i := 1; i <= 10; i++ {
			newLoc := generateUniqueLocations(occupiedLocations, maxX, maxY, 10)[0]
			predEmoji := &pEmojis[rand.Intn(len(pEmojis))]
			m.Predators.Add(NewPredator(&newLoc, predEmoji))
			occupiedLocations = append(occupiedLocations, newLoc)
		}
	}

	return m, m.tick()
}

func (m Model) View() string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("Savannah Scavenger Hunt"))
	sb.WriteByte('\n')

	viewport := make([]string, 0, m.Viewport.Y)

	for y := 0; y < m.Viewport.Y; y++ {
		var line strings.Builder

		for x := 0; x < m.Viewport.X; x++ {

			cellValue := ' '

			for _, o := range m.Predators {

				if o.Collides(Location{x, y}) {
					cellValue = *o.Emoji
					// cellValue = '@'
					break
				}
			}

			if m.Food.Collides(Location{x, y}) {
				cellValue = *m.Food.Emoji
				// cellValue = '#'
			}

			if m.scavenger.X == x && m.scavenger.Y == y {
				cellValue = []rune("🤠")[0]
				// cellValue = '*'
			}
			line.WriteRune(cellValue)
		}
		viewport = append(viewport, line.String())
	}

	// v := lipgloss.JoinHorizontal(
	// 	lipgloss.Top,
	// 	ViewportStyle.Render(strings.Join(viewport, "\n")),
	// 	ViewportStyle.Render(strings.Join(viewport, "\n")))
	// sb.WriteString(ViewportStyle.Render(v))

	sb.WriteString(ViewportStyle.Render(strings.Join(viewport, "\n")))
	sb.WriteString(fmt.Sprintf("\n%d point(s) ", m.Score))
	sb.WriteString(m.Help.View(m.Keys))

	if m.GameOver {
		sb.WriteString(GameOverStyle.Render("\n\n> Game over! <"))
	}

	return ViewStyle.Render(sb.String())
}
