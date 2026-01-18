package tui

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Particle represents a single firework particle
type Particle struct {
	x, y   float64
	vx, vy float64
	char   string
	color  string
	life   int
}

// Fireworks handles the fireworks animation
type Fireworks struct {
	particles []Particle
	width     int
	height    int
	active    bool
	frame     int
	maxFrames int
}

// FireworksTickMsg is sent to advance the animation
type FireworksTickMsg time.Time

// FireworksDoneMsg is sent when the animation completes
type FireworksDoneMsg struct{}

// Colors for fireworks
var fireworkColors = []string{
	"196", // Red
	"226", // Yellow
	"46",  // Green
	"51",  // Cyan
	"201", // Magenta
	"208", // Orange
	"255", // White
	"213", // Pink
}

// Particle characters
var particleChars = []string{"*", ".", "o", "+", "x", "~", "^", "'"}

// NewFireworks creates a new fireworks instance
func NewFireworks() *Fireworks {
	return &Fireworks{
		particles: make([]Particle, 0),
		maxFrames: 135, // ~4.5 seconds at 30fps
	}
}

// SetSize sets the display size
func (f *Fireworks) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// Start begins the fireworks animation
func (f *Fireworks) Start() tea.Cmd {
	f.active = true
	f.frame = 0
	f.particles = make([]Particle, 0)

	// Create initial explosions at random positions - more splash!
	for i := 0; i < 8; i++ {
		x := float64(f.width/6) + float64(rand.Intn(f.width*2/3))
		y := float64(f.height/4) + float64(rand.Intn(f.height/3))
		f.createExplosion(x, y)
	}

	return f.tick()
}

// createExplosion creates particles at a point
func (f *Fireworks) createExplosion(x, y float64) {
	color := fireworkColors[rand.Intn(len(fireworkColors))]
	numParticles := 25 + rand.Intn(20) // More particles!

	for i := 0; i < numParticles; i++ {
		speed := 0.4 + rand.Float64()*1.2
		f.particles = append(f.particles, Particle{
			x:     x,
			y:     y,
			vx:    speed * 3.0 * (rand.Float64() - 0.5), // Wider spread
			vy:    speed * 1.8 * (rand.Float64() - 0.6),
			char:  particleChars[rand.Intn(len(particleChars))],
			color: color,
			life:  20 + rand.Intn(25), // Longer life
		})
	}
}

// tick returns a command that sends a tick message
func (f *Fireworks) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(t time.Time) tea.Msg {
		return FireworksTickMsg(t)
	})
}

// Update handles animation updates
func (f *Fireworks) Update(msg tea.Msg) (*Fireworks, tea.Cmd) {
	switch msg.(type) {
	case FireworksTickMsg:
		if !f.active {
			return f, nil
		}

		f.frame++

		// Spawn new explosions frequently for more splash
		if f.frame%8 == 0 && f.frame < f.maxFrames-20 {
			// Spawn 1-2 explosions at a time
			numExplosions := 1 + rand.Intn(2)
			for i := 0; i < numExplosions; i++ {
				x := float64(f.width/6) + float64(rand.Intn(f.width*2/3))
				y := float64(f.height/4) + float64(rand.Intn(f.height/3))
				f.createExplosion(x, y)
			}
		}

		// Update particles
		newParticles := make([]Particle, 0)
		for _, p := range f.particles {
			p.x += p.vx
			p.y += p.vy
			p.vy += 0.06 // Gravity
			p.life--

			if p.life > 0 && p.y < float64(f.height)-1 && p.y > 0 &&
				p.x > 0 && p.x < float64(f.width)-1 {
				newParticles = append(newParticles, p)
			}
		}
		f.particles = newParticles

		// Check if animation is done
		if f.frame >= f.maxFrames {
			f.active = false
			return f, func() tea.Msg { return FireworksDoneMsg{} }
		}

		return f, f.tick()
	}

	return f, nil
}

// View renders the fireworks as a full-screen display
func (f *Fireworks) View() string {
	if !f.active || f.width == 0 || f.height == 0 {
		return ""
	}

	// Create styled grid
	grid := make([][]string, f.height)
	for i := range grid {
		grid[i] = make([]string, f.width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Place particles on grid with colors
	for _, p := range f.particles {
		px := int(p.x)
		py := int(p.y)
		if py >= 0 && py < f.height && px >= 0 && px < f.width {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(p.color))
			if p.life < 8 {
				style = style.Faint(true)
			} else {
				style = style.Bold(true)
			}
			grid[py][px] = style.Render(p.char)
		}
	}

	// Build output
	contentStyle := lipgloss.NewStyle().
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		Background(lipgloss.Color("0"))

	var content strings.Builder

	// Add celebration message at top
	msg := " OBJECTIVE COMPLETE! "
	padding := (f.width - len(msg)) / 2
	if padding > 0 {
		content.WriteString(strings.Repeat(" ", padding))
	}
	content.WriteString(titleStyle.Render(msg))
	content.WriteString("\n\n")

	// Add fireworks
	for _, row := range grid[2:] { // Skip first 2 rows for title
		for _, cell := range row {
			content.WriteString(cell)
		}
		content.WriteString("\n")
	}

	return lipgloss.Place(
		f.width,
		f.height,
		lipgloss.Center,
		lipgloss.Center,
		contentStyle.Render(content.String()),
	)
}

// IsActive returns whether fireworks are currently displaying
func (f *Fireworks) IsActive() bool {
	return f.active
}
