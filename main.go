package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	apiURL            = "https://utilisocial.io/datacapable/v2/p/NES/map/events"
	dataFileName      = ".nes-outage-history.json"
	retentionDays     = 10
	collectionInterval = 10 * time.Minute // Store data every 10 minutes
)

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
	Timestamp time.Time `json:"timestamp"`
	NumPeople int       `json:"numPeople"`
}

type storedData struct {
	DataPoints []dataPoint `json:"dataPoints"`
}

type model struct {
	spinner       spinner.Model
	loading       bool
	err           error
	lastChecked   time.Time
	blinkOn       bool
	statusBlink   bool
	showChart     bool
	history       []dataPoint
	totalAffected int
	eventCount    int
	lastSavedTime time.Time
	dataFilePath  string
}

type tickMsg time.Time
type blinkMsg time.Time
type fetchResultMsg struct {
	totalAffected int
	eventCount    int
	err           error
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
)

func getDataFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return dataFileName
	}
	return filepath.Join(homeDir, dataFileName)
}

// roundTo10Minutes rounds a time to the nearest 10-minute interval
func roundTo10Minutes(t time.Time) time.Time {
	minutes := t.Minute()
	roundedMinutes := (minutes / 10) * 10
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		roundedMinutes,
		0,
		0,
		t.Location(),
	)
}

func loadHistoricalData(filePath string) ([]dataPoint, time.Time) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		// File doesn't exist yet, return empty
		return []dataPoint{}, time.Time{}
	}

	var stored storedData
	if err := json.Unmarshal(data, &stored); err != nil {
		// Invalid data, return empty
		return []dataPoint{}, time.Time{}
	}

	// Filter to last 10 days
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	filtered := []dataPoint{}
	
	for _, point := range stored.DataPoints {
		// Round timestamp to 10-minute interval
		rounded := roundTo10Minutes(point.Timestamp)
		
		if rounded.After(cutoff) || rounded.Equal(cutoff) {
			// Check if we already have this 10-minute interval (deduplicate)
			exists := false
			for _, existing := range filtered {
				if existing.Timestamp.Equal(rounded) {
					exists = true
					break
				}
			}
			if !exists {
				filtered = append(filtered, dataPoint{
					Timestamp: rounded,
					NumPeople: point.NumPeople,
				})
			}
		}
	}

	// Sort by timestamp
	for i := 0; i < len(filtered)-1; i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[i].Timestamp.After(filtered[j].Timestamp) {
				filtered[i], filtered[j] = filtered[j], filtered[i]
			}
		}
	}

	// Get last saved time from the last data point if available
	var lastSaved time.Time
	if len(filtered) > 0 {
		lastSaved = filtered[len(filtered)-1].Timestamp
	}
	return filtered, lastSaved
}

func saveHistoricalData(filePath string, history []dataPoint) error {
	// Filter to last 10 days
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	filtered := []dataPoint{}
	
	for _, point := range history {
		// Round to 10-minute interval
		rounded := roundTo10Minutes(point.Timestamp)
		
		if rounded.After(cutoff) || rounded.Equal(cutoff) {
			filtered = append(filtered, dataPoint{
				Timestamp: rounded,
				NumPeople: point.NumPeople,
			})
		}
	}

	// Deduplicate by 10-minute interval (keep latest value for each interval)
	intervalMap := make(map[string]dataPoint)
	for _, point := range filtered {
		key := point.Timestamp.Format("2006-01-02T15:04")
		if existing, ok := intervalMap[key]; !ok || point.Timestamp.After(existing.Timestamp) {
			intervalMap[key] = point
		}
	}

	// Convert back to slice
	deduplicated := make([]dataPoint, 0, len(intervalMap))
	for _, point := range intervalMap {
		deduplicated = append(deduplicated, point)
	}

	// Sort by timestamp
	for i := 0; i < len(deduplicated)-1; i++ {
		for j := i + 1; j < len(deduplicated); j++ {
			if deduplicated[i].Timestamp.After(deduplicated[j].Timestamp) {
				deduplicated[i], deduplicated[j] = deduplicated[j], deduplicated[i]
			}
		}
	}

	stored := storedData{
		DataPoints: deduplicated,
	}

	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	dataFilePath := getDataFilePath()
	history, lastSaved := loadHistoricalData(dataFilePath)
	
	return model{
		spinner:       s,
		loading:       true,
		blinkOn:       true,
		statusBlink:   false,
		showChart:     false,
		history:       history,
		totalAffected: 0,
		eventCount:    0,
		lastSavedTime: lastSaved,
		dataFilePath:  dataFilePath,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		blinkCmd(),
		fetchAllEvents(), // Fetch current data immediately
		tickCmd(),
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

func fetchAllEvents() tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(apiURL)
		if err != nil {
			return fetchResultMsg{0, 0, fmt.Errorf("failed to fetch: %w", err)}
		}
		defer resp.Body.Close()

		var events []OutageEvent
		if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
			return fetchResultMsg{0, 0, fmt.Errorf("failed to parse JSON: %w", err)}
		}

		// Calculate total affected customers across all outages
		totalAffected := 0
		for _, e := range events {
			totalAffected += e.NumPeople
		}

		return fetchResultMsg{totalAffected, len(events), nil}
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
		return m, fetchAllEvents()
		case "c":
			m.showChart = !m.showChart
			return m, nil
		}

	case tickMsg:
		m.loading = true
		return m, tea.Batch(fetchAllEvents(), tickCmd())

	case blinkMsg:
		m.blinkOn = !m.blinkOn
		return m, blinkCmd()

	case fetchResultMsg:
		m.loading = false
		m.lastChecked = time.Now()
		if msg.err != nil {
			m.err = msg.err
			m.totalAffected = 0
			m.eventCount = 0
		} else {
			m.err = nil
			m.totalAffected = msg.totalAffected
			m.eventCount = msg.eventCount
			
			// Round current time to 10-minute interval
			now := time.Now()
			currentInterval := roundTo10Minutes(now)
			
			// Only save if it's a new 10-minute interval
			shouldSave := currentInterval.After(m.lastSavedTime)
			
			if shouldSave {
				// Add 10-minute interval data point
				newPoint := dataPoint{
					Timestamp: currentInterval,
					NumPeople: msg.totalAffected,
				}
				
				// Check if we already have this interval (replace if exists)
				found := false
				for i, point := range m.history {
					if point.Timestamp.Equal(currentInterval) {
						m.history[i] = newPoint
						found = true
						break
					}
				}
				if !found {
					m.history = append(m.history, newPoint)
				}
				
				// Sort by timestamp
				for i := 0; i < len(m.history)-1; i++ {
					for j := i + 1; j < len(m.history); j++ {
						if m.history[i].Timestamp.After(m.history[j].Timestamp) {
							m.history[i], m.history[j] = m.history[j], m.history[i]
						}
					}
				}
				
				// Remove data older than 10 days
				cutoff := time.Now().AddDate(0, 0, -retentionDays)
				filtered := []dataPoint{}
				for _, point := range m.history {
					if point.Timestamp.After(cutoff) || point.Timestamp.Equal(cutoff) {
						filtered = append(filtered, point)
					}
				}
				m.history = filtered
				
				// Save to file
				if err := saveHistoricalData(m.dataFilePath, m.history); err == nil {
					m.lastSavedTime = currentInterval
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
			// Format timestamp - show date if it's not today
			t := displayHistory[idx].Timestamp
			now := time.Now()
			var timeStr string
			if t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day() {
				timeStr = t.Format("3 PM")
			} else {
				timeStr = t.Format("1/2 3 PM")
			}
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
		chartContent := labelStyle.Render("Total Affected Customers Over Time") + "\n"
		chartContent += valueStyle.Render("All NES Outages") + "\n\n"
		
		// Get terminal size for chart dimensions
		chartContent += renderChart(m.history, 60, 15)
		
		s += boxStyle.Render(chartContent) + "\n"
		
		if m.loading {
			s += "\n" + m.spinner.View() + " Refreshing..."
		}
	} else if m.loading && m.totalAffected == 0 && m.eventCount == 0 {
		s += m.spinner.View() + " Fetching outage data...\n"
	} else if m.err != nil {
		s += errorStyle.Render("Error: "+m.err.Error()) + "\n"
	} else {
		// Summary view
		content := ""
		content += labelStyle.Render("Total Outages: ") + valueStyle.Render(fmt.Sprintf("%d", m.eventCount)) + "\n"
		content += labelStyle.Render("Total Affected: ") + valueStyle.Render(fmt.Sprintf("%d customers", m.totalAffected)) + "\n"

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
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
