package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const apiURL = "https://utilisocial.io/datacapable/v2/p/NES/map/events"

type OutageEvent struct {
	ID              int     `json:"id"`
	StartTime       int64   `json:"startTime"`
	LastUpdatedTime int64   `json:"lastUpdatedTime"`
	Title           string  `json:"title"`
	NumPeople       int     `json:"numPeople"`
	Status          string  `json:"status"`
	Cause           string  `json:"cause"`
	Identifier      string  `json:"identifier"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

type dataPoint struct {
	Timestamp time.Time
	NumPeople int
}

type model struct {
	eventID      int
	event        *OutageEvent
	spinner      spinner.Model
	loading      bool
	err          error
	lastChecked  time.Time
	blinkOn      bool
	statusBlink  bool
	showChart    bool
	history      []dataPoint
}

type tickMsg time.Time
type blinkMsg time.Time
type fetchResultMsg struct {
	event *OutageEvent
	err   error
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("252"))

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	statusUnassigned = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("196"))

	statusAssigned = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46"))

	statusAssignedDim = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("22"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	timeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Italic(true)

	chartStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))
)

func initialModel(eventID int) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		eventID:     eventID,
		spinner:     s,
		loading:     true,
		blinkOn:     true,
		statusBlink: false,
		showChart:   false,
		history:     make([]dataPoint, 0),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchEvent(m.eventID),
		tickCmd(),
		blinkCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func blinkCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return blinkMsg(t)
	})
}

func fetchEvent(eventID int) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(apiURL)
		if err != nil {
			return fetchResultMsg{nil, fmt.Errorf("failed to fetch: %w", err)}
		}
		defer resp.Body.Close()

		var events []OutageEvent
		if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
			return fetchResultMsg{nil, fmt.Errorf("failed to parse JSON: %w", err)}
		}

		for _, e := range events {
			if e.ID == eventID {
				return fetchResultMsg{&e, nil}
			}
		}

		return fetchResultMsg{nil, fmt.Errorf("event ID %d not found", eventID)}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, fetchEvent(m.eventID)
		case "c":
			m.showChart = !m.showChart
			return m, nil
		}

	case tickMsg:
		m.loading = true
		return m, tea.Batch(fetchEvent(m.eventID), tickCmd())

	case blinkMsg:
		m.blinkOn = !m.blinkOn
		return m, blinkCmd()

	case fetchResultMsg:
		m.loading = false
		m.lastChecked = time.Now()
		if msg.err != nil {
			m.err = msg.err
			m.event = nil
		} else {
			m.err = nil
			m.event = msg.event
			m.statusBlink = (m.event.Status != "Unassigned")
			// Add data point to history
			if m.event != nil {
				m.history = append(m.history, dataPoint{
					Timestamp: time.Now(),
					NumPeople: m.event.NumPeople,
				})
				// Keep only last 50 data points
				if len(m.history) > 50 {
					m.history = m.history[len(m.history)-50:]
				}
			}
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func renderChart(history []dataPoint, width, height int) string {
	if len(history) < 2 {
		return "Not enough data points yet. Wait for a few refreshes."
	}

	// Work with a copy to avoid modifying original
	displayHistory := history
	points := len(displayHistory)
	dataWidth := width - 8 // Leave space for Y-axis labels
	if dataWidth < 1 {
		dataWidth = 1
	}
	if points > dataWidth {
		displayHistory = displayHistory[len(displayHistory)-dataWidth:]
		points = len(displayHistory)
	}

	// Find min and max values
	minPeople := displayHistory[0].NumPeople
	maxPeople := displayHistory[0].NumPeople
	for _, point := range displayHistory {
		if point.NumPeople < minPeople {
			minPeople = point.NumPeople
		}
		if point.NumPeople > maxPeople {
			maxPeople = point.NumPeople
		}
	}

	// Add some padding
	rangePeople := maxPeople - minPeople
	if rangePeople == 0 {
		rangePeople = 1
	}
	padding := rangePeople / 10
	if padding == 0 {
		padding = 1
	}
	minPeople -= padding
	maxPeople += padding
	rangePeople = maxPeople - minPeople

	// Create chart grid
	dataHeight := height - 2 // Leave space for X-axis labels
	if dataHeight < 1 {
		dataHeight = 1
	}
	chart := make([][]rune, height)
	for i := range chart {
		chart[i] = make([]rune, width)
		for j := range chart[i] {
			chart[i][j] = ' '
		}
	}

	// Draw axes
	// Y-axis (left)
	for i := 0; i < height; i++ {
		chart[i][0] = '│'
	}
	// X-axis (bottom)
	for j := 0; j < width; j++ {
		chart[height-1][j] = '─'
	}
	// Corner
	chart[height-1][0] = '└'

	// Plot data points
	for i := 0; i < points-1; i++ {
		x1 := 1 + (i * dataWidth / points)
		x2 := 1 + ((i + 1) * dataWidth / points)
		y1 := dataHeight - 1 - int(float64(displayHistory[i].NumPeople-minPeople)*float64(dataHeight-1)/float64(rangePeople))
		y2 := dataHeight - 1 - int(float64(displayHistory[i+1].NumPeople-minPeople)*float64(dataHeight-1)/float64(rangePeople))

		// Draw line using Bresenham's algorithm
		dx := x2 - x1
		dy := y2 - y1
		absDx := dx
		if absDx < 0 {
			absDx = -absDx
		}
		absDy := dy
		if absDy < 0 {
			absDy = -absDy
		}
		sx := 1
		if x1 > x2 {
			sx = -1
		}
		sy := 1
		if y1 > y2 {
			sy = -1
		}
		err := absDx - absDy

		x, y := x1, y1
		for {
			if x >= 1 && x < width-7 && y >= 0 && y < dataHeight {
				chart[y][x] = '●'
			}
			if x == x2 && y == y2 {
				break
			}
			e2 := 2 * err
			if e2 > -absDy {
				err -= absDy
				x += sx
			}
			if e2 < absDx {
				err += absDx
				y += sy
			}
		}
	}

	// Add Y-axis labels
	labelY := []int{0, dataHeight / 2, dataHeight - 1}
	for _, y := range labelY {
		value := maxPeople - int(float64(y)*float64(rangePeople)/float64(dataHeight-1))
		label := fmt.Sprintf("%d", value)
		if len(label) <= 6 {
			for i, r := range label {
				if y < height && 1+i < width {
					chart[y][1+i] = r
				}
			}
		}
	}

	// Build chart string
	var result strings.Builder
	result.WriteString("\n")
	for i := 0; i < height; i++ {
		result.WriteString(string(chart[i]))
		result.WriteString("\n")
	}

	// Add X-axis time labels
	result.WriteString("      ")
	labelCount := 3
	if points < labelCount {
		labelCount = points
	}
	if labelCount > 1 {
		for i := 0; i < labelCount; i++ {
			idx := i * (points - 1) / (labelCount - 1)
			if idx >= len(displayHistory) {
				idx = len(displayHistory) - 1
			}
			timeStr := displayHistory[idx].Timestamp.Format("3:04 PM")
			result.WriteString(fmt.Sprintf("%-12s", timeStr))
		}
	} else if points == 1 {
		timeStr := displayHistory[0].Timestamp.Format("3:04 PM")
		result.WriteString(timeStr)
	}
	result.WriteString("\n")

	return result.String()
}

func (m model) View() string {
	var s string

	s += "\n"
	s += titleStyle.Render("NES Outage Status Checker") + "\n\n"

	if m.showChart && len(m.history) > 0 {
		// Chart view
		chartContent := labelStyle.Render("Affected Customers Over Time") + "\n"
		chartContent += valueStyle.Render(fmt.Sprintf("Event ID: %d", m.eventID)) + "\n\n"
		
		// Get terminal size for chart dimensions
		chartContent += renderChart(m.history, 60, 15)
		
		s += boxStyle.Render(chartContent) + "\n"
		
		if m.loading {
			s += "\n" + m.spinner.View() + " Refreshing..."
		}
	} else if m.loading && m.event == nil {
		s += m.spinner.View() + " Fetching outage data...\n"
	} else if m.err != nil {
		s += errorStyle.Render("Error: "+m.err.Error()) + "\n"
	} else if m.event != nil {
		content := ""

		content += labelStyle.Render("Event ID: ") + valueStyle.Render(fmt.Sprintf("%d", m.event.ID)) + "\n"
		content += labelStyle.Render("Identifier: ") + valueStyle.Render(m.event.Identifier) + "\n"
		content += labelStyle.Render("Title: ") + valueStyle.Render(m.event.Title) + "\n"
		content += labelStyle.Render("Affected: ") + valueStyle.Render(fmt.Sprintf("%d people", m.event.NumPeople)) + "\n"

		if m.event.Cause != "" {
			content += labelStyle.Render("Cause: ") + valueStyle.Render(m.event.Cause) + "\n"
		}

		startTime := time.UnixMilli(m.event.StartTime)
		content += labelStyle.Render("Started: ") + valueStyle.Render(startTime.Format("Mon Jan 2, 3:04 PM")) + "\n"

		lastUpdated := time.UnixMilli(m.event.LastUpdatedTime)
		content += labelStyle.Render("Last Updated: ") + valueStyle.Render(lastUpdated.Format("Mon Jan 2, 3:04 PM")) + "\n"

		content += "\n"

		var statusDisplay string
		if m.event.Status == "Unassigned" {
			statusDisplay = statusUnassigned.Render("STATUS: UNASSIGNED")
			content += statusDisplay + "\n"
			content += valueStyle.Render("No technician assigned yet") + "\n"
		} else {
			if m.blinkOn {
				statusDisplay = statusAssigned.Render("STATUS: " + m.event.Status)
			} else {
				statusDisplay = statusAssignedDim.Render("STATUS: " + m.event.Status)
			}
			content += statusDisplay + "\n"
			content += valueStyle.Render("A technician has been assigned!") + "\n"
		}

		s += boxStyle.Render(content) + "\n"

		if m.loading {
			s += "\n" + m.spinner.View() + " Refreshing..."
		}
	}

	s += "\n"
	if !m.lastChecked.IsZero() {
		s += timeStyle.Render(fmt.Sprintf("Last checked: %s", m.lastChecked.Format("3:04:05 PM"))) + "\n"
	}
	
	helpText := "Press 'r' to refresh • 'q' to quit"
	if m.showChart {
		helpText += " • 'c' to view details"
	} else {
		helpText += " • 'c' to view chart"
	}
	helpText += " • Auto-refreshes every 30s"
	s += helpStyle.Render(helpText) + "\n"

	return s
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nes-outage-status-checker <event-id>")
		fmt.Println("Example: nes-outage-status-checker 1971637")
		os.Exit(1)
	}

	eventID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid event ID: %s\n", os.Args[1])
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(eventID))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
